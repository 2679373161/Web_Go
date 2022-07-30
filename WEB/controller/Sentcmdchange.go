package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var sourcessssu = []string{"null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null",
	"null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null",
	"null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null", "null"}

type Htemper struct {
			Code    []string
			Pvule   []string
			Hvule   []string
			Pvule16 []string
			Hvule16 []string
			Flag    []int
			Flagerror int
			Index   []int
			Setflagerror11 string
		}

type Hfhtemp struct {
	Nowrow []string
	Hourow []string
	Perrow []string
}


/*=========================================================
 * 函数名称： Modifyfaultydeviceparameters
 * 功能描述:  恒温参数、熄火、风机、振荡修需要修改的值(传给前端)
 * 作    者：李贝贝
 * 修    改：李赞辉
 * 修改内容：先对row[16]赋值再存储修改后的下发数据，熄火次数记录函数去除，将其添加至下发参数的函数中
 =========================================================*/

func Modifyfaultydeviceparameters(ctx *gin.Context) {

	devType := ctx.DefaultQuery("dev_type", "5110173N")    //日表中监控城市号
	devId := ctx.DefaultQuery("devId", "182518931003942")  //输入需要改动的设备号
	timedate := ctx.DefaultQuery("timedate", "2022-06-06") //修改的时间

	db := common.GetIndexDB()
	var dailyMonitoring model.DailyMonitoring
	var htemper Htemper//e1整机关键参数和PL
	var windblockage Htemper//e1风堵等级大于8改参
	var htemperc4 Htemper//C4改参
	var htemperheng Htemper//恒温改参
	var htempzhen Htemper//振荡改参

	var xhrow Hfhtemp
	var htrow Hfhtemp

	var rowzhen Hfhtemp
	db.Table("daily_monitorings").Where("time_date = ? and dev_id = ? and dev_type = ?", timedate, devId, devType).Find(&dailyMonitoring)
	fmt.Println(dailyMonitoring)
	if dailyMonitoring.E1 != 0 { //判断熄火是否存在

		var numPerChange NumPerChanges
		db.Table("num_per_changes").Where("dev_id = ? AND fault_type= ?", devId, "E1").Find(&numPerChange)
		var e1fragementflag model.Daily_monitorings_error
		//******************************************风堵导致的熄火调参*************************************************//
       db.Table("daily_monitorings").Where("dev_id=? and time_date=?",devId,timedate).Find(&e1fragementflag)
       if e1fragementflag.E1fragementflag=="1"{
		num, sendMinFlag, sendMaxFlag, indexMax, indexMin, minDutyCycCoeff, maxDutyCycCoeff := FanParmquery(devId, dailyMonitoring.DevType,"e1")
		windblockage.Code = append(windblockage.Code, model.NonParaNmae[indexMax])
		windblockage.Code = append(windblockage.Code, model.NonParaNmae[indexMin])
		fmt.Println(num)
		if len(num)!=0{
			windblockage.Hvule = append(windblockage.Hvule, fmt.Sprintf("%02x",num[0]))
			windblockage.Hvule = append(windblockage.Hvule, fmt.Sprintf("%02x",num[1]))
			windblockage.Hvule = append(windblockage.Hvule, fmt.Sprintf("%02x",num[2]))
			windblockage.Hvule = append(windblockage.Hvule, fmt.Sprintf("%02x",num[3]))
		}

		windblockage.Pvule = append(windblockage.Pvule, strconv.Itoa(int(maxDutyCycCoeff)))
		windblockage.Pvule = append(windblockage.Pvule, strconv.Itoa(int(minDutyCycCoeff)))

        windblockage.Pvule16 = append(windblockage.Pvule16,fmt.Sprintf("%02x",maxDutyCycCoeff))
        windblockage.Pvule16 = append(windblockage.Pvule16, fmt.Sprintf("%02x",minDutyCycCoeff))
        windblockage.Hvule16 = append(windblockage.Hvule16, fmt.Sprintf("%02x",num[indexMax]))
        windblockage.Hvule16 = append(windblockage.Hvule16, fmt.Sprintf("%02x",(num[indexMin])))
        windblockage.Flagerror = 1
		windblockage.Flag = append(windblockage.Flag, sendMaxFlag)
		windblockage.Flag = append(windblockage.Flag, sendMinFlag)
		if sendMaxFlag==0&&sendMinFlag==0{
			windblockage.Flagerror = 0
		}
		
		windblockage.Index = append(windblockage.Index, indexMax)
		windblockage.Index = append(windblockage.Index, indexMin)
		fmt.Println("dayindeshi",windblockage)
	   }
//*******************************************************************************************************//
		row, over, starrow := QuerySetParameter0(devId, dailyMonitoring.DevType, "1")

		numper := 0

		if starrow[11] != "00"{
			htemper.Setflagerror11 = "1"
		}else{
			htemper.Setflagerror11 = "0"
		}
		fmt.Println("测试一下",len(over) )
		if len(over) != 0 { //关键参数出现问题
			for i, v := range row {

				if starrow[i] != v {
					htemper.Code = append(htemper.Code, over[numper])
					htemper.Pvule = append(htemper.Pvule, starrow[i])
					htemper.Hvule = append(htemper.Hvule, v)
					htemper.Flagerror = 1
					numper += 1
				}
			}

		} else if len(over) == 0 { //PL参数出现问题
			
			
			numflag, _ := strconv.ParseInt(numPerChange.ChangeNum, 10, 16)
			tempPL, _ := strconv.ParseInt(row[16], 16, 64)
			temp := tempPL + 5
			result := strconv.FormatInt(temp, 16)
			row[16] = result
			htemper.Code = append(htemper.Code, "PL")
			htemper.Pvule = append(htemper.Pvule, starrow[16])
			htemper.Hvule = append(htemper.Hvule, row[16])
			htemper.Flag = append(htemper.Flag, int(numflag))

			var e1num NumPerChanges
			db.Table("num_per_changes").Where("dev_id = ? AND fault_type= ?", devId, "E1").Find(&e1num)
			if e1num.ChangeNum=="1"{
				htemper.Flagerror = 2
			}
		}
		fmt.Println("下发值",row)
		xhrow.Nowrow = append(xhrow.Nowrow, row...)
		xhrow.Hourow = append(xhrow.Hourow, starrow...)
		xhrow.Perrow = append(xhrow.Perrow, over...)

	}
	if dailyMonitoring.C4 != 0 {

		num, sendMinFlag, sendMaxFlag, indexMax, indexMin, minDutyCycCoeff, maxDutyCycCoeff := FanParmquery(devId, dailyMonitoring.DevType,"c4")

		htemperc4.Code = append(htemperc4.Code, model.NonParaNmae[indexMax])
		htemperc4.Code = append(htemperc4.Code, model.NonParaNmae[indexMin])
		fmt.Println(num)
		if len(num)!=0{
			htemperc4.Hvule = append(htemperc4.Hvule, fmt.Sprintf("%02x",num[0]))
			htemperc4.Hvule = append(htemperc4.Hvule, fmt.Sprintf("%02x",num[1]))
			htemperc4.Hvule = append(htemperc4.Hvule, fmt.Sprintf("%02x",num[2]))
			htemperc4.Hvule = append(htemperc4.Hvule, fmt.Sprintf("%02x",num[3]))
		}

		htemperc4.Pvule = append(htemperc4.Pvule, strconv.Itoa(int(maxDutyCycCoeff)))
		htemperc4.Pvule = append(htemperc4.Pvule, strconv.Itoa(int(minDutyCycCoeff)))

        htemperc4.Pvule16 = append(htemperc4.Pvule16,fmt.Sprintf("%02x",maxDutyCycCoeff))
        htemperc4.Pvule16 = append(htemperc4.Pvule16, fmt.Sprintf("%02x",minDutyCycCoeff))
        htemperc4.Hvule16 = append(htemperc4.Hvule16, fmt.Sprintf("%02x",num[indexMax]))
        htemperc4.Hvule16 = append(htemperc4.Hvule16, fmt.Sprintf("%02x",(num[indexMin])))

		htemperc4.Flag = append(htemperc4.Flag, sendMaxFlag)
		htemperc4.Flag = append(htemperc4.Flag, sendMinFlag)
		htemperc4.Index = append(htemperc4.Index, indexMax)
		htemperc4.Index = append(htemperc4.Index, indexMin)
		fmt.Println("dayindeshi",htemperc4)
	}
//振荡调参
fmt.Println(dailyMonitoring.Zhendangflag)
if dailyMonitoring.Zhendangflag>0{
	fmt.Println("进来了")
	Wco , outcmt:= Volatilitydisplay(devId,devType)
	rowzhen.Nowrow = append(rowzhen.Nowrow, outcmt...) 
	
	htempzhen.Code = append(htempzhen.Code, "Wc")
	htempzhen.Code = append(htempzhen.Code, "W0")

	htempzhen.Pvule = append(htempzhen.Pvule, Wco[0])
	htempzhen.Pvule = append(htempzhen.Pvule, Wco[1])

	htempzhen.Hvule = append(htempzhen.Hvule, Wco[2])
	htempzhen.Hvule = append(htempzhen.Hvule, Wco[3])
    fmt.Println(Wco[0],Wco[2],Wco[1],Wco[3])

	var shocknum NumPerChanges
		db.Table("num_per_changes").Where("dev_id = ? AND fault_type= ?", devId, "shock").Find(&shocknum)
     fmt.Println("振荡次数",shocknum.ChangeNum)
	if shocknum.ChangeNum>="3" {
		htempzhen.Flag =append(htempzhen.Flag, 0)
	}else{
		htempzhen.Flag =append(htempzhen.Flag, 1)	
	}
	fmt.Println("666=",rowzhen)
}
	//恒温调参
	if dailyMonitoring.K!=0&&dailyMonitoring.Worst_temppattern != 0 && dailyMonitoring.Worst_temppattern != 11 && dailyMonitoring.Worst_score <= 80 {
		if !db.HasTable("param_history") {
			db.Table("param_history").AutoMigrate(&model.ParamHistory{})
		}
		//1、先判断整机关键参数是否超范围
		prow, overvaluep, s := QuerySetParameter0(dailyMonitoring.DevId, dailyMonitoring.DevType, "0") //查整机关键参数
		if s[11] != "00"{
			htemperheng.Setflagerror11 = "1"
		}else{
			htemperheng.Setflagerror11 = "0"
		}
		htrow.Nowrow = append(htrow.Nowrow, prow...)
		htrow.Hourow = append(htrow.Hourow, s...)
		htrow.Perrow = append(htrow.Perrow, overvaluep...)
		if len(overvaluep) != 0 && len(prow) != 1 { //改整机关键参数
			numper := 0
			for i, v := range prow {
				if s[i] != v {
					htemperheng.Code = append(htemperheng.Code, overvaluep[numper])
					htemperheng.Pvule = append(htemperheng.Pvule, s[i])
					htemperheng.Hvule = append(htemperheng.Hvule, v)
					numper += 1
				}
			}
		} else if len(overvaluep) == 0 && len(prow) != 1 {//此处需要判断一下热水器是否在燃烧，如果正在燃烧就没有xunyoude
			ka0, kb0, kc0, per := ParamModtempHquery(db, devId, dailyMonitoring.DevType, dailyMonitoring.Worst_start, dailyMonitoring.Worst_end, dailyMonitoring.F, dailyMonitoring.K)
			fmt.Println("返回", ka0, kb0, kc0)
			
			htemperheng.Code = append(htemperheng.Code, "ka")
			htemperheng.Code = append(htemperheng.Code, "kb")
			htemperheng.Code = append(htemperheng.Code, "kc")
			htemperheng.Pvule = append(htemperheng.Pvule, per...)
			htemperheng.Hvule = append(htemperheng.Hvule, fmt.Sprintf("%.4f", ka0))
			htemperheng.Hvule = append(htemperheng.Hvule, fmt.Sprintf("%.4f", kb0))
			htemperheng.Hvule = append(htemperheng.Hvule, fmt.Sprintf("%.4f", kc0))
            fmt.Println("恒温传递",htemperheng)
		}
	}
	response.Success(ctx, gin.H{"htemper": htemper, "windblockage":windblockage,"htemperc4": htemperc4, "htemperheng": htemperheng, 
	"xhrow": xhrow, "htrow": htrow ,"htempzhen":htempzhen , "rowzhen":rowzhen}, "恒温参数需要修改值")}

