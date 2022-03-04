package http_order

import (
	"github.com/go-chi/chi/v5"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/validation"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func NewOrderHandler(orderS domain.OrderService, deliveryS domain.DeliveryService, paymentS domain.PaymentService, itemS domain.ItemService, validate validation.Validater) Handler {
	return &handlerOrder{o: orderS, d: deliveryS, p: paymentS, i: itemS, v: validate}
}

const (
	endPoint = "/"
	DirTmpl  = "templates"
)

type handlerOrder struct {
	o domain.OrderService
	d domain.DeliveryService
	p domain.PaymentService
	i domain.ItemService
	v validation.Validater
}

func (h *handlerOrder) Register(r chi.Router) {

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/orders", func(r chi.Router) {
				r.Get("/", h.list)
				r.Get("/{order_uid}", h.get)
			})
		})
	})
}

func (h *handlerOrder) list(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(view("index"))
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}

	orders, err := h.o.All(r.Context())
	if err != nil {
		err = tmpl.Execute(w, err)
		if err != nil {
			log.Printf("err http handler:%s\n", err.Error())
			return
		}
		return
	}

	view := &ViewOrder{Len: len(orders), Orders: orders}

	err = tmpl.Execute(w, view)
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}
}

func (h *handlerOrder) get(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")

	tmpl, err := template.ParseFiles(view("detail"))
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}

	order, err := h.o.GetByUID(r.Context(), orderUID)
	if err != nil {
		err = tmpl.Execute(w, err.Error())
		if err != nil {
			log.Printf("err http handler:%s\n", err.Error())
			return
		}
		return
	}

	err = tmpl.Execute(w, order)
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}
}

type ViewOrder struct {
	Len    int
	Orders []*domain.Order
}

func view(name string) string {
	wd, _ := os.Getwd()
	ext := ".html"
	return filepath.Join(wd, DirTmpl, name+ext)
}
