package equipmentcontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/http/response"
	"centralService/internal/service/equipmentservice"
	apperr "centralService/internal/service/errors"
	"context"
	"errors"
	"net/http"
)

func TranslateServiceError(err error) *controllers.HTTPResponse {
	var appErr *apperr.AppError

	// Se for AppError, extraímos uma vez
	if errors.As(err, &appErr) {
		switch {
		case errors.Is(err, equipmentservice.ErrEquipmentNotFound):
			return controllers.Wrap(http.StatusNotFound, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, equipmentservice.ErrEquipmentAlreadyExists):
			return controllers.Wrap(http.StatusConflict, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, equipmentservice.ErrEquipmentInvalidId):
			return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, equipmentservice.ErrEquipmentHasActiveShipment):
			return controllers.Wrap(http.StatusConflict, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		default:
			return controllers.Wrap(http.StatusInternalServerError, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		}

	}

	// Erros de contexto (infra)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return controllers.Wrap(http.StatusRequestTimeout, &response.ErrForbidden)

	case errors.Is(err, context.Canceled):
		return controllers.Wrap(http.StatusBadRequest, &response.ErrForbidden)
	}

	// Fallback obrigatório
	return controllers.Wrap(http.StatusInternalServerError, &response.ErrInternal)
}
