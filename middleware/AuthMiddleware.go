package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"xmarket_gin/common"
	"xmarket_gin/model"
	"xmarket_gin/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取 authorization header
		tokenString := c.GetHeader("Authorization")
		//validate token format
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {

			response.Response(c, http.StatusUnauthorized, 401, nil, "权限不足")
			c.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			response.Response(c, http.StatusUnauthorized, 401, nil, "权限不足")
			c.Abort()
			return
		}

		//验证通过后获取claim中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)
		//用户不存在
		id, _ := strconv.Atoi(user.Telephone)
		if id == 0 {
			response.Response(c, http.StatusUnauthorized, 401, nil, "权限不足")
			c.Abort()
			return
		}
		//用户存在 将user信息写入上下文
		c.Set("user", user)
		c.Next()
	}
}
