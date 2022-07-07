package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sentrionic/ecommerce/common/events"
	"github.com/sentrionic/ecommerce/expiration/publishers"
)

const TypeOrderExpiration = "order:expiration"

type OrderScheduler interface {
	PublishExpirationTask(orderId uuid.UUID) (*asynq.Task, error)
	HandleExpirationTask(ctx context.Context, t *asynq.Task) error
}

type orderScheduler struct {
	publisher publishers.ExpirationPublisher
}

func NewOrderListener(publisher publishers.ExpirationPublisher) OrderScheduler {
	return &orderScheduler{
		publisher: publisher,
	}
}

func (s *orderScheduler) PublishExpirationTask(orderId uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(events.ExpirationCompleteData{OrderID: orderId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeOrderExpiration, payload), nil
}

func (s *orderScheduler) HandleExpirationTask(_ context.Context, t *asynq.Task) error {
	var p events.ExpirationCompleteData
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	s.publisher.PublishExpirationComplete(p.OrderID)

	return nil
}
