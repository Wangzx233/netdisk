package dao

import (
	"log"
	"strconv"
)

func Share(md5hash,power,username,secret string) error {
	p, err := strconv.Atoi(power)
	if err != nil {
		log.Println(err)
		return err
	}
	var file FileInfo
	err = DB.Where("username=? and md5hash=?", username, md5hash).First(&file).Error
	if err!=nil {
		return err
	}
	err = DB.Model(&file).Update("power", p,"secret",secret).Error
	return err
}