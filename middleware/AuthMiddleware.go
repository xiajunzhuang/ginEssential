package middleware

import (
	"ginEssential/common"
	"ginEssential/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		//获取authorization header
		tokenString := context.GetHeader("Authorization")

		//validate token format
		if tokenString == "" || !strings.HasPrefix(tokenString,"Bearer "){
			context.JSON(http.StatusUnauthorized,gin.H{"code": 401,"msg": "权限不足"})
			context.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid{
			context.JSON(http.StatusUnauthorized,gin.H{"code": 401,"msg": "权限不足"})
			context.Abort()
			return
		}

		//验证通过后获取claims中的 userId
		userid := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user,userid)


		//用户
		if user.ID==0{
			context.JSON(http.StatusUnauthorized,gin.H{"code": 401,"msg": "权限不足"})
			context.Abort()
			return
		}

		//用户存在，将信息写入上下文
		context.Set("user",user)
		context.Next()
	}
}
