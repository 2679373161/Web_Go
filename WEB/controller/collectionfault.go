package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IcollectionController interface{
	RestController


    Gettype(ctx *gin.Context)
	Getid(ctx *gin.Context)
    Getequ(ctx *gin.Context)
	Get_Change_Para_Info(ctx *gin.Context)

}
type collectionController struct {
	DB *gorm.DB
}
func (t collectionController) Create(ctx *gin.Context) {
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

func (t collectionController) Update(ctx *gin.Context) {
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

func (t collectionController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}


func (t collectionController) Gettype(ctx *gin.Context) {
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

	common.IndexDB.Raw("select distinct dev_type from failed_devices").Find(&tableStoreDates1)

	common.IndexDB.Raw("select distinct dev_id from failed_devices").Find(&tableStoreDates2) //设备号

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

func (t collectionController) Get_Change_Para_Info(ctx *gin.Context) {
	timeLow := ctx.DefaultQuery("timeLow", "")
	timeHigh := ctx.DefaultQuery("timeHigh", "")
	dev_id := ctx.DefaultQuery("dev_id", "")
	dev_change_state := ctx.DefaultQuery("dev_change_state", "")

	fields := []string{}
	values := []interface{}{}

	var Table_test_info []model.PerchangeNums

	if timeLow != "" {
		fields = append(fields, "updata_time >= ?", "updata_time <= ?")
		values = append(values, timeLow, timeHigh)
	}
	if dev_id != "" {
		fields = append(fields, "dev_id = ?")
		values = append(values, dev_id)
	}
	if dev_change_state != "" {
		fields = append(fields, "success_flag = ?")
		values = append(values, dev_change_state)
	}
	fmt.Println(dev_id)
	fmt.Println(dev_change_state)
	fmt.Println(values)
	fmt.Println(fields)
	common.IndexDB.Table("perchange_nums").Where(strings.Join(fields, " AND "), values...).Find(&Table_test_info)
	// tx.Find(&TopicList).RecordNotFound()

	response.Success(ctx, gin.H{"Para_data": Table_test_info}, "成功")
}

func (t collectionController) Getequ(ctx *gin.Context) {
	devtype := ctx.DefaultQuery("type", "")

	type Types2 struct {
		Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
	}

	type Type2 struct {
		Value string `json:"value"`
	}

	var tableStoreDates2 []Types2

	var type2 []Type2



	common.IndexDB.Raw("select distinct dev_id from failed_devices where dev_type=?",devtype).Find(&tableStoreDates2) //设备号


	for _, tableDate2 := range tableStoreDates2 {

		var tableStoreDates2 Type2

		tableStoreDates2.Value = tableDate2.Dev_id
		type2 = append(type2, tableStoreDates2)

	}

	fmt.Println(type2)


  response.Success(ctx,gin.H{"data":type2 },"读取成功")
}
func (t collectionController) Getid(ctx *gin.Context) {
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	devtype := ctx.DefaultQuery("type", "0")
	devid := ctx.DefaultQuery("id", "0")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	flag:= ctx.DefaultQuery("flag", "")
	//timeStr := "\"" + time1 + "\""
starttime=  starttime+" 00:00:00"
endtime=  endtime+" 23:59:59"
fmt.Println(starttime,endtime)

	var count int
	var idfault []model.Failed_devices
	fmt.Println(perPage,currentPage)

	if devtype=="" && devid==""{
		fmt.Println("666")
		common.IndexDB.Raw("select * from failed_devices where data_time>=?&&data_time<=?&& fault_type=?  ",starttime,endtime,flag).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
		common.IndexDB.Table("failed_devices").Where(" data_time>=?&&data_time<=?&& fault_type=?  ",starttime,endtime,flag).Count(&count)


	}else if devtype==""{
		common.IndexDB.Table("failed_devices").Where("dev_id=?&& data_time>=?&&data_time<=?&& fault_type=?  ",devid,starttime,endtime,flag).Count(&count)

		common.IndexDB.Raw("select * from failed_devices where dev_id=? &&data_time>=?&&data_time<=?&& fault_type=?  ",devid,starttime,endtime,flag).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else if devid==""{
		common.IndexDB.Table("failed_devices").Where("dev_type=?&& data_time>=?&&data_time<=?&& fault_type=?  ",devtype,starttime,endtime,flag).Count(&count)
		common.IndexDB.Raw("select * from failed_devices where dev_type=? &&data_time>=?&&data_time<=?&& fault_type=?  ",devtype,starttime,endtime,flag).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else {
		common.IndexDB.Table("failed_devices").Where("dev_typr=?&&dev_id=?&& data_time>=?&&data_time<=?&& fault_type=?  ",devtype,devid,starttime,endtime,flag).Count(&count)

		common.IndexDB.Raw("select * from failed_devices where dev_type=?&&dev_id=?&&data_time>=?&&data_time<=?&& fault_type=?  ",devtype,devid,starttime,endtime,flag).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}

	fmt.Println("idfault",idfault)
	//fmt.Println("count",count)
	//count=len(idfault)
	response.Success(ctx,gin.H{"data":idfault,"count":count },"读取成功")
}







func (t collectionController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}









	func NewcollectionController ()IcollectionController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return collectionController{DB:db}
}