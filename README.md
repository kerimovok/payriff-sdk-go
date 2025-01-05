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
import (
	"payriff-sdk-go/pkg/payriff"
)

payriff := payriff.NewSDK()
```

### Custom Configuration

```go
import (
	"payriff-sdk-go/pkg/payriff"
)

payriff := payriff.NewSDK(payriff.Config{
	// optional, defaults to https://api.payriff.com/api/v3
	BaseURL: "https://api.payriff.com/api/v3",
	// optional, defaults to PAYRIFF_SECRET_KEY environment variable
	SecretKey: "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
})
```

## Features

### Create Order

Create a new payment order:

```go
order, err := payriff.CreateOrder(payriff.CreateOrderRequest{
	Amount:      10.99,
	Language:    payriff.LanguageEN,
	Currency:    payriff.CurrencyUSD,
	Description: "Product purchase",
	CallbackURL: "https://example.com/webhook",
	CardSave:    true,
	Operation: payriff.OperationPurchase,
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

```go
autoPay, err := payriff.AutoPay(payriff.AutoPayRequest{
	CardUUID: "CARD_UUID",
	Amount:   10.99,
	Currency: payriff.CurrencyUSD,
	Description: "Subscription renewal",
	CallbackURL: "https://example.com/webhook",
	Operation:   payriff.OperationPurchase,
})
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
