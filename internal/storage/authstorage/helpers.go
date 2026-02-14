package authstorage

import (
	"centralService/internal/storage"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func TranslateTokenErr(err error) error {

	switch {
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return storage.ErrTokenClaims
	case errors.Is(err, jwt.ErrTokenMalformed):
		return storage.ErrTokenMalformed
	case errors.Is(err, jwt.ErrTokenExpired):
		return storage.ErrTokenExpired
	case errors.Is(err, storage.ErrTokenSignatureInvalid):
		return storage.ErrTokenSignatureInvalid
	case errors.Is(err, storage.ErrTokenInvalidKey):
		return storage.ErrTokenInvalidKey
	default:
		return storage.ErrTokenInternalError
	}
}
