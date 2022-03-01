package nats

import "github.com/nats-io/stan.go"

type FuncNats func(m *stan.Msg)