/*=========================================================
 * 函数名称：ChangenumWoWc
 * 功能描述: 调整参数次数
 * 输入参数: 设备号
 * 函数输出: -----
 * 作者   :  李贝贝
 * 时间   :  2022.07.11
 =========================================================*/
 func ChangenumWoWc(DevId string ,devtype string) (ret int) {
	
	db := common.GetIndexDB()
	var numPerChange NumPerChanges
	db.Table("num_per_changes").Where("fault_type=? and dev_id=?","shock",DevId).Find(&numPerChange)
   var ChangeNum =0
	ChangeNum,_=strconv.Atoi(numPerChange.ChangeNum)
	ChangeNum=ChangeNum+1
	fmt.Println("第几次修改wcwo",ChangeNum)
	var adjwcwonum NumPerChanges
	
	if ChangeNum==1{
		adjwcwonum.ChangeNum=strconv.Itoa(ChangeNum)
	   adjwcwonum.DataTime = time.Now().Format("2006-01-02 15:04:05")
       adjwcwonum.FaultType="shock"
	   adjwcwonum.DevId=DevId
	   db.Create(&adjwcwonum)
	}  else if ChangeNum>1{
		db.Table("num_per_changes").Where("dev_id = ? and fault_type = ?", DevId, "shock").
			Updates(map[string]interface{}{ 
				"data_time":           time.Now().Format("2006-01-02 15:04:05"),
				"change_num":          ChangeNum,
			})
	}
	if ChangeNum>3{
		/*需要上门设备输出*/
		var failedDevice FailedDevice
		failedDevice.DevId = DevId
		failedDevice.DevType = devtype
		failedDevice.FaultType = "shock"
		failedDevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedDevice.ReviseFlag = 1
		fmt.Println("需要上门修复",ChangeNum)
		db.Create(failedDevice)
	}
	return ChangeNum
}

/*=========================================================
 * 函数名称： Volatilitydisplay
 * 功能描述: 恒温震荡参数下发显示
 * 作    者：李贝贝
 * 时    间：2022-07-11
 * 修    改：李赞辉
 * 时    间：2022-07-12
 * 修改内容：该函数内不进行记录自加，只进行次数查询
// =========================================================*/
func Volatilitydisplay(DevId string,  devtype string)(ret [4]string ,outc []string){

	var backtotal [4]string
	var outcmd []string
	var flage int
	var numPerChange NumPerChanges
	common.IndexDB.Table("num_per_changes").Where("fault_type=? and dev_id=?","shock",DevId).Find(&numPerChange)
   var ChangeNum =0
	ChangeNum,_=strconv.Atoi(numPerChange.ChangeNum)
	ChangeNum++
	outcmd,flage = QuerySuanFaData0(DevId,"00")

	if flage ==0{
		for i:=0;i<10;i++ {
			outcmd,flage = QuerySuanFaData0(DevId,"00")
			if flage ==1{
				break
			}
		}		
	}else {
		backtotal[0] = outcmd[24]
		backtotal[1] = outcmd[25]
		fmt.Println("改前",outcmd[24],outcmd[25])
		wcstr, _ := strconv.ParseInt(outcmd[24], 16,32)
		wostr, _ := strconv.ParseInt(outcmd[25], 16,32)
		wcf := float64(wcstr) / float64(10)
		wof := float64(wostr)/ float64(10)
		if ChangeNum==1{

			wof=0.1
		}else if ChangeNum==2{
			wof=0.1
			wcf=1
		}else if ChangeNum==3{
			wof=0.1
			wcf=0.8
		}
		wostrchange := fmt.Sprintf("%02s", strconv.FormatInt(int64(wof*10),16)) 
		wcstrchange := fmt.Sprintf("%02s",  strconv.FormatInt(int64(wcf*10),16)) 

		outcmd[24] = wcstrchange
		outcmd[25] = wostrchange

		fmt.Println("改后",outcmd[24],outcmd[25])
		backtotal[2]= wcstrchange
		backtotal[3]= wostrchange
	}
	return backtotal ,outcmd
}
/*=========================================================
 * 函数名称： ReceiveillTemp
 * 功能描述: 恒温震荡参数修改
 * 作    者：李贝贝
 * 时    间：2022-07-11
 * 修    改：李赞辉
 * 时    间：2022-07-12
 * 修改内容：添加修改次数记录函数
// =========================================================*/
func ReceiveillTemp(devId string,devtype string,row []string , perback []string)(res int){
	//改写参数信息
	db := common.GetIndexDB()
	reSource := sendCommand(devId, 6, "00", row[14], row[15], row[16], row[17], row[18], row[19], row[20], row[21], row[22], row[23], row[24], row[25])
	reRow := strings.Split(reSource, ",")
	if len(reRow) != 48 {
		//返回值位数不对，进行查询
		secondSource := sendCommand(devId, 5, "00")
		secondRow := strings.Split(secondSource, ",")
		//若查询结果，位数依然错误，则直接返回0
		if len(secondRow) != 48 {
			fmt.Println("第二次读取恒温参数失败，无效数据")
			return 0
		} else if row[14] != secondRow[14] ||  row[15] != secondRow[15] || row[16] != secondRow[16] || row[17] != secondRow[17] || row[18] != secondRow[18] || row[19] != secondRow[19] ||
			row[20] != secondRow[20] || row[21] != secondRow[21] || row[22] != secondRow[22] || row[23] != secondRow[23] || row[24] != secondRow[24] || row[25] != secondRow[25] {
			fmt.Println("下发值与返回值不相等，改写失败")
			//改写失败,存储参数变化记录
			SeveRevisetemp(db,devId,row,perback,"0")
			return 0
		} else {
			//改写成功,存储参数变化记录
			SeveRevisetemp(db,devId,row,perback,"1")		
			return 1
		}
	} else if row[14] != reRow[14] || row[15] != reRow[15] || row[16] != reRow[16] || row[17] != reRow[17] || row[18] != reRow[18] || row[19] != reRow[19] ||
		row[20] != reRow[20] || row[21] != reRow[21] || row[22] != reRow[22] || row[23] != reRow[23] || row[24] != reRow[24] || row[25] != reRow[25] {
		fmt.Println("下发值与返回值不相等，改写失败")
		//改写失败,存储参数变化记录
		SeveRevisetemp(db,devId,row,perback,"0")
		return 0
	}
	//改写成功,存储参数变化记录
	SeveRevisetemp(db,devId,row,perback,"1")
	//改写成功，存储修改次数
	ChangenumWoWc(devId,devtype)
	return 1	
}
/*=========================================================
 * 函数名称： SeveRevisetemp
 * 功能描述: 恒温震荡参数存储
 * 作    者：李贝贝
 * 时    间：2022-07-11
// =========================================================*/
func SeveRevisetemp(db *gorm.DB,devId string ,row []string , perback []string,flag string){
	reviseTime := time.Now().Format("2006-01-02 15:04:05")
	var Parachange PerchangeNums
	Parachange.DevId = devId
	Parachange.FaultType = "恒温震荡修改"
	Parachange.UpdataTime = reviseTime
		//Wc参数存储
		Parachange.Pername = "Wc" 
		Parachange.UpValue = perback[0]
		Parachange.CurreValue = row[24]
		Parachange.SuccessFlag = flag
		db.Create(&Parachange)
		//Wo参数存储
		Parachange.Pername = "Wo" 
		Parachange.UpValue = perback[1]
		Parachange.CurreValue = row[25]
		Parachange.SuccessFlag = flag
		db.Create(&Parachange)
}


