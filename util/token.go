package util

import (
	"github.com/dgrijalva/jwt-go"
	"netdisk/conf"
	"time"
)

var jwtSecret=[]byte(conf.JwtSecret)

type  Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 产生token的函数
func GenerateToken(username string)(string,error){
	nowTime :=time.Now()
	expireTime:=nowTime.Add(3*time.Hour)

	claims:=Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer: "net-disk",
		},
	}
	//
	tokenClaims:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	token,err:=tokenClaims.SignedString(jwtSecret)

	return token,err
}


// 验证token的函数
func ParseToken(token string)(*Claims,error){
	tokenClaims,err:=jwt.ParseWithClaims(token,&Claims{},func(token *jwt.Token)(interface{},error){
		return jwtSecret,nil
	})

	if tokenClaims!=nil{
		if claims,ok:=tokenClaims.Claims.(*Claims);ok && tokenClaims.Valid{
			return claims,nil
		}
	}
	//
	return nil,err
}