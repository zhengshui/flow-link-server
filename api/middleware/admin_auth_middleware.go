package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
	"github.com/zhengshui/flow-link-server/internal/tokenutil"
)

// AdminAuthMiddleware 管理员权限验证中间件
// 需要在 JwtAuthMiddleware 之后使用
func AdminAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) != 2 {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "未授权访问"))
			c.Abort()
			return
		}

		authToken := t[1]
		role, err := tokenutil.ExtractRoleFromToken(authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "无效的令牌"))
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, domain.NewErrorResponse(403, "需要管理员权限"))
			c.Abort()
			return
		}

		c.Set("x-user-role", role)
		c.Next()
	}
}
