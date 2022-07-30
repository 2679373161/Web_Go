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

type ICleanDateExistController interface {
	RestController
	PageList(ctx *gin.Context)
}

type CleanDateExistController struct {
	DB *gorm.DB
}

func (t CleanDateExistController) Create(ctx *gin.Context) {
	var requestCleanDateExist vo.CleanDateExistRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCleanDateExist); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建CleanDateExist
	CleanDateExist := model.TaskList{
		Citycode:  requestCleanDateExist.Citycode,
		StartTime: requestCleanDateExist.Citycode,
		EndTime:   requestCleanDateExist.Dataday,
	}
	if err := t.DB.Create(&CleanDateExist).Error; err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"CleanDateExist": CleanDateExist}, "创建成功")

}

func (t CleanDateExistController) Update(ctx *gin.Context) {
	var requestCleanDateExist vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCleanDateExist); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	//获取path中的id
	CleanDateExistId := ctx.Params.ByName("id") //从上下文中解析

	var CleanDateExist model.TableDate
	if t.DB.Where("id=?", CleanDateExistId).First(&CleanDateExist).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	//更新文章
	if err := t.DB.Model(&CleanDateExist).Update(requestCleanDateExist).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"CleanDateExist": CleanDateExist}, "更新成功")

}

func (t CleanDateExistController) Show(ctx *gin.Context) {
	CleanDateExistId := ctx.Params.ByName("id") //从上下文中解析

	var CleanDateExist model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", CleanDateExistId).First(&CleanDateExist).RecordNotFound() {
		response.Fail(ctx, nil, "数据不存在")
		return
	}

	response.Success(ctx, gin.H{"CleanDateExist": CleanDateExist}, "读取成功")
}

func (t CleanDateExistController) Delete(ctx *gin.Context) {

	var requestCleanDateExist vo.CleanDateExistRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCleanDateExist); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	// db := common.GetDB()
	// db.AutoMigrate(model.CleanDateExist{})

	DeleteCleanDateLog := model.CleanDateExist{
		Citycode: requestCleanDateExist.Citycode,
		Dataday:  requestCleanDateExist.Dataday,
	}
	// var CleanDateExistone model.CleanDateExist Where("citycode=?&&dataday", requestCleanDateExist.Citycode, requestCleanDateExist.Dataday).
	if t.DB.First(&DeleteCleanDateLog).RecordNotFound() {
		response.Fail(ctx, nil, "数据不存在")
		return
	}
	//Find(DeleteCleanDateLog). Where("citycode=?&&dataday", requestCleanDateExist.Citycode, requestCleanDateExist.Dataday).
	if err := t.DB.Delete(&DeleteCleanDateLog).Error; err != nil {
		panic(err)
	}
	t.DB.DropTable("data" + requestCleanDateExist.Citycode + "_" + requestCleanDateExist.Dataday) //清除表格

	response.Fail(ctx, gin.H{"CleanDateExist": DeleteCleanDateLog, "citycode": requestCleanDateExist.Citycode}, "删除成功")

}

func (t CleanDateExistController) PageList(ctx *gin.Context) {
	//获取参数
	category := ctx.DefaultQuery("category", "112150186047071")
	timeLow := ctx.DefaultQuery("timeLow", "2015-12-31")
	timeHigh := ctx.DefaultQuery("timeHigh", "3015-12-31")

	var datatime []string //时间轴
	var flow []int     //水流量
	var flame []string    //火焰反馈
	var settemp []string  //设定温度
	var outtemp []string  //输出温度

	//分页
	var CleanDateExists []model.TableDate //"112150186047071"
	t.DB.Table(category).Where("DataTime BETWEEN ? AND ?", timeLow, timeHigh).Find(&CleanDateExists)

	for _, tableDate := range CleanDateExists {
		datatime = append(datatime, tableDate.Datatime) //数据放入切片中
		flow = append(flow, tableDate.Flow)
		flame = append(flame, tableDate.Flame)
		settemp = append(settemp, tableDate.Settemp)
		outtemp = append(outtemp, tableDate.Outtemp)
	}

	response.Success(ctx, gin.H{"data": CleanDateExists, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp}, "成功")
}

func NewCleanDateExistController() ICleanDateExistController {
	db := common.GetDB()
	db.AutoMigrate(model.CleanDateExist{})
	return CleanDateExistController{DB: db}
}
