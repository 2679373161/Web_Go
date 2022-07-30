package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"

	//"ginEssential/response"
	"io/ioutil"
	"net/http"
	"runtime"
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
}

type MultipleParaRewrite struct {
	ApplianceId string `json:"appliance_id" gorm:"column:appliance_id"`
	HandleFlag  string ` json:"handleflag" gorm:"column:handleflag"`
	SucceedFlag string ` json:"succeedflag" gorm:"column:succeedflag"`
	RewriteFlag string ` json:"rewriteflag" gorm:"column:rewriteflag`
	Index       string ` json:"index" gorm:"column:index`
	Value       string ` json:"value" gorm:"column:value`
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

/*=========================================================
 * 函数名称：QueryParaSetting
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/
func QueryParaSetting(applianceId string, index string, value string, swit bool) (outbyte []string) {
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

	// var m = make(map[string]interface{})

	// _ = json.Unmarshal(respBytes, &m) //数据处理
	// //fmt.Println(m["retcode"])
	// fmt.Println(m["retCode"])
	// fmt.Println(m["desc"])
	// fmt.Println(m["result"])

	// temp, _ := json.Marshal(m["result"])
	// var reply = make(map[string]string)
	// _ = json.Unmarshal(temp, &reply)
	reply := Jsonyinhao(respBytes)
	fmt.Println("输出", reply)
	fmt.Println("输出", reply["reply"])

	tr := strings.Split(reply["reply"], ",") //转为不带逗号的切片
	fmt.Println("tr", tr)
	if swit == true {
		sendCommand(applianceId, index, value)

	}
	return tr
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
	return string(bank)
}

/*=========================================================
 * 函数名称： BatchRewriteSingleParaCmd
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
	/************如果本次改写的值与上次的参数值相同，则不发送改写指令**************/
	TaskAssignment(7, MultipleParaRewritesu, index, value)

	//sendCommand(te)
}

func TaskAssignment(thread int, equ []MultipleParaRewrite, index, value string) {

	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)

	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)

	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go worker(jobs, results, equ, index, value)
	}

	// 每台设备一个任务
	for j := 0; j < len(equ); j++ {
		jobs <- j
		work <- equ[j].ApplianceId
	}
	close(jobs)

	// 输出结果
	for a := 0; a < len(equ); a++ {
		<-results
	}
}
func worker(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite, index, value string) {
	swi := true
	for i := range jobs {

		QueryParaSetting(equ[i].ApplianceId, index, value, swi)

		results <- i
		fmt.Println("线程i:", i)
	}
}

/*=========================================================
 * 函数名称：sendCommand
 * 功能描述: 发送指令 wifi->mcu
 * 参    数：applianceId：设备ID号，index：指令在MechineCode中的索引,parameter: 人为指定的参数
 * 返回参数：机器码
 =========================================================*/
func sendCommand(applianceId string, index string, value string) {
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
	// //byte数组直接转成string,优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("返回内容：")
	fmt.Println(*str)
	// /*******获取数据********/
	// // 反序列化
	// var m map[string]interface{}
	// _ = json.Unmarshal(respBytes, &m)
	// // 再次反序列化
	// temp, _ := json.Marshal(m["result"])
	// var reply map[string]string
	// _ = json.Unmarshal(temp, &reply)
	// fmt.Println("机器码：")
	// fmt.Println(reply["reply"])
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
		swit := false
		out := QueryParaSetting(applianceId, index, value, swit)
		if len(out) != 48 {
			fmt.Println("设备无联网")
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
				Updates(map[string]interface{}{
					"succeedflag": "0",
					"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
					"handleflag":  "1",
				})
			return
		} else if lastparm == out[13] {
			fmt.Println("成功")
			UpdateSigleParameter(applianceId, model.SingleParaIndex[index], value)
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
				Updates(map[string]interface{}{
					"succeedflag": "1",
					"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
					"handleflag":  "1",
				})
			return
		}
	} else if row[12] != index || row[13] != value {
		fmt.Println("改写失败")
		fmt.Println(row[12])
		fmt.Println(row[13])
		db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
			Updates(map[string]interface{}{
				"succeedflag": "0",
				"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
				"handleflag":  "1",
			})
		return
	} else if lastparm != value {
		fmt.Println("改写成功")
		db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
			Updates(map[string]interface{}{
				"succeedflag": "1",
				"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
				"handleflag":  "1",
			})
		// 1.将改写的数据存入 参数变化记录表 中
		// UpdateSigleParameter(applianceId, model.SingleParaIndex[index], value)
		return
	}
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

/*=========================================================
 * 函数名称： Jsonyinhao
 * 功能描述: 数据解析函数
 =========================================================*/
