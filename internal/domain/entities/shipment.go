package entities

import (
	"centralService/internal/domain/enums"
	"time"
)

type Mail struct {
	OperationType string
	SentAt        time.Time
	DeliveredAt   time.Time
	PostCode      string
	TrackCode     string
}

type Shipment struct {
	ID           string
	EquipmentIDs []string
	AssignedTo   string
	Status       enums.ShipmentStatus
	Mail         Mail
	Type         enums.ShipmentType
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
