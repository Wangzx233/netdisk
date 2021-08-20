package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NeedLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, e := context.Request.Cookie("user_cookie")
		if e == nil {
			//刷新cookie时间
			context.SetCookie(cookie.Name, cookie.Value, 1000, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

			context.Next()
		} else {
			context.Abort()
			context.String(http.StatusUnauthorized,"请先登录")
		}
	}
}

func NeedNoLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, e := context.Request.Cookie("user_cookie")
		if e == nil {
			//刷新cookie时间
			context.SetCookie(cookie.Name, cookie.Value, 1000, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

			context.Abort()
			context.String(http.StatusOK,"已处于登录状态")
		} else {
			context.Next()
		}
	}
}