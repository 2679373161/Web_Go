package common

import (
	"ginEssential/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey=[]byte("a_secret_creat")

type Claims struct{
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User)(string,error){
	expirationTime:=time.Now().Add(24*time.Hour)//token有效期
	claims:=&Claims{
		UserId: user.ID,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),//token有效时间
			IssuedAt:time.Now().Unix(),//token发放时间
			Issuer: "dongmingtao",//发放者
			Subject:"user token",

		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err:=token.SignedString(jwtKey)
	//token的三部分:使用加密协议+荷载+前两部分哈希值
	if err!=nil{
		return "",err
	}

	return tokenString,nil
}

func ParseToken(tokenString string)(*jwt.Token,*Claims,error){
	claims:=&Claims{}

	token,err:=jwt.ParseWithClaims(tokenString,claims, func(token *jwt.Token) (i interface{},err error) {
		return jwtKey,nil
	})
	return token,claims,err
}

