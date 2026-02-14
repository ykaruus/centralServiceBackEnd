package entities

import (
	"centralService/internal/domain"
	"centralService/internal/domain/enums"
	"time"
)

type EquipmentFilter struct {
	ID         string
	ExternalId string
	Name       string
	Type       enums.EquipmentType
	Serial     string
	Status     enums.EquipmentStatus // em estoque | enviado | manutenção
	ShipmentId string
	Region     enums.RegionFlag // hex do Shipment
}

type EquipmentHistory struct {
	ShipmentID string
	AssignedAt time.Time
}
type Equipment struct {
	ID                string
	Name              string
	Type              enums.EquipmentType
	Serial            string
	Status            enums.EquipmentStatus
	EquipmentsHistory []EquipmentHistory
	CurrentShipmentId string
	AssignedToRegion  enums.RegionFlag
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (e *Equipment) EquipmentCanBeShipped() error {
	if e.CurrentShipmentId != "" || e.Status != enums.EQUIPMENT_STATUS_AVAILABLE {
		return domain.ErrEquipmentNotAvailable
	}

	return nil
}
