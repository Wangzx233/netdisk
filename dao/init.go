package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB
func InitDB()  {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/netdisk?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("dao err : init err : open mysql err:",err)
	}
	db.AutoMigrate(&User{},&FileInfo{})

	DB=db
}

