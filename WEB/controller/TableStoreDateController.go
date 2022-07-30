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
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ITaleStoreDateController interface{
	RestController
	PageList(ctx *gin.Context)
}

type TableStoreDateController struct {
	DB *gorm.DB
}

func (t TableStoreDateController) Create(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateTableStoreRequest
	//数据验证
	if err:=ctx.ShouldBind(&requestTableStoreDate);err!=nil{
		log.Println(err.Error())
		response.Fail(ctx,nil,"数据验证错误，分类名称必填")
		return
	}
	//创建tablestoredate
	tableStoreDate:=model.TableDate{
		//CategoryId: requestTableStoreDate.CategoryId,
		Datatime:      requestTableStoreDate.Datatime,
		Flame:    requestTableStoreDate.Flame,
		Outtemp: requestTableStoreDate.Outtemp,
		Settemp: requestTableStoreDate.Settemp,
		//Flow: requestTableStoreDate.Flow,
		//Model: requestTableStoreDate.Model,
	}
	if err:=t.DB.Create(&tableStoreDate).Error;err!=nil{
		panic(err)
		return
	}
	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"创建成功")

}

func (t TableStoreDateController) Update(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateTableStoreRequest
	//数据验证
	if err:=ctx.ShouldBind(&requestTableStoreDate);err!=nil{
		log.Println(err.Error())
		response.Fail(ctx,nil,"数据验证错误，分类名称必填")
		return
	}


	//获取path中的id
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	//更新文章
	if err:=t.DB.Model(&tableStoreDate).Update(requestTableStoreDate).Error;err!=nil{
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx,nil,"更新失败")
		return
	}


	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"更新成功")


}

