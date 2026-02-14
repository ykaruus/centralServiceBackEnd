package shipmentstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/storage"
	"centralService/internal/storage/shipmentstorage/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func EntityToModel(es *entities.Shipment) (*models.Shipment, error) {

	if len(es.EquipmentIDs) != 0 || es.ID != "" || es.AssignedTo != "" {
		equipmentsIds := make([]bson.ObjectID, 0, len(es.EquipmentIDs))

		for _, id := range es.EquipmentIDs {
			ID, err := bson.ObjectIDFromHex(id)

			if err != nil {
				return nil, storage.ErrEquipmentInvalidId
			}

			equipmentsIds = append(equipmentsIds, ID)
		}

		IDobj, err := bson.ObjectIDFromHex(es.ID)

		if err != nil {
			return nil, storage.ErrEquipmentInvalidId
		}
		assignedToId, err := bson.ObjectIDFromHex(es.AssignedTo)

		if err != nil {
			return nil, storage.ErrEquipmentInvalidId
		}

		return &models.Shipment{
			ID:           IDobj,
			EquipmentIDs: equipmentsIds,
			Status:       es.Status,
			Type:         es.Type,
			Mail: models.Mail{
				SentAt:        es.Mail.SentAt,
				DeliveredAt:   es.Mail.DeliveredAt,
				OperationType: es.Mail.OperationType,
				PostCode:      es.Mail.PostCode,
				TrackCode:     es.Mail.TrackCode,
			},
			AssignedTo: assignedToId,
			CreatedAt:  es.CreatedAt,
			UpdatedAt:  es.UpdatedAt,
		}, nil
	}

	return &models.Shipment{
		Status: es.Status,
		Type:   es.Type,
		Mail: models.Mail{
			SentAt:        es.Mail.SentAt,
			DeliveredAt:   es.Mail.DeliveredAt,
			OperationType: es.Mail.OperationType,
			PostCode:      es.Mail.PostCode,
			TrackCode:     es.Mail.TrackCode,
		},
		CreatedAt: es.CreatedAt,
		UpdatedAt: es.UpdatedAt,
	}, nil

}

func ModelToEntity(ms *models.Shipment) *entities.Shipment {
	equipmentsIds := make([]string, 0, len(ms.EquipmentIDs))

	for _, id := range ms.EquipmentIDs {

		equipmentsIds = append(equipmentsIds, id.Hex())
	}

	return &entities.Shipment{
		ID:           ms.ID.Hex(),
		EquipmentIDs: equipmentsIds,
		AssignedTo:   ms.AssignedTo.Hex(),
		Type:         ms.Type,
		Status:       ms.Status,
		Mail: entities.Mail{
			SentAt:        ms.Mail.SentAt,
			DeliveredAt:   ms.Mail.DeliveredAt,
			OperationType: ms.Mail.OperationType,
			PostCode:      ms.Mail.PostCode,
			TrackCode:     ms.Mail.TrackCode,
		},
		CreatedAt: ms.CreatedAt,
		UpdatedAt: ms.UpdatedAt,
	}
}
