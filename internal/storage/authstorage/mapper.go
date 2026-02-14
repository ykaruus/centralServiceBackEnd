package authstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/storage"
	"centralService/internal/storage/authstorage/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ModelUserToEntityUser(mu *models.User) entities.User {

	return entities.User{
		ID:               mu.ID.Hex(),
		Name:             mu.Name,
		Email:            mu.Email,
		IsFirstLogin:     mu.IsFirstLogin,
		AssignedToRegion: mu.AssignedToRegion,
		Roles:            mu.Role,
		Picture:          mu.Picture,
		CreatedAt:        mu.CreatedAt,
		UpdatedAt:        mu.UpdatedAt,
	}

}

func EntityUserToModelUser(eu *entities.User) (*models.User, error) {
	mu := models.User{

		Name:             eu.Name,
		Email:            eu.Email,
		IsFirstLogin:     eu.IsFirstLogin,
		AssignedToRegion: eu.AssignedToRegion,
		Role:             eu.Roles,
		Picture:          eu.Picture,
		CreatedAt:        eu.CreatedAt,
		UpdatedAt:        eu.UpdatedAt,
	}

	if eu.ID != "" {
		userId, err := bson.ObjectIDFromHex(eu.ID)

		if err != nil {
			return nil, storage.ErrInvalidUserId
		}

		mu.ID = userId

	}

	return &mu, nil
}
