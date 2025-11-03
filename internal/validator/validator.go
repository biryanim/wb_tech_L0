package validator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/hashicorp/go-uuid"
)

// ValidationError represents validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationResult represents result of order struct validation
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// ValidateOrder validates order struct
func ValidateOrder(order *model.Order) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	if order == nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "order",
			Message: "order is nil",
		})
		return result
	}

	validateOrderID(order.OrderUID, result)
	validateTrackNumber(order.TrackNumber, result)
	validateEntry(order.Entry, result)
	validateLocale(order.Locale, result)
	validateCustomerID(order.CustomerID, result)
	validateDeliverService(order.DeliveryService, result)
	validateShardKey(order.ShardKey, result)
	validateSmID(order.SmID, result)
	validateDateCreated(order.DateCreated, result)
	validateOofShard(order.OofShard, result)
	validateDelivery(&order.Delivery, result)
	validatePayment(&order.Payment, result)
	validateItems(order.Items, result)

	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

func validateOrderID(orderUID string, result *ValidationResult) {
	_, err := uuid.ParseUUID(orderUID)
	if err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "order_uid",
			Message: err.Error(),
		})
	}
}

func validateTrackNumber(trackNumber string, result *ValidationResult) {
	if strings.TrimSpace(trackNumber) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "track_number",
			Message: "track_number is required",
		})
		return
	}

	if len(trackNumber) < 5 || len(trackNumber) > 50 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "track_number",
			Message: "track_number must be between 5 and 50 characters",
		})
	}
}

func validateEntry(entry string, result *ValidationResult) {
	if strings.TrimSpace(entry) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "entry",
			Message: "entry is required",
		})
	}
}

func validateLocale(locale string, result *ValidationResult) {
	if len(strings.TrimSpace(locale)) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "locale",
			Message: "locale is required",
		})
		return
	}
}

func validateCustomerID(customerID string, result *ValidationResult) {
	if len(strings.TrimSpace(customerID)) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "customer_id",
			Message: "customer_id is required",
		})
	}
}

func validateDeliverService(service string, result *ValidationResult) {
	if len(strings.TrimSpace(service)) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery_service",
			Message: "delivery_service is required",
		})
		return
	}
}

func validateShardKey(shardKey string, result *ValidationResult) {
	if len(strings.TrimSpace(shardKey)) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "shard_key",
			Message: "shard_key is required",
		})
		return
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(shardKey) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "shard_key",
			Message: "shard_key must match ^d+$",
		})
	}
}

func validateSmID(smID int, result *ValidationResult) {
	if smID <= 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "sm_id",
			Message: "sm_id must be greater than 0",
		})
	}
}

func validateDateCreated(dateCreated time.Time, result *ValidationResult) {
	if dateCreated.IsZero() {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "date_created",
			Message: "date_created is required",
		})
		return
	}
}

func validateOofShard(oofShard string, result *ValidationResult) {
	if strings.TrimSpace(oofShard) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "oof_shard",
			Message: "oof_shard is required",
		})
		return
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(oofShard) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "oof_shard",
			Message: "oof_shard must be numeric",
		})
	}
}

func validateDelivery(delivery *model.Delivery, result *ValidationResult) {
	if delivery == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery",
			Message: "delivery is required",
		})
		return
	}

	if strings.TrimSpace(delivery.Name) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.name",
			Message: "delivery name is required",
		})
	}

	if err := validatePhone(delivery.Phone); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.phone",
			Message: err.Error(),
		})
	}

	if err := validateEmail(delivery.Email); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.email",
			Message: err.Error(),
		})
	}

	if strings.TrimSpace(delivery.Zip) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.zip",
			Message: "delivery zip is required",
		})
	}

	if strings.TrimSpace(delivery.City) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.city",
			Message: "delivery city is required",
		})
	}

	if strings.TrimSpace(delivery.Address) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.address",
			Message: "delivery address is required",
		})
	}

	if strings.TrimSpace(delivery.Region) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "delivery.region",
			Message: "delivery region is required",
		})
	}
}

func validatePhone(phone string) error {
	if strings.TrimSpace(phone) == "" {
		return fmt.Errorf("phone is required")
	}

	if !regexp.MustCompile(`^\d{7,15}$`).MatchString(phone) {
		return fmt.Errorf("phone '%s' is invalid (must be +XXXXXXX format with 7-15 digits)", phone)
	}

	return nil
}

func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email is required")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("email '%s' is invalid: %w", email, err)
	}

	return nil
}

func validatePayment(payment *model.Payment, result *ValidationResult) {
	if payment == nil {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment",
			Message: "payment is required",
		})
		return
	}

	if strings.TrimSpace(payment.Transaction) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.transaction",
			Message: "payment transaction is required",
		})
	}

	if strings.TrimSpace(payment.Currency) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.currency",
			Message: "payment currency is required",
		})
	} else if len(payment.Currency) != 3 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.currency",
			Message: "payment currency must be 3-letter ISO code",
		})
	}

	if strings.TrimSpace(payment.Provider) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.provider",
			Message: "payment provider is required",
		})
	}

	if payment.Amount <= 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.amount",
			Message: "payment amount must be greater than 0",
		})
	}

	if payment.PaymentDt == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.payment_dt",
			Message: "payment_dt is required",
		})
	} else {
		paymentTime := time.Unix(payment.PaymentDt, 0)
		if paymentTime.After(time.Now().Add(24 * time.Hour)) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "payment.payment_dt",
				Message: "payment_dt cannot be significantly in the future",
			})
		}
	}

	if strings.TrimSpace(payment.Bank) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.bank",
			Message: "payment bank is required",
		})
	}

	if payment.DeliveryCost < 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.delivery_cost",
			Message: "delivery_cost cannot be negative",
		})
	}

	if payment.GoodsTotal < 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.goods_total",
			Message: "goods_total cannot be negative",
		})
	}

	expectedTotal := float64(payment.GoodsTotal) + payment.DeliveryCost + payment.CustomFee
	if payment.Amount != expectedTotal {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "payment.amount",
			Message: fmt.Sprintf("amount (%f) doesn't match goods_total + delivery_cost + custom_fee", payment.Amount),
		})
	}
}

func validateItems(items []model.Item, result *ValidationResult) {
	if len(items) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "items",
			Message: "items cannot be empty",
		})
		return
	}

	for i, item := range items {
		prefix := fmt.Sprintf("items[%d]", i)

		if item.ChrtID <= 0 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".chrt_id",
				Message: "chrt_id must be greater than 0",
			})
		}

		if strings.TrimSpace(item.TrackNumber) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".track_number",
				Message: "track_number is required",
			})
		}

		if item.Price <= 0 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".price",
				Message: "price must be greater than 0",
			})
		}

		if strings.TrimSpace(item.Rid) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".rid",
				Message: "rid is required",
			})
		}

		if strings.TrimSpace(item.Name) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".name",
				Message: "name is required",
			})
		}

		if item.Sale < 0 || item.Sale > 100 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".sale",
				Message: "sale must be between 0 and 100",
			})
		}

		if item.TotalPrice <= 0 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".total_price",
				Message: "total_price must be greater than 0",
			})
		}

		if item.NmID <= 0 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".nm_id",
				Message: "nm_id must be greater than 0",
			})
		}

		if strings.TrimSpace(item.Brand) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".brand",
				Message: "brand is required",
			})
		}

		if item.Status < 100 || item.Status >= 600 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   prefix + ".status",
				Message: "status must be a valid HTTP status code (100-599)",
			})
		}
	}
}
