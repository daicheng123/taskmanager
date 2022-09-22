package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

const (
	TypeWelcomeEmail = "email:welcome"
	//TypeReminderEmail = "email:reminder"
)

// Task payload for any email related tasks.
type emailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

func NewWelcomeEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, payload), nil
}

//func NewReminderEmailTask(id int) (*asynq.Task, error) {
//	payload, err := json.Marshal(emailTaskPayload{UserID: id})
//	if err != nil {
//		return nil, err
//	}
//
//	return asynq.NewTask(TypeReminderEmail, payload), nil
//}

type EmailProcessor struct {
	payLoad *emailTaskPayload
}

func (ep *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	ep.payLoad = new(emailTaskPayload)
	if err := json.Unmarshal(t.Payload(), ep.payLoad); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	return nil
}

func NewEmailProcessor() *EmailProcessor {
	return &EmailProcessor{}
}
