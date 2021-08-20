package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"netdisk/util"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *util.Claims
		token:=c.GetHeader("token")
		if token==""{
			c.String(http.StatusUnauthorized,"请先登录")
			c.Abort()
		}else{
			var err error
			claims,err = util.ParseToken(token)
			fmt.Println("解析出来的claims:",claims)
			fmt.Println("解析出来的err:",err)
			if err!=nil{
				c.String(http.StatusUnauthorized,"token解析错误")
				c.Abort()
			}else if time.Now().Unix() > claims.ExpiresAt{
				c.String(http.StatusUnauthorized,"已超时，请重新登录")
				c.Abort()
			}
		}
		c.Next()
	}
}