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
	_ "runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
)

type MultipleParaRewrite struct {
	ApplianceId string `json:"appliance_id" gorm:"column:appliance_id"`
	HandleFlag  string ` json:"handleflag" gorm:"column:handleflag"`
	SucceedFlag string ` json:"succeedflag" gorm:"column:succeedflag"`
	RewriteFlag string ` json:"rewriteflag" gorm:"column:rewriteflag"`
}

/*=========================================================
 * 函数名称：QuerySetParametertwenty
 * 功能描述: 参数查询（FA/FF）
 =========================================================*/
func QuerySetParametertwenty(ctx *gin.Context) {
	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite

	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)

	fmt.Println("输出内容", MultipleParaRewritesu)
	TaskAssignmenttwenty(7, MultipleParaRewritesu, db)

}

/*=========================================================
 * 函数名称： RewriteSingleParaCmd
 * 功能描述:  批量单个改写  参数设置参数值
 =========================================================*/
func BatchRewriteSingleParaCmd(ctx *gin.Context) {

	index := ctx.Query("index") //改写序号HS---01
	value := ctx.Query("value") //改写值
	SingleMin := ctx.Query("singlemin")
	SingleMax := ctx.Query("singlemax")

	fmt.Println("下线：", SingleMin)
	fmt.Println("上线：", SingleMax)

	if index == "" || value == "" {
		fmt.Println("格式错误")
		response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
		return
	}
	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite

	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)

	fmt.Println("输出内容", MultipleParaRewritesu)

	TaskAssignment(7, MultipleParaRewritesu, index, value)

}

/*=========================================================
 * 函数名称： BatchRewriteSumParaCmd
 * 功能描述:  批量改写总参数（只更改一个参数）
 * 目    的:  在设备中没有单个改写参数的协议
 =========================================================*/
func BatchRewriteSumParaCmd(ctx *gin.Context) {

	index := ctx.Query("index") //改写序号HS---01
	value := ctx.Query("value") //改写值
	SingleMin := ctx.Query("singlemin")
	SingleMax := ctx.Query("singlemax")
	Starent := ctx.Query("Starent") //backups
	Backups := ctx.Query("Backups")
	fmt.Println("下线：", SingleMin)
	fmt.Println("上线：", SingleMax)
	if index == "" || value == "" {
		fmt.Println("格式错误")
		response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
		return
	}
	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite
	var pare model.ParameterSerials
	var indexsrevu int
	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)
	db.Table("parameter_serials").Where("serial_number = ?", index).Find(&pare)
	fmt.Println("参数值：", pare.Parameter)
	for ind, val := range model.ParaSetting {
		if pare.Parameter == val {
			indexsrevu = ind + 1
		}
	}
	fmt.Println("对应的序号:", indexsrevu)

	fmt.Println("输出内容", MultipleParaRewritesu)
	if Starent == "1" {
		TaskAssignmentsumStarent(7, MultipleParaRewritesu, indexsrevu, value, index, Backups)
	} else if Starent == "0" {
		TaskAssignmentsum(7, MultipleParaRewritesu, indexsrevu, value, index, Backups)
	}

}

/*=========================================================
 * 函数名称： RewriteQueryParaSettingCmd(ctx *gin.Context)
 * 功能描述: 改写恒温算法2.0相关参数接口
 =========================================================*/
func BatchRewriteQueryParaSettingCmd(ctx *gin.Context) {
	//var sev [12]string
	//applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
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
	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite

	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)
	ParaTaskAssignment(7, MultipleParaRewritesu, segment, ka, kb, kc, kf, T1a, T1c, T2a, T2c, Tda, Tdc, Wc, Wo)
}

/*=========================================================
 * 函数名称： RewriteParameterSettingCmd(ctx *gin.Context)
 * 功能描述: 改写参数设置（20）parameter setting接口
 =========================================================*/
func BatchRewriteParameterSettingCmd(ctx *gin.Context) {
	//applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
	Script := ctx.DefaultQuery("script", "10")
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
	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite

	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)
	//fmt.Println("输出内容", applianceid, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn)
	TaskAssignmentParameter(7, MultipleParaRewritesu, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn, Script)
}

/*=========================================================
 * 函数名称： RewriteNoDebugFirstCmd(ctx *gin.Context)
 * 功能描述: 改写非调试模式参数第一组接口
 =========================================================*/
func BatchRewriteNoDebugFirstCmd(ctx *gin.Context) {

	ReWaterFlow := ctx.DefaultQuery("rewaterflow", "00")               //回水水流值
	WindPressureSensor := ctx.DefaultQuery("windpressuresensor", "00") //风压传感器报警点补偿值

	/************如果本次改写的值与上次的参数值相同，则不发送改写指令**************/
	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite

	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)

	TaskAssignment(7, MultipleParaRewritesu, ReWaterFlow, WindPressureSensor)
}

