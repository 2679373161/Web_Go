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
	"strconv"
)

type ImenuController interface{
	RestController

	Trend(ctx *gin.Context)
}
type menuController struct {
	DB *gorm.DB
}
func (t menuController) Create(ctx *gin.Context) {
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
		return
	}
	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"创建成功")

}

func (t menuController) Update(ctx *gin.Context) {
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

func (t menuController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t menuController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t menuController) Trend(ctx *gin.Context) {



	//迁移时间
	var transfer[]model.Migration_information_record
	var transfer_time[]string
	var transfer_xdata[]string
	//挖掘时间
	var excavate[]model.Datafeature
	var excavate_time[]string
	var excavate_xdata[]string
	//地图使用
	var geotableplace []model.Tableplace
	var geoweekdate []model.Weekdate
	var weekdatecount string
	var tableStoreDates4 model.TableDate3
	var province []string
	var equipment_num []int
	var tempscore []int
	var waterscore []int
	var province1 []string
	var equipment_num1 []int
	var tempscore1 []int
	var waterscore1 []int
	//剩余空间及温度评分分布
	var feature model.Datafeature
    var feature2[]model.Datafeature
	// 温度评分趋势图
	var typequshi1 []model.Type_oneday_avgscores
	var typequshi2 []model.Type_oneday_avgscores
	var typequshibest []model.Type_oneday_avgscores
	var typequshiworst []model.Type_oneday_avgscores
	var regionqushi1   []model.Region_oneday_average_scores
	var regionqushi2   []model.Region_oneday_average_scores
	var regionqushiworst []model.Region_oneday_average_scores
	var regionqushibest  []model.Region_oneday_average_scores
	var cityqushi1   []model.City_oneday_average_scores
	var cityqushi2   []model.City_oneday_average_scores
	var cityqushiworst []model.City_oneday_average_scores
	var cityqushibest  []model.City_oneday_average_scores
    var china_ydata[] string
    // 全国一周型号评分
	var typezhuzhuang7 []model.Type_oneday_avgscores
	var typezhuzhuang1 []model.Type_oneday_avgscores
	//全国区域及省份榜单
	var region[]string
    var region1[]model.Region_oneday_average_scores
	var tableStoreDates3 model.TableDate3
    //用水片段统计
    var water_frag [2]model.Dayinformation

	var dayinformation model.Dayinformation
	var yestinformation model.Dayinformation



//var abnormalinformation []  Day_information
 //var dayinformation [] Day_information
 var migratenum int
 var monitornum int
 var daydata []model.Day_data
 //var devscore []model.Weekdate
// var yestinformation [] Day_information

 //var tempscore []float32
 //var day []string
// common.IndexDB.Raw("select dev_id,round(avg(temp_score),0) as temp_score from week_dates where temp_score>=0 group by dev_id order by temp_score limit 5 ").Find(&tableplace)
 common.IndexDB.Table("data_features").Order("update_time desc").Last(&dayinformation)
abnormal_count:=dayinformation.Abnormal_count
	  common.IndexDB.Raw("select   update_time,day_all_dev_avg_score  from  data_features  order   by   data_features.update_time desc limit 0,7").Find(&daydata)
common.DB.Table("midea_device_all_select_onlines").Where("opt= ? ", 1).Count(&migratenum)
common.DB.Table("midea_device_all_select_onlines").Where("opt= ? and monitoring_flag= ? ", 1,1).Count(&monitornum)

common.IndexDB.Raw("select abnormal_count  from  data_features  order by data_features.update_time desc limit 1,1").Find(&yestinformation)
new_create:=dayinformation.Abnormal_count-yestinformation.Abnormal_count

//统计用水片段情况（一周和一天）
//一周
common.IndexDB.Raw("select sum(temp_all_normal) as temp_all_normal,sum(constant_temp_abnormal) as constant_temp_abnormal,sum(elevate_temp_abnormal) as elevate_temp_abnormal  ,sum(temp_all_abnormal) as temp_all_abnormal from week_dates ").Find(&water_frag[0])
//一天
common.IndexDB.Table("week_dates").Select(" sum(temp_all_normal) as temp_all_normal,sum(constant_temp_abnormal) as constant_temp_abnormal,sum(elevate_temp_abnormal) as elevate_temp_abnormal  ,sum(temp_all_abnormal) as temp_all_abnormal").Group("time_date").Order("time_date desc").Last(&water_frag[1])



   //计算型号均分
	//计算型号均分
	var typeavgscore[2] model.Type_oneday_avgscores
//右下角榜单
var National_ranking struct{
	First_region      string `json:"first_region" gorm:"type:varchar(255);not null"`
    Second_region      string `json:"second_region" gorm:"type:varchar(255);not null"`
    Third_region     string `json:"third_region" gorm:"type:varchar(255);not null"`
    First_province     string `json:"first_province" gorm:"type:varchar(255);not null"`
    Second_province     string `json:"second_province" gorm:"type:varchar(255);not null"`
    Third_province     string `json:"third_province" gorm:"type:varchar(255);not null"`
	Fourth_province     string `json:"fourth_province" gorm:"type:varchar(255);not null"`
    Fifth_province     string `json:"fifth_province" gorm:"type:varchar(255);not null"`
    Sixth_province     string `json:"sixth_province" gorm:"type:varchar(255);not null"`

}
common.IndexDB.Table("data_features").Order("update_time desc").Last(&feature)
common.IndexDB.Table("national_ranking_tables").Where("date=?",feature.Update_time).Find(&National_ranking)
//fmt.Println("nationnal_ranking",National_ranking)
var score_equipment []model.Score_equipment
common.IndexDB.Table("low_score_equ_rankings").Where("date=?",feature.Update_time).Find(&score_equipment)
//从score_equipment里提取出省份code
//fmt.Println("low_score_equ_rankings",score_equipment[0].Province )
//注意结构体切片的序号是在.前面
var province_code []string = []string{score_equipment[0].Province, score_equipment[1].Province, score_equipment[2].Province,score_equipment[3].Province,score_equipment[4].Province}
//fmt.Println(province_code)
//榜单上的评分
var score_province []model.Score_province
common.IndexDB.Table("province_days").Where("province_code in (?) AND time_date=?",province_code,feature.Update_time).Order("temp_score desc").Find(&score_province)
fmt.Println(score_province)

//监控片段表
var monitor_fragement []model.Monitor_fragement
var monitor_fragement_city []model.TableDate3
common.IndexDB.Raw("select *  from  temp_abnormal_modes order by start_time limit 10 ").Find(&monitor_fragement)
for _, tableDate := range monitor_fragement {
	var tableStoreDates4 model.TableDate3
	t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
	monitor_fragement_city = append(monitor_fragement_city, tableStoreDates4)
	//fmt.Println(monitor_fragement1)
	//fmt.Println("城市=",tableStoreDates7)
}
var fault_fragement []model.MideaFault
var fault_fragement_city []model.TableDate3
common.IndexDB.Raw("select * from midea_fault2 where e0!=0|| e1!=0|| e2!=0|| e3!=0||e4!=0||e5!=0|| e6!=0|| e8!=0|| ea!=0|| ee!=0||f2!=0|| c0!=0|| c1!=0|| c2!=0|| c3!=0|| c4!=0|| c5!=0||c6!=0|| c7!=0|| c8!=0||eh!=0||ef!=0 order by time_date desc").Limit(10).Find(&fault_fragement)
fmt.Println(fault_fragement)
for _, tableDate := range fault_fragement {
	var citycode model.Statistics2
	t.DB.Table("bo").Where("dev_id=?",tableDate.DevID).Find(&citycode)
	var tableStoreDates5 model.TableDate3
	t.DB.Table("midea_loc_code").Where("city_code = ?", citycode.City_code).Find(&tableStoreDates5)
	fault_fragement_city = append(fault_fragement_city, tableStoreDates5)
	//fmt.Println("设备",citycode)
	fmt.Println("城市",citycode)
	//fmt.Println("城市=",tableStoreDates7)
}


//地图信息获取
	common.IndexDB.Table("week_dates").Find(&geoweekdate)
	common.IndexDB.Table("week_dates").Count(&weekdatecount)
	time1:=geoweekdate[0].Time_date

    count,_:= strconv.Atoi(weekdatecount)
	time2:=geoweekdate[count-1].Time_date
    fmt.Println(time1)
    fmt.Println(time2)
	common.IndexDB.Table("province_days").Select("province_code,round(avg(temp_score),0) as temp_score,round(avg(water_score),0) as water_score,round(avg(equipment_num),0) as equipment_num").Group("province_code").Where("time_date>=? AND time_date<=?",time1,time2).Find(&geotableplace)
	for _, tableDate := range geotableplace {
		t.DB.Table("midea_loc_code").Where("province_code=?", tableDate.Province_code).Find(&tableStoreDates4)
		province = append(province, tableStoreDates4.Dev_province)
		equipment_num = append(equipment_num, tableDate.Equipment_num)
		tempscore = append(tempscore, tableDate.Temp_score)
		waterscore = append(waterscore, tableDate.Water_score)
	}
	common.IndexDB.Table("province_days").Select("province_code,round(avg(temp_score),0) as temp_score,round(avg(water_score),0) as water_score,round(avg(equipment_num),0) as equipment_num").Group("province_code").Where("time_date=?",time2).Find(&geotableplace)
	for _, tableDate := range geotableplace {
		t.DB.Table("midea_loc_code").Where("province_code=?", tableDate.Province_code).Find(&tableStoreDates4)
		province1 = append(province1, tableStoreDates4.Dev_province)
		equipment_num1 = append(equipment_num1, tableDate.Equipment_num)
		tempscore1 = append(tempscore1, tableDate.Temp_score)
		waterscore1= append(waterscore1, tableDate.Water_score)
	}
	//剩余空间及评分分布获取
	
	common.IndexDB.Table("data_features").Order("update_time desc").Limit(30).Find(&feature2)
	for _, tableDate := range feature2 {
		china_ydata = append(china_ydata, tableDate.Day_all_dev_avg_score)
	}

	// 温度评分趋势图
	common.IndexDB.Table("type_oneday_avgscores").Select("dev_type,round(avg(avgscore),0) as avgscore").Group("dev_type").Where("time_date>=? AND time_date<=? and avgscore>? ",time1,time2,0).Order("avgscore desc").Last(&typequshi1)
	common.IndexDB.Table("type_oneday_avgscores").Select("dev_type,round(avg(avgscore),0) as avgscore").Group("dev_type").Where("time_date>=? AND time_date<=? and avgscore>?",time1,time2,0).Order("avgscore asc").First(&typequshi2)
    common.IndexDB.Table("type_oneday_avgscores").Where("dev_type=?",typequshi1[0].Dev_type).Where("time_date>=? AND time_date<=? ",time1,time2).Find(&typequshibest)
	common.IndexDB.Table("type_oneday_avgscores").Where("dev_type=?",typequshi2[0].Dev_type).Where("time_date>=? AND time_date<=?",time1,time2).Find(&typequshiworst)
	common.IndexDB.Table("region_oneday_average_scores").Select("region_code,round(avg(avg_score),0) as avg_score").Group("region_code").Where(" avg_score>? ",0).Order("avg_score desc").Last(&regionqushi1)
	common.IndexDB.Table("region_oneday_average_scores").Select("region_code,round(avg(avg_score),0) as avg_score").Group("region_code").Where(" avg_score>? ",0).Order("avg_score asc").First(&regionqushi2)
	common.IndexDB.Table("region_oneday_average_scores").Where("region_code=?",regionqushi1[0].Region_code).Where("time_date>=? AND time_date<=?",time1,time2).Find(&regionqushibest)
	common.IndexDB.Table("region_oneday_average_scores").Where("region_code=?",regionqushi2[0].Region_code).Where("time_date>=? AND time_date<=?",time1,time2).Find(&regionqushiworst)
		var tableStore1 model.Region_oneday_average_scores
		t.DB.Table("midea_loc_code").Where("region_code = ?",regionqushibest[0].Region_code).Find(&tableStore1)
		regionqushibest[0].Dev_region=tableStore1.Dev_region

	t.DB.Table("midea_loc_code").Where("region_code = ?",regionqushiworst[0].Region_code).Find(&tableStore1)
	regionqushiworst[0].Dev_region=tableStore1.Dev_region
	common.IndexDB.Table("city_oneday_average_scores").Select("city_code,round(avg(avg_score),0) as avg_score").Group("city_code").Where(" avg_score>? ",0).Order("avg_score desc").Last(&cityqushi1)
	common.IndexDB.Table("city_oneday_average_scores").Select("city_code,round(avg(avg_score),0) as avg_score").Group("city_code").Where(" avg_score>? ",0).Order("avg_score asc").First(&cityqushi2)
	common.IndexDB.Table("city_oneday_average_scores").Where("city_code=?",cityqushi1[0].City_code).Where("time_date>=? AND time_date<=?",time1,time2).Find(&cityqushibest)
	common.IndexDB.Table("city_oneday_average_scores").Where("city_code=?",cityqushi2[0].City_code).Where("time_date>=? AND time_date<=?",time1,time2).Find(&cityqushiworst)
	var tableStore2 model.City_oneday_average_scores
	t.DB.Table("midea_loc_code").Where("city_code = ?",cityqushibest[0].City_code).Find(&tableStore2)
	cityqushibest[0].Dev_city=tableStore2.Dev_city
	t.DB.Table("midea_loc_code").Where("city_code = ?",cityqushiworst[0].City_code).Find(&tableStore2)
	cityqushiworst[0].Dev_city=tableStore2.Dev_city

	// 全国型号评分排名
var typetotal[2] int
var fenfen[2] int
var avgscore int
	common.IndexDB.Table("type_oneday_avgscores").Select("dev_type,round(avg(avgscore),0) as avgscore").
	Where("time_date>=? AND time_date<=? and avgscore>?",time1,time2,0).
		Group("dev_type").Order("avgscore desc").Find(&typezhuzhuang7)
	common.IndexDB.Table("type_oneday_avgscores").
		Select("dev_type,round(avg(avgscore),0) as avgscore,time_date ").
		Where("time_date>=? AND time_date<=? and avgscore>?",time1,time2,0).
		Group("dev_type").Order("avgscore desc").Count(&typetotal[0])
	for _,table :=range typezhuzhuang7{
		avgscore,_=strconv.Atoi(table.Avgscore)
		fenfen[0]=fenfen[0]+avgscore
	}
	fenfen[0]=fenfen[0]/typetotal[0]
	common.IndexDB.Table("type_oneday_avgscores").Order("time_date desc").Find(&typezhuzhuang1)
	common.IndexDB.Table("type_oneday_avgscores").Select("round(avg(avgscore),0) as avgscore").Where("time_date=?",typezhuzhuang1[0].Time_date).Find(&typeavgscore[1])
	common.IndexDB.Table("type_oneday_avgscores").Select("dev_type,round(avg(avgscore),0) as avgscore").Where("time_date=?  and avgscore>?",typezhuzhuang1[0].Time_date,0).Group("dev_type").Order("avgscore desc").Find(&typezhuzhuang1)
    typetotal[1]=len(typezhuzhuang1)
	for _,table :=range typezhuzhuang1{
		avgscore,_=strconv.Atoi(table.Avgscore)
		fenfen[1]=fenfen[1]+avgscore
	}
	fenfen[1]=fenfen[1]/typetotal[1]
	typeavgscore[0].Avgscore=strconv.Itoa(fenfen[0])
	typeavgscore[1].Avgscore=strconv.Itoa(fenfen[1])
   //全国区域及省份榜单
	common.IndexDB.Table("region_oneday_average_scores").Select("region_code,round(avg(avg_score),0) as avg_score").Group("region_code").Order("avg_score desc").Limit(3).Find(&region1)
	for _, tableDate := range region1 {
		t.DB.Table("midea_loc_code").Where("region_code=?", tableDate.Region_code).Limit(1).Find(&tableStoreDates3)
		region = append(region, tableStoreDates3.Dev_region)

	}
//迁移用时趋势
var qianyishijian struct{
	Date   string `json:"date" gorm:"type:varchar(255);not null"`
}
t.DB.Table("migration_information_record").Order("date desc").
		Limit(30).Find(&qianyishijian)
	t.DB.Table("migration_information_record").Where("date>?",qianyishijian.Date).
		Find(&transfer)
	for _, tableDate := range transfer {
		transfer_time = append(transfer_time , tableDate.Total_time)
		transfer_xdata = append(transfer_xdata , tableDate.Date)
	}
	//挖掘用时趋势
	 //update_time:="2022-02-19"
	common.IndexDB.Table("data_features").
	 	Order("update_time desc").
	 	Limit(30).Find(&excavate)
	for _, tableDate := range excavate {
		excavate_time = append(excavate_time , tableDate.Total_processing_time)
		excavate_xdata = append(excavate_xdata, tableDate.Update_time)
	}
	fmt.Println("excavate_xdata ",excavate_xdata )
	fmt.Println("num",len(excavate ))
	fmt.Println(typequshiworst)
	response.Success(ctx, gin.H{"excavate_xdata":excavate_xdata,"excavate_time":excavate_time,"transfer_xdata":transfer_xdata,"transfer_time":transfer_time,"data4": feature2,"daydata":daydata,"migratenum": migratenum,"monitornum":monitornum,"new_create":new_create,
		"data":water_frag,"region":region,"zhuzhuangtu7":typezhuzhuang7,"zhuzhuangtu1":typezhuzhuang1,"cityqushiworst":cityqushiworst,"cityqushibest":cityqushibest,"regionqushiworst":regionqushiworst,"regionqushibest":regionqushibest,
		"typequshiworst":typequshiworst,"typequshibest":typequshibest,"china_ydata":china_ydata,"province":province,"equipment_num":equipment_num,"tempscore":tempscore,"waterscore":waterscore,"province1":province1,"equipment_num1":equipment_num1,"tempscore1":tempscore1,"waterscore1":waterscore1,
		"feature":feature,"typetotal":typetotal,"typeavgscore1":typeavgscore[0].Avgscore,"typeavgscore2":typeavgscore[1].Avgscore,"ranking":National_ranking,"province_score":score_province,"equ":score_equipment,"monitor_fragement":monitor_fragement,
		"monitor_fragement_city":monitor_fragement_city,"abnormal_count":abnormal_count,"fault_monitor":fault_fragement,"fault_city":fault_fragement_city}, "成功")
}

	func NewmenuController ()ImenuController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return menuController{DB:db}
}