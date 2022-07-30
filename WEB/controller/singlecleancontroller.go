package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ISingleCleanTaskController interface {
	SingleCreate(ctx *gin.Context)
}

type SingleCleanTaskController struct {
	DB *gorm.DB
}

func(t SingleCleanTaskController)  SingleCreate (ctx *gin.Context){

	var requestCleanTask vo.SingleCleanTaskRequest
	requestCleanTask.Applianceid=ctx.DefaultQuery("applianceid","1101")
	requestCleanTask.StartTime=ctx.DefaultQuery("starttime","2021-01-01 00:00:00")
	requestCleanTask.EndTime=ctx.DefaultQuery("endtime","2021-01-01 00:00:00")
	requestCleanTask.AdjustTaskFlag=ctx.DefaultQuery("adjust_task_flag","0")
	//数据验证
	//if err := ctx.ShouldBind(&requestCleanTask); err != nil {
	//	log.Println(err.Error())
	//	response.Fail(ctx, nil, "数据验证错误，分类名称必填")
	//	return
	//}
	//创建CleanTask
	singleCleanTask := model.SignleApplicance{
		Applianceid:  requestCleanTask.Applianceid,
		StartTime: requestCleanTask.StartTime,
		EndTime:   requestCleanTask.EndTime,
		AdjustTaskFlag: requestCleanTask.AdjustTaskFlag,
	}
	if err := t.DB.Create(&singleCleanTask).Error; err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"CleanTask": singleCleanTask}, "创建成功")

}
func NewsingleCleanTaskController() ISingleCleanTaskController {
	db := common.GetDB()
	result:=db.HasTable(&model.SignleApplicance{})
	if result==true{//若有 则不创建表
		fmt.Printf("表已存在\n" )
	}else
	{//若无 则创建表
		fmt.Printf("无表已创建\n" )
		db.AutoMigrate(&model.SignleApplicance{})
	}
	return SingleCleanTaskController{DB: db}
}