/*=========================================================
 * 函数名称： BatchRewriteNoDebugSecondCmd
 * 功能描述: 改写非调试模式参数第二组接口
 =========================================================*/
func BatchRewriteNoDebugSecondCmd(ctx *gin.Context) {
	MaxCurrCoeff := ctx.DefaultQuery("maxcurrcoeff", "34")

	MinCurrCoeff := ctx.DefaultQuery("mincurrcoeff", "34")

	MaxDutyCycCoeff := ctx.DefaultQuery("maxdutycyccoeff", "34")

	MinDutyCycCoeff := ctx.DefaultQuery("mindutycyccoeff", "34")

	db := common.GetDB()
	var MultipleParaRewritesu []MultipleParaRewrite
	db.Table("multiple_para_rewrite").Where("rewriteflag = ?", "1").Find(&MultipleParaRewritesu)
	TaskAssignmentNoDebugSecond(7, MultipleParaRewritesu, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff)
}

/*=========================================================
 * 函数名称：QueryParaSetting
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/
func QueryParaSettingBa(applianceId string) (outbyte []string) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[9] + ","
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
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容")
	fmt.Println(*str)

	reply := Jsonyinhao(respBytes)
	fmt.Println("输出", reply["reply"])

	tr := strings.Split(reply["reply"], ",") //转为不带逗号的切片
	fmt.Println("tr", tr)
	return tr
}

/*=========================================================
 * 函数名称：sendCommand
 * 功能描述: 发送指令 wifi->mcu
 * 参    数：applianceId：设备ID号，index：指令在MechineCode中的索引,parameter: 人为指定的参数
 * 返回参数：机器码
 =========================================================*/
func BasendCommand(applianceId string, index string, value string) {
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
	fmt.Println("id:", applianceId)
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
	url := "http://47.111.4.75:13148/kh/mcloud/ctrl/v1"
	request, err := http.NewRequest("POST", url, reader) //创建url
	//fmt.Println("request:", request)
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
	reply := Jsonyinhao(respBytes)

	result := reply["reply"]

	row := strings.Split(result, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		//response.Success(ctx, gin.H{"errflag": "3"}, "改写失败，格式错误")
		//return
	}
	kay(index, applianceId, row, value) //验证

}

/*参数变量，校验函数*/
func kay(index string, applianceId string, row []string, value string) {
	db := common.GetDB()
	var lastparm string
	var code model.ParameterCodes
	var vaul model.ParameterFinalSettings
	// 获取上一次的参数值
	CodeName := model.SingleParaIndex[index]
	db.Where("parameter = ?", CodeName).First(&code)
	db.Table("parameter_final_settings").
		Where("appliance_id = ? and code = ?", applianceId, code.Code).First(&vaul)
	tem, _ := strconv.ParseInt(vaul.CurrentValue, 10, 64)
	st := strconv.FormatInt(tem, 16)
	fmt.Println(st)

	lastparm = fmt.Sprintf("%02s", st)
	if len(row) != 48 {
		fmt.Println("无效数据")
		//二次查询
		out := QueryParaSettingBa(applianceId)
		if len(out) != 48 {
			NoNetwork(applianceId)
			return
		} else if lastparm == out[13] {
			Network(applianceId)
			return
		}
	} else if lastparm != value {
		Network(applianceId)
		return
	}
}

