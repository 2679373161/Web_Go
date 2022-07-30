package houduan

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"time"
)

func Summary(start string, end string, AutoFlag bool, Day bool, Database string) {
	//连接数据库
	db, err := gorm.Open("mysql", Database) //数据库连接
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if AutoFlag == true {
		db.DropTable("province_months")
		db.DropTable("city_months")
		db.DropTable("type_months")
	} else {
		if Day == true {
			db.DropTable("province_days")
			db.DropTable("city_days")
			db.DropTable("type_days")
		} else {
			db.DropTable("province_months")
			db.DropTable("city_months")
			db.DropTable("type_months")
		}
	}
	var sheetName string
	var year int
	var month int
	var day int
	var Province string
	var City string
	var Type string
	if Day == true {
		sheetName = "days_summaries"
		year = 0
		month = 0
		day = 1
		Province = "province_days"
		City = "city_days"
		Type = "type_days"
	} else {
		sheetName = "month_summaries"
		year = 0
		month = 1
		day = 0
		Province = "province_months"
		City = "city_months"
		Type = "type_months"
	}

	var m []Multimodality
	db.AutoMigrate(ProvinceName(Province))
	db.AutoMigrate(CityName(City))
	db.AutoMigrate(TypeName(Type))
	//db.Table(sheetName).Find(&m)
	//城市信息汇总
	citySummaryFlag := true
	//省份信息汇总
	provinceSummaryFlag := true
	//型号信息汇总
	typeSummaryFlag := true

	category := sheetName //表格名字
	timeFormat := "2006-01-02"
	startTime, _ := time.Parse(timeFormat, start) //开始时间
	endTime, _ := time.Parse(timeFormat, end)     //结束时间
	for endTime.After(startTime) {
		time1 := startTime
		time2 := startTime.AddDate(year, month, day) //年指标统计=>月指标统计=>日指标统计

		db.Table(category).Where("Time_date >= ? AND Time_date< ? ", time1, time2).Find(&m)
		if citySummaryFlag == true {
			out := UrbanSummary(m, time1, timeFormat)
			for i, n := 0, len(out); i < n; i++ {
				db.Table(City).Create(&out[i])
			}
		}
		if provinceSummaryFlag == true {
			Pro := ProvinceSummary(m, time1, timeFormat)
			for i, n := 0, len(Pro); i < n; i++ {
				db.Table(Province).Create(&Pro[i])
			}
		}
		if typeSummaryFlag == true {
			Machine := TypeSummary(m, time1, timeFormat)
			for i, n := 0, len(Machine); i < n; i++ {
				db.Table(Type).Create(&Machine[i])
			}
		}

		startTime = time2
	}
}

/*=========================================================
 * 功能描述: 对指标统计按城市进行汇总
 * 输出指标: 输出MultimodalCity结构体
 =========================================================*/

