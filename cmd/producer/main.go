package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

// Структуры данных (те же, что и в основном сервисе)
type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   int     `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtID      int64   `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int64   `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

// Тестовые данные
var (
	cities = []string{"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg", "Kazan", "Nizhny Novgorod", "Chelyabinsk", "Samara", "Omsk", "Rostov-on-Don"}
	names  = []string{"Ivan Petrov", "Maria Sidorova", "Alexey Kozlov", "Elena Volkova", "Dmitry Smirnov", "Anna Kuznetsova", "Pavel Morozov", "Olga Novikova"}
	brands = []string{"Nike", "Adidas", "Puma", "Reebok", "New Balance", "Converse", "Vans", "Under Armour"}
	items  = []string{"T-Shirt", "Jeans", "Sneakers", "Jacket", "Hoodie", "Shorts", "Dress", "Backpack", "Watch", "Sunglasses"}
	banks  = []string{"Sberbank", "VTB", "Gazprombank", "Alfa Bank", "Raiffeisen", "Tinkoff"}
)

// generateRandomOrder создает случайный заказ
func generateRandomOrder() Order {
	orderUID := fmt.Sprintf("order_%d_%d", time.Now().Unix(), rand.Intn(10000))
	trackNumber := fmt.Sprintf("TRACK%d", rand.Intn(1000000))

	// Генерация товаров
	itemCount := rand.Intn(3) + 1 // 1-3 товара
	var orderItems []Item
	totalPrice := 0.0

	for i := 0; i < itemCount; i++ {
		price := rand.Float64() // 500-5500
		sale := rand.Intn(50)   // 0-49% скидка
		finalPrice := price * (100 - float64(sale)) / 100
		totalPrice += finalPrice

		item := Item{
			ChrtID:      int64(rand.Intn(1000000) + 1000000),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         fmt.Sprintf("rid_%d", rand.Intn(1000000)),
			Name:        items[rand.Intn(len(items))],
			Sale:        sale,
			Size:        []string{"XS", "S", "M", "L", "XL"}[rand.Intn(5)],
			TotalPrice:  finalPrice,
			NmID:        int64(rand.Intn(10000000) + 1000000),
			Brand:       brands[rand.Intn(len(brands))],
			Status:      []int{200, 201, 202}[rand.Intn(3)],
		}
		orderItems = append(orderItems, item)
	}

	deliveryCost := float64(rand.Intn(1000) + 200)

	return Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    names[rand.Intn(len(names))],
			Phone:   fmt.Sprintf("+7%d", rand.Int63n(9000000000)+1000000000),
			Zip:     fmt.Sprintf("%d", rand.Intn(900000)+100000),
			City:    cities[rand.Intn(len(cities))],
			Address: fmt.Sprintf("Street %d, Building %d", rand.Intn(100)+1, rand.Intn(200)+1),
			Region:  "Central",
			Email:   fmt.Sprintf("user%d@example.com", rand.Intn(10000)),
		},
		Payment: Payment{
			Transaction:  orderUID,
			RequestID:    fmt.Sprintf("req_%d", rand.Intn(1000000)),
			Currency:     "RUB",
			Provider:     "wbpay",
			Amount:       totalPrice + deliveryCost,
			PaymentDt:    time.Now().Unix(),
			Bank:         banks[rand.Intn(len(banks))],
			DeliveryCost: deliveryCost,
			GoodsTotal:   rand.Intn(100),
			CustomFee:    float64(rand.Intn(100)),
		},
		Items:             orderItems,
		Locale:            []string{"ru", "en"}[rand.Intn(2)],
		InternalSignature: "",
		CustomerID:        fmt.Sprintf("customer_%d", rand.Intn(100000)),
		DeliveryService:   []string{"cdek", "boxberry", "pickpoint", "dhl"}[rand.Intn(4)],
		ShardKey:          fmt.Sprintf("%d", rand.Intn(10)),
		SmID:              rand.Intn(1000) + 1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

const (
	brokerAddress = "localhost:9092"
	topicName     = "order-topic"
)

func main() {
	producer, err := newSyncProducer(strings.Split(brokerAddress, ","))
	if err != nil {
		log.Fatalf("failed to start producer: %v", err)
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("failed to close producer: %v\n", err)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		order := generateRandomOrder()
		data, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("failed to marshal order: %v", err)
		}
		msg := &sarama.ProducerMessage{
			Topic: topicName,
			Value: sarama.StringEncoder(data),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("failed to send message in Kafka: %v\n", err.Error())
			return
		}

		log.Printf("message sent to partition %d at offset %d\n", partition, offset)
	}

}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
