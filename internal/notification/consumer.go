package notification

import (
	"backend-bootcamp-assignment-2024/dto"
	"backend-bootcamp-assignment-2024/internal/house"
	"backend-bootcamp-assignment-2024/pkg/sender"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Subscriber interface {
	StartConsuming(ctx context.Context)
}

type subscriber struct {
	channel         *amqp.Channel
	sender          sender.Sender
	houseRepository house.HouseRepository
}

func NewSubscriber(channel *amqp.Channel, sender sender.Sender, houseRepository house.HouseRepository) Subscriber {
	return &subscriber{
		channel:         channel,
		sender:          sender,
		houseRepository: houseRepository,
	}
}

func (s *subscriber) StartConsuming(ctx context.Context) {
	msgs, err := s.channel.Consume(
		"SubscribeNotification",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var flat dto.Flat
			err := json.Unmarshal(d.Body, &flat)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			emails, err := s.houseRepository.GetEmailsByHouseID(ctx, flat.HouseId)
			if err != nil {
				log.Printf("Failed to get emails: %v", err)
				continue
			}

			for _, email := range emails {
				message := fmt.Sprintf("New flat created in house %d: %+v", flat.HouseId, flat)
				err := s.sender.SendEmail(ctx, string(email), message)
				if err != nil {
					log.Printf("Failed to send email: %v", err)
					continue
				}
			}
		}
	}()
	<-forever
}