func UrbanSummary(m []Multimodality, time1 time.Time, timeFormat string) []MultimodalityCity {
	var out []MultimodalityCity
	var mul [1]MultimodalityCity
	for i1, n := 0, len(cityNumber); i1 < n; i1++ {
		var number int
		var UnStableNumber int
		var AbnormalFlag int
		var avgTime float64
		var unStableTimePro float64
		var avg time.Duration
		var FlowScoreSum int
		var validTimeSum time.Duration
		var stableTimeSum time.Duration
		var unStableTimeSum time.Duration
		var MaximumTimeSlice []time.Duration
		var MinimumTimeSlice []time.Duration
		var equipmentNum int
		//var s2 time.Duration
		for i, n := 0, len(m); i < n; i++ {
			if cityNumber[i1] == m[i].CityCode {
				//用水事件个数
				number += m[i].PatternNum
				//水流不稳定个数
				UnStableNumber += m[i].UnStableBehavior
				//有效时间
				validTime, _ := time.ParseDuration(m[i].ValidTime)
				validTimeSum += validTime
				//最长用水时间
				MaximumTime, _ := time.ParseDuration(m[i].MaximumTime)
				MaximumTimeSlice = append(MaximumTimeSlice, MaximumTime)
				//最短用水时间
				MinimumTime, _ := time.ParseDuration(m[i].MinimumTime)
				MinimumTimeSlice = append(MinimumTimeSlice, MinimumTime)
				//稳定时长
				stableTime, _ := time.ParseDuration(m[i].StableTime)
				stableTimeSum += stableTime
				//不稳定时长
				UnStableTime, _ := time.ParseDuration(m[i].UnStableTime)
				unStableTimeSum += UnStableTime
				//设备数量
				equipmentNum++
				//水流波动评价
				FlowScoreSum += m[i].FlowMultipleScore * m[i].PatternNum
			}
		}
		//平均时间
		conversionValue := time.Duration(number) * time.Second
		if conversionValue != 0 {
			avgTime = validTimeSum.Seconds() / conversionValue.Seconds()
			avg = time.Duration(avgTime) * time.Second
		} else {
			avgTime = 0
		}
		//稳定占比
		if validTimeSum != 0 {
			unStableTimePro, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", unStableTimeSum.Seconds()/validTimeSum.Seconds()), 64)
		} else {
			unStableTimePro = 0
		}
		//异常标志位
		if unStableTimePro >= 0.5 {
			AbnormalFlag = 1
		} else {
			AbnormalFlag = 0
		}
		mul[0].ProvinceCode = ProvinceCodeName1(cityNumber[i1])
		mul[0].CityCode = cityNumber[i1]
		mul[0].EquipmentNum = equipmentNum
		mul[0].TimeDate = time1.Format(timeFormat)
		mul[0].ValidTime = validTimeSum.String()
		mul[0].PatternNum = number
		mul[0].AverageTime = avg.String()
		mul[0].UnStableProportion = PercentConversion(unStableTimePro)
		mul[0].MaximumTime = MaxAndMinimumTime1(MaximumTimeSlice, true).String()
		mul[0].MinimumTime = MaxAndMinimumTime1(MinimumTimeSlice, false).String()
		mul[0].StableTime = stableTimeSum.String()
		mul[0].UnStableTime = unStableTimeSum.String()
		mul[0].UnStableBehavior = UnStableNumber
		mul[0].AbnormalFlag = AbnormalFlag
		if number != 0 {
			mul[0].FlowMultipleScore = FlowScoreSum / number
			out = append(out, mul[0])
		}

		fmt.Printf("城市：%v汇总完成\n", cityNumber[i1])
	}

	return out
}

/*=========================================================
 * 功能描述: 对指标统计按省进行汇总
 * 输出指标: 输出MultimodalCity结构体
 =========================================================*/

