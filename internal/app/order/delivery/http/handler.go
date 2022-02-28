package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/internal/app/order/delivery"
	"net/http"
)

func NewHandlerOrder(service domain.OrderService) delivery.Handler {
	return &handlerOrder{service}
}

const (
	endPoint = "/"
)

type handlerOrder struct {
	s domain.OrderService
}

func (h *handlerOrder) Register(r chi.Router) {
	r.Route("api", func(r chi.Router) {
		r.Route("v1", func(r chi.Router) {
			r.Route("/orders", func(r chi.Router) {
				r.Get("{order_uid:[a-z]}", h.get)
				r.Post(endPoint, h.store)
			})
		})
	})
}

func (h *handlerOrder) get(w http.ResponseWriter, r *http.Request) {
	uidOrder := chi.URLParam(r, "order_uid")
	order, err := h.s.GetByUID(r.Context(), uidOrder)
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		return
	}
}

func (h *handlerOrder) store(w http.ResponseWriter, r *http.Request) {
	order := &domain.Order{}
	uid, err := h.s.StoreOrder(r.Context(), order)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(uid)
	if err != nil {
		return
	}
}
