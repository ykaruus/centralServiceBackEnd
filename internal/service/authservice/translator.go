package authservice

import (
	"centralService/internal/domain"
	apperr "centralService/internal/service/errors"
	"centralService/internal/storage"
	"errors"
)

func TranslateUserDomainErrors(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrUserEmailOutDomain):

		return ErrUserOutOfDomain
	case errors.Is(err, domain.ErrUserNotHasTargetRole):
		return domain.ErrUserNotHasTargetRole
	default:

		return apperr.Wrap("USER_DOMAIN_ERROR", "Ocorreu um erro desconhecido ao realizar esta operação.", err)

	}
}

func TranslateStorageErrs(err error) error {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return ErrUserNotFound
	case errors.Is(err, storage.ErrInvalidUserId):
		return ErrUserInvalidId
	case errors.Is(err, storage.ErrDuplicateKey):
		return ErrUserAlreadyExist
	default:
		return apperr.Wrap("USER_STORAGE_ERROR", "Ocorreu um erro desconhecido ao realizar esta operação.", err)
	}
}

func TranslateTokenErr(err error) error {
	switch {
	case errors.Is(err, storage.ErrTokenClaims):
		return ErrTokenClaimsInvalid
	case errors.Is(err, storage.ErrTokenExpired):
		return ErrTokenExpired
	case errors.Is(err, storage.ErrTokenSignatureInvalid):
		return ErrTokenInvalidSignature
	case errors.Is(err, storage.ErrTokenMalformed):
		return ErrTokenMalformed
	case errors.Is(err, storage.ErrTokenDecoding):
		return ErrTokenInvalidIssuer
	case errors.Is(err, storage.ErrTokenNotValid):
		return ErrTokenInvalidSignature

	default:
		return ErrTokenUnknown
	}
}
