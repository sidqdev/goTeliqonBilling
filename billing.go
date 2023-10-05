package teliqonBillling

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/levigross/grequests"
)

type BillingAPI struct {
	Config Config
	Urls   ApiUrls
}

func (b *BillingAPI) Load() error {
	if b.Config.ApiUrl == "" {
		b.Config.ApiUrl = API_URL
	}

	if !strings.HasSuffix(b.Config.ApiUrl, "/") {
		return ErrIncorrectApiUrl
	}

	b.LoadUrls()
	return b.Ping()
}

func (b *BillingAPI) LoadUrls() {
	url := b.Config.ApiUrl
	b.Urls = ApiUrls{
		User:              url + "user/",
		Deposit:           url + "deposit/",
		Ping:              url + "ping/",
		Withdrawal:        url + "withdrawal/",
		InSystemTransfer:  url + "transfer/in/",
		OutSystemTransfer: url + "transfer/out/",
		Transaction:       url + "transactions/",
		Subscriptions:     url + "subscriptions/",
	}
}

func (b BillingAPI) CheckResponseOnError(resp *grequests.Response) error {
	if resp.StatusCode/100 == 2 {
		return nil
	}
	err_body := ErrorResponse{}
	if err := resp.JSON(&err_body); err != nil {
		return err
	}
	err, ok := ErrorsRegister[err_body.Code]
	if !ok {
		return errors.New(err_body.Message)
	}
	return err
}

func (b BillingAPI) Headers() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Token %s", b.Config.ApiToken),
	}
}

func (b BillingAPI) parseResponse(resp *grequests.Response, err error) func(v interface{}) error {
	return func(v interface{}) error {
		if err != nil {
			return err
		}
		err = b.CheckResponseOnError(resp)
		if err != nil {
			return err
		}
		return resp.JSON(v)
	}
}

func (b BillingAPI) Ping() error {
	ro := &grequests.RequestOptions{
		Headers: b.Headers(),
		Params: map[string]string{
			"environment": strconv.Itoa(b.Config.EnvironmentID),
		},
		InsecureSkipVerify: true,
	}
	ping_response := struct {
		Status string `json:"status"`
	}{}
	err := b.parseResponse(grequests.Get(b.Urls.Ping, ro))(&ping_response)
	if err != nil {
		return err
	}
	if ping_response.Status != "ok" {
		return ErrPanic
	}
	return nil
}

func (b BillingAPI) CreateUser(billingUser BillingUser) (BillingUser, error) {
	data := billingUser.ToMap()
	data["environment"] = b.Config.EnvironmentID

	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		JSON:               data,
		InsecureSkipVerify: true,
	}
	createdBillingUser := BillingUser{}
	err := b.parseResponse(grequests.Post(b.Urls.User, ro))(&createdBillingUser)

	return createdBillingUser, err
}

func (b BillingAPI) GetUser(uniqueID string) (BillingUser, error) {
	data := map[string]string{
		"environment": strconv.Itoa(b.Config.EnvironmentID),
		"unique_id":   uniqueID,
	}

	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		Params:             data,
		InsecureSkipVerify: true,
	}
	billingUser := BillingUser{}
	err := b.parseResponse(grequests.Get(b.Urls.User, ro))(&billingUser)

	return billingUser, err
}

func (b BillingAPI) GetUsers() ([]BillingUser, error) {
	data := map[string]string{
		"environment": strconv.Itoa(b.Config.EnvironmentID),
	}

	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		Params:             data,
		InsecureSkipVerify: true,
	}

	users := struct {
		Users []BillingUser `json:"users"`
	}{}

	err := b.parseResponse(grequests.Get(b.Urls.User, ro))(&users)
	return users.Users, err
}

