package usercontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/http/response"
	"centralService/internal/service/authservice"
	apperr "centralService/internal/service/errors"
	"context"
	"errors"
	"net/http"
)

func TranslateServiceError(err error) *controllers.HTTPResponse {
	var appErr *apperr.AppError

	if errors.As(err, &appErr) {
		switch {
		case errors.Is(err, authservice.ErrUserAlreadyExist):
			return controllers.Wrap(http.StatusConflict, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		case errors.Is(err, authservice.ErrUserNotFound):
			return controllers.Wrap(http.StatusNotFound, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		case errors.Is(err, authservice.ErrUserOutOfDomain):
			return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
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

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return controllers.Wrap(http.StatusRequestTimeout, &response.ErrForbidden)

	case errors.Is(err, context.Canceled):
		return controllers.Wrap(http.StatusBadRequest, &response.ErrForbidden)
	}

	return controllers.Wrap(http.StatusInternalServerError, &response.ErrInternal)

}
