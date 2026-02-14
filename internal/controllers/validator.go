package controllers

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationErrors(validatorErr error) []string {
	var errs validator.ValidationErrors
	var messages = make([]string, 0)
	if errors.As(validatorErr, &errs) {

		for _, e := range errs {
			messages = append(messages, fmt.Sprintf("O %s está inválido (%s)", e.Field(), e.Tag()))

		}

	}

	return messages

}
