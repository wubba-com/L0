package nats_order

import "github.com/nats-io/stan.go"

type FuncNats func(m *stan.Msg)
