package models

import (
	"centralService/internal/controllers"
	"centralService/internal/http/response"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserFilter struct {
	Name  string `validate:"max=150"`
	Email string `validate:"email,max=150"`
	ID    string `validate:"max=50"`
	Roles []Role `validate:"dive"`
}

func (uf *UserFilter) Validate() error {
	validator := validator.New()

	if err := validator.Struct(uf); err != nil {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
			Code:    response.ErrBadRequest.Code,
			Message: response.ErrBadRequest.Message,
			Fields:  controllers.ValidationErrors(err),
		})
	}

	if len(uf.Roles) > 3 {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
			Code:    response.ErrBadRequest.Code,
			Message: response.ErrBadRequest.Message,
			Fields: []string{
				"O numero de argumentos excedeu o limite",
			},
		})
	}

	return nil
}
