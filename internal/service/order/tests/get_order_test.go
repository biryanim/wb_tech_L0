package tests

import (
	"context"
	"errors"
	"testing"
	"time"

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
)

func Test_GetOrder(t *testing.T) {
	type orderRepositoryMockFunc func(mc *minimock.Controller) repository.OrderRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type cacheMockFunc func(mc *minimock.Controller) cache.Client

	type args struct {
		ctx     context.Context
		orderID string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		orderUID        = gofakeit.UUID()
		trackNumber     = gofakeit.LetterN(10)
		entry           = gofakeit.LetterN(4)
		locale          = gofakeit.LanguageAbbreviation()
		customerID      = gofakeit.Username()
		shardKey        = gofakeit.DigitN(1)
		smID            = gofakeit.IntRange(1, 100)
		oofShard        = gofakeit.DigitN(1)
		deliveryService = gofakeit.Company()
		dateCreated     = gofakeit.Date()

		errOrderNotFound    = errors.New("order not found")
		errDeliveryNotFound = errors.New("delivery not found")
		errPaymentNotFound  = errors.New("payment not found")

		repoError = errors.New("repository error")

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

		res = &model.Order{
			OrderUID:          orderUID,
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

		repoOrderModel = &model.Order{
			OrderUID:          orderUID,
			TrackNumber:       trackNumber,
			Entry:             entry,
			Locale:            locale,
			InternalSignature: "",
			CustomerID:        customerID,
			DeliveryService:   res.DeliveryService,
			ShardKey:          shardKey,
			SmID:              smID,
			DateCreated:       res.DateCreated,
			OofShard:          oofShard,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *model.Order
		err                 error
		orderRepositoryMock orderRepositoryMockFunc
		txManagerMock       txManagerMockFunc
		cacheMock           cacheMockFunc
	}{
		{
			name: "success - order from cache",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: res,
			err:  nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(res)
				return mock
			},
		},
		{
			name: "success - order from db",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: res,
			err:  nil,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(repoOrderModel, nil)
				mock.GetDeliveryMock.Expect(ctx, orderUID).Return(delivery, nil)
				mock.GetPaymentMock.Expect(ctx, orderUID).Return(payment, nil)
				mock.ListItemsMock.Expect(ctx, orderUID).Return(items, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				mock.SetMock.Expect(orderUID, res).Return(true)
				return mock
			},
		},
		{
			name: "error - order not found in repository",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  errOrderNotFound,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(nil, errOrderNotFound)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
		{
			name: "error - failed to get order from repository",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  repoError,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(nil, repoError)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
		{
			name: "error - delivery not found",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  errDeliveryNotFound,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(repoOrderModel, nil)
				mock.GetDeliveryMock.Expect(ctx, orderUID).Return(nil, errDeliveryNotFound)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
		{
			name: "error - failed to get delivery from repository",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  repoError,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(repoOrderModel, nil)
				mock.GetDeliveryMock.Expect(ctx, orderUID).Return(nil, repoError)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
		{
			name: "error - payment not found in repository",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  errPaymentNotFound,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(repoOrderModel, nil)
				mock.GetDeliveryMock.Expect(ctx, orderUID).Return(delivery, nil)
				mock.GetPaymentMock.Expect(ctx, orderUID).Return(nil, errPaymentNotFound)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
		{
			name: "error - failed to get payment from repository",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  repoError,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(repoOrderModel, nil)
				mock.GetDeliveryMock.Expect(ctx, orderUID).Return(delivery, nil)
				mock.GetPaymentMock.Expect(ctx, orderUID).Return(nil, repoError)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
		{
			name: "error - failed to list items",
			args: args{
				ctx:     ctx,
				orderID: orderUID,
			},
			want: nil,
			err:  repoError,
			orderRepositoryMock: func(mc *minimock.Controller) repository.OrderRepository {
				mock := repoMock.NewOrderRepositoryMock(mc)
				mock.GetOrderMock.Expect(ctx, orderUID).Return(repoOrderModel, nil)
				mock.GetDeliveryMock.Expect(ctx, orderUID).Return(delivery, nil)
				mock.GetPaymentMock.Expect(ctx, orderUID).Return(payment, nil)
				mock.ListItemsMock.Expect(ctx, orderUID).Return(nil, repoError)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMock.NewTxManagerMock(mc)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.Client {
				mock := cacheMock.NewClientMock(mc)
				mock.GetMock.Expect(orderUID).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := tt.orderRepositoryMock(mc)
			txManager := tt.txManagerMock(mc)
			cacheLRU := tt.cacheMock(mc)

			service := order.NewService(orderRepo, txManager, cacheLRU)

			result, err := service.GetOrder(tt.args.ctx, tt.args.orderID)

			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
