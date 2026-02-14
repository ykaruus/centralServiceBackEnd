package enums

type ShipmentStatus int

const (
	SHIPMENT_STATUS_DRAFT      ShipmentStatus = 1
	SHIPMENT_STATUS_IN_TRANSIT ShipmentStatus = 2
	SHIPMENT_STATUS_COMPLETED  ShipmentStatus = 3
	SHIPMENT_STATUS_CANCELED   ShipmentStatus = 4
)

type ShipmentType int

const (
	SHIPMENT_TYPE_SENT   ShipmentStatus = 1
	SHIPMENT_TYPE_RETURN ShipmentStatus = 2
)

type OperationType int

const (
	SHIPMENT_OPERATION_TYPE_SENT   OperationType = 1
	SHIPMENT_OPERATION_TYPE_REFUND OperationType = 2
)
