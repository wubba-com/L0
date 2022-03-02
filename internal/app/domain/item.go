package domain

type Item struct {
	ChrtID      uint64 `json:"chrt_id" validate:"required,gt=0"`
	TrackNumber string `json:"track_number" validate:"required,max=256"`
	Price       int    `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required,max=128"`
	Sale        int    `json:"sale" validate:"required"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required"`
	NmID        uint64 `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required,max=256"`
	Status      int    `json:"status" validate:"required"`
}
