package main

import (
	"context"
	"fmt"
	"io/fs"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

var (
	version = "no version from LDFLAGS"

	source = flag.String("source", "", "Source directory for the images")
	dest   = flag.String("dest", "", "Destination directory for the images")
	x      = flag.Int("x", 224, "Size of the target tiles x")
	y      = flag.Int("y", 224, "Size of the target tiles y")

	resize      = flag.Int("resize", 0, "Resize percentage before tilling")
	workerCount = flag.Int("workerCount", 4, "Parallel worker count")

	logLevel        = flag.String("logLevel", "INFO", "DEBUG|INFO|WARN|ERROR")
	httpMetricsPort = flag.Int("httpMetricsPort", 34130, "http port")

	httpMetricsServer *http.Server
)

func main() {
	flag.Parse()

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "caller", log.Caller(5), "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "app", "ml-image-tile")
	logger = NewLevelFilterFromString(logger, *logLevel)

	stdlog.SetOutput(log.NewStdlibAdapter(logger))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// catch termination
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	// web server metrics
	g.Go(func() error {
		httpMetricsServer = &http.Server{
			Addr:         fmt.Sprintf(":%d", *httpMetricsPort),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		level.Info(logger).Log("msg", fmt.Sprintf("HTTP Metrics server listening at :%d", *httpMetricsPort))

		versionGauge.WithLabelValues(version).Add(1)

		// Register Prometheus metrics handler.
		http.Handle("/metrics", promhttp.Handler())

		if err := httpMetricsServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	queue := make(chan string, *workerCount)

	g.Go(func() error {
		var wg sync.WaitGroup

		worker := func(jobs <-chan string) {
			for {
				path, ok := <-jobs
				if !ok {
					break
				}
				level.Debug(logger).Log("msg", "processing", "path", path)

				err := processImageBimg(logger, path, *dest, *resize, *x, *y)
				if err != nil {
					level.Error(logger).Log("msg", "error processing", "path", path, "err", err)
					continue
				}
			}
			level.Debug(logger).Log("msg", "stopping worker")
			wg.Done()
		}

		for w := 0; w < *workerCount; w++ {
			wg.Add(1)

			go worker(queue)
		}

		wg.Wait()

		return fmt.Errorf("finished work")
	})

	g.Go(func() error {
		err := filepath.Walk(*source, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}

			if !info.IsDir() {
				level.Debug(logger).Log("msg", "queuing", "path", path)
				queue <- path
			}

			return nil
		})

		close(queue)
		return err
	})

	select {
	case <-interrupt:
		cancel()

		break
	case <-ctx.Done():
		break
	}

	level.Warn(logger).Log("msg", "received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if httpMetricsServer != nil {
		_ = httpMetricsServer.Shutdown(shutdownCtx)
	}

	err := g.Wait()
	if err != nil {
		level.Error(logger).Log("msg", "stopping", "error", err)
		os.Exit(2)
	}
}
