package interfaces

import (
	"centralService/internal/domain/enums"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	AuthHasRole(target enums.UserPermission) gin.HandlerFunc
	SaveClaimsInContext() gin.HandlerFunc
	AuthValididateToken() gin.HandlerFunc
}
