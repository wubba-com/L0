package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/wubba-com/L0/internal/app/domain"
	"log"
)

type Handler struct {
	s domain.OrderService
}

func (h *Handler) StoreOrder(m *stan.Msg) {
	order := &domain.Order{}
	ctx := context.Background()

	err := json.Unmarshal(m.Data, order)
	if err != nil {
		log.Printf("err: %s", err.Error())

		err = m.Ack()
		if err != nil {
			log.Printf("err: %s", err.Error())
			return
		}
		return
	}

	uid, err := h.s.StoreOrder(ctx, order)
	if err != nil {
		log.Printf("err: %s", err.Error())

		err = m.Ack()
		if err != nil {
			log.Printf("err: %s", err.Error())
			return
		}
		return
	}

	log.Printf("OrderUUID: %s", uid)
}


