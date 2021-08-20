package dao

import "fmt"

func DeleteFile(username,md5hash string) (path string,err error) {

	var file FileInfo
	DB.Where("username=? and md5hash=?",username,md5hash).First(&file).Debug()
	fmt.Println(file)
	path=file.Path
	err = DB.Where("username=? and md5hash=?",username,md5hash).Delete(&file).Debug().Error
	return
}