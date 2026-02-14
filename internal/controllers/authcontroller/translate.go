package authcontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/domain"
	"centralService/internal/http/response"
	"context"

	"centralService/internal/service/authservice"
	apperr "centralService/internal/service/errors"
	"errors"
	"net/http"
)

func TranslateServiceError(err error) *controllers.HTTPResponse {
	var appErr *apperr.AppError

	// Erros de domínio (AppError)
	if errors.As(err, &appErr) {
		switch {
		// ===== USER / AUTH =====
		case errors.Is(err, authservice.ErrUserNotFound):
			return controllers.Wrap(http.StatusUnauthorized, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		case errors.Is(err, authservice.ErrUserAlreadyExist):
			return controllers.Wrap(http.StatusConflict, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, authservice.ErrUserOutOfDomain):
			return controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, authservice.ErrUserNoHasTargetRole):
			return controllers.Wrap(http.StatusForbidden, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		case errors.Is(err, domain.ErrUserNotHasTargetRole):
			return controllers.Wrap(http.StatusForbidden, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		// ===== TOKEN =====
		case errors.Is(err, authservice.ErrTokenMissing):
			return controllers.Wrap(http.StatusUnauthorized, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, authservice.ErrTokenMalformed):
			return controllers.Wrap(http.StatusUnauthorized, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		case errors.Is(err, authservice.ErrTokenExpired):
			return controllers.Wrap(http.StatusUnauthorized, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, authservice.ErrTokenInvalidSignature):
			return controllers.Wrap(http.StatusUnauthorized, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err, authservice.ErrTokenClaimsInvalid):
			return controllers.Wrap(http.StatusUnauthorized, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		case errors.Is(err,
			authservice.ErrTokenInvalidIssuer), errors.Is(err, authservice.ErrTokenInvalidAlgorithm):
			return controllers.Wrap(http.StatusForbidden, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		case errors.Is(err, authservice.ErrTokenInvalidAlgorithm):
			return controllers.Wrap(http.StatusForbidden, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})

		// ===== FALLBACK DE DOMÍNIO =====
		default:
			return controllers.Wrap(http.StatusInternalServerError, &response.ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
		}
	}

	// ===== ERROS DE CONTEXTO / INFRA =====
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return controllers.Wrap(http.StatusRequestTimeout, &response.ErrForbidden)

	case errors.Is(err, context.Canceled):
		return controllers.Wrap(http.StatusBadRequest, &response.ErrForbidden)
	}

	// ===== FALLBACK FINAL =====
	return controllers.Wrap(http.StatusInternalServerError, &response.ErrInternal)
}
