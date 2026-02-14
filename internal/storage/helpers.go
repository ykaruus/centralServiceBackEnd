package storage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TranslateDBErr(err error) error {
	switch {
	case err == nil:
		return nil

	case errors.Is(err, context.Canceled):
		return err

	case errors.Is(err, context.DeadlineExceeded):
		return err

	case errors.Is(err, mongo.ErrNoDocuments):
		return ErrNotFound

	case mongo.IsDuplicateKeyError(err):
		return ErrDuplicateKey

	default:
		return err
	}
}
