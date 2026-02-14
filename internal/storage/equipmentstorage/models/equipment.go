package models

import (
	"centralService/internal/domain/enums"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EquipmentHistory struct {
	ShipmentID bson.ObjectID `bson:"shipment_id"`
	AssignedAt time.Time     `bson:"assigned_at"`
}

type Equipment struct {
	ID                bson.ObjectID         `bson:"_id,omitempty"`
	Name              string                `bson:"name"`
	Type              enums.EquipmentType   `bson:"type"`
	Serial            string                `bson:"serial"`
	Status            enums.EquipmentStatus `bson:"status"`
	AssignedToRegion  enums.RegionFlag      `bson:"region_flag"`
	EquipmentsHistory []EquipmentHistory    `bson:"equipments_history"`            // em estoque | enviado | manutenção
	CurrentShipmentId *bson.ObjectID        `bson:"current_shipment_id,omitempty"` // hex do Shipment
	CreatedAt         time.Time             `bson:"created_at"`
	UpdatedAt         time.Time             `bson:"updated_at"`
}
