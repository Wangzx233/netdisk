package dao

import (
	"log"
	"netdisk/conf"
	"os"
)

func ChangePath(username,md5hash,newPath string) (oldP,newP string,err error) {
	var file FileInfo
	err = DB.Where("username=? and md5hash=?", username, md5hash).First(&file).Error
	if err!=nil {
		return
	}
	// 确保路径存在
	if _, err = os.Stat(conf.FileSavePath+"/"+username+newPath); err != nil {
		// 不存在就创建
		if err = os.MkdirAll(conf.FileSavePath+"/"+username+newPath, os.ModePerm); err != nil {
			log.Println("alter_file err : Mkidir err",err)
		}
	}

	dst := conf.FileSavePath+"/"+username+newPath+file.FileName


	err = DB.Model(&file).Update("path", dst).Error

	return file.Path,dst,err
}

func Rename(username,md5hash,newName string) (oldP,newP string,err error) {
	var file FileInfo
	DB.Where("username=? and md5hash=?",username,md5hash).First(&file)

	newPath := file.Path[:len(file.Path)-len(file.FileName)]+newName
	err = DB.Where("username=? and md5hash=?", username, md5hash).First(&file).Error
	if err!=nil {
		return
	}
	err = DB.Model(file).Update("file_name", newName).Update("path",newPath).Error
	return file.Path,newPath,err
}