package infrastructure

import (
	"github.com/IBM/sarama"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/domain"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

type Producer struct {
	brokers  []string
	producer sarama.SyncProducer
}

func NewProducer(brokers []string, config *sarama.Config) (*Producer, error) {
	kProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Log.Error("Error connecting to Kafka", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Successfully connected to Kafka")

	return &Producer{
		brokers:  brokers,
		producer: kProducer,
	}, nil
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

func (p *Producer) SendToKafka(topic string, event domain.OutboxEvent) error {
	// создаем сообщение
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.EventID),
		Value: sarama.StringEncoder(event.Payload),
	}

	// повторяем попытку трижды
	for attempt := 1; attempt <= 3; attempt++ {
		_, _, err := p.producer.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Error sending message to Kafka", zap.Int("attempt", attempt), zap.Error(err))
			if attempt == 3 {
				return err
			}
			continue
		}
		// ели успешно - прерываем цикл
		break
	}

	logger.Log.Info("Message sended to Kafka successfully")
	return nil
}