func (t TableStoreDateController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t TableStoreDateController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func task() (){
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func (t TableStoreDateController) PageList(ctx *gin.Context) {
	var yuefen struct{
		TimeDate string `json:"time_date" gorm:"type:varchar(255);not null"`
	}
	common.IndexDB.Table("month_summaries").Order("time_date desc").Last(&yuefen)
	yue:=yuefen.TimeDate
fmt.Println(yue)


	//获取参数
	//pageNum,_:=strconv.Atoi(ctx.DefaultQuery("pageNum","1"))
	//pageSize,_:=strconv.Atoi(ctx.DefaultQuery("pageSize","200000"))
	category:=ctx.DefaultQuery("category","0000")
	dev_id:=ctx.DefaultQuery("dev_id","109951162778761")
	timeLow:=ctx.DefaultQuery("timeLow","2020-01-01")
	timeHigh:=ctx.DefaultQuery("timeHigh","2022-01-01")
	//averagelow := ctx.DefaultQuery("averagelow", "0")
	//averagehigh := ctx.DefaultQuery("averagehigh", "150")
	model2:=ctx.DefaultQuery("model2","1")
	model3 := ctx.DefaultQuery("model", "0")
	timelength := ctx.DefaultQuery("timelength", "0")
	month1:=ctx.DefaultQuery("month","1")+"-01"
	month2:=ctx.DefaultQuery("month1","1")+"-01"
	month3:=ctx.DefaultQuery("month2","1")
	//month4:=ctx.DefaultQuery("month3","1")
	pageType:=ctx.DefaultQuery("pagetype","2")
	categoryType:=ctx.QueryArray("type1[]")
	city:=ctx.DefaultQuery("city","0")
	city1:=ctx.DefaultQuery("city1","0")
	flag:=ctx.DefaultQuery("flag","0")
	flag1:=ctx.DefaultQuery("flag1","0")
	scorelow:= ctx.DefaultQuery("scorelow", "0")
	scorehigh:= ctx.DefaultQuery("scorehigh", "99")
	equipment_flag:=ctx.DefaultQuery("equipment_flag","0")

	//citySummaryFlag,_:=strconv.ParseBool(ctx.DefaultQuery("citySummaryFlag", "false"))
	//provinceSummaryFlag,_:= strconv.ParseBool(ctx.DefaultQuery("provinceSummaryFlag", "false"))
	//typeSummaryFlag,_:= strconv.ParseBool(ctx.DefaultQuery("typeSummaryFlag", "false"))
	//yearflag,_:= strconv.Atoi(ctx.DefaultQuery("yearflag", "0"))
	//monthflag,_:= strconv.Atoi(ctx.DefaultQuery("monthflag", "1"))
	//dayflag,_:= strconv.Atoi(ctx.DefaultQuery("dayflag", "0"))
	//start:= ctx.DefaultQuery("start", "2021-01-01")
	//end:= ctx.DefaultQuery("end", "2021-07-01")
	//Flag,_:=strconv.ParseBool(ctx.DefaultQuery("Flag", "false"))
	//timeStamp:= ctx.DefaultQuery("timeStamp", "15:02")
	//timeLowInt:=util.TimeParse(timeLow)
	//timeHighInt:=util.TimeParse(timeHigh)


	var datatime []string//时间轴
	var flow []int//水流量
	var flame []string//火焰反馈
	var settemp []string//设定温度
	var outtemp []string//输出温度
	var model1  []int
	var zone_id [] string
	//特征指标
	var Total_time []string
	var Min_time [] string
	var Max_time [] string
	var Total_num  [] int
	//单月份指标


	var modelname=[]string{"stablemode", "smallfluctuationmode", "upmode","downmode","","oscillatemode"}
	//var monthname=[]string{"jan","feb","mar","apr","may","jun","jul","aug","sept","oct","nove","dec"}
	var modelSelect int

	//分页
	var tableStoreDates []model.TableDate//"112150186047071"
	var tableStoreDates1 []model.TableDate1//"112150186047071"
	//var tableStoreDates2 []model.TableDate2
	var tableStoreDates3 []model.TableDate3
	//fmt.Println(categoryType)


	//var tableName string
	month1=strings.ToLower(month1)
	modelSelect, _ = strconv.Atoi(model2)
	//i, _ = strconv.Atoi(model2)


	model2="month_"+modelname[modelSelect]+"_info"
	//category1:=strings.Split(category,"_")

	//fmt.Println(model2)

	var tableStoreDates4 []model.TableDate3
	var tableplace []model.Tableplace
	var tableplace1 []model.Tableplace
	var feature model.Datafeature
	var feature1 model.Datafeature
	//var tablefragment []model.Tablefragment
	//var tablefragment1 [][]model.Tablefragment
	var Stable_proportion[] string
	var Un_stable_proportion[] float32
	//var Un_stable_behavior[] int
	var province[] string
	var province_code[] string
	var city_code[]  string
	var Dev_type[] string
	var time_date[] string
	//var equipment_search []model.Tableplace
	//var biaozhi int
	//if pageType =="test"{
	//	t.DB.Table("midea_loc_code").Where("city_code = ? ",city1).Find(&tableStoreDates4)
	//	for _,tableDate :=range tableStoreDates4 {
	//		city_code = append(city_code, tableDate.City_code)
	//		province_code = append(province_code, tableDate.Province_code)
	//		Dev_type = append(Dev_type, tableDate.Dev_city)
	//		province = append(province, tableDate.Dev_province)
	//	}
	//	fmt.Println(Dev_type)
	//	response.Success(ctx,gin.H{"city_code":city_code,"province_code":province_code,"province":province,"dev_type":Dev_type},"成功")
	//}else
	//if pageType =="chuli"{
	//	//可并发运行多个任务
	//	//Flag := false
	//
	//	if Flag == false {
	//		indexWriteFlag := true
	//		dataWriteFlag := true
	//		fragmentWriteFlag := true
	//		////timeStamp := "15:02"
	//		////
	//		yearsFlag := 0
	//		monthsFlag := 0
	//		DayFlag := 1
	//		timeNow := time.Now()
	//		dataStartTime := timeNow.Add(-time.Duration(24) * time.Hour).Format("2006-01-02")
	//		dataEndTime := timeNow.Format("2006-01-02")
	//		////
	//		gocron.Every(1).Day().At(timeStamp).Do(houduan.DataMining, dataStartTime, dataEndTime, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearsFlag, monthsFlag, DayFlag)
	//
	//		//fmt.Println(1)
	//		//gocron.Clear()
	//		//gocron.Every(1).Day().At(timeStamp).Do(task())
	//		<-gocron.Start()
	//	} else {
	//		//dataStartTime := "2021-01-01"
	//		//dataEndTime := "2021-07-01"
	//		indexWriteFlag := true
	//		dataWriteFlag := true
	//		fragmentWriteFlag := true
	//		//yearsFlag := 0
	//		//monthsFlag := 1
	//		//DayFlag := 0
	//        //fmt.Print(yearflag)
	//		//fmt.Print(monthflag)
	//		//fmt.Print(dayflag)
	//		houduan.DataMining(start, end, indexWriteFlag, dataWriteFlag, fragmentWriteFlag, yearflag, monthflag, dayflag)
	//	}
	//
	//}
	//if pageType =="summary"{
	//
	//	houduan.Summary(citySummaryFlag,provinceSummaryFlag,typeSummaryFlag,yearflag, monthflag, dayflag,start,end)
	//}
	if pageType =="selectAll"{
		t.DB.Table("midea_loc_code").Find(&tableStoreDates4)
		for _,tableDate :=range tableStoreDates4{
			city_code=append(city_code,tableDate.City_code)
			province_code=append(province_code,tableDate.Province_code)
			var flag bool=false
			for  i:=0;i< len(province);i++{
				if(province_code[len(province_code)-1]==province[i]){
					flag=true
				}
			}
			if flag==false{
				province = append(province, province_code[len(province_code)-1])
			}
		}
		response.Success(ctx,gin.H{"city_code":city_code,"province_code":province_code,"province":province},"成功")
	}else
	if pageType =="datafeatures"{
		var lastmonth string
		common.IndexDB.Table("data_features").Order("update_time asc").First(&feature1)
		common.IndexDB.Table("data_features").Order("update_time desc").Last(&feature)
		lastmonth=feature.Update_time[0:7]
		fmt.Println("666")
		fmt.Print(feature1)
		fmt.Print(feature.Update_time)
		response.Success(ctx,gin.H{"data":feature,"lastmonth":lastmonth,"data1":feature1},"成功")
	}else
	if pageType =="model2_vue"{
		// fmt.Println("00")
		var tableStoreDate model.TableDate3
		fmt.Println(city)
		if flag=="1"{
	       fmt.Println("进来了")
			common.IndexDB.Table("month_summaries").Where("dev_type = ? and time_date=?",city,yue).Find(&tableStoreDates4)
			//common.IndexDB.Table("month_summaries").Raw("select distinct dev_id,distinct dev_type where dev_type=?",city).Find(&tableStoreDates4)
			//Where("dev_type = ? and handle_flag=?",city,"1").Find(&tableStoreDates4)
			for _,tableDate1 :=range tableStoreDates4{
				var flag bool=false
				for  i:=0;i< len(city_code);i++{
					if(tableDate1.City_code==city_code[i]){
						flag=true
					}
				}
				if flag==false{
					city_code = append(city_code,tableDate1.City_code)
					t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDate)
					province=append(province,tableStoreDate.Dev_city)
				}

			}

		}else{
			if flag=="2"{
				fmt.Println("进来了2")
				var province1  []string
				t.DB.Table("midea_loc_code").Where("province_code =?",city).Find(&tableStoreDates4)
				for _, tableDate1 := range tableStoreDates4 {
					province1=append(province1,tableDate1.City_code)
				}
				common.IndexDB.Table("month_summaries").Where("city_code in (?) and time_date=?",province1,yue).Find(&tableStoreDates1)
				//t.DB.Table("bo").Where("city_code in (?) and handle_flag=?",province1,"1").Find(&tableStoreDates1)
			} else{
				fmt.Println("进来了3",city,yue)
				
				common.IndexDB.Table("month_summaries").Where("city_code = ? and time_date=?",city,yue).Find(&tableStoreDates1)
				fmt.Println("结果是",tableStoreDates1)
				//t.DB.Table("bo").Where("city_code = ? and handle_flag=? ",city,"1").Find(&tableStoreDates1)
			}

			for _,tableDate1 :=range tableStoreDates1{
				var flag bool=false
				for  i:=0;i< len(Min_time);i++{
					if(tableDate1.Dev_type==Min_time[i]){
						flag=true
					}
				}
				if flag==false{
					Min_time = append(Min_time,tableDate1.Dev_type)
				}
				Max_time = append(Max_time,tableDate1.Dev_type)

				//if equipment_flag=="0"{
				Total_time= append(Total_time,tableDate1.Dev_Id)
			//	}
				//if equipment_flag=="1"{
				//	t.DB.Table("bo").Where("city_code = ? ",city).Find(&equipment_search)
				//	for  i:=0;i<len(equipment_search);i++{
					
					
					//Total_time= append(Total_time,equipment_search[i].Dev_Id)}
				//}
				
			}
fmt.Println("equipment_flag=",equipment_flag)
		}

			fmt.Println(Total_time)
		response.Success(ctx,gin.H{"devid":Total_time,"devtype":Min_time,"total":Max_time,"city_code":city_code,"province":province},"成功")
	} else
	if pageType=="localized_vue"{
		if flag1 =="1"{
			common.IndexDB.Table("data_features").Order("update_time desc").Last(&feature)
			month1= feature.Update_time[0:7]+"-01"
		}
		if category=="city"{
			if flag=="1"{common.IndexDB.Table("month_summaries").Where("time_date=? AND city_code=?",month1,city).Find(&tableplace)
			}else{common.IndexDB.Table("month_summaries").Where("time_date=? AND city_code=? AND dev_type=?",month1,city,dev_id).Find(&tableplace)}
			for _,tableDate :=range tableplace{
				Dev_type=append(Dev_type,tableDate.Dev_type)//数据放入切片中
				Stable_proportion=append(Stable_proportion,tableDate.Stable_proportion)//数据放入切片中
				Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
				city_code=append(city_code,tableDate.City_code)
				t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
				for _,tableDate1 :=range tableStoreDates4 {
					province = append(province, tableDate1.Dev_city)
				}
			}
			//	fmt.Println(province)
			response.Success(ctx,gin.H{"data":tableplace,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province,"dev_type":Dev_type,"month1":month1[0:7]},"成功")
		}else if category=="province"{
			common.IndexDB.Table("city_months").Where("time_date=? AND province_code=? ",month1,city).Find(&tableplace)
			//fmt.Println(tableplace)
			var city1=city[0:3]+"100"
			common.IndexDB.Table("month_summaries").Where("time_date=? AND city_code=? ",month1,city1).Find(&tableplace1)
			for _,tableDate :=range tableplace{
				Dev_type=append(Dev_type,tableDate.Dev_type)//数据放入切片中
				Stable_proportion=append(Stable_proportion,tableDate.Stable_proportion)//数据放入切片中
				Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
				city_code=append(city_code,tableDate.City_code)
				t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
				//for _,tableDate1 :=range tableStoreDates4 {
				var flag bool=false
				for  i:=0;i< len(province);i++{
					if(tableStoreDates4[0].Dev_city==province[i]){
						flag=true
					}
				}
				if flag==false{
					province = append(province, tableStoreDates4[0].Dev_city)
				}
			}
			response.Success(ctx,gin.H{"data":tableplace,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province,"dev_type":Dev_type,"data1":tableplace1,"month1":month1[0:7]},"成功")
		}else if category=="month"{
			common.IndexDB.Table("days_summaries").Where("dev_id=?  AND time_date BETWEEN ? AND ?",dev_id,month1,month2).Find(&tableplace)
			// fmt.Println(tableplace)
			for _,tableDate :=range tableplace{
				Dev_type=append(Dev_type,tableDate.Dev_type)//数据放入切片中
				Stable_proportion=append(Stable_proportion,tableDate.Stable_proportion)//数据放入切片中
				Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
				city_code=append(city_code,tableDate.City_code)
				time_date=append(time_date,tableDate.Time_date)
			}
			response.Success(ctx,gin.H{"data":tableplace,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province,"dev_type":Dev_type,"data1":tableplace1,"time_date":time_date,"month1":month1[0:7]},"成功")
		}else if category=="day"{
			var tablename string
			fmt.Println("time=",timeLow)
			tablename="data"+city+"_"+timeLow[0:4]+timeLow[5:7]+timeLow[8:10]
			fmt.Println("tablename=",tablename)
			fmt.Println(dev_id)
			common.RunDB.Table(tablename).Where(" applianceid=?",dev_id).Find(&tableStoreDates)
			//common.RunDB.Table("data110100_20211020").Where(" applianceid=?",dev_id).Find(&tableStoreDates)
				fmt.Println(tableStoreDates)
			common.IndexDB.Table("days_summaries").Where("dev_id=?  AND time_date=?",dev_id,month3).Find(&tableplace)
			for _,tableDate :=range tableStoreDates{
				datatime=append(datatime,tableDate.Datatime)//数据放入切片中
				flow=append(flow,tableDate.Flow)
				flame=append(flame,tableDate.Flame)
				settemp=append(settemp,tableDate.Settemp)
				outtemp=append(outtemp,tableDate.Outtemp)
				model1=append(model1,tableDate.Water_pattern)
				zone_id=append(zone_id,tableDate.Zone_id)
			}
			//	response.Success(ctx,gin.H{"data":tableStoreDates,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province},"成功")
			fmt.Println(tableStoreDates)
			response.Success(ctx,gin.H{"data":tableStoreDates,"data_time":datatime,"flow":flow,"flame":flame,"set_temp":settemp,"out_temp":outtemp,"model":model1,
				"zone_id":zone_id,"data1":tableplace},"成功")
		} else {
			if flag=="1" {
				//var id3 [] string
			}else{
				common.IndexDB.Table("province_months").Where("time_date=? ",month1).Find(&tableplace)
				for _,tableDate :=range tableplace{
					Stable_proportion=append(Stable_proportion,tableDate.Stable_proportion)//数据放入切片中
					Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
					province_code=append(city_code,tableDate.Province_code)
					t.DB.Table("midea_loc_code").Where("province_code=?",province_code[len(province_code)-1]).Find(&tableStoreDates4)
					//for _,tableDate1 :=range tableStoreDates4 {
					//	province = append(province, tableDate1.Dev_province)
					var flag bool=false
					for  i:=0;i< len(province);i++{
						if(tableStoreDates4[0].Dev_province==province[i]){
							flag=true
						}
					}
					if flag==false{
						province = append(province, tableStoreDates4[0].Dev_province)
					}
				}
				response.Success(ctx,gin.H{"data":tableplace,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province,"month1":month1[0:7]},"成功")}
		}
	} else
	if pageType=="locationsummary"{
		var province1 []  string
		if flag=="provincesum"{
			for i:=0;i< len(categoryType);i++{
				//area1:=strings.Split(categoryType[i]," ")
				id2:=strings.Split(categoryType[i],"]")[0]
				id2=id2[2:len(id2)-1]
				//id3=append(id3,id2)
				//id3[i]=id2
				//id3:=strings.Split(id2,"_")
				//fmt.Println(tableStoreDates2)
				common.IndexDB.Table("province_months").Where("province_code=? And time_date BETWEEN ? AND ?",id2,month1,month2).Find(&tableplace)
				for i:=0;i< len(tableplace);i++{
					tableplace1=append(tableplace1,tableplace[i])
				}
				for _,tableDate :=range tableplace{
					//Un_stable_behavior=append(Un_stable_behavior,tableDate.Un_stable_behavior)//数据放入切片中
					//Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
					province_code=append(city_code,tableDate.Province_code)
					t.DB.Table("midea_loc_code").Where("province_code=?",province_code[len(province_code)-1]).Find(&tableStoreDates4)
					//for _,tableDate1 :=range tableStoreDates4 {
					//	province = append(province, tableDate1.Dev_province)
					province = append(province, tableStoreDates4[0].Dev_province)
					var flag bool=false
					for  i:=0;i< len(province1);i++{
						if(tableStoreDates4[0].Dev_province==province1[i]){
							flag=true
						}
					}
					if flag==false{
						province1 = append(province1, tableStoreDates4[0].Dev_province)
					}
				}
				for i:=0;i< len(tableplace1);i++{
					tableplace1[i].Time_date=tableplace1[i].Time_date[0:7]
				}
			}
			//fmt.Println(id3)
			response.Success(ctx,gin.H{"data":tableplace1,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province,"province1":province1},"成功")
			//fmt.Println(categoryType)
		}else {
			for i:=0;i< len(categoryType);i++{
				//area1:=strings.Split(categoryType[i]," ")
				id2:=strings.Split(categoryType[i],"]")[0]
				//id3=append(id3,id2)
				//id3[i]=id2
				//id3:=strings.Split(id2,"_")
				//fmt.Println(tableStoreDates2)
				common.IndexDB.Table("city_months").Where("city_code=? And time_date BETWEEN ? AND ?",id2,month1,month2).Find(&tableplace)
				for i:=0;i< len(tableplace);i++{
					tableplace1=append(tableplace1,tableplace[i])
				}
				for _,tableDate :=range tableplace{
					Stable_proportion=append(Stable_proportion,tableDate.Stable_proportion)//数据放入切片中
					Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
					city_code=append(city_code,tableDate.City_code)
					t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
					//for _,tableDate1 :=range tableStoreDates4 {
					//	province = append(province, tableDate1.Dev_province)
					province = append(province, tableStoreDates4[0].Dev_city)
					var flag bool=false
					for  i:=0;i< len(province1);i++{
						if(tableStoreDates4[0].Dev_city==province1[i]){
							flag=true
						}
					}
					if flag==false{
						province1 = append(province1, tableStoreDates4[0].Dev_city)
					}
				}
			}
			//fmt.Println(id3)
			for i:=0;i< len(tableplace1);i++{
				tableplace1[i].Time_date=tableplace1[i].Time_date[0:7]
			}
			response.Success(ctx,gin.H{"data":tableplace1,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province,"province1":province1},"成功")
			//fmt.Println(categoryType)
		}
	}else
	if pageType=="typesummary" {
		var timenew struct{
			Time_date      string   `json:"time_date" gorm:"type:varchar(50);not null"`
		}
		if flag=="1"{
			common.IndexDB.Table("type_months").Order("time_date desc").Last(&timenew)
			fmt.Println("timenew=",timenew)
			common.IndexDB.Table("type_months").Where("time_date = ?",timenew.Time_date).Find(&tableplace)
			response.Success(ctx,gin.H{"data":tableplace},"成功")
		}else{
			for i:=0;i< len(categoryType);i++{
				//area1:=strings.Split(categoryType[i]," ")
				id2:=strings.Split(categoryType[i],"]")[0]
				id2=id2[2:len(id2)-1]
				//id3=append(id3,id2)
				//id3[i]=id2
				//id3:=strings.Split(id2,"_")
				//fmt.Println(tableStoreDates2)
				common.IndexDB.Table("type_months").Where("dev_type=? And time_date BETWEEN ? AND ?", id2,month1,month2).Find(&tableplace)
				for i:=0;i< len(tableplace);i++{
					tableplace1=append(tableplace1,tableplace[i])
				}
				for _,tableDate :=range tableplace{
					Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
				}
				if len(tableplace)!=0{
					Dev_type = append(Dev_type, tableplace[0].Dev_type)
				}

			}
			for i:=0;i< len(tableplace1);i++{
				tableplace1[i].Time_date=tableplace1[i].Time_date[0:7]
			}
			response.Success(ctx,gin.H{"data":tableplace1,"dev_type":Dev_type},"成功")
			//fmt.Println(categoryType)
		}
	}else
	if pageType=="timesummary"{
		var month []  string
		var province1 []  string
		var Dev_Id []  string
		if flag=="1"{
			if flag1=="1"{
				if category=="1"{
					common.IndexDB.Table("province_months").Where("province_code=? AND time_date BETWEEN ? AND ?",city1,month1,month2).Find(&tableplace)
					for _,tableDate :=range tableplace {
						Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
						month = append(month, tableDate.Time_date)
						province_code=append(city_code,tableDate.Province_code)
						t.DB.Table("midea_loc_code").Where("province_code=?",province_code[len(province_code)-1]).Find(&tableStoreDates4)
						province = append(province, tableStoreDates4[0].Dev_province)
					}
					fmt.Print(province)
				}else{
					fmt.Print(month1)
					common.IndexDB.Table("province_days").Where("province_code=? AND time_date BETWEEN ? AND ?",city1,timeLow,timeHigh).Find(&tableplace)
					for _,tableDate :=range tableplace {
						Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
						month = append(month, tableDate.Time_date)
						city_code=append(city_code,tableDate.Province_code)
						t.DB.Table("midea_loc_code").Where("province_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
						province = append(province, tableStoreDates4[0].Dev_province)
					}
				}
			}else{
				if category=="1"{
					common.IndexDB.Table("city_months").Where("city_code=? AND time_date BETWEEN ? AND ?",city,month1,month2).Find(&tableplace)
					for _,tableDate :=range tableplace {
						Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
						month = append(month, tableDate.Time_date)
						city_code=append(city_code,tableDate.City_code)
						t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
						province = append(province, tableStoreDates4[0].Dev_city)
						province1 = append(province1, tableStoreDates4[0].Dev_province)
					}
				}else{
					common.IndexDB.Table("city_days").Where("city_code=? AND time_date BETWEEN ? AND ?",city,timeLow,timeHigh).Find(&tableplace)
					for _,tableDate :=range tableplace {
						Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
						month = append(month, tableDate.Time_date)
						city_code=append(city_code,tableDate.City_code)
						t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
						province = append(province, tableStoreDates4[0].Dev_city)
						province1 = append(province1, tableStoreDates4[0].Dev_province)
					}
				}
			}
		}else if flag=="2"{
			if flag1=="1"{
				common.IndexDB.Table("type_months").Where("dev_type=? And time_date BETWEEN ? AND ?",city,month1,month2).Find(&tableplace)
				for _,tableDate :=range tableplace {
					Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
					month = append(month, tableDate.Time_date)
					Dev_type=append(Dev_type,tableDate.Dev_type)
				}
			}else{
				common.IndexDB.Table("type_days").Where("dev_type=? And time_date BETWEEN ? AND ?",city,timeLow,timeHigh).Find(&tableplace)
				for _,tableDate :=range tableplace {
					Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
					month = append(month, tableDate.Time_date)
					Dev_type=append(Dev_type,tableDate.Dev_type)
				}
			}
		}else{
			if flag1=="1"{
				common.IndexDB.Table("month_summaries").Where("city_code=? AND dev_type=? AND dev_id=? AND time_date BETWEEN ? AND ?",city,category,dev_id,month1,month2).Find(&tableplace)
			}else{
				common.IndexDB.Table("days_summaries").Where("city_code=? AND dev_type=? AND dev_id=? AND time_date BETWEEN ? AND ?",city,category,dev_id,timeLow,timeHigh).Find(&tableplace)
			}
			for _,tableDate :=range tableplace {
				Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
				month = append(month, tableDate.Time_date)
				Dev_type=append(Dev_type,tableDate.Dev_type)
				Dev_Id=append(Dev_Id,tableDate.Dev_Id)
				city_code=append(city_code,tableDate.City_code)
				t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
				province = append(province, tableStoreDates4[0].Dev_city)
				province1 = append(province1, tableStoreDates4[0].Dev_province)
			}
		}
		//for _,tableDate :=range tableplace{
		//	Un_stable_proportion=append(Un_stable_proportion,tableDate.Un_stable_proportion)
		//	month=append(month,tableDate.Time_date)
		//	city_code=append(city_code,tableDate.City_code)
		//	Dev_type=append(Dev_type,tableDate.Dev_type)
		//t.DB.Table("midea_loc_code").Where("city_code=?",city_code[len(city_code)-1]).Find(&tableStoreDates4)
		////for _,tableDate1 :=range tableStoreDates4 {
		////	province = append(province, tableDate1.Dev_province)
		//var flag bool=false
		//for  i:=0;i< len(province);i++{
		//	if(tableStoreDates4[0].Dev_city==province[i]){
		//		flag=true
		//	}
		//
		//}
		//if flag==false{
		//
		//	province = append(province, tableStoreDates4[0].Dev_city)
		//}
		//}
		//fmt.Println(Un_stable_proportion)
		////fmt.Println(Dev_type)
		if category=="1"||(flag!="1"&&flag1=="1"){
			for i:=0;i< len(tableplace);i++{
				tableplace[i].Time_date=tableplace[i].Time_date[0:7]
				month[i]=month[i][0:7]
			}
		}
		response.Success(ctx,gin.H{"data":tableplace,"un_stable_proportion":Un_stable_proportion,"month":month,"province":province,"dev_type":Dev_type,"province1":province1,"dev_id":Dev_Id},"成功")
	}else
	if pageType == "indexdregion" {

		//fmt.Println(scorehigh)
		if flag == "1" {
			var tablename string
			tablename = "data" + city+"_"+timeLow[0:4]+timeLow[5:7]+timeLow[8:10]
			common.RunDB.Table(tablename).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_id,timeLow,timeHigh).Find(&tableStoreDates)
			//	fmt.Println(tableStoreDates)

			for _, tableDate := range tableStoreDates {

				datatime = append(datatime, tableDate.Datatime) //数据放入切片中
				flow = append(flow, tableDate.Flow)
				flame = append(flame, tableDate.Flame)
				settemp = append(settemp, tableDate.Settemp)
				outtemp = append(outtemp, tableDate.Outtemp)
				model1 = append(model1, tableDate.Water_pattern)
				zone_id = append(zone_id, tableDate.Zone_id)

			}

			//	response.Success(ctx,gin.H{"data":tableStoreDates,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province},"成功")
			response.Success(ctx, gin.H{"data": tableStoreDates, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp, "model": model1,
				"zone_id": zone_id, "data1": tableplace, "total_time": Total_time, "min_time": Min_time, "max_time": Max_time, "total_num": Total_num,
				"datanum": tableStoreDates3}, "成功")

		} else if flag == "2" {
			timeLow=timeLow+" 00:00:00"
			timeHigh=timeHigh+" 23:59:59"
			var tablefragment []model.Tablefragment
			var city1 = "fragment" + city
			//fmt.Println(averagelow)
			//fmt.Println(averagehigh)
			if model3 == "0" {
				if timelength == "1" {

					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND " +
						"water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id=?  AND start_time BETWEEN  ? AND ? " +
						"water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id=? AND start_time BETWEEN  ? AND ?" +
						" AND water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id=?  AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id=? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}

			}else  if model3=="5"{
				if timelength == "1" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
						" water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
						" AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}


			}else  if model3=="1"{
				if timelength == "1" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
						"  water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
						" AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}


			} else {
				if timelength == "1" {

					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? " +
						"AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {

					common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}
			}


			//	common.IndexDB.Table(city1).Where("dev_id=? AND pattern=? AND start_time BETWEEN  ? AND ?",dev_id,model3,timeLow,timeHigh).Find(&tablefragment)
			//for _, tableDate := range tablefragment {
			//
			//	datatime = append(datatime, tableDate.Multiple_score) //数据放入切片中
			//
			//
			//}
			//fmt.Println(datatime)
			//	response.Success(ctx,gin.H{"data":tablefragment},"成功")})=
			response.Success(ctx, gin.H{"data": tablefragment}, "成功")

		}else if flag == "3" {
			dev_type:=ctx.DefaultQuery("dev_type","0000")
			timeLow=timeLow+" 00:00:00"
			timeHigh=timeHigh+" 23:59:59"
			var tablefragment []model.Tablefragment
			var city1 = "fragment" + city
			var Dev_id [] string
			//fmt.Println(averagelow)
			//fmt.Println(averagehigh)
fmt.Println(dev_type,city,yue)
			common.IndexDB.Table("month_summaries").Where("dev_type = ? AND city_code = ? and time_date=?",dev_type,city,yue).Find(&tableStoreDates4)
fmt.Println("jieguoshi=",tableStoreDates4)
			//t.DB.Table("bo").Where("dev_type = ? AND city_code = ? and handle_flag=?",dev_type,city,"1").Find(&tableStoreDates4)
			for _,tableDate :=range tableStoreDates4{
				Dev_id=append(Dev_id,tableDate.Dev_Id)
			}
			if model3 == "0" {
				if timelength == "1" {

					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND " +
						"water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id in (?)  AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND start_time BETWEEN  ? AND ?" +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id in (?)  AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}

			}else  if model3=="5"{
				if timelength == "1" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
						" water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}


			}else  if model3=="1"{
				if timelength == "1" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
						" water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}


			} else {
				if timelength == "1" {

					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				} else if timelength == "2" {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "3" {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? " +
						"AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else if timelength == "4" {
					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						" AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
					//fmt.Println(tablefragment)
				} else {

					common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
						"AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
				}
			}


			//	common.IndexDB.Table(city1).Where("dev_id=? AND pattern=? AND start_time BETWEEN  ? AND ?",dev_id,model3,timeLow,timeHigh).Find(&tablefragment)
			//for _, tableDate := range tablefragment {
			//
			//	datatime = append(datatime, tableDate.Multiple_score) //数据放入切片中
			//
			//
			//}
			//fmt.Println(datatime)
			//	response.Success(ctx,gin.H{"data":tablefragment},"成功")})=
			if len(tablefragment)>1000&&flag1=="0"{
				ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"数据量超过1000条(建议具体查询条件)"})
				return
			}else{
				response.Success(ctx, gin.H{"data": tablefragment}, "成功")
			}

		}}else
	if pageType == "indexbehavior" {

		//fmt.Println(scorehigh)
		if flag == "1" {
			var tablename string
			tablename = "data" + city+"_"+timeLow[0:4]+timeLow[5:7]+timeLow[8:10]

			common.RunDB.Table(tablename).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_id,timeLow,timeHigh).Find(&tableStoreDates)
			//	fmt.Println(tableStoreDates)

			for _, tableDate := range tableStoreDates {

				datatime = append(datatime, tableDate.Datatime) //数据放入切片中
				flow = append(flow, tableDate.Flow)
				flame = append(flame, tableDate.Flame)
				settemp = append(settemp, tableDate.Settemp)
				outtemp = append(outtemp, tableDate.Outtemp)
				model1 = append(model1, tableDate.Water_pattern)
				zone_id = append(zone_id, tableDate.Zone_id)

			}
			//	response.Success(ctx,gin.H{"data":tableStoreDates,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province},"成功")
			// fmt.Println(len(datatime))
			response.Success(ctx, gin.H{"data": tableStoreDates, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp, "model": model1,
				"zone_id": zone_id, "data1": tableplace, "total_time": Total_time, "min_time": Min_time, "max_time": Max_time, "total_num": Total_num,
				"datanum": tableStoreDates3}, "成功")

		} else if flag == "2" {
			timeLow=timeLow+" 00:00:00"
			timeHigh=timeHigh+" 23:59:59"
			var tablefragment []model.Tablewaterbehavior
			//fmt.Println(averagelow)
			//fmt.Println(averagehigh)
			if model3 == "0" {
				common.IndexDB.Table("mode_behaviors").Where("dev_id=?   AND start_time BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh).Find(&tablefragment)
			}else  if model3=="1"{
				common.IndexDB.Table("mode_behaviors").Where("dev_id=?  AND start_time BETWEEN  ? AND ? AND effect_flag = ?", dev_id, timeLow, timeHigh,"1").Find(&tablefragment)
			}else  if model3=="2"{
				common.IndexDB.Table("mode_behaviors").Where("dev_id=?  AND start_time BETWEEN  ? AND ? AND effect_flag = ?", dev_id, timeLow, timeHigh,"0").Find(&tablefragment)
			}
			//	common.IndexDB.Table(city1).Where("dev_id=? AND pattern=? AND start_time BETWEEN  ? AND ?",dev_id,model3,timeLow,timeHigh).Find(&tablefragment)
			//for _, tableDate := range tablefragment {
			//
			//	datatime = append(datatime, tableDate.Multiple_score) //数据放入切片中
			//
			//
			//}
			//	response.Success(ctx,gin.H{"data":tablefragment},"成功")})=
			response.Success(ctx, gin.H{"data": tablefragment}, "成功")

		}} else
	if pageType == "indexmodeltype" {
		var filter[]string
		time := make(map[string] []string)
		xiaoshu := make(map[string] []float32)
		zhengshu := make(map[string] []int)
		if flag == "1" {
			if category=="0000"{
				common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow,timeHigh, city).Find(&tableplace)
			}else{
				common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow,timeHigh, city,category).Find(&tableplace)
			}

		} else {
			if category=="0000"{
				common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND dev_type=? ",timeLow,timeHigh, dev_id).Find(&tableplace)
			}else{
				common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow,timeHigh, category,dev_id).Find(&tableplace)
			}

		}
		//fmt.Println(tableplace)
		for _, tableDate := range tableplace {
			var flag bool=false
			for  i:=0;i< len(filter);i++{
				if(tableDate.Dev_Id==filter[i]){
					flag=true
				}
			}
			if flag==false{
				filter = append(filter, tableDate.Dev_Id)
				Dev_type = append(Dev_type, tableDate.Dev_type)
				city_code = append(city_code, tableDate.City_code)
				t.DB.Table("midea_loc_code").Where("city_code=?", city_code[len(city_code)-1]).Find(&tableStoreDates4)
				for _, tableDate1 := range tableStoreDates4 {
					province = append(province, tableDate1.Dev_city)
				}
				time[tableDate.Dev_Id]=[]string{"0","0","0"}
				xiaoshu[tableDate.Dev_Id]=[]float32{0.0,0.0}
				zhengshu[tableDate.Dev_Id]=[]int{0,0,0}
			}
			time0,_:=strconv.Atoi(time[tableDate.Dev_Id][0])
			time1,_:=strconv.Atoi(time[tableDate.Dev_Id][1])
			time2,_:=strconv.Atoi(time[tableDate.Dev_Id][2])
			time[tableDate.Dev_Id][0]=strconv.Itoa(time2sec(tableDate.Water_valid_time)+time0)
			time[tableDate.Dev_Id][1]=strconv.Itoa(time2sec(tableDate.Average_time)+time1)
			time[tableDate.Dev_Id][2]=strconv.Itoa(time2sec(tableDate.Maximum_time)+time2)
			xiaoshu[tableDate.Dev_Id][0]=tableDate.Un_stable_proportion+xiaoshu[tableDate.Dev_Id][0]
			xiaoshu[tableDate.Dev_Id][1]=xiaoshu[tableDate.Dev_Id][1]+1
			zhengshu[tableDate.Dev_Id][0]=tableDate.Water_score+zhengshu[tableDate.Dev_Id][0]
			zhengshu[tableDate.Dev_Id][1]=tableDate.Water_num+zhengshu[tableDate.Dev_Id][1]
			zhengshu[tableDate.Dev_Id][2]=zhengshu[tableDate.Dev_Id][2]+1
			//Dev_type = append(Dev_type, tableDate.Dev_type)                            //数据放入切片中
			//Stable_proportion = append(Stable_proportion, tableDate.Stable_proportion) //数据放入切片中
			//Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
			//city_code = append(city_code, tableDate.City_code)
			//t.DB.Table("midea_loc_code").Where("city_code=?", city_code[len(city_code)-1]).Find(&tableStoreDates4)
			//for _, tableDate1 := range tableStoreDates4 {
			//	province = append(province, tableDate1.Dev_city)
			//}

		}
		for key := range time {
			Time0,_:=strconv.Atoi(time[key][0])
			Time1,_:=strconv.Atoi(time[key][1])
			Time2,_:=strconv.Atoi(time[key][2])
			time[key][0]=strconv.Itoa(Time0/zhengshu[key][2])
			time[key][1]=strconv.Itoa(Time1/zhengshu[key][2])
			time[key][2]=strconv.Itoa(Time2/zhengshu[key][2])
			xiaoshu[key][0]=xiaoshu[key][0]/xiaoshu[key][1]
			zhengshu[key][0]=zhengshu[key][0]/zhengshu[key][2]
			zhengshu[key][1]=zhengshu[key][1]/zhengshu[key][2]
		}

		fmt.Println(Dev_type)
		fmt.Println(filter)
		fmt.Println(xiaoshu)
		response.Success(ctx, gin.H{"data": tableplace, "stable_proportion": Stable_proportion, "un_stable_proportion": Un_stable_proportion, "province": province, "dev_type": Dev_type,"filter":filter,"time":time,"xiaoshu":xiaoshu,"zhengshu":zhengshu}, "成功")

	}else
	if pageType =="menuhome"{
		var average_time []int
		var dev_id []string
		if flag =="1"{
			common.IndexDB.Table("month_summaries").Where("time_date =? AND province_code =?",month1,city).Find(&tableplace)
			for _,tableDate1 :=range tableplace{
				average_time=append(average_time,time2sec(tableDate1.Average_time))
				dev_id=append(dev_id,tableDate1.Dev_Id)
			}
		}else{
			//	fmt.Println(timeLow)
			common.IndexDB.Table("province_months").Where("time_date =?",month1).Find(&tableplace)
			//    fmt.Println(tableplace)
			for _,tableDate1 :=range tableplace{
				province_code=append(province_code,tableDate1.Province_code)
				t.DB.Table("midea_loc_code").Where("province_code =?",province_code[len(province_code)-1]).Find(&tableStoreDates4)
				province=append(province,tableStoreDates4[0].Dev_province)
			}
		}
		response.Success(ctx,gin.H{"data":tableplace,"province":province,"average_time":average_time,"dev_id":dev_id},"成功")
	} else
	if pageType == "menu" {
		var province1 []string


		var tabletimerate1 [] model.Tabletimerate
		var fragmentflow    []model.Fragmentflow
		t.DB.Table("midea_loc_code").Find(&tableStoreDates4)

		//循环省份表
		for _, tableDate := range tableStoreDates4 {
			city_code=append(city_code,tableDate.City_code)

			province=append(province,tableDate.Dev_province)

			var flag_1 bool=false
			for  i:=0;i< len(province1);i++{
				if(province[len(province)-1]==province1[i]){
					flag_1=true
				}
			}
			//若省份改变，进行该省份水流量片段分布查询
			if flag_1==false{
				province1 = append(province1, province[len(province)-1])
				province_code=append(province_code,tableDate.Province_code)
				common.IndexDB.Table("province_days").Where("province_code= ? AND time_date BETWEEN ? AND ?", province_code[len(province_code)-1], timeLow, timeHigh).Find(&fragmentflow)

				var three int = 0
				var five int = 0
				var ten int = 0
				var equipment_num int=0
				var tabletimerate model.Tabletimerate
				tabletimerate.Province=province1[len(province1)-1]
				for _, tableDate1 := range fragmentflow {
					three=three+tableDate1.Three
					five=five+tableDate1.Five
					ten=ten+tableDate1.Ten
					//fmt.Print(ten)
					equipment_num=equipment_num+tableDate1.Equipment_num


				}
				var three1 float64 = 0.0
				var five1 float64 = 0.0
				var ten1 float64 = 0.0
				three1=float64(three)/float64(equipment_num)
				five1=float64(five)/float64(equipment_num)
				ten1=float64(ten)/float64(equipment_num)
				//if equipment_num!=0{
				//	three=three/equipment_num
				//	five=five/equipment_num
				//	ten=ten/equipment_num
				//
				//}


				if three!=0&&five!=0&&ten!=0{
					dataparameter1 := make(map[string]float64)
					dataparameter1["three"]=three1
					dataparameter1["five"]=five1
					dataparameter1["ten"]=ten1
					tabletimerate.Duration_time=dataparameter1
					//fmt.Print(tabletimerate)
					tabletimerate.Equipment_num=equipment_num
					tabletimerate1=append(tabletimerate1,tabletimerate)
				}


			}

		}
		//查询全国，汇总计算
		if len(tabletimerate1)!=0{
			var three float64 = 0
			var five float64 = 0
			var ten float64 = 0
			var equipment_num int=0
			for i:=0;i<len(tabletimerate1);i++{
				three=three+tabletimerate1[i].Duration_time["three"]
				five=five+tabletimerate1[i].Duration_time["five"]
				ten=ten+tabletimerate1[i].Duration_time["ten"]
				equipment_num=equipment_num+tabletimerate1[i].Equipment_num
			}
			var tabletimerate model.Tabletimerate
			tabletimerate.Province="全国"
			dataparameter1 := make(map[string]float64)
			dataparameter1["three"]=three
			dataparameter1["five"]=five
			dataparameter1["ten"]=ten
			tabletimerate.Duration_time=dataparameter1
			tabletimerate.Equipment_num=equipment_num
			tabletimerate1=append(tabletimerate1,tabletimerate)
		}



		response.Success(ctx, gin.H{"data1": tabletimerate1}, "成功")
	}else if pageType == "newmenu" {
		var behavior [] model.Behavior
		var Tablebehavior1 [] model.Tablebehavior
		var id [] string
		SheBei := make(map[string] []float32)
		common.IndexDB.Table("behavior_summaries").Where("data_time BETWEEN ? AND ?",timeLow,timeHigh).Find(&behavior)
		//SheBei[behavior[0].Dev_Id]=[]float32{0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0}
		//id[0]=behavior[0].Dev_Id
		for _, tableDate := range behavior {
			var Tablebehavior2 model.Tablebehavior

			var flag_1 bool=false
			for i:=0;i< len(id);i++{
				if id[i]==tableDate.Dev_Id{
					SheBei[tableDate.Dev_Id][0]=SheBei[tableDate.Dev_Id][0]+tableDate.Sec0p
					SheBei[tableDate.Dev_Id][1]=SheBei[tableDate.Dev_Id][1]+tableDate.Sec30p
					SheBei[tableDate.Dev_Id][2]=SheBei[tableDate.Dev_Id][2]+tableDate.Min3p
					SheBei[tableDate.Dev_Id][3]=SheBei[tableDate.Dev_Id][3]+tableDate.Min10p
					flag_1=true
				}
			}
			if flag_1==false{
				id=append(id,tableDate.Dev_Id)
				SheBei[tableDate.Dev_Id]=[]float32{0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0}
				SheBei[tableDate.Dev_Id][0]=tableDate.Sec0p
				SheBei[tableDate.Dev_Id][1]=tableDate.Sec30p
				SheBei[tableDate.Dev_Id][2]=tableDate.Min3p
				SheBei[tableDate.Dev_Id][3]=tableDate.Min10p
			}
			if tableDate.Sec0p==0&&tableDate.Sec30p==0&&tableDate.Min3p==0&&tableDate.Min10p==0{
			}else{
				SheBei[tableDate.Dev_Id][4]=SheBei[tableDate.Dev_Id][4]+1
			}
			SheBei[tableDate.Dev_Id][5]=SheBei[tableDate.Dev_Id][0]/SheBei[tableDate.Dev_Id][4]
			SheBei[tableDate.Dev_Id][6]=SheBei[tableDate.Dev_Id][1]/SheBei[tableDate.Dev_Id][4]
			SheBei[tableDate.Dev_Id][7]=SheBei[tableDate.Dev_Id][2]/SheBei[tableDate.Dev_Id][4]
			SheBei[tableDate.Dev_Id][8]=SheBei[tableDate.Dev_Id][3]/SheBei[tableDate.Dev_Id][4]
			if flag_1==false{
				Tablebehavior2.Dev_id=tableDate.Dev_Id
				common.IndexDB.Table("month_summaries").Where("dev_id = ? and time_date=?",tableDate.Dev_Id,yue).Find(&tableStoreDates4)
				//t.DB.Table("bo").Where("dev_id = ? and handle_flag=? ",tableDate.Dev_Id,"1").Find(&tableStoreDates4)
				Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
				t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
				Tablebehavior2.City=tableStoreDates4[0].Dev_city
				Tablebehavior2.Duration_time=SheBei[tableDate.Dev_Id]
				Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
			}else{
				for i:=0;i< len(Tablebehavior1);i++{
					if Tablebehavior1[i].Dev_id==tableDate.Dev_Id{
						Tablebehavior1[i].Duration_time=SheBei[tableDate.Dev_Id]
					}
				}
			}

			//if id!=tableDate.Dev_Id{
			//	var Tablebehavior2 model.Tablebehavior
			//	Tablebehavior2.Dev_id=id
			//	t.DB.Table("id1").Where("dev_id = ? ",id).Find(&tableStoreDates4)
			//	Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
			//	t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
			//	Tablebehavior2.City=tableStoreDates4[0].Dev_city
			//	SheBei[id][0]=SheBei[id][0]/SheBei[id][4]
			//	SheBei[id][1]=SheBei[id][1]/SheBei[id][4]
			//	SheBei[id][2]=SheBei[id][2]/SheBei[id][4]
			//	SheBei[id][3]=SheBei[id][3]/SheBei[id][4]
			//	Tablebehavior2.Duration_time=SheBei[id]
			//	Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
			//	id=tableDate.Dev_Id
			//	SheBei[tableDate.Dev_Id]=[]float32{0.0,0.0,0.0,0.0,0.0}
			//	i++
			//}
			//SheBei[tableDate.Dev_Id][0]=SheBei[tableDate.Dev_Id][0]+tableDate.Sec0p
			//SheBei[tableDate.Dev_Id][1]=SheBei[tableDate.Dev_Id][1]+tableDate.Sec30p
			//SheBei[tableDate.Dev_Id][2]=SheBei[tableDate.Dev_Id][2]+tableDate.Min3p
			//SheBei[tableDate.Dev_Id][3]=SheBei[tableDate.Dev_Id][3]+tableDate.Min10p
			//if tableDate.Sec0p==0&&tableDate.Sec30p==0&&tableDate.Min3p==0&&tableDate.Min10p==0{
			//}else{
			//	SheBei[tableDate.Dev_Id][4]=SheBei[tableDate.Dev_Id][4]+1
			//}
			//fmt.Print(Tablebehavior1)
		}
		fmt.Print(Tablebehavior1)
		//for key,content:= range SheBei{
		//	fmt.Print(key,content)
		//	var Tablebehavior2 model.Tablebehavior
		//	Tablebehavior2.Dev_id=key
		//	t.DB.Table("id1").Where("dev_id = ? ",key).Find(&tableStoreDates4)
		//	Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
		//	t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
		//	Tablebehavior2.City=tableStoreDates4[0].Dev_city
		//	Tablebehavior2.Duration_time=content
		//	Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
		//}
		//var Tablebehavior2 model.Tablebehavior
		//Tablebehavior2.Dev_id=id
		//t.DB.Table("id1").Where("dev_id = ? ",id).Find(&tableStoreDates4)
		//Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
		//t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
		//Tablebehavior2.City=tableStoreDates4[0].Dev_city
		//SheBei[id][0]=SheBei[id][0]/SheBei[id][4]
		//SheBei[id][1]=SheBei[id][1]/SheBei[id][4]
		//SheBei[id][2]=SheBei[id][2]/SheBei[id][4]
		//SheBei[id][3]=SheBei[id][3]/SheBei[id][4]
		//Tablebehavior2.Duration_time=SheBei[id]
		//Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
		//fmt.Print(Tablebehavior1)
		response.Success(ctx, gin.H{"SheBei": Tablebehavior1}, "成功")
	}
}

