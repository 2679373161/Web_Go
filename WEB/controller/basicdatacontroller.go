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
type IBasicdatacontrollerController interface{
	RestController
	Datasave(ctx *gin.Context)
	Datatype(ctx *gin.Context)
	Modifydata(ctx *gin.Context)
	Modifytype(ctx *gin.Context)
	Dataequipment(ctx *gin.Context)
	Modifyequipment(ctx *gin.Context)
}
type BasicdatacontrollerController struct {
	DB *gorm.DB
}
func (t BasicdatacontrollerController) Create(ctx *gin.Context) {
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
func  typeadd( type1 []string,  newtype string )  (type2 []string,flag bool)  {
	flag =false
	for  i:=0;i< len(type1);i++{
		if(newtype==type1[i]){
			flag=true
		}
	}
	if flag==false{
		type1 = append(type1, newtype)

	}
	type2=type1
	return

}
func  numberadd( type1 []string, number1  []int, newtype string )  (type2 []string, number  []int,flag bool)  {
	flag =false
	for  i:=0;i< len(type1);i++{
		if(newtype==type1[i]){
			flag=true
			number1[i]++
		}
	}
	if flag==false{
		type1 = append(type1, newtype)
		number1 = append(number1, 1)
	}
	type2=type1
	number=number1
	return

}
func (t BasicdatacontrollerController) Update(ctx *gin.Context) {
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

func (t BasicdatacontrollerController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t BasicdatacontrollerController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t BasicdatacontrollerController) Datasave(ctx *gin.Context){

	flag1:=ctx.DefaultQuery("flag","0")
	citycode:=ctx.DefaultQuery("city","0")
	//var count int
	var tableStoreDate [] model.Tableselect
	var tableStoreDates4 [] model.TableDate3
	var tableStoreDates [] model.TableDate3
	//var dataparameter1 [] string
	//var dataparameter2 [] string
	var city1 []  string
	var province1 [] string
	var city []  string
	var province []string
	var flag2  bool
	//var tableStoreDate1 [] model.Tableplace
	SheBei := make(map[string] []int)
	XingHao := make(map[string] []string)
	if flag1=="0"{
		t.DB.Table("bo").Find(&tableStoreDate)
	}else if flag1=="1"{
		t.DB.Table("bo").Where("city_code=?",citycode).Find(&tableStoreDate)
	}else{
		t.DB.Table("midea_loc_code").Where("province_code =?",citycode).Find(&tableStoreDates)
		for _, tableDate1 := range tableStoreDates {
			province=append(province,tableDate1.City_code)
		}
		t.DB.Table("bo").Where("city_code in (?)",province).Find(&tableStoreDate)
	}


	for _, tableDate1 := range tableStoreDate {
		city=append(city,tableDate1.City_code)
		t.DB.Table("midea_loc_code").Where("city_code =?",city[len(city)-1]).Find(&tableStoreDates4)

		var flag bool=false
		for  i:=0;i< len(city1);i++{
			if(tableStoreDates4[0].Dev_city==city1[i]){
				flag=true
				XingHao[tableStoreDates4[0].Dev_city],flag2=typeadd(XingHao[tableStoreDates4[0].Dev_city],tableDate1.Dev_type)


			}
		}
		if flag==false{
			city1 = append(city1, tableStoreDates4[0].Dev_city)
			province1=append(province1, tableStoreDates4[0].Dev_province)
			XingHao[tableStoreDates4[0].Dev_city]=[]string{}
			SheBei[tableStoreDates4[0].Dev_city]=[]int{0,0,0}
			XingHao[tableStoreDates4[0].Dev_city],flag2=typeadd(XingHao[tableStoreDates4[0].Dev_city],tableDate1.Dev_type)
			//fmt.Println(XingHao)
		}
		SheBei[tableStoreDates4[0].Dev_city][0]++
		if flag2==false{
			SheBei[tableStoreDates4[0].Dev_city][1]++

		}
		if tableDate1.Disable=="1"{
			SheBei[tableStoreDates4[0].Dev_city][2]=1
		}
	}
	fmt.Println(SheBei)
	fmt.Println(XingHao)
	fmt.Println(province1)

	response.Success(ctx,gin.H{"SheBei":SheBei,"XingHao":XingHao,"city1":city1,"province1":province1},"成功")
}
/*=========================================================
 * 功能描述: 查询型号数据
 * 输出指标: 型号分布城市数，是否处理标志位
 =========================================================*/

func (t BasicdatacontrollerController) Datatype(ctx *gin.Context){

	flag1:=ctx.DefaultQuery("flag","0")
	citycode:=ctx.DefaultQuery("city","0")
	typecode:=ctx.DefaultQuery("type","0")
	//var count int
	var tableStoreDate [] model.Tableselect
	var tableStoreDates4 [] model.TableDate3
	var tableStoreDates [] model.TableDate3
	//var dataparameter1 [] string
	//var dataparameter2 [] string
	var type1 []  string
	var type2 []  string
	var city []  string

	var flag2  bool
	var flag3  bool
	//var tableStoreDate1 [] model.Tableplace
	SheBei := make(map[string] []int)
	XingHao := make(map[string] []string)
	province := make(map[string] []string)
	citynumber := make(map[string] []int)
	number:= make(map[string] []int)
	if flag1=="0"{
		t.DB.Table("bo").Find(&tableStoreDate)
	}else if flag1=="1"{
		t.DB.Table("bo").Where("dev_type=?",typecode).Find(&tableStoreDate)
	}else{
		if typecode=="0"{
			var province1  []string
			t.DB.Table("midea_loc_code").Where("province_code =?",citycode).Find(&tableStoreDates)
			for _, tableDate1 := range tableStoreDates {
				province1=append(province1,tableDate1.City_code)
			}
			t.DB.Table("bo").Where("city_code in (?) ",province1).Find(&tableStoreDate)
			//t.DB.Table("id1").Where("city_code=?",citycode).Find(&tableStoreDate)
		}else {
			var province1  []string
			t.DB.Table("midea_loc_code").Where("province_code =?",citycode).Find(&tableStoreDates)
			for _, tableDate1 := range tableStoreDates {
				province1=append(province1,tableDate1.City_code)
			}
			t.DB.Table("bo").Where("city_code in (?) AND dev_type=?",province1,typecode).Find(&tableStoreDate)
			//	t.DB.Table("id1").Where("city_code=? AND dev_type=?",citycode,typecode).Find(&tableStoreDate)
		}

	}

	//   fmt.Println("00")
	//选择省份下城市型号
	if flag1=="2"{
		for _, tableDate1 := range tableStoreDate {
			city=append(city,tableDate1.City_code)
			type1=append(type1,tableDate1.Dev_type)
			t.DB.Table("midea_loc_code").Where("city_code =?",city[len(city)-1]).Find(&tableStoreDates4)

			var flag bool=false
			for  i:=0;i< len(type2);i++{
				if(tableDate1.Dev_type==type2[i]){
					flag=true

					XingHao[tableDate1.Dev_type],flag2=typeadd(XingHao[tableDate1.Dev_type],tableStoreDates4[0].Dev_city)
					//   t.DB.Table("midea_loc_code").Where("dev_city =?",tableStoreDates4[0].Dev_city).Find(&tableStoreDates)

					province[tableDate1.Dev_type],number[tableDate1.Dev_type],flag3=numberadd(province[tableDate1.Dev_type],number[tableDate1.Dev_type],tableStoreDates4[0].Dev_city)
				}
			}
			if flag==false{
				type2 = append(type2, tableDate1.Dev_type)
				XingHao[tableDate1.Dev_type]=[]string{}
				province[tableDate1.Dev_type]=[]string{}

				SheBei[tableDate1.Dev_type]=[]int{0,0,0}
				number[tableDate1.Dev_type]=[]int{}
				XingHao[tableDate1.Dev_type],flag2=typeadd(XingHao[tableDate1.Dev_type],tableStoreDates4[0].Dev_city)
				//   t.DB.Table("midea_loc_code").Where("dev_city =?",tableStoreDates4[0].Dev_city).Find(&tableStoreDates)

				province[tableDate1.Dev_type],number[tableDate1.Dev_type],flag3=numberadd(province[tableDate1.Dev_type],number[tableDate1.Dev_type],tableStoreDates4[0].Dev_city)
				//  province[tableDate1.Dev_type],number[tableDate1.Dev_type],flag3=numberadd(province[tableDate1.Dev_type],number[tableDate1.Dev_type],tableStoreDates[0].Dev_province)

			}
			SheBei[tableDate1.Dev_type][0]++
			if flag2==false{
				SheBei[tableDate1.Dev_type][1]++

			}

			if tableDate1.Discityselect=="1"{
				SheBei[tableDate1.Dev_type][2]=1
			}
		}


		province1 := make(map[string] []string)
		for j:=0;j< len(type2);j++{
			//	fmt.Println(type2[0])
			var name string=type2[j]
			citynumber[name]=[]int{}
			province1[name]=[]string{}

			for i:=0;i< len(XingHao[name]);i++{
				t.DB.Table("midea_loc_code").Where("dev_city =?",XingHao[name][i]).Find(&tableStoreDates)
				province1[name],citynumber[name],flag3=numberadd(province1[name],citynumber[name],tableStoreDates[0].Dev_province)
			}

		}
		fmt.Println(flag3)

	}else{
		//只对型号选择
		for _, tableDate1 := range tableStoreDate {
			city=append(city,tableDate1.City_code)
			type1=append(type1,tableDate1.Dev_type)
			t.DB.Table("midea_loc_code").Where("city_code =?",city[len(city)-1]).Find(&tableStoreDates4)

			var flag bool=false
			for  i:=0;i< len(type2);i++{
				if(tableDate1.Dev_type==type2[i]){
					flag=true
					//	XingHao[tableStoreDates4[0].Dev_city],flag2=typeadd(XingHao[tableStoreDates4[0].Dev_city],tableDate1.Dev_type)
					XingHao[tableDate1.Dev_type],flag2=typeadd(XingHao[tableDate1.Dev_type],tableStoreDates4[0].Dev_city)
					t.DB.Table("midea_loc_code").Where("dev_city =?",tableStoreDates4[0].Dev_city).Find(&tableStoreDates)
					//	province[tableDate1.Dev_type],flag3=typeadd(province[tableDate1.Dev_type],tableStoreDates[0].Dev_province)
					province[tableDate1.Dev_type],number[tableDate1.Dev_type],flag3=numberadd(province[tableDate1.Dev_type],number[tableDate1.Dev_type],tableStoreDates[0].Dev_province)
				}
			}
			if flag==false{
				type2 = append(type2, tableDate1.Dev_type)
				XingHao[tableDate1.Dev_type]=[]string{}
				province[tableDate1.Dev_type]=[]string{}

				SheBei[tableDate1.Dev_type]=[]int{0,0,0}
				number[tableDate1.Dev_type]=[]int{}
				XingHao[tableDate1.Dev_type],flag2=typeadd(XingHao[tableDate1.Dev_type],tableStoreDates4[0].Dev_city)
				t.DB.Table("midea_loc_code").Where("dev_city =?",tableStoreDates4[0].Dev_city).Find(&tableStoreDates)


				province[tableDate1.Dev_type],number[tableDate1.Dev_type],flag3=numberadd(province[tableDate1.Dev_type],number[tableDate1.Dev_type],tableStoreDates[0].Dev_province)
				//fmt.Println(XingHao)
			}
			SheBei[tableDate1.Dev_type][0]++
			if flag2==false{
				SheBei[tableDate1.Dev_type][1]++

			}
			//if flag3==false{
			//	number[tableDate1.Dev_type][1]++
			//
			//}
			if tableDate1.Distype=="1"{
				SheBei[tableDate1.Dev_type][2]=1
			}
		}


		province1 := make(map[string] []string)
		for j:=0;j< len(type2);j++{
			//	fmt.Println(type2[0])
			var name string=type2[j]
			citynumber[name]=[]int{}
			province1[name]=[]string{}

			for i:=0;i< len(XingHao[name]);i++{
				t.DB.Table("midea_loc_code").Where("dev_city =?",XingHao[name][i]).Find(&tableStoreDates)
				province1[name],citynumber[name],flag3=numberadd(province1[name],citynumber[name],tableStoreDates[0].Dev_province)
			}

		}
		fmt.Println(flag3)
	}


	//fmt.Println(citynumber)

	response.Success(ctx,gin.H{"SheBei":SheBei,"XingHao":XingHao,"city1":type2,"province":province,"number":number,"citynumber":citynumber},"成功")
}

/*=========================================================
 * 功能描述: 查询设备数据
 * 输出指标: 设备分布，是否处理标志位
 =========================================================*/

func (t BasicdatacontrollerController) Dataequipment(ctx *gin.Context){

	flag1:=ctx.DefaultQuery("flag","0")
	provincecode:=ctx.DefaultQuery("provincecode","0")
	//dev_id:=ctx.DefaultQuery("dev_id","0")
	//typecode:=ctx.DefaultQuery("type","0")

	var tableStoreDate [] model.Tableselect
	var tableStoreDate1 [] model.TableDate3
	var tableStoreDates [] model.TableDate3
	var city [] string
	var province [] string
	var dev_type [] string
	if flag1=="0"{
		t.DB.Table("bo").Find(&tableStoreDate)
	}else if flag1=="1"{
		var province1  []string
		t.DB.Table("midea_loc_code").Where("province_code =?",provincecode).Find(&tableStoreDates)
		for _, tableDate1 := range tableStoreDates {
			province1=append(province1,tableDate1.City_code)
		}
		t.DB.Table("bo").Where("city_code in (?)",province1).Find(&tableStoreDate)
	}else{
		t.DB.Table("bo").Where("city_code= ? ",provincecode).Find(&tableStoreDate)
	}
	for _, tableDate1 := range tableStoreDate {

		var flag bool=false
		for  i:=0;i< len(dev_type);i++{
			if(tableDate1.Dev_type==dev_type[i]){
				flag=true
			}
		}
		if flag==false{
			dev_type = append(dev_type, tableDate1.Dev_type)
		}
		t.DB.Table("midea_loc_code").Where("city_code =?",tableDate1.City_code).Find(&tableStoreDate1)
		city=append(city,tableStoreDate1[0].Dev_city)
		province=append(province,tableStoreDate1[0].Dev_province)
	}
	fmt.Println(dev_type)
	response.Success(ctx,gin.H{"data":tableStoreDate,"province":province,"city":city,"dev_type":dev_type},"成功")
}

/*=========================================================
 * 功能描述: 更新城市标志位
 * 输出指标: 无
 =========================================================*/

func (t BasicdatacontrollerController) Modifydata(ctx *gin.Context){


	city:=ctx.DefaultQuery("city","0")
	selectflag:=ctx.DefaultQuery("selectflag","0")
	var  tableselect []model.Tableselect
	//var  tableselect1 []model.Tableselect
	var tableStoreDates4 model.TableDate3
	SheBei := make(map[string] []int)
	XingHao := make(map[string] []string)
	//selectdata:=strings.Split(categoryType,"-")
	//unselectdata:=strings.Split(categoryType1,"-")
	//fmt.Println(categoryType)
		t.DB.Table("midea_loc_code").Where("dev_city=?",city).Find(&tableStoreDates4)
		// t.DB.Table("id1").Where("city_code=?",tableStoreDates4[0].City_code).Find(&tableselect)
		//for j:=0;j< len(tableselect);j++ {
		// var  tableselects model.Tableselect
		// tableselects.Disable="1"
		// tableselects.City_code=tableselect[j].City_code
		// tableselects.Dev_Id=tableselect[j].Dev_Id
		// tableselects.Dev_type=tableselect[j].Dev_type
		//
		// tableselect1=append(tableselect1,tableselects)
		//}
		//fmt.Println(tableselect1)
		if err:=t.DB.Table("bo").Model(&tableselect).Where("city_code=?",tableStoreDates4.City_code).Update("disable", selectflag).Error;err!=nil{
			//panic(err)
			fmt.Println(err)
			response.Fail(ctx,nil,"更新失败")
			return
		}



	//fmt.Println(categoryType[0])
	//fmt.Println(len(categoryType1))

	response.Success(ctx,gin.H{"SheBei":SheBei,"XingHao":XingHao},"成功")
}

/*=========================================================
 * 功能描述: 更新型号标志位
 * 输出指标: 无
 =========================================================*/

func (t BasicdatacontrollerController) Modifytype(ctx *gin.Context){

	flag1:=ctx.DefaultQuery("flag","0")
	selectflag:=ctx.DefaultQuery("selectflag","0")
	city:=ctx.DefaultQuery("city","0")
	dev_type:=ctx.DefaultQuery("dev_type","0")
	var  tableselect []model.Tableselect
	var tableStoreDates model.TableDate3
	//var  tableselect1 []model.Tableselect

	SheBei := make(map[string] []int)
	XingHao := make(map[string] []string)
    fmt.Print(city)

		if flag1=="1"{
			t.DB.Table("midea_loc_code").Where("dev_city=?",city).Find(&tableStoreDates)
			fmt.Print(city)
			if err:=t.DB.Table("bo").Model(&tableselect).Where("dev_type=? AND city_code=?",dev_type,tableStoreDates.City_code).Update("discityselect", selectflag).Error;err!=nil{
				//panic(err)
				fmt.Println(err)
				response.Fail(ctx,nil,"更新失败")
				return
			}

		}else{
			if err:=t.DB.Table("bo").Model(&tableselect).Where("dev_type=?",dev_type).Update("distype", selectflag).Error;err!=nil{
				//panic(err)
				fmt.Println(err)
				response.Fail(ctx,nil,"更新失败")
				return
			}



	}

	//fmt.Println(categoryType[0])
	//fmt.Println(len(categoryType1))

	response.Success(ctx,gin.H{"SheBei":SheBei,"XingHao":XingHao},"成功")
}
/*=========================================================
 * 功能描述: 更新型号标志位
 * 输出指标: 无
 =========================================================*/

func (t BasicdatacontrollerController) Modifyequipment(ctx *gin.Context){
	id:=ctx.DefaultQuery("id","0")
	flag:=ctx.DefaultQuery("flag","0")
	var  tableselect []model.Tableselect

	t.DB.Table("bo").Model(&tableselect).Where("dev_id=?",id).Update("disequipment", flag)
}
func NewBasicdatacontrollerController ()IBasicdatacontrollerController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return BasicdatacontrollerController{DB:db}
}
