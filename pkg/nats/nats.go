package nats

import (
	"github.com/nats-io/stan.go"
	"github.com/wubba-com/L0/internal/app/order/delivery/nats"
	"log"
)

func NewStanConn(clusterID, clientID string) stan.Conn {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("[err] stan-streaming connect: %s\n", err.Error())
	}

	return sc
}

func NewSubscriber(sc stan.Conn, channel string, f nats_order.FuncNats) stan.Subscription {
	sub, err := sc.Subscribe(channel, stan.MsgHandler(f), stan.DurableName("durable-name"), stan.SetManualAckMode())
	if err != nil {
		log.Fatalf("[err] stan-streaming subscription: %s\n", err.Error())
		return nil
	}

	return sub
}

func NewPublisher() {

}
