package main

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	fileCounter uint64 = 0
	tileCounter uint64 = 0

	fileCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "file_counter",
			Help: "Counts number of file processed",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&fileCounter))
		})

	tileCounterP = promauto.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "tile_counter",
			Help: "Counts number of tiles generated",
		},
		func() float64 {
			return float64(atomic.LoadUint64(&tileCounter))
		})

	versionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "insided",
		Name:      "version",
		Help:      "App version.",
	}, []string{"version"})
)
