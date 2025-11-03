package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/cache"
	cacheMock "github.com/biryanim/wb_tech_L0/internal/client/cache/mocks"
	"github.com/biryanim/wb_tech_L0/internal/client/db"
	txMock "github.com/biryanim/wb_tech_L0/internal/client/db/mocks"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	kafkaMock "github.com/biryanim/wb_tech_L0/internal/client/kafka/mocks"
	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/biryanim/wb_tech_L0/internal/repository"
	repoMock "github.com/biryanim/wb_tech_L0/internal/repository/mocks"
	"github.com/biryanim/wb_tech_L0/internal/service/consumer/ordersaver"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_OrderSaveHandler(t *testing.T) {
	type orderRepositoryMockFunc func(mc *minimock.Controller) repository.OrderRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type cacheMockFunc func(mc *minimock.Controller) cache.Client
	type consumerKafkaMockFunc func(mc *minimock.Controller) kafka.Consumer
	type producerKafkaMockFunc func(mc *minimock.Controller) kafka.Producer

	type args struct {
		ctx     context.Context
		message *sarama.ConsumerMessage
	}

	dateCreated, _ := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")

	var (
		ctx = context.Background()

		orderUID = gofakeit.UUID()

		message = fmt.Sprintf(`{
   			"order_uid": "%s",
   			"track_number": "WBILMTESTTRACK",
   			"entry": "WBIL",
   			"delivery": {
      			"name": "Test Testov",
      			"phone": "9720000000",
      			"zip": "2639809",
      			"city": "Kiryat Mozkin",
      			"address": "Ploshad Mira 15",
      			"region": "Kraiot",
      			"email": "test@gmail.com"
   			},
   			"payment": {
				"transaction": "%s",
      			"request_id": "",
      			"currency": "USD",
      			"provider": "wbpay",
      			"amount": 1817,
      			"payment_dt": 1637907727,
      			"bank": "alpha",
      			"delivery_cost": 1500,
      			"goods_total": 317,
      			"custom_fee": 0
   			},
   			"items": [
      		{
         		"chrt_id": 9934930,
         		"track_number": "WBILMTESTTRACK",
         		"price": 453,
         		"rid": "ab4219087a764ae0btest",
         		"name": "Mascaras",
         		"sale": 30,
         		"size": "0",
         		"total_price": 317,
         		"nm_id": 2389212,
         		"brand": "Vivienne Sabo",
         		"status": 202
      		}
   			],
   			"locale": "en",
   			"internal_signature": "",
   			"customer_id": "test",
   			"delivery_service": "meest",
   			"shardkey": "9",
   			"sm_id": 99,
   			"date_created": "2021-11-26T06:22:19Z",
   			"oof_shard": "1"
			}`, orderUID, orderUID,
		)

		delivery = &model.Delivery{
			Name:    "Test Testov",
			Phone:   "9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		}
		payment = &model.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		}

		items = []*model.Item{
			&model.Item{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		}

		order = &model.Order{
			OrderUID:    orderUID,
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery:    *delivery,
			Payment:     *payment,
			Items: []model.Item{
				*items[0],
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			ShardKey:          "9",
			SmID:              99,
			DateCreated:       dateCreated,
			OofShard:          "1",
		}

		invalidJSON  = `{invalid json`
		invalidOrder = `{
          "order_uid": "",
          "track_number": "",
          "entry": "WBIL",
          "delivery": {},
          "payment": {},
          "items": []
        }`
	)

	tests := []struct {
		name                string
		args                args
		err                 error
		orderRepositoryMock orderRepositoryMockFunc
		txManagerMock       txManagerMockFunc
		consumerKafkaMock   consumerKafkaMockFunc
		producerKafkaMock   producerKafkaMockFunc
		cacheMock           cacheMockFunc
	}{
		{
			name: "success - order created",
			args: args{
				ctx:     ctx,
				message: &sarama.ConsumerMessage{Value: []byte(message)},
			},
			err: nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, order).Return(order.OrderUID, nil)
				mock.CreateDeliveryMock.Expect(ctx, order.OrderUID, delivery).Return(order.OrderUID, nil)
				mock.CreatePaymentMock.Expect(ctx, order.OrderUID, payment).Return(order.OrderUID, nil)
				mock.CreateItemMock.Expect(ctx, order.OrderUID, items[0]).Return(nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				return kafkaMock.NewProducerMock(mc)
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.SetMock.Expect(order.OrderUID, order).Return(true)
				return mock
			},
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				return kafkaMock.NewConsumerMock(mc)
			},
		},
		{
			name: "failure - invalid json",
			args: args{
				ctx:     ctx,
				message: &sarama.ConsumerMessage{Value: []byte(invalidJSON)},
			},
			err: nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				mock := kafkaMock.NewProducerMock(mc)
				mock.SendToDLQMock.Return(nil)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				return mock
			},
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				return kafkaMock.NewConsumerMock(mc)
			},
		},
		{
			name: "error - validation failed",
			args: args{
				ctx:     ctx,
				message: &sarama.ConsumerMessage{Value: []byte(invalidOrder)},
			},
			err: nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				mock := kafkaMock.NewProducerMock(mc)
				mock.SendToDLQMock.Return(nil)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				return mock
			},
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				return kafkaMock.NewConsumerMock(mc)
			},
		},
		{
			name: "error - DLQ send failed",
			args: args{
				ctx:     ctx,
				message: &sarama.ConsumerMessage{Value: []byte(invalidJSON)},
			},
			err: errors.New("dlq error"),
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				mock := kafkaMock.NewProducerMock(mc)
				mock.SendToDLQMock.Return(errors.New("dlq error"))
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				return mock
			},
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				return kafkaMock.NewConsumerMock(mc)
			},
		},
		{
			name: "error - transaction rollback",
			args: args{
				ctx:     ctx,
				message: &sarama.ConsumerMessage{Value: []byte(message)},
			},
			err: errors.New("transaction error"),
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.CreateOrderMock.Expect(ctx, order).Return(order.OrderUID, nil)
				mock.CreateDeliveryMock.Expect(ctx, order.OrderUID, delivery).Return(order.OrderUID, errors.New("db error"))
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
					_ = f(ctx)
					return errors.New("transaction error")
				})
				return mock
			},
			producerKafkaMock: func(mc *minimock.Controller) kafka.Producer {
				mock := kafkaMock.NewProducerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				return mock
			},
			consumerKafkaMock: func(mc *minimock.Controller) kafka.Consumer {
				return kafkaMock.NewConsumerMock(mc)
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

			err := service.OrderSaveHandler(tt.args.ctx, tt.args.message)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
