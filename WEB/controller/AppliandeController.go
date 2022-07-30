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

type IApplianceController interface {
	RestController
	PageList(ctx *gin.Context)
}

type ApplianceController struct {
	DB *gorm.DB
}

func (t ApplianceController) Create(ctx *gin.Context) {
	var requestAppliance vo.ApplianceSelectRequest
	//数据验证
	if err := ctx.ShouldBind(&requestAppliance); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建Appliance
	Appliance := model.ApplianceSelect{
		DevId:    requestAppliance.DevId,
		CityCode: requestAppliance.CityCode,
		Model:    requestAppliance.Model,
	}
	if err := t.DB.Create(&Appliance).Error; err != nil {
		panic(err)
	}
	response.Success(ctx, gin.H{"Appliance": Appliance}, "创建成功")

}

func (t ApplianceController) Update(ctx *gin.Context) {
	var requestAppliance vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestAppliance); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	//获取path中的id
	ApplianceId := ctx.Params.ByName("id") //从上下文中解析

	var Appliance model.TableDate
	if t.DB.Where("id=?", ApplianceId).First(&Appliance).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	//更新文章
	if err := t.DB.Model(&Appliance).Update(requestAppliance).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"Appliance": Appliance}, "更新成功")

}

func (t ApplianceController) Show(ctx *gin.Context) {
	ApplianceId := ctx.Params.ByName("id") //从上下文中解析

	var Appliance model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", ApplianceId).First(&Appliance).RecordNotFound() {
		response.Fail(ctx, nil, "数据不存在")
		return
	}

	response.Success(ctx, gin.H{"Appliance": Appliance}, "读取成功")
}

func (t ApplianceController) Delete(ctx *gin.Context) {

	var requestAppliance vo.ApplianceSelectRequest
	//数据验证
	if err := ctx.ShouldBind(&requestAppliance); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	// db := common.GetDB()
	// db.AutoMigrate(model.Appliance{})

	DeleteCleanDateLog := model.ApplianceSelect{
		CityCode: requestAppliance.CityCode,
		DevId:    requestAppliance.DevId,
		Model:    requestAppliance.Model,
	}
	// var Applianceone model.Appliance Where("citycode=?&&dataday", requestAppliance.Citycode, requestAppliance.Dataday).
	if t.DB.First(&DeleteCleanDateLog).RecordNotFound() {
		response.Fail(ctx, nil, "数据不存在")
		return
	}
	//Find(DeleteCleanDateLog). Where("citycode=?&&dataday", requestAppliance.Citycode, requestAppliance.Dataday).
	if err := t.DB.Delete(&DeleteCleanDateLog).Error; err != nil {
		panic(err)
	}
	// t.DB.DropTable("data" + requestAppliance.CityCode + "_" + requestAppliance.Dataday) //清除表格

	response.Fail(ctx, gin.H{"Appliance": DeleteCleanDateLog, "citycode": requestAppliance.CityCode}, "删除成功")

}

func (t ApplianceController) PageList(ctx *gin.Context) {
	//获取参数
	cityCode := ctx.DefaultQuery("citycode", "1101")
	// timeLow := ctx.DefaultQuery("timeLow", "2015-12-31")
	// timeHigh := ctx.DefaultQuery("timeHigh", "3015-12-31")
	var Appliances []model.ApplianceSelect

	var applianceid []string
	var citycode []string
	var model []string
	var devselect []string

	//分页
	// var Appliances []model.ApplianceSelect
	t.DB.Table("appliance_selects").Where("city_code=?", cityCode).Find(&Appliances)

	for _, tableDate := range Appliances {
		applianceid = append(applianceid, tableDate.DevId) //数据放入切片中
		citycode = append(citycode, tableDate.CityCode)
		model = append(model, tableDate.Model)
		devselect = append(devselect, tableDate.Select)
	}

	response.Success(ctx, gin.H{"data": Appliances, "applianceid": applianceid, "citycode": citycode, "model": model, "devselect": devselect}, "成功")
}

func NewApplianceController() IApplianceController {
	db := common.GetDB()
	db.AutoMigrate(model.ApplianceSelect{})
	return ApplianceController{DB: db}
}