func NewTableStoreController ()ITaleStoreDateController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return TableStoreDateController{DB:db}
}
func time2sec(time string) (second int) {
	//var timeArr = strings.Split(time,":")
	//var hour = timeArr[0]
	//var minute = timeArr[1]
	//var sec = timeArr[2]
	//hour1, _ := strconv.Atoi(hour)
	//minute1,_:=strconv.Atoi(minute)
	//sec1 ,_ :=strconv.Atoi(sec)
	//return hour1 * 3600 + minute1 * 60 + sec1
	hour:="";
	min:="";
	sec:="";
	j:=0;
	k:=0;
	// hour=time.split('h')[0];
	// min=time.split('h')[1].split('m')[0];
	// sec=time.split('h')[1].split('m')[1].split('s')[0];
	for i:=0;i<len(time);i++{
		if time[i:i+1]=="h"{
			hour=time[0:i]
			j=i+1;
		}
		if time[i:i+1]=="m"{
			if j!=0{
				min=time[j:i]
			}else{
				min=time[0:i]
			}
			k=i+1;
		}
		if time[i:i+1]=="s"{
			if k!=0{
				sec=time[k:i]
			}else{
				sec=time[0:i]
			}
			k=i+1;
		}
	}
	hour1,_:=strconv.Atoi(hour)
	min1,_:=strconv.Atoi(min)
	sec1,_:=strconv.Atoi(sec)
	second=hour1*3600+min1*60+sec1;
	// console.log(sec)
	return
}



