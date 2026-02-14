package interfaces

import (
	"centralService/internal/domain/entities"
	"context"
)

type ShipmentStorageInterface interface {
	Create(context.Context, entities.Shipment) (string, error)
	GetByID(context.Context, string) (*entities.Shipment, error)
}