/*=========================================================
 * 函数名称： HxcqueryRewrint
 * 功能描述:  恒温参数、熄火、风机修需要修改的值
 * 功能说明:  1：改写熄火关键参数
             2：改写恒温关键参数
			 3：改写恒温参数
			 4：改写风机参数
			 5: 改写恒温震荡参数
 * 作    者：李贝贝
 * 修    改：李赞辉
 * 修改内容：新增整机参数修改标志位。0表示无需修改，1表示需要修改。整机参数修改无需被记录次数。记录次数函数写在这个真实下发的函数里
 =========================================================*/
 func HxcqueryRewrint(ctx *gin.Context) {
	Merrox := ctx.DefaultQuery("merrox", "0")   //标志位
	Merrohg := ctx.DefaultQuery("merrohg", "0") //标志位
	Merrogt := ctx.DefaultQuery("merrogt", "0") //标志位
	Merroc4 := ctx.DefaultQuery("merroc4", "0") //标志位
	Merrozhen := ctx.DefaultQuery("Merrozhen", "0") //标志位
	windblockage:=ctx.DefaultQuery("windblockage", "0") //标志位
	//改写熄火关键参数传输
	devId := ctx.DefaultQuery("devId", "183618442176081") //输入需要改动的设备号
	devtype := ctx.DefaultQuery("devtype", "5110173N")
	//熄火标志位
	zhengjiflag:= ctx.DefaultQuery("zhengjiflag", "0")
	xflag := ctx.DefaultQuery("xflag", "0") //--->>>>htemper.flag
	xrow := ctx.QueryArray("xrow[]")        //---->>>xhrow.Nowrow
	xover := ctx.QueryArray("xover[]")      //---->>>xhrow.Perrow
	xstarow := ctx.QueryArray("xstarow[]")  //--->>>xhrow.Hourow
	//恒温关键参数
	Hgrow := ctx.QueryArray("Hgrow[]")     //---->>>htrow.Nowrow
	Hgover := ctx.QueryArray("Hgover[]")   //---->>>htrow.Perrow
	Hgstrow := ctx.QueryArray("Hgstrow[]") //---->>>htrow.Hourow
	//恒温参数改写
	starttime := ctx.DefaultQuery("starttime", "2022-06-06") //日表中监控城市号
	endtime := ctx.DefaultQuery("endtime", "2022-06-06")     //修改的时间
	ka := ctx.DefaultQuery("ka", "1")                        //htemperheng.Hvule[0]
	kb := ctx.DefaultQuery("kb", "3")                        //htemperheng.Hvule[1]
	kc := ctx.DefaultQuery("kc", "1")                        //htemperheng.Hvule[2]
	//风机参数改写
	num := ctx.QueryArray("num[]") //htemperc4.Hvule
	sendMinFlag := ctx.DefaultQuery("sendMinFlag", "1")
	sendMaxFlag := ctx.DefaultQuery("sendMaxFlag", "3")
	indexMax := ctx.DefaultQuery("indexMax", "1")
	indexMin := ctx.DefaultQuery("indexMin", "1")
	minDutyCycCoeff := ctx.DefaultQuery("minDutyCycCoeff", "3")
	maxDutyCycCoeff := ctx.DefaultQuery("maxDutyCycCoeff", "1")
	//恒温震荡参数改写
	downsenet := ctx.QueryArray("downsenet[]")
	backout := ctx.QueryArray("backout[]")
	
    //fmt.Println("num",num)
	db := common.GetIndexDB()
	var num1 []int64
	for _, v := range num {
		temp, _ := strconv.ParseInt(v, 16, 64)
		num1 = append(num1, temp)
	}
	var flagerror string
	if Merrox == "1" {
		fmt.Println("前端下发的",xflag)
		secc := FlameoutChange(db, xflag, devId, devtype, xrow, xover, xstarow)
		if secc == "1" {
			if zhengjiflag=="0"{
              Changenum(devId)
			}
			
			flagerror = "1"
		} else {
			flagerror = flagerror + "0"
		}
	} 
	 if Merrohg == "1" {
		secc := ParamModtempSetchange(db, Hgrow, Hgover, Hgstrow, devId, devtype)
		if secc == "1" {
			flagerror = flagerror + "1"
		} else {
			flagerror = flagerror + "0"
		}
	} 
  if Merrogt == "1" {
		ka0, _ := strconv.ParseFloat(ka, 64)
		kb0, _ := strconv.ParseFloat(kb, 64)
		kc0, _ := strconv.ParseFloat(kc, 64)
		secc := ReviseTempParaInter(devId, starttime, endtime, devtype, ka0, kb0, kc0)
		if secc == "1" {
			flagerror = flagerror + "1"
		} else {
			flagerror = flagerror + "0"
		}
	} 



	if windblockage=="1"||Merroc4 == "1" {
		fmt.Println("num1",num1)
		sendMinFlagi, _ := strconv.Atoi(sendMinFlag)
		sendMaxFlagi, _ := strconv.Atoi(sendMaxFlag)
		indexMaxi, _ := strconv.Atoi(indexMax)
		indexMini, _ := strconv.Atoi(indexMin)
		minDutyCycCoeffi, _ := strconv.Atoi(minDutyCycCoeff)
		maxDutyCycCoeffi, _ := strconv.Atoi(maxDutyCycCoeff)
		secc := FanParmchange(db, num1, sendMinFlagi, sendMaxFlagi, indexMaxi, indexMini, devId, devtype, int64(minDutyCycCoeffi), int64(maxDutyCycCoeffi))
		if secc == 1 {
			flagerror = flagerror + "1"
		} else {
			flagerror = flagerror + "0"
		}
	}
	if Merrozhen == "1"{
		secc:=ReceiveillTemp(devId,devtype,downsenet,backout)
		if secc == 1 {
			flagerror = flagerror + "1"
		} else {
			flagerror = flagerror + "0"
		}

	}
	response.Success(ctx, gin.H{"flagerror": flagerror}, "参数修改")

}

func ParamModification(devId string, devType string, starttime string, endtime string, newF float64, newK float64) (kaz0, kbz0, Kcz0 float64) {
	db := common.GetIndexDB()
	if !db.HasTable("param_history") {
		db.Table("param_history").AutoMigrate(&model.ParamHistory{})
	}

	//1、先判断整机关键参数是否超范围
	var dailyMonitoring DailyMonitoring
	var kaz float64
	var kbz float64
	var kcz float64

	db.Table("daily_monitorings").Where("time_date between ? and ? and dev_id = ?", starttime, endtime, devId).Find(&dailyMonitoring)
	prow, overvaluep, s := QuerySetParameter0(dailyMonitoring.DevId, dailyMonitoring.DevType, "0") //查整机关键参数
	if len(overvaluep) != 0 && len(prow) != 1 {                                                    //改整机关键参数
		RewriteParameterSettingCmd0(dailyMonitoring.DevId, prow, s, "0a", devType, overvaluep) //改参数
		//改的次数记录下来，update方法更新表中
		var paraChangeTable NumPerChanges
		e1 := db.Table("num_per_changes").Where("dev_id = ?&&fault_type=?", devId, "整机关键参数超范围").Find(&paraChangeTable).Error
		if errors.Is(e1, gorm.ErrRecordNotFound) { //有错误，则第一次修改参数
			fmt.Println(e1)
			var paraChangeTable1 NumPerChanges
			paraChangeTable1.DevId = devId
			paraChangeTable1.ChangeNum = "1"
			paraChangeTable1.FaultType = "整机关键参数超范围"
			paraChangeTable1.DataTime = time.Now().Format("2006-01-02 15:04:05")
			db.Table("num_per_changes").Create(&paraChangeTable1)
		} else { //更新表
			currentCount, _ := strconv.Atoi(paraChangeTable.ChangeNum)
			db.Table("num_per_changes").Where("dev_id = ?&&fault_type=?", devId, "整机关键参数超范围").Updates(map[string]interface{}{
				"data_time":  time.Now().Format("2006-01-02 15:04:05"),
				"change_num": currentCount + 1, //原次数加一
			})
		}
	} else if len(overvaluep) == 0 && len(prow) != 1 { //修改恒温算法参数
		//2、对分段进行初始化数据
		var basicdata BasicData
		var Ka, Kb, Kc float64
		var paramHistoryV model.ParamHistory
		e := db.Table("param_history").Where("dev_id = ?", devId).Find(&paramHistoryV).Error
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//第一次查参数
			var paramhistoryV model.ParamHistory
			Ka, Kb, Kc = QueryParameter0("00", devId)
			fmt.Println("第0段", Ka, Kb, Kc)
			basicdata.firstSegment.basicVectorT, basicdata.firstSegment.basicVectorK = BasicDataInit(Ka, Kb, Kc)
			fmt.Println("6666666",basicdata.firstSegment.basicVectorT, basicdata.firstSegment.basicVectorK)
			paramhistoryV.DevId = devId
			paramhistoryV.DataTime = time.Now().Format("2006-01-02 15:04:05")
			paramhistoryV.Wa = Ka
			paramhistoryV.Wb = Kb
			paramhistoryV.Wc = Kc
			db.Table("param_history").Create(&paramhistoryV)
		} else {
			//有历史记录
			Ka = paramHistoryV.Wa
			Kb = paramHistoryV.Wb
			Kc = paramHistoryV.Wc
			basicdata.firstSegment.basicVectorT, basicdata.firstSegment.basicVectorK = BasicDataInit(Ka, Kb, Kc)
		}
		basicdata.firstSegment.basicVectorT = append(basicdata.firstSegment.basicVectorT, newF)
		basicdata.firstSegment.basicVectorK = append(basicdata.firstSegment.basicVectorK, newK)
		fmt.Println("进入寻优函数中数据长度TK：", len(basicdata.firstSegment.basicVectorT), len(basicdata.firstSegment.basicVectorK))
		ka0, kb0, kc0 := GASearch(basicdata.firstSegment.basicVectorK, basicdata.firstSegment.basicVectorT, len(basicdata.firstSegment.basicVectorK))
		fmt.Println("寻优出的参数：", ka0, kb0, kc0)
		//将寻优出的参数下发
		kaz = ka0
		kbz = kb0
		kcz = kc0

		ReviseTempParaInter(devId, starttime, endtime, devType, ka0, kb0, kc0)

		fmt.Println("恒温算法参数下发成功！")
	}
	return kaz, kbz, kcz
}

