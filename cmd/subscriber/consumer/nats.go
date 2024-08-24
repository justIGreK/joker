package consumer

import (
	"encoding/json"
	"fmt"
	"jokegen/internal/models"
	"jokegen/internal/service"
	"log"

	"github.com/nats-io/nats.go"
)



func SubscribeNats(nc *nats.Conn, s *service.Publish) error {
	_, err := nc.Subscribe("joker", func(m *nats.Msg) {
		log.Println(string(m.Data))
		pl := &models.Payload{}
		json.Unmarshal(m.Data, pl)
		addNewAttempts(pl, s)
	})
	if err != nil {
		return fmt.Errorf("error during reading topic: %w", err)
	}
	log.Println("start waiting for messages")
	return nil
}

func addNewAttempts(payload *models.Payload, s *service.Publish) {
	newAttempts := int(payload.Price / 2)
	err := s.AddAttempts(payload.UserID, newAttempts)
	if err != nil {
		log.Printf("new attempts wasn't added because of this error: %v", err)
	}
	log.Printf("for user %v was added %v attempts", payload.UserID, newAttempts)
}
