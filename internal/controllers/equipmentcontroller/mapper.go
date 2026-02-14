package equipmentcontroller

import (
	"centralService/internal/controllers/equipmentcontroller/models"
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
)

func EntityToEntityRequest(equipment *entities.Equipment) models.Equipment {
	history := make([]models.EquipmentHistory, 0, len(equipment.EquipmentsHistory))
	for _, eh := range equipment.EquipmentsHistory {

		history = append(history, models.EquipmentHistory(eh))
	}

	regionFlag := enums.MapperRegionFlags[equipment.AssignedToRegion]

	return models.Equipment{
		AssignedToRegion:  regionFlag,
		ID:                equipment.ID,
		Name:              equipment.Name,
		Serial:            equipment.Serial,
		CurrentShipmentId: equipment.CurrentShipmentId,
		Type:              equipment.Type,
		Status:            equipment.Status,
		EquipmentsHistory: history,
		CreatedAt:         equipment.CreatedAt,
		UpdatedAt:         equipment.UpdatedAt,
	}
}

func EntitiesToEntitiesRequest(equipments []entities.Equipment) []models.Equipment {

	em := make([]models.Equipment, 0, len(equipments))

	if len(equipments) == 0 {
		return em
	}

	for _, e := range equipments {
		mod := EntityToEntityRequest(&e)
		em = append(em, mod)
	}

	return em
}

func EntityRequestToEntity(equipment *models.Equipment) *entities.Equipment {

	history := make([]entities.EquipmentHistory, 0, len(equipment.EquipmentsHistory))
	for _, eh := range equipment.EquipmentsHistory {

		history = append(history, entities.EquipmentHistory(eh))
	}

	regionFlag := enums.RemapperRegionFlags[equipment.AssignedToRegion]

	return &entities.Equipment{
		ID:                equipment.ID,
		Name:              equipment.Name,
		Serial:            equipment.Serial,
		AssignedToRegion:  regionFlag,
		CurrentShipmentId: equipment.CurrentShipmentId,
		Type:              equipment.Type,
		Status:            equipment.Status,
		EquipmentsHistory: history,
		CreatedAt:         equipment.CreatedAt,
		UpdatedAt:         equipment.UpdatedAt,
	}
}

func ModelEfToEntitiesEf(mdf models.EquipmentFilter) entities.EquipmentFilter {
	return entities.EquipmentFilter{
		Serial:     mdf.Serial,
		Name:       mdf.Name,
		Type:       mdf.Type,
		Status:     mdf.Status,
		ShipmentId: mdf.ShipmentId,
	}
}