/*=========================================================
 * 函数名称： ParamModtempSetquery
 * 功能描述: 查询恒温参数的关键参数是否超范围
 * 作   者：
 * 创建日期： 2022.06.02
 =========================================================*/
func ParamModtempSetquery(db *gorm.DB, devId string, devType string, starttime string, endtime string) (prow, overvaluep, s []string) {

	if !db.HasTable("param_history") {
		db.Table("param_history").AutoMigrate(&model.ParamHistory{})
	}

	//1、先判断整机关键参数是否超范围
	var dailyMonitoring DailyMonitoring

	db.Table("daily_monitorings").Where("time_date between ? and ? and dev_id = ?", starttime, endtime, devId).Find(&dailyMonitoring)
	prow, overvaluep, s = QuerySetParameter0(dailyMonitoring.DevId, dailyMonitoring.DevType, "0") //查整机关键参数
	return prow, overvaluep, s
}

/*=========================================================
 * 函数名称： ParamModtempSetchange
 * 功能描述: 对恒温参数的关键参数做修改
 * 作   者：
 * 创建日期： 2022.06.02
 =========================================================*/
func ParamModtempSetchange(db *gorm.DB, prow, overvaluep, s []string, devId, devType string) (seccess string) {

	RewriteParameterSettingCmd0(devId, prow, s, "0a", devType, overvaluep) //改参数
	//改的次数记录下来，update方法更新表中
	var paraChangeTable NumPerChanges
	e1 := db.Table("num_per_changes").Where("dev_id = ?&&fault_type=?", devId, "整机关键参数超范围").Find(&paraChangeTable).Error
	if errors.Is(e1, gorm.ErrRecordNotFound) { //有错误，则第一次修改参数
		fmt.Println(e1)
		var paraChangeTable1 NumPerChanges
		paraChangeTable1.DevId = devId
		paraChangeTable1.ChangeNum = "1"
		paraChangeTable1.FaultType = "整机关键参数超范围"
		paraChangeTable1.DataTime = time.Now().Format("2006-01-02 15:04:05")
		db.Table("num_per_changes").Create(&paraChangeTable1)
	} else { //更新表
		currentCount, _ := strconv.Atoi(paraChangeTable.ChangeNum)
		db.Table("num_per_changes").Where("dev_id = ?&&fault_type=?", devId, "整机关键参数超范围").Updates(map[string]interface{}{
			"data_time":  time.Now().Format("2006-01-02 15:04:05"),
			"change_num": currentCount + 1, //原次数加一
		})
	}
	seccess = "1"
	return seccess

}

/*=========================================================
 * 函数名称： ParamModtempHquery
 * 功能描述: 查询恒温参数需要修改的值
 * 作   者：
 * 创建日期： 2022.06.02
 =========================================================*/
func ParamModtempHquery(db *gorm.DB, devId, devType string, starttime string, endtime string, newF float64, newK float64) (a, b, c float64, perk []string) {
	var basicdata BasicData
	var Ka, Kb, Kc float64
	var paramHistoryV model.ParamHistory
	e := db.Table("param_history").Where("dev_id = ?", devId).Find(&paramHistoryV).Error
	if errors.Is(e, gorm.ErrRecordNotFound) {
		//第一次查参数
		var paramhistoryV model.ParamHistory
		Ka, Kb, Kc = QueryParameter0("00", devId)

		perk = append(perk, strconv.FormatFloat(Ka, 'f', 2, 64)) //strconv.ParseFloat
		perk = append(perk, strconv.FormatFloat(Kb, 'f', 2, 64))
		perk = append(perk, strconv.FormatFloat(Kc, 'f', 2, 64))
		fmt.Println("第0段", Ka, Kb, Kc)
		basicdata.firstSegment.basicVectorT, basicdata.firstSegment.basicVectorK = BasicDataInit(Ka, Kb, Kc)
		paramhistoryV.DevId = devId
		paramhistoryV.DataTime = time.Now().Format("2006-01-02 15:04:05")
		paramhistoryV.Wa = Ka
		paramhistoryV.Wb = Kb
		paramhistoryV.Wc = Kc
		db.Table("param_history").Create(&paramhistoryV)
	} else {
		//有历史记录
		Ka = paramHistoryV.Wa
		Kb = paramHistoryV.Wb
		Kc = paramHistoryV.Wc
		perk = append(perk, strconv.FormatFloat(Ka, 'f', 2, 64)) //strconv.ParseFloat
		perk = append(perk, strconv.FormatFloat(Kb, 'f', 2, 64))
		perk = append(perk, strconv.FormatFloat(Kc, 'f', 2, 64))
		basicdata.firstSegment.basicVectorT, basicdata.firstSegment.basicVectorK = BasicDataInit(Ka, Kb, Kc)
	}
	fmt.Println("水流量",basicdata.firstSegment.basicVectorT)
	fmt.Println("K值",basicdata.firstSegment.basicVectorK)
	basicdata.firstSegment.basicVectorT = append(basicdata.firstSegment.basicVectorT, newF)
	basicdata.firstSegment.basicVectorK = append(basicdata.firstSegment.basicVectorK, newK)
	fmt.Println( "新的实际点",newF,newK)
	fmt.Println("进入寻优函数中数据长度TK：", len(basicdata.firstSegment.basicVectorT), len(basicdata.firstSegment.basicVectorK))
	ka0, kb0, kc0 := GASearch(basicdata.firstSegment.basicVectorK, basicdata.firstSegment.basicVectorT, len(basicdata.firstSegment.basicVectorK))
	fmt.Println("寻优出的参数：", ka0, kb0, kc0)
	//将寻优出的参数下发

	//ReviseTempParaInter(devId,starttime,endtime,devType,ka0,kb0,kc0)

	fmt.Println("恒温算法参数寻优成功！")
	return ka0, kb0, kc0, perk
}

func GASearch(k []float64, f []float64, length int) (float64, float64, float64) {
	var max = []float64{1, 1, 0.5}
	var min = []float64{0.1, 0.1, 0.001}
	var save savegood
	var saveg []savegood
	GeneticObject := NewGeneticAlgorithm(Start(3, max, min, "0"))
	for i := 0; i < 10; i++ {
		res, jid := GeneticObject.genetic(k, f, uint64(length))
		save.jd = jid[length-1]
		save.ka = res[0]
		save.kb = res[1]
		save.kc = res[2]
		saveg = append(saveg, save)
	}
	var minjid savegood
	for i := 1; i < len(saveg); i++ {

		if saveg[i].jd <= saveg[i-1].jd {
			minjid = saveg[i]
		}
	}
	fmt.Println(minjid)
	return minjid.ka, minjid.kb, minjid.kc
}
func BasicDataInit(ka float64, kb float64, kc float64) ([]float64, []float64) {
	var basicVectorT = []float64{3, 3.5, 4, 4.5, 5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5, 9, 9.5, 10, 10.5, 11, 11.5, 12}
	var basicVectorK []float64
	X := 0.0
	for i := 0; i < 19; i++ {
		X = ka*math.Exp(-kb*basicVectorT[i]) + kc
		basicVectorK = append(basicVectorK, X)
	}
	return basicVectorT, basicVectorK
}

/*=========================================================
 * 函数名称： ReviseTempParaInter
 * 功能描述: 远程修复恒温算法参数接口函数
 * 作   者： 宫健
 * 创建日期： 2022.06.02
 =========================================================*/
