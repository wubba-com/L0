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

func NewOrderHandler(orderS domain.OrderService, deliveryS domain.DeliveryService, paymentS domain.PaymentService, itemS domain.ItemService, validater validation.Validater) *handler {
	return &handler{o: orderS, d: deliveryS, p: paymentS, i: itemS, v: validater}
}

type handler struct {
	o domain.OrderService
	d domain.DeliveryService
	p domain.PaymentService
	i domain.ItemService
	v validation.Validater
}

func (h *handler) StoreOrder(msg *stan.Msg) {
	order := &domain.Order{}

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
	uid, err := h.o.StoreOrder(context.Background(), order)
	if err != nil {
		log.Printf("[err] nats handler: %s\n", err.Error())

		err = msg.Ack()
		if err != nil {
			log.Printf("[err] nats handler: %s\n", err.Error())
			return
		}
		return
	}
	_, err = h.d.StoreDelivery(context.Background(), order.Delivery)
	if err != nil {
		log.Printf("[err] nats handler: %s\n", err.Error())

		err = msg.Ack()
		if err != nil {
			log.Printf("[err] nats handler: %s\n", err.Error())
			return
		}
		return
	}
	_, err = h.p.StorePayment(context.Background(), order.Payment)
	if err != nil {
		log.Printf("[err] nats handler: %s\n", err.Error())

		err = msg.Ack()
		if err != nil {
			log.Printf("[err] nats handler: %s\n", err.Error())
			return
		}
		return
	}
	for _, item := range order.Items {
		go func(item *domain.Item) {
			_, err = h.i.StoreItem(context.Background(), item)
			if err != nil {
				log.Printf("[err] nats handler: %s\n", err.Error())

				err = msg.Ack()
				if err != nil {
					log.Printf("[err] nats handler: %s\n", err.Error())
					return
				}
				return
			}
		}(item)
	}
	err = msg.Ack()
	if err != nil {
		log.Printf("[err] failed ACK msg: %d\n", msg.Sequence)
		return
	}
	fmt.Printf("OrderUUID: %s\n", uid)
}