/*=========================================================
 * 函数名称：QueryParaSetting
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/
func QueryParaSettingsum(applianceId string, index int, value, indexcanshu string, swit bool, Backups string) (outbyte []string) {
	db := common.GetDB()
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[9] + "," + "00"
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
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容查询查询")
	fmt.Println(*str)

	reply := Jsonyinhao(respBytes)
	tr := strings.Split(reply["reply"], ",") //转为不带逗号的切片
	fmt.Println("tr", tr, len(tr))
	if len(tr) != 1 {
		fmt.Println("进去creattwenty:", len(tr))
		creattwenty(db, applianceId, tr, indexcanshu, "0")

		if Backups == "1" {
			creattwenty(db, applianceId, tr, indexcanshu, "2")
		}

	}

	if swit == true {
		Copysame(tr, index, value, applianceId, indexcanshu, Backups)
	}
	return tr
}

func Copysame(source []string, index int, value string, applianceId string, indexcanshu, Backups string) {

	if len(source) == 1 {
		kaysum(indexcanshu, applianceId, source, value, index, Backups)

		return
	} else {
		Script := "0a"
		source[11+index] = value
		FA := source[12]
		FF := source[13]
		PH := source[14]
		FH := source[15]
		PL := source[16]
		FL := source[17]
		dH := source[18]
		Fd := source[19]
		CH := source[20]
		FC := source[21]
		CA := source[23]
		nE := source[22]
		FP := source[24]
		HS := source[26]
		Hb := source[27]
		HE := source[28]
		HL := source[29]
		HU := source[30]
		Fn := source[33]

		outcmd := RewriteParSettingsumsingle(applianceId, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn, Script)
		kaysum(indexcanshu, applianceId, outcmd, value, index, Backups)
	}

}

func kaysum(index string, applianceId string, row []string, value string, indexcanshu int, Backups string) {
	db := common.GetDB()
	if len(row) != 48 {
		fmt.Println("无效数据")
		//二次查询
		swit := false
		out := QueryParaSettingsum(applianceId, indexcanshu, value, index, swit, Backups)
		if len(out) != 48 {
			fmt.Println("设备无联网")
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
				Updates(map[string]interface{}{
					"succeedflag": "0",
					"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
					"handleflag":  "1",
				})
			return
		} else if value == out[11+indexcanshu] {
			fmt.Println("成功")
			UpdateSigleParameter(applianceId, model.SingleParaIndex[index], value)
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
				Updates(map[string]interface{}{
					"succeedflag": "1",
					"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
					"handleflag":  "1",
				})
			creattwenty(db, applianceId, row, index, "1")
			return
		}
	} else if row[11+indexcanshu] != value {
		fmt.Println("改写失败")
		fmt.Println(row[11+indexcanshu])
		fmt.Println(row[11+indexcanshu])
		db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
			Updates(map[string]interface{}{
				"succeedflag": "0",
				"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
				"handleflag":  "1",
			})
		return
	} else {
		fmt.Println("改写成功")
		db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
			Updates(map[string]interface{}{
				"succeedflag": "1",
				"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
				"handleflag":  "1",
			})
		creattwenty(db, applianceId, row, index, "1")
		return
	}
}

func creattwenty(db *gorm.DB, applianceId string, outcmd []string, indexcanshu string, numbit string) {
	var pare model.ParameterSerials
	db.Table("parameter_serials").Where("serial_number = ?", indexcanshu).Find(&pare)
	Result := db.Table("parameter_batch_sets").Where("appliance_id = ? and check_alter = ?", applianceId, numbit).First(&model.ParamenBatch{}).Error
	if len(outcmd)==48{
	if errors.Is(Result, gorm.ErrRecordNotFound) {
		var setaltle model.ParamenBatch
		setaltle.Appliance_id = applianceId
		setaltle.FA = outcmd[12]
		setaltle.FF = outcmd[13]
		setaltle.PH = outcmd[14]
		setaltle.FH = outcmd[15]
		setaltle.PL = outcmd[16]
		setaltle.FL = outcmd[17]
		setaltle.DH = outcmd[18]
		setaltle.Fd = outcmd[19]
		setaltle.CH = outcmd[20]
		setaltle.FC = outcmd[21]
		setaltle.CA = outcmd[23]
	//fmt.Println("zhelishi",outcmd[22])
	if outcmd[22]=="02"{
		outcmd[22]="00"
	}else if outcmd[22]=="03"{
		outcmd[22]="01"
	}
		setaltle.NE = outcmd[22]
		setaltle.FP = outcmd[24]
		setaltle.LF = outcmd[25]
		setaltle.HS = outcmd[26]
		setaltle.Hb = outcmd[27]
		setaltle.HE = outcmd[28]
		setaltle.HL = outcmd[29]
		setaltle.HU = outcmd[30]
		setaltle.UA = outcmd[31]
		setaltle.Ub = outcmd[32]
		setaltle.Fn = outcmd[33]
		setaltle.CheckAlter = numbit
		setaltle.AlterCode = pare.Parameter
		db.Table("parameter_batch_sets").Create(&setaltle)
		fmt.Println("创建成功")
	} else {
		fmt.Println("有没有值")
		if outcmd[22]=="02"{
			outcmd[22]="00"
		}else if outcmd[22]=="03"{
			outcmd[22]="01"
		}
		db.Table("parameter_batch_sets").
			Where("appliance_id = ? and check_alter = ?", applianceId, numbit).
			Updates(map[string]interface{}{
				"030012":     outcmd[12],
				"030013":     outcmd[13],
				"030014":     outcmd[14],
				"030015":     outcmd[15],
				"030016":     outcmd[16],
				"030017":     outcmd[17],
				"030018":     outcmd[18],
				"030019":     outcmd[19],
				"030020":     outcmd[20],
				"030021":     outcmd[21],
				"030022":     outcmd[22],
				"030023":     outcmd[23],
				"030024":     outcmd[24],
				"030025":     outcmd[25],
				"030026":     outcmd[26],
				"030027":     outcmd[27],
				"030028":     outcmd[28],
				"030029":     outcmd[29],
				"030030":     outcmd[30],
				"030031":     outcmd[31],
				"030032":     outcmd[32],
				"030033":     outcmd[33],
				"alter_code": pare.Parameter,
			})

		fmt.Println("更新成功")
	}
}
}
/*=========================================================
 * 函数名称：FastartingdownDatathread
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/
func FastartingdownDatathread(applianceId string, order int) (outcmd []string) {
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

	reply := Jsonyinhao(respBytes)
	strdatas := strings.Split(reply["reply"], ",") //返回字符串类型切片

	return strdatas
}

/*=========================================================
 * 函数名称： RewriteParSetting
 * 功能描述: 发送指令，改写参数设置（20）parameter setting
 =========================================================*/
