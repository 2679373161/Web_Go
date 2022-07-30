package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	ini "gopkg.in/ini.v1"
)

type IIndextempController interface {
	RestController
	Datasave(ctx *gin.Context)
	Menu(ctx *gin.Context)
	Tempequipment(ctx *gin.Context)
	Temporder(ctx *gin.Context)
	Parameter_change(ctx *gin.Context)
	Equipmentsearch(ctx *gin.Context)
	Provincecodesearch(ctx *gin.Context)
}
type IndextempController struct {
	DB *gorm.DB
}

func (t IndextempController) Create(ctx *gin.Context) {
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

func (t IndextempController) Update(ctx *gin.Context) {
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

func (t IndextempController) Show(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "读取成功")
}

func (t IndextempController) Delete(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx, gin.H{"tableStoreDate": tableStoreDate}, "删除成功")

}

func (t IndextempController) Datasave(ctx *gin.Context) {
	city := ctx.DefaultQuery("city", "0")
	timeLow := ctx.DefaultQuery("timeLow", "2020-01-01")
	timeHigh := ctx.DefaultQuery("timeHigh", "2022-01-01")
	flag := ctx.DefaultQuery("flag", "0")
	flag1 := ctx.DefaultQuery("flag1", "0")
	effect_flag := ctx.DefaultQuery("effect_flag", "0")
	fault_code:=ctx.DefaultQuery("fault_flag", "0")
	dev_id := ctx.DefaultQuery("dev_id", "0000")
	dev_type := ctx.DefaultQuery("dev_type", "0000")
    realflag:= ctx.DefaultQuery("realflag", "0")


	var tableStoreDates []model.TableDate
	var tableStoreDates4 []model.Tableselect
	var tablefragment []model.Tablefragment
	var datatime []string //时间轴
	var flow []int        //水流量
	var flame []string    //火焰反馈
	var settemp []string  //设定温度
	var intemp []string
	var outtemp []string  //输出温度
	var model1 []int
	var zone_id []string
	var effect_mark []string
	var temp_pattern []int
	var heat_temp_score []int
	var stable_temp_score []int
	var score []model.Score

	//var Message string

	//fmt.Println(scorehigh)
	if flag == "1" {
		var tablename string
		var tablename1 string
		
		var city1 = "fragment" + city
		var tablefragment []model.Tablefragment
		var aliyun_info []model.Aliyun_info
		if realflag=="0"{
			tablename1="data"+city[0:4]+"_"+timeLow[0:4] + timeLow[5:7] + timeLow[8:10]+"_Info"	
		tablename = "data" + city + "_" + timeLow[0:4] + timeLow[5:7] + timeLow[8:10]	
		}else if realflag=="1"{
	var city1 model.Tableplace
    t.DB.Table("bo").Where("dev_id=?",dev_id).Find(&city1)
     city=city1.City_code
			tablename1="e_data"+city[0:4]+"_"+timeLow[0:4] + timeLow[5:7] + timeLow[8:10]+"_Info"	
			tablename = "e_data" + dev_id
		
		}
		if realflag=="0"{
		t.DB.Table(tablename1).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_id, timeLow, timeHigh).Find(&aliyun_info)
		}else if realflag=="1"{
			common.RunDB.Table(tablename1).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_id, timeLow, timeHigh).Find(&aliyun_info)
		}
		fmt.Println(aliyun_info)
		common.RunDB.Table(tablename).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_id, timeLow, timeHigh).Find(&tableStoreDates)
		//	fmt.Println(tableStoreDates)
		if realflag=="0"{
		common.IndexDB.Table(city1).Where(" dev_id=? AND start_time>=?  AND end_time<=?", dev_id, timeLow, timeHigh).Find(&score)
}
		for _, tableDate := range tableStoreDates {

			datatime = append(datatime, tableDate.Datatime) //数据放入切片中
			flow = append(flow, tableDate.Flow)
			flame = append(flame, tableDate.Flame)
			intemp = append(intemp, tableDate.Intemp)
			settemp = append(settemp, tableDate.Settemp)
			outtemp = append(outtemp, tableDate.Outtemp)
			model1 = append(model1, tableDate.Water_pattern)
			zone_id = append(zone_id, tableDate.Zone_id)
			effect_mark = append(effect_mark, tableDate.Effect_mark)
			temp_pattern = append(temp_pattern, tableDate.Temp_pattern)

		}

		for _, tableDate := range score {
			heat_temp_score = append(heat_temp_score, tableDate.Heat_temp_score)
			stable_temp_score = append(stable_temp_score, tableDate.Stable_temp_score)
		}
		if realflag=="0"{
		common.IndexDB.Table(city1).Where("dev_id=? AND end_time =? ", dev_id, timeHigh).Find(&tablefragment)
	}

		//a:= []interface{}{
		//	aliyun_info[0].Datatime,aliyun_info[0].Actualpwm,
		
	
