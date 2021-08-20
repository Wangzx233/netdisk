package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"netdisk/dao"
	"netdisk/util"
)

func Register(c *gin.Context) {
	user := dao.User{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Println("controller err : Register err : shouldBind err",err)
	}

	ans := dao.Register(user.Username, user.Password)
	//c.SetCookie("user_cookie", user.Username, 1000, "/", "localhost", false, true)
	token, err := util.GenerateToken(user.Username)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest,"生成token失败")
	}
	c.Header("token",token)
	c.String(http.StatusOK,ans)
}

func Login(c *gin.Context) {
	user := dao.User{}
	err := c.ShouldBind(&user)
	if err != nil {
		log.Println("controller err : Register err : shouldBind err",err)
	}

	login, b := dao.Login(user.Username, user.Password)
	if !b {
		c.String(http.StatusNotAcceptable,login)
	}else {
		//c.SetCookie("user_cookie", user.Username, 1000, "/", "localhost", false, true)
		token, err := util.GenerateToken(user.Username)
		if err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest,"生成token失败")
		}
		c.Header("token",token)
		c.String(http.StatusOK,"登录成功")
	}
}