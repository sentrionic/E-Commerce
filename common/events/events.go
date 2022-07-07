package events

import (
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common/order"
	"time"
)

type Event struct {
	Subject Subject
	Data    any
}

type ProductData struct {
	ID      uuid.UUID
	Version uint
	Title   string
	Price   int
	UserId  uuid.UUID
}

func ProductCreatedEvent(data ProductData) Event {
	return Event{
		Subject: ProductCreated,
		Data:    data,
	}
}

func ProductUpdatedEvent(data ProductData) Event {
	return Event{
		Subject: ProductUpdated,
		Data:    data,
	}
}

type OrderCreatedData struct {
	ID        uuid.UUID
	Version   uint
	Status    order.Status
	UserId    uuid.UUID
	ExpiresAt time.Time
	Product   struct {
		ID    uuid.UUID
		Price int
	}
}

func OrderCreatedEvent(data OrderCreatedData) Event {
	return Event{
		Subject: OrderCreated,
		Data:    data,
	}
}

type OrderCancelledData struct {
	ID        uuid.UUID
	Version   uint
	ProductID uuid.UUID
}

func OrderCancelledEvent(data OrderCancelledData) Event {
	return Event{
		Subject: OrderCancelled,
		Data:    data,
	}
}

type ExpirationCompleteData struct {
	OrderID uuid.UUID
}

func ExpirationCompleteEvent(data ExpirationCompleteData) Event {
	return Event{
		Subject: ExpirationComplete,
		Data:    data,
	}
}

type PaymentCreatedData struct {
	ID       uuid.UUID
	OrderID  uuid.UUID
	StripeID string
}

func PaymentCreatedEvent(data PaymentCreatedData) Event {
	return Event{
		Subject: PaymentCreated,
		Data:    data,
	}
}
