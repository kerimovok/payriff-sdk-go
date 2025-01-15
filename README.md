# Unofficial Payriff SDK for Go

An unofficial Go SDK for integrating with the Payriff payment gateway.

## Installation

```bash
go get github.com/kerimovok/payriff-sdk-go
```

## Configuration

Initialize the SDK with your merchant credentials:

### Default Configuration

```go
import "github.com/kerimovok/payriff-sdk-go/payriff"

// Uses environment variables and default values:
// - PAYRIFF_SECRET_KEY for secret key
// - PAYRIFF_CALLBACK_URL for callback URL
// - "AZ" for language
// - "AZN" for currency
// - "https://api.payriff.com/api/v3" for base URL
sdk := payriff.NewSDK(payriff.Config{})
```

### Custom Configuration

```go
import "github.com/kerimovok/payriff-sdk-go/payriff"

sdk := payriff.NewSDK(payriff.Config{
	BaseURL:            "https://api.payriff.com/api/v3",
	SecretKey:          "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
	DefaultCallbackURL: "https://example.com/webhook",
	DefaultLanguage:    payriff.LanguageEN,
	DefaultCurrency:    payriff.CurrencyUSD,
})
```

## Features

### Create Order

Create a new payment order:

#### With defaults

```go
order, err := sdk.CreateOrder(payriff.CreateOrderRequest{
	Amount:      10.99,
	Description: "Product purchase",
	CardSave:    false,
})
```

#### With custom options

```go
order, err := sdk.CreateOrder(payriff.CreateOrderRequest{
    Amount:      10.99,
    Description: "Product purchase",
    CardSave:    false,
    Operation:   payriff.OperationPurchase,
    Language:    payriff.LanguageEN,
    Currency:    payriff.CurrencyUSD,
    CallbackURL: "https://example.com/custom-webhook",
})
```

### Get Order Information

Retrieve details about an existing order:

```go
orderInfo, err := payriff.GetOrderInfo("ORDER_ID")
```

### Process Refund

Refund a completed payment:

```go
refund, err := payriff.Refund(payriff.RefundRequest{
	OrderID: "ORDER_ID",
	Amount:  10.99,
})
```

### Complete Pre-authorized Payment

Complete a pre-authorized payment:

```go
err := payriff.Complete(payriff.CompleteRequest{
	OrderID: "ORDER_ID",
	Amount:  10.99,
})
```

### Automatic Payment

Process payment using saved card details:

#### With defaults

```go
autoPay, err := payriff.AutoPay(payriff.AutoPayRequest{
	CardUUID:    "CARD_UUID",
	Amount:      10.99,
	Description: "Subscription renewal",
})
```

#### With custom options

```go
autoPay, err := payriff.AutoPay(payriff.AutoPayRequest{
	CardUUID:    "CARD_UUID",
	Amount:      10.99,
	Currency:    payriff.CurrencyUSD,
	Description: "Subscription renewal",
	CallbackURL: "https://example.com/webhook",
	Operation:   payriff.OperationPurchase,
})
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