fmt.Println("结束了")
		response.Success(ctx, gin.H{"aliyun_info":aliyun_info,"data": tableStoreDates, "data_time": datatime, "flow": flow, "flame": flame, "in_temp":intemp,"set_temp": settemp, "out_temp": outtemp, "model": model1,
			"zone_id": zone_id, "effect_mark": effect_mark, "temp_pattern": temp_pattern, "data1": tablefragment, "heat_temp_score": heat_temp_score, "stable_temp_score": stable_temp_score}, "成功")
	
	} else if flag == "2" {

		//waterflowmodel := ctx.DefaultQuery("model", "0")
		tempmodel := ctx.DefaultQuery("model1", "0")
		watermodel := ctx.DefaultQuery("model", "0")
		averagelow := ctx.DefaultQuery("averagelow", "0")
		averagehigh := ctx.DefaultQuery("averagehigh", "100")
		timeLow = timeLow + " 00:00:00"
		timeHigh = timeHigh + " 23:59:59"
		var city1 = "fragment" + city

		//fmt.Println(averagehigh)
		var temp_model_number []string
		if tempmodel[1:] == "2" { //不恒温模式
			temp_model_number = []string{"12", "21", "22","13","23"} //查询温度模式编号
		} else if tempmodel[1:] == "12" { //加热异常
			temp_model_number = []string{"12", "22","13","23"} //查询温度模式编号
		} else if tempmodel[1:] == "21" { //恒温异常
			temp_model_number = []string{"21", "22","23"} //查询温度模式编号

		} else if tempmodel[1:] == "22" { //加热+恒温异常
			temp_model_number = []string{"22","23"} //查询温度模式编号
		} else if tempmodel[1:] == "11" { //恒温
			temp_model_number = []string{"11"} //查询温度模式编号

		}else if tempmodel[1:] == "13" { //恒温熄火
			temp_model_number = []string{"13","23"} //查询温度模式编号

		}else if tempmodel[1:] == "23" { //加热异常、恒温熄火
			temp_model_number = []string{"23"} //查询温度模式编号

		}else if tempmodel[1:] == "3" { //无效
			temp_model_number = []string{"0"} //查询温度模式编号
		}
		if tempmodel[0:1] == "0" { //不限温度模式
			if tempmodel[1:] == "0" {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh).Find(&tablefragment)

				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			} else {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, averagelow, averagehigh).Find(&tablefragment)

				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, "7", "8", averagelow, averagehigh).Find(&tablefragment)

				} else {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, watermodel, averagelow, averagehigh).Find(&tablefragment)

				}
			}
		} else if tempmodel[0:1] == "1" { //有效温度模式
			if tempmodel[1:] == "0" {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern !=? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, "0", averagelow, averagehigh).Find(&tablefragment)
				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern !=? AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, "0", "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern !=? AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, "0", watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			} else {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, averagelow, averagehigh).Find(&tablefragment)
				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			}
			//if watermodel=="1"{  //不限水流量模式
			//	common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ? " +
			//		"AND temp_pattern != ?", dev_id, timeLow, timeHigh,  averagelow, averagehigh, "0").Find(&tablefragment)
			//}else if watermodel=="5"{//不稳定水流量模式：震荡+阶跃
			//	common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ?" +
			//		" AND temp_pattern != ? AND water_pattern BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh,  averagelow, averagehigh, "0","7","8").Find(&tablefragment)
			//}else {//单独一个水流量模式
			//	common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ?" +
			//		" AND temp_pattern != ? AND water_pattern = ? ", dev_id, timeLow, timeHigh,  averagelow, averagehigh, "0",watermodel).Find(&tablefragment)
			//}

			//}else if tempmodel == "11" { //恒温模式
			//	if watermodel == "1" { //不限水流量模式
			//		common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ? "+
			//			"AND temp_pattern = ?", dev_id, timeLow, timeHigh, averagelow, averagehigh, tempmodel).Find(&tablefragment)
			//	} else if watermodel == "5" { //不稳定水流量模式：震荡+阶跃
			//		common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ? AND temp_pattern = ?"+
			//			" AND water_pattern BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh, tempmodel, "7", "8").Find(&tablefragment)
			//	} else { //单独一个水流量模式
			//		common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ?"+
			//			" AND temp_pattern = ? AND water_pattern = ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh, tempmodel, watermodel).Find(&tablefragment)
			//	}
		} else if tempmodel[0:1] == "2" { //无效模式
			if watermodel == "1" {
				common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) "+
					"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, averagelow, averagehigh).Find(&tablefragment)
			} else if watermodel == "5" {
				common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern BETWEEN  ? AND ? "+
					"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, "7", "8", averagelow, averagehigh).Find(&tablefragment)
			} else {
				common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern = ? "+
					"AND temp_score BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, temp_model_number, watermodel, averagelow, averagehigh).Find(&tablefragment)
			}
			//tempmodel="0"//温度编号改为0与数据表编号对应
			//if watermodel == "1" { //不限水流量模式
			//	common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ? "+
			//		"AND temp_pattern = ?", dev_id, timeLow, timeHigh, averagelow, averagehigh, tempmodel).Find(&tablefragment)
			//} else if watermodel == "5" { //不稳定水流量模式：震荡+阶跃
			//	common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ? AND temp_pattern = ?"+
			//		" AND water_pattern BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh, tempmodel, "7", "8").Find(&tablefragment)
			//} else { //单独一个水流量模式
			//	common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ?"+
			//		" AND temp_pattern = ? AND water_pattern = ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh, tempmodel, watermodel).Find(&tablefragment)
			//}

		}
		//else  {
		//	//var temp_model_number []string
		//	//if tempmodel == "2"{ //不恒温模式
		//	//	temp_model_number =[]string{"12","21","22"} //查询温度模式编号
		//	//}else if tempmodel=="12"{  //加热异常
		//	//	temp_model_number =[]string{"12","22"} //查询温度模式编号
		//	//}else if tempmodel=="21"{  //恒温异常
		//	//	temp_model_number =[]string{"21","22"} //查询温度模式编号
		//	//}else if tempmodel=="22"{  //加热+恒温异常
		//	//	temp_model_number =[]string{"22"} //查询温度模式编号
		//	//}
		//
		//
		//	if watermodel == "1" { //不限水流量模式
		//		common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ? "+
		//			"AND temp_pattern in (?)", dev_id, timeLow, timeHigh, averagelow, averagehigh, temp_model_number).Find(&tablefragment)
		//	} else if watermodel == "5" { //不稳定水流量模式：震荡+阶跃
		//		common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ?" +
		//			" AND temp_pattern in (?) AND water_pattern BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh, temp_model_number, "7", "8").Find(&tablefragment)
		//	} else { //单独一个水流量模式
		//		common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND temp_score BETWEEN  ? AND ?"+
		//			" AND temp_pattern in (?) AND water_pattern = ? ", dev_id, timeLow, timeHigh, averagelow, averagehigh, temp_model_number, watermodel).Find(&tablefragment)
		//	}
		//}
		fmt.Println(len(tablefragment))
		if len(tablefragment) > 1000 && flag1 == "0" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "数据量超过1000条(建议具体查询条件)"})
			return
		} else {


			var tablefragment1 []model.Tablefragment
			t.DB.Table("tablefragments").Find(&tablefragment1)
			t.DB.Table("tablefragments").Delete(&tablefragment1)
			for i := 0; i < len(tablefragment); i++ {

				t.DB.Table("tablefragments").Create(&tablefragment[i])
			}
			if effect_flag!="2"{
				if fault_code!="all"{
					t.DB.Table("tablefragments").Where("effect_flag=? and fault_code=?", effect_flag,fault_code).Find(&tablefragment)
				}else{
					t.DB.Table("tablefragments").Where("effect_flag=?", effect_flag).Find(&tablefragment)

				}
			}else {
				if fault_code!="all"{
					t.DB.Table("tablefragments").Where(" fault_code=?",fault_code).Find(&tablefragment)
				}else{
					t.DB.Table("tablefragments").Find(&tablefragment)
				}
			}
			response.Success(ctx, gin.H{"data": tablefragment}, "成功")
			t.DB.Raw("TRUNCATE TABLE tablefragments")


		}
	} else if flag == "3" {
		tempmodel := ctx.DefaultQuery("model1", "0")
		watermodel := ctx.DefaultQuery("model", "0")
		averagelow := ctx.DefaultQuery("averagelow", "0")
		averagehigh := ctx.DefaultQuery("averagehigh", "100")
		timeLow = timeLow + " 00:00:00"
		timeHigh = timeHigh + " 23:59:59"
		var city1 = "fragment" + city
		var temp_model_number []string
		if tempmodel[1:] == "2" { //不恒温模式

			temp_model_number = []string{"12", "21", "22","13","23"} //查询温度模式编号
		} else if tempmodel[1:] == "12" { //加热异常   
			temp_model_number = []string{"12", "22","13","23"} //查询温度模式编号
		} else if tempmodel[1:] == "21" { //恒温异常  
			temp_model_number = []string{"21", "22","23"} //查询温度模式编号

		} else if tempmodel[1:] == "22" { //加热+恒温异常
			temp_model_number = []string{"22","23"} //查询温度模式编号
		} else if tempmodel[1:] == "11" { //恒温
			temp_model_number = []string{"11"} //查询温度模式编号

		}else if tempmodel[1:] == "13" { //恒温熄火
			temp_model_number = []string{"13","23"} //查询温度模式编号

		}else if tempmodel[1:] == "23" { //加热异常、恒温熄火
			temp_model_number = []string{"23"} //查询温度模式编号

		}else if tempmodel[1:] == "3" { //无效
			temp_model_number = []string{"0"} //查询温度模式编号
		}

		var Dev_id []string
		t.DB.Table("bo").Where("dev_type = ? AND city_code = ? ", dev_type, city).Find(&tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			Dev_id = append(Dev_id, tableDate.Dev_Id)
		}
		if tempmodel[0:1] == "0" { //不限温度模式
			if tempmodel[1:] == "0" {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, averagelow, averagehigh).Find(&tablefragment)
				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			} else {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, averagelow, averagehigh).Find(&tablefragment)
				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			}
		} else if tempmodel[0:1] == "1" { //有效温度模式
			if tempmodel[1:] == "0" {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern !=? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, "0", averagelow, averagehigh).Find(&tablefragment)
				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern !=? AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, "0", "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern !=? AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, "0", watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			} else {
				if watermodel == "1" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, averagelow, averagehigh).Find(&tablefragment)
				} else if watermodel == "5" {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern BETWEEN  ? AND ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, "7", "8", averagelow, averagehigh).Find(&tablefragment)
				} else {
					common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern = ? "+
						"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, watermodel, averagelow, averagehigh).Find(&tablefragment)
				}
			}
		} else if tempmodel[0:1] == "2" { //无效模式
			if watermodel == "1" {
				common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) "+
					"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, averagelow, averagehigh).Find(&tablefragment)
			} else if watermodel == "5" {
				common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern BETWEEN  ? AND ? "+
					"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, "7", "8", averagelow, averagehigh).Find(&tablefragment)
			} else {
				common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND temp_pattern in(?) AND water_pattern = ? "+
					"AND temp_score BETWEEN  ? AND ? ", Dev_id, timeLow, timeHigh, temp_model_number, watermodel, averagelow, averagehigh).Find(&tablefragment)
			}
		}
		fmt.Println(len(tablefragment))
		if len(tablefragment) > 1000 && flag1 == "0" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "数据量超过1000条(建议具体查询条件)"})
			return
		} else {
			if effect_flag == "2" {
				response.Success(ctx, gin.H{"data": tablefragment}, "成功")
			} else {
				var tablefragment1 []model.Tablefragment
				common.IndexDB.Table("tablefragments").Find(&tablefragment1)
				common.IndexDB.Table("tablefragments").Delete(&tablefragment1)
				for i := 0; i < len(tablefragment); i++ {

					common.IndexDB.Table("tablefragments").Create(&tablefragment[i])
				}

				common.IndexDB.Table("tablefragments").Where("effect_flag=?", effect_flag).Find(&tablefragment)
				fmt.Println(len(tablefragment))

				response.Success(ctx, gin.H{"data": tablefragment}, "成功")

			}
		}
	}else {
	}
}
func (t IndextempController) Tempequipment(ctx *gin.Context) {
	flag := ctx.DefaultQuery("flag", "0")
	flag1 := ctx.DefaultQuery("flag1", "0")
	category := ctx.DefaultQuery("category", "0000")
	city := ctx.DefaultQuery("city", "0")
	timeLow := ctx.DefaultQuery("timeLow", "2020-01-01")
	timeHigh := ctx.DefaultQuery("timeHigh", "2022-01-01")
	province_code := ctx.DefaultQuery("province", "0000")
	id := ctx.DefaultQuery("id", "0000")


	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))

	var tableplace []model.Tableplace
	var tableStoreDates4 []model.TableDate3
	var avg_score[]model.Score_province
	var filter []string
	var province []string
	var city_code []string
	var Dev_type []string
	var num int
	time := make(map[string][]string)
	xiaoshu := make(map[string][]float32)
	zhengshu := make(map[string][]int64)
	println("category=", category)
	println("city=", city)
	println("province=", province_code)
	if flag == "4"{
		var Provincesearch struct {
			Province_code string `json:"province_code" gorm:"type:int;not null"`
			City_code     string `json:"city_code" gorm:"type:int;not null"`
		}
		common.IndexDB.Table("days_summaries").Where("dev_id= ?  ", id).Find(&Provincesearch)
		fmt.Println("我",Provincesearch)
		response.Success(ctx, gin.H{"data": Provincesearch}, "成功")
	}
	if flag1 == "1" {
		if flag == "1" {
			if category == "0000" {
				if city == "0" {
					common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? and province_code=?", timeLow, timeHigh,province_code).Find(&tableplace)
					common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && province_code=? GROUP BY dev_id ",timeLow,timeHigh,province_code).Find(&avg_score)

				} else {
					common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow, timeHigh, city).Find(&tableplace)
					common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && city_code=? GROUP BY dev_id ",timeLow,timeHigh,city).Find(&avg_score)

				}
			} else {
				common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow, timeHigh,city,category).Find(&tableplace)
				common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && dev_type=? GROUP BY dev_id ",timeLow,timeHigh,category).Find(&avg_score)

			}

		} else {
			if category == "0000" {
				if province_code == "0000" {
					common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ?  ", timeLow, timeHigh).Find(&tableplace).Offset(perPage * (currentPage - 1))
					//common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ?  ", timeLow, timeHigh).Count(&num)
				} else {
					common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND province_code=? ", timeLow, timeHigh, province_code).Find(&tableplace)
					common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && province_code=? GROUP BY dev_id ",timeLow,timeHigh,province_code).Find(&avg_score)
				}
			} else {
				common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow, timeHigh, city,category).Find(&tableplace)
				common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && city_code=? AND dev_type=? GROUP BY dev_id ",timeLow,timeHigh,city,category).Find(&avg_score)

			}

		}
	} else {
		if city == "0" {
			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND province_code=? ", timeLow, timeHigh, province_code).Find(&tableplace)
			common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && province_code=? GROUP BY dev_id ",timeLow,timeHigh,province_code).Find(&avg_score)
		}else{
			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow, timeHigh, city).Find(&tableplace)
			common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && city_code=? GROUP BY dev_id ",timeLow,timeHigh,city).Find(&avg_score)
		}
		//	common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND dev_id=? ", timeLow, timeHigh, id).Count(&num)
       if flag1=="3"{
		common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND dev_id=? ", timeLow, timeHigh, id).Find(&tableplace)
		common.IndexDB.Raw("select avg(temp_score) avg_score from days_summaries where time_date>=? && time_date<=? && dev_id=? GROUP BY dev_id ",timeLow,timeHigh,id).Find(&avg_score)

	   }
		fmt.Print(tableplace)
	}
	//fmt.Print(tableplace)
	//fmt.Println(tableplace)
	for _, tableDate := range tableplace {
		var flag bool = false
		for i := 0; i < len(filter); i++ {
			if tableDate.Dev_Id == filter[i] {
				flag = true
			}
		}
		if flag == false {
			filter = append(filter, tableDate.Dev_Id)
			Dev_type = append(Dev_type, tableDate.Dev_type)
			city_code = append(city_code, tableDate.City_code)
			t.DB.Table("midea_loc_code").Where("city_code=?", city_code[len(city_code)-1]).Find(&tableStoreDates4)
			for _, tableDate1 := range tableStoreDates4 {
				province = append(province, tableDate1.Dev_city)
			}
			time[tableDate.Dev_Id] = []string{"0", "0", "0"}
			xiaoshu[tableDate.Dev_Id] = []float32{0.0, 0.0}
			zhengshu[tableDate.Dev_Id] = []int64{0, 0, 0, 0, 0,0,0}//温度评分、片段数、天数、超调量、极差、升温段评分、恒温段评分
		}

		time0, _ := strconv.Atoi(time[tableDate.Dev_Id][0])
		time1, _ := strconv.Atoi(time[tableDate.Dev_Id][1])
		time2, _ := strconv.Atoi(time[tableDate.Dev_Id][2])
		time[tableDate.Dev_Id][0] = strconv.Itoa(time2sec(tableDate.Temp_valid_time) + time0)
		//fmt.Println(time[tableDate.Dev_Id][0])
		time[tableDate.Dev_Id][1] = strconv.Itoa(time2sec(tableDate.Ave_heat_duration) + time1)
		time[tableDate.Dev_Id][2] = strconv.Itoa(time2sec(tableDate.Ave_un_sable_duration) + time2)
		xiaoshu[tableDate.Dev_Id][0] = tableDate.Ave_unstable_proportion + xiaoshu[tableDate.Dev_Id][0]
		xiaoshu[tableDate.Dev_Id][1] = xiaoshu[tableDate.Dev_Id][1] + 1
		zhengshu[tableDate.Dev_Id][0] = int64(tableDate.Temp_score*tableDate.Temp_num) + zhengshu[tableDate.Dev_Id][0]
		zhengshu[tableDate.Dev_Id][1] = int64(tableDate.Temp_num) + zhengshu[tableDate.Dev_Id][1]
		zhengshu[tableDate.Dev_Id][2] = zhengshu[tableDate.Dev_Id][2] + 1
		zhengshu[tableDate.Dev_Id][3] = int64(tableDate.Overshoot_value) + zhengshu[tableDate.Dev_Id][3]
		zhengshu[tableDate.Dev_Id][4] = int64(tableDate.Temp_range) + zhengshu[tableDate.Dev_Id][4]
		zhengshu[tableDate.Dev_Id][5] = int64(tableDate.Heat_temp_score*tableDate.Temp_num) + zhengshu[tableDate.Dev_Id][5]
		zhengshu[tableDate.Dev_Id][6] = int64(tableDate.Stable_temp_score*tableDate.Temp_num) + zhengshu[tableDate.Dev_Id][6]
		//Dev_type = append(Dev_type, tableDate.Dev_type)                            //数据放入切片中
		//Stable_proportion = append(Stable_proportion, tableDate.Stable_proportion) //数据放入切片中
		//Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
		//city_code = append(city_code, tableDate.City_code)
		//t.DB.Table("midea_loc_code").Where("city_code=?", city_code[len(city_code)-1]).Find(&tableStoreDates4)
		//for _, tableDate1 := range tableStoreDates4 {
		//	province = append(province, tableDate1.Dev_city)
		//}
		fmt.Println("整数",zhengshu[tableDate.Dev_Id][0])
	}

	//fmt.Print(zhengshu)
	fmt.Print(time)
	//fmt.Print(zhengshu)
	for key := range time {
		if time[key][0] == "0" {
			zhengshu[key][0] = 0
			zhengshu[key][1] = 0
			zhengshu[key][3] = 0
			zhengshu[key][4] = 0
			zhengshu[key][5] = 0
			zhengshu[key][6] = 0
			xiaoshu[key][0] = 0
			time[key][0] = "0"
			time[key][1] = "0"
			time[key][2] = "0"
		} else {
			Time0, _ := strconv.Atoi(time[key][0])
			Time1, _ := strconv.Atoi(time[key][1])
			Time2, _ := strconv.Atoi(time[key][2])
			xiaoshu[key][0] = xiaoshu[key][0] / xiaoshu[key][1]
			zhengshu[key][0] = zhengshu[key][0] / zhengshu[key][1]//温度均分=温度总分/片段数
			zhengshu[key][5] = zhengshu[key][5] / zhengshu[key][1]
			zhengshu[key][6] = zhengshu[key][6] / zhengshu[key][1]
			zhengshu[key][1] = zhengshu[key][1] / zhengshu[key][2]//片段个数=片段总数/有数据的天数
			zhengshu[key][3] = zhengshu[key][3] / zhengshu[key][2]//温度均分=温度总分/片段数
			zhengshu[key][4] = zhengshu[key][4] / zhengshu[key][2]//片段个数=片段总数/有数据的天数


			time[key][0] = strconv.Itoa(Time0 / int(zhengshu[key][2]))
			time[key][1] = strconv.Itoa(Time1 / int(zhengshu[key][2]))
			time[key][2] = strconv.Itoa(Time2 / int(zhengshu[key][2]))
		}

	}
	fmt.Print(filter)

	fmt.Println("zhengshu=",zhengshu)
	fmt.Println("666")
	//common.DB.Table("equipment_search").Delete(&tableplace)
	//	for i := 0; i < len(tableplace); i++ {
	//		common.DB.Table("equipment_search").Create(&tableplace[i])
	//	}
	//common.DB.Table("equipment_search").Where("temp_score BETWEEN  ? AND ? ", 90, 100).Find(&tableplace)
	//	common.DB.Table("equipment_search").Where("temp_score BETWEEN  ? AND ? ", 90, 100).Count(&num)
	//	fmt.Println("666")
	//fmt.Println(Dev_type)
	//fmt.Println(filter)
	//fmt.Println("avg_score",avg_score)
	response.Success(ctx, gin.H{"data": tableplace, "count":num,"province": province, "dev_type": Dev_type, "filter": filter, "time": time, "xiaoshu": xiaoshu, "zhengshu": zhengshu,"avg_score":avg_score}, "成功")

}
func (t IndextempController) Temporder(ctx *gin.Context) {
	var flag string = ctx.DefaultQuery("flag", "0")
	var flag1 string = ctx.DefaultQuery("flag1", "0")
	var month1 string = ctx.DefaultQuery("month", "1") + "-01"
	var day1 string = ctx.DefaultQuery("day", "1")
	var feature model.Datafeature
	var tableplace2 []model.Tableplace
	var tableStoreDates4 model.TableDate3
	var city []string
	if flag == "1" {
		common.IndexDB.Table("data_features").Order("update_time desc").Last(&feature)
		//if flag1=="1"{
		month1 = feature.Update_time[0:7] + "-01"
		//}else{
		day1 = feature.Update_time
		//}
	}
	if flag1 == "1" {
		common.IndexDB.Table("month_summaries").Where("time_date = ? AND temp_score !=?", month1, "0").Order("temp_score asc").Limit(100).Find(&tableplace2)
	} else {
		common.IndexDB.Table("days_summaries").Where("time_date = ? AND temp_score !=?", day1, "0").Order("temp_score asc").Limit(100).Find(&tableplace2)
	}

	for i, tableDate := range tableplace2 {
		t.DB.Table("midea_loc_code").Where("city_code=?", tableDate.City_code).Find(&tableStoreDates4)
		city = append(city, tableStoreDates4.Dev_city)
		if flag1 == "1" {
			tableplace2[i].Time_date = tableDate.Time_date[0:7]
		}
	}
	response.Success(ctx, gin.H{"month1": month1[0:7], "day1": day1, "data": tableplace2, "city": city}, "成功")
}
func (t IndextempController) Provincecodesearch(ctx *gin.Context) {

	id := ctx.DefaultQuery("id", "0000")
	var Provincesearch struct {
		City_code     string `json:"city_code" gorm:"type:int;not null"`
		Dev_type          string `json:"dev_type" gorm:"type:int;not null"`
	}
	common.IndexDB.Table("days_summaries").Where("dev_id= ?  ", id).Find(&Provincesearch)
	fmt.Println("Provincesearch",Provincesearch)
	response.Success(ctx, gin.H{"data": Provincesearch}, "成功")
}

