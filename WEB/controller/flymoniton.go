package controller

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"ginEssential/common"
	"ginEssential/model"
	"strconv"
	"strings"
	"time"
)

//获取实时采集中的c4故障设备
func QueryRealC4Error(devId , startTime , endTime , devType string) {
	db := common.GetIndexDB()
	errCount := 0
	db.
		Table("e_midea_fault").
		Where("dev_id = ? and start_time >= ? and end_time <= ? and c4 != 0",devId,startTime,endTime).Count(&errCount)
		if errCount > 0 {
			alterFanParm(devId,devType)
		}
}

/*=========================================================
 * 函数名称： QueryRealC4ErrorTimer
 * 功能描述:  定时修改c4故障设备
 * 输入参数： dataTime-发生故障的时间,格式：2022-02-01
 * 创建日期： 2022.06.01
 =========================================================*/
func QueryRealC4ErrorTimer(dataTime string) {
	db := common.GetIndexDB()
	dor := common.GetDB()
	type C4ErrDev struct {
		DevId string
		DevType string
	}
	var c4ErrDevs []C4ErrDev
	db.Table("daily_monitorings").Where("time_date = ? and c4 != 0",dataTime).Scan(&c4ErrDevs)
	for _ , c4ErrDev := range c4ErrDevs {
		var bo model.Bo
		dor.Table("bo").Where("dev_id = ?",c4ErrDev.DevId).Scan(&bo)
		if bo.HandleFlag == "1" && bo.MonitoringFlag == "1" && bo.SentcmdFlag == "1"{
			alterFanParm(c4ErrDev.DevId,c4ErrDev.DevType)
		}
	}
}

/*=========================================================
 * 函数名称： alterFanParm
 * 功能描述:  修改风机自学习参数
 * 创建日期： 2022.05.10
 * 修改时间： 2022.05.25
 =========================================================*/
