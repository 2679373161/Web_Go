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
	"strconv"
)

type Ifaultdiagnosisne01Controller interface{
	RestController

	Search(ctx *gin.Context)
	Summaries(ctx *gin.Context)
	Getequipment(ctx *gin.Context)
	Gettype(ctx *gin.Context)
	Getday(ctx *gin.Context)
	Ceshi(ctx *gin.Context)
	Overtemp(ctx *gin.Context)
	Overtempgetdata(ctx *gin.Context)
	Overtempgettype(ctx *gin.Context)
	Overtempgetid(ctx *gin.Context)
	Overtempday(ctx *gin.Context)
	Overtempdaygetdata(ctx *gin.Context)
	Overtempdaygettype(ctx *gin.Context)
	Overtempdaygetid(ctx *gin.Context)

}
type faultdiagnosisne01Controller struct {
	DB *gorm.DB
}
func (t faultdiagnosisne01Controller) Create(ctx *gin.Context) {
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

func (t faultdiagnosisne01Controller) Update(ctx *gin.Context) {
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

func (t faultdiagnosisne01Controller) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t faultdiagnosisne01Controller) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t faultdiagnosisne01Controller) Search(ctx *gin.Context) {



	
	fmt.Println("历史详情")
	//downloadflag:=ctx.DefaultQuery("downloadflag", "0")
	flag := ctx.DefaultQuery("flag", "1")
	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")
	devtype := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")

	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	

	var count int

    var downloadplace [] model.Fault_summaries
   
	if provincecode != "0" {
		provincecode = provincecode
	}
	if citycode != "0" {
		citycode = citycode
	}
	if len(devtype) != 0 {
		devtype = devtype
	}
	if len(equipment) != 0 {
		equipment = equipment
	}
	if len(starttime) != 0 && len(endtime) != 0 {
		starttime = starttime
		endtime = endtime
	}
	
var table string //表名
	//flag=0是E1，1是E0
	if flag== "0" {
		 table="E1_day"
	}else if flag=="1"{
		table="E0_day"
	}

	if citycode!="0"&&len(devtype)!=0&&len(equipment)!=0{

		common.IndexDB.Table(table).Where("province_code= ? and city_code=? and dev_type= ? and dev_id= ? and time_date >= ? and time_date <= ?  ", provincecode,citycode,devtype, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=? and city_code=?  and dev_type=? and dev_id=? and time_date >= ? and time_date <=?  ", provincecode, citycode,devtype, equipment, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}



	if citycode == "0" {
		common.IndexDB.Table(table).Where("province_code= ?  and dev_type= ? and dev_id= ? and time_date >= ? and time_date <= ?  ", provincecode, devtype, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and dev_type=? and dev_id=? and time_date >= ? and time_date <=?  ", provincecode, devtype, equipment, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if len(devtype) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?  and city_code= ? and dev_id= ? and time_date >= ? and time_date <= ?  ", provincecode, citycode, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and city_code=? and dev_id=? and time_date >= ? and time_date <=?  ", provincecode, citycode, equipment, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?  and city_code= ? and dev_type= ? and time_date >= ? and time_date <= ?  ", provincecode, citycode, devtype, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and city_code=? and dev_type=? and time_date >= ? and time_date <=?  ", provincecode, citycode, devtype, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" {
		common.IndexDB.Table(table).Where(" dev_type= ? and dev_id= ? and time_date >= ? and time_date <= ?  ", devtype, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("dev_type=? and dev_id=? and time_date >= ? and time_date <=?  ", devtype, equipment, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if citycode == "0" && (len(devtype) == 0) {
		common.IndexDB.Table(table).Where("province_code= ?  and dev_id= ? and time_date >= ? and time_date <= ?  ", provincecode, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?   and dev_id=? and time_date >= ? and time_date <=?  ", provincecode, equipment, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if citycode == "0" && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?  and dev_type= ?  and time_date >= ? and time_date <= ?  ", provincecode, devtype, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=? and  dev_type=?  and time_date >= ? and time_date <=?  ", provincecode, devtype, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" && len(devtype) == 0 {
		common.IndexDB.Table(table).Where("  dev_id= ? and time_date >= ? and time_date <= ?  ", equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where(" dev_id=? and time_date >= ? and time_date <=?  ", equipment, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("dev_type= ? and time_date >= ? and time_date <= ?  ", devtype, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("dev_type=? and time_date >= ? and time_date <=?  ", devtype, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ? and city_code= ?  and  time_date >= ? and time_date <= ?  ", provincecode, citycode, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and city_code= ? and time_date >= ? and time_date <=?  ", provincecode, citycode, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if citycode == "0" && len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?   and time_date >= ? and time_date <= ?  ", provincecode, starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and time_date >= ? and time_date <=?  ", provincecode, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" && len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("time_date >= ? and time_date <= ?  ", starttime, endtime).Count(&count)
		common.IndexDB.Table(table).Where("time_date >= ? and time_date <=?  ", starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)

	}

	fmt.Println("count=", count)
	fmt.Println("tableplace=", downloadplace)

	
	

	var downloadplace1 []model.TableDate3
	for _, tableDate := range downloadplace {
		var tableStoreDates4 model.TableDate3
		t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
		downloadplace1 = append(downloadplace1, tableStoreDates4)

	}
	

	fmt.Println("count=", count)

	response.Success(ctx, gin.H{"data":downloadplace, "data1":downloadplace1, "count": count}, "成功")
}


func (t faultdiagnosisne01Controller) Summaries(ctx *gin.Context) {
	fmt.Println("故障诊断")
	//downloadflag:=ctx.DefaultQuery("downloadflag", "0")
	flag := ctx.DefaultQuery("flag", "1")
	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")
	devtype := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	timedate := ctx.DefaultQuery("time_date", "")
	

	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))


	//var tableplace1  []model.Fault_summaries
	var count int
	//var abnormal_count int
	
	//common.IndexDB.Table("days_summaries").Where("dev_type = ? AND dev_id = ? AND time_date BETWEEN ? AND ?",dev_type,dev_id,timeLow,timeHigh).Find(&tableplace)
	var downloadplace [] model.Fault_summaries

	if provincecode != "0" {
		provincecode = provincecode
	}
	if citycode != "0" {
		citycode = citycode
	}
	if len(devtype) != 0 {
		devtype = devtype
	}
	if len(equipment) != 0 {
		equipment = equipment
	}
	if len(timedate) != 0  {
		timedate = timedate
		
	}

	var table string //表名
	//flag=0是E1，1是E0
	if flag== "0" {
		 table="E1_recent"
	}else if flag=="1"{
		table="E0_recent"
	}
	
	if citycode!="0"&&len(devtype)!=0&&len(equipment)!=0{

		common.IndexDB.Table(table).Where("province_code= ? city_code=? and dev_type= ? and dev_id= ? and time_date = ?  ", provincecode,citycode, devtype, equipment, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=? city_code=? and dev_type=? and dev_id=? and  time_date = ?   ", provincecode,citycode, devtype, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}

	if citycode == "0" {
		common.IndexDB.Table(table).Where("province_code= ?  and dev_type= ? and dev_id= ? and time_date = ?  ", provincecode, devtype, equipment, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and dev_type=? and dev_id=? and  time_date = ?   ", provincecode, devtype, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if len(devtype) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?  and city_code= ? and dev_id= ? and  time_date = ?   ", provincecode, citycode, equipment, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and city_code=? and dev_id=? and  time_date = ?   ", provincecode, citycode, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?  and city_code= ? and dev_type= ? and  time_date = ?   ", provincecode, citycode, devtype, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and city_code=? and dev_type=? and  time_date = ?   ", provincecode, citycode, devtype, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" {
		common.IndexDB.Table(table).Where(" dev_type= ? and dev_id= ? and  time_date = ?   ", devtype, equipment, timedate).Count(&count)
		common.IndexDB.Table(table).Where("dev_type=? and dev_id=? and time_date = ?   ", devtype, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if citycode == "0" && (len(devtype) == 0) {
		common.IndexDB.Table(table).Where("province_code= ?  and dev_id= ? and  time_date = ?   ", provincecode, equipment,timedate ).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?   and dev_id=? and  time_date = ?   ", provincecode, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if citycode == "0" && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?  and dev_type= ? and  time_date = ?   ", provincecode, devtype, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=? and dev_type=? and time_date = ?  ", provincecode, devtype, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" && len(devtype) == 0 {
		common.IndexDB.Table(table).Where("  dev_id= ? and  time_date = ?   ", equipment, timedate).Count(&count)
		common.IndexDB.Table(table).Where(" dev_id=? and  time_date = ?   ", equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("dev_type= ? and  time_date = ?   ", devtype,timedate).Count(&count)
		common.IndexDB.Table(table).Where("dev_type=? and  time_date = ?   ", devtype,timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ? and city_code= ?  and   time_date = ?   ", provincecode, citycode, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and city_code= ? and  time_date = ?   ", provincecode, citycode, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if citycode == "0" && len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table(table).Where("province_code= ?   and  time_date = ?   ", provincecode, timedate).Count(&count)
		common.IndexDB.Table(table).Where("province_code=?  and  time_date = ?   ", provincecode, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	}
	if provincecode == "0" && len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table(table).Where(" time_date = ?   ", timedate).Count(&count)
		common.IndexDB.Table(table).Where(" time_date = ?   ", timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)

	}



	fmt.Println("count=", count)
	fmt.Println("tableplace=", downloadplace)

	



	var downloadplace1 []model.TableDate3
	for _, tableDate := range downloadplace {
		var tableStoreDates4 model.TableDate3
		t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
		downloadplace1 = append(downloadplace1, tableStoreDates4)

	}


	fmt.Println("count=", count)

	response.Success(ctx, gin.H{"data":downloadplace, "data1":downloadplace1, "count": count}, "成功")
}



func (t faultdiagnosisne01Controller) Getequipment(ctx *gin.Context) {
	flag := ctx.DefaultQuery("flag", "0")
	province_code := ctx.DefaultQuery("province_code", "0")
	city_code := ctx.DefaultQuery("city_code", "0")
	dev_type := ctx.DefaultQuery("type", "0")
	type Equipment struct {
		Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
	}
	type Equipment1 struct {
		Value string `json:"value"`
	}
	var tableStoreDates4 []Equipment

	var type1 []Equipment1
	if flag == "3" {

		common.IndexDB.Raw("select distinct dev_id from Neu_fault ").Find(&tableStoreDates4)
		//common.IndexDB.Table("fault_summaries").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "1" {

		common.IndexDB.Raw("select distinct dev_id from Neu_fault where province_code=? and dev_type=?", province_code, dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("fault_summaries").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "2" {
		common.IndexDB.Raw("select distinct dev_id from Neu_fault where city_code=? and dev_type=?", city_code, dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("fault_summaries").Where("city_code = ? and dev_type = ? ", city_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else {
		fmt.Println(dev_type)
		common.IndexDB.Raw("select distinct dev_id from Neu_fault where dev_type=?", dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("fault_summaries").Where("dev_type = ? ", dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	}
}

func (t faultdiagnosisne01Controller) Ceshi(ctx *gin.Context) {
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

	common.IndexDB.Raw("select distinct dev_type from Neu_fault").Find(&tableStoreDates1)

	common.IndexDB.Raw("select distinct dev_id from Neu_fault").Find(&tableStoreDates2) //设备号

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

func (t faultdiagnosisne01Controller) Gettype(ctx *gin.Context) {
	flag := ctx.DefaultQuery("flag", "0")
	province_code := ctx.DefaultQuery("province_code", "0")
	city_code := ctx.DefaultQuery("city_code", "0")

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
	if flag == "1" {
		common.IndexDB.Raw("select distinct dev_type from Neu_fault where province_code = ? ", province_code).Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id from Neu_fault where province_code = ? ", province_code).Find(&tableStoreDates2) //设备号
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
	} else if flag == "2" {
		common.IndexDB.Raw("select distinct dev_type from Neu_fault where city_code = ?", city_code).Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id   from Neu_fault where city_code = ?", city_code).Find(&tableStoreDates2) //设备号
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

	} else {
		common.IndexDB.Raw("select distinct dev_type from Neu_fault").Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id from Neu_fault").Find(&tableStoreDates2) //设备号
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



}
func (t faultdiagnosisne01Controller) Getday(ctx *gin.Context) {
	fmt.Println("666")
	var day_infor [] model.Fault_summaries
	flag := ctx.DefaultQuery("flag", "0")
	equipment := ctx.DefaultQuery("equipment", "")


	if flag=="1"{
		common.IndexDB.Table("Neu_fault").Where(" dev_id=?", equipment).Find(&day_infor)
	}else if flag=="0"{
		common.IndexDB.Table("Neu_fault").Where(" dev_id=?",equipment).Find(&day_infor)

	}

	response.Success(ctx, gin.H{"data": day_infor}, "成功")

}
func (t faultdiagnosisne01Controller) Overtemp(ctx *gin.Context) {
	fmt.Println("故障诊断")
	//flag := ctx.DefaultQuery("flag", "1")
	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")
	devtype := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	timedate := ctx.DefaultQuery("time_date", "")
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	var count int
	var overtempdata [] model.Fault_overtemp

	if provincecode != "0" {
	   provincecode = provincecode
	}
	if citycode != "0" {
	   citycode = citycode
	}
	if len(devtype) != 0 {
	   devtype = devtype
	}
	if len(equipment) != 0 {
	   equipment = equipment
	}
	if len(timedate) != 0  {
	   timedate = timedate

	}

	if citycode!="0"&&len(devtype)!=0&&len(equipment)!=0{
	   common.IndexDB.Table("E4_recent").Where(" dev_id= ? and time_date = ?  ", equipment, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("dev_id=? and  time_date = ?   ", equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}

	if citycode == "0" {
	   common.IndexDB.Table("E4_recent").Where("province_code= ?  and dev_type= ? and dev_id= ? and time_date = ?  ", provincecode, devtype, equipment, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=?  and dev_type=? and dev_id=? and  time_date = ?   ", provincecode, devtype, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if len(devtype) == 0 {
	   common.IndexDB.Table("E4_recent").Where("province_code= ?  and city_code= ? and dev_id= ? and  time_date = ?   ", provincecode, citycode, equipment, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=?  and city_code=? and dev_id=? and  time_date = ?   ", provincecode, citycode, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if len(equipment) == 0 {
	   common.IndexDB.Table("E4_recent").Where("province_code= ?  and city_code= ? and dev_type= ? and  time_date = ?   ", provincecode, citycode, devtype, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=?  and city_code=? and dev_type=? and  time_date = ?   ", provincecode, citycode, devtype, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if provincecode == "0" {
	   common.IndexDB.Table("E4_recent").Where(" dev_type= ? and dev_id= ? and  time_date = ?   ", devtype, equipment, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("dev_type=? and dev_id=? and time_date = ?   ", devtype, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if citycode == "0" && (len(devtype) == 0) {
	   common.IndexDB.Table("E4_recent").Where("province_code= ?  and dev_id= ? and  time_date = ?   ", provincecode, equipment,timedate ).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=?   and dev_id=? and  time_date = ?   ", provincecode, equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if citycode == "0" && len(equipment) == 0 {
	   common.IndexDB.Table("E4_recent").Where("province_code= ?  and dev_type= ? and  time_date = ?   ", provincecode, devtype, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=? and dev_type=? and time_date = ?  ", provincecode, devtype, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if provincecode == "0" && len(devtype) == 0 {
	   common.IndexDB.Table("E4_recent").Where("  dev_id= ? and  time_date = ?   ", equipment, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where(" dev_id=? and  time_date = ?   ", equipment, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if provincecode == "0" && len(equipment) == 0 {
	   common.IndexDB.Table("E4_recent").Where("dev_type= ? and  time_date = ?   ", devtype,timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("dev_type=? and  time_date = ?   ", devtype,timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if len(devtype) == 0 && len(equipment) == 0 {
	   common.IndexDB.Table("E4_recent").Where("province_code= ? and city_code= ?  and   time_date = ?   ", provincecode, citycode, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=?  and city_code= ? and  time_date = ?   ", provincecode, citycode, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if citycode == "0" && len(devtype) == 0 && len(equipment) == 0 {
	   common.IndexDB.Table("E4_recent").Where("province_code= ?   and  time_date = ?   ", provincecode, timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where("province_code=?  and  time_date = ?   ", provincecode, timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	if provincecode == "0" && len(devtype) == 0 && len(equipment) == 0 {
	   common.IndexDB.Table("E4_recent").Where(" time_date = ?   ", timedate).Count(&count)
	   common.IndexDB.Table("E4_recent").Where(" time_date = ?   ", timedate).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&overtempdata)
	}
	fmt.Println("count=", count)
	fmt.Println("tableplace=", overtempdata)

	var downloadplace1 []model.TableDate3
	for _, tableDate := range overtempdata {
	   var tableStoreDates4 model.TableDate3
	   t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
	   downloadplace1 = append(downloadplace1, tableStoreDates4)
	}
	fmt.Println("count=", count)
	response.Success(ctx, gin.H{"data":overtempdata, "data1":downloadplace1, "count": count}, "成功")
 }
func (t faultdiagnosisne01Controller) Overtempgetdata(ctx *gin.Context) {
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

 common.IndexDB.Raw("select distinct dev_type from E4_recent").Find(&tableStoreDates1)
 common.IndexDB.Raw("select distinct dev_id from E4_recent").Find(&tableStoreDates2) //设备号
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

func (t faultdiagnosisne01Controller) Overtempgettype(ctx *gin.Context) {
 flag := ctx.DefaultQuery("flag", "0")
 province_code := ctx.DefaultQuery("province_code", "0")
 city_code := ctx.DefaultQuery("city_code", "0")

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
 if flag == "1" {
	common.IndexDB.Raw("select distinct dev_type from E4_recent where province_code = ? ", province_code).Find(&tableStoreDates1)
	common.IndexDB.Raw("select distinct dev_id from E4_recent where province_code = ? ", province_code).Find(&tableStoreDates2) //设备号
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
 } else if flag == "2" {
	common.IndexDB.Raw("select distinct dev_type from E4_recent where city_code = ?", city_code).Find(&tableStoreDates1)
	common.IndexDB.Raw("select distinct dev_id   from E4_recent where city_code = ?", city_code).Find(&tableStoreDates2) //设备号
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

 } else {
	common.IndexDB.Raw("select distinct dev_type from E4_recent").Find(&tableStoreDates1)
	common.IndexDB.Raw("select distinct dev_id from E4_recent").Find(&tableStoreDates2) //设备号
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
}
func (t faultdiagnosisne01Controller) Overtempgetid(ctx *gin.Context) {
 flag := ctx.DefaultQuery("flag", "0")
 province_code := ctx.DefaultQuery("province_code", "0")
 city_code := ctx.DefaultQuery("city_code", "0")
 dev_type := ctx.DefaultQuery("type", "0")
 type Equipment struct {
	Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
 }
 type Equipment1 struct {
	Value string `json:"value"`
 }
 var tableStoreDates4 []Equipment

 var type1 []Equipment1
 if flag == "3" {
	common.IndexDB.Raw("select distinct dev_id from E4_recent ").Find(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
	fmt.Println("11",tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}
	fmt.Println("12",type1)
	response.Success(ctx, gin.H{"data": type1}, "成功")
 } else if flag == "1" {
	common.IndexDB.Raw("select distinct dev_id from E4_recent where province_code=? and dev_type=?", province_code, dev_type).Find(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
	fmt.Println("11",tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}
	fmt.Println("444")
	fmt.Println("type1=",type1)
	response.Success(ctx, gin.H{"data": type1}, "成功")
 } else if flag == "2" {
	common.IndexDB.Raw("select distinct dev_id from E4_recent where city_code=? and dev_type=?", city_code, dev_type).Find(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("city_code = ? and dev_type = ? ", city_code, dev_type).Find(&tableStoreDates4)
	fmt.Println(tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}

	response.Success(ctx, gin.H{"data": type1}, "成功")
 } else {
	common.IndexDB.Raw("select distinct dev_id from E4_recent where dev_type=?", dev_type).Find(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("dev_type = ? ", dev_type).Find(&tableStoreDates4)
	fmt.Println("0",tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}

	response.Success(ctx, gin.H{"data": type1}, "成功")
 }
}
func (t faultdiagnosisne01Controller) Overtempday(ctx *gin.Context) {
 fmt.Println("故障诊断")
 //flag := ctx.DefaultQuery("flag", "1")
 provincecode := ctx.DefaultQuery("province_code", "0")
 citycode := ctx.DefaultQuery("city_code", "0")
 devtype := ctx.DefaultQuery("type", "")
 equipment := ctx.DefaultQuery("equipment", "")
 starttime := ctx.DefaultQuery("timeLow", "")
 endtime := ctx.DefaultQuery("timeHigh", "")
 timedate := ctx.DefaultQuery("timeHigh", "")
 perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
 currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))


 var count int
 var overtempdata [] model.Fault_overtemp

 if provincecode != "0" {
	provincecode = provincecode
 }
 if citycode != "0" {
	citycode = citycode
 }
 if len(devtype) != 0 {
	devtype = devtype
 }
 if len(equipment) != 0 {
	equipment = equipment
 }
 if len(timedate) != 0  {
	timedate = timedate

 }

 if citycode!="0"&&len(devtype)!=0&&len(equipment)!=0{
	common.IndexDB.Table("E4_day").Where("dev_id= ? and time_date <= ? and  time_date >= ? ",equipment, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("dev_id=? and  time_date >= ? and  time_date <= ? ", equipment, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }

 if citycode == "0" {
	common.IndexDB.Table("E4_day").Where("province_code= ?  and dev_type= ? and dev_id= ? and time_date >= ? and  time_date <= ?  ", provincecode, devtype, equipment, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=?  and dev_type=? and dev_id=? and  time_date >= ? and  time_date <= ?   ", provincecode, devtype, equipment, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if len(devtype) == 0 {
	common.IndexDB.Table("E4_day").Where("province_code= ?  and city_code= ? and dev_id= ? and  time_date >= ? and  time_date <= ?   ", provincecode, citycode, equipment, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=?  and city_code=? and dev_id=? and  time_date >= ? and  time_date <= ?   ", provincecode, citycode, equipment, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if len(equipment) == 0 {
	common.IndexDB.Table("E4_day").Where("province_code= ?  and city_code= ? and dev_type= ? and  time_date >= ? and  time_date <= ?   ", provincecode, citycode, devtype, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=?  and city_code=? and dev_type=? and  time_date >= ? and  time_date <= ?   ", provincecode, citycode, devtype, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if provincecode == "0" {
	common.IndexDB.Table("E4_day").Where(" dev_type= ? and dev_id= ? and  time_date >= ? and  time_date <= ?   ", devtype, equipment, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("dev_type=? and dev_id=? and time_date >= ? and  time_date <= ?   ", devtype, equipment, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if citycode == "0" && (len(devtype) == 0) {
	common.IndexDB.Table("E4_day").Where("province_code= ?  and dev_id= ? and  time_date >= ? and  time_date <= ?   ", provincecode, equipment,starttime,endtime ).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=?   and dev_id=? and  time_date >= ? and  time_date <= ?   ", provincecode, equipment, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if citycode == "0" && len(equipment) == 0 {
	common.IndexDB.Table("E4_day").Where("province_code= ?  and dev_type= ? and  time_date >= ? and  time_date <= ?   ", provincecode, devtype, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=? and dev_type=? and time_date >= ? and  time_date <= ?  ", provincecode, devtype, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if provincecode == "0" && len(devtype) == 0 {
	common.IndexDB.Table("E4_day").Where("  dev_id= ? and  time_date >= ? and  time_date <= ?   ", equipment, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where(" dev_id=? and  time_date >= ? and  time_date <= ?   ", equipment, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if provincecode == "0" && len(equipment) == 0 {
	common.IndexDB.Table("E4_day").Where("dev_type= ? and  time_date >= ? and  time_date <= ?   ", devtype,starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("dev_type=? and  time_date >= ? and  time_date <= ?   ", devtype,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if len(devtype) == 0 && len(equipment) == 0 {
	common.IndexDB.Table("E4_day").Where("province_code= ? and city_code= ?  and   time_date >= ? and  time_date <= ?   ", provincecode, citycode, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=?  and city_code= ? and  time_date >= ? and  time_date <= ?   ", provincecode, citycode, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if citycode == "0" && len(devtype) == 0 && len(equipment) == 0 {
	common.IndexDB.Table("E4_day").Where("province_code= ?   and  time_date >= ? and  time_date <= ?   ", provincecode, starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where("province_code=?  and  time_date >= ? and  time_date <= ?   ", provincecode, starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 if provincecode == "0" && len(devtype) == 0 && len(equipment) == 0 {
	common.IndexDB.Table("E4_day").Where(" time_date >= ? and  time_date <= ?   ", starttime,endtime).Count(&count)
	common.IndexDB.Table("E4_day").Where(" time_date >= ? and  time_date <= ?   ", starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&overtempdata)
 }
 fmt.Println("count=", count)
 fmt.Println("tableplace=", overtempdata)

 var downloadplace1 []model.TableDate3
 for _, tableDate := range overtempdata {
	var tableStoreDates4 model.TableDate3
	t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Scan(&tableStoreDates4)
	downloadplace1 = append(downloadplace1, tableStoreDates4)

 }
 fmt.Println("count=", count)
 response.Success(ctx, gin.H{"data":overtempdata, "data1":downloadplace1, "count": count}, "成功")
}

func (t faultdiagnosisne01Controller) Overtempdaygetdata(ctx *gin.Context) {
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

 common.IndexDB.Raw("select distinct dev_type from E4_day").Scan(&tableStoreDates1)
 common.IndexDB.Raw("select distinct dev_id from E4_day").Scan(&tableStoreDates2) //设备号
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

func (t faultdiagnosisne01Controller) Overtempdaygettype(ctx *gin.Context) {
 flag := ctx.DefaultQuery("flag", "0")
 province_code := ctx.DefaultQuery("province_code", "0")
 city_code := ctx.DefaultQuery("city_code", "0")

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
 if flag == "1" {
	common.IndexDB.Raw("select distinct dev_type from E4_day where province_code = ? ", province_code).Scan(&tableStoreDates1)
	common.IndexDB.Raw("select distinct dev_id from E4_day where province_code = ? ", province_code).Scan(&tableStoreDates2) //设备号
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
 } else if flag == "2" {
	common.IndexDB.Raw("select distinct dev_type from E4_day where city_code = ?", city_code).Scan(&tableStoreDates1)
	common.IndexDB.Raw("select distinct dev_id   from E4_day where city_code = ?", city_code).Scan(&tableStoreDates2) //设备号
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

 } else {
	common.IndexDB.Raw("select distinct dev_type from E4_day").Scan(&tableStoreDates1)
	common.IndexDB.Raw("select distinct dev_id from E4_day").Scan(&tableStoreDates2) //设备号
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
}
func (t faultdiagnosisne01Controller) Overtempdaygetid(ctx *gin.Context) {
 flag := ctx.DefaultQuery("flag", "0")
 province_code := ctx.DefaultQuery("province_code", "0")
 city_code := ctx.DefaultQuery("city_code", "0")
 dev_type := ctx.DefaultQuery("type", "0")
 type Equipment struct {
	Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
 }
 type Equipment1 struct {
	Value string `json:"value"`
 }
 var tableStoreDates4 []Equipment

 var type1 []Equipment1
 if flag == "3" {
	common.IndexDB.Raw("select distinct dev_id from E4_day ").Scan(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Scan(&tableStoreDates4)
	fmt.Println("11",tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}
	fmt.Println("12",type1)
	response.Success(ctx, gin.H{"data": type1}, "成功")
 } else if flag == "1" {
	common.IndexDB.Raw("select distinct dev_id from E4_day where province_code=? and dev_type=?", province_code, dev_type).Scan(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Scan(&tableStoreDates4)
	fmt.Println("11",tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}
	fmt.Println("444")
	fmt.Println("type1=",type1)
	response.Success(ctx, gin.H{"data": type1}, "成功")
 } else if flag == "2" {
	common.IndexDB.Raw("select distinct dev_id from E4_day where city_code=? and dev_type=?", city_code, dev_type).Scan(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("city_code = ? and dev_type = ? ", city_code, dev_type).Scan(&tableStoreDates4)
	fmt.Println(tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}

	response.Success(ctx, gin.H{"data": type1}, "成功")
 } else {
	common.IndexDB.Raw("select distinct dev_id from E4_day where dev_type=?", dev_type).Scan(&tableStoreDates4)
	//common.IndexDB.Table("fault_summaries").Where("dev_type = ? ", dev_type).Scan(&tableStoreDates4)
	fmt.Println("0",tableStoreDates4)
	for _, tableDate := range tableStoreDates4 {
	   var tableStoreDates3 Equipment1
	   tableStoreDates3.Value = tableDate.Dev_id
	   type1 = append(type1, tableStoreDates3)
	}

	response.Success(ctx, gin.H{"data": type1}, "成功")
 }
}



	func Newfaultdiagnosisne01Controller ()Ifaultdiagnosisne01Controller{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return faultdiagnosisne01Controller{DB:db}
}