func ProvinceSummary(m []Multimodality, time1 time.Time, timeFormat string) []MultimodalityProvince {

	var province []MultimodalityProvince
	var mul2 [1]MultimodalityProvince
	for i1, n := 0, len(provinceNumber); i1 < n; i1++ {
		var number int
		var UnStableNumber int
		var AbnormalFlag int
		var avgTime float64
		var unStableTimePro float64
		var avg time.Duration
		var FlowScoreSum int
		var validTimeSum time.Duration
		var stableTimeSum time.Duration
		var unStableTimeSum time.Duration
		var MaximumTimeSlice []time.Duration
		var MinimumTimeSlice []time.Duration
		var equipmentNum int
		var three int
		var five int
		var ten int
		//var s2 time.Duration
		for i, n := 0, len(m); i < n; i++ {
			if provinceNumber[i1] == m[i].ProvinceCode {
				//用水事件个数
				number += m[i].PatternNum
				//水流不稳定个数
				UnStableNumber += m[i].UnStableBehavior
				//有效时间
				validTime, _ := time.ParseDuration(m[i].ValidTime)
				validTimeSum += validTime
				//最长用水时间
				MaximumTime, _ := time.ParseDuration(m[i].MaximumTime)
				MaximumTimeSlice = append(MaximumTimeSlice, MaximumTime)
				//最短用水时间
				MinimumTime, _ := time.ParseDuration(m[i].MinimumTime)
				MinimumTimeSlice = append(MinimumTimeSlice, MinimumTime)
				//稳定时长
				stableTime, _ := time.ParseDuration(m[i].StableTime)
				stableTimeSum += stableTime
				//不稳定时长
				UnStableTime, _ := time.ParseDuration(m[i].UnStableTime)
				unStableTimeSum += UnStableTime
				//设备数量
				equipmentNum++
				//水流波动评价
				FlowScoreSum += m[i].FlowMultipleScore * m[i].PatternNum
				//3-5分钟用水次数
				three += m[i].Three
				//5-10分钟用水次数
				five += m[i].Five
				//10分钟以上用水次数
				ten += m[i].Ten
			}
		}
		//平均时间
		conversionValue := time.Duration(number) * time.Second
		if conversionValue != 0 {
			avgTime = validTimeSum.Seconds() / conversionValue.Seconds()
			avg = time.Duration(avgTime) * time.Second
		} else {
			avgTime = 0
		}
		//稳定占比
		if validTimeSum != 0 {
			unStableTimePro, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", unStableTimeSum.Seconds()/validTimeSum.Seconds()), 64)
		} else {
			unStableTimePro = 0
		}
		//异常标志位
		if unStableTimePro >= 0.5 {
			AbnormalFlag = 1
		} else {
			AbnormalFlag = 0
		}
		mul2[0].ProvinceCode = provinceNumber[i1]
		mul2[0].EquipmentNum = equipmentNum
		mul2[0].TimeDate = time1.Format(timeFormat)
		mul2[0].ValidTime = validTimeSum.String()
		mul2[0].PatternNum = number
		mul2[0].AverageTime = avg.String()
		mul2[0].UnStableProportion = PercentConversion(unStableTimePro)
		mul2[0].MaximumTime = MaxAndMinimumTime1(MaximumTimeSlice, true).String()
		mul2[0].MinimumTime = MaxAndMinimumTime1(MinimumTimeSlice, false).String()
		mul2[0].StableTime = stableTimeSum.String()
		mul2[0].UnStableTime = unStableTimeSum.String()
		mul2[0].UnStableBehavior = UnStableNumber
		mul2[0].AbnormalFlag = AbnormalFlag
		mul2[0].Three = three
		mul2[0].Five = five
		mul2[0].Ten = ten
		if number != 0 {
			mul2[0].FlowMultipleScore = FlowScoreSum / number
			province = append(province, mul2[0])
		}
		fmt.Printf("省：%v汇总完成\n", provinceNumber[i1])
	}

	return province
}

/*=========================================================
 * 功能描述: 对指标统计按型号进行汇总
 * 输出指标: 输出MultimodalCity结构体
 =========================================================*/

