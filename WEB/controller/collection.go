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

type ICollectionController interface {
	RestController
	Getinfo(ctx *gin.Context)
	Create1(ctx *gin.Context)
	Tempequipment(ctx *gin.Context)
	Get_Change_Para_Info(ctx *gin.Context) 
	Get_Row_Para_Info(ctx *gin.Context)
}

type CollectionController struct {
	DB *gorm.DB
}

//存采集信息的表
//type Tableselect struct {
//	Dev_Id                 string `json:"dev_id" gorm:"type:varchar(255);not null"`
//	Start_time             string `json:"start_time" gorm:"type:varchar(255);not null"`
//	End_time               string `json:"end_time" gorm:"type:varchar(255);not null"`
//	Deviation              string `json:"deviation" gorm:"type:float;not null"`
//	Water_score            string `json:"water_score" gorm:"type:varchar(255);not null"`
//	Water_pattern          string `json:"water_pattern" gorm:"type:varchar(255);not null"`
//	Heat_duration          string `json:"heat_duration" gorm:"type:int;not null"`
//	Un_stable_temp_percent string `json:"un_stable_temp_percent" gorm:"type:int;not null"`
//	Overshoot_value        string `json:"overshoot_value" gorm:"type:int;not null"`
//	State_accuracy         string `json:"state_accuracy" gorm:"type:int;not null"`
//	Temp_score             string `json:"temp_score" gorm:"type:int;not null"`
//	Temp_pattern           string `json:"temp_pattern" gorm:"type:int;not null"`
//} //此处用来定义通信的一些内容 用于前后端的传输

func (t CollectionController) Create(ctx *gin.Context) {
	var requestCollection vo.CollectionRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCollection); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建Collection表格
	collection := model.Collection{
		//CategoryId: requestCollection.CategoryId,
		Applianceid:      requestCollection.Applianceid,
		StartTime:        requestCollection.Start_time,
		EndTime:          requestCollection.End_time,
		Move_task_flag:   requestCollection.Move_task_flag,
		Mining_task_flag: requestCollection.Mining_task_flag,
	}
	if err := t.DB.Create(&collection).Error; err != nil {
		panic(err)
		return
	}
	response.Success(ctx, gin.H{"data": collection}, "创建成功")

}

func (t CollectionController) Update(ctx *gin.Context) {
	var requestCollection vo.CollectionRequest
	//数据验证
	if err := ctx.ShouldBind(&requestCollection); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	fmt.Println(ctx.Params)
	//获取path中的id
	collectionId := ctx.Params.ByName("id") //从上下文中解析
	fmt.Println(collectionId)
	var collection model.Collection
	if t.DB.Where("label=?", collectionId).First(&collection).RecordNotFound() {
		Label := ctx.DefaultQuery("Label", "0000")
		Value := ctx.DefaultQuery("Value", "0000")
		fmt.Println("00")
		newUser := model.TableData{
			Label: Label,
			Value: Value,
		}
		requestCollection.Applianceid = ""
		fmt.Println(requestCollection)

		t.DB.Create(&newUser)
		//response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"参数不存在")
		//	response.Fail(ctx,nil,"文章不存在")
		return
	}
	fmt.Println(collection)
	requestCollection.Applianceid = ""
	fmt.Println(requestCollection)
	//更新文章
	if err := t.DB.Model(&collection).Update(requestCollection).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"data": collection}, "更新成功")

}

func (t CollectionController) Show(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var collection model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", tableStoreDateId).First(&collection).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"data": collection}, "读取成功")
}


func (t CollectionController) Delete(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var collection model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&collection).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	t.DB.Delete(&collection)

	response.Fail(ctx, gin.H{"data": collection}, "删除成功")

}


