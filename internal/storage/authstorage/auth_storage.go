package authstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/interfaces"
	"context"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
)

type AuthStorage struct {
	client_id string
}

func NewAuthStorage(client_id string) *AuthStorage {
	return &AuthStorage{
		client_id: client_id,
	}
}

var _ interfaces.AuthStorageInterface = (*AuthStorage)(nil)

func (a *AuthStorage) ValidateIdToken(context_count context.Context, id_token string) error {
	ctx, cancel := context.WithTimeout(context_count, 10*time.Second)

	defer cancel()
	_, errPayload := idtoken.Validate(ctx, id_token, a.client_id)

	if errPayload != nil {
		return errPayload
	}

	return nil
}

func (a *AuthStorage) ClaimIdToken(context_count context.Context, id_token string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context_count, 10*time.Second)

	defer cancel()
	payload, errPayload := idtoken.Validate(ctx, id_token, a.client_id)

	if errPayload != nil {
		return nil, errPayload
	}

	return &entities.User{
		Name:    payload.Claims["name"].(string),
		Email:   payload.Claims["email"].(string),
		Picture: payload.Claims["picture"].(string),
	}, nil

}
