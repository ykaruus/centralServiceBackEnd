package authstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
	"centralService/internal/domain/interfaces"
	"centralService/internal/storage"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenStorage struct {
	secretkey []byte
}

var _ interfaces.TokenStorageInterface = (*TokenStorage)(nil)

func NewTokenStorage(secretkey string) *TokenStorage {
	return &TokenStorage{
		secretkey: []byte(secretkey),
	}
}

func (t *TokenStorage) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return t.secretkey, nil
	})

	if err != nil {
		return storage.ErrTokenNotValid
	}

	if !token.Valid {
		return storage.ErrTokenNotValid
	}

	return nil
}

func (t *TokenStorage) CreateToken(u *entities.User) (string, error) {

	flags := make([]float64, 0, len(u.Roles))
	for _, p := range u.Roles {
		flags = append(flags, float64(p))
	}
	regionFlag := float64(u.AssignedToRegion)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"roles":   flags,
		"region":  regionFlag,
		"exp":     time.Now().Add(time.Hour * 10).Unix(),
	})

	tokenString, err := token.SignedString(t.secretkey)

	if err != nil {
		slog.Error("token-storage.CreateToken failed", "error", err.Error())
		err := TranslateTokenErr(err)
		return "", err
	}

	return tokenString, nil

}

func (t *TokenStorage) DecriptToken(token string) (*entities.User, error) {
	descriptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretkey), nil
	})

	if err != nil {
		slog.Info("token decoding error", "error", err.Error())
		return nil, err
	}

	claims, ok := descriptedToken.Claims.(jwt.MapClaims), descriptedToken.Valid

	if !ok {
		return nil, storage.ErrTokenNotValid
	}

	if !descriptedToken.Valid {
		return nil, storage.ErrTokenNotValid
	}

	slog.Info("token decoding", "claims", claims)

	userPerms := make([]enums.UserPermission, 0)
	var regionEnum enums.RegionFlag
	userId, ok := claims["user_id"].(string)

	if _, ok := claims["user_id"].(string); !ok {
		slog.Info("token decoding error", "error", "user_id not valid")
		return nil, storage.ErrTokenClaims
	}
	if _, ok := claims["region"]; !ok {
		slog.Info("token decoding error", "errors", "region not valid")
		return nil, storage.ErrTokenClaims
	}
	userId = fmt.Sprint(claims["user_id"])
	rolesRaw, ok := claims["roles"].([]interface{})

	if !ok {
		slog.Info("token decoding error", "error", "roles not valid")
		return nil, storage.ErrTokenClaims
	}

	for _, flag := range rolesRaw {

		flagFloat, taok := flag.(float64)

		if !taok {
			slog.Info("The role flag not is valid", "flag", flag)
			return nil, storage.ErrTokenClaims
		}

		userPerms = append(userPerms, enums.UserPermission(int(flagFloat)))

	}

	regionInterface, ok := claims["region"].(float64)
	regionEnum = enums.RegionFlag(int(regionInterface))

	if !ok {
		slog.Info("region claim in token not valid")
		return nil, storage.ErrTokenClaims
	}

	return &entities.User{
		ID:               userId,
		Roles:            userPerms,
		AssignedToRegion: regionEnum,
	}, nil

}

func (t *TokenStorage) ExtractToken(token string) (string, error) {
	parts := strings.SplitN(token, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", storage.ErrTokenMalformed
	}

	return parts[1], nil
}
