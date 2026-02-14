package models

import (
	"centralService/internal/controllers"
	"centralService/internal/domain/enums"
	"centralService/internal/http/response"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type EquipmentFilter struct {
	ID         string                `query:"id" validate:"len=50,omitempty"`
	Name       string                `query:"name" validate:"max=50, omitempty"`
	Type       enums.EquipmentType   `query:"type" validate:"oneof=1 2 3 4 5 6 7, omitempty"`
	Serial     string                `query:"serial" validate:"max=50, omitempty"`
	Status     enums.EquipmentStatus `query:"status" validate:"oneof=1 2 3 4 5, omitempty"` // em estoque | enviado | manutenção
	ShipmentId string                `query:"shipment_id" validate:"omitempty"`             // hex do Shipment
}

func (e *EquipmentFilter) Validate() error {

	validator := validator.New()

	if err := validator.Struct(validator); err != nil {
		return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{Code: "INVALID_BAD_REQUEST", Message: response.ErrBadRequest.Message, Fields: controllers.ValidationErrors(err)})
	}

	return nil
}
