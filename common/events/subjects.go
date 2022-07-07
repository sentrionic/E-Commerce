package events

type Subject string

const (
	ProductCreated     Subject = "product:created"
	ProductUpdated     Subject = "product:updated"
	OrderCreated       Subject = "order:created"
	OrderCancelled     Subject = "order:cancelled"
	ExpirationComplete Subject = "expiration:complete"
	PaymentCreated     Subject = "payment:created"
)