func alterFanParm(applianceid , devType string) {
	db := common.GetIndexDB()
	fmt.Println("-----------------C4~~风机故障参数修改进入-----------------")
	fmt.Println("修改设备：",applianceid)
	defer fmt.Println("-----------------C4~~风机故障参数修改退出-----------------")
	var readParmNum = 3 //设备参数读取次数上限
	var sendParmNum = 3 //设备参数修改次数上限
	/*1.读取设备风机自学习参数*/
	groupNum := "00"
	var row []string
	//最多读10次
	for i := 0 ; i < readParmNum ; i ++ {
		source := sendCommand(applianceid,2,groupNum)
		// 1.1 按逗号分割
		row1 := strings.Split(source, ",")
		if len(row1) != 48 {
			fmt.Println("风机参数查询数据失败，查询次数：",i+1)
			time.Sleep(2 * time.Second)//2秒后重新查
		} else {
			fmt.Println("风机参数查询成功")
			row = row1
			break
		}
	}
	if len(row) != 48 {
		fmt.Println("设备参数查询失败，风机参数修改退出~~~~")
	}
	// 1.2 将字符串转化为数字
	var num []int64
	//var str []string
	for i := 13; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	fmt.Println("设备风机参数：",num)
	/*判断设备风机类型*/
	fanType := struct {
		FanType string
	}{}
	db.Table("keyparamaters").Where("SN = ?",devType).Scan(&fanType)
	var indexMax , indexMin int
	if fanType.FanType == "正弦波风机" {
		indexMax = 2
		indexMin = 3
	} else if fanType.FanType == "方波风机" {
		indexMax = 0
		indexMin = 1
	} else {
		fmt.Println("风机类型错误，风机参数修改退出~~~~")
		return
	}

	/*2.判断设备自学习参数有无小于等于85*/
	maxDutyCycCoeff := num[indexMax]
	minDutyCycCoeff := num[indexMin]
	sendMaxFlag := 0
	sendMinFlag := 0

	if maxDutyCycCoeff > 85 {
		num[indexMax] -= 5
		sendMaxFlag = 1
	}
	if minDutyCycCoeff > 85 {
		num[indexMin] -= 5
		sendMinFlag = 1
	}
	/*3.参数下发*/
	if (sendMinFlag == 1) || (sendMaxFlag==1) {
		// 1.3 将数字变为十六进制字符串,得到参数值
		var str []string
		for _, v := range num {
			temp := strconv.FormatInt(v, 16)
			str = append(str, temp)
		}
		//向设备发送命令，改写参数
		//重复修改100次，不成功则退出
		var numDstF []int64
		for i := 0 ; i < sendParmNum ; i++ {
			sendCommand(applianceid, 4, str[0], str[1], str[2], str[3])
			//查询结果，验证有无修改成功
			sourceDst := sendCommand(applianceid,2,groupNum)
			rowDst := strings.Split(sourceDst, ",")
			if len(rowDst) != 48 {
				fmt.Println("回复无效数据,风机参数重新下发，次数：",i+1)
				time.Sleep(2 * time.Second)//2秒后重新查询
				continue
			}
			// 1.2 将字符串转化为数字
			var numDst []int64
			//var str []string
			for i := 13; i <= 18; i++ {
				temp, _ := strconv.ParseInt(rowDst[i], 16, 64)
				numDst = append(numDst, temp)
			}
			if (numDst[indexMax] == num[indexMax]) && (numDst[indexMin] == num[indexMin]) {
				numDstF = numDst
				fmt.Println("风机参数改写成功")
				break //改写成功则退出
			}
			//不成功则等2秒，再发一次
			fmt.Println("改写未成功,参数重新下发，次数：",i+1)
			time.Sleep(2 * time.Second)
		}
		if len(numDstF) == 0 {
			fmt.Println("回复参数获取失败,风机参数修改退出~~~~")
			return
		}
		//numDst = num
		/*4.结果记录*/
		var perchangeNum PerchangeNums
		changeFlag := 0
		perchangeNum.DevId = applianceid
		perchangeNum.UpdataTime = time.Now().Format("2006-01-02 15:04:05")
		if sendMaxFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMax]
			perchangeNum.UpValue = strconv.Itoa(int(maxDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDstF[indexMax]))
			perchangeNum.FaultType = "c4"
			if (numDstF[indexMax] == num[indexMax]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
				fmt.Println("改写未成功")
			}
			err := db.Create(&perchangeNum).Error
			if err != nil {
				fmt.Println("风机参数改写记录存储未成功~~~")
			} else {
				fmt.Println("风机参数改写记录存储成功~~~")
			}

		}
		if sendMinFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMin]
			perchangeNum.UpValue = strconv.Itoa(int(minDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDstF[indexMin]))
			perchangeNum.FaultType = "c4"
			if (numDstF[indexMin] == num[indexMin]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
				fmt.Println("改写未成功")
			}
			err := db.Create(&perchangeNum).Error
			if err != nil {
				fmt.Println("风机参数改写记录存储未成功~~~")
			} else {
				fmt.Println("风机参数改写记录存储成功~~~")
			}
		}
		if changeFlag == 1 {
			var numPerChange NumPerChanges
			err := db.Table("num_per_changes").
				Where("dev_id = ? AND fault_type= ?", applianceid, "C4").First(&numPerChange).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				numPerChange.DevId = applianceid
				numPerChange.FaultType = "C4"
				numPerChange.DataTime = time.Now().Format("2006-01-02 15:04:05")
				numPerChange.ChangeNum = "1"
				db.Create(&numPerChange)
			} else {
				ChangeNum , _ := strconv.Atoi(numPerChange.ChangeNum)
				ChangeNum++
				db.Table("num_per_changes").
					Where("dev_id = ? AND fault_type= ?", applianceid, "C4").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":           time.Now().Format("2006-01-02 15:04:05"),
						"change_num":          strconv.Itoa(ChangeNum),
					})
			}
		}
	} else {
		/*需要上门设备输出*/
		var failedDevice FailedDevice
		failedDevice.DevId = applianceid
		failedDevice.DevType = devType
		failedDevice.FaultType = "c4"
		failedDevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedDevice.ReviseFlag = 1
		err := db.Create(failedDevice).Error
		if err != nil {
			fmt.Println("风机故障设备存储未成功~~~")
		} else {
			fmt.Println("风机故障设备存储成功~~~")
		}
	}
}