func TypeSummary(m []Multimodality, time1 time.Time, timeFormat string) []MultimodalityType {
	var Machine []MultimodalityType
	var mul1 [1]MultimodalityType
	for i1, n := 0, len(typeNumber); i1 < n; i1++ {
		var number int
		var UnStableNumber int
		var AbnormalFlag int
		var avgTime float64
		var unStableTimePro float64
		var avg time.Duration
		var FlowScoreSum int
		var validTimeSum time.Duration
		var stableTimeSum time.Duration
		var unStableTimeSum time.Duration
		var MaximumTimeSlice []time.Duration
		var MinimumTimeSlice []time.Duration
		var equipmentNum int
		//var s2 time.Duration
		for i, n := 0, len(m); i < n; i++ {
			if typeNumber[i1] == m[i].DevType {
				//用水事件个数
				number += m[i].PatternNum
				//水流不稳定个数
				UnStableNumber += m[i].UnStableBehavior
				//有效时间
				validTime, _ := time.ParseDuration(m[i].ValidTime)
				validTimeSum += validTime
				//最长用水时间
				MaximumTime, _ := time.ParseDuration(m[i].MaximumTime)
				MaximumTimeSlice = append(MaximumTimeSlice, MaximumTime)
				//最短用水时间
				MinimumTime, _ := time.ParseDuration(m[i].MinimumTime)
				MinimumTimeSlice = append(MinimumTimeSlice, MinimumTime)
				//稳定时长
				stableTime, _ := time.ParseDuration(m[i].StableTime)
				stableTimeSum += stableTime
				//不稳定时长
				UnStableTime, _ := time.ParseDuration(m[i].UnStableTime)
				unStableTimeSum += UnStableTime
				//设备数量
				equipmentNum++
				//水流波动评价
				FlowScoreSum += m[i].FlowMultipleScore * m[i].PatternNum
			}
		}
		//平均时间
		conversionValue := time.Duration(number) * time.Second
		if conversionValue != 0 {
			avgTime = validTimeSum.Seconds() / conversionValue.Seconds()
			avg = time.Duration(avgTime) * time.Second
		} else {
			avgTime = 0
		}
		//稳定占比
		if validTimeSum != 0 {
			unStableTimePro, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", unStableTimeSum.Seconds()/validTimeSum.Seconds()), 64)
		} else {
			unStableTimePro = 0
		}
		//异常标志位
		if unStableTimePro >= 0.5 {
			AbnormalFlag = 1
		} else {
			AbnormalFlag = 0
		}
		mul1[0].DevType = typeNumber[i1]
		mul1[0].EquipmentNum = equipmentNum
		mul1[0].TimeDate = time1.Format(timeFormat)
		mul1[0].ValidTime = validTimeSum.String()
		mul1[0].PatternNum = number
		mul1[0].AverageTime = avg.String()
		mul1[0].UnStableProportion = PercentConversion(unStableTimePro)
		mul1[0].StableTime = stableTimeSum.String()
		mul1[0].MaximumTime = MaxAndMinimumTime1(MaximumTimeSlice, true).String()
		mul1[0].MinimumTime = MaxAndMinimumTime1(MinimumTimeSlice, false).String()
		mul1[0].UnStableTime = unStableTimeSum.String()
		mul1[0].UnStableBehavior = UnStableNumber
		mul1[0].AbnormalFlag = AbnormalFlag
		if number != 0 {
			mul1[0].FlowMultipleScore = FlowScoreSum / number
			Machine = append(Machine, mul1[0])
		}
		fmt.Printf("机型：%v汇总完成\n", typeNumber[i1])
	}

	return Machine
}

/*=========================================================
 * 功能描述: 时间类型转换（Time.time=>Time.Duration）
 * 输出指标: 输出Time.Duration类型
 =========================================================*/

func TimeTypeConversion(originalTime string) time.Duration {
	t0, _ := time.Parse("15:04:05", "00:00:00")
	t, _ := time.Parse("15:04:05", originalTime)
	finalTime := t.Sub(t0)
	return finalTime
}

/*=========================================================
 * 功能描述: 将计算得占比转换为百分数
 * 输出指标: 输出string类型
 =========================================================*/

func PercentConversion(originalValue float64) string {
	originalValue = originalValue * 100
	finalValue := strconv.FormatFloat(originalValue, 'f', 2, 64)
	return finalValue
}

/*=========================================================
 * 功能描述: 最大值最小值求解
 * 输出指标: 最大值或最小值（time,Duration）
 * 函数调用：标志位true=>最大值，标志位false=>最小值
 =========================================================*/

func MaxAndMinimumTime1(MaxAndMinSlice []time.Duration, Flag bool) time.Duration {
	var MaxAndMinTime time.Duration
	for i, n := 0, len(MaxAndMinSlice); i < n; i++ {
		if MaxAndMinSlice[i] != 0 {
			MaxAndMinTime = MaxAndMinSlice[i]
			break
		}
	}
	if MaxAndMinSlice != nil {
		if Flag == true {
			for i, n := 0, len(MaxAndMinSlice); i < n; i++ {
				if MaxAndMinSlice[i] > MaxAndMinTime {
					MaxAndMinTime = MaxAndMinSlice[i]
				}
			}
		}
		if Flag == false {
			for i, n := 0, len(MaxAndMinSlice); i < n; i++ {
				if (MaxAndMinSlice[i] < MaxAndMinTime) && (MaxAndMinSlice[i] != 0) {
					MaxAndMinTime = MaxAndMinSlice[i]
				}
			}
		}
	}

	return MaxAndMinTime
}