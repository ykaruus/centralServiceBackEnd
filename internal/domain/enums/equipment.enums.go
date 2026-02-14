package enums

type EquipmentStatus int

type EquipmentType int

const (
	EQUIPMENT_TYPE_NOTEBOOK EquipmentType = 1
	EQUIPMENT_TYPE_DESKTOP  EquipmentType = 2
	EQUIPMENT_TYPE_PHONE    EquipmentType = 3
	EQUIPMENT_TYPE_CHIP     EquipmentType = 4
	EQUIPMENT_TYPE_HEADSET  EquipmentType = 5
	EQUIPMENT_TYPE_TECLADO  EquipmentType = 6
	EQUIPMENT_TYPE_MOUSE    EquipmentType = 7
)

const (
	EQUIPMENT_STATUS_AVAILABLE    EquipmentStatus = 1
	EQUIPMENT_STATUS_DISPOSAL     EquipmentStatus = 2
	EQUIPMENT_STATUS_SENDED       EquipmentStatus = 3
	EQUIPMENT_STATUS_OUT_STANDARD EquipmentStatus = 4
	EQUIPMENT_STATUS_MAINTENANCE  EquipmentStatus = 5
)