//指标明细计算//个数,总时长,最大时间,最短时间,最大阶跃幅值,时间明细,水流量极差明细,水流量幅值明细,水流量均值明细,水流量标准差明细
func feature(u []model.TableDate, model int, gapTime time.Duration, zoneTime time.Duration) (number int, sum time.Duration, max time.Duration, min time.Duration, maxChange float64, timeLen []time.Duration, flowExtreme []int, flowMaxChange []float64, flowAvg []float64, flowDeviation []float64) {
	var length int
	var addT time.Duration
	var datatime []string //时间轴
	for _, u := range u {
		datatime = append(datatime, u.Datatime) //数据放入切片中
	}
	for i, n := 0, len(u)-1; i < n; i++ {
		length++
		t2, _ := time.Parse("2006-01-02 15:04:05", datatime[i+1])
		t1, _ := time.Parse("2006-01-02 15:04:05", datatime[i])
		deltaT := t2.Sub(t1)
		if deltaT > gapTime {
			if addT < zoneTime {
				//过滤无效区间
			} else {
				//fmt.Printf("有效区间--长度:%v,结束id:%v,节点时间差值:%v，区间总时间:%v\n", length, i+1, deltaT, addT)
				timeLen = append(timeLen, featureTime(u[i-length+1:i+1], model)...) //有效区间进入时间特征统计函数,左闭右开;将结果存放进一个切片u1中
				extreme, maxChange, avg, deviation := featureFlow(u[i-length+1:i+1], model)
				flowExtreme = append(flowExtreme, extreme...)
				flowMaxChange = append(flowMaxChange, maxChange...)
				flowAvg = append(flowAvg, avg...)
				flowDeviation = append(flowDeviation, deviation...)
			}
			length = 0
			addT = 0
		} else if i == n-1 { //最后一个区间
			addT += deltaT //加上最后一个时间差
			if addT < zoneTime {
				//过滤无效区间
			} else {
				//fmt.Printf("最后有效区间--长度:%v,结束id:%v,节点时间差值:%v，区间总时间:%v\n", length+1, i+2, deltaT, addT)
				timeLen = append(timeLen, featureTime(u[i-length+1:i+2], model)...) //有效区间进入时间特征统计函数,左闭右开;将结果存放进一个切片u1中
				extreme, maxChange, avg, deviation := featureFlow(u[i-length+1:i+2], model)
				flowExtreme = append(flowExtreme, extreme...)
				flowMaxChange = append(flowMaxChange, maxChange...)
				flowAvg = append(flowAvg, avg...)
				flowDeviation = append(flowDeviation, deviation...)
			}
		} else {
			addT += deltaT
		}
	}
	if len(timeLen) != 0 { //防止空数组（无该模式的情况）
		number, sum, max, min = featureTimeCalculate(timeLen)
		maxChange = featureFlowCalculate(flowMaxChange)
		//fmt.Println(FlowMaxChange)
	}
	return
}

