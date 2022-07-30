package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// CORSMiddleware
func CORSMiddleware()gin.HandlerFunc{
	return func(ctx *gin.Context){
		host:=viper.GetString("server.host")
		port:=viper.GetString("server.frontPort")
		if port !=""{
			ctx.Writer.Header().Set("Access-Control-Allow-Origin","http://"+host+":"+port)//后面填写域名，*代表所有域名都可以访问
		}else{
			ctx.Writer.Header().Set("Access-Control-Allow-Origin","http://"+host)
		}
		ctx.Writer.Header().Set("Access-Control-Max-Age","86400")//缓存时间
		ctx.Writer.Header().Set("Access-Control-Allow-Methods","*")//方法，*允许所有方法
		ctx.Writer.Header().Set("Access-Control-Allow-Headers","*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials","true")

		if ctx.Request.Method==http.MethodOptions{
			ctx.AbortWithStatus(200)
		}else{
			ctx.Next()//继续向下传递
		}
	}
}
