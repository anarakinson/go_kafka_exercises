package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Счетчик потребленных событий по типам
	EventsConsumed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_consumer_events_total",
		Help: "Total consumed events by type",
	}, []string{"event_type"})

	// Счетчик ошибок потребления
	ConsumeErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kafka_consumer_errors_total",
		Help: "Total consume errors",
	})

	// Гистограмма времени обработки событий
	ProcessingTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "kafka_processing_duration_seconds",
		Help:    "Event processing duration",
		Buckets: prometheus.DefBuckets,
	})
)
