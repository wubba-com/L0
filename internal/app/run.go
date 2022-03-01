package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	httpo "github.com/wubba-com/L0/internal/app/order/delivery/http"
	natso "github.com/wubba-com/L0/internal/app/order/delivery/nats"
	"github.com/wubba-com/L0/internal/app/order/repository/postgres"
	order "github.com/wubba-com/L0/internal/app/order/service"
	"github.com/wubba-com/L0/internal/config"
	cache "github.com/wubba-com/L0/pkg/cache"
	"github.com/wubba-com/L0/pkg/client/pg"
	"github.com/wubba-com/L0/pkg/nats"
	"log"
	"net/http"
	"time"
)

func Run() {
	// init of start vars
	var cfg *config.Config
	DefaultExpiration := 5 * time.Minute
	CleanupInterval := 10 * time.Minute
	TTL := 15 * time.Minute
	MaxAttempts := 3

	// init config
	cfg = config.GetConfig()
	log.Printf("init config\n")

	// init http router
	router := chi.NewRouter()
	log.Printf("init http-router")

	// init db client
	log.Println(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	PGSQLClient, err := pg.NewClient(context.TODO(), cfg, MaxAttempts)
	if err != nil {
		log.Fatalf("err: %s", err.Error())
	}

	// init cache
	cacheLocal := cache.NewCache(DefaultExpiration, CleanupInterval)
	log.Printf("init cache")

	// init order of handler, service, repository
	r := postgres.NewOrderRepository(PGSQLClient)
	s := order.NewOrderService(r, cacheLocal, TTL)
	h := httpo.NewOrderHandler(s)

	// init http handlers
	h.Register(router)

	// init nats handler
	n := natso.NewOrderHandler(s)
	sc := nats.NewStanConn(cfg.Nats.ClusterID, cfg.Nats.ClientID)

	// init nats-streaming subscriber
	nats.NewSubscriber(sc, cfg.Nats.Channel, n.StoreOrder)

	// init listen http host
	listen := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	log.Printf(listen)

	// start server
	log.Printf("init server")
	log.Fatal(http.ListenAndServe(listen, router))

}
