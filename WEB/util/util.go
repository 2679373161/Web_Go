package util

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(n int)string{
	var letters=[]byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result:=make([]byte,n)
	rand.Seed(time.Now().Unix())//随机数种子，保证rand.Intn是随机的
	for i:=range result{
		result[i]=letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// TimeParse 时间格式解析转换函数 string->int64
func TimeParse(timeString string)(timeInt int64){
	timeParse, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, time.Local)
	if err != nil{
		fmt.Println(err)
		return
	}
	timeInt=timeParse.Unix()//时间格式转化为时间戳（s）类型
	return timeInt
}