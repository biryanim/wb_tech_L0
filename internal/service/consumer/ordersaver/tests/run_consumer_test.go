package tests

import (
	"context"
	"errors"
	"github.com/biryanim/wb_tech_L0/internal/client/cache"
	cacheMock "github.com/biryanim/wb_tech_L0/internal/client/cache/mocks"
	"github.com/biryanim/wb_tech_L0/internal/client/db"
	txMock "github.com/biryanim/wb_tech_L0/internal/client/db/mocks"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	kafkaMock "github.com/biryanim/wb_tech_L0/internal/client/kafka/mocks"
	"github.com/biryanim/wb_tech_L0/internal/repository"
	repoMock "github.com/biryanim/wb_tech_L0/internal/repository/mocks"
	"github.com/biryanim/wb_tech_L0/internal/service/consumer/ordersaver"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_RunConsumer(t *testing.T) {
	type orderRepositoryMockFunc func(mc *minimock.Controller) repository.OrderRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type cacheMockFunc func(mc *minimock.Controller) cache.Client
	type producerKafkaMockFunc func(mc *minimock.Controller) kafka.Producer
	type consumerKafkaMockFunc func(mc *minimock.Controller) kafka.Consumer
	tests := []struct {
		name                string
		timeout             time.Duration
		expectErr           bool
		orderRepositoryMock orderRepositoryMockFunc
		txManagerMock       txManagerMockFunc
		consumerKafkaMock   consumerKafkaMockFunc
		producerKafkaMock   producerKafkaMockFunc
		cacheMock           cacheMockFunc
	}{
		{
			name:      "consumer runs until context cancelled",
			timeout:   50 * time.Millisecond,
			expectErr: true,
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				mock := kafkaMock.NewConsumerMock(mc)
				mock.ConsumeMock.Set(func(ctx context.Context, topicName string, handler kafka.Handler) (err error) {
					<-ctx.Done()
					return ctx.Err()
				})
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)

				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				return kafkaMock.NewProducerMock(mc)
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				return mock
			},
		},
		{
			name:      "consumer returns error",
			timeout:   1 * time.Second,
			expectErr: true,
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				mock := kafkaMock.NewConsumerMock(mc)
				mock.ConsumeMock.Return(errors.New("consumer error"))
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)

				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				return kafkaMock.NewProducerMock(mc)
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)

			service := ordersaver.NewService(
				tt.orderRepositoryMock(mc),
				tt.consumerKafkaMock(mc),
				tt.txManagerMock(mc),
				tt.cacheMock(mc),
				tt.producerKafkaMock(mc),
			)

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			err := service.RunConsumer(ctx)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