func ReviseTempParaInter(devId string, startTime string, endTime string, devType string, ka float64, kb float64, kc float64) (seccess string) {
	db := common.GetIndexDB()
	//1.整机关键参数未超范围，开始修改该设备恒温算法参数参数
	fmt.Println("开始改参")
	var paraChangeTable NumPerChanges
	e2 := db.Table("num_per_changes").Where("dev_id = ? && fault_type = ?", devId, "恒温参数修改").Find(&paraChangeTable).Error
	if errors.Is(e2, gorm.ErrRecordNotFound) { //第一次
		var youxiao int
		youxiao = ReadReviseSegmentAndPara(devId, startTime, endTime, devType, ka, kb, kc)
		if youxiao == 1 {
			var paraChangeTable1 NumPerChanges
			paraChangeTable1.DevId = devId
			paraChangeTable1.ChangeNum = "1"
			paraChangeTable1.FaultType = "恒温参数修改"
			paraChangeTable1.DataTime = time.Now().Format("2006-01-02 15:04:05")
			db.Table("num_per_changes").Create(&paraChangeTable1)
			seccess = "1"
			return seccess
		}
	} else {
		repairCountMax := 3
		changeNum, _ := strconv.Atoi(paraChangeTable.ChangeNum)
		if changeNum < repairCountMax {
			var youxiao int
			youxiao = ReadReviseSegmentAndPara(devId, startTime, endTime, devType, ka, kb, kc)
			if youxiao == 1 {
				db.Table("num_per_changes").Where("dev_id = ?&&fault_type=?", devId, "恒温参数修改").Updates(map[string]interface{}{
					"data_time":  time.Now().Format("2006-01-02 15:04:05"),
					"change_num": changeNum + 1, //原次数加一
				})
				seccess = "1"
				return seccess
			}
		} else {
			var failedevice FailedDevice
			er := db.Table("failed_devices").Where("dev_id = ? and fault_type = ?", devId, "恒温调参为解决，上门").Find(&failedevice).Error
			if errors.Is(er, gorm.ErrRecordNotFound) {
				var failedfevice FailedDevice
				failedfevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
				failedfevice.DevId = devId
				failedfevice.DevType = devType
				failedfevice.FaultType = "恒温调参为解决，上门"
				failedfevice.ReviseFlag = 1
				db.Table("failed_devices").Create(&failedfevice)
			} else {
				db.Table("failed_devices").Where("dev_id = ? and fault_type = ?", devId, "恒温调参为解决，上门").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":   time.Now().Format("2006-01-02 15:04:05"),
						"revise_flag": "1",
					})
			}
			seccess = "1"
			return seccess
		}

	}
	return
}

/*=========================================================
 * 函数名称： ReadReviseSegmentAndPara
 * 功能描述: 读取要修改的恒温算法参数
 * 作   者： 宫健
 * 创建日期： 2022.05.17
 =========================================================*/
func ReadReviseSegmentAndPara(devId string, startTime string, endTime string, devType string, ka float64, kb float64, kc float64) int {

	fmt.Println("恒温参数修改程序---------------------------------------------------------------------------------进入")

	//修改成功的负荷段数量
	var successNum int

	successNum = successNum + ReviseConstTempPara(devId, "00", ka, kb, kc)
	fmt.Println("第0段参数修改完成")

	successNum = successNum + ReviseConstTempPara(devId, "01", ka, kb, kc)
	fmt.Println("第1段参数修改完成")

	successNum = successNum + ReviseConstTempPara(devId, "02", ka, kb, kc)
	fmt.Println("第2段参数修改完成")

	successNum = successNum + ReviseConstTempPara(devId, "03", ka, kb, kc)
	fmt.Println("第3段参数修改完成")

	var youxiao int
	//是否所有负荷段全部改写成功
	if successNum == 4 {
		fmt.Println("4个负荷段全部改写成功")
		youxiao = 1
	} else {
		fmt.Println("4个负荷段没有全部改写成功，改写失败")
		youxiao = 0
	}
	return youxiao
}

/*=========================================================
 * 函数名称： ReviseConstTempPara
 * 功能描述: 修改恒温算法参数
 * 作   者： 宫健
 * 创建日期： 2022.05.17
 =========================================================*/
func ReviseConstTempPara(devId string, loadsegment string, ka float64, kb float64, kc float64) int {
	db := common.GetIndexDB()
	//参数进制、格式转换,去尾法，不知道是否精确
	kastr := strconv.FormatInt(int64(ka*100), 16)
	kbstr := strconv.FormatInt(int64(kb*100), 16)
	kcstr := strconv.FormatInt(int64(kc*1000), 16)
	if len(kastr) < 2 {
		kastr = "0" + kastr
	}
	if len(kbstr) < 2 {
		kbstr = "0" + kbstr
	}
	if len(kcstr) < 2 {
		kcstr = "0" + kcstr
	}

	//读取对应段的恒温算法参数
	source := sendCommand(devId, 5, loadsegment)
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("第一次读取恒温参数失败，无效数据")
		return 0
	}

	//改写参数信息
	reSource := sendCommand(devId, 6, loadsegment, kastr, kbstr, kcstr, row[17], row[18], row[19], row[20], row[21], row[22], row[23], row[24], row[25])
	reRow := strings.Split(reSource, ",")
	if len(reRow) != 48 {
		//返回值位数不对，进行查询
		secondSource := sendCommand(devId, 5, loadsegment)
		secondRow := strings.Split(secondSource, ",")
		//若查询结果，位数依然错误，则直接返回0
		if len(secondRow) != 48 {
			fmt.Println("第二次读取恒温参数失败，无效数据")
			return 0
		} else if kastr != secondRow[14] || kbstr != secondRow[15] || kcstr != secondRow[16] || row[17] != secondRow[17] || row[18] != secondRow[18] || row[19] != secondRow[19] ||
			row[20] != secondRow[20] || row[21] != secondRow[21] || row[22] != secondRow[22] || row[23] != secondRow[23] || row[24] != secondRow[24] || row[25] != secondRow[25] {
			fmt.Println("下发值与返回值不相等，改写失败")
			//改写失败,存储参数变化记录
			reviseTime := time.Now().Format("2006-01-02 15:04:05")
			var Parachange PerchangeNums
			Parachange.DevId = devId
			Parachange.FaultType = "恒温参数修改"
			Parachange.UpdataTime = reviseTime
			//ka参数存储
			Parachange.Pername = "ka" + loadsegment[1:2]
			Parachange.UpValue = row[14]
			Parachange.CurreValue = kastr
			Parachange.SuccessFlag = "0"
			db.Create(&Parachange)
			//kb参数存储
			Parachange.Pername = "kb" + loadsegment[1:2]
			Parachange.UpValue = row[15]
			Parachange.CurreValue = kbstr
			Parachange.SuccessFlag = "0"
			db.Create(&Parachange)
			//kc参数存储
			Parachange.Pername = "kc" + loadsegment[1:2]
			Parachange.UpValue = row[16]
			Parachange.CurreValue = kcstr
			Parachange.SuccessFlag = "0"
			db.Create(&Parachange)
			return 0
		} else {
			//改写成功,存储参数变化记录
			reviseTime := time.Now().Format("2006-01-02 15:04:05")
			var Parachange PerchangeNums
			Parachange.DevId = devId
			Parachange.FaultType = "恒温参数修改"
			Parachange.UpdataTime = reviseTime
			//ka参数存储
			Parachange.Pername = "ka" + loadsegment[1:2]
			Parachange.UpValue = row[14]
			Parachange.CurreValue = kastr
			Parachange.SuccessFlag = "1"
			db.Create(&Parachange)
			//kb参数存储
			Parachange.Pername = "kb" + loadsegment[1:2]
			Parachange.UpValue = row[15]
			Parachange.CurreValue = kbstr
			Parachange.SuccessFlag = "1"
			db.Create(&Parachange)
			//kc参数存储
			Parachange.Pername = "kc" + loadsegment[1:2]
			Parachange.UpValue = row[16]
			Parachange.CurreValue = kcstr
			Parachange.SuccessFlag = "1"
			db.Create(&Parachange)
			return 1
		}
	} else if kastr != reRow[14] || kbstr != reRow[15] || kcstr != reRow[16] || row[17] != reRow[17] || row[18] != reRow[18] || row[19] != reRow[19] ||
		row[20] != reRow[20] || row[21] != reRow[21] || row[22] != reRow[22] || row[23] != reRow[23] || row[24] != reRow[24] || row[25] != reRow[25] {
		fmt.Println("下发值与返回值不相等，改写失败")
		//改写失败,存储参数变化记录
		reviseTime := time.Now().Format("2006-01-02 15:04:05")
		var Parachange PerchangeNums
		Parachange.DevId = devId
		Parachange.FaultType = "恒温参数修改"
		Parachange.UpdataTime = reviseTime
		//ka参数存储
		Parachange.Pername = "ka" + loadsegment[1:2]
		Parachange.UpValue = row[14]
		Parachange.CurreValue = kastr
		Parachange.SuccessFlag = "0"
		db.Create(&Parachange)
		//kb参数存储
		Parachange.Pername = "kb" + loadsegment[1:2]
		Parachange.UpValue = row[15]
		Parachange.CurreValue = kbstr
		Parachange.SuccessFlag = "0"
		db.Create(&Parachange)
		//kc参数存储
		Parachange.Pername = "kc" + loadsegment[1:2]
		Parachange.UpValue = row[16]
		Parachange.CurreValue = kcstr
		Parachange.SuccessFlag = "0"
		db.Create(&Parachange)
		return 0
	}

	//改写成功,存储参数变化记录
	reviseTime := time.Now().Format("2006-01-02 15:04:05")
	var Parachange PerchangeNums
	Parachange.DevId = devId
	Parachange.FaultType = "恒温参数修改"
	Parachange.UpdataTime = reviseTime
	//ka参数存储
	Parachange.Pername = "ka" + loadsegment[1:2]
	Parachange.UpValue = row[14]
	Parachange.CurreValue = kastr
	Parachange.SuccessFlag = "1"
	db.Create(&Parachange)
	//kb参数存储
	Parachange.Pername = "kb" + loadsegment[1:2]
	Parachange.UpValue = row[15]
	Parachange.CurreValue = kbstr
	Parachange.SuccessFlag = "1"
	db.Create(&Parachange)
	//kc参数存储
	Parachange.Pername = "kc" + loadsegment[1:2]
	Parachange.UpValue = row[16]
	Parachange.CurreValue = kcstr
	Parachange.SuccessFlag = "1"
	db.Create(&Parachange)
	return 1
}

