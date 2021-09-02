package main

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	fileCounter           uint64 = 0
	tileCounter           uint64 = 0
	errCounter            uint64 = 0
	rejectedBlurryCounter uint64 = 0

	//nolint deadcode
	fileCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "file_counter",
			Help: "Counts number of file processed",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&fileCounter))
		})

	//nolint deadcode
	tileCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "tile_counter",
			Help: "Counts number of tiles generated",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&tileCounter))
		})

	//nolint deadcode
	rejectedBlurryCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "rejected_blurry_counter",
			Help: "Counts number of rejected pictures",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&rejectedBlurryCounter))
		})

	//nolint deadcode
	errCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "error_counter",
			Help: "Counts processing errors",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&errCounter))
		})

	versionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "insided",
		Name:      "version",
		Help:      "App version.",
	}, []string{"version"})
)
