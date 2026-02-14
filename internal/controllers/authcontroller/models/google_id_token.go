package models

import (
	"centralService/internal/controllers"
	"centralService/internal/http/response"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Google_ID struct {
	IdToken string `json:"id_token" validate:"required"`
}

func (g *Google_ID) Validate() error {

	validator := validator.New()

	if err := validator.Struct(g); err != nil {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
			Code:    response.ErrBadRequest.Code,
			Message: response.ErrBadRequest.Message,
			Fields:  controllers.ValidationErrors(err),
		})
	}

	return nil
}
