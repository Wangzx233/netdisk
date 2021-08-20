package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"io"
	"log"
	"net/http"
	"netdisk/conf"
	"netdisk/dao"
	"netdisk/util"
	"os"
	"path/filepath"
	"time"
)

var Client *redis.Client

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
func Upload(c *gin.Context) {
    InitRedis()
	//获取token
	token := c.GetHeader("token")
	//cookie, _ := c.Request.Cookie("user_cookie")
	cliams, err := util.ParseToken(token)
	if err != nil {
		log.Println("controller err : parseToken err :", err)
		c.String(http.StatusBadRequest, "获取token失败")
	}

	offset := 0
	//上传文件
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		fmt.Println("controller err : FormFile err :", err)
	}

	// 确保路径存在
	if _, err = os.Stat(conf.FileSavePath + "/" + cliams.Username); err != nil {
		// 不存在就创建
		if err = os.MkdirAll(conf.FileSavePath+"/"+cliams.Username, os.ModePerm); err != nil {
			log.Println("controller err : Mkidir err", err)
		}
	}

	dst := filepath.Join(conf.FileSavePath+"/"+cliams.Username, header.Filename)
	saveFile, err := os.OpenFile(dst, os.O_RDWR, os.ModePerm)
	if err != nil {
		saveFile, _ = os.Create(dst)
	}
	// 创建缓存区
	buf := make([]byte, 8192)
	// 检测redis里是否有上传记录
	res, err := Client.Get(dst).Int64()
	if err == nil {
		offset = int(res)

		// 设置上传文件seek，从文件头开始偏移
		file.Seek(int64(offset), io.SeekStart)

		fmt.Println("上次已传：", offset)
	}

	// 循环读取文件
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		// 从offset开始写起
		saveFile.WriteAt(buf, int64(offset))


		// 已上传的大小保存到redis
		offset += read
		fmt.Println(offset)
		err = Client.Set(dst, offset, time.Duration(0)).Err()
		if err != nil {
			log.Println("redis err: ", err)
		}

	}
	err = saveFile.Close()
	if err != nil {
		log.Println(err)
	}

	ans := dao.SaveFile(header.Filename, dst, cliams.Username, header.Size)
	c.String(http.StatusOK, ans)
}
