package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

type IScoreSummaryController interface{
	RestController

	Search(ctx *gin.Context)
	Getequipment(ctx *gin.Context)
	Gettype(ctx *gin.Context)
/*	Ceshi(ctx *gin.Context)*/
    Scorechart(ctx *gin.Context)
    Intempchart(ctx *gin.Context)
	Residualheatchart(ctx *gin.Context)
	Funnelchart(ctx *gin.Context)

}
type ScoreSummaryController struct {
	DB *gorm.DB
}
func (t ScoreSummaryController) Create(ctx *gin.Context) {
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

func (t ScoreSummaryController) Update(ctx *gin.Context) {
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

func (t ScoreSummaryController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t ScoreSummaryController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t ScoreSummaryController) Search(ctx *gin.Context) {
	CityCode :=ctx.DefaultQuery("city_code", "")
	DevType := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	starttime := ctx.DefaultQuery("start_time", "")
	endtime  :=ctx.DefaultQuery("end_time", "")
	var totaltable []model.Total
	if equipment !=""{
      fragmenttable:="fragment"+CityCode

      common.IndexDB.Table(fragmenttable).Select("dev_id,count(start_time) AS total_count").
		  Where(
			  "dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				  "and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
		  equipment,starttime, endtime).Group("dev_id").Find(&totaltable)
      if len(totaltable)!=0 {
		  common.IndexDB.Table(fragmenttable).Where(
			  "dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				  "and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			  equipment, starttime, endtime).Count(&totaltable[0].TotalNum)
	  }
		for i:=0;i<len(totaltable);i++ {
          var DevCity model.Table_all_select
          common.DB.Table("midea_loc_code").Where("city_code=?",CityCode).Find(&DevCity)
          totaltable[i].DevCity=DevCity.Dev_city
          totaltable[i].DevType=DevType
          totaltable[i].StartTime=starttime
          totaltable[i].EndTime=endtime
		}

	}else if DevType!="" {
		type Equipment struct {
			Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
		}
		var type1 []Equipment
		fragmenttable:="fragment"+CityCode
		common.DB.Table("bo").Select("dev_id").
			Where("dev_type=? AND city_code=? AND handle_flag=1",DevType,CityCode).
			Find(&type1)
		for _,tabledate:=range type1{
			var totaltable1 model.Total
			common.IndexDB.Table(fragmenttable).Select("dev_id,count(*) AS total_num").
				Where(
					"dev_id=? AND (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
						"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
					tabledate.Dev_id,starttime, endtime).Group("dev_id").Find(&totaltable1)
			if totaltable1.DevId!=""{
				common.IndexDB.Table(fragmenttable).Where(
					"dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
						"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
					tabledate.Dev_id, starttime, endtime).Count(&totaltable1.TotalNum)
			}
                if totaltable1.TotalNum>=1 {
					fmt.Println(totaltable1)
					totaltable = append(totaltable, totaltable1)
				}
		}
		for i:=0;i<len(totaltable);i++ {
			var DevCity model.Table_all_select
			common.DB.Table("midea_loc_code").Where("city_code=?",CityCode).Find(&DevCity)
			totaltable[i].DevCity=DevCity.Dev_city
			totaltable[i].DevType=DevType
			totaltable[i].StartTime=starttime
			totaltable[i].EndTime=endtime
		}
		fmt.Println(totaltable)
	}else {
		fragmenttable:="fragment"+CityCode
		common.IndexDB.Table(fragmenttable).Select("dev_id,count(*) AS total_num").
			Where(
				"(temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
					"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
				starttime, endtime).Group("dev_id").Find(&totaltable)
		for i:=0;i<len(totaltable);i++ {
			var DevCity model.Table_all_select
			common.DB.Table("midea_loc_code").Where("city_code=?",CityCode).Find(&DevCity)
			common.DB.Table("bo").Select("dev_type AS model_type").Where("dev_id=?",totaltable[i].DevId).Find(&DevCity)
			totaltable[i].DevCity=DevCity.Dev_city
			totaltable[i].DevType=DevCity.Model_type
			totaltable[i].StartTime=starttime
			totaltable[i].EndTime=endtime
		}
		fmt.Println(totaltable)
	}
	response.Success(ctx, gin.H{"data":totaltable}, "成功")

}


func (t ScoreSummaryController) Getequipment(ctx *gin.Context) {
	/*flag := ctx.DefaultQuery("flag", "0")*/
	table:=ctx.DefaultQuery("table", "0")
	/*province_code := ctx.DefaultQuery("province_code", "0")*/
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
	if dev_type!=""{
		common.DB.Table(table).Select("distinct dev_id").Where("dev_type=? AND handle_flag=1",dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if city_code!=""{
		common.DB.Table(table).Select("distinct dev_id").Where("city_code=? AND handle_flag=1",city_code).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	}

}
func (t ScoreSummaryController)Scorechart(ctx *gin.Context)  {
	equipment := ctx.DefaultQuery("equipment", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	CityCode:=ctx.DefaultQuery("city_code", "")
	fragmenttable:="fragment"+CityCode
	type Scorechart struct {
		Start_time              string  `json:"start_time" gorm:"type:varchar(255);not null"`
		End_time                string  `json:"end_time" gorm:"type:varchar(255);not null"`
		Temp_score              int     `json:"temp_score" gorm:"type:varchar(255);not null"`
		Intemp                  int `json:"intemp" gorm:"type:varchar(255);not null"`
		Outtemp                 int `json:"intemp" gorm:"type:varchar(255);not null"`

	}
	var fragment []Scorechart
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Order("start_time ASC").Find(&fragment)
	var score []int
	var intemp []int
	var outtemp []int
	var datetime[]string
	for i:=0;i<len(fragment);i++{
		score=append(score,fragment[i].Temp_score)
		intemp=append(intemp,fragment[i].Intemp)
		outtemp=append(outtemp,fragment[i].Outtemp-fragment[i].Intemp)
		datetime=append(datetime,fragment[i].Start_time[5:]+"至"+fragment[i].End_time[5:])
	}
	response.Success(ctx, gin.H{"score": score, "datetime": datetime,"intemp":intemp,"outtemp":outtemp}, "成功")

}
func (t ScoreSummaryController)Intempchart(ctx *gin.Context)  {
	equipment := ctx.DefaultQuery("equipment", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	CityCode:=ctx.DefaultQuery("city_code", "")
	fragmenttable:="fragment"+CityCode
	type Intemp struct {
		Intemp int `json:"intemp" gorm:"type:varchar(255);not null"`
	}
	var fragment []Intemp
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Order("intemp desc").Find(&fragment)
	var intemp []int
	for i:=0;i<len(fragment);i++{
		intemp=append(intemp,fragment[i].Intemp)
	}
	fmt.Println("intemp",intemp)
	response.Success(ctx, gin.H{"intemp": intemp}, "成功")

}
func (t ScoreSummaryController)Residualheatchart(ctx *gin.Context)  {
	equipment := ctx.DefaultQuery("equipment", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	CityCode:=ctx.DefaultQuery("city_code", "")
	fragmenttable:="fragment"+CityCode
	type Outtemp struct {
		Outtemp int `json:"intemp" gorm:"type:varchar(255);not null"`
		Intemp int `json:"intemp" gorm:"type:varchar(255);not null"`
	}
	var fragment []Outtemp
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Order("outtemp desc").Find(&fragment)
	var residualheat []int
	for i:=0;i<len(fragment);i++{
		residualheat=append(residualheat,fragment[i].Outtemp-fragment[i].Intemp)
	}
	fmt.Println("residualheat",residualheat)
	response.Success(ctx, gin.H{"residualheat": residualheat}, "成功")

}

func (t ScoreSummaryController)Funnelchart(ctx *gin.Context)  {
	equipment := ctx.DefaultQuery("equipment", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	CityCode:=ctx.DefaultQuery("city_code", "")
	fragmenttable:="fragment"+CityCode
	type Funnel struct {
		Pattern    int `json:"pattern" gorm:"type:varchar(255);not null"`
		Pattern11  int `json:"pattern11" gorm:"type:varchar(255);not null"`
		Pattern21  int `json:"pattern21" gorm:"type:varchar(255);not null"`
		Pattern12  int `json:"pattern12" gorm:"type:varchar(255);not null"`
		Pattern22  int `json:"pattern22" gorm:"type:varchar(255);not null"`
	}
	var fragment Funnel
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND  (temp_pattern = 11 or temp_pattern = 12 or temp_pattern = 21  or temp_pattern = 22) "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Count(&fragment.Pattern)
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND temp_pattern = 11  "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Count(&fragment.Pattern11)
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND temp_pattern = 21  "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Count(&fragment.Pattern21)
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND temp_pattern = 12  "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Count(&fragment.Pattern12)
	common.IndexDB.Table(fragmenttable).
		Where(
			"dev_id=? AND temp_pattern = 22  "+
				"and start_time >= ? and end_time <= ? and heat_flameout_duration_f <= 6 and fault_code LIKE '_0___' ",
			equipment,start_time, end_time ).Count(&fragment.Pattern22)
	fmt.Println("funnel",fragment)
	response.Success(ctx, gin.H{"funnel": fragment}, "成功")

}



func (t ScoreSummaryController) Gettype(ctx *gin.Context) {
	table_name:=ctx.DefaultQuery("table","0")
	city_code := ctx.DefaultQuery("city_code", "0")
	type Types1 struct {
		Dev_type string `json:"dev_type" gorm:"type:varchar(255);not null"`
	}
	type Type1 struct {
		Value string `json:"value"`
	}
	var type1 []Type1
	var tableStoreDates1 []Types1
      common.DB.Table(table_name).Select("distinct dev_type").Where("city_code=? AND handle_flag=1",city_code).Find(&tableStoreDates1)
	for _, tableDate1 := range tableStoreDates1 {
		var tableStoreDates2 Type1
		tableStoreDates2.Value = tableDate1.Dev_type
		type1 = append(type1, tableStoreDates2)

	}
	fmt.Println(type1)
	response.Success(ctx, gin.H{"data": type1}, "成功")
}




	func NewscoresummaryController () IScoreSummaryController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return ScoreSummaryController{DB:db}
}
