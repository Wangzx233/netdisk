package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"netdisk/dao"
	"netdisk/util"
)

func Share(c *gin.Context) {
	md5hash := c.Query("md5hash")
	power := c.Query("power")
	secret := c.Query("secret")
	if power == "" || md5hash == "" || secret == "" && power == "2" {
		c.String(http.StatusBadRequest, "参数缺少")
		return
	}

	token := c.GetHeader("token")
	cliams, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"url": "token有误",
		})
		return
	}

	err = dao.Share(md5hash, power, cliams.Username, secret)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"url": "文件不存在或您不是文件拥有者",
		})
		return
	}
	url:="http://localhost:8080" + "/file/download?md5hash=" + md5hash
	bytes, err := qrcode.Encode(url, qrcode.High, 256)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"url": url,
		"p":bytes,
	})

	//_, err = c.Writer.Write(bytes)
	//if err != nil {
	//	log.Println(err)
	//}
}
