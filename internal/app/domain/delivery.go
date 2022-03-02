package domain

type Delivery struct {
	OrderID int
	Name    string `json:"name,omitempty" validate:"required,max=128"`
	Phone   string `json:"phone,omitempty" validate:"required,max=32"`
	Zip     string `json:"zip,omitempty" validate:"required,max=128"`
	City    string `json:"city,omitempty" validate:"required,max=128"`
	Address string `json:"address,omitempty" validate:"required,max=128"`
	Region  string `json:"region,omitempty" validate:"required,max=128"`
	Email   string `json:"email,omitempty" validate:"required,max=128,email"`
}
