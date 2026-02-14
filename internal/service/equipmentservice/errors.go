package equipmentservice

import (
	errors "centralService/internal/service/errors"
)

var (
	ErrEquipmentNotFound = errors.New(
		"EQUIPMENT_NOT_FOUND",
		"Equipamento não encontrado",
	)

	ErrEquipmentInactive = errors.New(
		"EQUIPMENT_INACTIVE",
		"Equipamento está inativo",
	)

	ErrEquipmentAlreadyExists = errors.New(
		"EQUIPMENT_ALREADY_EXISTS",
		"Equipamento já cadastrado",
	)
	ErrEquipmentInvalidId         = errors.New("EQUIPMENT_INVALID_ID", "O id de equipamento está invalido")
	ErrEquipmentHasActiveShipment = errors.New("EQUIPMENT_HAS_ACTIVE_SHIPMENT", "O remessa já tem um equipamento ativo")
)
