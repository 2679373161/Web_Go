package houduan

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"math"
	"strconv"
	"time"
)

/*=========================================================
 * 功能描述: 数据挖掘
 * 输出指标: 输出指标
 =========================================================*/

func DataMining(dataStartTime string, dataEndTime string, indexWriteFlag bool, dataWriteFlag bool, fragmentWriteFlag bool, years int, months int, days int, applianceIdTableName string, Database string) {
	//连接mysql数据库
	db, err := gorm.Open("mysql", Database) //数据库连接
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var e []Equipment
	//数据处理输入参数

	yearsFlag, monthsFlag, DayFlag, timeFormat1 := StatisticsRangeSetting(years, months, days)
	db.Table(applianceIdTableName).Find(&e)

	for i, n := 0, len(e); i < n; i++ {

		id := e[i].DevId
		devtype := e[i].DevType
		citycode := e[i].CityCode
		applianceCityCode := citycode
		dataSheetName := OperationDataSheetName(applianceCityCode)
		fragmentTable := PatternFragmentTableName(applianceCityCode)
		db.AutoMigrate(NewUser2(fragmentTable))
		db.AutoMigrate(NewUser1("MultiModality"))
		outputTable := DataOutPutTableName(dataSheetName)
		db.AutoMigrate(NewUser(outputTable))
		//结构体索引
		var u []user
		var m [1]Multimodality
		var t time.Time

		//设备ID表名（input）
		//db.Table(applianceIdTableName).Where("City_code=?", applianceCityCode).Find(&e)

		category := dataSheetName //表格名字
		timeFormat := "2006-01-02"
		//timeFormat1 := "2006-01-02"
		startTime, _ := time.Parse(timeFormat, dataStartTime) //开始时间
		endTime, _ := time.Parse(timeFormat, dataEndTime)     //结束时间

		for endTime.After(startTime) {
			time1 := startTime
			time2 := startTime.AddDate(yearsFlag, monthsFlag, DayFlag) //年指标统计=>月指标统计=>日指标统计
			if time2.Sub(endTime) > 0 {
				time2 = endTime
			}
			db.Table(category).Where("Datatime >= ? AND Datatime< ? And applianceId=? ", time1, time2, id).Find(&u)

			repeatPointClean(u)
			uu := timeCut(u, 60*time.Second, 180*time.Second)
			frag6, number6, sum6, maxTime6, minTime6, avgFlowScore6, three6, five6, ten6 := feature(uu, 6, 1)
			frag7, number7, sum7, maxTime7, minTime7, avgFlowScore7, three7, five7, ten7 := feature(uu, 7, 1)
			frag8, number8, sum8, maxTime8, minTime8, avgFlowScore8, three8, five8, ten8 := feature(uu, 8, 1)
			if fragmentWriteFlag == true {
				for i, n := 0, len(frag6); i < n; i++ {
					db.Table(fragmentTable).Create(&frag6[i])
				}
				for i, n := 0, len(frag7); i < n; i++ {
					db.Table(fragmentTable).Create(&frag7[i])
				}
				for i, n := 0, len(frag8); i < n; i++ {
					db.Table(fragmentTable).Create(&frag8[i])
				}
			}
			/*=========================================================
			 * 功能描述: 计算水流量多模式指标
			 * 输出指标: 输出水流量多模式数据到MySQL数据表
			 =========================================================*/
			var UnStableRatio float64
			effectiveTime := sum6 + sum7 + sum8

			effectiveNumber := number6 + number7 + number8
			conversionValue := time.Duration(effectiveNumber) * time.Second
			fmt.Println(effectiveTime, effectiveNumber)
			m[0].DevId = id
			m[0].ProvinceCode = ProvinceCodeName(citycode)
			m[0].TimeDate = time1.Format(timeFormat1)
			m[0].DevType = devtype
			m[0].CityCode = citycode
			m[0].ValidTime = effectiveTime.String()
			if (number6 + number7 + number8) != 0 {
				m[0].FlowMultipleScore = (avgFlowScore6*number6 + avgFlowScore7*number7 + avgFlowScore8*number8) / (number6 + number7 + number8)
			}
			t.Add(effectiveTime).Format("15:04:05")
			m[0].PatternNum = effectiveNumber
			m[0].StableTime = sum6.String()
			t.Add(sum6).Format("15:04:05")
			m[0].UnStableTime = (sum7 + sum8).String()
			t.Add(sum7 + sum8).Format("15:04:05")
			m[0].UnStableBehavior = number7 + number8
			m[0].MaximumTime = MaxAndMinimumTime(maxTime6, maxTime7, maxTime8, true).String()
			m[0].MinimumTime = MaxAndMinimumTime(minTime6, minTime7, minTime8, false).String()
			m[0].Three = three6 + three7 + three8
			m[0].Five = five6 + five7 + five8
			m[0].Ten = ten6 + ten7 + ten8
			if conversionValue != 0 {
				avgTime, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", effectiveTime.Seconds()/conversionValue.Seconds()), 64)
				avgTimeDuration := time.Duration(avgTime) * time.Second
				m[0].AverageTime = avgTimeDuration.String()
			} else {
				m[0].AverageTime = "0s"
			}
			if effectiveTime != 0 {
				UnStableRatio, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", 100*((sum7+sum8).Seconds()/effectiveTime.Seconds())), 64)
				m[0].UnStableProportion = UnStableRatio
				fmt.Printf("%v,%v,%v\n", id, time1.Format(timeFormat1), UnStableRatio)
			} else {
				m[0].UnStableProportion = 0
			}
			if UnStableRatio >= 0.5 {
				m[0].AbnormalFlag = 1
			} else {
				m[0].AbnormalFlag = 0
			}
			if indexWriteFlag == true {
				if effectiveTime != 0 {
					db.Table("MultiModality").Create(&m[0])
				}
			}
			startTime = time2
			if dataWriteFlag == true {
				for i, n := 0, len(uu); i < n; i++ {
					db.Table(outputTable).Create(&uu[i])
				}
			}
		}

		fmt.Printf("当前第%v台\n", i+1)

	}

}

