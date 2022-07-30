package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/houduan"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
)

var flag=0
var dataMinedate int
type IDatasaveandmineController interface{
	RestController
	Datasave(ctx *gin.Context)
}
type DatasaveandmineController struct {
	DB *gorm.DB
}
func (t DatasaveandmineController) Create(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateDatasaveRequest
	//数据验证
	if err:=ctx.ShouldBind(&requestTableStoreDate);err!=nil{
		log.Println(err.Error())
		response.Fail(ctx,nil,"数据验证错误，分类名称必填")
		return
	}
	//创建tablestoredate
	tableStoreDate:=model.TableData{
		//CategoryId: requestTableStoreDate.CategoryId,
		Id:      requestTableStoreDate.Id,
		Label:    requestTableStoreDate.Label,
		Value: requestTableStoreDate.Value,

		//Flow: requestTableStoreDate.Flow,
		//Model: requestTableStoreDate.Model,
	}
	if err:=t.DB.Create(&tableStoreDate).Error;err!=nil{
		panic(err)
		return
	}
	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"创建成功")

}
func (t DatasaveandmineController) Update(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateDatasaveRequest
	//数据验证
	var timeStamp string
	if err:=ctx.ShouldBind(&requestTableStoreDate);err!=nil{
		log.Println(err.Error())
		response.Fail(ctx,nil,"数据验证错误，分类名称必填")
		return
	}

	fmt.Println(ctx.Params)
	//获取path中的id
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析
	fmt.Println(tableStoreDateId)
	var tableStoreDate model.TableData
		if t.DB.Where("label=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
			Label:=ctx.DefaultQuery("Label","0000")
			Value:=ctx.DefaultQuery("Value","0000")
			fmt.Println("00")
			newUser:=model.TableData{
				Label: Label,
				Value:Value,
			}
			requestTableStoreDate.Id=""
			fmt.Println(requestTableStoreDate)

			t.DB.Create(&newUser)
			//response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"参数不存在")
			//	response.Fail(ctx,nil,"文章不存在")
			return
		}
		fmt.Println(tableStoreDate)
		requestTableStoreDate.Id=""
		fmt.Println(requestTableStoreDate)
		//更新文章
		if err:=t.DB.Model(&tableStoreDate).Update(requestTableStoreDate).Error;err!=nil{
			//panic(err)
			fmt.Println(err)
			response.Fail(ctx,nil,"更新失败")
			return
		}

	var tableStoreDate1 [] model.TableData

	t.DB.Table("table_data").Find(&tableStoreDate1)
	for _,tableDate :=range tableStoreDate1{
        if tableDate.Label=="timestamp"{
			timeStamp=tableDate.Value
			fmt.Print(timeStamp)
		}
		if tableDate.Label=="dataminedate"{
			dataMinedate,_=strconv.Atoi(tableDate.Value)
			fmt.Print(dataMinedate)
		}
		//dataparameter1[tableDate.Label]=tableDate.Value

	}
	fmt.Print(timeStamp)
	Flag,_:=strconv.ParseBool(ctx.DefaultQuery("Flag", "false"))
	start:= ctx.DefaultQuery("start", "2021-01-01")
	end:= ctx.DefaultQuery("end", "2021-07-01")
	data_mining(timeStamp,Flag,start,end)
	//fmt.Print(1)
	//response.Success(ctx,gin.H{"data":timeStamp,"tableStoreDate":tableStoreDate},"成功")
}