/*=========================================================
 * 函数名称：Stallparameters
 * 功能描述: 远程控制熄火参数调整
 * 输入参数: 设备号，开始时间，结束时间，设备型号
 * 函数输出: -----
 * 作者   :  李贝贝
 * 时间   :  2022.05.13
 =========================================================*/
func Stallparameters(DevId string, satrtime, endtime, dectype string) (PLs string) {
	var MideaFault Real_mideaFault
	db := common.GetIndexDB()
	var numPerChange NumPerChanges
	var pl string
	db.Table("num_per_changes").Where("dev_id = ? AND fault_type= ?", DevId, "E1").Find(&numPerChange)

	err := db.Table("e_midea_fault").Where("dev_id = ? and e1 != ? and start_time = ? and end_time= ?", DevId, 0, satrtime, endtime).Find(&MideaFault).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("没有找到")
	} else {
		if numPerChange.ChangeNum != "1" {
			row, over, starrow := QuerySetParameter0(DevId, dectype, "1")

			if len(over) != 0 {

				RewriteParameterSettingCmd0(DevId, row, starrow, "0a", dectype, over)
			} else if len(over) == 0 {
				pl = row[16]
				tempPL, _ := strconv.ParseInt(row[16], 16, 64)
				temp := tempPL + 5
				result := strconv.FormatInt(temp, 16)
				row[16] = result
				Changenum(DevId)

				RewriteParameterSettingCmd0(DevId, row, starrow, "0a", dectype, over)

			}
		} else {
			var failedevice FailedDevice
			er := db.Table("failed_devices").Where("dev_id = ? and fault_type = ?", DevId, "E1").Find(&failedevice).Error
			if errors.Is(er, gorm.ErrRecordNotFound) {
				var failedfevice FailedDevice
				failedfevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
				failedfevice.DevId = DevId
				failedfevice.DevType = dectype
				failedfevice.FaultType = "E1"
				failedfevice.ReviseFlag = 1
				db.Table("failed_devices").Create(&failedfevice)
			} else {
				db.Table("failed_devices").Where("dev_id = ? and fault_type = ?", DevId, "E1").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":   time.Now().Format("2006-01-02 15:04:05"),
						"revise_flag": "1",
					})
			}

		}

	}
	return pl

}

/*=========================================================
 * 函数名称：Flameoutquery
 * 功能描述: 远程控制熄火参数调整
 * 输入参数: 设备号，开始时间，结束时间，设备型号
 * 函数输出: -----
 * 时间   :  2022.05.13
 =========================================================*/
func Flameoutquery(db *gorm.DB, DevId string, satrtime, endtime, dectype string) (row []string, starrow []string, over []string, flag string) {

	var numPerChange NumPerChanges

	db.Table("num_per_changes").Where("dev_id = ? AND fault_type= ?", DevId, "E1").Find(&numPerChange)

	row, over, starrow = QuerySetParameter0(DevId, dectype, "1")

	return row, starrow, over, numPerChange.ChangeNum
}

/*=========================================================
 * 函数名称：FlameoutChange
 * 功能描述: 远程控制熄火参数调整
 * 输入参数: 设备号，开始时间，结束时间，设备型号
 * 函数输出: -----
 * 时间   :  2022.05.13
 * 修改：李赞辉
 * 修改内容：flag表示pl修改次数，pl修改0次时可以再次修改，已经修改大于0次则存储记录
 =========================================================*/
func FlameoutChange(db *gorm.DB, flag string, DevId, dectype string, row []string, over []string, starrow []string) (seccess string) {
	fmt.Println("复刻一下",flag)
	if flag =="0" {
		//row, over , starrow:= QuerySetParameter0(DevId, dectype, "1")
		if len(over) != 0 {
			RewriteParameterSettingCmd0(DevId, row, starrow, "0a", dectype, over)
			seccess = "1"
			return seccess
		} else if len(over) == 0 {

			RewriteParameterSettingCmd0(DevId, row, starrow, "0a", dectype, over)
			seccess = "1"
			return seccess
		}
	} else {
		var failedevice FailedDevice
		er := db.Table("failed_devices").Where("dev_id = ? and fault_type = ?", DevId, "E1").Find(&failedevice).Error
		if errors.Is(er, gorm.ErrRecordNotFound) {
			var failedfevice FailedDevice
			failedfevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
			failedfevice.DevId = DevId
			failedfevice.DevType = dectype
			failedfevice.FaultType = "E1"
			failedfevice.ReviseFlag = 1
			db.Table("failed_devices").Create(&failedfevice)
		} else {
			db.Table("failed_devices").Where("dev_id = ? and fault_type = ?", DevId, "E1").
				Updates(map[string]interface{}{ //有记录，更新记录
					"data_time":   time.Now().Format("2006-01-02 15:04:05"),
					"revise_flag": "1",
				})
		}
		seccess = "0"
		return seccess
	}
	return
}

/*=========================================================
 * 函数名称：Changenum
 * 功能描述: 调整参数次数
 * 输入参数: 设备号
 * 函数输出: -----
 * 作者   :  李贝贝
 * 时间   :  2022.05.13
 * 修改   :  李赞辉
 * 修改内容:  修复changenum无法自加的bug，该函数只用在熄火参数下发处
 =========================================================*/
func Changenum(DevId string) {
	var numPerChange NumPerChanges
	db := common.GetIndexDB()
	er := db.Table("num_per_changes").Where("dev_id = ? AND fault_type= ?", DevId, "E1").Find(&numPerChange).Error
	if errors.Is(er, gorm.ErrRecordNotFound) {
		var failedfevice NumPerChanges
		failedfevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedfevice.DevId = DevId
		failedfevice.FaultType = "E1"
		failedfevice.ChangeNum = "1"
		db.Table("num_per_changes").Create(&failedfevice)
	} else {
		changee1, _ := strconv.Atoi(numPerChange.ChangeNum)
		ChangeE1 := strconv.Itoa(changee1+1)
		db.Table("num_per_changes").Where("dev_id = ? and fault_type = ?", DevId, "E1").
			Updates(map[string]interface{}{ //有记录，更新记录
				"data_time":  time.Now().Format("2006-01-02 15:04:05"),
				"change_num": ChangeE1,
			})
	}
}

/*=========================================================
 * 函数名称：queryParameter
 * 功能描述: 参数查询（FA/FF）
 * 输入参数: 设备号，设备型号，恒温熄火关键参数标志位
 * 函数输出: 查询结果，修改结果，修改参数
 * 作者   :  李贝贝
 * 时间   :  2022.05.13
 =========================================================*/
func QuerySetParameter0(applianceid string, devtype string, fineflag string) (prow []string, overvaluep []string, starrow []string) {
	fmt.Println("查询的设备：", applianceid)

	var Keyparamater model.Keyparamaters
	db := common.GetIndexDB()
	db.Table("keyparamaters").Where("SN = ?", devtype).Find(&Keyparamater)

	// 发送指令
	var swatch  bool
	var outcmd string
	outcmd = sendCommand(applianceid, 9)
	strow := strings.Split(outcmd, ",")
	row := strings.Split(outcmd, ",")
	var per []string
	if len(outcmd) == 0 {//注意：第11位返回信息不是0，是2表示正在燃烧，查询到的参数都是错的
		swatch = true
		for i:=0;swatch == true;i++{
			outcmd = sendCommand(applianceid, 9)
			if len(outcmd) ==0{
				swatch = true
			}else{
				swatch = false
			}
		}
		return
	}
	tempFA, _ := strconv.ParseInt(row[12], 16, 64)
	tempFF, _ := strconv.ParseInt(row[13], 16, 64)
	tempPH, _ := strconv.ParseInt(row[14], 16, 64)
	tempFH, _ := strconv.ParseInt(row[15], 16, 64)
	tempPL, _ := strconv.ParseInt(row[16], 16, 64)
	tempFL, _ := strconv.ParseInt(row[17], 16, 64)
	tempDH, _ := strconv.ParseInt(row[18], 16, 64)
	tempFD, _ := strconv.ParseInt(row[19], 16, 64)
	tempCH, _ := strconv.ParseInt(row[20], 16, 64)
	tempFC, _ := strconv.ParseInt(row[21], 16, 64)

	tempFArange, _ := strconv.ParseInt(Keyparamater.FA, 16, 64)
	tempFFrange, _ := strconv.ParseInt(Keyparamater.FF, 16, 64)
	tempPHrange, _ := strconv.ParseInt(Keyparamater.PH, 16, 64)
	tempFHrange, _ := strconv.ParseInt(Keyparamater.FH, 16, 64)
	tempPLrange, _ := strconv.ParseInt(Keyparamater.PL, 16, 64)
	tempFLrange, _ := strconv.ParseInt(Keyparamater.FL, 16, 64)
	tempDHrange, _ := strconv.ParseInt(Keyparamater.DH, 16, 64)
	tempFDrange, _ := strconv.ParseInt(Keyparamater.FD, 16, 64)
	tempCHrange, _ := strconv.ParseInt(Keyparamater.CH, 16, 64)
	tempFCrange, _ := strconv.ParseInt(Keyparamater.FC, 16, 64)
	
	if fineflag == "0" {
		if tempFA != tempFArange {
			row[12] = Keyparamater.FA
			per = append(per, "FA")
		}
		if tempFF != tempFFrange {
			row[13] = Keyparamater.FF
			per = append(per, "FF")
		}
		if tempPH > (tempPHrange+10) || tempPH < (tempPHrange-10) {
			row[14] = Keyparamater.PH
			per = append(per, "PH")
		}

		if tempFH > (tempFHrange+6) || tempFH < (tempFHrange-6) {
			row[15] = Keyparamater.FH
			per = append(per, "FH")
		}

		if tempDH > (tempDHrange+10) || tempDH < (tempDHrange-10) {
			row[18] = Keyparamater.DH
			per = append(per, "DH")
		}
		if tempFD > (tempFDrange+6) || tempFD < (tempFDrange-6) {
			row[19] = Keyparamater.FD
			per = append(per, "FD")
		}
	}

	if tempPL > (tempPLrange+8) || tempPL < (tempPLrange-8) {
		row[16] = Keyparamater.PL
		per = append(per, "PL")
	}
	if tempFL > (tempFLrange+6) || tempFL < (tempFLrange-6) {
		row[17] = Keyparamater.FL
		per = append(per, "FL")
	}
	if tempCH > (tempCHrange+8) || tempCH < (tempCHrange-8) {
		row[20] = Keyparamater.CH
		per = append(per, "CH")
	}
	if tempFC > (tempFCrange+6) || tempFC < (tempFCrange-6) {
		row[21] = Keyparamater.FC
		per = append(per, "FC")
	}

	return row, per, strow

}

