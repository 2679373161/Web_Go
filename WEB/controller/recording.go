package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IrecordingController interface{
	RestController


    Gettype(ctx *gin.Context)
	Getid(ctx *gin.Context)
    Getequ(ctx *gin.Context)


}
type recordingController struct {
	DB *gorm.DB
}
func (t recordingController) Create(ctx *gin.Context) {
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

func (t recordingController) Update(ctx *gin.Context) {
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

func (t recordingController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}


func (t recordingController) Gettype(ctx *gin.Context) {
	type Types1 struct {
		Dev_type string `json:"dev_type" gorm:"type:varchar(255);not null"`
	}
	type Types2 struct {
		Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
	}
	type Type1 struct {
		Value string `json:"value"`
	}
	type Type2 struct {
		Value string `json:"value"`
	}
	var tableStoreDates1 []Types1
	var tableStoreDates2 []Types2
	var type1 []Type1
	var type2 []Type2

	common.IndexDB.Raw("select distinct dev_type from e_midea_fault").Find(&tableStoreDates1)

	common.IndexDB.Raw("select distinct dev_id from e_midea_fault").Find(&tableStoreDates2) //设备号

	for _, tableDate1 := range tableStoreDates1 {

		var tableStoreDates1 Type1

		tableStoreDates1.Value = tableDate1.Dev_type
		type1 = append(type1, tableStoreDates1)

	}
	fmt.Println(type1)
	for _, tableDate2 := range tableStoreDates2 {

		var tableStoreDates2 Type2

		tableStoreDates2.Value = tableDate2.Dev_id
		type2 = append(type2, tableStoreDates2)

	}

	fmt.Println(type2)

	response.Success(ctx, gin.H{"data": type1, "data1": type2}, "成功")
}


func (t recordingController) Getequ(ctx *gin.Context) {
	devtype := ctx.DefaultQuery("type", "")

	type Types2 struct {
		Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
	}

	type Type2 struct {
		Value string `json:"value"`
	}

	var tableStoreDates2 []Types2

	var type2 []Type2



	common.IndexDB.Raw("select distinct dev_id from e_midea_fault where dev_type=?",devtype).Find(&tableStoreDates2) //设备号


	for _, tableDate2 := range tableStoreDates2 {

		var tableStoreDates2 Type2

		tableStoreDates2.Value = tableDate2.Dev_id
		type2 = append(type2, tableStoreDates2)

	}

	fmt.Println(type2)


  response.Success(ctx,gin.H{"data":type2 },"读取成功")
}
func (t recordingController) Getid(ctx *gin.Context) {
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	devtype := ctx.DefaultQuery("type", "0")
	devid := ctx.DefaultQuery("id", "0")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")

starttime=starttime+" 00:00:00"
endtime=endtime+" 23:59:59"

	var count int
	var idfault []model.Real_mideaFault
	fmt.Println(perPage,currentPage)
	if devtype=="" && devid==""{
		fmt.Println("666")
		common.IndexDB.Raw("select * from e_midea_fault where start_time>=?&&end_time<=?  ",starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
		common.IndexDB.Table("e_midea_fault").Where(" start_time>=?&&end_time<=?  ",starttime,endtime).Count(&count)


	}else if devtype==""{
		common.IndexDB.Table("e_midea_fault").Where("dev_id=?&& start_time>=?&&end_time<=?  ",devid,starttime,endtime).Count(&count)

		common.IndexDB.Raw("select * from e_midea_fault where dev_id=? &&start_time>=?&&end_time<=?  ",devid,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else if devid==""{
		common.IndexDB.Table("e_midea_fault").Where("dev_type=?&& start_time>=?&&end_time<=?  ",devtype,starttime,endtime).Count(&count)
		common.IndexDB.Raw("select * from e_midea_fault where dev_type=? &&start_time>=?&&end_time<=?  ",devtype,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else {
		common.IndexDB.Table("e_midea_fault").Where("dev_typr=?&&dev_id=?&& start_time>=?&&end_time<=?  ",devtype,devid,starttime,endtime).Count(&count)

		common.IndexDB.Raw("select * from e_midea_fault where dev_type=?&&dev_id=?&&start_time>=?&&end_time<=?  ",devtype,devid,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}

	fmt.Println("idfault",idfault)
	//fmt.Println("count",count)
	//count=len(idfault)
	response.Success(ctx,gin.H{"data":idfault,"count":count },"读取成功")
}







func (t recordingController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}









	func NewrecordingController ()IrecordingController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return recordingController{DB:db}
}