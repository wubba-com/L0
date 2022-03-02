package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	httpo "github.com/wubba-com/L0/internal/app/order/delivery/http"
	natso "github.com/wubba-com/L0/internal/app/order/delivery/nats"
	_ "github.com/wubba-com/L0/internal/app/order/repository/file"
	"github.com/wubba-com/L0/internal/app/order/repository/postgres"
	order "github.com/wubba-com/L0/internal/app/order/service"
	"github.com/wubba-com/L0/internal/config"
	cache "github.com/wubba-com/L0/pkg/cache"
	"github.com/wubba-com/L0/pkg/client/pg"
	"github.com/wubba-com/L0/pkg/nats"
	"github.com/wubba-com/L0/pkg/validation"
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
	PGSQLClient, err := pg.NewClient(context.TODO(), cfg, MaxAttempts)
	if err != nil {
		log.Fatalf("err: %s", err.Error())
	}
	log.Printf("init db: %s", "http://"+cfg.DB.Host+":"+cfg.DB.Port)

	// init cache
	cacheLocal := cache.NewCache(DefaultExpiration, CleanupInterval)
	log.Printf("init cache")

	v := validation.NewValidater()

	// init order of handler, service, repository
	r := postgres.NewOrderRepository(PGSQLClient)
	//r := file.NewOrderRepository()
	s := order.NewOrderService(r, cacheLocal, TTL)
	//err := s.LoadOrderCache(context.Background())
	//if err != nil {
	//	log.Printf("err run load cache: %s", err.Error())
	//	return
	//}

	h := httpo.NewOrderHandler(s, v)

	// init http handlers
	h.Register(router)

	// init nats handler
	n := natso.NewOrderHandler(s, v)

	// init nats-streaming connection
	sc := nats.NewStanConn(cfg.Nats.ClusterID, cfg.Nats.ClientID)
	// close nats connect
	defer sc.Close()

	// init nats-streaming subscriber
	subscription := nats.NewSubscriber(sc, cfg.Nats.Channel, n.StoreOrder)
	// unsubscribe removes subscription
	defer subscription.Unsubscribe()

	// init listen http host
	listen := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	log.Printf("http://%s", listen)

	// start server
	log.Printf("init server")
	log.Fatal(http.ListenAndServe(listen, router))
}
