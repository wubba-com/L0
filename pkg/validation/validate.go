package validation

import (
	"github.com/go-playground/validator/v10"
)

func NewValidater() Validater {
	v := validator.New()
	return &validater{v}
}

type Validater interface {
	Struct(interface{}) error
}

type validater struct {
	v *validator.Validate
}

func (v validater) Struct(s interface{}) error {
	err := v.v.Struct(s)
	if err != nil {
		// from here you can create your own error messages in whatever language you wish
		return err
	}
	return nil
}