func alterFanParm1(applianceid , devType string , limitParm int64) {
	db := common.GetIndexDB()
	/*1.读取设备风机自学习参数*/
	groupNum := "00"
	source := sendCommand(applianceid,2,groupNum)
	// 1.1 按逗号分割
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	// 1.2 将字符串转化为数字
	var num []int64
	//var str []string
	for i := 13; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	fmt.Println("设备参数：",num)
	/*判断设备风机类型*/
	fanType := struct {
		FanType string
	}{}
	db.Table("keyparamaters").Where("SN = ?",devType).Scan(&fanType)
	var indexMax , indexMin int
	if fanType.FanType == "正弦波风机" {
		indexMax = 2
		indexMin = 3
	} else if fanType.FanType == "方波风机" {
		indexMax = 0
		indexMin = 1
	} else {
		fmt.Println("风机类型错误")
		return
	}
	//最多改2次
	//var perchangeNum []PerchangeNums
	//IndexDB.Where("dev_id").Find(&perchangeNum)
	/*2.判断设备自学习参数有无小于等于85*/
	maxDutyCycCoeff := num[indexMax]
	minDutyCycCoeff := num[indexMin]
	sendMaxFlag := 0
	sendMinFlag := 0

	if maxDutyCycCoeff > limitParm {
		num[indexMax] -= 5
		sendMaxFlag = 1
	}
	if minDutyCycCoeff > limitParm {
		num[indexMin] -= 5
		sendMinFlag = 1
	}
	/*3.参数下发*/

	if (sendMinFlag == 1) || (sendMaxFlag==1) {
		// 1.3 将数字变为十六进制字符串,得到参数值
		var str []string
		for _, v := range num {
			temp := strconv.FormatInt(v, 16)
			str = append(str, temp)
		}
		//向设备发送命令，改写参数
		sendCommand(applianceid, 4, str[0], str[1], str[2], str[3])
		//查询结果，验证有无修改成功
		sourceDst := sendCommand(applianceid,2,groupNum)
		rowDst := strings.Split(sourceDst, ",")
		if len(rowDst) != 48 {
			fmt.Println("无效数据")
			return
		}
		// 1.2 将字符串转化为数字
		var numDst []int64
		//var str []string
		for i := 13; i <= 18; i++ {
			temp, _ := strconv.ParseInt(rowDst[i], 16, 64)
			numDst = append(numDst, temp)
		}
		/*4.结果记录*/
		var perchangeNum PerchangeNums
		changeFlag := 0
		perchangeNum.DevId = applianceid
		perchangeNum.UpdataTime = time.Now().Format("2006-01-02 15:04:05")
		if sendMaxFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMax]
			perchangeNum.UpValue = strconv.Itoa(int(maxDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDst[indexMax]))
			perchangeNum.FaultType = "c4"
			if (numDst[indexMax] == num[indexMax]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
			}
			db.Create(&perchangeNum)

		}
		if sendMinFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMin]
			perchangeNum.UpValue = strconv.Itoa(int(minDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDst[indexMin]))
			perchangeNum.FaultType = "c4"
			if (numDst[indexMin] == num[indexMin]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
			}
			db.Create(&perchangeNum)
		}
		if changeFlag == 1 {
			var numPerChange NumPerChanges
			err := db.Table("num_per_changes").
				Where("dev_id = ? AND fault_type= ?", applianceid, "C4").First(&numPerChange).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				numPerChange.DevId = applianceid
				numPerChange.FaultType = "C4"
				numPerChange.DataTime = time.Now().Format("2006-01-02 15:04:05")
				numPerChange.ChangeNum = "1"
				db.Create(&numPerChange)
			} else {
				ChangeNum , _ := strconv.Atoi(numPerChange.ChangeNum)
				ChangeNum++
				db.Table("num_per_changes").
					Where("dev_id = ? AND fault_type= ?", applianceid, "C4").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":           time.Now().Format("2006-01-02 15:04:05"),
						"change_num":          strconv.Itoa(ChangeNum),
					})
			}
		}
	} else {
		/*需要上门设备输出*/
		var failedDevice FailedDevice
		failedDevice.DevId = applianceid
		failedDevice.DevType = devType
		failedDevice.FaultType = "c4"
		failedDevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedDevice.ReviseFlag = 1
		db.Create(failedDevice)
	}
}

