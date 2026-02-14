package authcontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/controllers/authcontroller/models"
	"centralService/internal/controllers/usercontroller"
	"centralService/internal/domain/entities"
	"centralService/internal/http/response"
	"centralService/internal/infra/trace"
	"centralService/internal/service/authservice"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *authservice.AuthService
}

func NewAuthController(service *authservice.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	c := gen.CtxWithValue(ctx.Request.Context(), traceID)

	log := trace.LogWithTraceID("auth-controller", c)

	log.Info("request.started")

	var google_id models.Google_ID

	err := ctx.ShouldBindJSON(&google_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}

	slog.Info("received token id", "idtoken", google_id.IdToken)

	if err := google_id.Validate(); err != nil {

		var apperr *controllers.HTTPResponse

		if errors.As(err, &apperr) {
			ctx.JSON(apperr.Code, response.Fail(ctx, apperr.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}

	slog.Info("received token id", "idtoken", google_id.IdToken)

	err = ac.service.ValidateIdToken(ctx, google_id.IdToken)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrUnauthorized))
		return
	}

	payload, err := ac.service.DecriptIdToken(c, google_id.IdToken)

	token, err := ac.service.Login(c, payload)

	if err != nil {

		cr := TranslateServiceError(err)

		ctx.JSON(cr.Code, response.Fail(ctx, cr.ApiError))
		return
	}

	ctx.JSON(http.StatusOK, response.OK(ctx, token))

}

func (au *AuthController) Me(ctx *gin.Context) {
	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	c := gen.CtxWithValue(ctx.Request.Context(), traceID)

	log := trace.LogWithTraceID("auth-controller", c)

	log.Info("request.started")

	token := ctx.Request.Header.Get("Authorization")

	if token == "" {
		ctx.JSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
		return
	}

	user, err := au.service.GetUserFromToken(c, token)

	if err != nil {
		cr := TranslateServiceError(err)

		ctx.JSON(cr.Code, response.Fail(ctx, cr.ApiError))
		return
	}

	eu, err := usercontroller.EntityUserToModelsUser(user)

	if err != nil {

		log.Error("auth-controller.Me failed", "error", err.Error())
		var co *controllers.HTTPResponse

		if errors.As(err, &co) {
			ctx.JSON(co.Code, response.Fail(ctx, co.ApiError))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.Fail(ctx, &response.ErrInternal))
		return
	}

	ctx.JSON(http.StatusOK, response.OK(ctx, eu))

}

func (au *AuthController) GenerateToken(ctx *gin.Context) {
	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	c := gen.CtxWithValue(ctx.Request.Context(), traceID)

	log := trace.LogWithTraceID("auth-controller", c)

	log.Info("request.started")

	token, err := au.service.Login(ctx, &entities.User{
		Email:        "icarogagarincontato@gmail.com",
		IsFirstLogin: false,
	})

	if err != nil {

		cr := TranslateServiceError(err)

		ctx.JSON(cr.Code, response.Fail(ctx, cr.ApiError))
		return
	}

	ctx.JSON(http.StatusOK, response.OK(ctx, token))

}
