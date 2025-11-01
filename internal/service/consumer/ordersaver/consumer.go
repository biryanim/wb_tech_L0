package ordersaver

import (
	"context"

	"github.com/biryanim/wb_tech_L0/internal/client/cache"
	"github.com/biryanim/wb_tech_L0/internal/client/db"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	"github.com/biryanim/wb_tech_L0/internal/repository"
	def "github.com/biryanim/wb_tech_L0/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	consumer        kafka.Consumer
	producer        kafka.Producer
	txManager       db.TxManager
	cache           cache.Client
}

// NewService creates and returns a new order saver service with the provided dependencies.
func NewService(orderRepository repository.OrderRepository, consumer kafka.Consumer, txManager db.TxManager, cache cache.Client, producer kafka.Producer) *service {
	return &service{
		orderRepository: orderRepository,
		consumer:        consumer,
		txManager:       txManager,
		cache:           cache,
		producer:        producer,
	}
}

// RunConsumer starts the consumer loop that processes Kafka messages until the context is cancelled.
func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		errCh <- s.consumer.Consume(ctx, "order-topic", s.OrderSaveHandler)
	}()

	return errCh
}
