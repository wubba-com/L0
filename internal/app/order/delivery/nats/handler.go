package nats_order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/validation"
	"log"
)

func NewOrderHandler(service domain.OrderService, validater validation.Validater) *handler {
	return &handler{s: service, v: validater}
}

type handler struct {
	s domain.OrderService
	v validation.Validater
}

func (h *handler) StoreOrder(msg *stan.Msg) {
	order := &domain.Order{}
	ctx := context.Background()

	err := json.Unmarshal(msg.Data, order)
	if err != nil {
		log.Printf("[err] nats handler: %s\n", err.Error())

		err = msg.Ack()
		if err != nil {
			log.Printf("[err] nats handler: %s\n", err.Error())
			return
		}
		return
	}

	err = h.v.Struct(order)
	if err != nil {
		log.Printf("[err] validate: %s\n", err.Error())
		err = msg.Ack()
		if err != nil {
			log.Printf("[err] nats handler: %s\n", err.Error())
			return
		}
		return
	}
	uid, err := h.s.StoreOrder(ctx, order)
	if err != nil {
		log.Printf("[err] nats handler: %s\n", err.Error())

		err = msg.Ack()
		if err != nil {
			log.Printf("[err] nats handler: %s\n", err.Error())
			return
		}
		return
	}
	err = msg.Ack()
	if err != nil {
		log.Printf("[err] failed ACK msg: %d\n", msg.Sequence)
		return
	}
	fmt.Printf("OrderUUID: %s\n", uid)
}
