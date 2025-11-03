# L0

Микросервис для управления заказами с интеграцией Kafka, PostgreSQL и LRU кешированием.

## Описание

Разработан демонстрационный сервис, отображающий данные о заказах. 
Сервис получает данные заказов из очереди сообщений Kafka, 
сохраняет их в базу данных PostgreSQL и кеширует в памяти для быстрого доступа.
Есть метрики и трэйсинг

## Структура проекта
```
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
│   │   │   └── lrucache
│   │   │       ├── cache.go
│   │   │       └── cache_test.go
│   │   ├── db
│   │   │   ├── db.go
│   │   │   ├── pg
│   │   │   │   ├── client.go
│   │   │   │   └── pg.go
│   │   │   └── transaction
│   │   │       └── transaction.go
│   │   └── kafka
│   │       ├── consumer
│   │       │   ├── consumer.go
│   │       │   └── message_handler.go
│   │       └── kafka.go
│   ├── config
│   │   ├── config.go
│   │   └── env
│   │       ├── http_config.go
│   │       ├── kafka_consumer.go
│   │       └── pg.go
│   ├── model
│   │   └── model.go
│   ├── repository
│   │   ├── order
│   │   │   └── repository.go
│   │   └── repository.go
│   └── service
│       ├── consumer
│       │   └── ordersaver
│       │       ├── consumer.go
│       │       └── handler.go
│       ├── order
│       │   └── service.go
│       └── service.go
├── local.env
├── Makefile
├── migrations
│   └── 20250820193524_orders_tables.sql
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