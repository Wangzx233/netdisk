package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"netdisk/dao"
	"netdisk/model"
	"netdisk/util"
)

func DeleteFile(c *gin.Context) {
	token := c.GetHeader("token")
	md5hash := c.Query("md5hash")

	claims, err := util.ParseToken(token)
	if err != nil {
		c.String(http.StatusBadRequest, "token错误")
		return
	}

	path,err := dao.DeleteFile(claims.Username, md5hash)
	if err != nil {
		log.Println("dao err : deleteFile err:",err)
		c.String(http.StatusBadRequest, "删除失败,文件不存在或非文件拥有者")
		return
	}

	err = model.DeleteFile(path)
	if err != nil {
		log.Println("model err : deleteFile err:",err)
		c.String(http.StatusBadRequest, "删除失败，文件不存在或非文件拥有者")
		return
	}

	Client.Del(path)
	c.String(http.StatusOK, "删除成功")
}