/*=========================================================
 * 功能描述: 数据挖掘
 * 输出指标: 输出指标
 =========================================================*/

func MaxAndMinimumTime(Time1 time.Duration, Time2 time.Duration, Time3 time.Duration, Flag bool) time.Duration {
	var MaxAndMinArray = [...]time.Duration{Time1, Time2, Time3}
	var MaxAndMinTime time.Duration
	for i, n := 0, len(MaxAndMinArray); i < n; i++ {
		if MaxAndMinArray[i] != 0 {
			MaxAndMinTime = MaxAndMinArray[i]
			break
		}
	}
	if Flag == true {
		for i, n := 0, len(MaxAndMinArray); i < n; i++ {
			if MaxAndMinArray[i] > MaxAndMinTime {
				MaxAndMinTime = MaxAndMinArray[i]
			}
		}
	}
	if Flag == false {
		for i, n := 0, len(MaxAndMinArray); i < n; i++ {
			if (MaxAndMinArray[i] < MaxAndMinTime) && (MaxAndMinArray[i] != 0) {
				MaxAndMinTime = MaxAndMinArray[i]
			}
		}
	}
	return MaxAndMinTime
}

/*=========================================================
 * 功能描述: 数据挖掘
 * 输出指标: 输出指标
 =========================================================*/

func StatisticsRangeSetting(yearsFlag int, monthsFlag int, daysFlag int) (int, int, int, string) {
	var timeFormat1 string

	if yearsFlag != 0 {
		timeFormat1 = "2006"
	} else if monthsFlag != 0 {
		timeFormat1 = "2006-01"
	} else {
		timeFormat1 = "2006-01-02"
	}

	return yearsFlag, monthsFlag, daysFlag, timeFormat1
}

/*=========================================================
 * 功能描述: 逻辑错误清洗
 * 输出指标: 输出指标
 =========================================================*/

func repeatPointClean(u []user) []user {

	var datatime []time.Time //时间轴
	for i, n := 0, len(u); i < n; i++ {
		t, _ := time.Parse("2006-01-02 15:04:05", u[i].Datatime)
		datatime = append(datatime, t) //数据放入切片中
	}
	for i, n := 1, len(u)-1; i < n; i++ {
		deltaT := datatime[i].Sub(datatime[i-1])
		if ((deltaT <= 0*time.Second) || (datatime[i] == datatime[i+1])) && (u[i+1].Applianceid == u[i-1].Applianceid) {
			datatime[i] = datatime[i-1].Add(1 * time.Second)
			u[i].Datatime = datatime[i].Format("2006-01-02 15:04:05")
		}
	}

	return u
}

/*=========================================================
 * 功能描述: 用水行为切分
 * 输出指标: 输出指标
 =========================================================*/

func timeCut(u []user, gapTime time.Duration, zoneTime time.Duration) []user {
	var length int
	var addT time.Duration
	var datatime []string //时间轴
	var u3 = make([]user, 0, len(u))
	var t1m int

	var t3m int

	for _, u := range u {
		datatime = append(datatime, u.Datatime) //数据放入切片中
	}
	for i, n := 0, len(u)-1; i < n; i++ {
		length++
		t2, _ := time.Parse("2006-01-02 15:04:05", datatime[i+1])
		t1, _ := time.Parse("2006-01-02 15:04:05", datatime[i])
		deltaT := t2.Sub(t1)
		if deltaT > gapTime {

			if addT < zoneTime {
				//过滤无效区间
				t1m++
			} else {
				if addT >= 180*time.Second {
					t3m++
				}
				//插值
				u0 := insertPoint(u[i-length+1 : i+1])
				u1 := outletWaterTemperaturePattern(u0)
				u2 := patternRec(u1) //有效区间进入模式识别函数,左闭右开
				zoneID(u2)
				u3 = append(u3, u2...) //将模式识别结果存放进一个切片u1中
			}
			length = 0
			addT = 0
		} else if i == n-1 { //最后一个区间
			if deltaT < time.Duration(10)*time.Second {
				addT += deltaT
			}
			if addT < zoneTime {
				//过滤无效区间
				t1m++
			} else {
				if addT >= 180*time.Second {
					t3m++
				}
				u0 := insertPoint(u[i-length+1 : i+2])
				u1 := outletWaterTemperaturePattern(u0)
				u2 := patternRec(u1) //有效区间进入模式识别函数,左闭右开
				zoneID(u2)
				u3 = append(u3, u2...) //将模式识别结果存放进一个切片u1中
			}
		} else {
			if deltaT < time.Duration(10)*time.Second {
				addT += deltaT
			}
		}
	}
	fmt.Println(t1m, t3m)
	return u3
}

