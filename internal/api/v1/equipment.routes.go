package v1

import (
	"centralService/internal/controllers/equipmentcontroller"
	"centralService/internal/domain/interfaces"

	"github.com/gin-gonic/gin"
)

func RegisterEquipmentsRoutes(router *gin.RouterGroup, equipment_controller *equipmentcontroller.EquipmentController, m interfaces.Middleware) {
	equipments := router.Group("/")
	{
		//equipments.Use(middlewares.TraceIDMiddleware(trace.NewGenerator()))
		//equipments.Use(m.AuthValididateToken())
		//equipments.Use(m.AuthHasRole(enums.USER_PERMISSION_TECHNICAL))
		//equipments.Use(m.SaveClaimsInContext())
		equipments.POST("/equipments", equipment_controller.Create)
		equipments.GET("/equipments", equipment_controller.List)
		//equipments.DELETE("/equipments/:id", equipment_controller.Delete)
		//equipments.PUT("/equipments/:id", equipment_controller.Put)
		equipments.GET("/equipments/:id", equipment_controller.GetByID)
		//equipments.GET("/equipments/:id/shipment", equipment_controller.) // retorna o shipment atrelado ao id
	}
}
