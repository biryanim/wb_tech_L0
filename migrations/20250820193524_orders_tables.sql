-- +goose Up
-- +goose StatementBegin
create table orders (
    order_uid varchar(255) primary key,
    track_number varchar(255) not null,
    entry varchar(255) not null,
    locale varchar(10) not null,
    internal_signature varchar(255),
    customer_id varchar(255) not null,
    delivery_service varchar(255) not null,
    shardKey varchar(255) not null,
    sm_id int not null,
    date_created timestamp not null,
    oof_shard varchar(255) not null
);

create table deliveries(
    order_uid varchar(255) primary key,
    name varchar(255) not null,
    phone varchar(50) not null,
    zip varchar(20) not null,
    city varchar(255) not null,
    address text not null,
    region varchar(255) not null,
    email varchar(255) not null
);

create table payments(
    order_uid varchar(255) primary key,
    transaction varchar not null,
    request_id varchar(255),
    currency varchar(10) not null,
    provider varchar(255) not null,
    amount decimal not null,
    payment_dt int not null,
    bank varchar(255) not null,
    delivery_cost decimal not null,
    goods_total int not null,
    custom_fee decimal not null
);

create table items(
    id int generated always as identity primary key,
    order_uid varchar(255) not null,
    chrt_id int not null,
    track_number varchar(255) not null,
    price decimal not null,
    rid varchar(255) not null,
    name varchar(255) not null,
    sale int not null,
    size varchar(50) not null,
    total_price decimal not null,
    nm_id int not null,
    brand varchar(255) not null,
    status int not null,
    foreign key (order_uid) references orders(order_uid) on delete cascade
);

create index idx_orders_date_created on orders(date_created);
create index idx_orders_customer_id on orders(customer_id);
create index idx_deliveries_order_id on deliveries(order_uid);
create index idx_payments_order_id on payments(order_uid);
create index idx_items_order_id on items(order_uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_orders_date_created;
drop index idx_orders_customer_id;
drop index idx_deliveries_order_id;
drop index idx_payments_order_id;
drop index idx_items_order_id;

drop table items;
drop table payments;
drop table deliveries;
drop table orders;
-- +goose StatementEnd
