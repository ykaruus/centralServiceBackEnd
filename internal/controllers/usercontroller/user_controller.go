package usercontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/controllers/usercontroller/models"
	"centralService/internal/http/response"
	"centralService/internal/infra/trace"
	"centralService/internal/service/authservice"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *authservice.UserService
}

func NewUserController(service *authservice.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (u *UserController) Create(ctx *gin.Context) {
	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	c := gen.CtxWithValue(ctx.Request.Context(), traceID)

	log := trace.LogWithTraceID("user-controller", c).With("method", "Create")

	log.Info("request.started")

	var user models.User

	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		log.Error("request.failed json", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}

	if err := user.Validate(); err != nil {
		log.Error("request.failed", "error", err.Error())
		var apperr *controllers.HTTPResponse

		if errors.As(err, &apperr) {
			ctx.JSON(apperr.Code, response.Fail(ctx, apperr.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}

	eu, mErr := ModelsUserToEntityUser(&user)

	if mErr != nil {

		log.Error("request.failed", "error", mErr.Error())
		var httpErr *controllers.HTTPResponse

		if errors.As(mErr, &httpErr) {
			ctx.JSON(httpErr.Code, response.Fail(ctx, httpErr.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}
	userId, err := u.service.Create(c, *eu)

	if err != nil {
		log.Error("request.failed", "error", err.Error())
		httpRes := TranslateServiceError(err)

		ctx.JSON(httpRes.Code, response.Fail(ctx, httpRes.ApiError))
		return
	}

	ctx.JSON(http.StatusOK, response.OK(ctx, userId))

}

func (u *UserController) List(ctx *gin.Context) {
	gen := trace.NewGenerator()
	traceID := gen.GetTraceIDFromCtx(ctx)
	c := gen.CtxWithValue(ctx.Request.Context(), traceID)

	log := trace.LogWithTraceID("user-controller", c).With("method", "List")

	log.Info("request.started")

	uf := &models.UserFilter{
		Email: ctx.Query("email"),
		Name:  ctx.Query("name"),
		ID:    ctx.Query("id"),
	}

	if len(ctx.QueryArray("roles")) != 0 && len(ctx.QueryArray("roles")) > 3 {
		ctx.JSON(http.StatusBadRequest, controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
			Code:    response.ErrBadRequest.Code,
			Message: response.ErrBadRequest.Message,
			Fields: []string{
				"Argumentos demais nas querys otário zé mané",
			},
		}))

		return
	}

	if len(ctx.QueryArray("roles")) != 0 {
		for _, r := range ctx.QueryArray("roles") {
			uf.Roles = append(uf.Roles, models.Role{RoleName: r})
		}
	}

	if err := uf.Validate(); err != nil {
		log.Error("Validating  models.UserFilter failed", "err", err)

		var con *controllers.HTTPResponse

		if errors.As(err, &con) {

			ctx.JSON(con.Code, response.Fail(ctx, con.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}

	euf, err := ModelUserFilterToEntities(uf)

	if err != nil {
		log.Error("Mapping models.UserFilter failed", "err", err)
		var con *controllers.HTTPResponse

		if errors.As(err, &con) {

			ctx.JSON(con.Code, response.Fail(ctx, con.ApiError))
			return
		}

		ctx.JSON(http.StatusBadRequest, response.Fail(ctx, &response.ErrBadRequest))
		return
	}

	users, err := u.service.List(c, euf)

	if err != nil {

		err := TranslateServiceError(err)

		ctx.JSON(err.Code, response.Fail(ctx, err.ApiError))
		return

	}

	musers, err := EntitiesUsersToModelUsers(users)

	if err != nil {
		log.Error("Failed mapping entities.Users to models.Users", "error", err.Error())
		ctx.JSON(500, response.Fail(ctx, &response.ErrInternal))
		return
	}

	ctx.JSON(http.StatusOK, response.OK(ctx, musers))

}
