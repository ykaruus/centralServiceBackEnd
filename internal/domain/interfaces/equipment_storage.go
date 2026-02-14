package interfaces

import (
	"centralService/internal/domain/entities"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EquipmentStorageInterface interface {
	InitIndexes() error
	Create(context context.Context, equipament entities.Equipment) (string, error)
	GetByID(context context.Context, id string) (*entities.Equipment, error)
	Put(context context.Context, equipament_id string, update bson.M) error
	Patch(context context.Context, equipament_id string, update bson.M) error
	Delete(context context.Context, id string) error
	List(context context.Context, filter entities.EquipmentFilter) ([]entities.Equipment, error)
}
