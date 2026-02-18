package models

import (
	"centralService/internal/controllers"
	"centralService/internal/domain/enums"
	"centralService/internal/http/response"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type EquipmentFilter struct {
	ID                string                `json:"id"`
	Name              string                `json:"name" validate:"max=50,required"`
	Type              enums.EquipmentType   `json:"type" validate:"oneof=1 2 3 4"`
	Serial            string                `json:"serial" validate:"max=50,required"`
	AssignedToRegion  string                `json:"region" validate:"oneof=sudeste centro-oeste sul"`
	Status            enums.EquipmentStatus `json:"status" validate:"oneof=1 2 3 4 5"`
	CurrentShipmentId string                `json:"current_shipment_id" validate:"omitempty"` // hex do Shipment
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

func (e *EquipmentFilter) Validate() error {

	validator := validator.New()

	if err := validator.Struct(validator); err != nil {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{Code: "INVALID_BAD_REQUEST", Message: response.ErrBadRequest.Message, Fields: controllers.ValidationErrors(err)})
	}

	return nil
}
