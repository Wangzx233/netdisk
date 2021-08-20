package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"netdisk/dao"
	"netdisk/model"
	"netdisk/util"
)

func ChangePath(c *gin.Context) {
	newPath := c.Query("new_path")
	md5hash := c.Query("md5hash")
	token := c.GetHeader("token")

	if newPath == "" || md5hash == "" || token == "" {
		c.String(http.StatusBadRequest, "参数缺少")
		return
	}
	claims, err := util.ParseToken(token)
	if err != nil {
		c.String(http.StatusBadRequest, "token错误")
		return
	}

	//修改数据库
	oldP, newP, err := dao.ChangePath(claims.Username, md5hash, newPath)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "修改失败")
		return
	}
	//修改文件位置
	err = model.ChangePath(oldP, newP)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "修改失败")
		return
	}

	c.String(http.StatusOK, "修改成功")
}

func Rename(c *gin.Context) {
	newName := c.Query("new_name")
	md5hash := c.Query("md5hash")
	token := c.GetHeader("token")
	if newName == "" || md5hash == "" || token == "" {
		c.String(http.StatusBadRequest, "参数缺少")
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		c.String(http.StatusBadRequest, "token错误")
		return
	}

	oldP, newP, err := dao.Rename(claims.Username, md5hash, newName)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "重命名失败,非文件拥有者")
		return
	}

	err = model.ChangePath(oldP, newP)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "重命名失败,非文件拥有者")
		return
	}

	c.String(http.StatusOK, "重命名成功")
}
