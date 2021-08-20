package dao

import "time"

type FileInfo struct {
	FileName   string
	FileSize   int64
	Md5hash    string `gorm:"primary_key"` //文件唯一标识
	Path       string `gorm:"unique"`
	Username   string
	Power      int    //共享等级(0:只有自己可下载,1:需要密钥下载,2:所有获得链接的人均可下载）
	Secret     string //共享等级为1时的密钥
	UploadTime time.Time
}

type User struct {
	Username string `json:"username" gorm:"primary_key" form:"username"`
	Password string `json:"password" gorm:"type:varchar(64)" form:"password"`
}
