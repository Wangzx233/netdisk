package dao

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func SaveFile(fileName, dst, username string, fileSize int64) string {
	md5hash := createMD5(dst)
	err := DB.Create(&FileInfo{
		FileName:   fileName,
		FileSize:   fileSize,
		Md5hash:    md5hash,
		Path:       dst,
		Username:   username,
		Power:    0,
		Secret:     "",
		UploadTime: time.Now(),
	}).Error
	if err != nil {
		return "文件已存在"
	}
	return "上传成功"
}

// 根据文件内容生成MD5
func createMD5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println("model err : create md5 err:", err)
	}
	defer f.Close()

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		log.Println("model err : copy md5 err", err)
	}

	hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
