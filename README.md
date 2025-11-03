# L0

Микросервис для управления заказами с интеграцией Kafka, PostgreSQL и LRU кешированием.

## Описание

Разработан демонстрационный сервис, отображающий данные о заказах. 
Сервис получает данные заказов из очереди сообщений Kafka, 
сохраняет их в базу данных PostgreSQL и кеширует в памяти для быстрого доступа.
Есть метрики(prometheus), трейсинг(jaeger + opentelemetry), также реализована DLQ. 

## Структура проекта
```
.
├── alerts.yaml
├── bin
│   ├── golint
│   ├── goose
│   └── main
├── cmd
│   ├── producer
│   │   └── main.go
│   └── server
│       └── main.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   └── service.go
│   ├── client
│   │   ├── cache
│   │   │   ├── client.go
│   │   │   ├── generate.go
│   │   │   ├── lrucache
│   │   │   │   ├── cache.go
│   │   │   │   └── cache_test.go
│   │   │   └── mocks
│   │   │       └── client_minimock.go
│   │   ├── db
│   │   │   ├── db.go
│   │   │   ├── generate.go
│   │   │   ├── mocks
│   │   │   │   ├── client_minimock.go
│   │   │   │   └── tx_manager_minimock.go
│   │   │   ├── pg
│   │   │   │   ├── client.go
│   │   │   │   └── pg.go
│   │   │   └── transaction
│   │   │       └── transaction.go
│   │   └── kafka
│   │       ├── consumer
│   │       │   ├── consumer.go
│   │       │   └── message_handler.go
│   │       ├── generate.go
│   │       ├── kafka.go
│   │       ├── mocks
│   │       │   ├── consumer_minimock.go
│   │       │   └── producer_minimock.go
│   │       └── producer
│   │           └── producer.go
│   ├── config
│   │   ├── config.go
│   │   └── env
│   │       ├── http_config.go
│   │       ├── jaeger_config.go
│   │       ├── kafka_consumer.go
│   │       └── pg.go
│   ├── metric
│   │   └── metric.go
│   ├── middleware
│   │   └── metric.go
│   ├── model
│   │   └── model.go
│   ├── repository
│   │   ├── generate.go
│   │   ├── mocks
│   │   │   └── order_repository_minimock.go
│   │   ├── order
│   │   │   └── repository.go
│   │   └── repository.go
│   ├── service
│   │   ├── consumer
│   │   │   └── ordersaver
│   │   │       ├── consumer.go
│   │   │       ├── handler.go
│   │   │       └── tests
│   │   │           ├── order_save_handler_test.go
│   │   │           └── run_consumer_test.go
│   │   ├── order
│   │   │   ├── service.go
│   │   │   └── tests
│   │   │       ├── get_order_test.go
│   │   │       └── resotre_cache_test.go
│   │   └── service.go
│   ├── tracing
│   │   └── trace.go
│   └── validator
│       └── validator.go
├── local.env
├── Makefile
├── migrations
│   └── 20250820193524_orders_tables.sql
├── prometheus.yaml
├── README.md
└── templates
    └── index.html

```

## Технологический стек

| Компонент                         | Технология | Версия |
|-----------------------------------|-----------|----|
| [Язык программирования] | Go | 1.24.1 |
| [Веб-фреймворк]         | Gin | Latest |
| [База данных]           | PostgreSQL | 17.2-alpine3.21|
| [Очередь сообщений]     | Kafka | 7.6.1 |
| [Драйвер БД]            | pgx/v5 | Latest |
| [Конструктор запросов]  | Squirrel | Latest |
| [Клиент Kafka]                    | Sarama | Latest |


## Установка и запуск
### Шаг 1: Клонирование репозитория

git clone https://github.com/biryanim/wb_tech_L0.git
cd wb_tech_L0

### Шаг 2: Конфигурация окружения

Создайте файл `local.env`:

```
HTTP_HOST=localhost
HTTP_PORT=8080
PG_DSN=postgres://postgres:postgres@localhost:5432/orders_db
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=order-consumer-group
```



### Шаг 3: Запуск инфраструктуры

`docker-compose up -d`

Команда запустит:
- PostgreSQL на порту 5432
- Kafka на порту 9092
- Zookeeper на порту 2181

### Шаг 4: Применение миграций

`make migration-up`


### Шаг 5: Запуск приложения
Сборка

`make build`

Или запуск напрямую

`go run ./cmd/server/main.go`


Сервер запустится на `http://localhost:8080`

## Использование

### Веб-интерфейс

1. Откройте `http://localhost:8080` в браузере
2. Введите ID заказа (например: `b563feb7b2b84b6test`)
3. Нажмите кнопку поиска
4. Данные заказа отобразятся на странице