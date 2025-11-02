package tests

import (
	"context"
	"github.com/biryanim/wb_tech_L0/internal/client/cache"
	cacheMock "github.com/biryanim/wb_tech_L0/internal/client/cache/mocks"
	"github.com/biryanim/wb_tech_L0/internal/client/db"
	txMock "github.com/biryanim/wb_tech_L0/internal/client/db/mocks"
	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/biryanim/wb_tech_L0/internal/repository"
	repoMock "github.com/biryanim/wb_tech_L0/internal/repository/mocks"
	"github.com/biryanim/wb_tech_L0/internal/service/order"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"time"

	"testing"
)

func Test_RestoreCache(t *testing.T) {
	type orderRepositoryMockFunc func(mc *minimock.Controller) repository.OrderRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type cacheMockFunc func(mc *minimock.Controller) cache.Client

	type args struct {
		ctx   context.Context
		limit int
	}

	var (
		ctx = context.Background()

		orderUID1       = gofakeit.UUID()
		orderUID2       = gofakeit.UUID()
		trackNumber     = gofakeit.LetterN(10)
		entry           = gofakeit.LetterN(4)
		locale          = gofakeit.LanguageAbbreviation()
		customerID      = gofakeit.Username()
		shardKey        = gofakeit.DigitN(1)
		smID            = gofakeit.IntRange(1, 100)
		oofShard        = gofakeit.DigitN(1)
		deliveryService = gofakeit.Company()
		dateCreated     = gofakeit.Date()

		//repoError = errors.New("repository error")

		delivery = &model.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		}

		payment = &model.Payment{
			Transaction:  gofakeit.LetterN(19),
			RequestID:    gofakeit.UUID(),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.Company(),
			Amount:       gofakeit.Float64Range(10, 100000),
			PaymentDt:    time.Now().Unix(),
			Bank:         gofakeit.BankName(),
			DeliveryCost: gofakeit.Float64Range(100, 5000),
			GoodsTotal:   gofakeit.IntRange(10, 5000),
			CustomFee:    gofakeit.Float64Range(0, 1000),
		}

		item = model.Item{
			ChrtID:      int64(gofakeit.IntRange(1000000, 9999999)),
			TrackNumber: trackNumber,
			Price:       gofakeit.Price(100, 10000),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.IntRange(0, 100),
			Size:        gofakeit.DigitN(1),
			TotalPrice:  gofakeit.Float64Range(100, 10000),
			NmID:        int64(gofakeit.IntRange(1000000, 9999999)),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.IntRange(200, 300),
		}

		items = []*model.Item{
			&item,
		}

		res1 = &model.Order{
			OrderUID:          orderUID1,
			TrackNumber:       trackNumber,
			Entry:             entry,
			Delivery:          *delivery,
			Payment:           *payment,
			Items:             []model.Item{item},
			Locale:            locale,
			InternalSignature: "",
			CustomerID:        customerID,
			DeliveryService:   deliveryService,
			ShardKey:          shardKey,
			SmID:              smID,
			DateCreated:       dateCreated,
			OofShard:          oofShard,
		}

		res2 = &model.Order{
			OrderUID:          orderUID2,
			TrackNumber:       trackNumber,
			Entry:             entry,
			Delivery:          *delivery,
			Payment:           *payment,
			Items:             []model.Item{item},
			Locale:            locale,
			InternalSignature: "",
			CustomerID:        customerID,
			DeliveryService:   deliveryService,
			ShardKey:          shardKey,
			SmID:              smID,
			DateCreated:       dateCreated.AddDate(0, 0, 1),
			OofShard:          oofShard,
		}

		repoOrderModel1 = &model.Order{
			OrderUID:          orderUID1,
			TrackNumber:       trackNumber,
			Entry:             entry,
			Locale:            locale,
			InternalSignature: "",
			CustomerID:        customerID,
			DeliveryService:   res1.DeliveryService,
			ShardKey:          shardKey,
			SmID:              smID,
			DateCreated:       res1.DateCreated,
			OofShard:          oofShard,
		}

		repoOrderModel2 = &model.Order{
			OrderUID:          orderUID2,
			TrackNumber:       trackNumber,
			Entry:             entry,
			Locale:            locale,
			InternalSignature: "",
			CustomerID:        customerID,
			DeliveryService:   res2.DeliveryService,
			ShardKey:          shardKey,
			SmID:              smID,
			DateCreated:       res2.DateCreated,
			OofShard:          oofShard,
		}
	)

	tests := []struct {
		name                string
		args                args
		err                 error
		orderRepositoryMock orderRepositoryMockFunc
		txManagerMock       txManagerMockFunc
		cacheMock           cacheMockFunc
	}{
		{
			name: "success - restore cache with single order",
			args: args{
				ctx:   ctx,
				limit: 1,
			},
			err: nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.ListOrdersByLastAddedMock.Expect(ctx, 1).Return([]*model.Order{repoOrderModel1}, nil)
				mock.GetDeliveryMock.Expect(ctx, repoOrderModel1.OrderUID).Return(delivery, nil)
				mock.GetPaymentMock.Expect(ctx, repoOrderModel1.OrderUID).Return(payment, nil)
				mock.ListItemsMock.Expect(ctx, repoOrderModel1.OrderUID).Return(items, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.SetMock.Expect(res1.OrderUID, res1).Return(true)
				return mock
			},
		},
		{
			name: "success - restore cache with multiple orders",
			args: args{
				ctx:   ctx,
				limit: 2,
			},
			err: nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)

				mock.ListOrdersByLastAddedMock.Expect(ctx, 2).Return(
					[]*model.Order{repoOrderModel1, repoOrderModel2},
					nil,
				)

				mock.GetDeliveryMock.Set(func(ctx context.Context, orderID string) (*model.Delivery, error) {
					return delivery, nil
				})
				mock.GetPaymentMock.Set(func(ctx context.Context, orderID string) (*model.Payment, error) {
					return payment, nil
				})
				mock.ListItemsMock.Set(func(ctx context.Context, orderID string) ([]*model.Item, error) {
					return items, nil
				})

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)

				mock.SetMock.When(res1.OrderUID, res1).Then(true)
				mock.SetMock.When(res2.OrderUID, res2).Then(true)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)

			orderRepo := tt.orderRepositoryMock(mc)
			txManager := tt.txManagerMock(mc)
			cacheLRU := tt.cacheMock(mc)

			service := order.NewService(orderRepo, txManager, cacheLRU)

			err := service.RestoreCache(ctx, tt.args.limit)

			require.NoError(t, err)
		})
	}
}