/*=========================================================
 * 功能描述: 删除火焰为0的
 * 输出指标: 输出指标
 =========================================================*/

func deletePoint(u []user) []user {
	var udelete []user

	for i, n := 0, len(u); i < n; i++ {
		if u[i].Flame != 0 {
			udelete = append(udelete, u[i])
		}
	}
	return udelete
}

/*=========================================================
 * 功能描述: 插值
 * 输出指标: 输出指标
 =========================================================*/

func insertPoint(u []user) []user {
	var t1 time.Time
	var t2 time.Time
	us := make([]user, len(u), len(u))
	copy(us, u)

	var i int = 0
	for {
		if i == len(us)-1 {
			break
		}
		t1, _ = time.Parse("2006-01-02 15:04:05", us[i].Datatime)
		t2, _ = time.Parse("2006-01-02 15:04:05", us[i+1].Datatime)
		deltaT := t2.Sub(t1)
		if deltaT >= time.Duration(4)*time.Second && deltaT <= time.Duration(10)*time.Second {
			ui := us[i]
			ui.Datatime = t1.Add(2 * time.Second).Format("2006-01-02 15:04:05")
			us = append(us[:i+1], append([]user{ui}, us[i+1:]...)...)
		}
		i++
	}
	return us
}

/*=========================================================
 * 功能描述: 模式划分总
 * 输出指标: 输出指标
 =========================================================*/

func patternRec(u []user) []user {
	var deltaflow []float64
	for i, n := 0, len(u)-1; i < n; i++ {
		deltaflow = append(deltaflow, float64(u[i+1].Flow-u[i].Flow)/float64(u[i].Flow))
	}

	deltaflow[0] = 0

	deltaflow[len(deltaflow)-1] = 0

	priorityDivision(u, deltaflow, 0.1, 2, 3, 10) //优先级1，代码为10
	priority1Classification(u, deltaflow, 0.1, 10, 2, 3, 5)
	u1 := roughClassification(u, 8, 7, 6) //粗分类震荡7、阶跃8、稳定6

	return u1
}

/*=========================================================
 * 功能描述: 模式划分
 * 输出指标: 输出指标
 =========================================================*/

func roughClassification(u []user, model1 int, model2 int, model3 int) []user {
	var u1 []user
	var firstu user
	var lastu user
	var shockFlag bool
	var stepFlag bool
	for i, n := 0, len(u); i < n; i++ {
		if u[i].Model == 5 {
			shockFlag = true
		}
		if (u[i].Model == 2) || (u[i].Model == 3) {
			stepFlag = true
		}
	}
	//不稳定
	if shockFlag == true {
		for i, n := 0, len(u); i < n; i++ {
			u[i].Model = model1
		}
	} else if stepFlag == true {
		for i, n := 0, len(u); i < n; i++ {
			u[i].Model = model2

		}
	} else {
		for i, n := 0, len(u); i < n; i++ {
			u[i].Model = model3
		}
	}
	firstu = u[0]
	firstu.Flame = 0
	firstu.Flow = 0
	t0, _ := time.Parse("2006-01-02 15:04:05", firstu.Datatime)
	firstu.Datatime = t0.Add(-1 * time.Second).Format("2006-01-02 15:04:05")

	lastu = u[len(u)-1]
	lastu.Flame = 0
	lastu.Flow = 0
	t1, _ := time.Parse("2006-01-02 15:04:05", lastu.Datatime)
	lastu.Datatime = t1.Add(1 * time.Second).Format("2006-01-02 15:04:05")

	u1 = append(u1, firstu)
	u1 = append(u1, u...)
	u1 = append(u1, lastu)

	return u1
}

/*=========================================================
 * 功能描述: 用水行为标记
 * 输出指标: 输出指标
 =========================================================*/

func zoneID(u []user) []user {

	for i, n := 0, len(u); i < n; i++ {
		u[i].ZoneID = i + 1
	}

	return u
}

/*=========================================================
 * 功能描述: 优先级划分
 * 输出指标: 输出指标
 =========================================================*/

