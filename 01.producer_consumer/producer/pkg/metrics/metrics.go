package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Счетчик отправленных событий с разделением по типам
	EventsProduced = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_producer_events_total",
		Help: "Total produced events by type",
	}, []string{"event_type"}) // группировка по event_type

	// Счетчик ошибок при отправке
	ProduceErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kafka_producer_errors_total",
		Help: "Total produce errors",
	})

	// Гистограмма времени отправки сообщений
	ProduceDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "kafka_produce_duration_seconds",
		Help:    "Produce operation duration",
		Buckets: prometheus.DefBuckets, // Стандартные корзины для гистограммы
	})
)
