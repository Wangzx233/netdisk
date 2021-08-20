package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"netdisk/model"
	"netdisk/util"
	"time"
)

func Download(c *gin.Context) {
	token := c.GetHeader("token")
	md5hash := c.Query("md5hash")
	secret := c.Query("secret")
	if md5hash=="" {
		c.String(http.StatusBadRequest, "参数缺少")
		return
	}


	var allowDownload = false

	file := model.Download(md5hash)
	fmt.Println(file)
	switch file.Power {
	case -1:
		c.String(http.StatusBadRequest, file.Secret)
	case 0:
		claims, err := util.ParseToken(token)
		if err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest, "无下载权限")
			return
		} else {
			if claims.Username != file.Username {
				c.String(http.StatusBadRequest, "无下载权限")
				return
			} else {
				allowDownload = true
			}
		}
	case 1:
		if file.Secret!=secret {
			c.String(http.StatusBadRequest, "提取码错误")
			return
		}else {
			allowDownload=true
		}
	case 2:
		allowDownload=true
	}
	if allowDownload {
		c.Header("Content-Type", "application/octect-stream")
		c.Header("Content-Description", "attachment;filename=\""+file.FileName+"\"")

		// 限速下载
		bytes, _ := ioutil.ReadFile(file.Path)
		var offset int64
		buf :=make([]byte,1024)
		for  {
			if len(bytes[offset:])<1024 {
				buf=bytes[offset:]
				break
			}else {
				buf=bytes[offset:offset+1024]
			}
			offset+=1024
			c.Data(200,"application/x-gzip",buf)
			time.Sleep(time.Second)
		}

	}
}
