package equipmentservice

import (
	apperr "centralService/internal/service/errors"
	"centralService/internal/storage"
	"errors"
)

func TranslateStorageErrs(err error) error {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return ErrEquipmentNotFound
	case errors.Is(err, storage.ErrEquipmentInvalidId):
		return ErrEquipmentInvalidId
	case errors.Is(err, storage.ErrDuplicateKey):
		return ErrEquipmentAlreadyExists
	default:
		return apperr.Wrap("EQUIPMENT_STORAGE_ERROR", "Ocorreu um erro desconhecido ao realizar esta operação.", err)
	}
}
