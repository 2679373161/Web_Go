package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB
var RunDB *gorm.DB
var IndexDB *gorm.DB
func InitDB()*gorm.DB{
	driverName:=viper.GetString("database.driverName")
	host:=viper.GetString("database.host")
	port:=viper.GetString("database.port")
	database:=viper.GetString("database.database")
	indexdatabase:=viper.GetString("database.indexdatabase")
	rundatabase:=viper.GetString("database.rundatabase")
	username:=viper.GetString("database.username")
	password:=viper.GetString("database.password")
	charset:=viper.GetString("database.charset")
	//loc:=viper.GetString("database.loc")
	args:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		username,
		password,
		host,
		port,
		database,
		charset)
	indexargs:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		username,
		password,
		host,
		port,
		indexdatabase,
		charset)
	runargs:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		username,
		password,
		host,
		port,
		rundatabase,
		charset)
	db,err:=gorm.Open(driverName,args)
	RunDB,_=gorm.Open(driverName,runargs)
	IndexDB,_=gorm.Open(driverName,indexargs)
	if err!=nil{
		panic("failed to connect database err: "+err.Error())
	}
	//db.AutoMigrate(&model.User{})//gorm自动创建数据表
	DB=db
	return db
}

func GetDB()*gorm.DB{
	return DB
}
func GetRunDB()*gorm.DB{
	return RunDB
}
func GetIndexDB()*gorm.DB{
	return IndexDB
}