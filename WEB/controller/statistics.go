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

type IstatisticsController interface{
	RestController

	Search(ctx *gin.Context)
	Getequipment(ctx *gin.Context)
	Gettype(ctx *gin.Context)
	Getday(ctx *gin.Context)
	Getidnumber(ctx *gin.Context)
	Ceshi(ctx *gin.Context)


}
type statisticsController struct {
	DB *gorm.DB
}
func (t statisticsController) Create(ctx *gin.Context) {
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

func (t statisticsController) Update(ctx *gin.Context) {
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

func (t statisticsController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t statisticsController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t statisticsController) Search(ctx *gin.Context) {
	fmt.Println("历史详情")


	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")
	devtype := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")

	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))

	
	
	var count int
	//var abnormal_count int
	
	//common.IndexDB.Table("days_summaries").Where("dev_type = ? AND dev_id = ? AND time_date BETWEEN ? AND ?",dev_type,dev_id,timeLow,timeHigh).Find(&tableplace)

    var downloadplace [] model.Statistics
   
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
	//var tableplace  []model.Statistics
	fmt.Println(perPage,currentPage)
		
		
		if citycode!="0"&&len(devtype)!=0&&len(equipment)!=0{
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&city_code=?&&dev_type=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,citycode,devtype,equipment,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&city_code=?&&dev_type=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,citycode,devtype,equipment,starttime,endtime,0).Find(&downloadplace)


		}
	
	
	
		if citycode == "0" {
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&dev_type=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,devtype,equipment,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&dev_type=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,devtype,equipment,starttime,endtime,0).Find(&downloadplace)



		}
		if len(devtype) == 0 {
		//	common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&city_code=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,citycode,equipment,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&city_code=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,citycode,equipment,starttime,endtime,0).Find(&downloadplace)

		}
		if len(equipment) == 0 {
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&city_code=?&&dev_type=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,citycode,devtype,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&city_code=?&&dev_type=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,citycode,devtype,starttime,endtime,0).Find(&downloadplace)

			
		}
		if provincecode == "0" {
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where dev_type=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",devtype,equipment,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where dev_type=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",devtype,equipment,starttime,endtime,0).Find(&downloadplace)

		
		}
		if citycode == "0" && (len(devtype) == 0) {
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_markfrom days_summaries where province_code=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,equipment,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code=?&&dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,equipment,starttime,endtime,0).Find(&downloadplace)



		}
		if citycode == "0" && len(equipment) == 0 {


		//	common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code= ?&&dev_type= ?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,devtype,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code= ?&&dev_type= ?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,devtype,starttime,endtime,0).Find(&downloadplace)


		}
		if provincecode == "0" && len(devtype) == 0 {
			fmt.Println(equipment)
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",equipment,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where dev_id=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",equipment,starttime,endtime,0).Find(&downloadplace)

		}
		if provincecode == "0" && len(equipment) == 0 {
		//	common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where dev_type=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",devtype,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where dev_type=?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",devtype,starttime,endtime,0).Find(&downloadplace)

		}
		if len(devtype) == 0 && len(equipment) == 0 {
		//	common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where city_code=?&&province_code= ?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",citycode,provincecode,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where city_code=?&&province_code= ?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",citycode,provincecode,starttime,endtime,0).Find(&downloadplace)
		}
		if citycode == "0" && len(devtype) == 0 && len(equipment) == 0 {
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code= ?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where province_code= ?&&time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id",provincecode,starttime,endtime,0).Find(&downloadplace)
			
		}

		if provincecode == "0" && len(devtype) == 0 && len(equipment) == 0 {
			//common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_markfrom days_summaries where time_date>=? && time_date<=? and temp_score>?  GROUP BY dev_id",starttime,endtime,0).Find(&tableplace)
			common.IndexDB.Raw("select province_code,city_code,dev_type,dev_id,COUNT(un_stable_mark<60 or NULL)as un_stable_mark60,COUNT(heat_mark<60 or NULL)as heat_mark60,COUNT(un_heat_dev_mark<60 or NULL)as un_heat_dev_mark60,COUNT(over_shoot_mark<60 or NULL)as over_shoot_mark60,COUNT(temp_score>=80 or NULL)as high80,COUNT(temp_score>60&&temp_score<80 or NULL) as sixto80,COUNT(temp_score<60&&temp_score>0 or NULL) as low60,FLOOR(avg(temp_score)) as temp_score,FLOOR(avg(un_stable_mark)) as un_stable_mark,FLOOR(avg(un_heat_dev_mark)) as un_heat_dev_mark,FLOOR(avg(over_shoot_mark)) as over_shoot_mark,FLOOR(avg(heat_mark)) as heat_mark from days_summaries where time_date>=? && time_date<=? and temp_score>? GROUP BY dev_id order by temp_score",starttime,endtime,0).Find(&downloadplace)

		}
	

		

