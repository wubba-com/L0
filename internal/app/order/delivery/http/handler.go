package http_order

import (
	"encoding/json"
	"fmt"
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
	DirTmpl = "templates"
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
				r.Post(endPoint, h.store)
			})
		})
	})
}

func (h *handlerOrder) list(w http.ResponseWriter, r *http.Request) {
	orders, err := h.o.AllOrders(r.Context())
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}

	tmpl, err := template.ParseFiles(view("index"))
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}

	err = tmpl.Execute(w, orders)
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}
}

func (h *handlerOrder) get(w http.ResponseWriter, r *http.Request) {
	uidOrder := chi.URLParam(r, "order_uid")
	fmt.Println("uidOrder", uidOrder)
	order, err := h.o.GetByUID(r.Context(), uidOrder)
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}
	tmpl, err := template.ParseFiles(view("detail"))
	if err != nil {
		log.Printf("err http handler:%s\n", err.Error())
		return
	}

	err = tmpl.Execute(w, order)
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
	uid, err := h.o.StoreOrder(r.Context(), order)
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

func view(name string) string {
	wd, _ := os.Getwd()
	ext := ".html"
	return filepath.Join(wd, DirTmpl, name+ext)
}