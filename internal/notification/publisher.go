package notification

import (
	"backend-bootcamp-assignment-2024/dto"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type Publisher interface {
	NewFlatNotification(flat dto.Flat) error
}

type publisher struct {
	channel *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) (Publisher, error) {
	return &publisher{channel: ch}, nil
}

func (p *publisher) NewFlatNotification(flat dto.Flat) error {
	reqJson, err := json.Marshal(flat)
	if err != nil {
		return fmt.Errorf("failed to marshal flat to JSON: %w", err)
	}
	return p.channel.Publish(
		"NotificationQueue",
		"FlatCreated",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        reqJson,
		},
	)
}
