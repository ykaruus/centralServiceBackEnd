package authstorage

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/interfaces"
	"centralService/internal/infra/trace"
	"centralService/internal/storage"
	"centralService/internal/storage/authstorage/models"
	"centralService/internal/storage/db"
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserStorage struct {
	clientDB   *db.ClientDB
	collection *mongo.Collection
}

func NewUserStorage(clientDB *db.ClientDB, collname string) interfaces.UserStorageInterface {
	collection := clientDB.LoadCollection(collname)
	return &UserStorage{
		clientDB:   clientDB,
		collection: collection,
	}
}

func (e *UserStorage) InitIndexes() error {
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	err := e.clientDB.CreateUniqueIndexes(e.collection, indexModels)

	if err != nil {

		return err
	}

	return nil

}

func (u *UserStorage) Create(cont context.Context, user entities.User) (string, error) {
	ctx, cancel := context.WithTimeout(cont, 10*time.Second)

	defer cancel()

	mu, err := EntityUserToModelUser(&user)

	if err != nil {
		return "", err
	}

	mu.CreatedAt = time.Now()
	mu.UpdatedAt = time.Now()

	result, err := u.collection.InsertOne(ctx, mu)

	if err != nil {
		slog.Error("user-storage.Create failed", "error", err.Error())

		err := storage.TranslateDBErr(err)

		return "", err
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil

}

func (u *UserStorage) GetByEmail(c context.Context, email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var result models.User

	log := trace.LogWithTraceID("user-storage", c)

	log.Info("user-storage.GetByEmail called")

	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {

		log.Error("user-storage.GetByEmail failed", "error", err.Error())

		return nil, storage.TranslateDBErr(err)

	}

	eUser := ModelUserToEntityUser(&result)

	return &eUser, nil
}

func (u *UserStorage) GetByID(c context.Context, id string) (*entities.User, error) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var result models.User

	log := trace.LogWithTraceID("user-storage", c)

	log.Info("user-storage.GetById called")

	idObj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Error("user-storage.failed, id invalid", "error", err.Error())
		return nil, storage.ErrEquipmentInvalidId
	}

	err = u.collection.FindOne(ctx, bson.M{"_id": idObj}).Decode(&result)
	if err != nil {

		log.Error("user-storage.GetByEmail failed", "error", err.Error())

		err := storage.TranslateDBErr(err)

		return nil, err

	}

	eUser := ModelUserToEntityUser(&result)

	return &eUser, nil
}

func (u *UserStorage) List(c context.Context, userfilter *entities.UserFilter) ([]entities.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("user-storage", c)

	log.Info("user-storage.Update called")

	query := bson.M{}

	if userfilter.Email != "" {
		query["email"] = userfilter.Email
	}
	if userfilter.Name != "" {
		query["name"] = userfilter.Name
	}
	if userfilter.ID != "" {
		query["_id"] = userfilter.ID
	}
	if len(userfilter.Roles) != 0 {
		query["role"] = userfilter.Roles
	}

	opts := options.Find().SetLimit(100)

	cursor, err := u.collection.Find(ctx, query, opts)

	if err != nil {
		return nil, storage.TranslateDBErr(err)
	}

	users := make([]entities.User, 0)

	for cursor.Next(ctx) {

		var user models.User

		err := cursor.Decode(&user)

		if err != nil {
			return nil, storage.TranslateDBErr(err)
		}

		users = append(users, ModelUserToEntityUser(&user))

	}

	return users, nil

}

func (u *UserStorage) Update(c context.Context, eu *entities.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("user-storage", c).With("method", "Update")

	log.Info("user-storage.Update called")

	mu, err := EntityUserToModelUser(eu)

	if err != nil {
		return err
	}

	opts := options.FindOneAndReplace()

	up := models.User{

		Role:         mu.Role,
		Email:        mu.Email,
		IsFirstLogin: mu.IsFirstLogin,
		Name:         mu.Name,
		Picture:      mu.Picture,
		CreatedAt:    mu.CreatedAt,
		UpdatedAt:    mu.UpdatedAt,
	}
	update := bson.D{{Key: "set", Value: mu}}

	err = u.collection.FindOneAndReplace(ctx, bson.M{"_id": up}, update, opts).Err()

	if err != nil {
		log.Error("user-storage.Update failed", "error", err.Error())

		err := storage.TranslateDBErr(err)

		return err
	}

	return nil

}

func (u *UserStorage) UpdateNameAndPicture(c context.Context, id string, name string, picture string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("user-storage", c).With("method", "UpdateNameAndPicture")

	log.Info("user-storage.Update called")

	idObj, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return storage.ErrInvalidUserId
	}

	update := bson.D{{Key: "$set", Value: bson.M{
		"name":        name,
		"picture":     picture,
		"first_login": false,
	}}}

	opts := options.FindOneAndUpdate()

	err = u.collection.FindOneAndUpdate(ctx, bson.M{"_id": idObj}, update, opts).Err()

	if err != nil {
		log.Error("user-storage.Update failed", "error", err.Error())

		err := storage.TranslateDBErr(err)

		return err
	}

	return nil

}

func (u *UserStorage) UpdateByFilter(ctx context.Context, userID string, euf *entities.UserFilter) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log := trace.LogWithTraceID("user-storage", ctx).With("method", "UpdateByFilter")

	log.Info("user-storage.Update called")

	idObj, err := bson.ObjectIDFromHex(userID)

	if err != nil {
		return storage.ErrInvalidUserId
	}

	query := bson.M{}

	if euf.Email != "" {
		query["email"] = euf.Email
	}
	if euf.Name != "" {
		query["name"] = euf.Name
	}
	if len(euf.Roles) != 0 {
		query["roles"] = euf.Roles
	}
	if euf.LastAccessAt != nil {
		query["lastAccess_at"] = euf.LastAccessAt
	}
	if euf.Picture != "" {
		query["picture"] = euf.Picture
	}

	now := time.Now()

	query["updated_at"] = now

	update := bson.D{{Key: "$set", Value: query}}

	opts := options.FindOneAndUpdate()

	err = u.collection.FindOneAndUpdate(ctx, bson.M{"_id": idObj}, update, opts).Err()

	if err != nil {
		log.Error("user-storage.Update failed", "error", err.Error())

		err := storage.TranslateDBErr(err)

		return err
	}

	return nil
}

func (u *UserStorage) Delete(c context.Context, id string) error {
	return nil
}