func priorityDivision(u []user, judgment []float64, threshold float64, M int, N int, model int) []user {
	for i, n := 0, len(judgment); i < n; i++ {
		if (math.Abs(float64(judgment[i])) >= threshold) && (u[i].Model == 0 || u[i].Model == model) {
			if i >= M && (len(u)-i-1) >= N { //前后点数都够
				for i, n := i-M, i+N; i <= n; i++ {
					if u[i].Model == 0 {
						u[i].Model = model
					}
				}
			} else if i < M { //前面不够点数
				for i, n := 0, M+N; i <= n; i++ {
					if u[i].Model == 0 {
						u[i].Model = model
					}
				}
			} else if (len(u) - i - 1) < N { //后面不够点数
				for i, n := len(u)-1-M-N, len(u)-1; i <= n; i++ {
					if u[i].Model == 0 {
						u[i].Model = model
					}
				}
			}
		} //根据流量差值确定优先级，无else不涉及部分不操作
	} //遍历切片
	return u
}

/*=========================================================
 * 功能描述: 优先级细分
 * 输出指标: 输出指标
 =========================================================*/

func priority1Classification(u []user, judgment []float64, threshold float64, model int, model1 int, model2 int, model3 int) []user {
	var number = 0 //片段长度
	var up = 0
	var down = 0
	for i, n := 0, len(judgment); i < n; i++ {
		if u[i].Model == model {
			number++
			if judgment[i] >= threshold {
				up++
			} else if judgment[i] <= -threshold {
				down++
			} else {
			}
			if i == len(judgment)-1 { //有效区间末尾
				if up > 0 && down == 0 {
					for i1, n1 := len(judgment)-number, len(judgment); i1 <= n1; i1++ {
						u[i1].Model = model1
					}
				} else if down > 0 && up == 0 {
					for i1, n1 := len(judgment)-number, len(judgment); i1 <= n1; i1++ {
						u[i1].Model = model2
					}
				} else if down > 0 && up > 0 {
					for i1, n1 := len(judgment)-number, len(judgment); i1 <= n1; i1++ {
						u[i1].Model = model3
					}
				} else {
				} //其他优先级
			}
		} else {
			//fmt.Printf("总的%v,%v,%v\n", up, down, number)
			if up > 0 && down == 0 {
				for i1, n1 := i-number, i; i1 < n1; i1++ {
					u[i1].Model = model1
				}
			} else if down > 0 && up == 0 {
				for i1, n1 := i-number, i; i1 < n1; i1++ {
					u[i1].Model = model2
				}
			} else if down > 0 && up > 0 {
				for i1, n1 := i-number, i; i1 < n1; i1++ {
					u[i1].Model = model3
				}
			} else {
			} //其他优先级
			number = 0
			up = 0
			down = 0
		}
	}
	return u
}

/*=========================================================
 * 功能描述: 特征指标
 * 输出指标: 输出指标
 =========================================================*/

func feature(u []user, model int, gapNumber int) ([]ModeFragment, int, time.Duration, time.Duration, time.Duration, int, int, int, int) {
	var frag []ModeFragment
	var number int
	var sum time.Duration
	var max time.Duration
	var min time.Duration
	var avgFlowScore int
	var length int
	var timeLen []time.Duration
	var multipleScoreSlice []int
	var three int
	var five int
	var ten int
	for i, n := 0, len(u)-1; i < n; i++ {
		length++
		if u[i+1].ZoneID-u[i].ZoneID != gapNumber {
			//fmt.Printf("有效区间--长度:%v,结束id:%v,节点时间差值:%v，区间总时间:%v\n", length, i+1, deltaT, addT)
			time1 := zoneTime(u[i-length+1:i+1], model)
			if time1 != 0 {
				timeLen = append(timeLen, time1) //有效区间进入时间特征统计函数,左闭右开;将结果存放进一个切片u1中
			}
			extreme, maxChange, avg, deviation := featureFlow(u[i-length+1:i+1], model)
			up, down := flowChangeCalculate(u[i-length+1:i+1], 6)

			if u[i-length+1].Model == model {
				three, five, ten = WaterDuration(time1, three, five, ten)
				heatDuration, unStableTempDuration, unStableTempPercent, unHeatDev, heatDev, tempNum := outletWaterFeature(u[i-length+1:i+1], time1)
				Coefficient := deviation[0] / avg[0]
				multipleScore := FlowFluctuationEvaluation(Coefficient, model, up, down)
				s := ModeFragmentIndexWrite(u[i-length+1:i+1], frag, model, time1, extreme, maxChange, avg, deviation, up, down, multipleScore, heatDuration, unStableTempDuration, unStableTempPercent, unHeatDev, heatDev, tempNum)
				frag = append(frag, s[0])
				multipleScoreSlice = append(multipleScoreSlice, multipleScore)
			}
			length = 0
		} else if i == n-1 { //最后一个区间
			//fmt.Printf("最后有效区间--长度:%v,结束id:%v,节点时间差值:%v，区间总时间:%v\n", length+1, i+2, deltaT, addT)
			time2 := zoneTime(u[i-length+1:i+2], model)
			if time2 != 0 {
				timeLen = append(timeLen, time2) //有效区间进入时间特征统计函数,左闭右开;将结果存放进一个切片u1中
			}
			extreme, maxChange, avg, deviation := featureFlow(u[i-length+1:i+2], model)
			up, down := flowChangeCalculate(u[i-length+1:i+2], 6)

			if u[i-length+1].Model == model {
				three, five, ten = WaterDuration(time2, three, five, ten)
				heatDuration, unStableTempDuration, unStableTempPercent, unHeatDev, heatDev, tempNum := outletWaterFeature(u[i-length+1:i+2], time2)
				Coefficient := deviation[0] / avg[0]
				multipleScore := FlowFluctuationEvaluation(Coefficient, model, up, down)
				s := ModeFragmentIndexWrite(u[i-length+1:i+2], frag, model, time2, extreme, maxChange, avg, deviation, up, down, multipleScore, heatDuration, unStableTempDuration, unStableTempPercent, unHeatDev, heatDev, tempNum)
				frag = append(frag, s[0])
				multipleScoreSlice = append(multipleScoreSlice, multipleScore)
			}
		} else {
		}
	}

	if len(timeLen) != 0 { //防止空数组（无该模式的情况）
		number, sum, max, min = featureTimeCalculate(timeLen)
		avgFlowScore = fluctuationAvg(multipleScoreSlice)
	}
	return frag, number, sum, max, min, avgFlowScore, three, five, ten
}

