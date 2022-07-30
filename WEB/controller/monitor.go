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

type IMonitorController interface {
	RestController
	Datasave(ctx *gin.Context)
	Getequipment(ctx *gin.Context)
	Gettype(ctx *gin.Context)
	Datachart(ctx *gin.Context)
	Ceshi(ctx *gin.Context)
	Data_Error(ctx *gin.Context)
}
type MonitorController struct {
	DB *gorm.DB
}

func (t MonitorController) Create(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateDatasaveRequest
	//数据验证
	if err := ctx.ShouldBind(&requestTableStoreDate); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建tablestoredate
	tableStoreDate := model.TableData{
		//CategoryId: requestTableStoreDate.CategoryId,
		Id:    requestTableStoreDate.Id,
		Label: requestTableStoreDate.Label,
		Value: requestTableStoreDate.Value,

		//Flow: requestTableStoreDate.Flow,
		//Model: requestTableStoreDate.Model,
	}
	if err := t.DB.Create(&tableStoreDate).Error; err != nil {
		panic(err)
		return
	}
	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "创建成功")

}

func (t MonitorController) Update(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateDatasaveRequest
	//数据验证
	if err := ctx.ShouldBind(&requestTableStoreDate); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	fmt.Println(ctx.Params)
	//获取path中的id
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析
	fmt.Println(tableStoreDateId)
	var tableStoreDate model.TableData
	if t.DB.Where("label=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		Label := ctx.DefaultQuery("Label", "0000")
		Value := ctx.DefaultQuery("Value", "0000")
		fmt.Println("00")
		newUser := model.TableData{
			Label: Label,
			Value: Value,
		}
		requestTableStoreDate.Id = ""
		fmt.Println(requestTableStoreDate)

		t.DB.Create(&newUser)
		//response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"参数不存在")
		//	response.Fail(ctx,nil,"文章不存在")
		return
	}
	fmt.Println(tableStoreDate)
	requestTableStoreDate.Id = ""
	fmt.Println(requestTableStoreDate)
	//更新文章
	if err := t.DB.Model(&tableStoreDate).Update(requestTableStoreDate).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "更新成功")

}

func (t MonitorController) Show(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "读取成功")
}

func (t MonitorController) Delete(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx, gin.H{"tableStoreDate": tableStoreDate}, "删除成功")

}
func (t MonitorController) Data_Error(ctx *gin.Context) {

	downloadflag := ctx.DefaultQuery("downloadflag", "0") // 下载数据使用

	 flag := ctx.DefaultQuery("flag", "1") // 异常标志位

	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")

	devtype := ctx.DefaultQuery("type", "")

	equipment := ctx.DefaultQuery("equipment", "")

	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")

	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))

	var tableplace []model.Daily_monitorings_error

	var count int

	var tableStoreDates3 []model.TableDate3

	fields := []string{}
	values := []interface{}{}

	if provincecode != "0" {
		fields = append(fields, "province_code = ?")
		values = append(values, provincecode)
	}
	if citycode != "0" {
		fields = append(fields, "city_code = ?")
		values = append(values, citycode)
	}
	if len(devtype) != 0 {
		fields = append(fields, "dev_type = ?")
		values = append(values, devtype)
	}
	if len(equipment) != 0 {
		fields = append(fields, "dev_id = ?")
		values = append(values, equipment)
	}
	if len(starttime) != 0 && len(endtime) != 0 {
		
		fields = append(fields, "time_date >= ?", "time_date <= ?")
		values = append(values, starttime, endtime)
	}
	if flag=="1"{
		fields = append(fields, "abnormalflag= ?")
		values = append(values, 1)
	}
	if flag=="0"{
		fmt.Println("666")
		fields = append(fields, "abnormalflag= ?")
		values = append(values, 0)
	}
	//fields = append(fields, "temp_score != ?")
	//values = append(values, 0)
	fmt.Println(fields,values)
	common.IndexDB.Table("daily_monitorings").Where(strings.Join(fields, " AND "), values...).Count(&count)
	if downloadflag == "1" {
		common.IndexDB.Table("daily_monitorings").Where(strings.Join(fields, " AND "), values...).Find(&tableplace)
	} else {
		common.IndexDB.Table("daily_monitorings").Where(strings.Join(fields, " AND "), values...).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
	}

	fmt.Println("count=", count)
	fmt.Println("tableplace=", tableplace)

		var tableStoreDates4 model.TableDate3
	for _, tableDate := range tableplace {

		t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)

		tableStoreDates3 = append(tableStoreDates3, tableStoreDates4)

	}
	fmt.Println(tableplace)

	fmt.Println("count=", count)

	response.Success(ctx, gin.H{"data": tableplace, "count": count, "data1": tableStoreDates3}, "成功")
}



