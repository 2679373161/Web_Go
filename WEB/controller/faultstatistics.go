package controller

import (
	"encoding/json"
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IfaultstatisticsController interface{
	RestController

	Search(ctx *gin.Context)
    Gettype(ctx *gin.Context)
	Getid(ctx *gin.Context)
    Getequ(ctx *gin.Context)
	GetTypeErrorInfo(ctx *gin.Context)
	GetTypeErrorInfosum(ctx *gin.Context)
	Prodmonthfaultsum(ctx *gin.Context)
    Equnum(ctx *gin.Context)
    Equnumsum(ctx *gin.Context)

}
type faultstatisticsController struct {
	DB *gorm.DB
}
func (t faultstatisticsController) Create(ctx *gin.Context) {
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

func (t faultstatisticsController) Update(ctx *gin.Context) {
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

func (t faultstatisticsController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t faultstatisticsController) Prodmonthfaultsum(ctx *gin.Context) {
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	prodmonth1:=ctx.DefaultQuery("prodmonth1", "")
	prodmonth2:=ctx.DefaultQuery("prodmonth2", "")
	var typefaultsummary []model.MideaFault
	if prodmonth1=="NaN-NaN"{
		common.IndexDB.Raw("select prod_month,time_date,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef from midea_prodtimefault where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&time_date>=? &&time_date<=? group by prod_month order by e0  DESC,e1  DESC,e2  DESC,e3  DESC,e4 DESC ,e5  DESC,e6  DESC,e8  DESC,ea  DESC,ee DESC ,f2 DESC ,c0 DESC ,c1 DESC ,c2 DESC ,c3 DESC ,c4  DESC,c5 DESC ,c6 DESC ,c7 DESC ,c8 DESC ,eh DESC ,ef desc",starttime,endtime).Find(&typefaultsummary)
		  }else{
		common.IndexDB.Raw("select prod_month,time_date,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef from midea_prodtimefault where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&prod_month>=?&&prod_month<=?&&time_date>=? &&time_date<=? group by prod_month order by  e0  DESC,e1  DESC,e2  DESC,e3  DESC,e4 DESC ,e5  DESC,e6  DESC,e8  DESC,ea  DESC,ee DESC ,f2 DESC ,c0 DESC ,c1 DESC ,c2 DESC ,c3 DESC ,c4  DESC,c5 DESC ,c6 DESC ,c7 DESC ,c8 DESC ,eh DESC ,ef desc",prodmonth1,prodmonth2,starttime,endtime).Find(&typefaultsummary)
		  }  
      fmt.Println(typefaultsummary)
	  count:=len(typefaultsummary)
	  response.Success(ctx,gin.H{"data":typefaultsummary,"count":count },"读取成功")
}



func (t faultstatisticsController) Search(ctx *gin.Context) {
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	devtype := ctx.DefaultQuery("type", "")
	flag:= ctx.DefaultQuery("flag", "")
    //starttime:="2022-03-22"
    //endtime:="2022-03-22"
	var typefaultsummary []model.MideaFault

	if flag=="1"{
		if len(devtype)==0{
	  common.IndexDB.Raw("select dev_type,time_date,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&time_date>=? &&time_date<=? group by dev_type order by e0  DESC,e1  DESC,e2  DESC,e3  DESC,e4 DESC ,e5  DESC,e6  DESC,e8  DESC,ea  DESC,ee DESC ,f2 DESC ,c0 DESC ,c1 DESC ,c2 DESC ,c3 DESC ,c4  DESC,c5 DESC ,c6 DESC ,c7 DESC ,c8 DESC ,eh DESC ,ef desc",starttime,endtime).Find(&typefaultsummary)
		}else{
	  common.IndexDB.Raw("select dev_type,time_date,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=?&&time_date>=? &&time_date<=? group by dev_type order by  e0  DESC,e1  DESC,e2  DESC,e3  DESC,e4 DESC ,e5  DESC,e6  DESC,e8  DESC,ea  DESC,ee DESC ,f2 DESC ,c0 DESC ,c1 DESC ,c2 DESC ,c3 DESC ,c4  DESC,c5 DESC ,c6 DESC ,c7 DESC ,c8 DESC ,eh DESC ,ef desc",devtype,starttime,endtime).Find(&typefaultsummary)
	
		}
		}



	if flag=="2"{
	if len(devtype)==0{
  common.IndexDB.Raw("select dev_type,time_date,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef from midea_fault2 where time_date>=? &&time_date<=? group by dev_type order by e0  DESC,e1  DESC,e2  DESC,e3  DESC,e4 DESC ,e5  DESC,e6  DESC,e8  DESC,ea  DESC,ee DESC ,f2 DESC ,c0 DESC ,c1 DESC ,c2 DESC ,c3 DESC ,c4  DESC,c5 DESC ,c6 DESC ,c7 DESC ,c8 DESC ,eh DESC ,ef desc",starttime,endtime).Find(&typefaultsummary)
	}else{
		common.IndexDB.Raw("select dev_type,time_date,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef from midea_fault2 where dev_type=?&&time_date>=? &&time_date<=? group by dev_type order by  e0  DESC,e1  DESC,e2  DESC,e3  DESC,e4 DESC ,e5  DESC,e6  DESC,e8  DESC,ea  DESC,ee DESC ,f2 DESC ,c0 DESC ,c1 DESC ,c2 DESC ,c3 DESC ,c4  DESC,c5 DESC ,c6 DESC ,c7 DESC ,c8 DESC ,eh DESC ,ef desc",devtype,starttime,endtime).Find(&typefaultsummary)

	}
	}
	var count=len(typefaultsummary)
	response.Success(ctx,gin.H{"data":typefaultsummary,"count":count },"读取成功")
}

func (t faultstatisticsController) Gettype(ctx *gin.Context) {
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

	common.IndexDB.Raw("select distinct dev_type from midea_typefault").Find(&tableStoreDates1)

	common.IndexDB.Raw("select distinct dev_id from midea_fault2").Find(&tableStoreDates2) //设备号

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


func (t faultstatisticsController) Getequ(ctx *gin.Context) {
	devtype := ctx.DefaultQuery("type", "")

	type Types2 struct {
		Dev_id string `json:"dev_id" gorm:"type:varchar(255);not null"`
	}

	type Type2 struct {
		Value string `json:"value"`
	}

	var tableStoreDates2 []Types2

	var type2 []Type2



	common.IndexDB.Raw("select distinct dev_id from midea_fault2 where dev_type=?",devtype).Find(&tableStoreDates2) //设备号


	for _, tableDate2 := range tableStoreDates2 {

		var tableStoreDates2 Type2

		tableStoreDates2.Value = tableDate2.Dev_id
		type2 = append(type2, tableStoreDates2)

	}

	fmt.Println(type2)


  response.Success(ctx,gin.H{"data":type2 },"读取成功")
}
func (t faultstatisticsController) Getid(ctx *gin.Context) {
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	devtype := ctx.DefaultQuery("type", "0")
	devid := ctx.DefaultQuery("id", "0")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	prodmonth1:=ctx.DefaultQuery("prodmonth1", "")
	prodmonth2:=ctx.DefaultQuery("prodmonth2", "")
	var count int
	var idfault []model.MideaFault
	fmt.Println(perPage,currentPage)
	//common.IndexDB.Raw("select dev_id, dev_type,time_date,e0  ,e1  , e2  , e3  , e4  , e5  , e6  , e8  , ea  , ee  , f2  , c0  , c1  , c2  , c3  , c4  , c5  , c6  , c7  , c8  , eh  , ef   from midea_fault2 where time_date>=? &&time_date<=? ",starttime,endtime). Scan(&idfault)
 if prodmonth1=="NaN-NaN"{
	if devtype=="" && devid==""{
		fmt.Println("666")
		common.IndexDB.Raw("select * from midea_fault2 where(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&& time_date>=? &&time_date<=?  ",starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&& time_date>=? &&time_date<=?  ",starttime,endtime).Count(&count)


	}else if devtype==""{
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_id=?&& time_date>=? &&time_date<=?  ",devid,starttime,endtime).Count(&count)

		common.IndexDB.Raw("select * from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_id=? &&time_date>=? &&time_date<=?  ",devid,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else if devid==""{
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=?&& time_date>=? &&time_date<=?  ",devtype,starttime,endtime).Count(&count)
		common.IndexDB.Raw("select * from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=? &&time_date>=? &&time_date<=?  ",devtype,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else {
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_typr=?&&dev_id=?&& time_date>=? &&time_date<=?  ",devtype,devid,starttime,endtime).Count(&count)

		common.IndexDB.Raw("select * from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=?&&dev_id=?&&time_date>=? &&time_date<=?  ",devtype,devid,starttime,endtime).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}
}else{
	prodmonth1=prodmonth1+"-00"
	prodmonth2=prodmonth2+"-31"
	if devtype=="" && devid==""{
		common.IndexDB.Raw("select * from midea_fault2 where(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&& time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",starttime,endtime,prodmonth1,prodmonth2).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&& time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",starttime,endtime,prodmonth1,prodmonth2).Count(&count)


	}else if devtype==""{
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_id=?&& time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",devid,starttime,endtime,prodmonth1,prodmonth2).Count(&count)

		common.IndexDB.Raw("select * from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_id=? &&time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",devid,starttime,endtime,prodmonth1,prodmonth2).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else if devid==""{
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=?&& time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",devtype,starttime,endtime,prodmonth1,prodmonth2).Count(&count)
		common.IndexDB.Raw("select * from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=? &&time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",devtype,starttime,endtime,prodmonth1,prodmonth2).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}else {
		common.IndexDB.Table("midea_fault2").Where("(e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_typr=?&&dev_id=?&& time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",devtype,devid,starttime,endtime,prodmonth1,prodmonth2).Count(&count)

		common.IndexDB.Raw("select * from midea_fault2 where (e0!=0 || e1!=0 || e2!=0 || e3!=0 || e4!=0 || e5!=0 || e6!=0 || e8!=0 || ea!=0 || ee!=0 || f2!=0 || c0!=0|| c1!=0 || c2!=0 || c3!=0 || c4!=0 || c5!=0 || c6!=0 || c7!=0 || c8!=0 || eh!=0 || ef!=0)&&dev_type=?&&dev_id=?&&time_date>=? &&time_date<=?&& prod_time>=? &&prod_time<=?  ",devtype,devid,starttime,endtime,prodmonth1,prodmonth2).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&idfault)
	}
}
	fmt.Println("idfault",idfault)
	//fmt.Println("count",count)
	//count=len(idfault)
	response.Success(ctx,gin.H{"data":idfault,"count":count },"读取成功")
}

func (t faultstatisticsController) GetTypeErrorInfo(ctx *gin.Context) {

	var MideaTypeFault_info []model.MideaTypeFault
	dev_type := ctx.DefaultQuery("dev_type", "")
	error_type := ctx.DefaultQuery("error_type", "")
	end_time := ctx.DefaultQuery("end_time", "")
	start_time := ctx.DefaultQuery("start_time", "")
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	Info_num := 0
	if dev_type == "" && error_type == "" && end_time == "" && start_time == "" {
		common.IndexDB.Table("midea_typefault").Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Scan(&MideaTypeFault_info)
	} else if dev_type != "" && error_type == "" && end_time == "" && start_time == "" {

		common.IndexDB.Table("midea_typefault").Where("dev_type=?", dev_type).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where("dev_type=?", dev_type).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)
	} else if dev_type != "" && error_type != "" && end_time == "" && start_time == "" {

		string_test := "dev_type = ? and  " + error_type + " != ? "
		common.IndexDB.Table("midea_typefault").Where(string_test, dev_type, 0).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where(string_test, dev_type, 0).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)

	} else if dev_type != "" && error_type == "" && end_time != "" && start_time != "" {
		common.IndexDB.Table("midea_typefault").Where("dev_type=? and time_date BETWEEN ? and ? ", dev_type, start_time, end_time).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where("dev_type=? and time_date BETWEEN ? and ? ", dev_type, start_time, end_time).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)

	} else if dev_type != "" && error_type != "" && end_time != "" && start_time != "" {
		string_test := "dev_type=? and time_date BETWEEN ? and ? and " + error_type + " != ?"
		common.IndexDB.Table("midea_typefault").Where(string_test, dev_type, start_time, end_time, 0).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where(string_test, dev_type, start_time, end_time, 0).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)

	} else if dev_type == "" && error_type != "" && end_time == "" && start_time == "" {
		string_test := error_type + " != ?"
		common.IndexDB.Table("midea_typefault").Where(string_test, 0).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where(string_test, 0).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)

	} else if dev_type == "" && error_type == "" && end_time != "" && start_time != "" {
		common.IndexDB.Table("midea_typefault").Where("time_date BETWEEN ? and ? ", start_time, end_time).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where("time_date BETWEEN ? and ? ", start_time, end_time).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)

	} else if dev_type == "" && error_type != "" && end_time != "" && start_time != "" {
		string_test := "time_date BETWEEN ? and ? and " + error_type + " !=?"
		common.IndexDB.Table("midea_typefault").Where(string_test, start_time, end_time, 0).Count(&Info_num)
		common.IndexDB.Table("midea_typefault").Where(string_test, start_time, end_time, 0).Offset(perPage * (currentPage - 1)).Limit(perPage).Order("time_date desc").Find(&MideaTypeFault_info)
	}

	for index, Tablefragceshi := range MideaTypeFault_info {
		t := reflect.TypeOf(Tablefragceshi)
		v := reflect.ValueOf(Tablefragceshi)
		total := 0

		//此处与故障数挂钩 切记
		for k := 2; k < t.NumField()-4; k++ {
			// 结构体内容转出
			n, err := strconv.Atoi(Strval(v.Field(k).Interface()))
			if err != nil {
				panic(err)
			}
			total = total + n
			// 初始值处理
			if k == 2 {
				MideaTypeFault_info[index].MAX_Error = t.Field(k).Name
				MideaTypeFault_info[index].Max_Error_count = n
			} else if MideaTypeFault_info[index].Max_Error_count < n {
				MideaTypeFault_info[index].MAX_Error = t.Field(k).Name
				MideaTypeFault_info[index].Max_Error_count = n
			}
		}
		MideaTypeFault_info[index].Max_Other_count = total - MideaTypeFault_info[index].Max_Error_count
		MideaTypeFault_info[index].All_Error_count = total
	}
	fmt.Println(MideaTypeFault_info)
	fmt.Println(Info_num)
	response.Success(ctx, gin.H{"MideaTypeFault_info": MideaTypeFault_info, "total_info_num": Info_num}, "读取成功")
}
func (t faultstatisticsController) GetTypeErrorInfosum(ctx *gin.Context) {

	var MideaTypeFault_info []model.MideaTypeFault
	dev_type := ctx.DefaultQuery("dev_type", "")
	error_type := ctx.DefaultQuery("error_type", "")
	end_time := ctx.DefaultQuery("end_time", "")
	start_time := ctx.DefaultQuery("start_time", "")

	if dev_type == "" && error_type == "" && end_time == "" && start_time == "" {
		common.IndexDB.Table("midea_typefault").Select("dev_type, sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Scan(&MideaTypeFault_info)
	} else if dev_type != "" && error_type == "" && end_time == "" && start_time == "" {
		common.IndexDB.Table("midea_typefault").Where("dev_type=?", dev_type).Select("dev_type, sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Find(&MideaTypeFault_info)
	} else if dev_type != "" && error_type != "" && end_time == "" && start_time == "" {

		string_test := "dev_type = ? and  " + error_type + " != ? "
		common.IndexDB.Table("midea_typefault").Where(string_test, dev_type, 0).Select(" dev_type,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Find(&MideaTypeFault_info)

	} else if dev_type != "" && error_type == "" && end_time != "" && start_time != "" {
		common.IndexDB.Table("midea_typefault").Where("dev_type=? and time_date BETWEEN ? and ? ", dev_type, start_time, end_time).Select(" dev_type,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Find(&MideaTypeFault_info)
	} else if dev_type != "" && error_type != "" && end_time != "" && start_time != "" {
		string_test := "dev_type=? and time_date BETWEEN ? and ? and " + error_type + " != ?"
		common.IndexDB.Table("midea_typefault").Where(string_test, dev_type, start_time, end_time, 0).Select("dev_type, sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Find(&MideaTypeFault_info)
	} else if dev_type == "" && error_type != "" && end_time == "" && start_time == "" {
		string_test := error_type + " != ?"
		common.IndexDB.Table("midea_typefault").Where(string_test, 0).Find(&MideaTypeFault_info)
	} else if dev_type == "" && error_type == "" && end_time != "" && start_time != "" {
		common.IndexDB.Table("midea_typefault").Where("time_date BETWEEN ? and ? ", start_time, end_time).Select(" dev_type,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Find(&MideaTypeFault_info)
	} else if dev_type == "" && error_type != "" && end_time != "" && start_time != "" {
		string_test := "time_date BETWEEN ? and ? and " + error_type + " !=?"
		common.IndexDB.Table("midea_typefault").Where(string_test, start_time, end_time, 0).Select(" dev_type,sum(e0) e0,sum(e1) e1,sum(e2) e2,sum(e3) e3,sum(e4) e4,sum(e5) e5,sum(e6) e6,sum(e8) e8,sum(ea) ea,sum(ee) ee,sum(f2) f2,sum(c0) c0,sum(c1) c1,sum(c2) c2,sum(c3) c3,sum(c4) c4,sum(c5) c5,sum(c6) c6,sum(c7) c7,sum(c8) c8,sum(eh) eh,sum(ef) ef ").Group("dev_type").Order("e0 desc").Order("e1 desc").Order("e2 desc").Order("e3 desc").Order("e4 desc").Order("e5 desc").Order("e6 desc").Order("ea desc").Order("ee desc").Order("f2 desc").Order("c0 desc").Order("c1 desc").Order("c2 desc").Order("c3 desc").Order("c4 desc").Order("c5 desc").Order("c6 desc").Order("c7 desc").Order("c8 desc").Order("eh desc").Order("ef desc").Find(&MideaTypeFault_info)
	}
fmt.Println(MideaTypeFault_info)
	for index, Tablefragceshi := range MideaTypeFault_info {
		t := reflect.TypeOf(Tablefragceshi)
		v := reflect.ValueOf(Tablefragceshi)
		total := 0

		//此处与故障数挂钩 切记
		for k := 2; k < t.NumField()-4; k++ {
			// 结构体内容转出
			n, err := strconv.Atoi(Strval(v.Field(k).Interface()))
			if err != nil {
				panic(err)
			}
			total = total + n
			// 初始值处理
			if k == 2 {
				MideaTypeFault_info[index].MAX_Error = t.Field(k).Name
				MideaTypeFault_info[index].Max_Error_count = n
			} else if MideaTypeFault_info[index].Max_Error_count < n {
				MideaTypeFault_info[index].MAX_Error = t.Field(k).Name
				MideaTypeFault_info[index].Max_Error_count = n
			}
		}
		MideaTypeFault_info[index].Max_Other_count = total - MideaTypeFault_info[index].Max_Error_count
		MideaTypeFault_info[index].All_Error_count = total
	}
	fmt.Println(MideaTypeFault_info)
	var count=len(MideaTypeFault_info)
	response.Success(ctx, gin.H{"MideaTypeFault_info": MideaTypeFault_info,"total_info_num": count}, "读取成功")

}
// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func (t faultstatisticsController) Equnum(ctx *gin.Context) {
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	devtype := ctx.DefaultQuery("type", "0")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	var equfaultnum []model.MideaNumFault
	var count=0
	if devtype==""{
		common.IndexDB.Table("midea_numfault").Where("time_date>=? and time_date<=?",starttime,endtime).Count(&count)
		common.IndexDB.Table("midea_numfault").Where("time_date>=? and time_date<=?",starttime,endtime).Order("e0num desc").Order("e1num desc").Order("e2num desc").Order("e3num desc").Order("e4num desc").Order("e5num desc").Order("e6num desc").Order("eanum desc").Order("eenum desc").Order("f2num desc").Order("c0num desc").Order("c1num desc").Order("c2num desc").Order("c3num desc").Order("c4num desc").Order("c5num desc").Order("c6num desc").Order("c7num desc").Order("c8num desc").Order("ehnum desc").Order("efnum desc").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&equfaultnum)
	}else{
		common.IndexDB.Table("midea_numfault").Where("dev_type=? and time_date>=? and time_date<=?",devtype,starttime,endtime).Count(&count)
		common.IndexDB.Table("midea_numfault").Where("dev_type=? and time_date>=? and time_date<=?",devtype,starttime,endtime).Order("e0num desc").Order("e1num desc").Order("e2num desc").Order("e3num desc").Order("e4num desc").Order("e5num desc").Order("e6num desc").Order("eanum desc").Order("eenum desc").Order("f2num desc").Order("c0num desc").Order("c1num desc").Order("c2num desc").Order("c3num desc").Order("c4num desc").Order("c5num desc").Order("c6num desc").Order("c7num desc").Order("c8num desc").Order("ehnum desc").Order("efnum desc").Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&equfaultnum)
	}
	response.Success(ctx, gin.H{"data": equfaultnum,"count": count}, "读取成功")

}
func (t faultstatisticsController) Equnumsum(ctx *gin.Context) {

	devtype := ctx.DefaultQuery("type", "0")
	starttime := ctx.DefaultQuery("timeLow", "")
	endtime := ctx.DefaultQuery("timeHigh", "")
	var equfaultnumsum []model.MideaNumFault
	var count=0
	if devtype==""{
		common.IndexDB.Raw("select dev_type,time_date,COUNT( DISTINCT dev_id) as total,COUNT(DISTINCT CASE WHEN e0>0 THEN dev_id END) as e0num,COUNT(DISTINCT CASE WHEN e1>0 THEN dev_id END) as e1num,COUNT(DISTINCT CASE WHEN e2>0 THEN dev_id END) as e2num,COUNT(DISTINCT CASE WHEN e3>0 THEN dev_id END) as e3num, COUNT(DISTINCT CASE WHEN e4>0 THEN dev_id END) as e4num,COUNT(DISTINCT CASE WHEN e5>0 THEN dev_id END) as e5num,COUNT(DISTINCT CASE WHEN e6>0 THEN dev_id END) as e6num,COUNT(DISTINCT CASE WHEN e8>0 THEN dev_id END) as e8num, COUNT(DISTINCT CASE WHEN ea>0 THEN dev_id END) as eanum,COUNT(DISTINCT CASE WHEN ee>0 THEN dev_id END) as eenum,COUNT(DISTINCT CASE WHEN f2>0 THEN dev_id END) as f2num,COUNT(DISTINCT CASE WHEN c0>0 THEN dev_id END) as c0num,COUNT(DISTINCT CASE WHEN c1>0 THEN dev_id END) as c1num,COUNT(DISTINCT CASE WHEN c2>0 THEN dev_id END) as c2num,COUNT(DISTINCT CASE WHEN c3>0 THEN dev_id END) as c3num,COUNT(DISTINCT CASE WHEN c4>0 THEN dev_id END) as c4num,COUNT(DISTINCT CASE WHEN c5>0 THEN dev_id END) as c5num,COUNT(DISTINCT CASE WHEN c6>0 THEN dev_id END) as c6num,COUNT(DISTINCT CASE WHEN c7>0 THEN dev_id END) as c7num,COUNT(DISTINCT CASE WHEN c8>0 THEN dev_id END) as c8num,COUNT(DISTINCT CASE WHEN eh>0 THEN dev_id END) as ehnum,COUNT(DISTINCT CASE WHEN ef>0 THEN dev_id END) as efnum from midea_fault2 where time_date>=? and time_date<=? group by dev_type HAVING e0num!=0 || e1num!=0 || e2num!=0 || e3num!=0 || e4num!=0 || e5num!=0 || e6num!=0 || e8num!=0 || eanum!=0 || eenum!=0 || f2num!=0 || c0num!=0|| c1num!=0 || c2num!=0 || c3num!=0 || c4num!=0 || c5num!=0 || c6num!=0 || c7num!=0 || c8num!=0 || ehnum!=0 || efnum!=0 ",starttime,endtime).Find(&equfaultnumsum)
	}else{
		common.IndexDB.Raw("select dev_type,time_date,COUNT( DISTINCT dev_id) as total,COUNT(DISTINCT CASE WHEN e0>0 THEN dev_id END) as e0num,COUNT(DISTINCT CASE WHEN e1>0 THEN dev_id END) as e1num,COUNT(DISTINCT CASE WHEN e2>0 THEN dev_id END) as e2num,COUNT(DISTINCT CASE WHEN e3>0 THEN dev_id END) as e3num, COUNT(DISTINCT CASE WHEN e4>0 THEN dev_id END) as e4num,COUNT(DISTINCT CASE WHEN e5>0 THEN dev_id END) as e5num,COUNT(DISTINCT CASE WHEN e6>0 THEN dev_id END) as e6num,COUNT(DISTINCT CASE WHEN e8>0 THEN dev_id END) as e8num, COUNT(DISTINCT CASE WHEN ea>0 THEN dev_id END) as eanum,COUNT(DISTINCT CASE WHEN ee>0 THEN dev_id END) as eenum,COUNT(DISTINCT CASE WHEN f2>0 THEN dev_id END) as f2num,COUNT(DISTINCT CASE WHEN c0>0 THEN dev_id END) as c0num,COUNT(DISTINCT CASE WHEN c1>0 THEN dev_id END) as c1num,COUNT(DISTINCT CASE WHEN c2>0 THEN dev_id END) as c2num,COUNT(DISTINCT CASE WHEN c3>0 THEN dev_id END) as c3num,COUNT(DISTINCT CASE WHEN c4>0 THEN dev_id END) as c4num,COUNT(DISTINCT CASE WHEN c5>0 THEN dev_id END) as c5num,COUNT(DISTINCT CASE WHEN c6>0 THEN dev_id END) as c6num,COUNT(DISTINCT CASE WHEN c7>0 THEN dev_id END) as c7num,COUNT(DISTINCT CASE WHEN c8>0 THEN dev_id END) as c8num,COUNT(DISTINCT CASE WHEN eh>0 THEN dev_id END) as ehnum,COUNT(DISTINCT CASE WHEN ef>0 THEN dev_id END) as efnum from midea_fault2 where dev_type=? and time_date>=? and time_date<=? group by dev_type HAVING e0num!=0 || e1num!=0 || e2num!=0 || e3num!=0 || e4num!=0 || e5num!=0 || e6num!=0 || e8num!=0 || eanum!=0 || eenum!=0 || f2num!=0 || c0num!=0|| c1num!=0 || c2num!=0 || c3num!=0 || c4num!=0 || c5num!=0 || c6num!=0 || c7num!=0 || c8num!=0 || ehnum!=0 || efnum!=0",devtype,starttime,endtime).Find(&equfaultnumsum)
	}
	count=len(equfaultnumsum)
	response.Success(ctx, gin.H{"data": equfaultnumsum,"count": count}, "读取成功")

}

func (t faultstatisticsController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}









	func NewfaultstatisticsController ()IfaultstatisticsController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return faultstatisticsController{DB:db}
}