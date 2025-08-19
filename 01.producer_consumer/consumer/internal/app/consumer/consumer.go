package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

var ErrCanceled = errors.New("context canceled")

// структура для работы с кнешним миром
type Consumer struct {
	brokers []string
	config  *sarama.Config
}

// NewConsumer создает новый экземпляр Consumer
func NewConsumer(brokers []string, config *sarama.Config) *Consumer {
	return &Consumer{
		brokers: brokers, // преобразуем строку в слайс
		config:  config,
	}
}

// RunConsumer запускает потребителя для указанной темы и группы
func (c *Consumer) RunConsumer(ctx context.Context, topics []string, groupName string, logFile string) error {
	// создаем нового потребителя
	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, groupName, c.config)
	if err != nil {
		logger.Log.Error("Error creating consumer group", zap.Error(err))
		return err
	}
	// закрываем соединение при завершении
	defer consumerGroup.Close()

	// создаем обработчик
	handler := &consumerHandler{logFile: logFile}

	// запускаем цикл потребления
	logger.Log.Info("Starting consuming circle")
	for {
		select {
		case <-ctx.Done():
			// Контекст отменен, завершаем работу
			logger.Log.Info("Context canceled, stopping consumer")
			return ErrCanceled
		default:
			// Потребляем сообщения
			err := consumerGroup.Consume(ctx, topics, handler)
			if err != nil {
				logger.Log.Error("Error consuming circle", zap.Error(err))
				return err
			}
		}
	}

}
