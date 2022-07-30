package main

import (
	"ginEssential/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"os"
)



func main(){
	InitConfig()//读取配置
	db:=common.InitDB()//数据库初始化
	defer db.Close()//延迟关闭

	r := gin.Default()
	CollectRoute(r)
	//data_mining(r)
	port:=viper.GetString("server.backPort")
	if port!=""{
		panic(r.Run(":"+port))
	}
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}

func InitConfig()  {
	workDir,_:=os.Getwd()//读取工作路径
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir+"/config")
	err:=viper.ReadInConfig()
	if err!=nil{
		panic(err)
	}
}
//func data_mining(r *gin.Engine)*gin.Engine{
//	timeStamp:=
//	indexWriteFlag := true
//	dataWriteFlag := true
//	fragmentWriteFlag := true
//	////timeStamp := "15:02"
//	////
//	yearsFlag := 0
//	monthsFlag := 0
//	DayFlag := 1
//	timeNow := time.Now()
//	dataStartTime := timeNow.Add(-time.Duration(24) * time.Hour).Format("2006-01-02")
//	dataEndTime := timeNow.Format("2006-01-02")
//	////
//	gocron.Every(1).Day().At(timeStamp).Do(houduan.DataMining, dataStartTime, dataEndTime, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearsFlag, monthsFlag, DayFlag)
//
//	//fmt.Println(1)
//	//gocron.Clear()
//	//gocron.Every(1).Day().At(timeStamp).Do(task())
//	<-gocron.Start()
//}






