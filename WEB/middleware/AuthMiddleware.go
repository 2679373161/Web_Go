package middleware

import (
	"ginEssential/common"
	"ginEssential/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)
//中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context){
		//获取authorization header
		tokenString:=ctx.GetHeader("Authorization")

		//validata token formate
		if tokenString==""||!strings.HasPrefix(tokenString,"Bearer "){//tokenstring为空 或 不是以bearer开头
			ctx.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}

		tokenString=tokenString[7:]//由于bearer+空格占七位，从第七位开始取

		token,claims,err:=common.ParseToken(tokenString)//解析token
		if err!=nil||!token.Valid{
			ctx.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}

		//验证通过后获取claim中的userid
		userId:=claims.UserId

		DB:=common.GetDB()
		var user model.User
		DB.First(&user,userId)

		//验证用户是否存在
		if user.ID==0{
			ctx.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}
		//用户存在 将user信息写入上下文
		ctx.Set("user",user)

		ctx.Next()


	}
}
