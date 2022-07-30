package houduan

import (
"fmt"
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/mysql"
"strconv"
"time"
)

func DaySummary(start string, end string, AutoFlag bool, applianceIdTableName string, Database string) {
	db, err := gorm.Open("mysql", Database) //数据库连接
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if AutoFlag == false {
		db.DropTable("days_summaries")
	}
	db.AutoMigrate(&DaysSummary{})

	var e []Equipment

	db.Table(applianceIdTableName).Find(&e)
	var fragment []ModeFragment
	timeFormat := "2006-01-02"
	startTime, _ := time.Parse(timeFormat, start) //开始时间
	endTime, _ := time.Parse(timeFormat, end)     //结束时间

	for endTime.After(startTime) {
		time1 := startTime
		time2 := startTime.AddDate(0, 0, 1) //年指标统计=>月指标统计=>日指标统计
		for i, n := 0, len(e); i < n; i++ {
			category := FragmentTableName(e[i].CityCode)
			id := e[i].DevId
			db.Table(category).Where("start_time BETWEEN ? AND ? And dev_id=? ", time1, time2, id).Find(&fragment)
			if len(fragment) != 0 {
				summaryDays := FragmentSummary(fragment, e[i].CityCode, e[i].DevType)
				db.Create(&summaryDays[0])
			}
		}
		startTime = time2
	}
}

//表名
func FragmentTableName(patternName string) string {
	value := "fragment"
	finalValue := value + patternName
	return finalValue
}

func FragmentSummary(fragment []ModeFragment, cityCode string, typeCode string) [1]DaysSummary {
	var day [1]DaysSummary
	var validTime time.Duration
	var patternNum int
	var averageTime time.Duration
	var unStableTime time.Duration
	var stableTime time.Duration
	var timeLen []time.Duration
	var unStableBehavior int
	var abnormalFlag int
	var ScoreSum int
	var FlowScoreSum int
	var three int
	var five int
	var ten int
	for i, n := 0, len(fragment); i < n; i++ {
		//有效时长
		valid, _ := time.ParseDuration(fragment[i].DurationTime)
		validTime += valid
		//用水次数
		patternNum++
		//不稳定时长
		if fragment[i].Pattern == 7 || fragment[i].Pattern == 8 {
			unStable, _ := time.ParseDuration(fragment[i].DurationTime)
			unStableTime += unStable
			unStableBehavior++
		}
		//最大最小时长
		timeLen = append(timeLen, valid)
		//水流量评价
		ScoreSum += fragment[i].MultipleScore
		//时长统计
		if valid >= time.Duration(10)*time.Minute {
			ten++
		} else if valid >= time.Duration(5)*time.Minute {
			five++
		} else {
			three++
		}

	}
	//平均时长
	conversionValue := time.Duration(patternNum) * time.Second
	avgTime, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", validTime.Seconds()/conversionValue.Seconds()), 64)
	averageTime = time.Duration(avgTime) * time.Second
	//不稳定占比
	unStableProportion, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", 100*(unStableTime.Seconds()/validTime.Seconds())), 64)
	//不稳定标志位
	if unStableProportion > 50 {
		abnormalFlag = 1
	}
	//稳定时长
	stableTime = validTime - unStableTime
	//最大最小时长
	maxTime, minTime := MaxAndMinTime(timeLen)
	//水流量分数
	FlowScoreSum = ScoreSum / patternNum
	day[0].ProvinceCode = ProvinceCode(cityCode)
	day[0].CityCode = cityCode
	day[0].DevType = typeCode
	day[0].DevId = fragment[0].DevId
	day[0].TimeDate = fragment[0].StartTime[0:10]
	day[0].ValidTime = validTime.String()
	day[0].PatternNum = patternNum
	day[0].AverageTime = averageTime.String()
	day[0].UnStableProportion = unStableProportion
	day[0].StableTime = stableTime.String()
	day[0].UnStableTime = unStableTime.String()
	day[0].MaximumTime = maxTime.String()
	day[0].MinimumTime = minTime.String()
	day[0].AbnormalFlag = abnormalFlag
	day[0].UnStableBehavior = unStableBehavior
	day[0].FlowMultipleScore = FlowScoreSum
	day[0].Three = three
	day[0].Five = five
	day[0].Ten = ten
	return day
}

func MaxAndMinTime(timeLen []time.Duration) (max time.Duration, min time.Duration) {
	max = timeLen[0]
	min = timeLen[0]
	for i, n := 0, len(timeLen); i < n; i++ {
		if timeLen[i] > max {
			max = timeLen[i]
		} else if timeLen[i] < min {
			min = timeLen[i]
		}
	}
	return
}
