package payriff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Config holds the configuration for the Payriff SDK
type Config struct {
	BaseURL            string
	SecretKey          string
	DefaultCallbackURL string
	DefaultLanguage    Language
	DefaultCurrency    Currency
}

// SDK represents the Payriff payment gateway client
type SDK struct {
	baseURL            string
	secretKey          string
	defaultCallbackURL string
	defaultLanguage    Language
	defaultCurrency    Currency
	client             *http.Client
}

// Language represents supported language codes
type Language string

// Currency represents supported currency codes
type Currency string

// Operation represents supported payment operations
type Operation string

// Status represents supported payment statuses
type Status string

// ResultCode represents the possible API result codes
type ResultCode string

const (
	LanguageAZ Language = "AZ"
	LanguageEN Language = "EN"
	LanguageRU Language = "RU"
)

const (
	CurrencyAZN Currency = "AZN"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

const (
	OperationPurchase Operation = "PURCHASE"
	OperationPreAuth  Operation = "PRE_AUTH"
)

const (
	StatusCreated         Status = "CREATED"
	StatusApproved        Status = "APPROVED"
	StatusCanceled        Status = "CANCELED"
	StatusDeclined        Status = "DECLINED"
	StatusRefunded        Status = "REFUNDED"
	StatusPreAuthApproved Status = "PREAUTH_APPROVED"
	StatusExpired         Status = "EXPIRED"
	StatusReverse         Status = "REVERSE"
	StatusPartialRefund   Status = "PARTIAL_REFUND"
)

const (
	ResultCodeSuccess           ResultCode = "00000"
	ResultCodeSuccessGateway    ResultCode = "00"
	ResultCodeSuccessApprove    ResultCode = "APPROVED"
	ResultCodeSuccessPreauth    ResultCode = "PREAUTH-APPROVED"
	ResultCodeWarning           ResultCode = "01000"
	ResultCodeError             ResultCode = "15000"
	ResultCodeInvalidParameters ResultCode = "15400"
	ResultCodeUnauthorized      ResultCode = "14010"
	ResultCodeTokenNotPresent   ResultCode = "14013"
	ResultCodeInvalidToken      ResultCode = "14014"
)

// OrderPayload represents the response payload for order creation
type OrderPayload struct {
	OrderID       string `json:"orderId"`
	PaymentURL    string `json:"paymentUrl"`
	TransactionID int64  `json:"transactionId"`
}

// CardDetails represents saved card information
type CardDetails struct {
	MaskedPan      string `json:"maskedPan"`
	Brand          string `json:"brand"`
	CardHolderName string `json:"cardHolderName"`
}

// Transaction represents a payment transaction
type Transaction struct {
	UUID             string      `json:"uuid"`
	CreatedDate      string      `json:"createdDate"`
	Status           Status      `json:"status"`
	Channel          string      `json:"channel"`
	ChannelType      string      `json:"channelType"`
	RequestRRN       string      `json:"requestRrn"`
	ResponseRRN      *string     `json:"responseRrn"`
	Pan              string      `json:"pan"`
	PaymentWay       string      `json:"paymentWay"`
	CardDetails      CardDetails `json:"cardDetails"`
	CardUUID         *string     `json:"cardUuid,omitempty"`
	MerchantCategory string      `json:"merchantCategory"`
	Installment      struct {
		Type   *string `json:"type"`
		Period *string `json:"period"`
	} `json:"installment"`
	DeliveryAddress *string `json:"deliveryAddress"`
}

// OrderInfo represents detailed order information
type OrderInfo struct {
	OrderID        string        `json:"orderId"`
	InvoiceUUID    *string       `json:"invoiceUuid"`
	Amount         float64       `json:"amount"`
	CurrencyType   Currency      `json:"currencyType"`
	MerchantName   string        `json:"merchantName"`
	CommissionRate *float64      `json:"commissionRate,omitempty"`
	OperationType  Operation     `json:"operationType"`
	PaymentStatus  Status        `json:"paymentStatus"`
	Auto           bool          `json:"auto"`
	CreatedDate    string        `json:"createdDate"`
	Description    string        `json:"description"`
	Transactions   []Transaction `json:"transactions,omitempty"`
}

// CreateOrderRequest represents parameters for creating a new order
type CreateOrderRequest struct {
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Operation   Operation `json:"operation"`
	CardSave    bool      `json:"cardSave"`
	Language    Language  `json:"language,omitempty"`
	Currency    Currency  `json:"currency,omitempty"`
	CallbackURL string    `json:"callbackUrl,omitempty"`
}

// RefundRequest represents parameters for refund operation
type RefundRequest struct {
	Amount  float64 `json:"amount"`
	OrderID string  `json:"orderId"`
}

// CompleteRequest represents parameters for complete operation
type CompleteRequest struct {
	Amount  float64 `json:"amount"`
	OrderID string  `json:"orderId"`
}

// AutoPayRequest represents parameters for automatic payment
type AutoPayRequest struct {
	CardUUID    string    `json:"cardUuid"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Operation   Operation `json:"operation"`
	Currency    Currency  `json:"currency,omitempty"`
	CallbackURL string    `json:"callbackUrl,omitempty"`
}

// Response represents the base API response structure
type Response struct {
	Code            ResultCode      `json:"code"`
	Message         string          `json:"message"`
	Route           string          `json:"route"`
	InternalMessage *string         `json:"internalMessage"`
	ResponseID      string          `json:"responseId"`
	Payload         json.RawMessage `json:"payload"`
}

// ApiResponse represents a generic API response with typed payload
type ApiResponse[T any] struct {
	Code            ResultCode `json:"code"`
	Message         string     `json:"message"`
	Route           string     `json:"route"`
	InternalMessage *string    `json:"internalMessage"`
	ResponseID      string     `json:"responseId"`
	Payload         T          `json:"payload"`
}

// NewSDK creates a new instance of the Payriff SDK
func NewSDK(config Config) *SDK {
	// Set default base URL
	if config.BaseURL == "" {
		config.BaseURL = "https://api.payriff.com/api/v3"
	}

	// Set default secret key from environment
	if config.SecretKey == "" {
		config.SecretKey = os.Getenv("PAYRIFF_SECRET_KEY")
	}

	// Set default callback URL from environment
	if config.DefaultCallbackURL == "" {
		config.DefaultCallbackURL = os.Getenv("PAYRIFF_CALLBACK_URL")
	}

	// Set default language
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = LanguageAZ
	}

	// Set default currency
	if config.DefaultCurrency == "" {
		config.DefaultCurrency = CurrencyAZN
	}

	return &SDK{
		baseURL:            config.BaseURL,
		secretKey:          config.SecretKey,
		defaultCallbackURL: config.DefaultCallbackURL,
		defaultLanguage:    config.DefaultLanguage,
		defaultCurrency:    config.DefaultCurrency,
		client:             &http.Client{},
	}
}

// makeRequest handles HTTP requests to the Payriff API
func (s *SDK) makeRequest(endpoint string, method string, body interface{}) (*Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, s.baseURL+endpoint, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", s.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// CreateOrder creates a new payment order
func (s *SDK) CreateOrder(req CreateOrderRequest) (*ApiResponse[OrderPayload], error) {
	// Apply defaults if values are not provided
	if req.Language == "" {
		req.Language = s.defaultLanguage
	}
	if req.Currency == "" {
		req.Currency = s.defaultCurrency
	}
	if req.CallbackURL == "" {
		req.CallbackURL = s.defaultCallbackURL
	}

	resp, err := s.makeRequest("/orders", http.MethodPost, req)
	if err != nil {
		return nil, err
	}

	var result ApiResponse[OrderPayload]
	if err := json.Unmarshal(resp.Payload, &result.Payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order payload: %w", err)
	}

	// Copy response metadata
	result.Code = resp.Code
	result.Message = resp.Message
	result.Route = resp.Route
	result.InternalMessage = resp.InternalMessage
	result.ResponseID = resp.ResponseID

	return &result, nil
}

// GetOrderInfo retrieves information about an existing order
func (s *SDK) GetOrderInfo(orderID string) (*ApiResponse[OrderInfo], error) {
	resp, err := s.makeRequest(fmt.Sprintf("/orders/%s", orderID), http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var result ApiResponse[OrderInfo]
	if err := json.Unmarshal(resp.Payload, &result.Payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order info: %w", err)
	}

	// Copy response metadata
	result.Code = resp.Code
	result.Message = resp.Message
	result.Route = resp.Route
	result.InternalMessage = resp.InternalMessage
	result.ResponseID = resp.ResponseID

	return &result, nil
}

// Refund initiates a refund for an order
func (s *SDK) Refund(req RefundRequest) (*ApiResponse[json.RawMessage], error) {
	resp, err := s.makeRequest("/refund", http.MethodPost, req)
	if err != nil {
		return nil, err
	}

	var result ApiResponse[json.RawMessage]
	result.Payload = resp.Payload
	result.Code = resp.Code
	result.Message = resp.Message
	result.Route = resp.Route
	result.InternalMessage = resp.InternalMessage
	result.ResponseID = resp.ResponseID

	return &result, nil
}

// Complete completes a pre-authorized payment
func (s *SDK) Complete(req CompleteRequest) error {
	_, err := s.makeRequest("/complete", http.MethodPost, req)
	if err != nil {
		return err
	}

	return nil
}

// AutoPay processes an automatic payment using saved card details
func (s *SDK) AutoPay(req AutoPayRequest) (*ApiResponse[OrderInfo], error) {
	// Apply defaults if values are not provided
	if req.Currency == "" {
		req.Currency = s.defaultCurrency
	}
	if req.CallbackURL == "" {
		req.CallbackURL = s.defaultCallbackURL
	}

	resp, err := s.makeRequest("/autoPay", http.MethodPost, req)
	if err != nil {
		return nil, err
	}

	var result ApiResponse[OrderInfo]
	if err := json.Unmarshal(resp.Payload, &result.Payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order info: %w", err)
	}

	// Copy response metadata
	result.Code = resp.Code
	result.Message = resp.Message
	result.Route = resp.Route
	result.InternalMessage = resp.InternalMessage
	result.ResponseID = resp.ResponseID

	return &result, nil
}

// IsSuccessful checks if an operation was successful based on the response code
func (s *SDK) IsSuccessful(code ResultCode) bool {
	return code == ResultCodeSuccess || code == ResultCodeSuccessGateway
}
