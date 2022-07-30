package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ICleanTaskController interface {
	RestController
	PageList(ctx *gin.Context)
}

type CleanTaskController struct {
	DB *gorm.DB
}

func (t CleanTaskController) Create(ctx *gin.Context) {
	var requestCleanTask vo.CleanTaskRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCleanTask); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建CleanTask
	CleanTask := model.TaskList{
		Citycode:  requestCleanTask.Citycode,
		StartTime: requestCleanTask.Starttime,
		EndTime:   requestCleanTask.Endtime,
	}
	if err := t.DB.Create(&CleanTask).Error; err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"CleanTask": CleanTask}, "创建成功")

}

func (t CleanTaskController) Update(ctx *gin.Context) {
	var requestCleanTask vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCleanTask); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	//获取path中的id
	CleanTaskId := ctx.Params.ByName("id") //从上下文中解析

	var CleanTask model.TableDate
	if t.DB.Where("id=?", CleanTaskId).First(&CleanTask).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	//更新文章
	if err := t.DB.Model(&CleanTask).Update(requestCleanTask).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"CleanTask": CleanTask}, "更新成功")

}

func (t CleanTaskController) Show(ctx *gin.Context) {
	CleanTaskId := ctx.Params.ByName("id") //从上下文中解析

	var CleanTask model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", CleanTaskId).First(&CleanTask).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"CleanTask": CleanTask}, "读取成功")
}

func (t CleanTaskController) Delete(ctx *gin.Context) {
	CleanTaskId := ctx.Params.ByName("id") //从上下文中解析

	var CleanTask model.TableDate
	if t.DB.Where("id=?", CleanTaskId).First(&CleanTask).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	t.DB.Delete(&CleanTask)

	response.Fail(ctx, gin.H{"CleanTask": CleanTask}, "删除成功")

}

func (t CleanTaskController) PageList(ctx *gin.Context) {
	//获取参数
	category := ctx.DefaultQuery("category", "112150186047071")
	timeLow := ctx.DefaultQuery("timeLow", "2015-12-31")
	timeHigh := ctx.DefaultQuery("timeHigh", "3015-12-31")
   // city:=ctx.DefaultQuery("city", "")
	var datatime []string //时间轴
	var flow []int     //水流量
	var flame []string    //火焰反馈
	var settemp []string  //设定温度
	var outtemp []string  //输出温度
    fmt.Println("666")
	//分页
	var CleanTasks []model.TableDate //"112150186047071"
	//var typetable  []model.TableDate
	t.DB.Table(category).Where("DataTime BETWEEN ? AND ?", timeLow, timeHigh).Find(&CleanTasks)
    //common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow, timeHigh, city).Find(&typetable)
	for _, tableDate := range CleanTasks {
		datatime = append(datatime, tableDate.Datatime) //数据放入切片中
		flow = append(flow, tableDate.Flow)
		flame = append(flame, tableDate.Flame)
		settemp = append(settemp, tableDate.Settemp)
		outtemp = append(outtemp, tableDate.Outtemp)
	}

	response.Success(ctx, gin.H{"data": CleanTasks, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp}, "成功")
}

func NewCleanTaskController() ICleanTaskController {
	db := common.GetDB()
	db.AutoMigrate(model.TaskList{})
	return CleanTaskController{DB: db}
}