func (t CollectionController) Get_Row_Para_Info(ctx *gin.Context) {
	dev_id := ctx.DefaultQuery("dev_id", "")
	time := ctx.DefaultQuery("time", "")

	fields := []string{}
	values := []interface{}{}

	var Table_Row_info []model.PerchangeNums

	if time != "" {
		fields = append(fields, "updata_time = ?")
		values = append(values, time)
	}
	if dev_id != "" {
		fields = append(fields, "dev_id = ?")
		values = append(values, dev_id)
	}

	//common.IndexDB.Table("perchange_nums").Where(strings.Join(fields, " AND "), values...).Group("dev_id,updata_time").Find(&Table_Row_info)
	common.IndexDB.Table("perchange_nums").Where(strings.Join(fields, " AND "), values...).Find(&Table_Row_info)
	response.Success(ctx, gin.H{"Row_data": Table_Row_info}, "成功")
}


func (t CollectionController) Get_Change_Para_Info(ctx *gin.Context) {
	timeLow := ctx.DefaultQuery("timeLow", "")
	timeHigh := ctx.DefaultQuery("timeHigh", "")
	dev_id := ctx.DefaultQuery("dev_id", "")
	dev_change_state := ctx.DefaultQuery("dev_change_state", "")

	fields := []string{}
	values := []interface{}{}

	var Table_test_info []model.PerchangeNums
	var Table_test_info_all []model.PerchangeNums

	common.IndexDB.Table("perchange_nums").Scan(&Table_test_info_all)
	//common.IndexDB.Table("perchange_nums").Scan(&Table_test_info_all)
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
	//common.IndexDB.Table("perchange_nums").Where(strings.Join(fields, " AND "), values...).Group("dev_id,updata_time").Find(&Table_test_info)
	common.IndexDB.Table("perchange_nums").Where(strings.Join(fields, " AND "), values...).Group("dev_id,updata_time").Find(&Table_test_info)
	// common.IndexDB.Raw("select * from perchange_nums  group by dev_id,updata_time).Find(&avg_score)
	// tx.Find(&TopicList).RecordNotFound()

	response.Success(ctx, gin.H{"Para_data": Table_test_info, "Para_data_all": Table_test_info_all}, "成功")
}