// featureTime 时间特征明细
func featureTime(u []model.TableDate, model int) []time.Duration {
	var number = 0
	var timeLen []time.Duration
	for i, n := 0, len(u); i < n; i++ {
		if u[i].Water_pattern == model {
			number++
			if i == n-1 { //区间末尾模式片段
				t2, _ := time.Parse("2006-01-02 15:04:05", u[i].Datatime)
				t1, _ := time.Parse("2006-01-02 15:04:05", u[i-number+1].Datatime)
				timeLen = append(timeLen, t2.Sub(t1))
			} else {
			} //不做处理
		} else {
			if number != 0 {
				t2, _ := time.Parse("2006-01-02 15:04:05", u[i-1].Datatime)
				t1, _ := time.Parse("2006-01-02 15:04:05", u[i-number].Datatime)
				timeLen = append(timeLen, t2.Sub(t1))
			}
			number = 0
		}
	}
	return timeLen
}

// featureTimeCalculate 时间特征指标计算
func featureTimeCalculate(timeLen []time.Duration) (number int, sum time.Duration, max time.Duration, min time.Duration) {
	max = timeLen[0]
	min = timeLen[0]
	for i, n := 0, len(timeLen); i < n; i++ {
		sum += timeLen[i]
		if timeLen[i] > max {
			max = timeLen[i]
		} else if timeLen[i] < min {
			min = timeLen[i]
		}
	}
	number = len(timeLen)
	return
}

