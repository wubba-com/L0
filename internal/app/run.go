package app

import (
	"github.com/wubba-com/L0/internal/config"
	"log"
)

func Run() {
	var cfg *config.Config

	cfg = config.GetConfig()
	log.Printf(cfg.Listen.Port)
	//r := chi.NewRouter()
}
