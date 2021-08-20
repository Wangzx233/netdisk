package dao

func Register(username,password string) string {
	var user User
	err := DB.Where(User{Username: username}).Find(&user).Error
	if err != nil {
		DB.Create(&User{
			Username: username,
			Password: password,
		})
		return "创建成功"
	}
	return "用户已存在"
}

func Login(username, password string) (string,bool) {
	var user User
	err := DB.Where("username=? and password=?",username,password).Find(&user).Error
	if err != nil {
		return "用户不存在或密码错误",false
	}
	return "",true
}