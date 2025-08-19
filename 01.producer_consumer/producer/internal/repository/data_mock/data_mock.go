package data_mock

import (
	"math/rand"
	"time"

	"github.com/anarakinson/go_kafka_exercises/01.producer_consumer/producer/internal/domain"
	"github.com/google/uuid"
)

type DataMock struct {
	users      []string
	products   []string
	eventTypes []domain.UserEventType
}

func NewDataMock() *DataMock {
	return &DataMock{
		users: []string{
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
		},
		products:   []string{"product1", "product2", "product3", "product4", "product5", "product6"},
		eventTypes: []domain.UserEventType{domain.EventView, domain.EventAddToCart, domain.EventPurchase},
	}
}

func (m *DataMock) GetRandomAction() domain.UserAction {
	return domain.UserAction{
		UserID:    m.GetRandomUser(), // выбираем случайного пользователя из мока
		ProductID: m.GetRandomProduct(),
		EventType: m.GetRandomEventType(),
		Timestamp: time.Now().Unix(),
	}
}

func (m *DataMock) GetRandomUser() string {
	return m.users[rand.Intn(len(m.users))]
}

func (m *DataMock) GetRandomProduct() string {
	return m.products[rand.Intn(len(m.products))]
}

func (m *DataMock) GetRandomEventType() domain.UserEventType {
	return m.eventTypes[rand.Intn(len(m.eventTypes))]
}
