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

//To be sent to SQSPayments
type UpdatedPaymentEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

//to be recived from PATCH /payments
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
	Status     string `json:"status"`
}

//To be persisted into DynamoDB 'payments'
type PaymentItem struct {
	PaymentID  string `json:"payment_id"`
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
	Status     string `json:"status"`
}

//To be sent back to the client
type HttpMessageResponse struct {
	Message string `json:"message"`
}

//Order/Payment Status
const (
	OrderStatusIncomplete      = "incomplete"
	PaymentStatusPending       = "pending"
	PaymentStatusSuccess       = "success"
	OrderStatusRedyForShipping = "ready_for_shipping"
	PaymentStatusRejected      = "rejected"
	OrderStatusRejected        = "rejected"
)

const (
	PaymentsDynamoDBTableName = "payments"
	OrdersDynamoDBTableName   = "orders"
	PaymentsSQSQueue          = "SQSPayments"
	OrdersSQSQueue            = "SQSOrders"
)

/* deprected staus enum. May be confusing during tech challenge validation
//Order/Payment Status enum
type Status int

const (
	OrderStatusIncomplete Status = iota
	PaymentStatusPending
	PaymentStatusSuccess
	OrderStatusRedyForShipping
	PaymentStatusRejected
)*/