func (t DatasaveandmineController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t DatasaveandmineController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t DatasaveandmineController) Datasave(ctx *gin.Context){
	var timeStamp string
	var tableStoreDate1 [] model.TableData

	t.DB.Table("table_data").Find(&tableStoreDate1)
	for _,tableDate :=range tableStoreDate1{
		if tableDate.Label=="timestamp"{
			timeStamp=tableDate.Value
			fmt.Print(timeStamp)
		}
		//dataparameter1[tableDate.Label]=tableDate.Value

	}
	fmt.Print(timeStamp)

	response.Success(ctx,gin.H{"data":timeStamp},"成功")
}
func task1(dataStartTime string,dataEndTime string,indexWriteFlag bool,dataWriteFlag bool,fragmentWriteFlag bool,yearsFlag int,monthsFlag int,DayFlag int,AutoFlag bool,applianceIdTableName string,Database string){
    fmt.Println(dataMinedate)
	if time.Now().Day()!=dataMinedate{
		fmt.Println(time.Now().Day())
		flag=0
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		houduan.DataMining(dataStartTime, dataEndTime, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearsFlag, monthsFlag, DayFlag, applianceIdTableName, Database)
		//日指标汇总
		houduan.DaySummary(dataStartTime, dataEndTime, AutoFlag, applianceIdTableName, Database)
		//月指标汇总
		houduan.MonthsSummary1(dataStartTime, dataEndTime, AutoFlag, applianceIdTableName, Database)
		//日省份、城市、型号汇总
		houduan.Summary(dataStartTime, dataEndTime, AutoFlag, true, Database)
		//月省份、城市、型号汇总
		houduan.Summary(dataStartTime, dataEndTime, AutoFlag, false, Database)
		db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/test1dev?charset=utf8mb4")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		var tableStoreDate model.TableData
		var requestTableStoreDate model.TableData
		db.Table("table_data").Where("label=?","dataminedate").First(&tableStoreDate)
		requestTableStoreDate.Value=strconv.Itoa(time.Now().Day())
		if err:=db.Table("table_data").Model(&tableStoreDate).Update(requestTableStoreDate).Error;err!=nil{
			//panic(err)
			fmt.Println(err)
			//response.Fail(ctx,nil,"更新失败")
			return
		}
		flag=1
	}
}

func data_mining(timeStamp string,AutoFlag bool,start string,end string){

	//fmt.Println(timeStamp)
	//gocron.Clear()
	//gocron.Every(1).Day().At(timeStamp).Do(task1)
	//<-gocron.Start()

	//设备表名
	var applianceIdTableName = "bo"
	//数据库信息
	var Database = "root:123456@(127.0.0.1:3306)/table_test?charset=utf8mb4"

	gocron.Clear()
	fmt.Print("11111111111111111111")
	if AutoFlag == true {
		indexWriteFlag := true
		dataWriteFlag := true
		fragmentWriteFlag := true
		//timeStamp := "15:02"

		yearsFlag := 0
		monthsFlag := 0
		DayFlag := 1
		timeNow := time.Now()
		dataStartTime := timeNow.Add(-time.Duration(24) * time.Hour).Format("2006-01-02")
		dataEndTime := timeNow.Format("2006-01-02")


		gocron.Every(1).Day().At(timeStamp).Do(task1,dataStartTime, dataEndTime, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearsFlag, monthsFlag, DayFlag, AutoFlag, applianceIdTableName, Database)
		<-gocron.Start()
		//houduan.DataMining(dataStartTime, dataEndTime, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearsFlag, monthsFlag, DayFlag)
		////日指标汇总
		//houduan.DaySummary(dataStartTime, dataEndTime, AutoFlag)
		////月指标汇总
		//houduan.MonthsSummary1(dataStartTime, dataEndTime, AutoFlag)
		////日省份、城市、型号汇总
		//houduan.Summary(dataStartTime, dataEndTime, AutoFlag, true)
		////月省份、城市、型号汇总
		//houduan.Summary(dataStartTime, dataEndTime, AutoFlag, false)
	} else {
		fmt.Print("11111111111111111111")
		//dataStartTime := "2021-01-01"
		//dataEndTime := "2021-07-01"
		dataStartTime := start
		dataEndTime := end
		indexWriteFlag := true
		dataWriteFlag := true
		fragmentWriteFlag := true
		yearsFlag := 0
		monthsFlag := 1
		DayFlag := 0

		houduan.DataMining(dataStartTime, dataEndTime, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearsFlag, monthsFlag, DayFlag, applianceIdTableName, Database)
		//日指标汇总
		houduan.DaySummary(dataStartTime, dataEndTime, AutoFlag, applianceIdTableName, Database)
		//月指标汇总
		houduan.MonthsSummary1(dataStartTime, dataEndTime, AutoFlag, applianceIdTableName, Database)
		//日省份、城市、型号汇总
		houduan.Summary(dataStartTime, dataEndTime, AutoFlag, true, Database)
		//月省份、城市、型号汇总
		houduan.Summary(dataStartTime, dataEndTime, AutoFlag, false, Database)
	}
}

func NewDatasaveandmineController ()IDatasaveandmineController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return DatasaveandmineController{DB:db}
}