/*=========================================================
 * 功能描述: 模式片段写入
 * 输出指标: 输出指标
 =========================================================*/

func ModeFragmentIndexWrite(u []user, frag []ModeFragment, model int, timeLen time.Duration, extreme []int, maxChange []float64, avg []float64, deviation []float64, up int, down int, multipleScore int, heatDuration time.Duration, unStableTempDuration time.Duration, unStableTempPercent float64, unHeatDev float64, heatDev float64, tempNum int) [1]ModeFragment {
	var f [1]ModeFragment
	maxChangeValue, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", maxChange[0]), 64)
	averageValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg[0]), 64)
	deviationValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", deviation[0]), 64)
	unHeatDev, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", unHeatDev), 64)
	heatDev, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", heatDev), 64)
	f[0].DevId = u[0].Applianceid
	f[0].Pattern = model
	f[0].StartTime = u[0].Datatime
	f[0].EndTime = u[len(u)-1].Datatime
	f[0].DurationTime = timeLen.String()
	f[0].Extreme = extreme[0]
	f[0].MaxChange = maxChangeValue
	f[0].Average = averageValue
	f[0].Deviation = deviationValue
	f[0].UpNumber = up
	f[0].DownNumber = down
	f[0].MultipleScore = multipleScore
	f[0].HeatDuration = heatDuration.String()
	f[0].UnStableTempDuration = unStableTempDuration.String()
	f[0].UnStableTempPercent = unStableTempPercent
	f[0].UnHeatDev = unHeatDev
	f[0].HeatDev = heatDev
	f[0].TempNum = tempNum

	return f
}

/*=========================================================
 * 功能描述: 用水时长计算
 * 输出指标: 输出指标
 =========================================================*/

func WaterDuration(timeLen time.Duration, three int, five int, ten int) (int, int, int) {
	if timeLen >= time.Duration(600)*time.Second {
		ten++
	} else if timeLen >= time.Duration(300)*time.Second {
		five++
	} else {
		three++
	}
	return three, five, ten
}

/*=========================================================
 * 功能描述: 水流量波动评价
 * 输出指标: 输出指标
 =========================================================*/

func FlowFluctuationEvaluation(Coefficient float64, model int, up int, down int) int {
	var StableMark = 100.0
	var StepMark = 80.0
	var shockMark = 60.0
	var variationScore float64
	var modelScore float64
	var mutationScore float64
	var multipleScore int
	//变异得分
	variationScore = 100 - Coefficient*200 //变异系数得分=100-变异系数×200
	if variationScore < 0 {
		variationScore = 0
	}
	//模式得分
	if model == 6 {
		modelScore = StableMark
	} else if model == 7 {
		modelScore = StepMark
	} else if model == 8 {
		modelScore = shockMark
	}
	//突变次数得分
	sum := up + down
	mutationScore = 100.0 - float64(5*sum)
	if mutationScore < 0 {
		mutationScore = 0
	}
	//综合得分
	multipleScore = int(variationScore*0.4 + modelScore*0.4 + mutationScore*0.2)

	return multipleScore
}

/*=========================================================
 * 功能描述: 区间时长计算
 * 输出指标: 输出指标
 =========================================================*/

