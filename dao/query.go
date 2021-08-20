package dao

func Query(username string) ([]FileInfo,error) {
	var files []FileInfo
	err := DB.Where(FileInfo{Username: username}).Find(&files).Error

	return files,err
}

func Download(md5hash string) (FileInfo,error) {
	var file FileInfo
	err := DB.Where(FileInfo{Md5hash: md5hash}).Find(&file).Error
	return file,err
}