package outbox_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Глобальные метрики
var (
	EventsSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "booking_events_sent_total",
		Help: "Total number of events sent to Kafka",
	})

	EventsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "booking_events_failed_total",
		Help: "Total number of events that failed to send",
	})
)

