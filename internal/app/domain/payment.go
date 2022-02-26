package domain

type Payment struct {
	Transaction  string `json:"transaction,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	Currency     string `json:"currency,omitempty"`
	Provider     string `json:"provider,omitempty"`
	Amount       uint64 `json:"amount,omitempty"`
	PaymentDt    uint64 `json:"payment_dt,omitempty"`
	Bank         string `json:"bank,omitempty"`
	DeliveryCost uint64 `json:"delivery_cost,omitempty"`
	GoodsTotal   uint64 `json:"goods_total,omitempty"`
	CustomFee    uint64 `json:"custom_fee,omitempty"`
}