func (t IndextempController) Menu(ctx *gin.Context) {
	var tableplace []model.Tableplace
	var tableplace1 []model.Tableplace
	var tableplace2 []model.Tableplace
	var tableplace3 []model.Tableplace
	var tableplace4 []model.Tableplace
	var tableplace5 []model.Tableplace
	var tableplace6 []model.Tableplace
	var Behavior []model.Behavior
	var tableStoreDates4 model.TableDate3
	var city []string
	var province1 []string
	var months []string
	var months1 []string
	var tempscore5 []float32
	var waterscore5 []float32
	nationaltrendinmonth := make(map[string][]int)
	tempregion := make(map[string][]int)
	waterregion := make(map[string][]int)
	monthregion := make(map[string][]string)
	var region []string
	var tempscore4 []int
	var waterscore4 []int
	var dev_type []string
	var dev_id []string
	var tempscore []int
	var waterscore []int
	var tempscore1 []int
	var waterscore1 []int
	var tempscore2 []int
	var waterscore2 []int
	var province []string
	var equipment_num []int
	var tempscore3 []int
	var waterscore3 []int
	var behaviorTime = []float32{0.0, 0.0, 0.0, 0.0, 0.0}
	var Regiontempwater []model.Regiontempwater
	var feature model.Datafeature
	var feature1 model.Datafeature
	var feature2 []model.Datafeature
	var feature3 []model.Datafeature
	common.IndexDB.Table("data_features").Order("update_time desc").Last(&feature)
	common.IndexDB.Table("data_features").Order("update_time asc").First(&feature1)
	mon, _ := strconv.Atoi(feature.Update_time[5:7])
	var mon1 string
	year := feature.Update_time[0:4]
	if mon == 1 {
		temp_year,_ := strconv.Atoi(feature.Update_time[0:4])
		temp_year1:=temp_year - 1
		year=strconv.Itoa(temp_year1)
		mon1="12"

	} else {
		mon1 = strconv.Itoa(mon - 1)
		fmt.Print(mon1)
		if len(mon1) == 1 {
			mon1 = "0" + mon1
		}
	}
	time := year + "-" + mon1 + "-01"
	fmt.Print(feature.Free_space)
	//time:=ctx.DefaultQuery("time","2020-01-01")
	common.IndexDB.Table("city_months").Where("time_date = ? AND temp_score !=?", time, "0").Order("temp_score asc").Limit(10).Find(&tableplace)
	common.IndexDB.Table("type_months").Where("time_date = ? AND temp_score !=?", time, "0").Order("temp_score asc").Limit(10).Find(&tableplace1)
	common.IndexDB.Table("province_months").Where("time_date = ? AND temp_score !=?", time, "0").Order("temp_score asc").Limit(10).Find(&tableplace4)
	common.IndexDB.Table("month_summaries").Where("time_date = ? AND temp_score !=?", time, "0").Order("temp_score asc").Limit(5).Find(&tableplace2)
	common.IndexDB.Table("province_months").Where("time_date = ?", time).Find(&tableplace3)
	for _, tableDate := range tableplace4 {
		t.DB.Table("midea_loc_code").Where("province_code=?", tableDate.Province_code).Limit(1).Find(&tableStoreDates4)
		province1 = append(province1, tableStoreDates4.Dev_province)
		tempscore4 = append(tempscore4, tableDate.Temp_score)
		waterscore4 = append(waterscore4, tableDate.Water_score)
	}
	for _, tableDate := range tableplace {
		t.DB.Table("midea_loc_code").Where("city_code=?", tableDate.City_code).Find(&tableStoreDates4)
		city = append(city, tableStoreDates4.Dev_city)
		tempscore = append(tempscore, tableDate.Temp_score)
		waterscore = append(waterscore, tableDate.Water_score)
	}
	for _, tableDate := range tableplace1 {
		dev_type = append(dev_type, tableDate.Dev_type)
		tempscore1 = append(tempscore1, tableDate.Temp_score)
		waterscore1 = append(waterscore1, tableDate.Water_score)
	}
	for _, tableDate := range tableplace2 {
		dev_id = append(dev_id, tableDate.Dev_Id)
		tempscore2 = append(tempscore2, tableDate.Temp_score)
		waterscore2 = append(waterscore2, tableDate.Water_score)
	}
	for _, tableDate := range tableplace3 {
		t.DB.Table("midea_loc_code").Where("province_code=?", tableDate.Province_code).Find(&tableStoreDates4)
		province = append(province, tableStoreDates4.Dev_province)
		equipment_num = append(equipment_num, tableDate.Equipment_num)
		tempscore3 = append(tempscore3, tableDate.Temp_score)
		waterscore3 = append(waterscore3, tableDate.Water_score)
	}
	//fmt.Print(province1)

	common.IndexDB.Table("province_months").Where("time_date <= ? ", time).Find(&tableplace4)

	for _, tableDate := range tableplace4 {
		var flag1 bool = false
		for i := 0; i < len(months); i++ {
			if tableDate.Time_date == months[i] {
				flag1 = true

			}
		}
		if flag1 == false {
			months = append(months, tableDate.Time_date)
			nationaltrendinmonth[tableDate.Time_date] = []int{0, 0, 0}

		}
		nationaltrendinmonth[tableDate.Time_date][0]++
		nationaltrendinmonth[tableDate.Time_date][1] = nationaltrendinmonth[tableDate.Time_date][1] + tableDate.Temp_score
		nationaltrendinmonth[tableDate.Time_date][2] = nationaltrendinmonth[tableDate.Time_date][2] + tableDate.Water_score

	}
	//fmt.Println(nationaltrendinmonth)

	common.IndexDB.Table("region_months").Where("time_date <= ? ", time).Find(&tableplace5)
	fmt.Println(tableplace5)
	fmt.Println("12345678")
	for _, tableDate := range tableplace5 {
		t.DB.Table("midea_loc_code").Where("region_code=?", tableDate.Region_code).Find(&tableStoreDates4)
		var flag1 bool = false

		for i := 0; i < len(region); i++ {
			if tableStoreDates4.Dev_region == region[i] {
				flag1 = true
			}
		}
		if flag1 == false {
			region = append(region, tableStoreDates4.Dev_region)
			tempregion[tableStoreDates4.Dev_region] = []int{}
			waterregion[tableStoreDates4.Dev_region] = []int{}
			monthregion[tableStoreDates4.Dev_region] = []string{}
		}
		tempregion[tableStoreDates4.Dev_region] = append(tempregion[tableStoreDates4.Dev_region], tableDate.Temp_score)
		waterregion[tableStoreDates4.Dev_region] = append(waterregion[tableStoreDates4.Dev_region], tableDate.Water_score)
		monthregion[tableStoreDates4.Dev_region] = append(monthregion[tableStoreDates4.Dev_region], tableDate.Time_date[0:7])
	}
	fmt.Println(region)
	for i := 0; i < len(region); i++ {
		var regiontempwater model.Regiontempwater
		regiontempwater.Region = region[i]
		regiontempwater.Tempscore = tempregion[region[i]]
		regiontempwater.Waterscore = waterregion[region[i]]
		regiontempwater.Month = monthregion[region[i]]
		Regiontempwater = append(Regiontempwater, regiontempwater)
	}
	//fmt.Println(Regiontempwater)

	for _, month := range months {
		//nationaltrendinmonth[month][1]=nationaltrendinmonth[month][1]/nationaltrendinmonth[month][0]
		// nationaltrendinmonth[month][2]=nationaltrendinmonth[month][2]/nationaltrendinmonth[month][0]
		tempscore5 = append(tempscore5, float32(nationaltrendinmonth[month][1])/float32(nationaltrendinmonth[month][0]))
		waterscore5 = append(waterscore5, float32(nationaltrendinmonth[month][2])/float32(nationaltrendinmonth[month][0]))
		months1 = append(months1, month[0:7])

	}

	common.IndexDB.Table("type_months").Where("time_date <= ? AND dev_type = ?", time, tableplace1[0].Dev_type).Find(&tableplace6)
	//fmt.Println(tableplace6)
	common.IndexDB.Table("behavior_summaries").Where("data_time BETWEEN  ? AND ? AND dev_id = ?", time, time[0:7]+"-31", tableplace2[0].Dev_Id).Find(&Behavior)

	for _, tableDate := range Behavior {
		behaviorTime[0] = tableDate.Sec0p + behaviorTime[0]
		behaviorTime[1] = tableDate.Sec30p + behaviorTime[1]
		behaviorTime[2] = tableDate.Min3p + behaviorTime[2]
		behaviorTime[3] = tableDate.Min10p + behaviorTime[3]
		behaviorTime[4] = behaviorTime[4] + 1
	}
	//behaviorTime[0]=behaviorTime[0]/behaviorTime[4]
	//behaviorTime[1]=behaviorTime[1]/behaviorTime[4]
	//behaviorTime[2]=behaviorTime[2]/behaviorTime[4]
	//behaviorTime[3]=behaviorTime[3]/behaviorTime[4]
	//fmt.Println(tableplace2[0].Dev_Id)
	//fmt.Println(behaviorTime)
	common.IndexDB.Table("data_features").Order("update_time desc").Limit(31).Find(&feature2)
	fmt.Println(feature2)
	fmt.Println(feature3)
	response.Success(ctx, gin.H{"data": feature, "data1": feature1, "data2": Regiontempwater, "data4": feature2, "region": region, "data3": tableplace6, "city": city, "tempscore": tempscore, "waterscore": waterscore, "dev_id": dev_id, "tempscore2": tempscore2, "waterscore2": waterscore2, "dev_type": dev_type,
		"tempscore1": tempscore1, "waterscore1": waterscore1, "province": province, "equipment_num": equipment_num, "tempscore3": tempscore3, "waterscore3": waterscore3, "province1": province1,
		"tempscore4": tempscore4, "waterscore4": waterscore4, "months": months1, "tempscore5": tempscore5, "waterscore5": waterscore5, "behaviorTime": behaviorTime}, "成功")
}
func NewIndextempController() IIndextempController {
	db := common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return IndextempController{DB: db}
}
func (t IndextempController) Parameter_change(ctx *gin.Context) {

	OriginalDB := ctx.DefaultQuery("OriginalDB", "")
	RunDB := ctx.DefaultQuery("RunDB", "")
	IndexDB := ctx.DefaultQuery("IndexDB", "")
	Username := ctx.DefaultQuery("Username", "")
	Password := ctx.DefaultQuery("Password", "")
	Equipments := ctx.DefaultQuery("Equipments", "")
	println("OriginalDB:", OriginalDB)

	FragmentInterval := ctx.DefaultQuery("FragmentInterval", "")
	InvalidTime := ctx.DefaultQuery("InvalidTime", "")
	BehaviorInterval := ctx.DefaultQuery("BehaviorInterval", "")
	FirstInterpolationLower := ctx.DefaultQuery("FirstInterpolationLower", "")
	FirstInterpolationUpper := ctx.DefaultQuery("FirstInterpolationUpper", "")
	SecondInterpolationLower := ctx.DefaultQuery("SecondInterpolationLower", "")
	SecondInterpolationUpper := ctx.DefaultQuery("SecondInterpolationUpper", "")
	WaterValidTime := ctx.DefaultQuery("WaterValidTime", "")
	TempValidTime := ctx.DefaultQuery("TempValidTime", "")
	EffectiveBehaviorTime := ctx.DefaultQuery("EffectiveBehaviorTime", "")
	HeatStableThreshold := ctx.DefaultQuery("HeatStableThreshold", "")
	ConstantTempThreshold := ctx.DefaultQuery("ConstantTempThreshold", "")
	OvershootThreshold := ctx.DefaultQuery("OvershootThreshold", "")
	FlowChangeThreshold := ctx.DefaultQuery("FlowChangeThreshold", "")
	BeforeCollectionPoint := ctx.DefaultQuery("BeforeCollectionPoint", "")
	AfterCollectionPoint := ctx.DefaultQuery("AfterCollectionPoint", "")
	DisturbanceDuration := ctx.DefaultQuery("DisturbanceDuration", "")
	DisturbanceProportion := ctx.DefaultQuery("DisturbanceProportion", "")
	NonConstantTempDuration := ctx.DefaultQuery("NonConstantTempDuration", "")
	HeatingDuration := ctx.DefaultQuery("HeatingDuration", "")
	AutoFlag := ctx.DefaultQuery("AutoFlag", "")
	DataWriteFlag := ctx.DefaultQuery("DataWriteFlag", "")
	FragmentWriteFlag := ctx.DefaultQuery("FragmentWriteFlag", "")
	Thread := ctx.DefaultQuery("Thread", "")
	StartTime := ctx.DefaultQuery("datatime_s_e[0]", "2022-01-01")
	EndTime := ctx.DefaultQuery("datatime_s_e[1]", "2022-01-01")
	CycleTime := ctx.DefaultQuery("CycleTime", "")
	Timing := ctx.DefaultQuery("Timing", "")
	TimingPeriod := ctx.DefaultQuery("TimingPeriod", "")
	BasicDataDelete := ctx.DefaultQuery("BasicDataDelete", "0")
	OperationDataSaveTimed := ctx.DefaultQuery("OperationDataSaveTimed", "0")

	cfg, err := ini.Load("E:/study/go/src/ginEssential/最新/config.ini")
	if err != nil {
		log.Fatal("Fail to read file:", err)
	}
	cfg.Section("Mysql").Key("OriginalDB").SetValue(OriginalDB)
	cfg.Section("Mysql").Key("RunDB").SetValue(RunDB)
	cfg.Section("Mysql").Key("IndexDB").SetValue(IndexDB)
	cfg.Section("Mysql").Key("Username").SetValue(Username)
	cfg.Section("Mysql").Key("Password").SetValue(Password)
	cfg.Section("Mysql").Key("Equipments").SetValue(Equipments)

	if AutoFlag != "" {
		cfg.Section("Task").Key("AutoFlag").SetValue(AutoFlag)
	}
	if DataWriteFlag != "" {
		cfg.Section("Task").Key("DataWriteFlag").SetValue(DataWriteFlag)
	}
	if FragmentWriteFlag != "" {
		cfg.Section("Task").Key("FragmentWriteFlag").SetValue(FragmentWriteFlag)
	}
	cfg.Section("Task").Key("Thread").SetValue(Thread)
	cfg.Section("Task").Key("StartTime").SetValue(StartTime)
	cfg.Section("Task").Key("EndTime").SetValue(EndTime)
	cfg.Section("Task").Key("CycleTime").SetValue(CycleTime)
	cfg.Section("Task").Key("Timing").SetValue(Timing)
	cfg.Section("Task").Key("TimingPeriod").SetValue(TimingPeriod)
	cfg.Section("Task").Key("BasicDataDelete").SetValue(BasicDataDelete)
	cfg.Section("Task").Key("OperationDataSaveTimed").SetValue(OperationDataSaveTimed)

	cfg.Section("Clean").Key("FragmentInterval").SetValue(FragmentInterval)
	cfg.Section("Clean").Key("InvalidTime").SetValue(InvalidTime)
	cfg.Section("Clean").Key("BehaviorInterval").SetValue(BehaviorInterval)
	cfg.Section("Clean").Key("FirstInterpolationLower").SetValue(FirstInterpolationLower)
	cfg.Section("Clean").Key("FirstInterpolationUpper").SetValue(FirstInterpolationUpper)
	cfg.Section("Clean").Key("SecondInterpolationLower").SetValue(SecondInterpolationLower)
	cfg.Section("Clean").Key("SecondInterpolationUpper").SetValue(SecondInterpolationUpper)
	cfg.Section("Clean").Key("WaterValidTime").SetValue(WaterValidTime)
	cfg.Section("Clean").Key("TempValidTime").SetValue(TempValidTime)
	cfg.Section("Clean").Key("EffectiveBehaviorTime").SetValue(EffectiveBehaviorTime)

	cfg.Section("Mining").Key("HeatStableThreshold").SetValue(HeatStableThreshold)
	cfg.Section("Mining").Key("ConstantTempThreshold").SetValue(ConstantTempThreshold)
	cfg.Section("Mining").Key("OvershootThreshold").SetValue(OvershootThreshold)
	cfg.Section("Mining").Key("FlowChangeThreshold").SetValue(FlowChangeThreshold)
	cfg.Section("Mining").Key("BeforeCollectionPoint").SetValue(BeforeCollectionPoint)
	cfg.Section("Mining").Key("AfterCollectionPoint").SetValue(AfterCollectionPoint)

	cfg.Section("Temperature").Key("DisturbanceDuration").SetValue(DisturbanceDuration)
	cfg.Section("Temperature").Key("DisturbanceProportion").SetValue(DisturbanceProportion)
	cfg.Section("Temperature").Key("NonConstantTempDuration").SetValue(NonConstantTempDuration)
	cfg.Section("Temperature").Key("HeatingDuration").SetValue(HeatingDuration)
	cfg.SaveTo("E:/study/go/src/ginEssential/最新/config.ini")

}
func (t IndextempController) Equipmentsearch(ctx *gin.Context) {
	category := ctx.DefaultQuery("category", "")
	city := ctx.DefaultQuery("city", "")
	timeLow := ctx.DefaultQuery("timeLow", "2020-01-01")
	timeHigh := ctx.DefaultQuery("timeHigh", "2022-01-01")
	var equipmentplace []model.Equipment_search
	var effectiveplace []model.Effective_statistics
	var equipmentplace1 []model.Equipment_search
	var tableStoreDates3 []model.TableDate3
	var num = 1
	if len(city) == 0 && len(category) == 0 {
		common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ?  ", timeLow, timeHigh).Find(&equipmentplace)
		fmt.Println("tableplace1=", equipmentplace)
		common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? ", timeLow, timeHigh).Count(&num)
	}

	if len(city) != 0 {
		if len(category) == 0 {
			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow, timeHigh, city).Find(&equipmentplace)
			//fmt.Println("tableplace1=",equipmentplace)
			//common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow, timeHigh, city).Count(&num)
		} else if len(category) != 0 {
			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow, timeHigh, city, category).Find(&equipmentplace)
			//fmt.Println("tableplace2=",equipmentplace)
			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow, timeHigh, city, category).Count(&num)
		}
	}
	if len(city) == 0 && len(category) != 0 {

		common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND dev_type=? ", timeLow, timeHigh, category).Find(&equipmentplace)
	}

	common.IndexDB.Table("equipmentsearch").Delete(&equipmentplace)
	for i := 0; i < len(equipmentplace); i++ {
		common.IndexDB.Table("equipmentsearch").Create(&equipmentplace[i])
	}
	common.IndexDB.Raw("select city_code,dev_id,dev_type,sum(temp_score) as temp_score,FLOOR(avg(un_stable_proportion))as un_stable_proportion,FLOOR(avg(temp_valid_time)) as temp_valid_time,FLOOR(avg(ave_heat_duration))as ave_heat_duration,FLOOR(avg(ave_un_sable_duration)) as ave_un_sable_duration,FLOOR(avg(temp_pattern)) as temp_pattern from equipmentsearch  group by dev_id ").Find(&equipmentplace1)

	common.IndexDB.Raw("select distinct dev_id,effective_day from equipmentsearch ").Find(&effectiveplace)

	for j := 0; j < len(effectiveplace); j++ {
		for i := 0; i < len(equipmentplace); i++ {
			if equipmentplace[i].Dev_Id == effectiveplace[j].Dev_id && equipmentplace[i].Temp_score != 0 {
				effectiveplace[j].Effective_day++
			}
		}
	}
	for i := 0; i < len(effectiveplace); i++ {

		if effectiveplace[i].Effective_day != 0 {
			equipmentplace1[i].Temp_score = equipmentplace1[i].Temp_score / effectiveplace[i].Effective_day
		}
		if effectiveplace[i].Effective_day == 0 {
			equipmentplace1[i].Temp_score = 0
		}
	}
	fmt.Println("zuizhong=", equipmentplace)
	// var tableplace [] model.Equipment_search
	for _, tableDate := range equipmentplace1 {
		var tableStoreDates4 model.TableDate3
		t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
		tableStoreDates3 = append(tableStoreDates3, tableStoreDates4)
		//fmt.Println("666")
		//fmt.Println("城市=",tableStoreDates7)
	}
	//fmt.Println("666")
	//fmt.Println("城市=",tableStoreDates3)
	fmt.Println("equipmentplace1=", equipmentplace1)

	response.Success(ctx, gin.H{"data": equipmentplace1, "tabledates3": tableStoreDates3}, "成功")
}
