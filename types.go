package teliqonBillling

import (
	"encoding/json"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Code    string `json:"error_code"`
	Message string `json:"error"`
}

type Config struct {
	ApiUrl        string
	EnvironmentID int
	ApiToken      string
}

type ApiUrls struct {
	User              string
	Deposit           string
	Ping              string
	Withdrawal        string
	InSystemTransfer  string
	OutSystemTransfer string
	Transaction       string
	Subscriptions     string
	OutSystemService  string
}

type BillingUser struct {
	UniqueID     string  `json:"unique_id"`
	CreditLimit  float64 `json:"credit_limit"`
	Balance      float64 `json:"balance"` // Readonly field
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	MobileNumber string  `json:"mobile_number"`
	WorkNumber   string  `json:"work_number"`
	Email        string  `json:"email"`
	CompanyName  string  `json:"company_name"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Address      string  `json:"address"`
	ZipCode      string  `json:"zip_code"`
	Description  string  `json:"description"`
}

func (u BillingUser) ToMap() map[string]interface{} {
	data := map[string]interface{}{}
	body, _ := json.Marshal(u)
	json.Unmarshal(body, &data)
	return data
}

type Transfer struct {
	Income bool        `json:"income"`
	User   BillingUser `json:"user"`
}

type OutSystemService struct {
	UniqueID    string   `json:"unique_id"`
	Prefix      string   `json:"prefix"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Cost        *float64 `json:"cost"`
}

type Transaction struct {
	ID               int               `json:"id"`
	Type             string            `json:"type"`
	Status           string            `json:"string"`
	Comment          string            `json:"comment"`
	CreatedAt        time.Time         `json:"created_at"`
	Amount           float64           `json:"amount"`
	Fee              float64           `json:"fee"`
	BalanceBefore    float64           `json:"balance_before"`
	BalanceAfter     float64           `json:"balance_after"`
	Transfer         *Transfer         `json:"transfer"`
	OutSystemService *OutSystemService `json:"out_system_service"`
}

type DepositParams struct {
	UniqueID string
	Amount   float64
	Comment  string
}

type WithdrawalParams struct {
	UniqueID    string
	Amount      float64
	Comment     string
	FeeOnSender bool
}

type ProcessWithdrawalParams struct {
	UniqueID      string
	TransactionID int
	Status        bool
}

type InSystemTransferParams struct {
	FromUniqueID string
	ToUniqueID   string
	Amount       float64
	Comment      string
	FeeOnSender  bool
}

type OutSystemTransferParams struct {
	UniqueID        string
	Amount          float64
	Comment         string
	FeeOnSender     bool
	ServiceUniqueID string
	Quantity        int
}

type TransactionFilterConfig struct {
	Limit    *int
	Offset   *int
	FromDate *time.Time
	ToDate   *time.Time
	OrderBy  *string
	Query    *string
}

func (t TransactionFilterConfig) ToParamsMap() map[string]string {
	params := map[string]string{}
	if t.Limit != nil {
		params["limit"] = strconv.Itoa(*t.Limit)
	}
	if t.Offset != nil {
		params["offset"] = strconv.Itoa(*t.Offset)
	}
	if t.OrderBy != nil {
		params["order_by"] = *t.OrderBy
	}
	if t.FromDate != nil {
		params["from_date"] = t.FromDate.Format(time.RFC3339)
	}
	if t.ToDate != nil {
		params["to_date"] = t.ToDate.Format(time.RFC3339)
	}
	if t.Query != nil {
		params["query"] = *t.Query
	}

	return params
}

type OutSystemServicesFilterConfig struct {
	Prefix string
}

type SubscriptionPlan struct {
	UniqueID           string            `json:"unique_id"`
	OutSystemService   *OutSystemService `json:"out_system_service"`
	FirstPaymentAmount float64           `json:"first_payment_amount"`
	Amount             float64           `json:"amount"`
	DaysFrequency      int               `json:"days_frequency"`
	Type               string            `json:"type"`
}

type Subscription struct {
	Plan           SubscriptionPlan `json:"plan"`
	Count          int              `json:"count"`
	Amount         float64          `json:"amount"`
	Comment        string           `json:"comment"`
	FirstPaymentAt time.Time        `json:"first_payment_at"`
	NextPaymentAt  time.Time        `json:"next_payment_at"`
	IsActive       bool             `json:"is_active"`
	IsDisabled     bool             `json:"is_disabled"`
}

type SubscriptionsFilterConfig struct {
	Prefix     string
	OrderBy    string
	Limit      int
	Offset     int
	IsDisabled *bool
}