//前版程序
func alterFanParmOld(applianceid , devType string) {
	db := common.GetIndexDB()
	/*1.读取设备风机自学习参数*/
	groupNum := "00"
	source := sendCommand(applianceid,2,groupNum)
	// 1.1 按逗号分割
	row := strings.Split(source, ",")
	if len(row) != 48 {
		fmt.Println("无效数据")
		return
	}
	// 1.2 将字符串转化为数字
	var num []int64
	//var str []string
	for i := 13; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	fmt.Println("设备参数：",num)
	/*判断设备风机类型*/
	fanType := struct {
		FanType string
	}{}
	db.Table("keyparamaters").Where("SN = ?",devType).Scan(&fanType)
	var indexMax , indexMin int
	if fanType.FanType == "正弦波风机" {
		indexMax = 2
		indexMin = 3
	} else if fanType.FanType == "方波风机" {
		indexMax = 0
		indexMin = 1
	} else {
		fmt.Println("风机类型错误")
		return
	}

	/*2.判断设备自学习参数有无小于等于85*/
	maxDutyCycCoeff := num[indexMax]
	minDutyCycCoeff := num[indexMin]
	sendMaxFlag := 0
	sendMinFlag := 0
	if maxDutyCycCoeff > 85 {
		num[indexMax] -= 5
		sendMaxFlag = 1
	}
	if minDutyCycCoeff > 85 {
		num[indexMin] -= 5
		sendMinFlag = 1
	}
	/*3.参数下发*/
	if (sendMinFlag == 1) || (sendMaxFlag==1) {
		// 1.3 将数字变为十六进制字符串,得到参数值
		var str []string
		for _, v := range num {
			temp := strconv.FormatInt(v, 16)
			str = append(str, temp)
		}
		//向设备发送命令，改写参数
		sendCommand(applianceid, 4, str[0], str[1], str[2], str[3])
		//查询结果，验证有无修改成功
		sourceDst := sendCommand(applianceid,2,groupNum)
		rowDst := strings.Split(sourceDst, ",")
		if len(rowDst) != 48 {
			fmt.Println("无效数据")
			return
		}
		// 1.2 将字符串转化为数字
		var numDst []int64
		//var str []string
		for i := 13; i <= 18; i++ {
			temp, _ := strconv.ParseInt(rowDst[i], 16, 64)
			numDst = append(numDst, temp)
		}
		/*4.结果记录*/
		var perchangeNum PerchangeNums
		changeFlag := 0
		perchangeNum.DevId = applianceid
		perchangeNum.UpdataTime = time.Now().Format("2006-01-02 15:04:05")
		if sendMaxFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMax]
			perchangeNum.UpValue = strconv.Itoa(int(maxDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDst[indexMax]))
			perchangeNum.FaultType = "c4"
			if (numDst[indexMax] == num[indexMax]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
			}
			db.Create(&perchangeNum)

		}
		if sendMinFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMin]
			perchangeNum.UpValue = strconv.Itoa(int(minDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDst[indexMin]))
			perchangeNum.FaultType = "c4"
			if (numDst[indexMin] == num[indexMin]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
			}
			db.Create(&perchangeNum)
		}
		if changeFlag == 1 {
			var numPerChange NumPerChanges
			err := db.Table("num_per_changes").
				Where("dev_id = ? AND fault_type= ?", applianceid, "C4").First(&numPerChange).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				numPerChange.DevId = applianceid
				numPerChange.FaultType = "C4"
				numPerChange.DataTime = time.Now().Format("2006-01-02 15:04:05")
				numPerChange.ChangeNum = "1"
				db.Create(&numPerChange)
			} else {
				ChangeNum , _ := strconv.Atoi(numPerChange.ChangeNum)
				ChangeNum++
				db.Table("num_per_changes").
					Where("dev_id = ? AND fault_type= ?", applianceid, "C4").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":           time.Now().Format("2006-01-02 15:04:05"),
						"change_num":          strconv.Itoa(ChangeNum),
					})
			}
		}
	} else {
		/*需要上门设备输出*/
		var failedDevice FailedDevice
		failedDevice.DevId = applianceid
		failedDevice.DevType = devType
		failedDevice.FaultType = "c4"
		failedDevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedDevice.ReviseFlag = 1
		db.Create(failedDevice)
	}
}