func (t CollectionController) Getinfo(ctx *gin.Context) {
	type ModeFragment1 struct {
		DevId        string  `json:"dev_id" gorm:"type:varchar(50);not null;index:dev_id_idx"`
		StartTime    string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx"`
		EndTime      string  `json:"end_time" gorm:"type:varchar(20);not null"`
		DurationTime string  `json:"duration_time" gorm:"type:varchar(20);not null"`
		WaterPattern int     `json:"water_pattern" gorm:"type:int;not null"`
		FlowAvg      float64 `json:"flow_avg" gorm:"type:float;not null"`
		SmallWater   float64 `json:"small_water" gorm:"type:float;not null"`
		//Extreme              int     `json:"extreme" gorm:"type:int;not null"`
		//MaxChange            float64 `json:"max_change" gorm:"type:float;not null"`
		//Average              float64 `json:"average" gorm:"type:float;not null"`
		Deviation            float64 `json:"deviation" gorm:"type:float;not null"`
		UpNumber             int     `json:"up_number" gorm:"type:int;not null"`
		DownNumber           int     `json:"down_number" gorm:"type:int;not null"`
		WaterScore           int     `json:"water_score" gorm:"int;not null"`
		HeatDuration         string  `json:"heat_duration" gorm:"type:varchar(20);not null"`
		UnStableTempDuration string  `json:"un_stable_temp_duration" gorm:"type:varchar(20);not null"`
		UnStableTempPercent  float64 `json:"un_stable_temp_percent" gorm:"type:float;not null"`
		UnHeatDev            float64 `json:"un_heat_dev" gorm:"type:float;not null"`
		TempPattern          int     `json:"temp_pattern" gorm:"type:int;not null"`
		OvershootValue       int     `json:"overshoot_value" gorm:"type:int;"`
		StateAccuracy        int     `json:"state_accuracy" gorm:"type:int;"`
		TempScore            int     `json:"temp_score" gorm:"type:int;not null"`
		NewTempScore         int     `json:"new_temp_score" gorm:"type:int;not null"`
		HeatTempScore        int     `json:"heat_temp_score" gorm:"type:int;not null"`
		StableTempScore      int     `json:"stable_temp_score" gorm:"type:int;not null"`
		TempJudgeFlag        int    `json:"temp_judge_flag" gorm:"type:int;not null"`
		WaterFlag            int     `json:"water_flag" gorm:"type:int;not null"`
		TempFlag             int     `json:"temp_flag" gorm:"type:int;not null"`
		AbnormalState        int     `json:"abnormal_state" gorm:"type:int;not null"`
		TabName              string  `gorm:"-"`
	}
	var tableStoreDates4 []ModeFragment1
	var count int

	type Users struct {
		Applianceid      string
		Start_time       string
		End_time         string
		Move_task_flag   string
		Mining_task_flag string
	}
	downloadflag:=ctx.DefaultQuery("downloadflag", "0")
	DevId := ctx.DefaultQuery("equipment", "0")
	starttime := ctx.DefaultQuery("timeLow", "0")
	endtime := ctx.DefaultQuery("timeHigh", "0")
	//timeLow := ctx.DefaultQuery("timeLow", "2020-01-01")
	//timeHigh := ctx.DefaultQuery("timeHigh", "2022-01-01")
	//dev_id := ctx.DefaultQuery("dev_id", "0000")

	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))

	fragmentTableName := "e_" + "fragment" + DevId
	var user Users
	t.DB.Table("signle_applicances").Where("applianceid = ? and start_time = ? and end_time =?", DevId, starttime, endtime).Find(&user)
	fmt.Println(DevId)
	fmt.Println(starttime)
	fmt.Println(endtime)
	fmt.Println(user)
	fmt.Println(user.Move_task_flag)
	fmt.Println(user.Mining_task_flag)


	if user.Mining_task_flag == "1" { //挖掘成功
		if downloadflag=="0"{
		common.IndexDB.Table(fragmentTableName).Where("dev_id = ? and start_time >= ? and end_time <=?", DevId, starttime, endtime).Count(&count)
		// common.IndexDB.Table(fragmentTableName).Where("dev_id = ? and start_time >= ? and end_time <=?", DevId, starttime, endtime).Find(&tableStoreDates4)
		
		common.IndexDB.Table(fragmentTableName).Where("dev_id = ? and start_time >= ? and end_time <=?", DevId, starttime, endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableStoreDates4)
		}
		if downloadflag=="1"{
			common.IndexDB.Table(fragmentTableName).Where("dev_id = ? and start_time >= ? and end_time <=?", DevId, starttime, endtime).Count(&count)
		// common.IndexDB.Table(fragmentTableName).Where("dev_id = ? and start_time >= ? and end_time <=?", DevId, starttime, endtime).Find(&tableStoreDates4)
		
		common.IndexDB.Table(fragmentTableName).Where("dev_id = ? and start_time >= ? and end_time <=?", DevId, starttime, endtime).Find(&tableStoreDates4)
		}
		response.Success(ctx, gin.H{"flag_state": "2", "data": tableStoreDates4, "count": count}, "成功")
	} else if user.Mining_task_flag == "2" { //挖掘失败 可能设备信息错误等
		response.Success(ctx, gin.H{"flag_state": "4"}, "失败")
	} else if user.Mining_task_flag == "3" { //挖掘失败 时间段内无设备信息
		response.Success(ctx, gin.H{"flag_state": "5"}, "失败")
	} else if user.Move_task_flag == "3" { //迁移失败 迁移数据为空
		response.Success(ctx, gin.H{"flag_state": "6"}, "失败")
	} else if user.Move_task_flag == "2" { //迁移失败 迁移设备id出错
		response.Success(ctx, gin.H{"flag_state": "3"}, "失败")
	} else if user.Move_task_flag == "1" { //迁移成功
		response.Success(ctx, gin.H{"flag_state": "1"}, "成功")
	} else { //等待中
		response.Success(ctx, gin.H{"flag_state": "0"}, "成功")
	}
}

