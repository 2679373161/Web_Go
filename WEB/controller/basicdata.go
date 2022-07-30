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

type IBasicdataController interface {
	RestController
	PageList(ctx *gin.Context)
	Modifyequipment(ctx *gin.Context)
	Getequipment(ctx *gin.Context)
	Gettype(ctx *gin.Context)
}

type BasicdataController struct {
	DB *gorm.DB
}

type Tableselect struct {
	Dev_Id          string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Model_type      string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code       string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_province    string `json:"dev_province" gorm:"type:varchar(255);not null"`
	Dev_city        string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Opt             string `json:"handle_flag" gorm:"type:int;not null"`
	Monitoring_flag string `json:"monitoring_flag" gorm:"type:int;not null"`
	City_migrate_code string `json:"city_migrate_code" gorm:"type:int;not null"`
} //此处用来定义通信的一些内容 用于前后端的传输
type Loc_code struct {
	Province_code string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code     string `json:"city_code" gorm:"type:varchar(255);not null"`
	Region_code   string `json:"region_code" gorm:"type:varchar(255);not null"`
	Dev_province  string `json:"dev_province" gorm:"type:varchar(255);not null"`
	Dev_city      string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Dev_region    string `json:"dev_region" gorm:"type:varchar(255);not null"`
}
type Modeltype_and_nums struct {
	Model_type string `json:"model_type" gorm:"type:varchar(255);not null"`
	Num        string `json:"num" gorm:"type:varchar(255);not null"`
	Opt        string `json:"opt" gorm:"type:varchar(255);not null"`
}

func (t BasicdataController) Create(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestTableStoreDate); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建tablestoredate
	tableStoreDate := model.TableDate{
		//CategoryId: requestTableStoreDate.CategoryId,
		Datatime: requestTableStoreDate.Datatime,
		Flame:    requestTableStoreDate.Flame,
		Outtemp:  requestTableStoreDate.Outtemp,
		Settemp:  requestTableStoreDate.Settemp,
		Flow:     requestTableStoreDate.Flow,
	}
	if err := t.DB.Create(&tableStoreDate).Error; err != nil {
		panic(err)
		// return
	}
	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "创建成功")

}

func (t BasicdataController) Update(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestTableStoreDate); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	//获取path中的id
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	//更新文章
	if err := t.DB.Model(&tableStoreDate).Update(requestTableStoreDate).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "更新成功")

}

func (t BasicdataController) Show(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "读取成功")
}

func (t BasicdataController) Delete(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx, gin.H{"tableStoreDate": tableStoreDate}, "删除成功")

}

func (t BasicdataController) PageList(ctx *gin.Context) {
	flag := ctx.DefaultQuery("flag", "2")
	orderflag := ctx.DefaultQuery("orderflag", "")
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	province_code := ctx.DefaultQuery("province_code", "0")
	city_code := ctx.DefaultQuery("city_code", "0")
	dev_type := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	temp_flag := ctx.DefaultQuery("monitoring_flag", "0")
	var tableStoreDates4 []Tableselect
	var district Loc_code
	var count int
	select_option := make(map[string]interface{})

	if province_code != "0" {
		select_option["province_code"] = province_code
	}
	if city_code != "0" {
		select_option["city_code"] = city_code
	}
	if len(dev_type) != 0 {
		select_option["model_type"] = dev_type
	}
	if len(equipment) != 0 {
		select_option["dev_id"] = equipment
	}
	if flag != "2" {
		if temp_flag == "0" {
			select_option["opt"] = flag
		} else if temp_flag == "1" {
			select_option["monitoring_flag"] = flag
		}

	}
	if temp_flag == "1" {
		select_option["opt"] = "1"
	}
	fmt.Println(select_option)
	if orderflag == "1"{
		t.DB.Table(table_overall_equipment).Where(select_option).Count(&count)
		t.DB.Table(table_overall_equipment).Where(select_option).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("opt DESC").Find(&tableStoreDates4)
	} else if orderflag == "2"{
		t.DB.Table(table_overall_equipment).Where(select_option).Count(&count)
		t.DB.Table(table_overall_equipment).Where(select_option).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("monitoring_flag DESC").Find(&tableStoreDates4)
	}
	for i, tableDate := range tableStoreDates4 {
		t.DB.Table(table_loc_code).Where("city_code = ? ", tableDate.City_code).Find(&district)
		tableStoreDates4[i].Dev_province = district.Dev_province
		tableStoreDates4[i].Dev_city = district.Dev_city
	}
	fmt.Println(tableStoreDates4)

	response.Success(ctx, gin.H{"data": tableStoreDates4, "total": count}, "成功")
}

