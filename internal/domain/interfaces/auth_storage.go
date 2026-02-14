package interfaces

import (
	"centralService/internal/domain/entities"
	"context"
)

type TokenStorageInterface interface {
	VerifyToken(string) error
	CreateToken(*entities.User) (string, error)
	DecriptToken(string) (*entities.User, error)
	ExtractToken(string) (string, error)
}

type AuthStorageInterface interface {
	ValidateIdToken(context.Context, string) error
	ClaimIdToken(context.Context, string) (*entities.User, error)
}
