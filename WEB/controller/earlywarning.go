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

type IEarlyWarningController interface{
	RestController

	Search(ctx *gin.Context)
	/*	Ceshi(ctx *gin.Context)*/
	ScoreChart(ctx *gin.Context)
}
type EarlyWarningController struct {
	DB *gorm.DB
}

func (t EarlyWarningController) Create(ctx *gin.Context) {
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

	}
	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"创建成功")

}

func (t EarlyWarningController) Update(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateDatasaveRequest
	//数据验证
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


	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"更新成功")


}

func (t EarlyWarningController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t EarlyWarningController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t EarlyWarningController) Search(ctx *gin.Context) {
	DataTime :=ctx.DefaultQuery("datatime", "")
	var TempScoreDownDev []model.TempScoreDownDevs
    common.IndexDB.Table("early_warning_devs").Where("end_time=?",DataTime).Find(&TempScoreDownDev)
	var City model.TableDate3
	for i:=0;i<len(TempScoreDownDev);i++{
		TempScoreDownDev[i].CityCode=TempScoreDownDev[i].City
		common.DB.Table("midea_loc_code").Where("city_code = ?", TempScoreDownDev[i].City).Find(&City)
		TempScoreDownDev[i].City=City.Dev_city
	}
	fmt.Println(TempScoreDownDev)
	response.Success(ctx, gin.H{"data":TempScoreDownDev}, "成功")

}
func (t EarlyWarningController)ScoreChart(ctx *gin.Context)  {
	equipment := ctx.DefaultQuery("equipment", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	CityCode:=ctx.DefaultQuery("city_code", "")
	fragmenttable:="fragment"+CityCode
	var fragment []model.Tablefragment
	var TempScoreDownDev model.TempScoreDownDevs
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Order("start_time ASC").Find(&fragment)
	common.IndexDB.Table("early_warning_devs").
		Where("dev_id=? AND start_time >= ? AND end_time <= ?",equipment,start_time, end_time).
		Find(&TempScoreDownDev)
	var score []int
	var TempScoreAvg []float64
	var UpperLimit   []float64
	var LowerLimit   []float64
	var datetime[]string
	for i:=0;i<len(fragment);i++{
		score=append(score,fragment[i].Temp_score)
		TempScoreAvg=append(TempScoreAvg,TempScoreDownDev.TempScoreAvg)
		UpperLimit=append(UpperLimit,TempScoreDownDev.UpperLimit)
		LowerLimit=append(LowerLimit,TempScoreDownDev.LowerLimit)
		datetime=append(datetime,fragment[i].Start_time[5:]+"至"+fragment[i].End_time[5:])
	}

	response.Success(ctx, gin.H{"score": score, "datetime": datetime,"tempscoreavg":TempScoreAvg,"upperlimit":UpperLimit,"lowerlimit":LowerLimit}, "成功")

}

func NewEarlyWarningController () IEarlyWarningController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return EarlyWarningController {DB:db}
}