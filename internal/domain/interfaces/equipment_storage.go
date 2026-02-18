package interfaces

import (
	"centralService/internal/domain/entities"
	"context"
)

type EquipmentStorageInterface interface {
	InitIndexes() error
	Create(context context.Context, equipament entities.Equipment) (string, error)
	GetByID(context context.Context, id string) (*entities.Equipment, error)
	UpdateByFilter(context.Context, *entities.EquipmentFilter)
	Replace(context.Context, *entities.Equipment)
	Delete(context context.Context, id string) error
	List(context context.Context, filter entities.EquipmentFilter) ([]entities.Equipment, error)
}
