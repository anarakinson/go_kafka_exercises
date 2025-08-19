package producer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_kafka_exercises/01.producer_consumer/producer/internal/repository/data_mock"
	"github.com/anarakinson/go_kafka_exercises/01.producer_consumer/producer/pkg/metrics"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

var ErrCanceled = errors.New("context canceled")

type Producer struct {
	brokers []string
	configs *sarama.Config
}

func NewProducer(brokers []string, configs *sarama.Config) *Producer {
	return &Producer{
		brokers: brokers,
		configs: configs,
	}
}

// создание топика в кафке,
func (p *Producer) CreateTopic(topicName string) {
	admin, err := sarama.NewClusterAdmin(
		p.brokers, // адреса узлов
		sarama.NewConfig(),
	)
	if err != nil {
		logger.Log.Error("Error creating sarama admin client", zap.Error(err))
		return
	}
	defer admin.Close()

	// получаем список существующих топиков
	topics, err := admin.ListTopics()
	if err != nil {
		logger.Log.Error("Error listing topics", zap.Error(err))
		return
	}

	// проверяем существование топика
	_, exists := topics[topicName]
	if exists {
		logger.Log.Info("Topic already exists", zap.String("topic name", topicName))
		return
	}

	// если не существует
	// создаем топик с 2 партициями и фактором репликации 1
	err = admin.CreateTopic(
		topicName,
		&sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		},
		false,
	)
	if err != nil {
		logger.Log.Error("Error creating topic", zap.String("topic name", topicName), zap.Error(err))
		return
	}
}

func (p *Producer) Produce(ctx context.Context, topicName string, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	mock := data_mock.NewDataMock()
	// Создаем синхронный продюсер
	producer, err := sarama.NewSyncProducer(p.brokers, p.configs)
	if err != nil {
		logger.Log.Error("Failed to start producer", zap.Error(err))
		return fmt.Errorf("error making sarama producer %w", err)
	}
	defer producer.Close()

	for {
		select {
		case <-ctx.Done():
			return ErrCanceled
		case <-ticker.C:
			// генерируем пользовательскую активность
			event := mock.GetRandomAction()
			logger.Log.Info(
				"New user action created",
				zap.String("userID", event.UserID),
				zap.String("productID", event.ProductID),
				zap.String("eventType", string(event.EventType)),
			)

			// сериализуем в json
			eventJson, err := json.Marshal(event)
			if err != nil {
				logger.Log.Error("Error serializing eveng", zap.Error(err))
				// обновление метрики
				metrics.ProduceErrors.Inc()
				continue
			}

			// начало замера времени
			start := time.Now()

			// отправка события в кафку
			partition, offset, err := producer.SendMessage(
				&sarama.ProducerMessage{
					Topic: topicName,
					Key:   sarama.StringEncoder(event.UserID),
					Value: sarama.ByteEncoder(eventJson),
				},
			)
			if err != nil {
				logger.Log.Error("error sending message to kafka", zap.Error(err))
				// обновление метрики
				metrics.ProduceErrors.Inc()
			} else {
				// обновление метрики
				metrics.ProduceDuration.Observe(float64(time.Since(start).Seconds()))
				logger.Log.Info(
					"Event succesfully sent",
					zap.Int32("Partition", partition),
					zap.Int64("Offset", offset),
					zap.String("UserID", event.UserID),
					zap.String("EventType", string(event.EventType)),
				)
				// обновление метрики
				metrics.EventsProduced.WithLabelValues(string(event.EventType)).Inc()
			}
		}
	}
}
