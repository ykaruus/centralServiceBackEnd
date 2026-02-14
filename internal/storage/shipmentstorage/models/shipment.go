package models

import (
	"centralService/internal/domain/enums"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Mail struct {
	OperationType string    `bson:"operation_type"`
	SentAt        time.Time `bson:"sent_at,omitempty"`
	DeliveredAt   time.Time `bson:"delivered_at,omitempty"`
	PostCode      string    `bson:"post_code"`
	TrackCode     string    `bson:"track_code"`
}

type Shipment struct {
	ID           bson.ObjectID        `bson:"_id"`
	EquipmentIDs []bson.ObjectID      `bson:"equipment_ids"`
	AssignedTo   bson.ObjectID        `bson:"assigned_to,omitempty"`
	Status       enums.ShipmentStatus `bson:"status" json:"status"`
	Mail         Mail                 `bson:"mail" json:"mail" validate:"required"`
	Type         enums.ShipmentType   `bson:"type" json:"type"`
	CreatedAt    time.Time            `bson:"created_at,omitempty"`
	UpdatedAt    time.Time            `bson:"updated_at,omitempty"`
}
