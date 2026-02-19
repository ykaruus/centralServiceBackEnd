package equipmentstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/interfaces"
	"centralService/internal/infra/trace"
	"centralService/internal/storage"
	"centralService/internal/storage/db"
	"centralService/internal/storage/equipmentstorage/models"
	"context"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EquipmentStorage struct {
	clientDB   *db.ClientDB
	collection *mongo.Collection
}

func NewStorage(clientDB *db.ClientDB, collname string) interfaces.EquipmentStorageInterface {
	collection := clientDB.LoadCollection(collname)
	return &EquipmentStorage{
		clientDB:   clientDB,
		collection: collection,
	}
}

var _ interfaces.EquipmentStorageInterface = (*EquipmentStorage)(nil)

func (e *EquipmentStorage) InitIndexes() error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "serial", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	err := e.clientDB.CreateUniqueIndexes(e.collection, indexModels)

	if err != nil {

		return err
	}

	return nil

}

func (s *EquipmentStorage) Create(context_count context.Context, equipment entities.Equipment) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("equipment-storage", context_count)

	log.Info("equipment-storage.create")

	equipmentModel, err := EntityToModel(&equipment)

	if err != nil {
		return "", err
	}

	equipmentModel.CreatedAt = time.Now()
	equipmentModel.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, equipmentModel)
	if err != nil {
		log.Error("equipment-storage.create.failed", "err", err.Error())

		err = storage.TranslateDBErr(err)

		return "", err
	}
	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (eq *EquipmentStorage) UpdateByFilter(c context.Context, f *entities.EquipmentFilter) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)

	defer cancel()

	log := trace.LogWithTraceID("equipment-storage", ctx)

	log.Info("equipment-storage.called")

	idObj, err := bson.ObjectIDFromHex(f.ID)

	if err != nil {
		return storage.ErrEquipmentInvalidId
	}
	now := time.Now()

	query := bson.M{}

	query["updated_at"] = now

	if f.Name != "" {
		query["name"] = f.Name
	}
	if f.Serial != "" {
		query["serial"] = f.Serial
	}

	if f.ShipmentId != "" {
		if shipmentId, err := bson.ObjectIDFromHex(f.ShipmentId); err == nil {
			query["current_shipment_id"] = shipmentId
		}

		return storage.ErrShipmentInvalidID
	}

	if f.Region != 0 {
		query["region_flag"] = f.Region
	}

	if f.Type != 0 {
		query["type"] = f.Type
	}

	if f.Status != 0 {
		query["status"] = f.Status
	}

	update := bson.D{{Key: "$set", Value: query}}

	opts := options.FindOneAndUpdate()

	err = eq.collection.FindOneAndUpdate(ctx, bson.M{"_id": idObj}, update, opts).Err()

	if err != nil {
		return storage.TranslateDBErr(err)
	}

	return nil
}

func (eq *EquipmentStorage) Replace(c context.Context, ee *entities.Equipment) error {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)

	defer cancel()

	log := trace.LogWithTraceID("equipment-storage", ctx)

	log.Info("equipment-storage.called")

	idObj, err := bson.ObjectIDFromHex(ee.ID)

	if err != nil {
		return storage.ErrEquipmentInvalidId
	}
	now := time.Now()
	ee.UpdatedAt = now

	em, err := EntityToModel(ee)

	if err != nil {
		return err
	}

	err = eq.collection.FindOneAndReplace(ctx, bson.M{"_id": idObj}, em).Err()

	if err != nil {
		return storage.TranslateDBErr(err)
	}

	return nil
}

func (s *EquipmentStorage) GetByID(context_count context.Context, id string) (*entities.Equipment, error) {
	ctx, cancel := context.WithTimeout(context_count, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("equipment-storage", context_count)

	log.Info("equipment-storage.called")

	var result models.Equipment
	idObj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Error("equipment-storage.failed", "error", err.Error())
		return nil, storage.ErrEquipmentInvalidId
	}

	err = s.collection.FindOne(ctx, bson.M{"_id": idObj}).Decode(&result)
	if err != nil {
		log.Error("equipment-storage.failed", "error", err.Error())

		err := storage.TranslateDBErr(err)

		return nil, err
	}

	equipment := ModelToEntity(&result)

	return &equipment, nil
}

func (s *EquipmentStorage) Delete(context_count context.Context, id string) error {
	ctx, cancel := context.WithTimeout(context_count, 10*time.Second)
	defer cancel()

	log := slog.With("layer", "equipment-storage", "traceID", context_count.Value("traceID").(string))

	log.Info("equipment.storage.delete started")

	idObj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = s.collection.FindOneAndDelete(ctx, bson.M{"_id": idObj}).Err()

	if err != nil {
		log.Error("equipment.storage.delete failed", "err", err.Error())

		return err
	}

	return nil
}

func (s *EquipmentStorage) List(context_count context.Context, filter entities.EquipmentFilter) ([]entities.Equipment, error) {
	ctx, cancel := context.WithTimeout(context_count, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("equipment-controller", context_count)

	log.Info("equipment.storage.list started")

	query := bson.M{}

	if filter.ID != "" {
		query["_id"] = filter.ID
	}
	if filter.Name != "" {
		query["name"] = filter.Name
	}
	if filter.Serial != "" {
		query["serial"] = filter.Serial
	}
	if filter.Type != 0 {
		query["type"] = filter.Type
	}
	if filter.ShipmentId != "" {
		query["shipment_id"] = filter.ShipmentId
	}

	if filter.Status != 0 {
		query["status"] = filter.Status
	}

	if filter.Region != 0 {
		query["region_flag"] = filter.Region
	}

	cursor, err := s.collection.Find(ctx, query)
	equipments := make([]entities.Equipment, 0)

	if err != nil {
		log.Error("equipment.storage.list failed", "err", err.Error())
		err = storage.TranslateDBErr(err)
		return []entities.Equipment{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var equipmentModel models.Equipment

		err := cursor.Decode(&equipmentModel)

		if err != nil {
			log.Error("equipment.storage.list failed", "err", err.Error())
			err = storage.TranslateDBErr(err)
			return nil, err
		}

		equipment := ModelToEntity(&equipmentModel)

		equipments = append(equipments, equipment)

	}

	log.Info("equipment.storage.list passed")

	return equipments, nil
}

func (s *EquipmentStorage) ListEquipmentsId(context_cont context.Context, equipmentsId []string) ([]entities.Equipment, error) {
	ctx, cancel := context.WithTimeout(context_cont, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("equipment-storage", context_cont)

	log.Info("equipment-storage.ListEquipmentsId called")

	result := make([]entities.Equipment, 0)
	var query bson.M = bson.M{"_id": bson.M{"$in": equipmentsId}}

	cursor, err := s.collection.Find(ctx, query)
	if err != nil {

		log.Info("equipment-storage.ListEquipmentsId called", "error", err.Error())

		err = storage.TranslateDBErr(err)

		return nil, err
	}

	for cursor.Next(context_cont) {
		var em models.Equipment
		err := cursor.Decode(&em)

		if err != nil {
			err := storage.TranslateDBErr(err)
			return nil, err
		}

	}

	return result, nil
}

func (s *EquipmentStorage) UpdateByEquipmentsId(context_count context.Context, equipmentsId []string, up bson.M) error {
	ctx, cancel := context.WithTimeout(context_count, 10*time.Second)
	defer cancel()

	update := bson.D{{Key: "$set", Value: up}}

	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": bson.M{"$in": equipmentsId}}, update).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil
		}
		return err
	}

	return nil
}
