package consumer

import (
	"encoding/json"
	"fmt"
	"jokegen/internal/models"
	"jokegen/internal/service"
	"log"

	"github.com/nats-io/nats.go"
)

type Publisher interface {
	AddAttempts(userID, count int) error
}
type Publish struct {
	Publ Publisher
}

func NewPublish(users *service.UsersService) *Publish {
	return &Publish{Publ: users}
}

func (p *Publish)SubscribeNats(nc *nats.Conn) error {
	_, err := nc.Subscribe("joker", func(m *nats.Msg) {
		log.Println(string(m.Data))
		pl := &models.Payload{}
		json.Unmarshal(m.Data, pl)
		p.addNewAttempts(pl)
	})
	if err != nil {
		return fmt.Errorf("error during reading topic: %w", err)
	}
	log.Println("start waiting for messages")
	return nil
}

func (p *Publish)addNewAttempts(payload *models.Payload) {
	newAttempts := int(payload.Price / 2)
	err := p.Publ.AddAttempts(payload.ServiceID, newAttempts)
	if err != nil {
		log.Printf("new attempts wasn't added because of this error: %v", err)
	}
	log.Printf("for user %v was added %v attempts", payload.ServiceID, newAttempts)
}
