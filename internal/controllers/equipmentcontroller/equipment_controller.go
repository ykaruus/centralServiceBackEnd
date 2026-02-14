package equipmentcontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/controllers/equipmentcontroller/models"
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
	"centralService/internal/http/response"
	"centralService/internal/infra/trace"
	"centralService/internal/service/equipmentservice"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EquipmentController struct {
	service *equipmentservice.EquipmentService
}

func NewController(service *equipmentservice.EquipmentService) *EquipmentController {
	return &EquipmentController{
		service: service,
	}
}

func (ec *EquipmentController) Create(ctx *gin.Context) {

	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	context := gen.CtxWithValue(ctx.Request.Context(), traceID)

	log := trace.LogWithTraceID("equipment-controller", context)

	var equipment *models.Equipment

	bodyErr := ctx.ShouldBindJSON(&equipment)

	if bodyErr != nil {

		log.Error("request failed", "err", bodyErr.Error())
		ctx.JSON(http.StatusBadRequest, response.BadRequest(ctx))
		return
	}

	if err := equipment.Validate(); err != nil {

		var apperr *controllers.HTTPResponse

		if errors.As(err, &apperr) {
			ctx.JSON(apperr.Code, response.Fail(ctx, apperr.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.ErrBadRequest)
		return
	}
	ee := EntityRequestToEntity(equipment)
	newEquipmentId, err := ec.service.Create(context, *ee)

	if err != nil {
		httpRes := TranslateServiceError(err)

		ctx.JSON(httpRes.Code, response.Fail(ctx, httpRes.ApiError))
		return
	}

	ctx.JSON(http.StatusOK, response.OK(ctx, newEquipmentId))

}

func (ec *EquipmentController) List(ctx *gin.Context) {

	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	context := gen.CtxWithValue(ctx.Request.Context(), traceID)

	equipmentFilter := models.EquipmentFilter{
		Name:       ctx.Query("name"),
		Serial:     ctx.Query("serial"),
		ShipmentId: ctx.Query("shipmentId"),
	}

	log := trace.LogWithTraceID("equipment-controller", context).With("method", "List")

	u, ok := ctx.Get("user")

	if !ok {
		log.Info("user not has save in ctx")
		ctx.JSON(http.StatusForbidden, response.Fail(ctx, &response.ErrForbidden))
		return
	}

	user, ok := u.(entities.User)

	if !ok {
		log.Info("user not has a valid user struct")
		ctx.JSON(http.StatusForbidden, response.Fail(ctx, &response.ErrForbidden))
		return
	}

	if err := equipmentFilter.Validate(); err != nil {
		var apperr *controllers.HTTPResponse

		if errors.As(err, &apperr) {
			ctx.JSON(apperr.Code, response.Fail(ctx, apperr.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.BadRequest(ctx))

		return
	}

	log.Info("request.started")

	if statusStr := ctx.Query("status"); statusStr != "" {
		status, err := strconv.Atoi(statusStr)

		if err != nil || status > 5 || status <= 0 {
			slog.Error("request.failed", "err", err.Error())

			responseModel := response.BadRequest(ctx)

			ctx.JSON(400, responseModel)
			return
		}

		equipmentFilter.Status = enums.EquipmentStatus(status)
	}

	if typeStr := ctx.Query("type"); typeStr != "" {
		typeN, err := strconv.Atoi(typeStr)

		if err != nil {

			log.Error("request.failed", "err", err.Error())

			responseModel := response.BadRequest(ctx)

			ctx.JSON(400, responseModel)
			return
		}

		if typeN > 4 || typeN <= 0 {
			log.Error("request.failed", "err", "A query 'type' está invalida")
			responseModel := response.BadRequest(ctx)

			ctx.JSON(400, responseModel)
			return
		}

		equipmentFilter.Type = enums.EquipmentType(typeN)
	}

	equipments := make([]entities.Equipment, 0)

	if err := user.HasPermission(enums.USER_PERMISSION_COORDINATOR); err == nil {
		equipments, err = ec.service.List(context, ModelEfToEntitiesEf(equipmentFilter))

		if err != nil {

			httpRes := TranslateServiceError(err)

			ctx.JSON(httpRes.Code, response.Fail(ctx, httpRes.ApiError))

			return
		}

	} else {
		equipments, err = ec.service.ListByRegion(context, user.AssignedToRegion)

		if err != nil {

			httpRes := TranslateServiceError(err)

			ctx.JSON(httpRes.Code, response.Fail(ctx, httpRes.ApiError))

			return
		}
	}

	equipmentRequest := EntitiesToEntitiesRequest(equipments)

	ctx.JSON(http.StatusOK, response.OK(ctx, equipmentRequest))

}

func (ec *EquipmentController) GetByID(c *gin.Context) {
	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(c)
	context := gen.CtxWithValue(c.Request.Context(), traceID)
	log := trace.LogWithTraceID("equipment-controller", context)

	id := c.Param("id")

	log.Info("request.started")

	equipment, err := ec.service.GetByID(context, id)

	if err != nil {
		log.Error("request.failed")
		httpRes := TranslateServiceError(err)

		c.JSON(httpRes.Code, response.Fail(c, httpRes.ApiError))
		return
	}

	c.JSON(http.StatusOK, response.OK(c, EntityToEntityRequest(equipment)))
}
