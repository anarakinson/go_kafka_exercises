package consumer

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_kafka_exercises/01.producer_consumer/consumer/pkg/metrics"
	"github.com/anarakinson/go_stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

// consumerHandler реализует интерфейс sarama.ConsumerGroupHandler
type consumerHandler struct {
	logFile string
}

// Setup вызывается перед началом потребления
func (h *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	logger.Log.Info("Consumer group setup", zap.String("member_ID", session.MemberID()))
	return nil
}

// Cleanup вызывается после завершения потребления
func (h *consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	logger.Log.Info("Consumer group cleanup", zap.Any("offsets", session.Claims()))
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиции
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// открываем файл, в который будут сохраняться сообщения
	// os.O_APPEND - запись в конец файла (не перезаписывать существующее)
	// os.O_CREATE - создать файл, если не существует
	// os.O_WRONLY - только запись (без чтения)
	file, err := os.OpenFile(h.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Error("Error opening file", zap.Error(err))
		return err
	}
	defer file.Close() // закрываем файл в конце

	// создаем писателя для буферизованной записи
	writer := bufio.NewWriter(file)
	defer writer.Flush() // сбрасываем буфер в конце

	// обрабатываем каждое сообщение из партиции
	for message := range claim.Messages() {
		// засекаем время
		start := time.Now()

		// валидируем json
		if !json.Valid(message.Value) {
			logger.Log.Error("Invalid json received", zap.Error(err))
			// обновляем метрики
			metrics.ConsumeErrors.Inc()
			continue
		}

		// записываем дынные в файл
		if _, err := writer.Write(append(message.Value, '\n')); err != nil {
			logger.Log.Error("File write error", zap.Error(err))
			metrics.ConsumeErrors.Inc()
			continue
		}

		// сбрасываем буфер
		err = writer.Flush()
		if err != nil {
			logger.Log.Error("Error flushing buffer", zap.Error(err))
			// обновляем метрики
			metrics.ConsumeErrors.Inc()
			continue
		}

		// подтверждаем обработку сообщения (коммит оффсета)
		session.MarkMessage(message, "")

		// обновляем метрики
		metrics.ProcessingTime.Observe(time.Since(start).Seconds())

		// логируем факт записи
		logger.Log.Info(
			"Event successfully written to file",
			zap.Int32("Partition", message.Partition),
			zap.Int64("Offset", message.Offset),
		)
	}

	return nil
}