func RewriteParSettingsumsingle(applianceId string, fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string, scrit string) (ord []string) {
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

	reply := Jsonyinhao(respBytes)
	ord = strings.Split(reply["reply"], ",")
	//存储到数据库
	//res = ParaSettingSave(applianceId, reply["reply"])
	return ord
}

/*=========================================================
 * 函数名称：BasendCommandPara
 * 功能描述: 验证恒温改写的设备是否联网
 =========================================================*/
func BasendCommandPara(applianceId string, groupNum string, ka string, kb string, kc string, kf string, T1a string, T1c string, T2a string, T2c string, Tda string, Tdc string, Wc string, Wo string) {
	var sev [12]string
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
	strdatas := RewriteQueryParaSetting(applianceId, groupNum, ka, kb, kc, kf, T1a, T1c, T2a, T2c, Tda, Tdc, Wc, Wo)
	if len(strdatas) != 48 {
		fmt.Println("设备返回无效数据")
		outcmd1 := QuerySuanFaData(applianceId, groupNum) //二次查询
		fmt.Println(outcmd1)
		ReDataBase(outcmd1, applianceId, groupNum, sev)
		if len(outcmd1) != 12 {
			NoNetwork(applianceId)
			return
		} else if len(outcmd1) == 12 {
			Network(applianceId)
			return
		}
	} else if len(strdatas) == 48 {
		Network(applianceId)
		return
	}
}

/*=========================================================
 * 函数名称：BasendCommandParameter
 * 功能描述: 验证多个参数改写的设备是否联网
 =========================================================*/
func BasendCommandParameter(applianceId string, fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string, scrit string) {
	outcmd, gou := RewriteParSetting(applianceId, fa, ff, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn, scrit)
	fmt.Println(outcmd)
	if len(gou) != 48 {
		fmt.Println("无效数据")
		outcmd1 := sendCommand(applianceId, 9)
		row := strings.Split(outcmd1, ",")
		if len(row) != 48 {
			//fmt.Println("无效数据")
			NoNetwork(applianceId)
			return
		} else {
			Network(applianceId)
			return
		}
	} else {
		Network(applianceId)
		return
	}
}

/*=========================================================
 * 函数名称： BasendNoDebugFirst
 * 功能描述: 改写非调试模式参数第一组
 =========================================================*/
func BasendNoDebugFirst(applianceId string, reWaterFlow string, windPressureSensor string) {
	source := RewriteNoDebugFirst(applianceId, reWaterFlow, windPressureSensor)
	row := strings.Split(source, ",")
	if len(row) != 48 {
		outcmd := sendCommand(applianceId, 2, "00")
		if len(outcmd) != 48 {
			NoNetwork(applianceId)
			return
		}
	} else {
		Network(applianceId)
		return
	}
}

/*=========================================================
 * 函数名称： BasendNoDebugSecond
 * 功能描述: 改写非调试模式参数第二组
 =========================================================*/
func BasendNoDebugSecond(applianceId string, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff string) {
	//向设备发送命令，改写参数，并将结果存入参数最终表
	source := sendCommand(applianceId, 4, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff)
	row := strings.Split(source, ",")
	if len(row) != 48 {
		outcmd := sendCommand(applianceId, 2, "00")
		if len(outcmd) != 48 {
			NoNetwork(applianceId)
			return
		}
	} else {
		Network(applianceId)
		return
	}
}