/*=========================================================
 * 函数名称： alterFanParmInterface
 * 功能描述:  修改风机自学习参数,后端接口
 * 返回参数： 0-修改失败；1-修改成功；2-没有修改（修改条件不成立）
 * 创建日期： 2022.05.10
 * 修改时间： 2022.05.25
 =========================================================*/

func alterFanParmInterface(applianceid , devType ,rewrinteflag string) (res int ,maxDut ,minDut int64 ) {
	db := common.GetIndexDB()
	fmt.Println("-----------------C4~~风机故障参数修改进入-----------------")
	fmt.Println("修改设备：",applianceid)
	defer fmt.Println("-----------------C4~~风机故障参数修改退出-----------------")
	var readParmNum = 3 //设备参数读取次数上限
	var sendParmNum = 3 //设备参数修改次数上限
	/*1.读取设备风机自学习参数*/
	groupNum := "00"
	var row []string
	//最多读10次
	for i := 0 ; i < readParmNum ; i ++ {
		source := sendCommand(applianceid,2,groupNum)
		// 1.1 按逗号分割
		row1 := strings.Split(source, ",")
		if len(row1) != 48 {
			fmt.Println("风机参数查询数据失败，查询次数：",i+1)
			time.Sleep(2 * time.Second)//2秒后重新查
		} else {
			fmt.Println("风机参数查询成功")
			row = row1
			break
		}
	}
	if len(row) != 48 {
		fmt.Println("设备参数查询失败，风机参数修改退出~~~~")
		return 0 ,0 , 0
	}
	// 1.2 将字符串转化为数字
	var num []int64
	//var str []string
	for i := 13; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	fmt.Println("设备风机参数：",num)
	/*判断设备风机类型*/
	fanType := struct {
		FanType string
	}{}
	db.Table("keyparamaters").Where("SN = ?",devType).Scan(&fanType)
	var indexMax , indexMin int
	if fanType.FanType == "正弦波风机" {
		indexMax = 2
		indexMin = 3
	} else if fanType.FanType == "方波风机" {
		indexMax = 0
		indexMin = 1
	} else {
		fmt.Println("风机类型错误，风机参数修改退出~~~~")
		return 0 , 0,0
	}

	/*2.判断设备自学习参数有无小于等于85*/
	maxDutyCycCoeff := num[indexMax]
	minDutyCycCoeff := num[indexMin]
	sendMaxFlag := 0
	sendMinFlag := 0

	if maxDutyCycCoeff > 85 {
		num[indexMax] -= 5
		sendMaxFlag = 1
	}
	if minDutyCycCoeff > 85 {
		num[indexMin] -= 5
		sendMinFlag = 1
	}
	/*3.参数下发*/
	if ((sendMinFlag == 1) || (sendMaxFlag==1))&&rewrinteflag=="1" {
		
		// 1.3 将数字变为十六进制字符串,得到参数值
		var str []string
		for _, v := range num {
			temp := strconv.FormatInt(v, 16)
			str = append(str, temp)
		}
		//向设备发送命令，改写参数
		//重复修改100次，不成功则退出
		var numDstF []int64
		for i := 0 ; i < sendParmNum ; i++ {
			sendCommand(applianceid, 4, str[0], str[1], str[2], str[3])
			//查询结果，验证有无修改成功
			sourceDst := sendCommand(applianceid,2,groupNum)
			rowDst := strings.Split(sourceDst, ",")
			if len(rowDst) != 48 {
				fmt.Println("回复无效数据,风机参数重新下发，次数：",i+1)
				time.Sleep(2 * time.Second)//2秒后重新查询
				continue
			}
			// 1.2 将字符串转化为数字
			var numDst []int64
			//var str []string
			for i := 13; i <= 18; i++ {
				temp, _ := strconv.ParseInt(rowDst[i], 16, 64)
				numDst = append(numDst, temp)
			}
			if (numDst[indexMax] == num[indexMax]) && (numDst[indexMin] == num[indexMin]) {
				numDstF = numDst
				fmt.Println("风机参数改写成功")
				res = 1
				break //改写成功则退出
			}
			//不成功则等2秒，再发一次
			fmt.Println("改写未成功,参数重新下发，次数：",i+1)
			time.Sleep(2 * time.Second)
		}
		if len(numDstF) == 0 {
			fmt.Println("回复参数获取失败,风机参数修改退出~~~~")
			return 0 , 0,0
		}
		/*4.结果记录*/
		var perchangeNum PerchangeNums
		changeFlag := 0
		perchangeNum.DevId = applianceid
		perchangeNum.UpdataTime = time.Now().Format("2006-01-02 15:04:05")
		if sendMaxFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMax]
			perchangeNum.UpValue = strconv.Itoa(int(maxDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDstF[indexMax]))
			perchangeNum.FaultType = "c4"
			if (numDstF[indexMax] == num[indexMax]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
				fmt.Println("改写未成功")
			}
			err := db.Create(&perchangeNum).Error
			if err != nil {
				fmt.Println("风机参数改写记录存储未成功~~~")
			} else {
				fmt.Println("风机参数改写记录存储成功~~~")
			}

		}
		if sendMinFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMin]
			perchangeNum.UpValue = strconv.Itoa(int(minDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDstF[indexMin]))
			perchangeNum.FaultType = "c4"
			if (numDstF[indexMin] == num[indexMin]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
				fmt.Println("改写未成功")
			}
			err := db.Create(&perchangeNum).Error
			if err != nil {
				fmt.Println("风机参数改写记录存储未成功~~~")
			} else {
				fmt.Println("风机参数改写记录存储成功~~~")
			}
		}
		if changeFlag == 1 {
			var numPerChange NumPerChanges
			err := db.Table("num_per_changes").
				Where("dev_id = ? AND fault_type= ?", applianceid, "C4").First(&numPerChange).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				numPerChange.DevId = applianceid
				numPerChange.FaultType = "C4"
				numPerChange.DataTime = time.Now().Format("2006-01-02 15:04:05")
				numPerChange.ChangeNum = "1"
				db.Create(&numPerChange)
			} else {
				ChangeNum , _ := strconv.Atoi(numPerChange.ChangeNum)
				ChangeNum++
				db.Table("num_per_changes").
					Where("dev_id = ? AND fault_type= ?", applianceid, "C4").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":           time.Now().Format("2006-01-02 15:04:05"),
						"change_num":          strconv.Itoa(ChangeNum),
					})
			}
		}
		return
	} else {
		/*需要上门设备输出*/
		var failedDevice FailedDevice
		failedDevice.DevId = applianceid
		failedDevice.DevType = devType
		failedDevice.FaultType = "c4"
		failedDevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedDevice.ReviseFlag = 1
		err := db.Create(failedDevice).Error
		if err != nil {
			fmt.Println("风机故障设备存储未成功~~~")
		} else {
			fmt.Println("风机故障设备存储成功~~~")
		}
		return 2 ,num[indexMax],num[indexMin]
	}
}


