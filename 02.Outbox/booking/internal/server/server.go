package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/app/delivery"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/pkg/metrics"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/pkg/middleware"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

var ErrNilServer = errors.New("server not initialized")

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(port string, handler *delivery.Handler) error {
	// роутинг
	mux := http.NewServeMux()
	mux.HandleFunc("/bookings", handler.BookingHandler)

	// Обертываем в middleware
	var middlewareHandler http.Handler = mux

	// добавляем middleware для rate limiting
	rateLimiter := middleware.NewRateLimiter(100, 100)
	middlewareHandler = rateLimiter.RateLimiting(middlewareHandler)

	// Добавляем middleware для сбора метрик
	middlewareHandler = metrics.MetricsMiddleware(middlewareHandler)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: middlewareHandler,
	}

	// запускаем сервер
	logger.Log.Info("Starting server", zap.String("port", port))
	return s.server.ListenAndServe()

}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Log.Info("Shutting down server")
	if s.server == nil {
		return ErrNilServer
	}
	return s.server.Shutdown(ctx)
}
