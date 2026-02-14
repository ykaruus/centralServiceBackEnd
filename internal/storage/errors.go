package storage

import "errors"

var (
	ErrEquipmentInvalidId    = errors.New("storage:Invalid equipment ID")
	ErrShipmentInvalidID     = errors.New("storage:Invalid shipment ID")
	ErrTimeout               = errors.New("storage: timeout expired")
	ErrNotFound              = errors.New("storage: not found")
	ErrDuplicateKey          = errors.New("storage: duplicated key")
	ErrContextCanceled       = errors.New("storage: context canceled")
	ErrTokenDecoding         = errors.New("storage: error ocurred at decoding token")
	ErrInvalidUserId         = errors.New("storage: invalid user ID")
	ErrTokenMalformed        = errors.New("storage: the token is malformed")
	ErrTokenSignatureInvalid = errors.New("storage: the token has a invalid signature")
	ErrTokenExpired          = errors.New("storage: the token has expired")
	ErrTokenClaims           = errors.New("storage: error at extracting the claims from token")
	ErrTokenNotValid         = errors.New("storage: token not valid")
	ErrTokenInvalidKey       = errors.New("storage: token invalid key")
	ErrTokenInternalError    = errors.New("storage: unfound error ocurred at token operation")
	ErrTokenInvalidRoleFlag  = errors.New("storage: The role claim at token is invalid")
)
