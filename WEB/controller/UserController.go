package controller

import (
	"ginEssential/common"
	"ginEssential/dto"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	DB:=common.GetDB()
	//法一：使用map获取请求参数
	//var requestMap=make(map[string]string)//
	//json.NewDecoder(ctx.Request.Body).Decode(&requestMap)//将json数据解析到requestMap中

	//法二：结构体获取请求参数
	//var requestUser=model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)

	//法三：gin框架中的bind函数
	var requestUser=model.User{}
	ctx.Bind(&requestUser)
	//获取参数
	name:=requestUser.Name//ctx.PostForm("name")//发送post请求
	telephone:=requestUser.Telephone//ctx.PostForm("telephone")
	password:=requestUser.Password//ctx.PostForm("password")

	//数据验证
	if len(telephone)!=11{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password)<6{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}
	//如果名称为传，给一个10位随机字符串
	if len(name)==0{
		name=util.RandomString(10)
	}
	log.Println(name,telephone,password)//日志记录，打印信息
	//判断手机号是否存在
	if isTelephoneExist(DB,telephone){
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户已经存在")
		return
	}
	//创建用户
	hasedPassword,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)//密码加密
	if err!=nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,"加密错误")//系统级错误
		return
	}

	newUser:=model.User{
		Name:name,
		Telephone: telephone,
		Password: string(hasedPassword),
	}

	DB.Create(&newUser)
	//发放token
	token,err:=common.ReleaseToken(newUser)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"code":500,"msg":"系统异常"})
		log.Printf("token generate err: %v",err)
		return
	}
	//返回结果
	response.Success(ctx,gin.H{"token":token},"注册成功")
}

func Login(ctx *gin.Context){
	DB:=common.GetDB()
	var requestUser=model.User{}
	ctx.Bind(&requestUser)
	//获取参数
	telephone:=requestUser.Telephone//ctx.PostForm("telephone")
	password:=requestUser.Password//ctx.PostForm("password")

	//数据验证
	if len(telephone)!=11{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password)<6{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	DB.Where("telephone=?",telephone).First(&user)
	if user.ID==0{
		ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户不存在"})
		return
	}
	//判断密码是否正确
	if err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password));err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"code":400,"msg":"密码错误"})
		return
	}
	//发放token
	token,err:=common.ReleaseToken(user)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"code":500,"msg":"系统异常"})
		log.Printf("token generate err: %v",err)
		return
	}
	//返回结果
	response.Success(ctx,gin.H{"token":token},"登录成功")

}

func Info(ctx *gin.Context)  {
	user,_:=ctx.Get("user")

	ctx.JSON(http.StatusOK,gin.H{"code":200,"data":gin.H{"user":dto.ToUserDto(user.(model.User))}})
}

func isTelephoneExist(db *gorm.DB,telephone string)bool{
	var user model.User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID!=0{
		return true
	}
	return false
}