func zoneTime(u []user, model int) time.Duration {
	var timeLen time.Duration

	if u[0].Model == model {
		startTime, _ := time.Parse("2006-01-02 15:04:05", u[0].Datatime)
		endTime, _ := time.Parse("2006-01-02 15:04:05", u[len(u)-1].Datatime)
		timeLen = endTime.Sub(startTime)
	}
	return timeLen
}

/*=========================================================
 * 功能描述: 时间特征指标计算
 * 输出指标: 输出指标
 =========================================================*/

func featureTimeCalculate(timeLen []time.Duration) (number int, sum time.Duration, max time.Duration, min time.Duration) {
	max = timeLen[0]
	min = timeLen[0]
	for i, n := 0, len(timeLen); i < n; i++ {
		sum += timeLen[i]
		if timeLen[i] > max {
			max = timeLen[i]
		} else if timeLen[i] < min {
			min = timeLen[i]
		}
	}
	number = len(timeLen)
	return
}

/*=========================================================
 * 功能描述: 平均值计算
 * 输出指标: 输出指标
 =========================================================*/

func fluctuationAvg(flowAvg []int) int {
	var sum int
	for i, n := 0, len(flowAvg); i < n; i++ {
		sum += flowAvg[i]
	}

	avgFlowScore := sum / len(flowAvg)

	return avgFlowScore
}

/*=========================================================
 * 功能描述: 水流量变化量计算
 * 输出指标: 输出指标
 =========================================================*/

func flowChangeCalculate(u []user, threshold int) (up int, down int) {
	up = 0
	down = 0
	for i, n := 2, len(u)-3; i < n; i++ {

		if u[i+1].Flow-u[i].Flow >= threshold {
			up++
		} else if u[i+1].Flow-u[i].Flow <= (-threshold) {
			down++
		} else {

		}
	}
	return
}

/*=========================================================
 * 功能描述: 水流量特征
 * 输出指标: 输出指标
 =========================================================*/

func featureFlow(u []user, model int) (flowExtreme []int, flowMaxChange []float64, flowAvg []float64, flowDeviation []float64) {
	var number = 0
	var deltaflow []int
	for i, n := 0, len(u)-1; i < n; i++ {
		deltaflow = append(deltaflow, u[i+1].Flow-u[i].Flow)
	}

	for i, n := 1, len(u)-1; i < n; i++ {
		if u[i].Model == model {
			number++
			if i == n-1 { //区间末尾模式片段
				max := u[i-number+1].Flow
				min := u[i-number+1].Flow
				var sumFlow = 0
				var sumFF = 0.0
				maxChange := math.Abs(float64(deltaflow[i-number+1]))
				for i1, n1 := i-number+1, i+1; i1 < n1; i1++ {
					sumFlow += u[i1].Flow
					if u[i1].Flow > max {
						max = u[i1].Flow
					} else if u[i1].Flow < min {
						min = u[i1].Flow
					}
				}
				for i1, n1 := i-number+1, i+1; i1 < n1; i1++ {
					sumFF += (float64(sumFlow)/float64(number) - float64(u[i1].Flow)) * (float64(sumFlow)/float64(number) - float64(u[i1].Flow))
				}
				flowExtreme = append(flowExtreme, max-min)
				avg, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sumFlow)/float64(number)), 64)
				flowAvg = append(flowAvg, avg)
				deviation, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Sqrt(sumFF/float64(number))), 64)
				flowDeviation = append(flowDeviation, deviation)
				for i1, n1 := i-number+1, i; i1 < n1; i1++ { //水流量差值
					if math.Abs(float64(deltaflow[i1])) > math.Abs(maxChange) {
						maxChange = float64(deltaflow[i1])
					}
				}
				flowMaxChange = append(flowMaxChange, maxChange)
			} else {
			}
		} else {
			if number != 0 {
				max := u[i-number].Flow
				min := u[i-number].Flow
				var sumFlow int
				var sumFF float64
				maxChange := math.Abs(float64(deltaflow[i-number]))
				for i1, n1 := i-number, i; i1 < n1; i1++ { //水流量
					sumFlow += u[i1].Flow
					if u[i1].Flow > max {
						max = u[i1].Flow
					} else if u[i1].Flow < min {
						min = u[i1].Flow
					}
				}
				for i1, n1 := i-number, i; i1 < n1; i1++ {
					sumFF += (float64(sumFlow)/float64(number) - float64(u[i1].Flow)) * (float64(sumFlow)/float64(number) - float64(u[i1].Flow))
				}
				flowExtreme = append(flowExtreme, max-min)
				avg, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sumFlow)/float64(number)), 64)
				flowAvg = append(flowAvg, avg)
				deviation, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Sqrt(sumFF/float64(number))), 64)
				flowDeviation = append(flowDeviation, deviation)

				for i1, n1 := i-number, i-1; i1 < n1; i1++ { //水流量差值
					if math.Abs(float64(deltaflow[i1])) > math.Abs(maxChange) {
						maxChange = float64(deltaflow[i1])
					}
				}
				flowMaxChange = append(flowMaxChange, maxChange)
			}
			number = 0
		}
	}

	return
}

