package model

import (
	"log"
	"netdisk/dao"
)

func Download(md5hash string) dao.FileInfo {
	file, err := dao.Download(md5hash)
	if err != nil {
		log.Println("dao download err:",err)
		return dao.FileInfo{Power: -1,Secret: "文件不存在"}
	}
	return file
}