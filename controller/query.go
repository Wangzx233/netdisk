package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netdisk/dao"
	"netdisk/util"
)

func Query(c *gin.Context)  {
	//cookie, _ := c.Request.Cookie("user_cookie")
	token := c.GetHeader("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSONP(http.StatusNoContent,"token验证失败")
		return
	}
	files, err := dao.Query(claims.Username)
	if err != nil {
		c.JSONP(http.StatusNoContent,err)
	}else {
		c.JSONP(http.StatusOK,files)
	}

}