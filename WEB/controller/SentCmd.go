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
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MechineCode = []string{
	"aa,1E,e3,00,00,00,00,00,00,02,0E,03", //0 一键启动风机自学习
	"AA,1E,E3,00,00,00,00,00,00,02,0E,05", //1 查询参数设置默认参数信息,12bit组序号
	"AA,1E,E3,00,00,00,00,00,00,02,0E,0A", //2 查询非参数设置参数信息,12bit组序号
	"AA,1E,E3,00,00,00,00,00,00,02,0E,55", //3 改写非调试模式参数第一组,12bit组序号
	"AA,1E,E3,00,00,00,00,00,00,02,0E,56", //4 改写非调试模式参数第二组
	"AA,1E,E3,00,00,00,00,00,00,02,0E,A0", //5 查询恒温算法相关数据,12bit段序号
	"AA,1E,E3,00,00,00,00,00,00,02,0E,A5", //6 改写恒温算法2.0相关数据
	"AA,1E,E3,00,00,00,00,00,00,02,0E,C5", //7 单个改写参数设置参数值
	"AA,1E,E3,00,00,00,00,00,00,02,0E,09", //8 一键恢复默认恒温参数
	"AA,1E,E3,00,00,00,00,00,00,02,08,03", //9 查询参数设置20个
	"AA,2F,E3,00,00,00,00,00,00,02,08,02", //10 改写参数设置20个
	"AA,1E,E3,00,00,00,00,00,00,02,08,01", //11 恢复出厂设置（参数设置参数）
	"AA,1E,E3,00,00,00,00,00,00,02,02,01", //12 关机请求指令
	"AA,1E,E3,00,00,00,00,00,00,02,01,01", //13 开机机请求指令
}

/******************************************************以下是一键恢复默认值接口*********************************/
/******************************************************以下是一键恢复默认值接口*********************************/
/******************************************************以下是一键恢复默认值接口*********************************/

/*=========================================================
 * 函数名称： OneKeyResConTempPara
 * 功能描述:  一键恢复恒温算法默认值
 =========================================================*/
func OneKeyResConTempParaCmd(ctx *gin.Context) {
	var sev [12]string
	applianceId := ctx.DefaultQuery("applianceid", "188016488514318")
	segment := ctx.DefaultQuery("segment", "00") //段序号
	sev[0] = ctx.DefaultQuery("ka", "00")        //模型参数Ka*100
	sev[1] = ctx.DefaultQuery("kb", "00")        //模型参数Kb*100
	sev[2] = ctx.DefaultQuery("kc", "00")        //模型参数Kc*1000
	sev[3] = ctx.DefaultQuery("kf", "00")        //bit0-6:模型参数Kf*10,bit7:正负(0：正1：负)
	sev[4] = ctx.DefaultQuery("T1a", "00")       //模型参数 T1a*1000
	sev[5] = ctx.DefaultQuery("T1c", "00")       //模型参数T1c*10
	sev[6] = ctx.DefaultQuery("T2a", "00")       //模型参数T2a*1000
	sev[7] = ctx.DefaultQuery("T2c", "00")       //模型参数T2c*100
	sev[8] = ctx.DefaultQuery("Tda", "00")       //模型参数Tda*1000
	sev[9] = ctx.DefaultQuery("Tdc", "00")       //模型参数Tdc*10
	sev[10] = ctx.DefaultQuery("Wc", "00")       //Wc*10
	sev[11] = ctx.DefaultQuery("Wo", "00")       //Wo*10
	outcmd := OneKeyResConTempPara(applianceId)
	if outcmd == "success" {
		outcmd1 := QuerySuanFaData(applianceId, segment)
		ReDataBase(outcmd1, applianceId, segment, sev)
		response.Success(ctx, gin.H{"outcmd": 1}, "恢复默认参数成功")
	} else if outcmd == "fail" {
		response.Success(ctx, gin.H{"errflag": 3}, "恢复默认参数失败")
	}
}

/*=========================================================
 * 函数名称： OneKeyResConTempPara
 * 功能描述: 使用协议向设备发送指令，恢复恒温算法默认值
// =========================================================*/
func OneKeyResConTempPara(applianceId string) (outcmd string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[8]
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
	reply :=  Jsonyinhao(respBytes)
	//解析数据
	// var m = make(map[string]interface{})

	// _ = json.Unmarshal(respBytes, &m)
	// fmt.Println("retcode", m["retCode"])
	// fmt.Println("desc", m["desc"])
	// fmt.Println("result", m["result"])
	// //fmt.Println(m["retCode"])
	// //fmt.Println(m["desc"])
	// //fmt.Println(m["result"])
	// temp, _ := json.Marshal(m["result"])
	// var reply = make(map[string]string)
	// _ = json.Unmarshal(temp, &reply)

	strdatas := strings.Split(reply["reply"], ",") //返回字符串类型切片
	out := strdatas[13:26]                         //截取有用数据,将段序号也截进来
	fmt.Println(out)

	if reply["desc"] == "success" {
		return "success"
	} else {
		return "fail"
	}
}
func ReDataBase(outcmd []string, applianceid string, segment string, shu [12]string) {
	db := common.GetDB()
	datas := make([]string, 13, 30) //定义切片
	for i := 0; i < len(outcmd); i++ {
		intdata, _ := strconv.ParseInt(outcmd[i], 16, 64) //转化为10进制
		datas[i] = strconv.FormatInt(intdata, 10)         //再转化为字符串，与数据表类型对应
	}
	fmt.Println(datas)
	live := make([]string, 13, 30)
	for i := 0; i < len(shu); i++ {
		intdata, _ := strconv.ParseInt(shu[i], 16, 64) //转化为10进制
		live[i] = strconv.FormatInt(intdata, 10)       //再转化为字符串，与数据表类型对应
	}
	var maps = make(map[string]string)
	maps["ApplianceId"] = applianceid
	var SegmentFlag int
	if segment == "00" {
		SegmentFlag = 0
	} else if segment == "01" {
		SegmentFlag = 1
	} else if segment == "02" {
		SegmentFlag = 2
	} else if segment == "03" {
		SegmentFlag = 3
	}
	j := 0
	for i := SegmentFlag * 12; i < (12 + SegmentFlag*12); i++ {
		//fmt.Println(i,v)
		var parametercode model.ParameterCode

		db.Table("parameter_codes").Where("Parameter = ?", model.Paranames[i]).First(&parametercode)
		code := parametercode.Code
		//写入最新修改的参数
		if datas[i-SegmentFlag*12] != live[j] {

			//写入最新修改的参数，将该设备的之前的该参数的最新参数标识位置0
			var ParaChangeRecordList model.ParameterChangesSetting
			err := db.Table("parameter_changes_settings").Where(" appliance_id = ? AND  code = ? AND latest_parameter_flag = ?", applianceid, code, "1").First(&ParaChangeRecordList).Error
			if err == nil {
				db.Table("parameter_changes_settings").Where(" appliance_id = ? AND  code = ? AND latest_parameter_flag = ?", applianceid, code, "1").Update("latest_parameter_flag", "0")
			}

			ParaChangeRecord := model.ParameterChangesSettings{
				ApplianceId:         applianceid,
				Code:                code,
				LastValue:           live[j],
				Value:               datas[i-SegmentFlag*12],
				Updatetime:          time.Now().Format("2006-01-02 15:04:05"),
				LatestParameterFlag: "1", //以true写入
			}
			db.Table("parameter_changes_settings").Create(&ParaChangeRecord)
		}
		//写入最终参数表
		var QuFinalPara model.ParameterFinalSetting
		err := db.Table("parameter_final_settings").Where("appliance_id = ? AND  code = ? ", applianceid, code).First(&QuFinalPara).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			FinalPara := model.ParameterFinalSetting{
				ApplianceId:        applianceid,
				Code:               code,
				CurrentValue:       datas[i-SegmentFlag*12],
				RewriteSuccessFlag: "1",
				Updatetime:         time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&FinalPara)
		} else {
			db.Table("parameter_final_settings").Where("appliance_id = ? AND  code = ? ", applianceid, code).Updates(map[string]interface{}{
				"current_value":        datas[i-SegmentFlag*12],
				"updatetime":           time.Now().Format("2006-01-02 15:04:05"),
				"rewrite_success_flag": "1",
			})
		}
		maps[model.Paranames[i]] = datas[i-SegmentFlag*12]
		j++
		if j == 12 {
			j = 0
		}
	}
	maps["Updatetime"] = time.Now().Format("2006-01-02 15:04:05")
}

/******************************************************以下是恢复出厂设置接口*********************************/
/******************************************************以下是恢复厂设置接口*********************************/
/******************************************************以下是恢复出厂设置接口*********************************/

/*=========================================================
 * 函数名称： HResConTempPara
 * 功能描述:  恢复出厂设置接口
 =========================================================*/
