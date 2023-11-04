package models

import (
	"github.com/go-playground/validator"
)

var Validate *validator.Validate

func InitializeValidators() {
	Validate = validator.New()
}
