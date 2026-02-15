package authservice

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/interfaces"
	"centralService/internal/infra/trace"
	"context"
	"time"
)

type AuthService struct {
	user  *UserService
	auth  interfaces.AuthStorageInterface
	token *TokenService
}

func NewAuthService(
	user *UserService,
	auth interfaces.AuthStorageInterface,
	token *TokenService,
) *AuthService {
	return &AuthService{
		user:  user,
		auth:  auth,
		token: token,
	}
}

// procura o email do user na db, senão achar retorna err se achar retorna o token

func (as *AuthService) ValidateIdToken(ctx context.Context, token string) error {
	log := trace.LogWithTraceID("auth-service", ctx)

	log.Info("auth-service.ValidateToken started")

	err := as.auth.ValidateIdToken(ctx, token)

	if err != nil {
		log.Error("auth-service.ValidateToken started")
		return TranslateTokenErr(err)
	}

	return nil
}

func (as *AuthService) DecriptIdToken(ctx context.Context, token string) (*entities.User, error) {
	log := trace.LogWithTraceID("auth-service", ctx)

	log.Info("auth-service.DecriptToken started")

	payload, err := as.auth.ClaimIdToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (as *AuthService) Login(ctx context.Context, u *entities.User) (string, error) {

	log := trace.LogWithTraceID("auth-service", ctx)

	log.Info("auth-service.Login started")
	if err := u.Validate(); err != nil {

		err := TranslateUserDomainErrors(err)

		return "", err

	}

	log.Info("auth-service.Login started searching email", "email", u.Email)

	user, err := as.user.GetByEmail(ctx, u.Email)

	if err != nil {
		log.Error("auth-service.Login failed", "error", err.Error())

		return "", err
	}

	now := time.Now()
	log.Info("auth-service.Login update user", "IsFirstLogin?", user.IsFirstLogin)
	if user.LastAccessAt != now || user.Picture == "" || user.Name == "" {
		log.Info("auth-service.Login updating user")

		update := &entities.UserFilter{
			Name:         u.Name,
			Picture:      u.Picture,
			LastAccessAt: &now,
		}

		err := as.user.UpdateByFilter(ctx, user.ID, update)

		if err != nil {
			log.Error("auth-service.Login update user failed", "error", err.Error())
		}

	}

	token, err := as.token.CreateToken(ctx, user)

	if err != nil {

		log.Error("auth-service.Login failed", "error", err.Error())

		return "", err

	}
	return token, nil

}

func (as *AuthService) GetUserFromToken(ctx context.Context, token string) (*entities.User, error) {
	log := trace.LogWithTraceID("auth-service", ctx)

	log.Info("auth-service.GetUserFromToken started")

	tokenExtracted, err := as.token.ExtractToken(token)

	if err != nil {
		log.Error("auth-service.GetUserFromToken failed", "error", err.Error())
		return nil, err
	}

	log.Info("token extracted", "token", tokenExtracted)

	err = as.token.VerifyToken(ctx, tokenExtracted)

	if err != nil {
		log.Error("auth-service.GetUserFromToken failed", "error", err.Error())
		return nil, err
	}

	payload, err := as.token.DecriptToken(ctx, tokenExtracted)

	if err != nil {
		log.Error("auth-service.GetUserFromToken failed", "error", err.Error())
		return nil, err
	}

	user, err := as.user.GetByID(ctx, payload.ID)

	if err != nil {
		log.Error("auth-service.GetUserFromToken failed", "error", err.Error())
		return nil, err
	}

	return user, nil

}
