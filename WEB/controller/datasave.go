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
type IDatasaveController interface{
	RestController
	Datasave(ctx *gin.Context)
}
type DatasaveController struct {
	DB *gorm.DB
}
func (t DatasaveController) Create(ctx *gin.Context) {
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

func (t DatasaveController) Update(ctx *gin.Context) {
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

func (t DatasaveController) Show(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}

	response.Success(ctx,gin.H{"tableStoreDate":tableStoreDate},"读取成功")
}

func (t DatasaveController) Delete(ctx *gin.Context) {
	tableStoreDateId:=ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?",tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx,nil,"文章不存在")
		return
	}


	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx,gin.H{"tableStoreDate":tableStoreDate},"删除成功")

}

func (t DatasaveController) Datasave(ctx *gin.Context){

	//pagetype:=ctx.DefaultQuery("pagetype","0000")


	var tableStoreDate [] model.TableData

	dataparameter1 := make(map[string]string)

	t.DB.Table("table_data").Find(&tableStoreDate)
	for _,tableDate :=range tableStoreDate{

		dataparameter1[tableDate.Label]=tableDate.Value

	}
	response.Success(ctx,gin.H{"data":dataparameter1},"成功")
}
func NewDatasaveController ()IDatasaveController{
	db:=common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return DatasaveController{DB:db}
}