func (t BasicdataController) Modifyequipment(ctx *gin.Context) {
	province_code := ctx.DefaultQuery("province_code", "0")
	city_code := ctx.DefaultQuery("city_code", "0")
	dev_type := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	id := ctx.DefaultQuery("id", "0")
	flag := ctx.DefaultQuery("flag", "0")
	flag1 := ctx.DefaultQuery("flag1", "0")
	pageselect := ctx.QueryArray("pageselect[]")
	temp_flag := ctx.DefaultQuery("monitoring_flag", "0")
	migrate_code:=ctx.DefaultQuery("city_migrate_code", "0")
	pageselect_city := ctx.QueryArray("pageselect_city[]")

	var tableselect []Tableselect
	select_option := make(map[string]interface{})
	var migrate_table model.Citycode_migrate
	var num=0
	if province_code != "0" {
		select_option["province_code"] = province_code
	}
	if city_code != "0" {
		select_option["city_code"] = city_code
	}
	if len(dev_type) != 0 {
		select_option["model_type"] = dev_type
	}
	if len(equipment) != 0 {
		select_option["dev_id"] = equipment
	}
	fmt.Println("666")
	fmt.Println("pageselect_city=",pageselect_city)
	fmt.Println("pageselect=",pageselect)
fmt.Println("code=",migrate_code)
	if temp_flag == "0" {
		if flag1 == "selectall" {
			t.DB.Table(table_overall_equipment).Model(&tableselect).Where(select_option).Update(map[string]interface{}{"opt": flag, "handle_flag":flag})
			t.DB.Table("city_select_codes").Model(&migrate_table).Where(select_option).Update("opt", flag)
		} else if flag1 == "pageall" {
		t.DB.Table(table_overall_equipment).Model(&tableselect).Where("dev_id in (?)", pageselect).Update(map[string]interface{}{"opt": flag, "handle_flag":flag})
		t.DB.Table("city_select_codes").Model(&migrate_table).Where("code in (?)", pageselect_city).Update("opt", flag)


		} else {
			t.DB.Table(table_overall_equipment).Model(&tableselect).Where("dev_id=?", id).Update(map[string]interface{}{"opt": flag, "handle_flag":flag})
			t.DB.Table(table_overall_equipment).Model(&tableselect).Where("city_migrate_code= ? and opt= ?", migrate_code,"1").Count(&num)
			fmt.Println("num=",num)
			if num==0{
				t.DB.Table("city_select_codes").Model(&migrate_table).Where("code=?", migrate_code).Update("opt", "0")
			}else{
			t.DB.Table("city_select_codes").Model(&migrate_table).Where("code=?", migrate_code).Update("opt", "1")
		}}
	//	t.DB.Raw("select city_migrate_code,sum(opt) as opt from midea_device_all_select_onlines where opt>0 group by city_migrate_code").Find(&migrate_table)
	//	fmt.Println("666")
	//	fmt.Println(migrate_table)


	} else if temp_flag == "1" {
		if flag1 == "selectall" {
			t.DB.Table(table_overall_equipment).Model(&tableselect).Where("opt=?","1").Update("monitoring_flag", flag)
			t.DB.Table("bo").Model(&tableselect).Where(select_option).Update("monitoring_flag", flag)
		} else if flag1 == "pageall" {
			t.DB.Table(table_overall_equipment).Model(&tableselect).Where("dev_id in (?)", pageselect).Update("monitoring_flag", flag)
			t.DB.Table("bo").Model(&tableselect).Where("dev_id in (?)", pageselect).Update("monitoring_flag", flag)
		} else {
			t.DB.Table(table_overall_equipment).Model(&tableselect).Where("dev_id=?", id).Update("monitoring_flag", flag)
			t.DB.Table("bo").Model(&tableselect).Where("dev_id=?", id).Update("monitoring_flag", flag)
		}
	}

}

func (t BasicdataController) Getequipment(ctx *gin.Context) {

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


	if flag == "1" {
		t.DB.Table(table_overall_equipment).Where("province_code = ? and model_type = ?", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "2" {
		t.DB.Table(table_overall_equipment).Where("city_code = ? and model_type = ?", city_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else {
	
		t.DB.Table(table_overall_equipment).Where("model_type = ?", dev_type).Find(&tableStoreDates4)
		fmt.Println("tableStoreDates4=",tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	}
}

func (t BasicdataController) Gettype(ctx *gin.Context) {
	flag := ctx.DefaultQuery("flag", "0")
	province_code := ctx.DefaultQuery("province_code", "0")
	city_code := ctx.DefaultQuery("city_code", "0")
	type Type struct {
		Model_type string `json:"model_type" gorm:"type:varchar(255);not null"`
	}
	type Type1 struct {
		Value string `json:"value"`
	}
	var tableStoreDates4 []Type

	var type1 []Type1
	if flag == "1" {
		t.DB.Raw("select distinct model_type from midea_device_all_select_onlines where province_code = ?", province_code).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Type1
			tableStoreDates3.Value = tableDate.Model_type
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "2" {
		t.DB.Raw("select distinct model_type from midea_device_all_select_onlines where city_code = ?", city_code).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Type1
			tableStoreDates3.Value = tableDate.Model_type
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else {
		t.DB.Table(table_alltype).Find(&tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Type1
			tableStoreDates3.Value = tableDate.Model_type
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	}

}

func NewBasicdataController() IBasicdataController {
	db := common.GetDB()
	// db.AutoMigrate(model.TableDate{})
	return BasicdataController{DB: db}
}
