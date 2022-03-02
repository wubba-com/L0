package http_order

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/validation"
	"log"
	"net/http"
)

func NewOrderHandler(service domain.OrderService, validate validation.Validater) Handler {
	return &handlerOrder{service, validate}
}

const (
	endPoint = "/"
)

type handlerOrder struct {
	s domain.OrderService
	v validation.Validater
}

func (h *handlerOrder) Register(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/orders", func(r chi.Router) {
				r.Get("/{order_uid:[a-z]}", h.get)
				r.Post(endPoint, h.store)
			})
		})
	})
}

func (h *handlerOrder) get(w http.ResponseWriter, r *http.Request) {
	uidOrder := chi.URLParam(r, "order_uid")
	order, err := h.s.GetByUID(r.Context(), uidOrder)
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}
}

func (h *handlerOrder) store(w http.ResponseWriter, r *http.Request) {
	order := &domain.Order{}
	err := json.NewDecoder(r.Body).Decode(order)
	defer r.Body.Close()
	if err != nil {
		log.Printf("[err] http handler:%s\n", err.Error())
		return
	}

	err = h.v.Struct(order)
	if err != nil {
		log.Printf("[err] http handler:%s\n", err.Error())
		return
	}
	uid, err := h.s.StoreOrder(r.Context(), order)
	if err != nil {
		log.Printf("[err] http handler:%s\n", err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(uid)
	if err != nil {
		log.Printf("[err] http handler:%s\n", err.Error())
		return
	}
}
