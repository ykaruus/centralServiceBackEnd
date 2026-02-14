package interfaces

import (
	"centralService/internal/domain/entities"
	"context"
)

type UserStorageInterface interface {
	Create(context.Context, entities.User) (string, error)
	GetByID(context.Context, string) (*entities.User, error)
	GetByEmail(context.Context, string) (*entities.User, error)
	List(context.Context, *entities.UserFilter) ([]entities.User, error)
	Delete(context.Context, string) error
	Update(context.Context, *entities.User) error
	InitIndexes() error
	UpdateNameAndPicture(context.Context, string, string, string) error
}
