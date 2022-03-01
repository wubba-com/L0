package main

import (
	publisher "github.com/wubba-com/L0/cmd/nats_publisher"
	"github.com/wubba-com/L0/internal/app"
)

func main() {
	app.Run()
	publisher.Run()
}
