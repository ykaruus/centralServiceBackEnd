package models

import (
	"centralService/internal/controllers"
	"centralService/internal/http/response"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type Role struct {
	RoleName string `json:"role" validate:"required,oneof=admin, coordinator, technical"`
}

type User struct {
	Name             string    `json:"name" validate:"omitempty,min=3"`
	Email            string    `json:"email" validate:"required,email"`
	Roles            []Role    `json:"roles" validate:"required"`
	AssignedToRegion string    `json:"region" validate:"oneof=centro-oeste sul sudeste"`
	Picture          string    `json:"picture" validate:"omitempty,min=3"`
	ID               string    `json:"user_id,omitempty" validate:"omitempty"`
	LastAccessAt     time.Time `json:"lastAccessAt" validate:"omitempty"`
	CreatedAt        time.Time `json:"createdAt"  validate:"omitempty"`
	UpdatedAt        time.Time `json:"updatedAt"  validate:"omitempty"`
}

func (u *User) Validate() error {

	validator := validator.New()

	if err := validator.Struct(u); err != nil {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
			Code:    response.ErrBadRequest.Code,
			Message: response.ErrBadRequest.Message,
			Fields:  controllers.ValidationErrors(err),
		})
	}

	if len(u.Roles) == 0 {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
			Code:    response.ErrBadRequest.Code,
			Message: response.ErrBadRequest.Message,
			Fields: []string{
				"O usuário precisa estar associado a algum cargo",
			},
		})
	}

	return nil
}