/*=========================================================
 * 功能描述: 温度过程计算
 * 输出指标: 输出指标
 =========================================================*/

func outletWaterTemperaturePattern(u []user) []user {
	//升温过程判定
	heatingUpPattern(u, 1)
	//超温过程判定
	overTemperaturePattern(u, 3)
	//超温关火判定
	shutdownPattern(u, 4)
	//回温过程判定
	underTemperaturePattern(u, 5)
	//恒温过程判定
	return u
}

/*=========================================================
 * 功能描述: 出水温度指标
 * 输出指标: 输出指标
 =========================================================*/

func outletWaterFeature(u []user, totalDuration time.Duration) (time.Duration, time.Duration, float64, float64, float64, int) {

	//升温过程时长
	heatDuration := tempFeatureTime(u, 1)
	//开机过冲时长
	overShootDuration := tempFeatureTime(u, 2)
	//超温过程时长
	overTempDuration := tempFeatureTime(u, 3)
	//超温关火时长
	//shutdownDuration := tempFeatureTime(u, 4)
	//回温过程时长
	underTempDuration := tempFeatureTime(u, 5)

	//不恒温时长
	unStableTempDuration := overTempDuration + underTempDuration
	//不恒温占比
	unStableTempPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", 100*(unStableTempDuration.Seconds()/totalDuration.Seconds())), 64)
	//温度模式
	tempNum := tempPattern(heatDuration, overShootDuration, overTempDuration, underTempDuration)
	///相对目标温度的标准差
	unHeatDev, heatDev := tempFeatureDeviation(u)
	return heatDuration, unStableTempDuration, unStableTempPercent, unHeatDev, heatDev, tempNum
}

/*=========================================================
 * 功能描述: 温度模式判定
 * 输出指标: 输出指标
 =========================================================*/
func tempPattern(heatDuration time.Duration, overShootDuration time.Duration, overTempDuration time.Duration, underTempDuration time.Duration) int {
	var heatNum int
	var constantNum int
	//升温阶段判断
	if heatDuration+overShootDuration > time.Duration(35)*time.Second && overShootDuration > time.Duration(0)*time.Second {
		heatNum = 40
	} else if heatDuration+overShootDuration > time.Duration(35)*time.Second {
		heatNum = 20
	} else if overShootDuration > time.Duration(0)*time.Second {
		heatNum = 30
	} else {
		heatNum = 10
	}
	//恒温阶段判断
	if overTempDuration > time.Duration(0)*time.Second && underTempDuration > time.Duration(0)*time.Second {
		constantNum = 4
	} else if overTempDuration > time.Duration(0)*time.Second {
		constantNum = 2
	} else if underTempDuration > time.Duration(0)*time.Second {
		constantNum = 3
	} else {
		constantNum = 1
	}
	tempNum := heatNum + constantNum
	return tempNum
}

/*=========================================================
 * 功能描述: 温度特征指标
 * 输出指标: 输出指标
 =========================================================*/

func tempFeatureDeviation(u []user) (float64, float64) {
	var sumUnHeat int
	var sumHeat int
	var numberUnHeat int
	var numberHeat int
	var unHeatDeviation float64
	var heatDeviation float64
	for i, n := 0, len(u); i < n; i++ {
		if u[i].TempPattern != 1 {
			numberUnHeat++
			sumUnHeat += (u[i].Outtemp - u[i].Settemp) * (u[i].Outtemp - u[i].Settemp)
		} else {
			numberHeat++
			sumHeat += (u[i].Outtemp - u[i].Settemp) * (u[i].Outtemp - u[i].Settemp)
		}
	}
	if numberUnHeat != 0 {
		unHeatDeviation = math.Sqrt(float64(sumUnHeat / numberUnHeat))
	}
	if numberHeat != 0 {
		heatDeviation = math.Sqrt(float64(sumHeat / numberHeat))
	}
	return unHeatDeviation, heatDeviation
}

/*=========================================================
 * 功能描述: 时长计算函数
 * 输出指标: 输出指标
 =========================================================*/

func tempFeatureTime(u []user, model int) time.Duration {
	var number = 0
	var totalTime time.Duration
	for i, n := 0, len(u); i < n; i++ {
		if u[i].TempPattern == model {
			number++
			if i == n-1 { //区间末尾模式片段
				t2, _ := time.Parse("2006-01-02 15:04:05", u[i].Datatime)
				t1, _ := time.Parse("2006-01-02 15:04:05", u[i-number+1].Datatime)
				totalTime += t2.Sub(t1)
			} else {
			} //不做处理
		} else {
			if number != 0 {
				t2, _ := time.Parse("2006-01-02 15:04:05", u[i-1].Datatime) //模式间隔时间u[i-1].Datetime-->u[i]
				t1, _ := time.Parse("2006-01-02 15:04:05", u[i-number].Datatime)
				totalTime += t2.Sub(t1)
			}
			number = 0
		}
	}
	return totalTime
}

