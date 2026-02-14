package equipmentservice

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
	"centralService/internal/domain/interfaces"
	"centralService/internal/infra/trace"
	"context"
)

type EquipmentService struct {
	storage interfaces.EquipmentStorageInterface
}

func NewService(storage interfaces.EquipmentStorageInterface) *EquipmentService {
	return &EquipmentService{
		storage: storage,
	}
}

func (es *EquipmentService) List(context context.Context, filter entities.EquipmentFilter) ([]entities.Equipment, error) {

	log := trace.LogWithTraceID("equipment-service", context)

	log.Info("equipment.service.list calling")

	equipments, err := es.storage.List(context, filter)

	if err != nil {

		log.Info("equipment.service.list failed", "err", err.Error())
		err = TranslateStorageErrs(err)

		return nil, err
	}

	return equipments, nil
}

func (es *EquipmentService) Create(context context.Context, eq entities.Equipment) (string, error) {

	log := trace.LogWithTraceID("equipment-service", context).With("method", "Create")

	log.Info("equipment.service calling")

	newEquipmentId, err := es.storage.Create(context, eq)

	if err != nil {

		log.Info("equipment.service failed", "err", err.Error())
		err = TranslateStorageErrs(err)

		return "", err
	}

	return newEquipmentId, nil
}

func (es *EquipmentService) GetByID(context context.Context, id string) (*entities.Equipment, error) {
	log := trace.LogWithTraceID("equipment-service", context)

	log.Info("equipment.service.getById calling", "id", id)

	equipment, err := es.storage.GetByID(context, id)

	if err != nil {
		log.Error("equipment.service.getById failed", "id", id)

		err := TranslateStorageErrs(err)

		return nil, err
	}

	return equipment, nil
}

func (es *EquipmentService) ListByRegion(ctx context.Context, region enums.RegionFlag) ([]entities.Equipment, error) {
	log := trace.LogWithTraceID("equipment-service", ctx)

	log.Info("equipment.service.ListByRegion calling", "region", region)

	filter := entities.EquipmentFilter{
		Region: region,
	}

	equipments, err := es.storage.List(ctx, filter)

	if err != nil {
		log.Error("equipment.service.ListByRegion failed", "region", region)

		return nil, TranslateStorageErrs(err)
	}

	return equipments, nil
}