func HResConTempPara(ctx *gin.Context) {
	applianceId := ctx.DefaultQuery("applianceid", "188016488514318")

	outcmd := HuiResConTempPara(applianceId)
	if outcmd == "success" {
		response.Success(ctx, gin.H{"outcmd": 1}, "恢复默认参数成功")
	} else if outcmd == "fail" {
		response.Success(ctx, gin.H{"errflag": 3}, "恢复默认参数失败")
	}
}

/*=========================================================
 * 函数名称： HResConTempPara
 * 功能描述: 使用协议向设备发送指令，恢复恒温算法默认值
// =========================================================*/
func HuiResConTempPara(applianceId string) (outcmd string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[11]
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

	// //解析数据
	// var m = make(map[string]interface{})

	// _ = json.Unmarshal(respBytes, &m)
	// fmt.Println("retcode", m["retCode"])
	// fmt.Println("desc", m["desc"])
	// fmt.Println("result", m["result"])
	// temp, _ := json.Marshal(m["result\\"])
	// var reply map[string]string
	// _ = json.Unmarshal(temp, &reply)
	reply :=  Jsonyinhao(respBytes)
	fmt.Println(reply["reply"])
	ord := strings.Split(reply["reply"], ",")
	//存储到数据库
	res := SettingSave(applianceId, reply["reply"])
	fmt.Println(ord)
	fmt.Println(res)
	if reply["desc"] == "success" {
		return "success"
	} else {
		return "fail"
	}
}

func SettingSave(applianceId string, source string) (result model.ParamenSetting) {
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
		num = append(num, temp)
	}
	// 1.3 将数字变为十进制字符串,得到参数值
	for _, v := range num {
		temp := strconv.FormatInt(v, 10)
		str = append(str, temp)
	}

	//	2.将参数值存入到表中
	db := common.GetDB()
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
			// 2.3 将值存入到参数最终表
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
 * 函数名称：
 * 功能描述: 一键启动风机自学习
 =========================================================*/

/******************************************************以下是多个参数查询接口*********************************/
/******************************************************以下是多个参数查询接口*********************************/
/******************************************************以下是多个参数查询接口*********************************/

/*=========================================================
 * 函数名称：DefaultParaInformationCmd(ctx *gin.Context)
 * 功能描述: 查询参数设置默认参数信息接口
 =========================================================*/
func QueryParaSettingCmd(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	segment := ctx.DefaultQuery("segment", "00")
	fmt.Println("输出内容", applianceid, segment)

	assoucmd := FindNoParaFromDB(applianceid)
	response.Success(ctx, gin.H{"assoucmd": assoucmd}, "成功")
}

/*=========================================================
 * 函数名称：FindNoParaFromDB
 * 功能描述: 从 最终表 中返回 参数设置默认参数
 =========================================================*/
func FindNoParaFromDB(applianceId string) (res model.ParameterSettings) {
	var str = make([]string, 8, 8)
	db := common.InitDB()
	var code model.ParameterCodes
	var final model.ParameterFinalSettings
	for i, v := range model.Setpar {
		db.Where("parameter = ?", v).First(&code)
		db.Where("appliance_id = ? and code = ?", applianceId, code.Code).
			First(&final)
		str[i] = final.CurrentValue
	}
	res.PH0 = str[0]
	res.FH0 = str[1]
	res.PL0 = str[2]
	res.FL0 = str[3]
	res.DH0 = str[4]
	res.Fd0 = str[5]
	res.CH0 = str[6]
	res.FC0 = str[7]
	res.ApplianceId = applianceId
	res.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	return res
}

/*=========================================================
 * 函数名称： QueryNoParaSettingCmd
 * 功能描述: 查询非参数设置参数信息
 =========================================================*/
func QueryNoParaSettingCmd(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	//组序号
	groupNum := ctx.DefaultQuery("groupnum", "00")
	fmt.Println("输出内容", applianceid, groupNum)

	//从设备获取参数
	//outcmd := sendCommand(applianceid,2,groupNum)
	//res := NoParaSettingSave(applianceid, outcmd)
	//从参数最终表获取参数
	outcmd := QueryNoParaFromDB(applianceid)
	/*	ID:
		ApplianceId	---> applianceId
		最大负荷风机电流偏差系数:
		MaxCurrCoeff ---> 13
		最小负荷风机电流偏差系数:
		MinCurrCoeff ---> 14
		最大负荷风机占空比偏差系数:
		MaxDutyCycCoeff ---> 15
		最小负荷风机占空比偏差系数:
		MinDutyCycCoeff ---> 16
		回水水流值:
		BackwaterFlow ---> 17
		风压传感器报警点频率补偿值
		FreqWind ---> 18
		Updatetime
	*/
	//response.Success(ctx, gin.H{"outcmd": res}, "成功")
	response.Success(ctx, gin.H{"outcmd": outcmd}, "成功")

}

/*=========================================================
 * 函数名称： QueryNoParaFromDB
 * 功能描述:  从 参数最终表 中获取 非参数设置参数信息
 =========================================================*/
func QueryNoParaFromDB(applianceId string) (res model.NonParamenterSetting) {
	var str = make([]string, 6, 6)
	db := common.GetDB()
	var code model.ParameterCode
	var final model.ParameterFinalSetting
	for i, v := range model.NonParaNmae {
		// 1.查询参数对应的编码
		db.Where("parameter = ?", v).First(&code)
		// 2.从参数最终表中获取参数
		db.Where("appliance_id = ? and code = ?", applianceId, code.Code).
			First(&final)
		str[i] = final.CurrentValue
	}
	res.MaxCurrCoeff = str[0]
	res.MinCurrCoeff = str[1]
	res.MaxDutyCycCoeff = str[2]
	res.MinDutyCycCoeff = str[3]
	res.BackwaterFlow = str[4]
	res.FreqWind = str[5]
	res.ApplianceId = applianceId
	res.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	return res
}

/*=========================================================
 * 函数名称： QuerySuanFaDataCmd(ctx *gin.Context)
 * 功能描述: 在 最终表中 查询恒温算法相关参数
 =========================================================*/
func QuerySuanFaDataCmd(ctx *gin.Context) {
	db := common.GetDB()
	applianceid := ctx.DefaultQuery("applianceid", "179220395415410")
	BankMap := make(map[string]string)
	BankMap["ApplianceId"] = applianceid
	for _, v := range model.Paranames {
		var paracode model.ParameterCode
		var finalpara model.ParameterFinalSetting
		db.Table("parameter_codes").Where("parameter = ?", v).First(&paracode)
		db.Table("parameter_final_settings").Where("appliance_id = ? AND code = ?", applianceid, paracode.Code).First(&finalpara)
		BankMap[v] = finalpara.CurrentValue
		fmt.Println(BankMap[v])
	}
	response.Success(ctx, gin.H{"outcmd": BankMap}, "成功")
}

/*=========================================================
 * 函数名称：queryParameter
 * 功能描述: 参数查询（FA/FF）
 =========================================================*/
func QuerySetParameter(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	Index_Para_Find := ctx.DefaultQuery("Index_Para_Find", "4")
	fmt.Println("查询的设备：", applianceid)
	// 发送指令
	if Index_Para_Find == "4"{
		outcmd := sendCommand(applianceid, 9)
		// 存储查询到的参数（FA/FF）到最终表
		SaveSetParameter(applianceid, outcmd)
		var temp_msg string
		fmt.Println(outcmd[33:35])
      if outcmd == ""||outcmd[33:35]!="00"{
        temp_msg = "fail"
       } else if outcmd[33:35]=="00"{
      temp_msg = "ok"
      }
	  response.Success(ctx, gin.H{
           "outcmd":outcmd,
          "Info_Falg": temp_msg,
            }, "成功")
	}

}

/*=========================================================
 * 函数名称：SaveSetParameter
 * 功能描述: 存储查询到的参数（FA/FF）到最终表
 =========================================================*/