/*=========================================================
 * 函数名称：QueryParaSetting
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/
func QueryParaSettingtwenty(applianceId string, db *gorm.DB) {
	CodeLen := 31 //机器码总位数
	cmdCodeOut := MechineCode[9] + "," + "00"
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
	//byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("输出内容查询查询")
	fmt.Println(*str)

	reply := Jsonyinhao(respBytes)
	fmt.Println("输出", reply)
	fmt.Println("输出", reply["reply"])

	tr := strings.Split(reply["reply"], ",") //转为不带逗号的切片
	var set model.ParamenBatch
	set.Appliance_id = applianceId
	set.FA = tr[12]
	set.FF = tr[13]
	set.PH = tr[14]
	set.FH = tr[15]
	set.PL = tr[16]
	set.FL = tr[17]
	set.DH = tr[18]
	set.Fd = tr[19]
	set.CH = tr[20]
	set.FC = tr[21]
	set.CA = tr[23]
	set.NE = tr[22]
	set.FP = tr[24]
	set.HS = tr[26]
	set.Hb = tr[27]
	set.HE = tr[28]
	set.HL = tr[29]
	set.HU = tr[30]
	set.Fn = tr[33]
	set.CheckAlter = "0"
	db.Table("parameter_batch_sets").Create(&set)
	fmt.Println("读取参数20值", set)

}

/*=========================================================
 * 函数名称： SumNoparasettingfirstCheck
 * 功能描述: 检查参数是否正确
 =========================================================*/
func SumNoparasettingfirstCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//回水水流值
		applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
		Pattern := ctx.Query("pattern")
		fmt.Println(Pattern)
		if len(applianceid) != 15 {
			response.Success(ctx, gin.H{"errflag": "1"}, "设备号错误")
			ctx.Abort()
		} else {
			ctx.Next()
		}
		var Format []string //用来保存参数格式的错误（输入的参数无效）
		var Long []string   //用来保存长度出错
		var top []string    //用来保存超出上线
		var Down []string   //用来保存超出下线

		PatternInt, err := strconv.ParseInt(Pattern, 16, 64)
		switch PatternInt {
		/*=========================================================
		* 功能描述:  验证 非参数设置可调参数第一组 的逻辑
		* 模式选择：1
		* 创建日期： 2021.12.25
		=========================================================*/
		case 1:
			{
				ReWaterFlow := ctx.Query("rewaterflow")
				//风压传感器报警点补偿值
				WindPressureSensor := ctx.Query("windpressuresensor")
				//_, err := strconv.ParseInt(ReWaterFlow, 16, 64)
				if err != nil {
					Format = append(Format, "ReWaterFlow")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(ReWaterFlow) > 2 {
					Long = append(Long, "ReWaterFlow")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				_, err = strconv.ParseInt(WindPressureSensor, 16, 64)
				if err != nil {
					Format = append(Format, "WindPressureSensor")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(WindPressureSensor) > 2 {
					Long = append(Long, "WindPressureSensor")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//	return
				}
				if len(Format) == 0 && len(Long) == 0 {
					response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
				} else if len(Format) != 0 && len(Long) == 0 {
					response.Success(ctx, gin.H{"errflag": "2", "无效参数": Format}, "参数格式错误")
				} else if len(Format) == 0 && len(Long) != 0 {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long}, "输入参数长度超过范围")
				} else {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long, "无效参数": Format}, "输入参数超过范围")
				}
				//response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
				//ctx.Next()
			}
			/*=========================================================
			* 功能描述:  验证 非参数设置可调参数第二组 的逻辑
			* 模式选择：2
			* 创建日期： 2021.12.25
			=========================================================*/
		case 2:
			{
				MaxCurrCoeff := ctx.Query("maxcurrcoeff")
				//最小负荷风机电流偏差系数MinCurrCoeff
				MinCurrCoeff := ctx.Query("mincurrcoeff")
				//最大负荷风机占空比偏差系数MaxDutyCycCoeff
				MaxDutyCycCoeff := ctx.Query("maxdutycyccoeff")
				//最小负荷风机占空比偏差系数MinDutyCycCoeff
				MinDutyCycCoeff := ctx.Query("mindutycyccoeff")

				MaxCurrCoeffInt, err := strconv.ParseInt(MaxCurrCoeff, 16, 64)
				if err != nil {
					Format = append(Format, "MaxCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(MaxCurrCoeff) > 2 {
					Long = append(Long, "MaxCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if MaxCurrCoeffInt < 80 {
					Down = append(Down, "MaxCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if MaxCurrCoeffInt > 120 {
					top = append(top, "MaxCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}

				MinCurrCoeffInt, err := strconv.ParseInt(MinCurrCoeff, 16, 64)
				if err != nil {
					Format = append(Format, "MinCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(MinCurrCoeff) > 2 {
					Long = append(Long, "MinCurrCoeff")
					///response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if MinCurrCoeffInt < 80 {
					Down = append(Down, "MinCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if MinCurrCoeffInt > 120 {
					top = append(top, "MinCurrCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}

				MaxDutyCycCoeffInt, err := strconv.ParseInt(MaxDutyCycCoeff, 16, 64)
				if err != nil {
					Format = append(Format, "MaxDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(MaxDutyCycCoeff) > 2 {
					Long = append(Long, "MaxDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if MaxDutyCycCoeffInt < 80 {
					Down = append(Down, "MaxDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if MaxDutyCycCoeffInt > 120 {
					top = append(top, "MaxDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				MinDutyCycCoeffInt, err := strconv.ParseInt(MinDutyCycCoeff, 16, 64)
				if err != nil {
					Format = append(Format, "MinDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(MinDutyCycCoeff) > 2 {
					Long = append(Long, "MinDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if MinDutyCycCoeffInt < 80 {
					Down = append(Down, "MinDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if MinDutyCycCoeffInt > 120 {
					top = append(top, "MinDutyCycCoeff")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//		return
				}
				//ctx.Next()
				if len(Format) == 0 && len(Long) == 0 {
					if len(Down) != 0 && len(top) == 0 {
						response.Success(ctx, gin.H{"errflag": "4", "超出下线参数": Down}, "输入参数超出下限")
					} else if len(Down) == 0 && len(top) != 0 {
						response.Success(ctx, gin.H{"errflag": "4", "超出上线参数": top}, "输入参数超出上限")
					} else if len(Down) != 0 && len(top) != 0 {
						response.Success(ctx, gin.H{"errflag": "2", "超出下线参数": Down, "超出上线参数": top}, "输入参数值超出范围")
					} else {
						response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
					}
				} else if len(Format) != 0 && len(Long) == 0 {
					response.Success(ctx, gin.H{"errflag": "2", "无效参数": Format}, "参数格式错误")
				} else if len(Format) == 0 && len(Long) != 0 {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long}, "输入参数长度超过范围")
				} else {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long, "无效参数": Format}, "输入参数超过范围")
				}
				//r
				//response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
			}
			/*=========================================================
			* 模式选择：3
			* 功能描述:  验证 单个改写参数设置参数值 的逻辑
			* 创建日期： 2021.12.25
			=========================================================*/
		case 3:
			{
				index := ctx.Query("index") //改写序号
				value := ctx.Query("value") //改写值
				SingleMin := ctx.Query("singlemin")
				SingleMax := ctx.Query("singlemax")
				SingleMinInt, _ := strconv.ParseInt(SingleMin, 16, 64)
				SingleMaxInt, _ := strconv.ParseInt(SingleMax, 16, 64)

				_, err := strconv.ParseInt(index, 16, 64)
				if err != nil {
					response.Success(ctx, gin.H{"errflag": "2", "参数名称": index}, "参数格式错误")
					ctx.Abort()
					return
				}
				if len(index) > 2 {
					response.Success(ctx, gin.H{"errflag": "2", "参数名称": index}, "输入参数超过范围")
					ctx.Abort()
					return
				}

				ValueInt, _ := strconv.ParseInt(value, 16, 64)
				if err != nil {
					response.Success(ctx, gin.H{"errflag": "2", "参数名称": index}, "参数格式错误")
					ctx.Abort()
					return
				}
				if len(value) > 2 {
					response.Success(ctx, gin.H{"errflag": "2", "参数名称": index}, "输入参数超过范围")
					ctx.Abort()
					return
				}
				if ValueInt < SingleMinInt {
					response.Success(ctx, gin.H{"errflag": "4", "参数名称": index}, "输入参数超出下限")
					ctx.Abort()
					return
				} else if ValueInt > SingleMaxInt {
					response.Success(ctx, gin.H{"errflag": "4", "参数名称": index}, "输入参数超出上限")
					ctx.Abort()
					return
				}
				//ctx.Next()
				response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
			}
		/*=========================================================
		* 模式选择： 4
		* 功能描述:  验证 恒温算法2.0相关数据设置默认参数信息 的逻辑
		* 创建日期： 2021.12.25
		=========================================================*/

		case 4:
			{
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

				segmentInt, err := strconv.ParseInt(segment, 16, 64) //段序号
				if err != nil {
					response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					ctx.Abort()
					return
				}
				if segmentInt != 1 && segmentInt != 2 && segmentInt != 3 && segmentInt != 0 {
					response.Success(ctx, gin.H{"errflag": "2"}, "输入不存在的段序号")
					ctx.Abort()
					return
				}
				kaInt, err := strconv.ParseInt(ka, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "ka")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(ka) > 2 {
					Long = append(Long, "ka")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if kaInt < 0 {
					Down = append(Down, "ka")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				kbInt, err := strconv.ParseInt(kb, 16, 64)
				if err != nil {
					Format = append(Format, "kb")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(kb) > 2 {
					Long = append(Long, "kb")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if kbInt < 0 {
					Down = append(Down, "kb")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				kcInt, err := strconv.ParseInt(kc, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "kc")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(kc) > 2 {
					Long = append(Long, "kc")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if kcInt < 0 {
					Down = append(Down, "kc")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				kfInt, err := strconv.ParseInt(kf, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "kf")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(kf) > 2 {
					Long = append(Long, "kf")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if kfInt < 0 {
					Down = append(Down, "kf")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				T1aInt, err := strconv.ParseInt(T1a, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "T1a")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(T1a) > 2 {
					Long = append(Long, "T1a")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if T1aInt < 0 {
					Down = append(Down, "T1a")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				T1cInt, err := strconv.ParseInt(T1c, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "T1c")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(T1c) > 2 {
					Long = append(Long, "T1c")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if T1cInt < 0 {
					Down = append(Down, "T1c")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				T2aInt, err := strconv.ParseInt(T2a, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "T2a")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(T2a) > 2 {
					Long = append(Long, "T2a")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if T2aInt < 0 {
					Down = append(Down, "T2a")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				T2cInt, err := strconv.ParseInt(T2c, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "T2c")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(T2c) > 2 {
					Long = append(Long, "T2c")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if T2cInt < 0 {
					Down = append(Down, "T2c")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				TdaInt, err := strconv.ParseInt(Tda, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "Tda")
					//esponse.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//tx.Abort()
					//eturn
				}
				if len(Tda) > 2 {
					Long = append(Long, "Tda")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if TdaInt < 0 {
					Down = append(Down, "Tda")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				TdcInt, err := strconv.ParseInt(Tdc, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "Tdc")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(Tdc) > 2 {
					Long = append(Long, "Tdc")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if TdcInt < 0 {
					Down = append(Down, "Tdc")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				WcInt, err := strconv.ParseInt(Wc, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "Wc")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(Wc) > 2 {
					Long = append(Long, "Wc")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if WcInt < 0 {
					Down = append(Down, "Wc")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				WoInt, err := strconv.ParseInt(Wo, 16, 64) //段序号
				if err != nil {
					Format = append(Format, "Wo")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(Wo) > 2 {
					Long = append(Long, "Wo")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if WoInt < 0 {
					Down = append(Down, "Wo")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}

				if len(Format) == 0 && len(Long) == 0 {
					if len(Down) != 0 {
						response.Success(ctx, gin.H{"errflag": "4", "超出下线参数": Down}, "输入参数超出下限")
					} else {
						response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
					}
				} else if len(Format) != 0 && len(Long) == 0 {
					response.Success(ctx, gin.H{"errflag": "2", "无效参数": Format}, "参数格式错误")
				} else if len(Format) == 0 && len(Long) != 0 {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long}, "输入参数长度超过范围")
				} else {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long, "无效参数": Format}, "输入参数超过范围")
				}
				//ctx.Next()
				//response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
			}
		/*=========================================================
		* 模式选择： 5
		* 功能描述:  验证 多个改写参数设置的参数值 的逻辑
		* 创建日期： 2022.01.06
		=========================================================*/
		case 5:
			{
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
				_, err := strconv.ParseInt(FA, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "FA")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(FA) > 2 {
					Long = append(Long, "FA")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				_, err = strconv.ParseInt(FF, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "FF")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(FF) > 2 {
					Long = append(Long, "FF")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				PHInt, err := strconv.ParseInt(PH, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "PH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(PH) > 2 {
					Long = append(Long, "PH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if PHInt < 10 {
					Down = append(Down, "PH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if PHInt > 200 {
					top = append(top, "PH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				FHInt, err := strconv.ParseInt(FH, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "FH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(FH) > 2 {
					Long = append(Long, "FH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if FHInt < 10 {
					Down = append(Down, "FH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if FHInt > 160 {
					top = append(top, "FH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				PLInt, err := strconv.ParseInt(PL, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "PL")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(PL) > 2 {
					Long = append(Long, "PL")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if PLInt < 10 {
					Down = append(Down, "PL")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if PHInt > 200 {
					top = append(top, "PL")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				FLInt, err := strconv.ParseInt(FL, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "FL")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(FL) > 2 {
					Long = append(Long, "FL")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if FLInt < 10 {
					Down = append(Down, "FL")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if FLInt > 160 {
					top = append(top, "FL")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				dHInt, err := strconv.ParseInt(dH, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "dH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(dH) > 2 {
					Long = append(Long, "dH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if dHInt < 10 {
					Down = append(Down, "dH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if dHInt > 200 {
					top = append(top, "dH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				FdInt, err := strconv.ParseInt(Fd, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "Fd")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(Fd) > 2 {
					Long = append(Long, "Fd")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if FdInt < 10 {
					Down = append(Down, "Fd")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if FdInt > 160 {
					top = append(top, "Fd")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				CHInt, err := strconv.ParseInt(CH, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "CH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(CH) > 2 {
					Long = append(Long, "CH")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if CHInt < 10 {
					Down = append(Down, "CH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if CHInt > 200 {
					top = append(top, "CH")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				FCInt, err := strconv.ParseInt(FC, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "FC")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(FC) > 2 {
					Long = append(Long, "FC")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if FCInt < 10 {
					Down = append(Down, "FC")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if FCInt > 160 {
					top = append(top, "FC")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				CAInt, err := strconv.ParseInt(CA, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "CA")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(CA) > 2 {
					Long = append(Long, "CA")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if CAInt < 0 {
					Down = append(Down, "CA")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if CAInt > 5 {
					top = append(top, "CA")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				nEInt, err := strconv.ParseInt(nE, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "nE")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(nE) > 2 {
					Long = append(Long, "nE")
					//response.Success(ctx, gin.H{"errflag": "2","参数名称": "nE" }, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if nEInt < 0 {
					Down = append(Down, "nE")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if nEInt > 1 {
					top = append(top, "nE")
					//response.Success(ctx, gin.H{"errflag": "4","参数名称": "nE" }, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				FPInt, err := strconv.ParseInt(FP, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "FP")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(FP) > 2 {
					Long = append(Long, "FP")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if FPInt < 4 {
					Down = append(Down, "FP")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				}
				HSInt, err := strconv.ParseInt(HS, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "HS")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(HS) > 2 {
					Long = append(Long, "HS")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if HSInt < 50 {
					Down = append(Down, "HS")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if HSInt > 100 {
					top = append(top, "HS")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				HbInt, err := strconv.ParseInt(Hb, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "Hb")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(Hb) > 2 {
					Long = append(Long, "Hb")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if HbInt < 0 {
					Down = append(Down, "Hb")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if HbInt > 5 {
					top = append(top, "Hb")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				HEInt, err := strconv.ParseInt(HE, 16, 64) //参数值
				if err != nil {
					//Format = append(Format, "HE")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					return
				}
				if len(HE) > 2 {
					Long = append(Long, "HE")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if HEInt < 20 {
					Down = append(Down, "HE")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if HEInt > 50 {
					top = append(top, "HE")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				HLInt, err := strconv.ParseInt(HL, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "HL")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(HL) > 2 {
					Long = append(Long, "HL")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if HLInt < 40 {
					Down = append(Down, "HL")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if HLInt > 60 {
					top = append(top, "HL")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				HUInt, err := strconv.ParseInt(HU, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "HU")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(HU) > 2 {
					Long = append(Long, "HU")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if HUInt < 0 {
					Down = append(Down, "HU")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if HUInt > 99 {
					top = append(top, "HU")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				FnInt, err := strconv.ParseInt(Fn, 16, 64) //参数值
				if err != nil {
					Format = append(Format, "Fn")
					//response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
					//ctx.Abort()
					//return
				}
				if len(Fn) > 2 {
					Long = append(Long, "Fn")
					//response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
					//ctx.Abort()
					//return
				}
				if FnInt < 0 {
					Down = append(Down, "Fn")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
					//ctx.Abort()
					//return
				} else if FnInt > 2 {
					top = append(top, "Fn")
					//response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
					//ctx.Abort()
					//return
				}
				if len(Format) == 0 && len(Long) == 0 {
					if len(Down) != 0 && len(top) == 0 {
						response.Success(ctx, gin.H{"errflag": "4", "超出下线参数": Down}, "输入参数超出下限")
					} else if len(Down) == 0 && len(top) != 0 {
						response.Success(ctx, gin.H{"errflag": "4", "超出上线参数": top}, "输入参数超出上限")
					} else if len(Down) != 0 && len(top) != 0 {
						response.Success(ctx, gin.H{"errflag": "2", "超出下线参数": Down, "超出上线参数": top}, "输入参数值超出范围")
					} else {
						response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
					}

				} else if len(Format) != 0 && len(Long) == 0 {
					response.Success(ctx, gin.H{"errflag": "2", "无效参数": Format}, "参数格式错误")
				} else if len(Format) == 0 && len(Long) != 0 {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long}, "输入参数长度超过范围")
				} else {
					response.Success(ctx, gin.H{"errflag": "2", "长度超出参数": Long, "无效参数": Format}, "输入参数超过范围")
				}
				//ctx.Next()
				//response.Success(ctx, gin.H{"errflag": "1"}, "校验通过")
			}
		}

	}
}