/*=========================================================
 * 函数名称： RewriteParameterSettingCmd
 * 功能描述: 改写参数设置（20）parameter setting接口
 * 输入参数: 设备号，上次查询结果，修改结果，改写选择标志位
 * 函数输出:  -----
 * 作者   :  李贝贝
 * 时间   :  2022.05.13
 =========================================================*/
func RewriteParameterSettingCmd0(applianceid string, rowsme []string, starrow []string, Script string, fault string, per []string) {
	FA := rowsme[12]
	FF := rowsme[13]
	PH := rowsme[14]
	FH := rowsme[15]
	PL := rowsme[16]
	FL := rowsme[17]
	dH := rowsme[18]
	Fd := rowsme[19]
	CH := rowsme[20]
	FC := rowsme[21]
	CA := rowsme[23]
	nE := rowsme[22]
	FP := rowsme[24]
	HS := rowsme[26]
	Hb := rowsme[27]
	HE := rowsme[28]
	HL := rowsme[29]
	HU := rowsme[30]
	Fn := rowsme[33]

	fmt.Println("总参数下发输出内容", applianceid, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn)

	gou := RewriteParSetting0(applianceid, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn, Script)
	//SaveChParameter(applianceid,starrow,gou,"0",fault)

	fmt.Println("总参数下发输出内容", gou)

	if Script == "0a" {
		if len(gou) != 48 {
			fmt.Println("无效数据")
			outcmd := sendCommand(applianceid, 9)
			row := strings.Split(outcmd, ",")
			fmt.Println("返回值长度", len(outcmd))
			if len(outcmd) == 0 {
				fmt.Println("无效数据2")
				SaveChParameter0(applianceid, starrow, row, "0", fault, per)
				return
			} else {
				fmt.Println("存入")
				SaveChParameter0(applianceid, starrow, row, "1", fault, per)
				return
			}
		} else {
			SaveChParameter0(applianceid, starrow, gou, "1", fault, per)
		}

	}

}

/*=========================================================
 * 函数名称： SaveChParameter
 * 功能描述: 改写参数设置（20）parameter setting接口
 * 输入参数: 设备号，上次查询结果，修改结果，成功失败把标志位
 * 函数输出:  -----
 * 作者   :  李贝贝
 * 时间   :  2022.05.13
 =========================================================*/
func SaveChParameter0(applianceId string, source []string, Nowrow []string, sunflag string, fault string, per []string) {
	if sunflag == "0" {
		source = sourcessssu
		Nowrow = sourcessssu
	}
	fmt.Println("参数改写成功记录")
	db := common.GetIndexDB()
	for i, v := range model.Setpar3 {
		for _, m := range per {
			if v == m {
				temp1 := PerchangeNums{
					DevId:       applianceId,
					Pername:     v,
					FaultType:   fault,
					UpValue:     source[i+12], //等等
					CurreValue:  Nowrow[i+12],
					SuccessFlag: "1",
					UpdataTime:  time.Now().Format("2006-01-02 15:04:05"),
				}
				db.Create(&temp1)
			} else {
				temp := PerchangeNums{
					DevId:       applianceId,
					Pername:     v,
					FaultType:   fault,
					UpValue:     source[i+12], //等等
					CurreValue:  Nowrow[i+12],
					SuccessFlag: "0",
					UpdataTime:  time.Now().Format("2006-01-02 15:04:05"),
				}
				db.Create(&temp)
			}

		}

	}

	//获取参数代码
}

/*=========================================================
 * 函数名称： RewriteParSetting
 * 功能描述: 发送指令，改写参数设置（20）parameter setting
 =========================================================*/
func RewriteParSetting0(applianceId string, fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string, scrit string) (ord []string) {
	var CodeLen int

	var outCmd string

	te, _ := strconv.ParseInt(scrit, 16, 64)
	fmt.Println("s:", te)
	var cmdsplits []string
	var cmdCodeOut string
	if te == 10 {
		cmdCodeOut = MechineCode[te] + "," + fa + "," + ff + "," + PH + "," + FH + "," + PL + "," + FL + "," + dH + "," + Fd + "," + CH + "," + FC + "," + CA + "," + nE + "," + FP + "," + HS + "," + Hb + "," + HE + "," + HL + "," + HU + "," + Fn
		cmdsplits = StandardizedMachine(cmdCodeOut)
		CodeLen = 48
		outCmd = cmdsplits[0]
		for num := 1; num < CodeLen; num++ {
			outCmd = StrAdd(outCmd, ",")
			outCmd = StrAdd(outCmd, cmdsplits[num])
		}
	} else if te == 11 {
		cmdCodeOut = MechineCode[te]
		cmdsplits = StandardizedMachineCode(cmdCodeOut)
		CodeLen = 31
		outCmd = cmdsplits[0]
		for num := 1; num < CodeLen; num++ {
			outCmd = StrAdd(outCmd, ",")
			outCmd = StrAdd(outCmd, cmdsplits[num])
		}
	}
	// cmdsplits[12] = groupNum

	fmt.Println("sendcmd:", outCmd)
	song := make(map[string]interface{})
	song["proType"] = "e3"
	song["applianceId"] = applianceId
	song["env"] = "prod"
	song["cmd"] = outCmd
	// fmt.Println("机器码", "AA,1f,E3,00,00,00,00,00,00,02,0E,05,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,e9")
	// fmt.Println("机器码", StandardizedMachineCode(MechineCode[1]))
	bytesData, err := json.Marshal(song) //将数据编码成json字符串
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	timeStartDiff := time.Now()
	reader := bytes.NewReader(bytesData) //将json改为byte格式，作为body传给http请求
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader) //创建url
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{} //客户端发起请求，接收返回值
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	timeEndDiff := time.Now()
	timeLength := timeEndDiff.Sub(timeStartDiff)       //两个时间相减，计算时间差
	fmt.Println("接受数据，共耗时", timeLength.Seconds(), "s") //输出时间长度，单位s
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*str)

	reply := Jsonyinhao(respBytes)

	ord = strings.Split(reply["reply"], ",")
	//存储到数据库

	return ord
}

func ParaSettingSave0(applianceId string, source string) (result model.ParamenSetting) {
	// 1.1 按逗号分割
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	// 1.2 将字符串转化为数字
	var num []int64
	var str []string
	for i := 12; i <= 33; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		if i == 22 {
			if temp == 2 {
				temp = 0
			} else if temp == 3 {
				temp = 1
			}
		}
		num = append(num, temp)
	}
	fmt.Println("修改后", num)
	// 1.3 将数字变为十进制字符串,得到参数值
	for _, v := range num {
		temp := strconv.FormatInt(v, 10)
		str = append(str, temp)
	}

	//	2.将参数值存入到表中
	db := common.GetRunDB()
	var code model.ParameterCodes
	var vaul model.ParameterFinalSettings
	for i, v := range model.ParaSetting {
		// 2.1 查询参数对应的编码
		db.Where("parameter = ?", v).First(&code)
		//2.2将相同编号的标志位清零
		db.Table("parameter_changes_settings").Where("appliance_id = ? AND code= ?", applianceId, code.Code).
			Update("latest_parameter_flag", "0")
		// 2.2 将值存入到变化表
		db.Table("parameter_final_settings").
			Where("appliance_id = ? AND code = ?", applianceId, code.Code).
			First(&vaul)
		if vaul.CurrentValue != str[i] {
			temp := model.ParameterChangesSettings{
				ApplianceId:         applianceId,
				Code:                code.Code,
				Value:               str[i],
				LastValue:           vaul.CurrentValue,
				LatestParameterFlag: "1",
				Updatetime:          time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&temp)

			// 2.3 将值存入到参数最终表
			err := db.Where("appliance_id = ? and code = ?", applianceId, code.Code).
				First(&model.ParameterFinalSettings{}).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				temp := model.ParameterFinalSettings{
					ApplianceId:        applianceId,
					Code:               code.Code,
					CurrentValue:       str[i],
					Updatetime:         time.Now().Format("2006-01-02 15:04:05"),
					RewriteSuccessFlag: "1",
				}
				db.Create(&temp) //没有找见记录，创建记录
			} else {
				db.Table("parameter_final_settings").
					Where("appliance_id = ? and code = ?", applianceId, code.Code).
					Updates(map[string]interface{}{ //有记录，更新记录
						"current_value":        str[i],
						"updatetime":           time.Now().Format("2006-01-02 15:04:05"),
						"rewrite_success_flag": "1",
					})
			}
		}
	}
	fmt.Println("获取", str)
	//3.返回信息
	result.FA = str[0]
	result.FF = str[1]
	result.PH = str[2]
	result.FH = str[3]
	result.PL = str[4]
	result.FL = str[5]
	result.DH = str[6]
	result.Fd = str[7]
	result.CH = str[8]
	result.FC = str[9]
	result.NE = str[10]
	result.CA = str[11]
	result.FP = str[12]
	result.LF = str[13]
	result.HS = str[14]
	result.Hb = str[15]
	result.HE = str[16]
	result.HL = str[17]
	result.HU = str[18]
	result.UA = str[19]
	result.Ub = str[20]
	result.Fn = str[21]
	result.ApplianceId = applianceId
	result.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("存取成功", result)
	return result
}