///升温过程
func heatingUpPattern(u []user, pattern int) []user {
	var temperatureDiff1 int
	var temperatureDiff2 int
	var threshold0 = -2
	var threshold1 = 2
	var heatingUpNumber int
	for i, n := 0, len(u)-1; i < n; i++ {
		temperatureDiff1 = u[i].Outtemp - u[i].Settemp
		temperatureDiff2 = u[i+1].Outtemp - u[i+1].Settemp
		if (temperatureDiff1 < threshold0 || temperatureDiff1 > threshold1) || (temperatureDiff2 < threshold0 || temperatureDiff2 > threshold1) {
			heatingUpNumber++
			if i == n-1 {
				for i1, n1 := 0, i+1; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
			StartupOvershoot(u[0:i+1], 2)
		} else {
			if heatingUpNumber != 0 {
				for i1, n1 := 0, heatingUpNumber; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
			StartupOvershoot(u[0:heatingUpNumber], 2)
			heatingUpNumber = 0
			break
		}
	}
	return u
}

/*=========================================================
 * 功能描述: 开机过冲计算
 * 输出指标: 输出指标
 =========================================================*/

func StartupOvershoot(u []user, pattern int) []user {
	var temperatureDiff1 int
	var temperatureDiff2 int
	var threshold = 4
	var number int
	for i, n := 0, len(u)-1; i < n; i++ {
		temperatureDiff1 = u[i].Outtemp - u[i].Settemp
		temperatureDiff2 = u[i+1].Outtemp - u[i+1].Settemp
		if temperatureDiff1 > threshold && temperatureDiff2 > threshold {
			number++
			if i == n-1 {
				for i1, n1 := i-number+1, i+1; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
		} else {
			if number != 0 {
				for i1, n1 := i-number, i; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
			number = 0
			break
		}
	}
	return u
}

/*=========================================================
 * 功能描述: 超温过程计算
 * 输出指标: 输出指标
 =========================================================*/

func overTemperaturePattern(u []user, pattern int) []user {
	var temperatureDiff1 int
	var temperatureDiff2 int
	//var threshold0 = -2
	var threshold1 = 2
	var overTemperatureNumber int
	for i, n := 0, len(u)-1; i < n; i++ {
		temperatureDiff1 = u[i].Outtemp - u[i].Settemp
		temperatureDiff2 = u[i+1].Outtemp - u[i+1].Settemp
		if (temperatureDiff1 > threshold1) && (temperatureDiff2 > threshold1) && (u[i].TempPattern != 1) {
			overTemperatureNumber++
			if i == n-1 {
				for i1, n1 := i-overTemperatureNumber+1, i+1; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
		} else {
			if overTemperatureNumber != 0 {
				for i1, n1 := i-overTemperatureNumber, i; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
			overTemperatureNumber = 0
		}
	}
	return u

}

/*=========================================================
 * 功能描述: 关火
 * 输出指标: 输出指标
 =========================================================*/

func shutdownPattern(u []user, pattern int) []user {
	var temperatureDiff1 int
	var temperatureDiff2 int
	//var threshold0 = -2
	var threshold1 = 2
	var shutdownNumber int
	for i, n := 0, len(u)-1; i < n; i++ {
		temperatureDiff1 = u[i].Outtemp - u[i].Settemp
		temperatureDiff2 = u[i+1].Outtemp - u[i+1].Settemp
		if (temperatureDiff1 > threshold1) && (temperatureDiff2 > threshold1) && (u[i].Flame == 0) {
			shutdownNumber++
			if i == n-1 {
				for i1, n1 := i-shutdownNumber+1, i+1; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
		} else {
			if shutdownNumber != 0 {
				for i1, n1 := i-shutdownNumber, i; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
			shutdownNumber = 0
		}
	}
	return u
}

/*=========================================================
 * 功能描述: 回温判定
 * 输出指标: 输出指标
 =========================================================*/

func underTemperaturePattern(u []user, pattern int) []user {
	var temperatureDiff1 int
	var temperatureDiff2 int
	var threshold0 = -2
	//var threshold1 = 2
	var underTemperatureNumber int
	for i, n := 0, len(u)-1; i < n; i++ {
		temperatureDiff1 = u[i].Outtemp - u[i].Settemp
		temperatureDiff2 = u[i+1].Outtemp - u[i+1].Settemp
		if (temperatureDiff1 < threshold0) && (temperatureDiff2 < threshold0) && (u[i].TempPattern == 0) {
			underTemperatureNumber++
			if i == n-1 {
				for i1, n1 := i-underTemperatureNumber+1, i+1; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
		} else {
			if underTemperatureNumber != 0 {
				for i1, n1 := i-underTemperatureNumber, i; i1 <= n1; i1++ {
					u[i1].TempPattern = pattern
				}
			} else {
			}
			underTemperatureNumber = 0
		}
	}
	return u
}