func (b BillingAPI) Deposit(params DepositParams) (Transaction, error) {
	data := map[string]interface{}{
		"environment": b.Config.EnvironmentID,
		"unique_id":   params.UniqueID,
		"amount":      params.Amount,
		"comment":     params.Comment,
	}
	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		JSON:               data,
		InsecureSkipVerify: true,
	}
	transaction := Transaction{}
	err := b.parseResponse(grequests.Post(b.Urls.Deposit, ro))(&transaction)

	return transaction, err
}

func (b BillingAPI) Withdrawal(params WithdrawalParams) (Transaction, error) {
	data := map[string]interface{}{
		"environment":   b.Config.EnvironmentID,
		"unique_id":     params.UniqueID,
		"amount":        params.Amount,
		"comment":       params.Comment,
		"fee_on_sender": params.FeeOnSender,
	}
	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		JSON:               data,
		InsecureSkipVerify: true,
	}
	transaction := Transaction{}
	err := b.parseResponse(grequests.Post(b.Urls.Withdrawal, ro))(&transaction)

	return transaction, err
}

func (b BillingAPI) ProcessWithdrawal(params ProcessWithdrawalParams) (Transaction, error) {
	data := map[string]interface{}{
		"environment":    b.Config.EnvironmentID,
		"unique_id":      params.UniqueID,
		"transaction_id": params.TransactionID,
		"status":         params.Status,
	}
	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		JSON:               data,
		InsecureSkipVerify: true,
	}
	transaction := Transaction{}
	err := b.parseResponse(grequests.Patch(b.Urls.Withdrawal, ro))(&transaction)

	return transaction, err
}

func (b BillingAPI) InSystemTransfer(params InSystemTransferParams) (Transaction, error) {
	data := map[string]interface{}{
		"environment":    b.Config.EnvironmentID,
		"from_unique_id": params.FromUniqueID,
		"to_unique_id":   params.ToUniqueID,
		"amount":         params.Amount,
		"comment":        params.Comment,
		"fee_on_sender":  params.FeeOnSender,
	}
	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		JSON:               data,
		InsecureSkipVerify: true,
	}
	transaction := Transaction{}
	err := b.parseResponse(grequests.Post(b.Urls.InSystemTransfer, ro))(&transaction)

	return transaction, err
}

func (b BillingAPI) OutSystemTransfer(params OutSystemTransfer) (Transaction, error) {
	outSystemService := &params.ServiceUniqueID
	if params.ServiceUniqueID == "" {
		outSystemService = nil
	}
	if params.Quantity == 0 {
		params.Quantity = 1
	}

	data := map[string]interface{}{
		"environment":       b.Config.EnvironmentID,
		"unique_id":         params.UniqueID,
		"amount":            params.Amount,
		"comment":           params.Comment,
		"service_unique_id": outSystemService,
		"quantity":          params.Quantity,
	}
	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		JSON:               data,
		InsecureSkipVerify: true,
	}
	transaction := Transaction{}
	err := b.parseResponse(grequests.Post(b.Urls.OutSystemTransfer, ro))(&transaction)

	return transaction, err
}

func (b BillingAPI) GetUserTransactions(uniqueID string, filterConfig *TransactionFilterConfig) ([]Transaction, int, error) {
	data := map[string]string{
		"environment": strconv.Itoa(b.Config.EnvironmentID),
		"unique_id":   uniqueID,
	}

	if filterConfig != nil {
		for k, v := range filterConfig.ToParamsMap() {
			data[k] = v
		}
	}

	ro := &grequests.RequestOptions{
		Headers:            b.Headers(),
		Params:             data,
		InsecureSkipVerify: true,
	}
	transactions := struct {
		Transactions []Transaction `json:"transactions"`
		Count        int           `json:"count"`
	}{}
	err := b.parseResponse(grequests.Get(b.Urls.Transaction, ro))(&transactions)

	return transactions.Transactions, transactions.Count, err
}

func NewBillingAPI(config Config) (BillingAPI, error) {
	billingApi := BillingAPI{
		Config: config,
	}
	err := billingApi.Load()
	return billingApi, err
}
