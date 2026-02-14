package equipmentstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/storage"
	"centralService/internal/storage/equipmentstorage/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func EntityToModel(equipment *entities.Equipment) (*models.Equipment, error) {
	//slog.Info("valores do equipment", "ID", equipment.ID, "SHIPMENT_ID", equipment.CurrentShipmentId, "tam do history", len(equipment.EquipmentsHistory))
	em := models.Equipment{
		AssignedToRegion: equipment.AssignedToRegion,
		Name:             equipment.Name,
		Serial:           equipment.Serial,
		Type:             equipment.Type,
		Status:           equipment.Status,
		CreatedAt:        equipment.CreatedAt,
		UpdatedAt:        equipment.UpdatedAt,
	}

	em.EquipmentsHistory = make([]models.EquipmentHistory, 0, len(equipment.EquipmentsHistory))

	if equipment.ID != "" {
		equipmentId, err := bson.ObjectIDFromHex(equipment.ID)

		if err != nil {
			return nil, storage.ErrEquipmentInvalidId
		}

		em.ID = equipmentId

	}

	if equipment.CurrentShipmentId != "" {
		shipmentId, err := bson.ObjectIDFromHex(equipment.CurrentShipmentId)

		if err != nil {
			return nil, storage.ErrShipmentInvalidID
		}

		em.CurrentShipmentId = &shipmentId
	}

	if len(equipment.EquipmentsHistory) != 0 {
		for _, h := range equipment.EquipmentsHistory {

			shipmentId, err := bson.ObjectIDFromHex(h.ShipmentID)

			if err != nil {
				return nil, storage.ErrShipmentInvalidID
			}

			em.EquipmentsHistory = append(em.EquipmentsHistory, models.EquipmentHistory{
				ShipmentID: shipmentId,
				AssignedAt: h.AssignedAt,
			})
		}
	}

	return &em, nil

}

func ModelToEntity(equipmentModel *models.Equipment) entities.Equipment {

	es := entities.Equipment{
		AssignedToRegion: equipmentModel.AssignedToRegion,
		ID:               equipmentModel.ID.Hex(),
		Name:             equipmentModel.Name,
		Serial:           equipmentModel.Serial,
		Status:           equipmentModel.Status,
		CreatedAt:        equipmentModel.CreatedAt,
		UpdatedAt:        equipmentModel.UpdatedAt,
	}

	if len(equipmentModel.EquipmentsHistory) != 0 {
		for _, h := range equipmentModel.EquipmentsHistory {
			es.EquipmentsHistory = append(es.EquipmentsHistory, entities.EquipmentHistory{
				ShipmentID: h.ShipmentID.Hex(),
				AssignedAt: h.AssignedAt,
			})
		}
	}
	if equipmentModel.CurrentShipmentId != nil {
		es.CurrentShipmentId = equipmentModel.CurrentShipmentId.Hex()
	}

	return es
}

func ModelsToEntitys(models []models.Equipment) []entities.Equipment {
	equipment := make([]entities.Equipment, 0, len(models))
	for _, m := range models {
		equipment = append(equipment, ModelToEntity(&m))
	}

	return equipment
}
