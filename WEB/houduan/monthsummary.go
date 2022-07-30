package houduan

import (
"fmt"
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/mysql"
"strconv"
"time"
)

func MonthsSummary1(start string, end string, AutoFlag bool, applianceIdTableName string, Database string) {
	db, err := gorm.Open("mysql", Database) //数据库连接
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if AutoFlag == false {
		db.DropTable("month_summaries")
	}
	db.AutoMigrate(&MonthSummary{})

	var e []Equipment

	db.Table(applianceIdTableName).Find(&e)
	var day []DaysSummary
	timeFormat := "2006-01-02"
	startTime, _ := time.Parse(timeFormat, start) //开始时间
	endTime, _ := time.Parse(timeFormat, end)     //结束时间
	var summaryMonth [1]MonthSummary
	for endTime.After(startTime) {
		time1 := startTime
		time2 := startTime.AddDate(0, 1, 0) //年指标统计=>月指标统计=>日指标统计
		for i, n := 0, len(e); i < n; i++ {

			id := e[i].DevId
			db.Table("days_summaries").Where("time_date >= ? AND time_date< ? And dev_id=? ", time1, time2, id).Find(&day)
			if len(day) != 0 {
				summaryMonth = MonthsSummary(day, time1)
				db.Create(&summaryMonth[0])
			}
		}
		startTime = time2
	}
}

//表名

func MonthsSummary(day []DaysSummary, dataTime time.Time) [1]MonthSummary {
	var month [1]MonthSummary
	var validTime time.Duration
	var patternNum int
	var averageTime time.Duration
	var unStableTime time.Duration
	var stableTime time.Duration
	var MaxTime []time.Duration
	var MinTime []time.Duration
	var unStableBehavior int
	var abnormalFlag int
	var ScoreSum int
	var three int
	var five int
	var ten int
	for i, n := 0, len(day); i < n; i++ {
		//有效时长
		valid, _ := time.ParseDuration(day[i].ValidTime)
		validTime += valid
		//用水次数
		patternNum += day[i].PatternNum
		//不稳定时长
		unStable, _ := time.ParseDuration(day[i].UnStableTime)
		unStableTime += unStable
		//不稳定用水次数
		unStableBehavior += day[i].UnStableBehavior
		//最大用水时长
		max, _ := time.ParseDuration(day[i].MaximumTime)
		MaxTime = append(MaxTime, max)
		//最小用水时长
		min, _ := time.ParseDuration(day[i].MinimumTime)
		MinTime = append(MinTime, min)
		//水流量评价
		ScoreSum += day[i].FlowMultipleScore * day[i].PatternNum
		//时长统计
		//3-5分钟用水次数
		three += day[i].Three
		//5-10分钟用水次数
		five += day[i].Five
		//10分钟以上用水次数
		ten += day[i].Ten
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
	month[0].ProvinceCode = day[0].ProvinceCode
	month[0].CityCode = day[0].CityCode
	month[0].DevType = day[0].DevType
	month[0].DevId = day[0].DevId
	month[0].TimeDate = dataTime.String()[0:10]
	month[0].ValidTime = validTime.String()
	month[0].PatternNum = patternNum
	month[0].UnStableProportion = unStableProportion
	month[0].StableTime = stableTime.String()
	month[0].AverageTime = averageTime.String()
	month[0].UnStableTime = unStableTime.String()
	month[0].MaximumTime = MaxAndMinimumTime1(MaxTime, true).String()
	month[0].MinimumTime = MaxAndMinimumTime1(MinTime, false).String()
	month[0].AbnormalFlag = abnormalFlag
	month[0].UnStableBehavior = unStableBehavior
	month[0].FlowMultipleScore = ScoreSum / patternNum
	month[0].Three = three
	month[0].Five = five
	month[0].Ten = ten

	return month
}