/*=========================================================
 * 函数名称： FanParmquery
 * 功能描述:  风机参数需要修改值后端接口，涉及到风堵系数超限修改与C4修改
 * 返回参数： 
 * 创建日期： 2022.07.22
 * 作    者： 李赞辉

 =========================================================*/
 func FanParmquery(applianceid , devType,flag string) (nm []int64,sendMin ,sendMax ,indMax,indMin int ,minDutyC ,maxDutyC int64){
	db := common.GetIndexDB()
	fmt.Println("-----------------C4~~风机故障参数修改进入-----------------")
	fmt.Println("修改设备：",applianceid)
	defer fmt.Println("-----------------C4~~风机故障参数修改退出-----------------")
	var readParmNum = 3 //设备参数读取次数上限
	
	/*1.读取设备风机自学习参数*/
	groupNum := "00"
	var row []string
	//最多读10次
	for i := 0 ; i < readParmNum ; i ++ {
		source := sendCommand(applianceid,2,groupNum)
		// 1.1 按逗号分割
		row1 := strings.Split(source, ",")
		if len(row1) != 48 {
			fmt.Println("风机参数查询数据失败，查询次数：",i+1)
			time.Sleep(2 * time.Second)//2秒后重新查
		} else {
			fmt.Println("风机参数查询成功")
			row = row1
			break
		}
	}
	if len(row) != 48 {
		fmt.Println("设备参数查询失败，风机参数修改退出~~~~")
		return 
	}
	// 1.2 将字符串转化为数字
	var num []int64
	//var str []string
	for i := 13; i <= 18; i++ {
		temp, _ := strconv.ParseInt(row[i], 16, 64)
		num = append(num, temp)
	}
	fmt.Println("设备风机参数：",num)
	/*判断设备风机类型*/
	fanType := struct {
		FanType string
	}{}
	db.Table("keyparamaters").Where("SN = ?",devType).Scan(&fanType)
	var indexMax , indexMin int
	if fanType.FanType == "正弦波风机" {
		indexMax = 2
		indexMin = 3
	} else if fanType.FanType == "方波风机" {
		indexMax = 0
		indexMin = 1
	} else {
		fmt.Println("风机类型错误，风机参数修改退出~~~~")
		return 
	}

	/*2.判断设备自学习参数有无小于等于90*/
	maxDutyCycCoeff := num[indexMax]
	minDutyCycCoeff := num[indexMin]
	sendMaxFlag := 0
	sendMinFlag := 0
if flag=="c4"{
	if maxDutyCycCoeff >=110{
		num[indexMax] = 100
		sendMaxFlag = 1
	}else 	if maxDutyCycCoeff > 90 &&maxDutyCycCoeff <110{
		num[indexMax] -= 5
		sendMaxFlag = 1
		if num[indexMax]<90{
			num[indexMax]=90
			sendMaxFlag = 1
		}
	}else if maxDutyCycCoeff<=90{
		sendMaxFlag = 0
	}
	if minDutyCycCoeff >= 110 {
		num[indexMin] = 100
		sendMinFlag = 1
	} else if minDutyCycCoeff > 90&&minDutyCycCoeff < 110 {
		num[indexMin] -= 5
		sendMinFlag = 1
	if num[indexMin]<90{
		num[indexMin] = 90
		sendMinFlag =1 
	}
}else if minDutyCycCoeff<=90{
	sendMinFlag =0
}
}else if flag=="e1"{
	if maxDutyCycCoeff > 110 {
		num[indexMax] = 100
		sendMaxFlag = 1
	}
	if minDutyCycCoeff > 110 {
		num[indexMin] = 100
		sendMinFlag = 1
	}
}
	return num,sendMinFlag,sendMaxFlag,indexMax,indexMin,minDutyCycCoeff,maxDutyCycCoeff

 }

 /*=========================================================
 * 函数名称： FanParmquery
 * 功能描述:  风机参数需要修改值后端接口
 * 返回参数： 
 * 创建日期： 2022.05.10
 * 修改时间： 2022.05.25
 =========================================================*/

 func  FanParmchange(db *gorm.DB,num []int64 , sendMinFlag ,sendMaxFlag ,indexMax,indexMin int,applianceid,devType string,minDutyCycCoeff ,maxDutyCycCoeff int64)(res int){
	var sendParmNum = 3 //设备参数修改次数上限
	groupNum := "00"
	if ((sendMinFlag == 1) || (sendMaxFlag==1)){
		
		// 1.3 将数字变为十六进制字符串,得到参数值
		var str []string
		for _, v := range num {
			temp := strconv.FormatInt(v, 16)
			str = append(str, temp)
		}
		//向设备发送命令，改写参数
		//重复修改100次，不成功则退出
		var numDstF []int64
		for i := 0 ; i < sendParmNum ; i++ {
			sendCommand(applianceid, 4, str[0], str[1], str[2], str[3])
			//查询结果，验证有无修改成功
			sourceDst := sendCommand(applianceid,2,groupNum)
			rowDst := strings.Split(sourceDst, ",")
			if len(rowDst) != 48 {
				fmt.Println("回复无效数据,风机参数重新下发，次数：",i+1)
				time.Sleep(2 * time.Second)//2秒后重新查询
				continue
			}
			// 1.2 将字符串转化为数字
			var numDst []int64
			//var str []string
			for i := 13; i <= 18; i++ {
				temp, _ := strconv.ParseInt(rowDst[i], 16, 64)
				numDst = append(numDst, temp)
			}
			if (numDst[indexMax] == num[indexMax]) && (numDst[indexMin] == num[indexMin]) {
				numDstF = numDst
				fmt.Println("风机参数改写成功")
				res = 1
				break //改写成功则退出
			}
			//不成功则等2秒，再发一次
			fmt.Println("改写未成功,参数重新下发，次数：",i+1)
			time.Sleep(2 * time.Second)
		}
		if len(numDstF) == 0 {
			fmt.Println("回复参数获取失败,风机参数修改退出~~~~")
			return  0
		}
		/*4.结果记录*/
		var perchangeNum PerchangeNums
		changeFlag := 0
		perchangeNum.DevId = applianceid
		perchangeNum.UpdataTime = time.Now().Format("2006-01-02 15:04:05")
		if sendMaxFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMax]
			perchangeNum.UpValue = strconv.Itoa(int(maxDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDstF[indexMax]))
			perchangeNum.FaultType = "c4"
			if (numDstF[indexMax] == num[indexMax]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
				fmt.Println("改写未成功")
			}
			err := db.Create(&perchangeNum).Error
			if err != nil {
				fmt.Println("风机参数改写记录存储未成功~~~")
			} else {
				fmt.Println("风机参数改写记录存储成功~~~")
			}

		}
		if sendMinFlag == 1 {
			perchangeNum.Pername = model.NonParaNmae[indexMin]
			perchangeNum.UpValue = strconv.Itoa(int(minDutyCycCoeff))
			perchangeNum.CurreValue = strconv.Itoa(int(numDstF[indexMin]))
			perchangeNum.FaultType = "c4"
			if (numDstF[indexMin] == num[indexMin]) {
				perchangeNum.SuccessFlag = "1"
				changeFlag = 1
			} else {
				perchangeNum.SuccessFlag = "0"
				fmt.Println("改写未成功")
			}
			err := db.Create(&perchangeNum).Error
			if err != nil {
				fmt.Println("风机参数改写记录存储未成功~~~")
			} else {
				fmt.Println("风机参数改写记录存储成功~~~")
			}
		}
		if changeFlag == 1 {
			var numPerChange NumPerChanges
			err := db.Table("num_per_changes").
				Where("dev_id = ? AND fault_type= ?", applianceid, "C4").First(&numPerChange).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				numPerChange.DevId = applianceid
				numPerChange.FaultType = "C4"
				numPerChange.DataTime = time.Now().Format("2006-01-02 15:04:05")
				numPerChange.ChangeNum = "1"
				db.Create(&numPerChange)
			} else {
				ChangeNum , _ := strconv.Atoi(numPerChange.ChangeNum)
				ChangeNum++
				db.Table("num_per_changes").
					Where("dev_id = ? AND fault_type= ?", applianceid, "C4").
					Updates(map[string]interface{}{ //有记录，更新记录
						"data_time":           time.Now().Format("2006-01-02 15:04:05"),
						"change_num":          strconv.Itoa(ChangeNum),
					})
			}
		}
		return
	} else {
		/*需要上门设备输出*/
		var failedDevice FailedDevice
		failedDevice.DevId = applianceid
		failedDevice.DevType = devType
		failedDevice.FaultType = "c4"
		failedDevice.DataTime = time.Now().Format("2006-01-02 15:04:05")
		failedDevice.ReviseFlag = 1
		err := db.Create(failedDevice).Error
		if err != nil {
			fmt.Println("风机故障设备存储未成功~~~")
		} else {
			fmt.Println("风机故障设备存储成功~~~")
		}
		return 2 
	}
 }
