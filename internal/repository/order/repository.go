package order

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/biryanim/wb_tech_L0/internal/client/db"
	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/google/uuid"

	def "github.com/biryanim/wb_tech_L0/internal/repository"
)

var _ def.OrderRepository = (*repo)(nil)

type repo struct {
	db db.Client
	qb squirrel.StatementBuilderType
}

func NewRepository(db db.Client) *repo {
	return &repo{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *repo) CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	query, args, err := r.qb.
		Insert("orders").
		Columns(
			"order_uid",
			"track_number",
			"entry",
			"locale",
			"internal_signature",
			"customer_id",
			"delivery_service",
			"shardKey",
			"sm_id",
			"date_created",
			"oof_shard",
		).
		Values(
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SmID,
			order.DateCreated,
			order.OofShard,
		).
		Suffix("RETURNING \"order_uid\"").
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id uuid.UUID
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert order: %w", err)
	}

	return id, nil
}

func (r *repo) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	query, args, err := r.qb.
		Select(
			"track_number",
			"entry",
			"locale",
			"internal_signature",
			"customer_id",
			"delivery_service",
			"shardKey",
			"sm_id",
			"date_created",
			"oof_shard",
		).
		From("orders").
		Where(squirrel.Eq{"order_uid": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var order model.Order
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query order: %w", err)
	}
	order.OrderUID = orderID

	return &order, nil
}

func (r *repo) CreateDelivery(ctx context.Context, orderID uuid.UUID, delivery *model.Delivery) (int64, error) {
	query, args, err := r.qb.
		Insert("deliveries").
		Columns(
			"order_uid",
			"name",
			"phone",
			"zip",
			"city",
			"address",
			"region",
			"email",
		).
		Values(
			orderID,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email,
		).
		Suffix("RETURNING \"order_uid\"").
		ToSql()
	if err != nil {
		return int64(0), fmt.Errorf("failed to build insert query: %w", err)
	}
	var id int64
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return int64(0), fmt.Errorf("failed to insert delivery: %w", err)
	}

	return id, nil
}

func (r *repo) GetDelivery(ctx context.Context, orderID uuid.UUID) (*model.Delivery, error) {
	query, args, err := r.qb.
		Select(
			"name",
			"phone",
			"zip",
			"city",
			"address",
			"region",
			"email",
		).
		From("deliveries").
		Where(squirrel.Eq{"order_uid": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var delivery model.Delivery
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query delivery: %w", err)
	}

	return &delivery, nil
}

func (r *repo) CreatePayment(ctx context.Context, orderID uuid.UUID, payment *model.Payment) (int64, error) {
	query, args, err := r.qb.
		Insert("payments").
		Columns(
			"order_uid",
			"transaction",
			"request_id",
			"currency",
			"provider",
			"amount",
			"payment_dt",
			"bank",
			"delivery_cost",
			"goods_total",
			"custom_fee",
		).
		Values(
			orderID,
			payment.Transaction,
			payment.RequestID,
			payment.Currency,
			payment.Provider,
			payment.Amount,
			payment.PaymentDt,
			payment.Bank,
			payment.DeliveryCost,
			payment.GoodsTotal,
			payment.CustomFee,
		).
		Suffix("RETURNING order_uid").
		ToSql()
	if err != nil {
		return int64(0), fmt.Errorf("failed to build insert query: %w", err)
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return int64(0), fmt.Errorf("failed to insert payment: %w", err)
	}

	return id, nil
}

func (r *repo) GetPayment(ctx context.Context, orderID uuid.UUID) (*model.Payment, error) {
	query, args, err := r.qb.
		Select(
			"transaction",
			"request_id",
			"currency",
			"provider",
			"amount",
			"payment_dt",
			"bank",
			"delivery_cost",
			"goods_total",
			"custom_fee",
		).
		From("payments").
		Where(squirrel.Eq{"order_uid": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var payment model.Payment
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment: %w", err)
	}
	return &payment, nil
}

func (r *repo) CreateItem(ctx context.Context, orderID uuid.UUID, item *model.Item) error {
	query, args, err := r.qb.
		Insert("items").
		Columns(
			"order_uid",
			"chrt_id",
			"track_number",
			"price",
			"rid",
			"name",
			"sale",
			"size",
			"total_price",
			"nm_id",
			"brand",
			"status",
		).
		Values(
			orderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert item: %w", err)
	}

	return nil
}

func (r *repo) ListItems(ctx context.Context, orderID uuid.UUID) ([]*model.Item, error) {
	query, args, err := r.qb.
		Select(
			"chrt_id",
			"track_number",
			"price",
			"rid",
			"name",
			"sale",
			"size",
			"total_price",
			"nm_id",
			"brand",
			"status",
		).
		From("items").
		Where(squirrel.Eq{"order_uid": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []*model.Item
	for rows.Next() {
		item := &model.Item{}
		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to query items: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	return items, nil
}