// featureFlow 水流量特征明细（模式片段极差、最大幅值（有符号）、标准差、平均值）
func featureFlow(u []model.TableDate, model int) (flowExtreme []int, flowMaxChange []float64, flowAvg []float64, flowDeviation []float64) {
	var number = 0
	var deltaflow []int
	for i, n := 0, len(u)-1; i < n; i++ {
		deltaflow = append(deltaflow, u[i+1].Flow-u[i].Flow)
	}
	for i, n := 0, len(u); i < n; i++ {
		if u[i].Water_pattern == model {
			number++
			if i == n-1 { //区间末尾模式片段
				max := u[i-number+1].Flow
				min := u[i-number+1].Flow
				var sumFlow = 0
				var sumFF = 0.0
				maxChange := math.Abs(float64(deltaflow[i-number+1]))
				for i1, n1 := i-number+1, i+1; i1 < n1; i1++ {
					sumFlow += u[i1].Flow
					if u[i1].Flow > max {
						max = u[i1].Flow
					} else if u[i1].Flow < min {
						min = u[i1].Flow
					}
				}
				for i1, n1 := i-number+1, i+1; i1 < n1; i1++ {
					sumFF += (float64(sumFlow)/float64(number) - float64(u[i].Flow)) * (float64(sumFlow)/float64(number) - float64(u[i].Flow))
				}
				flowExtreme = append(flowExtreme, max-min)
				avg, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sumFlow)/float64(number)), 64)
				flowAvg = append(flowAvg, avg)
				deviation, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Sqrt(sumFF/float64(number))), 64)
				flowDeviation = append(flowDeviation, deviation)
				for i1, n1 := i-number+1, i; i1 < n1; i1++ { //水流量差值
					if math.Abs(float64(deltaflow[i1])) > math.Abs(maxChange) {
						maxChange = float64(deltaflow[i1])
					}
				}
				flowMaxChange = append(flowMaxChange, maxChange)
			} else {
			}
		} else {
			if number != 0 {
				max := u[i-number].Flow
				min := u[i-number].Flow
				var sumFlow int
				var sumFF float64
				maxChange := math.Abs(float64(deltaflow[i-number]))
				for i1, n1 := i-number, i; i1 < n1; i1++ { //水流量
					sumFlow += u[i1].Flow
					if u[i1].Flow > max {
						max = u[i1].Flow
					} else if u[i1].Flow < min {
						min = u[i1].Flow
					}
				}
				for i1, n1 := i-number, i; i1 < n1; i1++ {
					sumFF += (float64(sumFlow)/float64(number) - float64(u[i].Flow)) * (float64(sumFlow)/float64(number) - float64(u[i].Flow))
				}
				flowExtreme = append(flowExtreme, max-min)
				avg, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sumFlow)/float64(number)), 64)
				flowAvg = append(flowAvg, avg)
				deviation, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Sqrt(sumFF/float64(number))), 64)
				flowDeviation = append(flowDeviation, deviation)

				for i1, n1 := i-number, i-1; i1 < n1; i1++ { //水流量差值
					if math.Abs(float64(deltaflow[i1])) > math.Abs(maxChange) {
						maxChange = float64(deltaflow[i1])
					}
				}
				flowMaxChange = append(flowMaxChange, maxChange)
			}
			number = 0
		}
	}

	return
}

// featureFlowCalculate 求水流量幅值最大模（有符号）
func featureFlowCalculate(u []float64) float64 {

	var max = math.Abs(u[0])
	for i, n := 0, len(u); i < n; i++ { //水流量差值
		if math.Abs(u[i]) >= math.Abs(max) {
			max = u[i]
		}
	}
	return max

}
