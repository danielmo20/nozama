package nozama

//To be recived from POST /orders/
type CreateOrderRequest struct {
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

//To be sent to SQSPayments
type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

//to be recived from POST /payments
type ProcessPaymentRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

//To be persisted into DynamoDB 'orders'
type OrderItem struct {
	OrderID    string `json:"order_id"`
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
	Status     Status `json:"status"`
}

//To be persisted into DynamoDB 'payments'
type PaymentItem struct {
	PaymentID  string `json:"payment_id"`
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
	Status     Status `json:"status"`
}

//Order/Payment Status enum
type Status int

const (
	OrderStatusIncomplete Status = iota
	PaymentStatusPending
	PaymentStatusSuccess
	OrderStatusRedyForShipping
	PaymentStatusRejected
)
