package authservice

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/interfaces"
	"centralService/internal/infra/trace"
	"context"
)

type TokenService struct {
	storage interfaces.TokenStorageInterface
}

func NewTokenService(storage interfaces.TokenStorageInterface) *TokenService {
	return &TokenService{
		storage: storage,
	}
}

func (ts *TokenService) VerifyToken(context context.Context, tokenString string) error {

	log := trace.LogWithTraceID("token-service", context)

	log.Info("token-service.VerifyToken calling")
	err := ts.storage.VerifyToken(tokenString)

	if err != nil {
		log.Error("token-service.VerifyToken failed", "error", err.Error())
		err := TranslateTokenErr(err)
		return err
	}

	return nil
}
func (ts *TokenService) CreateToken(context context.Context, u *entities.User) (string, error) {

	log := trace.LogWithTraceID("token-service", context)

	log.Info("token-service.CreateToken calling")

	token, err := ts.storage.CreateToken(u)

	if err != nil {
		log.Error("token-service.CreateToken failed", "error", err)
		err := TranslateTokenErr(err)
		return "", err
	}

	return token, nil

}

func (ts *TokenService) DecriptToken(context context.Context, token string) (*entities.User, error) {
	log := trace.LogWithTraceID("token-service", context)

	log.Info("token-service.DecriptToken calling")

	claims, err := ts.storage.DecriptToken(token)

	if err != nil {
		log.Error("token-service.DecriptToken failed", "error", err)
		err := TranslateTokenErr(err)
		return nil, err
	}

	return claims, nil
}

func (ts *TokenService) ExtractToken(token string) (string, error) {

	token, err := ts.storage.ExtractToken(token)

	if err != nil {
		return "", ErrTokenMalformed
	}

	return token, nil
}