func (t CollectionController) Create1(ctx *gin.Context) {
	type Users struct {
		Applianceid      string
		Start_time       string
		End_time         string
		Adjust_task_flag string
		Move_task_flag   string
		Mining_task_flag string
	}

	DevId := ctx.DefaultQuery("equipment", "0")
	starttime := ctx.DefaultQuery("timeLow", "0")
	endtime := ctx.DefaultQuery("timeHigh", "0")
    adjust_task_flag:=ctx.DefaultQuery("adjust_task_flag", "0")
	fmt.Println(endtime)

	type_ceshi := Users{Applianceid: DevId, Start_time: starttime, End_time: endtime,Adjust_task_flag: adjust_task_flag, Move_task_flag: "0", Mining_task_flag: "0"}

	t.DB.Table("signle_applicances").Create(&type_ceshi)
	fmt.Println(type_ceshi)
	response.Success(ctx, gin.H{"data": type_ceshi}, "成功")

}

func (t CollectionController) Tempequipment(ctx *gin.Context) {
	timeLow := ctx.DefaultQuery("timeLow", "2020-01-01")
	timeHigh := ctx.DefaultQuery("timeHigh", "2022-01-01")
	dev_id := ctx.DefaultQuery("dev_id", "0000")

	var Dev_Id []string
	var Start_time []string
	var End_time []string
	var Deviation []string
	var Water_score []string
	var Water_pattern []string
	var Heat_duration []string
	var Un_stable_temp_percent []string
	var Overshoot_value []string
	var State_accuracy []string
	var Temp_score []string
	var NewTempscore []string
	var Temp_pattern []string

	var TableFragceshi []model.Tablefragceshi

	t.DB.Table("fragment110100").Where("start_time >=  ? AND end_time  <=  ? AND dev_id= ? ", timeLow, timeHigh, dev_id).Find(&TableFragceshi)
	//fmt.Print(TableFragceshi)
	for _, Tablefragceshi := range TableFragceshi {
		Dev_Id = append(Dev_Id, Tablefragceshi.Dev_Id) //数据放入切片中
		Start_time = append(Start_time, Tablefragceshi.Start_time)
		End_time = append(End_time, Tablefragceshi.End_time) //数据放入切片中
		Deviation = append(Deviation, Tablefragceshi.Deviation)
		Water_score = append(Water_score, Tablefragceshi.Water_score) //数据放入切片中
		Water_pattern = append(Water_pattern, Tablefragceshi.Water_pattern)
		Heat_duration = append(Heat_duration, Tablefragceshi.Heat_duration) //数据放入切片中
		Un_stable_temp_percent = append(Un_stable_temp_percent, Tablefragceshi.Un_stable_temp_percent)
		Overshoot_value = append(Overshoot_value, Tablefragceshi.Overshoot_value) //数据放入切片中
		State_accuracy = append(State_accuracy, Tablefragceshi.State_accuracy)
		Temp_score = append(Temp_score, Tablefragceshi.Temp_score)
		NewTempscore = append(NewTempscore, Tablefragceshi.NewTempScore)
		Temp_pattern = append(Temp_pattern, Tablefragceshi.Temp_pattern)
	}
	response.Success(ctx, gin.H{"data": TableFragceshi, "Dev_Id": Dev_Id, "Start_time": Start_time, "End_time": End_time, "Deviation": Deviation, "Water_score": Water_score, "Water_pattern": Water_pattern, "Heat_duration": Heat_duration, "Un_stable_temp_percent": Un_stable_temp_percent, "Overshoot_value": Overshoot_value, "Temp_score": Temp_score, "NewTempscore": NewTempscore, "Temp_pattern": Temp_pattern}, "成功")

}

//func NewCollectionController() ICollectionController {
//	db := common.GetDB()
//	// db.AutoMigrate(model.TableDate{})
//	return CollectionController{DB: db}
//}
func NewCollectionController() ICollectionController {
	db := common.GetDB()
	//db.AutoMigrate(model.Collection{})
	return CollectionController{DB: db}
}
