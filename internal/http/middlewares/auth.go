package middlewares

import (
	"centralService/internal/controllers/authcontroller"
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
	"centralService/internal/http/response"
	"centralService/internal/infra/trace"
	"centralService/internal/service/authservice"
	apperr "centralService/internal/service/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Auth : middlewares para autenticação das rotas
// esse middleware vai se o token recebido tem a role "technical"

type AuthenticateRole struct {
	token *authservice.TokenService
	user  *authservice.UserService
}

func NewAuthenticateRole(
	token *authservice.TokenService,
	user *authservice.UserService,
) *AuthenticateRole {

	return &AuthenticateRole{

		user:  user,
		token: token,
	}

}

func (r *AuthenticateRole) AuthValididateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gen := trace.NewGenerator()
		traceID := gen.GetTraceIDFromCtx(ctx)
		ctxWithTrace := gen.CtxWithValue(ctx.Request.Context(), traceID)

		log := trace.LogWithTraceID("auth-controller", ctxWithTrace).With("middleware", "AuthValidateToken")

		log.Info("autenticate token middleware started")

		token := ctx.Request.Header.Get("Authorization")

		tokenExtracted, err := r.token.ExtractToken(token)

		if err != nil {
			log.Error("token extract failed", "error", err)
			var apperr *apperr.AppError
			if errors.As(err, &apperr) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrorBody{
					Code:    apperr.Code,
					Message: apperr.Message,
				}))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
			return
		}

		err = r.token.VerifyToken(ctx, tokenExtracted)

		if err != nil {
			log.Error("token validation failed", "error", err)
			var apperr *apperr.AppError
			if errors.As(err, &apperr) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrorBody{
					Code:    apperr.Code,
					Message: apperr.Message,
				}))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
			return
		}

		_, err = r.token.DecriptToken(ctx, tokenExtracted)

		if err != nil {
			log.Error("token decripting  failed", "error", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))

			return
		}

		ctx.Next()
		log.Info("autenticate token middleware passed")
	}
}

func (r *AuthenticateRole) AuthHasRole(target_role enums.UserPermission) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gen := trace.NewGenerator()
		traceID := gen.GetTraceIDFromCtx(ctx)
		ctxWithTrace := gen.CtxWithValue(ctx.Request.Context(), traceID)

		log := trace.LogWithTraceID("auth-controller", ctxWithTrace).With("middleware", "AuthHasRole")

		log.Info("authenticate role started", "target_role", target_role)

		token := ctx.Request.Header.Get("Authorization")

		tokenExtracted, err := r.token.ExtractToken(token)

		if err != nil {
			log.Error("token extract failed", "error", err)
			var apperr *apperr.AppError
			if errors.As(err, &apperr) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrorBody{
					Code:    apperr.Code,
					Message: apperr.Message,
				}))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
			return
		}

		err = r.token.VerifyToken(ctxWithTrace, tokenExtracted)

		if err != nil {
			log.Error("token validation failed", "error", err)
			var apperr *apperr.AppError
			if errors.As(err, &apperr) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrorBody{
					Code:    apperr.Code,
					Message: apperr.Message,
				}))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
			return
		}

		payload, err := r.token.DecriptToken(ctxWithTrace, tokenExtracted)

		if err != nil {
			log.Error("token decripting  failed", "error", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))

			return
		}

		if err := r.user.VerifyRole(payload, target_role); err != nil {

			err := authcontroller.TranslateServiceError(err)

			ctx.AbortWithStatusJSON(err.Code, response.Fail(ctx, err.ApiError))
			return
		}

		ctx.Next()
		log.Info("authenticate role passed", "target_role", target_role)

	}
}

func (r *AuthenticateRole) SaveClaimsInContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gen := trace.NewGenerator()
		traceID := gen.GetTraceIDFromCtx(ctx)
		ctxWithTrace := gen.CtxWithValue(ctx.Request.Context(), traceID)

		log := trace.LogWithTraceID("auth-controller", ctxWithTrace).With("middleware", "SaveClaimsInContext")

		token := ctx.Request.Header.Get("Authorization")

		tokenExtracted, err := r.token.ExtractToken(token)

		if err != nil {
			log.Error("token extract failed", "error", err)
			var apperr *apperr.AppError
			if errors.As(err, &apperr) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrorBody{
					Code:    apperr.Code,
					Message: apperr.Message,
				}))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
			return
		}

		err = r.token.VerifyToken(ctxWithTrace, tokenExtracted)

		if err != nil {
			log.Error("token validation failed", "error", err)
			var apperr *apperr.AppError
			if errors.As(err, &apperr) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrorBody{
					Code:    apperr.Code,
					Message: apperr.Message,
				}))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))
			return
		}

		payload, err := r.token.DecriptToken(ctxWithTrace, tokenExtracted)

		if err != nil {
			log.Error("token decripting  failed", "error", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ctx, &response.ErrUnauthorized))

			return
		}

		u := entities.User{
			ID:               payload.ID,
			Roles:            payload.Roles,
			AssignedToRegion: payload.AssignedToRegion,
		}

		log.Info("saved claims of token in context", "claims", u)

		ctx.Set("user", u)

		ctx.Next()

	}
}