/*=========================================================
 * 函数名称：SaveSetParameter
 * 功能描述: 存储查询到的参数（FA/FF）到最终表
 =========================================================*/
func SaveSetParameter0(applianceId string, source string) {
	/***1.解析字符串***/
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	/***2.遍历参数***/
	for index, value := range model.Setpar2 {
		/***3.存储到最终表***/
		// 第12->33位
		SaveSigleParameter(applianceId, value, row[index+12])
	}
}

/*=========================================================
 * 函数名称： SaveSigleParameter
 * 功能描述:  将单个更改的参数存入参数最终表
 =========================================================*/
func SaveSigleParameter0(applianceId string, CodeName string, value string) (finalpara model.ParameterFinalSetting) {
	db := common.GetRunDB()
	//获取参数代码
	parametercode := model.ParameterCode{}
	db.Where("parameter = ?", CodeName).First(&parametercode)
	// 将16进制数转为10进制数
	num, _ := strconv.ParseInt(value, 16, 64)
	fmt.Println("数据为", num)
	if parametercode.Code == "030022" {
		if num == 2 {
			num = 0
		} else if num == 3 {
			num = 1
		}
	}
	str := strconv.FormatInt(num, 10)
	//存入数据
	finalpara.ApplianceId = applianceId
	finalpara.Code = parametercode.Code
	finalpara.CurrentValue = str
	finalpara.RewriteSuccessFlag = "1"
	finalpara.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	Result := db.Where("appliance_id = ? and code = ?", applianceId, parametercode.Code).
		First(&model.ParameterFinalSetting{})
	if errors.Is(Result.Error, gorm.ErrRecordNotFound) {
		println("没有找到")
		db.Create(&finalpara)
		fmt.Println("创建成功")
	} else {
		fmt.Println("已经存在")
		db.Table("parameter_final_settings").
			Where("appliance_id = ? and code = ?", applianceId, parametercode.Code).
			Updates(map[string]interface{}{
				"current_value":        str,
				"updatetime":           time.Now().Format("2006-01-02 15:04:05"),
				"rewrite_success_flag": "1",
			})
		fmt.Println("更新成功")
	}
	return
}

/*=========================================================
 * 函数名称：sendCommand
 * 功能描述: 发送指令 wifi->mcu
 * 参    数：applianceId：设备ID号，index：指令在MechineCode中的索引,parameter: 人为指定的参数
 * 返回参数：机器码
 =========================================================*/

func sendCommand0(applianceId string, index int, parameter ...string) (outcmd string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[index]
	for _, v := range parameter {
		cmdCodeOut = StrAdd(cmdCodeOut, ",")
		cmdCodeOut = StrAdd(cmdCodeOut, v)
	}
	cmdsplits := StandardizedMachineCode(cmdCodeOut)
	var outCmd string
	outCmd = cmdsplits[0]
	for num := 1; num < CodeLen; num++ {
		outCmd = StrAdd(outCmd, ",")
		outCmd = StrAdd(outCmd, cmdsplits[num])
	}

	fmt.Println("sendcmd:", outCmd)
	song := make(map[string]interface{})
	song["proType"] = "e3"
	song["applianceId"] = applianceId
	song["env"] = "prod"
	song["cmd"] = outCmd
	bytesData, err := json.Marshal(song) //将数据编码成json字符串
	if err != nil {
		fmt.Println(err.Error())
		return
	} //reader是指向结构体的指针
	reader := bytes.NewReader(bytesData) //将json改为byte格式，作为body传给http请求
	fmt.Println("reader:", reader)
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader) //创建url
	fmt.Println("request:", request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}         //客户端发起请求，接收返回值
	resp, err := client.Do(request) //resp也是结构体指针
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("返回内容：")
	fmt.Println(*str)
	/*******获取数据********/
	// 反序列化
	// var m map[string]interface{}
	// _ = json.Unmarshal(respBytes, &m)
	// // 再次反序列化
	// temp, _ := json.Marshal(m["result"])
	// var reply map[string]string
	// _ = json.Unmarshal(temp, &reply)
	reply := Jsonyinhao(respBytes)
	fmt.Println("机器码：")
	fmt.Println(reply["reply"])
	return reply["reply"]
}

/*=========================================================
 * 函数名称： StandardizedMachineCode
 * 功能描述: 机器码标准化
 =========================================================*/
func StandardizedMachineCode0(inCmd string) (outCmd []string) {
	CodeLen := 31 //机器码总位数
	outCmdSplits := strings.Split(inCmd, ",")
	for i := len(outCmdSplits); i < CodeLen; i++ { //位数补齐
		outCmdSplits = append(outCmdSplits, "00")
	}
	var outCmdSplitInts []int64
	for _, outCmdSplit := range outCmdSplits {
		outCmdSplitInt, _ := strconv.ParseInt(outCmdSplit, 16, 64)
		outCmdSplitInts = append(outCmdSplitInts, outCmdSplitInt)
	}
	var CheckSum int64 = 0
	for num := 1; num < CodeLen-1; num++ {
		CheckSum += outCmdSplitInts[num]
	}
	CheckSumComplement := 256 - CheckSum%256
	CheckSumComplementStr := strconv.FormatInt(CheckSumComplement, 16)
	CheckSumComplementStr1 := fmt.Sprintf("%02s", CheckSumComplementStr)
	outCmdSplits[CodeLen-1] = CheckSumComplementStr1 //填充校验位
	return outCmdSplits
}

/*=========================================================
 * 函数名称： StrAdd
 * 功能描述: 字符串拼接
 =========================================================*/
func StrAdd0(s1, s2 string) string { //字符串拼接函数
	s1 += s2
	return s1
}

/*=========================================================
 * 函数名称： StandardizedMachine
 * 功能描述: 机器码标准化
 =========================================================*/
func StandardizedMachine0(inCmd string) (outCmd []string) {
	CodeLen := 48 //机器码总位数
	outCmdSplits := strings.Split(inCmd, ",")
	for i := len(outCmdSplits); i < CodeLen; i++ { //位数补齐
		outCmdSplits = append(outCmdSplits, "00")
	}
	var outCmdSplitInts []int64
	for _, outCmdSplit := range outCmdSplits {
		outCmdSplitInt, _ := strconv.ParseInt(outCmdSplit, 16, 64)
		outCmdSplitInts = append(outCmdSplitInts, outCmdSplitInt)
	}
	var CheckSum int64 = 0
	for num := 1; num < CodeLen-1; num++ {
		CheckSum += outCmdSplitInts[num]
	}
	CheckSumComplement := 256 - CheckSum%256
	CheckSumComplementStr := strconv.FormatInt(CheckSumComplement, 16)
	outCmdSplits[CodeLen-1] = CheckSumComplementStr //填充校验位

	return outCmdSplits
}

func QueryParameter0(segment string, applianceid string) (ka float64, kb float64, kc float64) {

	outcmd1 := QuerySuanFaData(applianceid, segment) //查询组序号为00的参数

	kax, _ := strconv.ParseInt(outcmd1[0], 16, 32)
	kbx, _ := strconv.ParseInt(outcmd1[1], 16, 32)
	kcx, _ := strconv.ParseInt(outcmd1[2], 16, 32)

	kaxfloa := float64(kax) / float64(100)
	kbxfloa := float64(kbx) / float64(100)
	kcxfloa := float64(kcx) / float64(1000)

	return kaxfloa, kbxfloa, kcxfloa

}

/*=========================================================
 * 函数名称： QuerySuanFaData0
 * 功能描述: 使用协议向设备查询恒温算法相关数据
// =========================================================*/
func QuerySuanFaData0(applianceId string, groupNum string) (outcmd []string , flag int) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[5] + "," + groupNum
	cmdsplits := StandardizedMachineCode(cmdCodeOut)
	// cmdsplits[12] = groupNum
	var outCmd string
	outCmd = cmdsplits[0]
	for num := 1; num < CodeLen; num++ {
		outCmd = StrAdd(outCmd, ",")
		outCmd = StrAdd(outCmd, cmdsplits[num])
	}
	fmt.Println("sendcmd:", outCmd)
	song := make(map[string]interface{})
	song["proType"] = "e3"
	song["applianceId"] = applianceId
	song["env"] = "prod"
	song["cmd"] = outCmd
	// fmt.Println("机器码", "AA,1f,E3,00,00,00,00,00,00,02,0E,05,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,00,e9")
	// fmt.Println("机器码", StandardizedMachineCode(MechineCode[1]))
	bytesData, err := json.Marshal(song) //将数据编码成json字符串
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	reader := bytes.NewReader(bytesData) //将json改为byte格式，作为body传给http请求
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader) //创建url
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{} //客户端发起请求，接收返回值
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*str)

	m := Jsonyinhao(respBytes)
	// 再次反序列化
	//temp, _ := json.Marshal(m["result"])
	//	var reply map[string]string
	//	_ = json.Unmarshal(temp, &reply)
	fmt.Println("解析", m["reply"])

	strdatas := strings.Split(m["reply"], ",") //返回字符串类型切片
	if len(strdatas) != 48 {
		fmt.Println("设备返回无效数据,数据位数错误")
		return cmdsplits,0
	}

	fmt.Println(strdatas)
	return strdatas , 1
}