func (t MonitorController) Datasave(ctx *gin.Context) {



	//var tableselect [] model.Tableselect
	//if err:=t.DB.Table("bo").Model(&tableselect).Where("dev_type = ?","00016HGS").Update("monitor_flag", 0).Error;err!=nil{
	//	//panic(err)
	//	fmt.Println(err)
	//	response.Fail(ctx,nil,"更新失败")
	//	return
	//}
	downloadflag:=ctx.DefaultQuery("downloadflag", "0")
	flag := ctx.DefaultQuery("flag", "1")
	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")
	devtype := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")

	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))

	var tableplace []model.Daily_monitorings
	//var tableplace1  []model.Daily_monitorings
	var count int
	//var abnormal_count int
	var tableStoreDates3 []model.TableDate3
	//common.IndexDB.Table("days_summaries").Where("dev_type = ? AND dev_id = ? AND time_date BETWEEN ? AND ?",dev_type,dev_id,timeLow,timeHigh).Find(&tableplace)

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
	//select_option := make(map[string]interface{})
	//if province_code != "0" {
	//	select_option["province_code"] = province_code
	//}
	//if city_code != "0" {
	//	select_option["city_code"] = city_code
	//}
	//if len(dev_type) != 0 {
	//	select_option["dev_type"] = dev_type
	//}
	//if len(equipment) != 0 {
	//	select_option["dev_id"] = equipment
	//}
	//if len(starttime)!=0&&len(endtime) !=0{
	//	select_option["starttime"] = starttime
	//	select_option["endtime"]=endtime
	//}

	//common.IndexDB.Table("daily_monitorings").Where("province_code= ? and city_code= ? and dev_type= ? and distinct dev_id= ? and time_date >= ? and time_date <= ?", provincecode,citycode,devtype,equipment,starttime,endtime).Count(&count)
	//common.IndexDB.Table("daily_monitorings").Where("province_code=? and city_code=? and dev_type=? and dev_id=? and time_date >= ? and time_date <=?", provincecode,citycode,devtype,equipment,starttime,endtime).Offset(perPage*(currentPage-1)).Limit(perPage).Find(&tableplace)

	if citycode == "0" {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ?  and dev_type= ? and dev_id= ? and time_date >= ? and time_date <= ?", provincecode, devtype, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=?  and dev_type=? and dev_id=? and time_date >= ? and time_date <=?", provincecode, devtype, equipment, starttime, endtime).Find(&tableplace)
	}
	if len(devtype) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ?  and city_code= ? and dev_id= ? and time_date >= ? and time_date <= ?", provincecode, citycode, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=?  and city_code=? and dev_id=? and time_date >= ? and time_date <=?", provincecode, citycode, equipment, starttime, endtime).Find(&tableplace)
	}
	if len(equipment) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ?  and city_code= ? and dev_type= ? and time_date >= ? and time_date <= ?", provincecode, citycode, devtype, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=?  and city_code=? and dev_type=? and time_date >= ? and time_date <=?", provincecode, citycode, devtype, starttime, endtime).Find(&tableplace)
	}
	if provincecode == "0" {
		common.IndexDB.Table("daily_monitorings").Where(" dev_type= ? and dev_id= ? and time_date >= ? and time_date <= ?", devtype, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("dev_type=? and dev_id=? and time_date >= ? and time_date <=?", devtype, equipment, starttime, endtime).Find(&tableplace)
	}
	if citycode == "0" && (len(devtype) == 0) {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ?  and dev_id= ? and time_date >= ? and time_date <= ?", provincecode, equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=?   and dev_id=? and time_date >= ? and time_date <=?", provincecode, equipment, starttime, endtime).Find(&tableplace)
	}
	if citycode == "0" && len(equipment) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ?  and dev_type= ?  and time_date >= ? and time_date <= ?", provincecode, devtype, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=? and  dev_type=?  and time_date >= ? and time_date <=?", provincecode, devtype, starttime, endtime).Find(&tableplace)
	}
	if provincecode == "0" && len(devtype) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("  dev_id= ? and time_date >= ? and time_date <= ?", equipment, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where(" dev_id=? and time_date >= ? and time_date <=?", equipment, starttime, endtime).Find(&tableplace)
	}
	if provincecode == "0" && len(equipment) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("dev_type= ? and time_date >= ? and time_date <= ?", devtype, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("dev_type=? and time_date >= ? and time_date <=?", devtype, starttime, endtime).Find(&tableplace)
	}
	if len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ? and city_code= ?  and  time_date >= ? and time_date <= ?", provincecode, citycode, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=?  and city_code= ? and time_date >= ? and time_date <=?", provincecode, citycode, starttime, endtime).Find(&tableplace)
	}
	if citycode == "0" && len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("province_code= ?   and time_date >= ? and time_date <= ?", provincecode, starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("province_code=?  and time_date >= ? and time_date <=?", provincecode, starttime, endtime).Find(&tableplace)
	}
	if provincecode == "0" && len(devtype) == 0 && len(equipment) == 0 {
		common.IndexDB.Table("daily_monitorings").Where("time_date >= ? and time_date <= ?", starttime, endtime).Count(&count)
		common.IndexDB.Table("daily_monitorings").Where("time_date >= ? and time_date <=?", starttime, endtime).Find(&tableplace)

	}

	fmt.Println("count=", count)
	fmt.Println("tableplace=", tableplace)

	t.DB.Table("fragment_monitorings").Delete(&tableplace)
	for i := 0; i < len(tableplace); i++ {
		if tableplace[i].Temp_score != "0" {
			t.DB.Table("fragment_monitorings").Create(&tableplace[i])
		}
	}

	var tableplace2 []model.Daily_monitorings
	if downloadflag=="0"{
	if flag == "1" {
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag!=0").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag!=0").Find(&tableplace2)
		//count=yichang
	}
	if flag == "0" {

		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag=0").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag=0").Find(&tableplace2)
	}
	if flag == "2" {
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score from fragment_monitorings where temp_score>=0 group by dev_id ").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score from fragment_monitorings where temp_score>=0 group by dev_id ").Find(&tableplace2)
	}
	count = len(tableplace2)
}else if downloadflag=="1"{
	if flag == "1" {
		//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag!=0").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag!=0").Find(&tableplace)
		//count=yichang
	}
	if flag == "0" {

		//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag=0").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fragment_monitorings where temp_score>=0 group by dev_id having abnormal_flag=0").Find(&tableplace)
	}
	if flag == "2" {
		//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score from fragment_monitorings where temp_score>=0 group by dev_id ").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableplace)
		t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score from fragment_monitorings where temp_score>=0 group by dev_id ").Find(&tableplace2)
	}
	count = len(tableplace)
}
	fmt.Println(tableplace)
	for _, tableDate := range tableplace {
		var tableStoreDates4 model.TableDate3
		t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
		tableStoreDates3 = append(tableStoreDates3, tableStoreDates4)

	}
	fmt.Println(tableplace)

	fmt.Println("count=", count)

	response.Success(ctx, gin.H{"data": tableplace, "count": count, "data1": tableStoreDates3}, "成功")
}
func (t MonitorController) Getequipment(ctx *gin.Context) {
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

		common.IndexDB.Raw("select distinct dev_id from daily_monitorings ").Find(&tableStoreDates4)
		//common.IndexDB.Table("daily_monitorings").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "1" {

		common.IndexDB.Raw("select distinct dev_id from daily_monitorings where province_code=? and dev_type=?", province_code, dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("daily_monitorings").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "2" {
		common.IndexDB.Raw("select distinct dev_id from daily_monitorings where city_code=? and dev_type=?", city_code, dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("daily_monitorings").Where("city_code = ? and dev_type = ? ", city_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else {
		fmt.Println(dev_type)
		common.IndexDB.Raw("select distinct dev_id from daily_monitorings where dev_type=?", dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("daily_monitorings").Where("dev_type = ? ", dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	}
}
func (t MonitorController) Datachart(ctx *gin.Context) {
	equipment := ctx.DefaultQuery("equipment", "")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	var tableplace []model.Daily_monitorings
	if len(equipment) != 0 {
		equipment = equipment
	}
	if len(starttime) != 0 && len(endtime) != 0 {
		starttime = starttime
		endtime = endtime
	}
	common.IndexDB.Table("daily_monitorings").Where(" dev_id=? and time_date >= ? and time_date <=?", equipment, starttime, endtime).Order("time_date ASC").Find(&tableplace)
	fmt.Println(tableplace)
	response.Success(ctx, gin.H{"data": tableplace}, "成功")
}

func (t MonitorController) Ceshi(ctx *gin.Context) {
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

	common.IndexDB.Raw("select distinct dev_type from daily_monitorings").Find(&tableStoreDates1)

	common.IndexDB.Raw("select distinct dev_id from daily_monitorings").Find(&tableStoreDates2) //设备号

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

func (t MonitorController) Gettype(ctx *gin.Context) {




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
		common.IndexDB.Raw("select distinct dev_type from daily_monitorings where province_code = ? ", province_code).Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id from daily_monitorings where province_code = ? ", province_code).Find(&tableStoreDates2) //设备号
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
		common.IndexDB.Raw("select distinct dev_type from daily_monitorings where city_code = ?", city_code).Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id   from daily_monitorings where city_code = ?", city_code).Find(&tableStoreDates2) //设备号
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
		common.IndexDB.Raw("select distinct dev_type from daily_monitorings").Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id from daily_monitorings").Find(&tableStoreDates2) //设备号
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

	//if flag == "1" {
	//common.IndexDB.Raw("select distinct dev_type from daily_monitorings where province_code = ? ", province_code).Find(&tableStoreDates4)
	//fmt.Println(tableStoreDates4)
	//	for _, tableDate := range tableStoreDates4 {
	//		var tableStoreDates3 Type1
	//		tableStoreDates3.Value = tableDate.Dev_type
	//		type1 = append(type1, tableStoreDates3)
	//	}

	//response.Success(ctx, gin.H{"data": type1}, "成功")
	//} else if flag == "2" {
	//	common.IndexDB.Raw("select distinct dev_type from daily_monitorings where city_code = ?", city_code).Find(&tableStoreDates4)
	//	fmt.Println(tableStoreDates4)
	//	for _, tableDate := range tableStoreDates4 {
	//		var tableStoreDates3 Type1
	//		tableStoreDates3.Value = tableDate.Dev_type
	//		type1 = append(type1, tableStoreDates3)
	//	}

	//	response.Success(ctx, gin.H{"data": type1}, "成功")
	//} else {
	//	common.IndexDB.Raw("select distinct dev_type from daily_monitorings").Find(&tableStoreDates4)
	//	for _, tableDate := range tableStoreDates4 {
	//		var tableStoreDates3 Type1
	//		tableStoreDates3.Value = tableDate.Dev_type
	//		type1 = append(type1, tableStoreDates3)
	//}
	//
	//	response.Success(ctx, gin.H{"data": type1}, "成功")
	//	}

}

func NewMonitorController() IMonitorController {
	db := common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return MonitorController{DB: db}
}