func SaveSetParameter(applianceId string, source string) {
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

/******************************************************以下是总查询接口*********************************/
/******************************************************以下是总查询接口*********************************/
/******************************************************以下是总查询接口*********************************/

/*=========================================================
 * 函数名称： QueryParameter
 * 功能描述:  总查询接口，查询设备的所有参数，并将结果存入默认参数表
 =========================================================*/
func QueryParameter(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	//applianceid := ctx.DefaultQuery("applianceid", "179220395415410")
	Index_Para_Find := ctx.DefaultQuery("Index_Para_Find", "1")

	//查询非参数设置参数信息
	if Index_Para_Find =="3"||Index_Para_Find =="2"{
		outcmd := sendCommand(applianceid, 2, "00")
		NoParaSettingOutcmd := NoParaSettingSave(applianceid, outcmd)
		var temp_msg string
     if NoParaSettingOutcmd.BackwaterFlow == "" && NoParaSettingOutcmd.MaxCurrCoeff == "" {
          temp_msg = "fail"
       } else {
            temp_msg = "ok"
           }
          response.Success(ctx, gin.H{
          "noparasettingoutcmd": NoParaSettingOutcmd, //返回非参数设置参数信息
          "Info_Falg":  temp_msg,
          }, "成功")

	}else if Index_Para_Find =="9"{
		//查询参数设置默认参数信息
		gou := QueryParaSetting(applianceid, "00")
		parasettingoutcmd := tablestor(gou, applianceid)	
		response.Success(ctx, gin.H{
			"parasettingoutcmd":   parasettingoutcmd,   //返回参数设置参数信息
			"Info_Falg":   "ok",    //返回恒温算法相关数据
		}, "成功")	
	}else if Index_Para_Find == "1"{
		//查询恒温算法相关数据
		segment := "00"
		outcmd1 := QuerySuanFaData(applianceid, segment) //查询组序号为00的参数
		segmentlist := []string{"01", "02", "03"}        //建立剩下三组组序号列表00，01，02
		//查询组序号为01，02，03的参数
		for _, v := range segmentlist {
			outcmd := QuerySuanFaData(applianceid, v) //返回字符串切片类型
			outcmd1 = append(outcmd1, outcmd...)      //切片拼接
		}
		suanfadataoutcmd := StoreDataBase(outcmd1, applianceid) //存储数据		
		var temp_msg string
         if suanfadataoutcmd == "" {
            temp_msg = "fail"
           } else {
            temp_msg = "ok"
           }
		
           response.Success(ctx, gin.H{
          "suanfadataoutcmd": suanfadataoutcmd, //返回恒温算法相关数据
           "Info_Falg":temp_msg,
        }, "成功")

	}
	//将全部数据返回
	// response.Success(ctx, gin.H{
	// 	"noparasettingoutcmd": NoParaSettingOutcmd, //返回非参数设置参数信息
	// 	"parasettingoutcmd":   parasettingoutcmd,   //返回参数设置参数信息
	// 	"suanfadataoutcmd":    suanfadataoutcmd,    //返回恒温算法相关数据
	// }, "成功")
}

/*=========================================================
 * 函数名称： NoParaSettingSave
 * 功能描述:  非参数设置参数信息 存储在 默认参数表、参数最终表 中
				第一次查完之后存到默认参数表和参数最终表里，之后再用总查询接口查就是只存到参数最终表里
 =========================================================*/
func NoParaSettingSave(applianceId string, source string) (result model.NonParamenterSetting) {
	//1.解析source字符串，获取参数
	/*	ID:
		ApplianceId	---> applianceId
		最大负荷风机电流偏差系数:
		MaxCurrCoeff ---> 13
		最小负荷风机电流偏差系数:
		MinCurrCoeff ---> 14
		最大负荷风机占空比偏差系数:
		MaxDutyCycCoeff ---> 15
		最小负荷风机占空比偏差系数:
		MinDutyCycCoeff ---> 16
		回水水流值:
		BackwaterFlow ---> 17
		风压传感器报警点频率补偿值
		FreqWind ---> 18
		Updatetime
	*/
	// 1.1 按逗号分割
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	// 1.2 将字符串转化为数字
	var num []int64
	var str []string
	for i := 13; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	// 1.3 将数字变为十进制字符串,得到参数值
	for _, v := range num {
		temp := strconv.FormatInt(v, 10)
		str = append(str, temp)
	}

	//	2.将参数值存入到表中
	db := common.GetDB()
	var code model.ParameterCode
	for i, v := range model.NonParaNmae {
		// 2.1 查询参数对应的编码
		db.Where("parameter = ?", v).First(&code)
		// 2.2 将值存入到默认参数表
		err := db.Where("appliance_id = ? and code = ?", applianceId, code.Code).
			First(&model.ParameterDefault{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			temp := model.ParameterDefault{
				ApplianceId:  applianceId,
				Code:         code.Code,
				DefaultValue: str[i],
				Updatetime:   time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&temp)
		}
		// 2.3 将值存入到参数最终表
		err = db.Where("appliance_id = ? and code = ?", applianceId, code.Code).
			First(&model.ParameterFinalSetting{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			temp := model.ParameterFinalSetting{
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
	//3.返回信息
	result.MaxCurrCoeff = str[0]
	result.MinCurrCoeff = str[1]
	result.MaxDutyCycCoeff = str[2]
	result.MinDutyCycCoeff = str[3]
	result.BackwaterFlow = str[4]
	result.FreqWind = str[5]
	result.ApplianceId = applianceId
	result.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(result)
	return result
}

/*=========================================================
 * 函数名称：QueryParaSetting
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/

func QueryParaSetting(applianceId string, groupNum string) (outbyte []string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[1] + "," + groupNum
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
	} //reader是指向结构体的指针
	reader := bytes.NewReader(bytesData) //将json改为byte格式，作为body传给http请求
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader) //创建url
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

	for i := 0; i < len(respBytes); i++ {
		if respBytes[i] == 92 {
			respBytes = append(respBytes[:i], respBytes[i+1:]...)
		}
	}
	if respBytes[0] == 34{
		respBytes =respBytes[1:len(respBytes)-2]
	}
	st := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*st)
	var m = make(map[string]string)	 
     json.Unmarshal(respBytes, &m)
	 fmt.Println("打印打印",m["reply"])
	tr := strings.Split(m["reply"], ",") //转为不带逗号的切片
	fmt.Println("tr", tr)
	return tr
}

/*=========================================================
 * 函数名称：tablestor
 * 功能描述: 查询参数设置默认参数信息存储在数据库中
 =========================================================*/
func tablestor(word []string, applianceId string) (bank string) {

	var back []int64
	var numstr []string

	for i := 13; i < 21; i++ {
		temp, _ := strconv.ParseInt(word[i], 16, 64) //转为10进制
		back = append(back, temp)
	}
	for _, u := range back {
		temp := strconv.FormatInt(u, 10)
		numstr = append(numstr, temp)
	}

	fmt.Println(numstr)
	db := common.InitDB()

	db.AutoMigrate(&model.ParameterDefaults{})
	use := model.ParameterDefaults{} //定义一个结构体
	for i := 0; i < 8; i++ {
		Result := db.Table("parameter_defaults").Where("code = ? and appliance_id = ? ", model.SetCode[i], applianceId).First(&model.ParameterDefaults{})
		if errors.Is(Result.Error, gorm.ErrRecordNotFound) {
			use.ApplianceId = applianceId
			use.Code = model.SetCode[i]
			use.DefaultValue = numstr[i]
			use.Updatetime = time.Now().Format("2006-01-02 15:04:05")
			db.Create(&use)
		} else {
			fmt.Println("默认表已经存在")
		}
	}
	//保存到最终表中
	db.AutoMigrate(&model.ParameterFinalSettings{})
	Parfis := model.ParameterFinalSettings{}
	//fl:=1
	for k := 0; k < 8; k++ {
		Re := db.Table("parameter_final_settings").Where("code = ? and appliance_id = ? ", model.SetCode[k], applianceId).First(&model.ParameterFinalSettings{})
		if errors.Is(Re.Error, gorm.ErrRecordNotFound) {
			Parfis.ApplianceId = applianceId
			Parfis.Code = model.SetCode[k]
			Parfis.CurrentValue = numstr[k]
			Parfis.Updatetime = time.Now().Format("2006-01-02 15:04:05")
			Parfis.RewriteSuccessFlag = "1"
			db.Create(&Parfis)
			fmt.Println("最终表创建成功")
		} else {
			fmt.Println("最终表已经存在")
			db.Table("parameter_final_settings").
				Where("appliance_id = ? and code = ?", applianceId, model.SetCode[k]).
				Update("current_value", numstr[k]) //有记录，更新记录
		}
	}
	findback := FindNoParaFromDB(applianceId)
	b, _ := json.Marshal(findback)
	return string(b)
}

/*=========================================================
 * 函数名称： QuerySuanFaData
 * 功能描述: 使用协议向设备查询恒温算法相关数据
// =========================================================*/
func QuerySuanFaData(applianceId string, groupNum string) (outcmd []string) {
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

	m:=Jsonyinhao(respBytes)
	// 再次反序列化
	//temp, _ := json.Marshal(m["result"])
	//	var reply map[string]string
	//	_ = json.Unmarshal(temp, &reply)
	fmt.Println("解析", m["reply"])

	strdatas := strings.Split(m["reply"], ",") //返回字符串类型切片
	fmt.Println(len(strdatas))
	if len(strdatas) != 48 {
		fmt.Println("设备返回无效数据,数据位数错误")
		return
	}
	outstrdatas := strdatas[14:26] //截取有用数据
	fmt.Println(outstrdatas)
	fmt.Println(outstrdatas[0],outstrdatas[1],outstrdatas[2])
	return outstrdatas
}

/*=========================================================
 * 函数名称： StoreDataBase
 * 功能描述: 存储数据到默认表及最终表中
=========================================================*/
func StoreDataBase(outcmd1 []string, applianceid string) string {
	db := common.GetDB()
	datas := make([]string, 50, 60) //定义切片
	for i := 0; i < len(outcmd1); i++ {
		intdata, _ := strconv.ParseInt(outcmd1[i], 16, 64) //转化为10进制
		datas[i] = strconv.FormatInt(intdata, 10)          //再转化为字符串，与数据表类型对应
	}
 

	//存储到数据表中
	for i, v := range model.Paranames {
		//fmt.Println(i,v)
		var parametercode model.ParameterCode
		db.Table("parameter_codes").Where("Parameter = ?", v).First(&parametercode)
		code := parametercode.Code
		//fmt.Println(code)
		//存储到默认值表中
		var qudefaultpara model.ParameterDefault
		err := db.Table("parameter_defaults").Where(" appliance_id = ? AND  code = ? ", applianceid, code).First(&qudefaultpara).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultpara := model.ParameterDefault{
				ApplianceId:  applianceid,
				Code:         code,
				DefaultValue: datas[i],
				Updatetime:   time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Table("parameter_defaults").Create(&defaultpara)
		}
		//写入最终参数表
		var QuFinalPara model.ParameterFinalSetting
		err = db.Table("parameter_final_settings").Where("appliance_id = ? AND  code = ? ", applianceid, code).First(&QuFinalPara).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			FinalPara := model.ParameterFinalSetting{
				ApplianceId:        applianceid,
				Code:               code,
				CurrentValue:       datas[i],
				RewriteSuccessFlag: "1",
				Updatetime:         time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&FinalPara)
		} else {
			db.Table("parameter_final_settings").Where("appliance_id = ? and code = ?", applianceid, code).
				Update(map[string]string{"current_value": datas[i],
					"updatetime":           time.Now().Format("2006-01-02 15:04:05"),
					"rewrite_success_flag": "1",
				}) //有记录，更新记录
		}
	}
	bankmap := make(map[string]string)
	bankmap["ApplianceId"] = applianceid
	for i, v := range model.Paranames {
		bankmap[v] = datas[i]
	}
	bankmap["Updatetime"] = time.Now().Format("2006-01-02 15:04:05")

	bank, _ := json.Marshal(&bankmap)
	 if datas[1] == "" {
      return ""
   } else {
       return string(bank)
       }
}

/******************************************************以下是改写接口*********************************/
/******************************************************以下是改写接口*********************************/
/******************************************************以下是改写接口*********************************/

/*=========================================================
 * 函数名称： RewriteNoDebugFirstCmd(ctx *gin.Context)
 * 功能描述: 改写非调试模式参数第一组接口
 =========================================================*/
func RewriteNoDebugFirstCmd(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")

	ReWaterFlow := ctx.DefaultQuery("rewaterflow", "00")               //回水水流值
	WindPressureSensor := ctx.DefaultQuery("windpressuresensor", "00") //风压传感器报警点补偿值
	fmt.Println("输出内容", applianceid, ReWaterFlow, WindPressureSensor)
	/************如果本次改写的值与上次的参数值相同，则不发送改写指令**************/
	db := common.GetDB()
	lastparms := make([]string, 2, 2)
	var code model.ParameterCodes
	var vaul model.ParameterFinalSettings
	for i, v := range model.FlyNonParaNmae {
		// 获取上一次的参数值
		db.Where("parameter = ?", v).First(&code)
		db.Table("parameter_final_settings").
			Where("appliance_id = ? and code = ?", applianceid, code.Code).First(&vaul)
		temp, _ := strconv.ParseInt(vaul.CurrentValue, 10, 64)
		str := strconv.FormatInt(temp, 16)
		lastparms[i] = fmt.Sprintf("%02s", str)
	}
	if ReWaterFlow != lastparms[0] || WindPressureSensor != lastparms[1] {
		source := RewriteNoDebugFirst(applianceid, ReWaterFlow, WindPressureSensor)
		row := strings.Split(source, ",")
		if len(row) != 48 || row[17] != ReWaterFlow || row[18] != WindPressureSensor  {
			fmt.Println("无效数据")
			outcmd := sendCommand(applianceid, 2, "00")
			row2 := strings.Split(outcmd, ",")
			if row2[17] != ReWaterFlow || row2[18] != WindPressureSensor {
				fmt.Println("改写失败")
				fmt.Println(row)
				response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
			}else {
				fmt.Println("改写成功")
				v, _ := strconv.ParseInt(ReWaterFlow, 16, 64)
				k, _ := strconv.ParseInt(WindPressureSensor, 16, 64)
				var feng [2]string
				feng[0] = strconv.FormatInt(v, 10)
				feng[1] = strconv.FormatInt(k, 10)
				//存储到数据库
				res := TwoNoParaSettingSave(applianceid, source, feng)
				fmt.Println(res)
				response.Success(ctx, gin.H{"outcmd": 1}, "成功")
			}
			//response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，无效数据")
			return
		} else {
			fmt.Println("改写成功")
			v, _ := strconv.ParseInt(ReWaterFlow, 16, 64)
			k, _ := strconv.ParseInt(WindPressureSensor, 16, 64)
			var feng [2]string
			feng[0] = strconv.FormatInt(v, 10)
			feng[1] = strconv.FormatInt(k, 10)
			//存储到数据库
			res := TwoNoParaSettingSave(applianceid, source, feng)
			fmt.Println(res)
			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
		}
	} else {
		response.Success(ctx, gin.H{"outcmd": 1}, "改写值与设备值相同，忽略")
	}
}

/*=========================================================
 * 函数名称： RewriteNoDebugFirstCmd(ctx *gin.Context)
 * 功能描述:向设备发送改写指令，改写非调试模式参数第一组
 =========================================================*/
func RewriteNoDebugFirst(applianceId string, ReWaterFlow string, WindPressureSensor string) (outcmd string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[3] + "," + ReWaterFlow + "," + WindPressureSensor
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
		return
	}
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*str)
	//获取数据
	// 反序列化
	//temp := strings.TrimLeft(*str, "\\")
	//jsonByteData := []byte(temp)
	//var m map[string]string
	m :=  Jsonyinhao(respBytes)
	//_ = json.Unmarshal(jsonByteData, &m)
	// 再次反序列化
	//temp, _ := json.Marshal(m["result"])
	//	var reply map[string]string
	//	_ = json.Unmarshal(temp, &reply)
	fmt.Println("解析", m["reply"])
	return m["reply"]
}

/*=========================================================
 * 函数名称： RewriteNoDebugFirstCmd(ctx *gin.Context)
 * 功能描述:存储非调试模式参数第一组参数到数据库
 =========================================================*/
func TwoNoParaSettingSave(applianceId string, source string, Feng [2]string) (result model.TowFlyParamenterSetting) {
	//1.解析source字符串，获取参数
	/*	ID:
		回水水流值:
		BackwaterFlow ---> 17
		风压传感器报警点频率补偿值
		FreqWind ---> 18
		Updatetime
	*/
	// 1.1 按逗号分割
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	// 1.2 将字符串转化为数字
	var num []int64
	var str []string
	for i := 17; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	// 1.3 将数字变为十进制字符串,得到参数值
	for _, v := range num {
		temp := strconv.FormatInt(v, 10)
		str = append(str, temp)
	}

	//	2.将参数值存入到表中
	db := common.GetDB()
	var code model.ParameterCodes
	var vaul model.ParameterFinalSettings
	for i, v := range model.FlyNonParaNmae {
		// 2.1 查询参数对应的编码
		db.Where("parameter = ?", v).First(&code)
		// 2.2 将值存入到变化表
		db.Table("parameter_final_settings").
			Where("appliance_id = ? AND code = ?", applianceId, code.Code).
			First(&vaul)
		if vaul.CurrentValue != Feng[i] { ///等等
			temp := model.ParameterChangesSettings{
				ApplianceId:         applianceId,
				Code:                code.Code,
				Value:               str[i],
				LastValue:           vaul.CurrentValue,
				LatestParameterFlag: "1",
				Updatetime:          time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&temp)
			db.Table("parameter_changes_settings").
				Where("appliance_id = ? and code = ? and updatetime <> ?", applianceId, temp.Code, temp.Updatetime).
				Update("latest_parameter_flag", "0")
		}
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
	//3.返回信息
	result.BackwaterFlow = str[0]
	result.FreqWind = str[1]
	result.ApplianceId = applianceId
	result.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("存取成功", result)
	return result
}

/*=========================================================
 * 函数名称： RewriteNoDebugSecondCmd(ctx *gin.Context)
 * 功能描述: 改写非调试模式参数第二组接口
 =========================================================*/
func RewriteNoDebugSecondCmd(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	//最大负荷风机电流偏差系数MaxCurrCoeff
	MaxCurrCoeff := ctx.DefaultQuery("maxcurrcoeff", "34")
	//最小负荷风机电流偏差系数MinCurrCoeff
	MinCurrCoeff := ctx.DefaultQuery("mincurrcoeff", "34")
	//最大负荷风机占空比偏差系数MaxDutyCycCoeff
	MaxDutyCycCoeff := ctx.DefaultQuery("maxdutycyccoeff", "34")
	//最小负荷风机占空比偏差系数MinDutyCycCoeff
	MinDutyCycCoeff := ctx.DefaultQuery("mindutycyccoeff", "34")
	fmt.Println("输出内容", applianceid, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff)
	/************如果本次改写的值与上次的参数值相同，则不发送改写指令**************/
	db := common.GetDB()
	lastparms := make([]string, 6, 6)
	var code model.ParameterCodes
	var vaul model.ParameterFinalSettings
	for i, v := range model.NonParaNmae {
		// 获取上一次的参数值
		db.Where("parameter = ?", v).First(&code)
		db.Table("parameter_final_settings").
			Where("appliance_id = ? and code = ?", applianceid, code.Code).First(&vaul)
		temp, _ := strconv.ParseInt(vaul.CurrentValue, 10, 64)
		str := strconv.FormatInt(temp, 16)
		lastparms[i] = fmt.Sprintf("%02s", str)
	}
	if lastparms[0] != MaxCurrCoeff || lastparms[1] != MinCurrCoeff || lastparms[2] != MaxDutyCycCoeff || lastparms[3] != MinDutyCycCoeff {
		//向设备发送命令，改写参数，并将结果存入参数最终表
		source := sendCommand(applianceid, 4, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff)
		row := strings.Split(source, ",")
		if len(row) != 48 ||row[13] != MaxCurrCoeff || row[14] != MinCurrCoeff || row[15] != MaxDutyCycCoeff || row[16] != MinDutyCycCoeff{
			fmt.Println("无效数据")
			outcmd := sendCommand(applianceid, 2, "00")
			row2 := strings.Split(outcmd, ",")
			if row2[13] != MaxCurrCoeff || row2[14] != MinCurrCoeff || row2[15] != MaxDutyCycCoeff || row2[16] != MinDutyCycCoeff {
				fmt.Println("改写失败")
				response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
			} else {
				fmt.Println("改写成功")
				//向参数变化表里添加记录
				UpdateSigleParameter(applianceid, model.NonParaNmae[0], MaxCurrCoeff)
				UpdateSigleParameter(applianceid, model.NonParaNmae[1], MinCurrCoeff)
				UpdateSigleParameter(applianceid, model.NonParaNmae[2], MaxDutyCycCoeff)
				UpdateSigleParameter(applianceid, model.NonParaNmae[3], MinDutyCycCoeff)
				//向参数最终表里添加记录
				outcmd := NoParaSettingSave(applianceid, source)
				fmt.Println(outcmd)
				response.Success(ctx, gin.H{"outcmd": 1}, "成功")
			}
			//response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，无效数据")
			return
		} else {
			fmt.Println("改写成功")
			//向参数变化表里添加记录
			UpdateSigleParameter(applianceid, model.NonParaNmae[0], MaxCurrCoeff)
			UpdateSigleParameter(applianceid, model.NonParaNmae[1], MinCurrCoeff)
			UpdateSigleParameter(applianceid, model.NonParaNmae[2], MaxDutyCycCoeff)
			UpdateSigleParameter(applianceid, model.NonParaNmae[3], MinDutyCycCoeff)
			//向参数最终表里添加记录
			outcmd := NoParaSettingSave(applianceid, source)
			fmt.Println(outcmd)
			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
		}
	} else {
		response.Success(ctx, gin.H{"outcmd": 1}, "改写值与设备值相同，忽略")
	}
}

/*=========================================================
 * 函数名称：sendCommand
 * 功能描述: 发送指令 wifi->mcu
 * 参    数：applianceId：设备ID号，index：指令在MechineCode中的索引,parameter: 人为指定的参数
 * 返回参数：机器码
 =========================================================*/
func sendCommand(applianceId string, index int, parameter ...string) (outcmd string) {
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
	reply:=Jsonyinhao(respBytes)
	fmt.Println("机器码：")
	fmt.Println(reply["reply"])
	return reply["reply"]
}

/*=========================================================
 * 函数名称： UpdateSigleParameter
 * 功能描述:  将单个更改的数据存入参数变化记录表
 =========================================================*/
func UpdateSigleParameter(applianceId string, CodeName string, value string) {
	var vaul model.ParameterFinalSettings
	db := common.GetDB()
	//获取参数代码
	parametercode := model.ParameterCode{}
	db.Where("parameter = ?", CodeName).First(&parametercode)
	db.Table("parameter_final_settings").Where("appliance_id = ? AND code = ?", applianceId, parametercode.Code).First(&vaul)
	// 将16进制数转为10进制数
	num, _ := strconv.ParseInt(value, 16, 64)
	str := strconv.FormatInt(num, 10)
	if str != vaul.CurrentValue { //判断当前值与传入的修改值是否一样，不一样保存变化表
		recode := model.ParameterChangesSettings{}
		recode.ApplianceId = applianceId
		recode.Code = parametercode.Code
		recode.Value = str
		recode.LastValue = vaul.CurrentValue
		recode.Updatetime = time.Now().Format("2006-01-02 15:04:05")
		recode.LatestParameterFlag = "1"
		//创建更新记录
		db.Create(&recode)
		//将该设备此参数之前为1的记录改为0
		db.Table("parameter_changes_settings").
			Where("appliance_id = ? and code = ? and updatetime <> ?", applianceId, parametercode.Code, recode.Updatetime).
			Update("latest_parameter_flag", "0")
		fmt.Println("更新成功")
	}
}

/*=========================================================
 * 函数名称： SaveSigleParameter
 * 功能描述:  将单个更改的参数存入参数最终表
 =========================================================*/
func SaveSigleParameter(applianceId string, CodeName string, value string) (finalpara model.ParameterFinalSetting) {
	db := common.GetDB()
	//获取参数代码
	parametercode := model.ParameterCode{}
	db.Where("parameter = ?", CodeName).First(&parametercode)
	// 将16进制数转为10进制数
	num, _ := strconv.ParseInt(value, 16, 64)
	fmt.Println("数据为",num)
	if parametercode.Code=="030022"{
		if num==2{
		num=0
	}else if num==3{
		num=1
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

func SaveChParameter(applianceId string, source string) (finalpara model.ParameterFinalSetting) {
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	// var datas []string
	// for i := 0; i < len(row); i++ {
	// 	intdata, _ := strconv.ParseInt(row[i], 16, 64) //转化为10进制
	// 	datas[i] = strconv.FormatInt(intdata, 10)      //再转化为字符串，与数据表类型对应
	// }
	db := common.GetDB()
	//获取参数代码
	parametercode := model.ParameterCode{}
	var vaul model.ParameterFinalSettings
	for i, v := range model.Setpar2 {
		// 2.1 查询参数对应的编码
		db.Where("parameter = ?", v).First(&parametercode)
		//2.2将相同编号的标志位清零
		db.Table("parameter_changes_settings").Where("appliance_id = ? AND code= ?", applianceId, parametercode.Code).
			Update("latest_parameter_flag", "0")
		// 2.2 将值存入到变化表
		db.Table("parameter_final_settings").
			Where("appliance_id = ? AND code = ?", applianceId, parametercode.Code).
			First(&vaul)
		temp_curret_int, _ := strconv.ParseInt(row[i+12], 16, 64)
		temp_curret_str := strconv.FormatInt(temp_curret_int, 10)
		if vaul.CurrentValue != temp_curret_str {
			temp := model.ParameterChangesSettings{
				ApplianceId:         applianceId,
				Code:                parametercode.Code,
				Value:               temp_curret_str, //等等
				LastValue:           vaul.CurrentValue,
				LatestParameterFlag: "1",
				Updatetime:          time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&temp)
		}
	}
	return
}

/*=========================================================
 * 函数名称： RewriteNoDebugSecond
 * 功能描述:  向设备发送指令，改写 非参数设置参数第二组
 =========================================================*/
func RewriteNoDebugSecond(applianceId, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff string) (res model.NonParamenterSetting, ord []string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[4] + "," + MaxCurrCoeff + "," + MinCurrCoeff + "," + MaxDutyCycCoeff + "," + MinDutyCycCoeff
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
	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	fmt.Println("reader:", reader)
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
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
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*str)
	//获取数据
	// 反序列化
	//var m map[string]interface{}
	//_ = json.Unmarshal(respBytes, &m)
	// 再次反序列化
	//temp, _ := json.Marshal(m["result"])
	//var reply map[string]string
	//_ = json.Unmarshal(temp, &reply)
	reply:=Jsonyinhao(respBytes)
	fmt.Println(reply["reply"])
	ord = strings.Split(reply["reply"], ",")
	//存储到最终表
	res = NoParaSettingSave(applianceId, reply["reply"])
	return res, ord
}

/*=========================================================
 * 函数名称： RewriteQueryParaSettingCmd(ctx *gin.Context)
 * 功能描述: 改写恒温算法2.0相关参数接口
 =========================================================*/
func RewriteQueryParaSettingCmd(ctx *gin.Context) {
	var sev [12]string
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	segment := ctx.DefaultQuery("segment", "00") //段序号
	ka := ctx.DefaultQuery("ka", "00")           //模型参数Ka*100
	kb := ctx.DefaultQuery("kb", "00")           //模型参数Kb*100
	kc := ctx.DefaultQuery("kc", "00")           //模型参数Kc*1000
	kf := ctx.DefaultQuery("kf", "00")           //bit0-6:模型参数Kf*10,bit7:正负(0：正1：负)
	T1a := ctx.DefaultQuery("T1a", "00")         //模型参数 T1a*1000
	T1c := ctx.DefaultQuery("T1c", "00")         //模型参数T1c*10
	T2a := ctx.DefaultQuery("T2a", "00")         //模型参数T2a*1000
	T2c := ctx.DefaultQuery("T2c", "00")         //模型参数T2c*100
	Tda := ctx.DefaultQuery("Tda", "00")         //模型参数Tda*1000
	Tdc := ctx.DefaultQuery("Tdc", "00")         //模型参数Tdc*10
	Wc := ctx.DefaultQuery("Wc", "00")           //Wc*10
	Wo := ctx.DefaultQuery("Wo", "00")           //Wo*10
	fmt.Println("输出内容", applianceid, segment, ka, kb, kc, kf, T1a, T1c, T2a, T2c, Tda, Tdc, Wc, Wo)
	sev[0] = ka
	sev[1] = kb
	sev[2] = kc
	sev[3] = kf
	sev[4] = T1a
	sev[5] = T1c
	sev[6] = T2a
	sev[7] = T2c
	sev[8] = Tda
	sev[9] = Tdc
	sev[10] = Wc
	sev[11] = Wo
	strdatas := RewriteQueryParaSetting(applianceid, segment, ka, kb, kc, kf, T1a, T1c, T2a, T2c, Tda, Tdc, Wc, Wo)
    outcmd := strdatas[13:26] //截取有用数据,将段序号也截进来	
    if len(strdatas) != 48 || ka != outcmd[1] || kb != outcmd[2] || kc != outcmd[3] || T1a != outcmd[5] || T1c != outcmd[6] || T2a != outcmd[7] || T2c != outcmd[8] || Tda != outcmd[9] || Tdc != outcmd[10] || Wc != outcmd[11] || Wo != outcmd[12]{
		fmt.Println("设备返回无效数据")
		outcmd1 := QuerySuanFaData(applianceid, segment)
		fmt.Println(outcmd1)
		ReDataBase(outcmd1, applianceid, segment, sev)
		if ka != outcmd1[0] || kb != outcmd1[1] || kc != outcmd1[2] || T1a != outcmd1[4] || T1c != outcmd1[5] || T2a != outcmd1[6] || T2c != outcmd1[7] || Tda != outcmd1[8] || Tdc != outcmd1[9] || Wc != outcmd1[10] || Wo != outcmd1[11] {
			response.Success(ctx, gin.H{"errflag": "3"}, "设备返回无效数据")
		} else {
			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
		}
	} else {
		fmt.Println("改写成功")
		//将改写的参数存入最终表
		bank := StoreReDataBase(outcmd, applianceid, segment) //存储数据
		fmt.Println(bank)
		response.Success(ctx, gin.H{"outcmd": 1}, "成功")
	}

}

/*=========================================================
 * 函数名称： RewriteQueryParaSetting
 * 功能描述: 发送改写指令，改写恒温算法2.0相关参数
 =========================================================*/
func RewriteQueryParaSetting(applianceId string, groupNum string, ka string, kb string, kc string, kf string, T1a string, T1c string, T2a string, T2c string, Tda string, Tdc string, Wc string, Wo string) (outcmd []string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[6] + "," + groupNum + "," + ka + "," + kb + "," + kc + "," + kf + "," + T1a + "," + T1c + "," + T2a + "," + T2c + "," + Tda + "," + Tdc + "," + Wc + "," + Wo
	cmdsplits := StandardizedMachineCode(cmdCodeOut) //机器码标准化
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
	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
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
	//byte数组直接转成string,优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println("输出内容")
	//fmt.Println(*str)
	//return *str

	//解析数据
	//var m = make(map[string]interface{})
	//_ = json.Unmarshal(respBytes, &m)
	//fmt.Println(m["retCode"])
	//fmt.Println(m["desc"])
	//fmt.Println(m["result"])
	//temp, _ := json.Marshal(m["result"])
	//var reply = make(map[string]string)
	//_ = json.Unmarshal(temp, &reply)
	reply:=Jsonyinhao(respBytes)
	strdatas := strings.Split(reply["reply"], ",") //返回字符串类型切片
	fmt.Println(strdatas)
	return strdatas
}

/*=========================================================
 * 函数名称： StoreReDataBase
 * 功能描述: 存储恒温算法2.0相关参数到数据库
// =========================================================*/
func StoreReDataBase(outcmd []string, applianceid string, segment string) string {
	db := common.GetDB()
	datas := make([]string, 13, 30) //定义切片
	for i := 0; i < len(outcmd); i++ {
		intdata, _ := strconv.ParseInt(outcmd[i], 16, 64) //转化为10进制
		datas[i] = strconv.FormatInt(intdata, 10)         //再转化为字符串，与数据表类型对应
	}
	fmt.Println(datas)

	var maps = make(map[string]string)
	maps["ApplianceId"] = applianceid
	var SegmentFlag int
	if segment == "00" {
		SegmentFlag = 0
	} else if segment == "01" {
		SegmentFlag = 1
	} else if segment == "02" {
		SegmentFlag = 2
	} else if segment == "03" {
		SegmentFlag = 3
	}
	for i := SegmentFlag * 12; i < (12 + SegmentFlag*12); i++ {
		//fmt.Println(i,v)
		var parametercode model.ParameterCode
		var vaul model.ParameterFinalSettings
		db.Table("parameter_codes").Where("Parameter = ?", model.Paranames[i]).First(&parametercode)
		code := parametercode.Code

		//写变化记录表
		//在最终表中查找当前参数的值
		err := db.Table("parameter_final_settings").
			Where("appliance_id = ? AND code = ?", applianceid, code).
			First(&vaul).Error
		if err == nil {

			//写入最新修改的参数
			if datas[i-SegmentFlag*12+1] != vaul.CurrentValue {

				//写入最新修改的参数，将该设备的之前的该参数的最新参数标识位置0
				var ParaChangeRecordList model.ParameterChangesSetting
				err := db.Table("parameter_changes_settings").Where(" appliance_id = ? AND  code = ? AND latest_parameter_flag = ?", applianceid, code, "1").First(&ParaChangeRecordList).Error
				if err == nil {
					db.Table("parameter_changes_settings").Where(" appliance_id = ? AND  code = ? AND latest_parameter_flag = ?", applianceid, code, "1").Update("latest_parameter_flag", "0")
				}

				ParaChangeRecord := model.ParameterChangesSettings{
					ApplianceId:         applianceid,
					Code:                code,
					LastValue:           vaul.CurrentValue,
					Value:               datas[i-SegmentFlag*12+1],
					Updatetime:          time.Now().Format("2006-01-02 15:04:05"),
					LatestParameterFlag: "1", //以true写入
				}
				db.Table("parameter_changes_settings").Create(&ParaChangeRecord)
			}
		} else {

			ParaChangeRecord := model.ParameterChangesSettings{
				ApplianceId:         applianceid,
				Code:                code,
				Value:               datas[i-SegmentFlag*12+1],
				Updatetime:          time.Now().Format("2006-01-02 15:04:05"),
				LatestParameterFlag: "1", //以true写入
			}
			db.Table("parameter_changes_settings").Create(&ParaChangeRecord)

		}
		//写入最终参数表
		var QuFinalPara model.ParameterFinalSetting
		err = db.Table("parameter_final_settings").Where("appliance_id = ? AND  code = ? ", applianceid, code).First(&QuFinalPara).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			FinalPara := model.ParameterFinalSetting{
				ApplianceId:        applianceid,
				Code:               code,
				CurrentValue:       datas[i-SegmentFlag*12+1],
				RewriteSuccessFlag: "1",
				Updatetime:         time.Now().Format("2006-01-02 15:04:05"),
			}
			db.Create(&FinalPara)
		} else {
			db.Table("parameter_final_settings").Where("appliance_id = ? AND  code = ? ", applianceid, code).Updates(map[string]interface{}{
				"current_value":        datas[i-SegmentFlag*12+1],
				"updatetime":           time.Now().Format("2006-01-02 15:04:05"),
				"rewrite_success_flag": "1",
			})
		}
		maps[model.Paranames[i]] = datas[i-SegmentFlag*12+1]
	}
	maps["Updatetime"] = time.Now().Format("2006-01-02 15:04:05")

	bank, _ := json.Marshal(&maps)
	return string(bank)
}

/*=========================================================
 * 函数名称： RewriteParameterSettingCmd(ctx *gin.Context)
 * 功能描述: 改写参数设置（20）parameter setting接口
 =========================================================*/
func RewriteParameterSettingCmd(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	Script := ctx.DefaultQuery("script", "0a")
	FA := ctx.DefaultQuery("FA", "00")
	FF := ctx.DefaultQuery("FF", "00")
	PH := ctx.DefaultQuery("PH", "00")
	FH := ctx.DefaultQuery("FH", "00")
	PL := ctx.DefaultQuery("PL", "00")
	FL := ctx.DefaultQuery("FL", "00")
	dH := ctx.DefaultQuery("dH", "00")
	Fd := ctx.DefaultQuery("Fd", "00")
	CH := ctx.DefaultQuery("CH", "00")
	FC := ctx.DefaultQuery("FC", "00")
	CA := ctx.DefaultQuery("CA", "00")
	nE := ctx.DefaultQuery("nE", "00")
	FP := ctx.DefaultQuery("FP", "00")
	HS := ctx.DefaultQuery("HS", "00")
	Hb := ctx.DefaultQuery("Hb", "00")
	HE := ctx.DefaultQuery("HE", "00")
	HL := ctx.DefaultQuery("HL", "00")
	HU := ctx.DefaultQuery("HU", "00")
	Fn := ctx.DefaultQuery("Fn", "00")

	fmt.Println("总参数下发输出内容", applianceid, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn)

	outcmd, gou := RewriteParSetting(applianceid, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn, Script)
	fmt.Println(outcmd)
	fmt.Println("总参数下发输出内容",gou)
	if len(gou) != 48 {
		fmt.Println("无效数据")
		outcmd1 := sendCommand(applianceid, 9)
		row := strings.Split(outcmd1, ",")

		if Script == "0b" {
			change := SaveChParameter(applianceid, outcmd1)
			fmt.Println(change)
			SaveSetParameter(applianceid, outcmd1)

			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
			return
		}

		if len(row) != 48 {
			fmt.Println("无效数据")
			return
		}
		// 删除了FC/FD的判断 引入FH FL的处理
		if row[12] != FA || row[13] != FF || row[14] != PH || row[15] != FH || row[16] != PL || row[17] != FL || row[18] != dH || row[20] != CH ||
			row[23] != CA || row[26] != HS || row[27] != Hb || row[28] != HE || row[29] != HL || row[30] != HU || row[33] != Fn {
			if (PH < PL && row[14] == row[16]) || (FH < FL && row[15] == row[17]) || (row[15] == "0" || row[17] == "0") {
				//fmt.Println("改写成功", gou[23])
				response.Success(ctx, gin.H{"outcmd": 1}, "成功")
			} else {
				fmt.Println("改写失败")
				response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
			}
		} else {
			change := SaveChParameter(applianceid, outcmd1)
			fmt.Println(change)
			SaveSetParameter(applianceid, outcmd1)

			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
			return
		}
	} //||gou[23]!=CA||gou[22]!=nE
	if Script == "0a" { //大坑
		// 删除了FC/FD的判断 引入FH FL的处理
		if gou[12] != FA || gou[13] != FF || gou[14] != PH || gou[15] != FH || gou[16] != PL || gou[17] != FL || gou[18] != dH || gou[20] != CH ||
			gou[23] != CA || gou[26] != HS || gou[27] != Hb || gou[28] != HE || gou[29] != HL || gou[30] != HU || gou[33] != Fn {
			if (PH < PL && gou[14] == gou[16]) || (FH < FL && gou[15] == gou[17]) || (gou[15] == "0" || gou[17] == "0") {
				fmt.Println("改写成功", gou[23])
				response.Success(ctx, gin.H{"outcmd": 1}, "成功")
			} else {
				fmt.Println("改写失败")
				//二次查询
				outcmd1 := sendCommand(applianceid, 9)
				row := strings.Split(outcmd1, ",")
				if row[12] != FA || row[13] != FF || row[14] != PH || row[15] != FH || row[16] != PL || row[17] != FL || row[18] != dH || row[20] != CH ||
					row[23] != CA || row[26] != HS || row[27] != Hb || row[28] != HE || row[29] != HL || row[30] != HU || row[33] != Fn {
					if (PH < PL && row[14] == row[16]) || (FH < FL && row[15] == row[17]) || (row[15] == "0" || row[17] == "0") {
						fmt.Println("改写成功", row[23])
						response.Success(ctx, gin.H{"outcmd": 1}, "成功")
					} else {
						fmt.Println("改写失败")
						response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
					}
				} else {
					fmt.Println("改写成功", row[23])
					response.Success(ctx, gin.H{"outcmd": 1}, "成功")
				}
				//response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
			}
		} else {
			fmt.Println("改写成功", gou[23])
			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
		}
	} else {

		//等等
		response.Success(ctx, gin.H{"outcmd": 1}, "恢复默认成功")
	}

}

/*=========================================================
 * 函数名称： RewriteParSetting
 * 功能描述: 发送指令，改写参数设置（20）parameter setting
 =========================================================*/
func RewriteParSetting(applianceId string, fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string, scrit string) (res model.ParamenSetting, ord []string) {
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
	//获取数据
	// 反序列化
//	var m map[string]interface{}
//	_ = json.Unmarshal(respBytes, &m)
	// 再次反序列化
	//temp, _ := json.Marshal(m["result"])
	//var reply map[string]string
	//_ = json.Unmarshal(temp, &reply)
	reply:=Jsonyinhao(respBytes)
	fmt.Println(reply["reply"])
	ord = strings.Split(reply["reply"], ",")
	//存储到数据库
	res = ParaSettingSave(applianceId, reply["reply"])
	return res, ord
}

/*=========================================================
 * 函数名称： RewriteParSetting
 * 功能描述: 存储参数设置（20）parameter setting到数据库
 =========================================================*/
func ParaSettingSave(applianceId string, source string) (result model.ParamenSetting) {
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
		if i==22{
			if temp==2{
				temp=0
			}else if temp==3{
				temp=1
			}
		}
		num = append(num, temp)
	}
	fmt.Println("修改后",num)
	// 1.3 将数字变为十进制字符串,得到参数值
	for _, v := range num {
		temp := strconv.FormatInt(v, 10)
		str = append(str, temp)
	}

	//	2.将参数值存入到表中
	db := common.GetDB()
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
 * 函数名称： RewriteSingleParaCmd
 * 功能描述:  单个改写  参数设置参数值
 =========================================================*/
func RewriteSingleParaCmd(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	index := ctx.Query("index") //改写序号
	value := ctx.Query("value") //改写值
	if index == "" || value == "" {
		fmt.Println("格式错误")
		response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
		return
	}
	fmt.Println("输出内容", applianceid, index, value)
	/************如果本次改写的值与上次的参数值相同，则不发送改写指令**************/
	db := common.GetDB()
	var lastparm string
	var code model.ParameterCodes
	var vaul model.ParameterFinalSettings
	// 获取上一次的参数值
	db.Where("parameter = ?", model.SingleParaIndex[index]).First(&code)
	db.Table("parameter_final_settings").
		Where("appliance_id = ? and code = ?", applianceid, code.Code).First(&vaul)
	temp, _ := strconv.ParseInt(vaul.CurrentValue, 10, 64)
	str := strconv.FormatInt(temp, 16)
	lastparm = fmt.Sprintf("%02s", str)
	if lastparm != value {
		// 2.发送改写指令
		result := sendCommand(applianceid, 7, index, value)
		//result := SendRewriteSinglePara(applianceid, index, value)
		// 3.数据校验并存入参数最终表
		// 3.1 按逗号分割
		row := strings.Split(result, ",")
		if len(row) != 48 {
			fmt.Println("无效数据")
			response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
			return
		}
		if row[12] != index || row[13] != value {
			fmt.Println("改写失败")
			fmt.Println(row[12])
			fmt.Println(row[13])
			response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
		} else {
			fmt.Println("改写成功")
			// 1.将改写的数据存入 参数变化记录表 中
			UpdateSigleParameter(applianceid, model.SingleParaIndex[index], value)
			//将改写的参数存入最终表
			res := SaveSigleParameter(applianceid, model.SingleParaIndex[index], value)
			fmt.Println(res)
			response.Success(ctx, gin.H{"outcmd": 1}, "成功")
		}
	} else {
		response.Success(ctx, gin.H{"outcmd": 1}, "改写值与设备值相同，忽略")
	}
}

/*=========================================================
 * 函数名称： SendRewriteSinglePara
 * 功能描述:  改写单个参数协议
 =========================================================*/
func SendRewriteSinglePara(applianceId string, index string, value string) (outcmd string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[7] + "," + index + "," + value
	cmdsplits := StandardizedMachineCode(cmdCodeOut) //机器码标准化
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
	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
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
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*str)
	//获取数据
	// 反序列化
	//var m map[string]interface{}
	//_ = json.Unmarshal(respBytes, &m)
	// 再次反序列化
	//temp, _ := json.Marshal(m["result"])
	//var reply map[string]string
	//_ = json.Unmarshal(temp, &reply)
	reply:=Jsonyinhao(respBytes)
	fmt.Println(reply["reply"])
	return reply["reply"]
}

/******************************************************以下是置0及删除变化表数据接口*********************************/
/******************************************************以下是置0及删除变化表数据接口*********************************/
/******************************************************以下是置0及删除变化表数据接口*********************************/

/*=========================================================
 * 函数名称： RewriteFind
 * 功能描述: 最终表归零接口
 =========================================================*/
func RewriteFind(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	outcmd := UpdateFindZone(applianceid)
	fmt.Println(outcmd)
	response.Success(ctx, gin.H{"outcmd": 1}, "成功")
}

/*=========================================================
 * 函数名称： UpdateFindZone
 * 功能描述: 最终表归零函数
 =========================================================*/
func UpdateFindZone(applianceId string) int64 {
	db := common.GetDB()
	var k int64
	Result := db.Where("appliance_id = ? ", applianceId).
		First(&model.ParameterFinalSetting{})
	if errors.Is(Result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("没有找到")
		k = 0
	} else {
		fmt.Println("已经存在")
		db.Table("parameter_final_settings").
			Where("appliance_id = ? ", applianceId).
			Updates(map[string]interface{}{
				"rewrite_success_flag": "0",
				"updatetime":           time.Now().Format("2006-01-02 15:04:05"),
			})
		fmt.Println("更新成功")
		k = 1
	}
	return k
}

/*=========================================================
 * 函数名称：DeleChangeParameter
 * 功能描述: 删除变化表中的记录接口
 * 创建日期：2022/1/7
 =========================================================*/
func DeleChangeParameter(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	changeCode := ctx.Query("changecode")
	DelTime := ctx.DefaultQuery("starttime", "2021-01-01 00:00:00")
	fmt.Println("删除的设备：", applianceid, changeCode, DelTime)
	// 发送指令
	outcmd := DeleChangeTable(applianceid, changeCode, DelTime)
	// 存储查询到的参数（FA/FF）到最终
	if outcmd == "5" {
		response.Success(ctx, gin.H{"outcmd": outcmd}, "删除成功")
	} else {
		response.Success(ctx, gin.H{"outcmd": outcmd}, "删除失败")
	}
}

/*=========================================================
 * 函数名称：DeleChangeParameter
 * 功能描述:删除变化表中数据函数
 * 创建日期：2022/1/7
 =========================================================*/
func DeleChangeTable(id string, gcode string, timedel string) (re string) {

	db := common.GetDB()

	//找到所给数据是否存在
	Resu := db.Where("appliance_id = ? and code = ? and updatetime = ?", id, gcode, timedel).
		First(&model.ParameterChangesSetting{})
	if errors.Is(Resu.Error, gorm.ErrRecordNotFound) {
		fmt.Println("数据不存在")
		re = "7"
	} else {
		//删除数据
		db.Where("appliance_id = ? and code = ? and updatetime = ?", id, gcode, timedel).Delete(&model.ParameterChangesSettings{})
	}
	//再次查找数据，判断是否删除成功
	Result := db.Where("appliance_id = ? and code = ? and updatetime = ?", id, gcode, timedel).
		First(&model.ParameterChangesSetting{})
	if errors.Is(Result.Error, gorm.ErrRecordNotFound) {
		if re != "7" {
			fmt.Println("删除成功")
			re = "5"
		}
	} else {
		fmt.Println("删除失败")
		re = "6"
	}
	return re
}

/******************************************************以下是开机关机接口函数*********************************/
/******************************************************以下是开机关机接口函数*********************************/
/******************************************************以下是开机关机接口函数*********************************/

/*=========================================================
 * 函数名称： Downdiv
 * 功能描述: 关机接口函数
// =========================================================*/
func Downdiv(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	fmt.Println("查询的设备：", applianceid)
	// 发送指令
	outcmd := FastartingdownData(applianceid, 12)
	// 存储查询到的参数（FA/FF）到最终表
	if len(outcmd) != 64 {
		fmt.Println("设备返回无效数据,数据位数错误")
		response.Success(ctx, gin.H{"errflag": 3}, "关机失败，格式错误")
		return
	}


	orderdown, _ := strconv.ParseInt(outcmd[13], 16, 64)
	orderdowntwo := strconv.FormatInt(orderdown, 2)
	if len(orderdowntwo) != 8 {
		for i := len(orderdowntwo); i < 8; i++ {
			orderdowntwo = "0" + orderdowntwo
		}
	}
	fmt.Println("返回值orderdowntwo", orderdowntwo)
	if orderdowntwo[7:] == "0" {
		fmt.Println("返回值12位中0位", orderdowntwo[7:])
		response.Success(ctx, gin.H{"outcmd": 1}, "成功")
	}

}

/*=========================================================
 * 函数名称： Downdiv
 * 功能描述: 开机接口函数
// =========================================================*/
func Straupdiv(ctx *gin.Context) {
	applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	fmt.Println("查询的设备：", applianceid)
	// 发送指令
	outcmd := FastartingdownData(applianceid, 13)
	// 存储查询到的参数（FA/FF）到最终表
	if len(outcmd) != 64 {
		fmt.Println("设备返回无效数据,数据位数错误")
		response.Success(ctx, gin.H{"errflag": 3}, "关机失败，格式错误")
		return
	}
	fmt.Println("返回值", outcmd)
	orderdown, _ := strconv.ParseInt(outcmd[12], 16, 64)
	orderdowntwo := strconv.FormatInt(orderdown, 2)
	
	if len(orderdowntwo) != 8 {
		for i := len(orderdowntwo); i < 8; i++ {
			orderdowntwo = "0" + orderdowntwo
		}
	}
	fmt.Println("返回值orderdowntwo", orderdowntwo)
	if orderdowntwo[7:] == "1" {
		fmt.Println("返回值12位中0位", orderdowntwo[7:])
		response.Success(ctx, gin.H{"outcmd": 1}, "成功")
	}

}

/*=========================================================
 * 函数名称： FastartingdownData
 * 功能描述: 开关机指令代码
// =========================================================*/
func FastartingdownData(applianceId string, order int) (outcmd []string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[order]
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
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println("输出内容")
	//fmt.Println(*str)
	//return *str

	//解析数据
	//var m = make(map[string]interface{})

//	_ = json.Unmarshal(respBytes, &m)
//	fmt.Println(m["retCode"])
//	fmt.Println(m["desc"])
//	fmt.Println(m["result"])

//	temp, _ := json.Marshal(m["result"])
//	var reply = make(map[string]string)
//	_ = json.Unmarshal(temp, &reply)
reply:=Jsonyinhao(respBytes)
	strdatas := strings.Split(reply["reply"], ",") //返回字符串类型切片
	// if len(strdatas) != 48 {
	// 	fmt.Println("设备返回无效数据,数据位数错误")
	// 	return
	// }
	// outstrdatas := strdatas[14:26] //截取有用数据
	// fmt.Println(outstrdatas)
	return strdatas
}

/*=========================================================
 * 函数名称： StandardizedMachineCode
 * 功能描述: 机器码标准化
 =========================================================*/
func StandardizedMachineCode(inCmd string) (outCmd []string) {
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
func StrAdd(s1, s2 string) string { //字符串拼接函数
	s1 += s2
	return s1
}

/*=========================================================
 * 函数名称： StandardizedMachine
 * 功能描述: 机器码标准化
 =========================================================*/
func StandardizedMachine(inCmd string) (outCmd []string) {
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


func Jsonyinhao(respBytes []byte) map[string]string{
	for i := 0; i < len(respBytes); i++ {
		if respBytes[i] == 92 {
			respBytes = append(respBytes[:i], respBytes[i+1:]...)
		}
	}
   
		if respBytes[0] == 34 {
			respBytes =respBytes[1:len(respBytes)-2]
		}
        
	
	
		//respBytes =respBytes[1:len(respBytes)-1]
	
	var m = make(map[string]string)	 
     json.Unmarshal(respBytes, &m)
	 return m
}