count=len(downloadplace)
	fmt.Println("count=", count)
	//fmt.Println("tableplace=", downloadplace)

	
	
  // if downloadflag=="1"{
//	if flag == "1" {
	//	t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fault_help where temp_score>=0 and temp_score>? GROUP BY dev_id having abnormal_flag!=0").Find(&tableplace)
	
	//}
	//if flag == "0" {

	//	t.DB.Raw("select province_code,city_code,dev_type,dev_id,FLOOR(avg(temp_score)) as temp_score,sum(abnormal_flag) as abnormal_flag from fault_help where temp_score>=0 and temp_score>? GROUP BY dev_id having abnormal_flag=0").Find(&tableplace)
//	}

//	count = len(tableplace)
//}


	var downloadplace1 []model.TableDate3
	for _, tableDate := range downloadplace {
		var tableStoreDates4 model.TableDate3
		t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
		downloadplace1 = append(downloadplace1, tableStoreDates4)

	}
	

	//fmt.Println("count=", count)

	response.Success(ctx, gin.H{"data":downloadplace, "data1":downloadplace1, "count": count}, "成功")
}






func (t statisticsController) Getequipment(ctx *gin.Context) {
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

		common.IndexDB.Raw("select distinct dev_id from days_summaries ").Find(&tableStoreDates4)
		//common.IndexDB.Table("statistics").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "1" {

		common.IndexDB.Raw("select distinct dev_id from days_summaries where province_code=? and dev_type=?", province_code, dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("statistics").Where("province_code = ? and dev_type = ? ", province_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else if flag == "2" {
		common.IndexDB.Raw("select distinct dev_id from days_summaries where city_code=? and dev_type=?", city_code, dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("statistics").Where("city_code = ? and dev_type = ? ", city_code, dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	} else {
		fmt.Println(dev_type)
		common.IndexDB.Raw("select distinct dev_id from days_summaries where dev_type=?", dev_type).Find(&tableStoreDates4)
		//common.IndexDB.Table("statistics").Where("dev_type = ? ", dev_type).Find(&tableStoreDates4)
		fmt.Println(tableStoreDates4)
		for _, tableDate := range tableStoreDates4 {
			var tableStoreDates3 Equipment1
			tableStoreDates3.Value = tableDate.Dev_id
			type1 = append(type1, tableStoreDates3)
		}

		response.Success(ctx, gin.H{"data": type1}, "成功")
	}
}
//func (t statisticsController) Datachart(ctx *gin.Context) {
//	equipment := ctx.DefaultQuery("equipment", "")
//	starttime := ctx.DefaultQuery("timeLow", "")
//	endtime := ctx.DefaultQuery("timeHigh", "")
//	var tableplace []model.statistics
//	if len(equipment) != 0 {
//		equipment = equipment
//	}
//	if len(starttime) != 0 && len(endtime) != 0 {
//		starttime = starttime
//		endtime = endtime
///	}
///	common.IndexDB.Table("statistics").Where(" dev_id=? and time_date >= ? and time_date <=?", equipment, starttime, endtime).Order("time_date ASC").Find(&tableplace)
//	fmt.Println(tableplace)
//	response.Success(ctx, gin.H{"data": tableplace}, "成功")
func (t statisticsController) Getidnumber(ctx *gin.Context) {
	fmt.Println("历史详情")
	provincecode := ctx.DefaultQuery("province_code", "0")
	citycode := ctx.DefaultQuery("city_code", "0")
	devtype := ctx.DefaultQuery("type", "")
	equipment := ctx.DefaultQuery("equipment", "")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	searchflag :=ctx.DefaultQuery("searchflag", "1")
	//perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	//currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	var count int
	var downloadplace [] model.Statistics

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
	if searchflag=="1" {
	common.IndexDB.Raw("SELECT province_code,COUNT(avg_score>=80 or NULL) high80device,COUNT(avg_score>=60&&avg_score<=80 or NULL) sixto80device,COUNT(avg_score<=60 or NULL) low60device from (select DISTINCT province_code as province_code,IFNULL(sum(temp_score*temp_num)/sum(temp_num),0) as avg_score, dev_id from days_summaries where time_date>=? && time_date<=? GROUP BY dev_id  ) zz group by province_code",starttime,endtime).Find(&downloadplace)
	//common.IndexDB.Raw("SELECT province_code,COUNT(avg_score>=80 or NULL) high80device,COUNT(avg_score>=60&&avg_score<=80 or NULL) sixto80device,COUNT(avg_score<=60 or NULL) low60device from (select DISTINCT province_code as province_code,IFNULL(sum(temp_score*temp_num)/sum(temp_num),0) as avg_score, dev_id from days_summaries where time_date>=? && time_date<=? GROUP BY dev_id  ) zz group by province_code",starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	fmt.Println("table",downloadplace)
	count=len(downloadplace)
	var downloadplace1 []model.TableDate3
	for _, tableDate := range downloadplace {
	var tableStoreDates4 model.TableDate3
	t.DB.Table("midea_loc_code").Where("province_code = ?", tableDate.Province_code).Find(&tableStoreDates4)
	downloadplace1 = append(downloadplace1, tableStoreDates4)
	}
	response.Success(ctx, gin.H{"data":downloadplace, "data1":downloadplace1, "count": count}, "成功")
	}else if searchflag=="2" {
	common.IndexDB.Raw("SELECT city_code,COUNT(avg_score>=80 or NULL) high80device,COUNT(avg_score>=60&&avg_score<=80 or NULL) sixto80device,COUNT(avg_score<=60 or NULL) low60device from (select DISTINCT city_code as city_code,IFNULL(sum(temp_score*temp_num)/sum(temp_num),0) as avg_score, dev_id from days_summaries where time_date>=? && time_date<=? GROUP BY dev_id ) zz group by city_code",starttime,endtime).Find(&downloadplace)
	//common.IndexDB.Raw("SELECT city_code,COUNT(avg_score>=80 or NULL) high80device,COUNT(avg_score>=60&&avg_score<=80 or NULL) sixto80device,COUNT(avg_score<=60 or NULL) low60device from (select DISTINCT city_code as city_code,IFNULL(sum(temp_score*temp_num)/sum(temp_num),0) as avg_score, dev_id from days_summaries where time_date>=? && time_date<=? GROUP BY dev_id ) zz group by city_code",starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&downloadplace)
	fmt.Println("table",downloadplace)
	count=len(downloadplace)
	var downloadplace1 []model.TableDate3
	for _, tableDate := range downloadplace {
	var tableStoreDates4 model.TableDate3
	t.DB.Table("midea_loc_code").Where("city_code = ?", tableDate.City_code).Find(&tableStoreDates4)
	downloadplace1 = append(downloadplace1, tableStoreDates4)
	}
	response.Success(ctx, gin.H{"data":downloadplace, "data1":downloadplace1, "count": count}, "成功")
	}
}

func (t statisticsController) Ceshi(ctx *gin.Context) {
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

	common.IndexDB.Raw("select distinct dev_type from days_summaries").Find(&tableStoreDates1)

	common.IndexDB.Raw("select distinct dev_id from days_summaries").Find(&tableStoreDates2) //设备号

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

func (t statisticsController) Gettype(ctx *gin.Context) {
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
		common.IndexDB.Raw("select distinct dev_type from days_summaries where province_code = ? ", province_code).Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id from days_summaries where province_code = ? ", province_code).Find(&tableStoreDates2) //设备号
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
		common.IndexDB.Raw("select distinct dev_type from days_summaries where city_code = ?", city_code).Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id   from days_summaries where city_code = ?", city_code).Find(&tableStoreDates2) //设备号
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
		common.IndexDB.Raw("select distinct dev_type from days_summaries").Find(&tableStoreDates1)
		common.IndexDB.Raw("select distinct dev_id from days_summaries").Find(&tableStoreDates2) //设备号
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
func (t statisticsController) Getday(ctx *gin.Context) {
	fmt.Println("666")
	var day_infor [] model.Statistics
	flag := ctx.DefaultQuery("flag", "0")
	equipment := ctx.DefaultQuery("equipment", "")


	if flag=="1"{
		common.IndexDB.Table("days_summaries").Where(" dev_id=?", equipment).Find(&day_infor)
	}else if flag=="0"{
		common.IndexDB.Table("days_summaries").Where(" dev_id=?",equipment).Find(&day_infor)

	}

	response.Success(ctx, gin.H{"data": day_infor}, "成功")

}



	func NewstatisticsController ()IstatisticsController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return statisticsController{DB:db}
}