func Jsonyinhao(respBytes []byte) map[string]string {
	for i := 0; i < len(respBytes); i++ {
		if respBytes[i] == 92 {
			respBytes = append(respBytes[:i], respBytes[i+1:]...)
		}
	}

	if respBytes[0] == 34 {
		respBytes = respBytes[1 : len(respBytes)-2]
	}

	//respBytes =respBytes[1:len(respBytes)-1]

	var m = make(map[string]string)
	json.Unmarshal(respBytes, &m)
	return m
}

/*=========================================================
 * 函数名称： BatchRewriteSumParaCmd
 * 功能描述:  批量改写总参数  参数设置参数值
 =========================================================*/
func BatchRewriteSumParaCmd(ctx *gin.Context) {

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
	/************如果本次改写的值与上次的参数值相同，则不发送改写指令**************/
	TaskAssignmentsum(7, MultipleParaRewritesu, indexsrevu, value, index)

	//sendCommand(te)
}

func TaskAssignmentsum(thread int, equ []MultipleParaRewrite, index int, value, indexcanshu string) {

	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)

	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)

	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go workersum(jobs, results, equ, index, value, indexcanshu)
	}

	// 每台设备一个任务
	for j := 0; j < len(equ); j++ {
		jobs <- j
		work <- equ[j].ApplianceId
	}
	close(jobs)

	// 输出结果
	for a := 0; a < len(equ); a++ {
		<-results
	}
}
func workersum(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite, index int, value, indexcanshu string) {
	swit := true
	for i := range jobs {

		QueryParaSettingsum(equ[i].ApplianceId, index, value, indexcanshu, swit)

		results <- i
		fmt.Println("线程i:", i)
	}
}

/*=========================================================
 * 函数名称：QueryParaSetting
 * 功能描述: 向设备发送查询指令，查询参数设置默认参数信息
 =========================================================*/
func QueryParaSettingsum(applianceId string, index int, value, indexcanshu string, swit bool) (outbyte []string) {
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
	fmt.Println("tr", tr)
	if swit == true {
		Copysame(tr, index, value, applianceId, indexcanshu)
	}
	return tr
}

func Copysame(source []string, index int, value string, applianceId string, indexcanshu string) {
	fmt.Println("进去了:", source)
	fmt.Println("进去了2:", len(source))
	if len(source) == 1 {
		kaysum(indexcanshu, applianceId, source, value, index)

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
		outcmd := RewriteParSetting(applianceId, FA, FF, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn, Script)
		kaysum(indexcanshu, applianceId, outcmd, value, index)
	}

}

/*=========================================================
 * 函数名称： RewriteParSetting
 * 功能描述: 发送指令，改写参数设置（20）parameter setting
 =========================================================*/
func RewriteParSetting(applianceId string, fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string, scrit string) (ord []string) {
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
	fmt.Println(reply["reply"])
	ord = strings.Split(reply["reply"], ",")
	//存储到数据库
	//res = ParaSettingSave(applianceId, reply["reply"])
	return ord
}

func kaysum(index string, applianceId string, row []string, value string, indexcanshu int) {
	db := common.GetDB()
	//var lastparm string
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

	//lastparm = fmt.Sprintf("%02s", st)
	if len(row) != 48 {
		fmt.Println("无效数据")
		//二次查询
		swit := false
		out := QueryParaSettingsum(applianceId, indexcanshu, value, index, swit)
		if len(out) != 48 {
			fmt.Println("设备无联网")
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
				Updates(map[string]interface{}{
					"succeedflag": "0",
					"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
					"handleflag":  "1",
				})
			return
		} else { // if lastparm == out[12+indexcanshu]
			fmt.Println("成功")
			UpdateSigleParameter(applianceId, model.SingleParaIndex[index], value)
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
				Updates(map[string]interface{}{
					"succeedflag": "1",
					"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
					"handleflag":  "1",
				})
			return
		}
		// } else if row[12] != index || row[13] != value {
		// 	fmt.Println("改写失败")
		// 	fmt.Println(row[12])
		// 	fmt.Println(row[13])
		// 	db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
		// 		Updates(map[string]interface{}{
		// 			"succeedflag": "0",
		// 			"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
		// 			"handleflag":  "1",
		// 		})
		// 	return
	} else {
		fmt.Println("改写成功")
		db.Table("multiple_para_rewrite").Where("appliance_id = ?", applianceId).
			Updates(map[string]interface{}{
				"succeedflag": "1",
				"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
				"handleflag":  "1",
			})
		// 1.将改写的数据存入 参数变化记录表 中
		// UpdateSigleParameter(applianceId, model.SingleParaIndex[index], value)
		return
	}
}
