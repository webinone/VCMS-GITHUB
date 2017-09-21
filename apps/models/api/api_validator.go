package api

import "gopkg.in/go-playground/validator.v9"

type APIValidator struct {
	Validator *validator.Validate
}

func (av *APIValidator) Validate(i interface{}) error {
	return av.Validator.Struct(i)
}