package controller

import (
	"fmt"
	"ginEssential/common"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
)

/*线程处理分配管理线程------------>>>>>>查询20参数*/
func TaskAssignmenttwenty(thread int, equ []MultipleParaRewrite, db *gorm.DB) {

	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)

	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)

	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go workersumtwenty(jobs, results, equ, db)
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
func workersumtwenty(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite, db *gorm.DB) {
	
	for i := range jobs {

		QueryParaSettingtwenty(equ[i].ApplianceId,db)

		results <- i
		fmt.Println("线程i:", i)
	}
}







/*线程处理分配管理线程------------>>>>>>恒温改写参数线程分配任务*/
func ParaTaskAssignment(thread int, equ []MultipleParaRewrite, groupNum string, ka string, kb string, kc string, kf string, T1a string, T1c string, T2a string, T2c string, Tda string, Tdc string, Wc string, Wo string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine	
	for w := 1; w <= thread; w++ {
		go workerPara( jobs, results, equ,groupNum, ka, kb, kc, kf, T1a, T1c, T2a, T2c, Tda, Tdc, Wc, Wo)
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

/*分配工作*/
func workerPara(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite,groupNum string, ka string, kb string, kc string, kf string, T1a string, T1c string, T2a string, T2c string, Tda string, Tdc string, Wc string, Wo string) {

	for i := range jobs {
		BasendCommandPara(equ[i].ApplianceId,groupNum, ka, kb, kc, kf, T1a, T1c, T2a, T2c, Tda, Tdc, Wc, Wo)

		results <- i
		fmt.Println("线程i:", i)
	}
}


/*线程处理分配管理线程------------>>>>>>单个改写参数线程分配任务*/
func TaskAssignment(thread int, equ []MultipleParaRewrite, index,value string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine	
	for w := 1; w <= thread; w++ {
		go workerSingle( jobs, results, equ,index,value)
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

/*分配工作*/
func workerSingle(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite,index,value string) {

	for i := range jobs {
		BasendCommand(equ[i].ApplianceId,index,value)

		results <- i
		fmt.Println("线程i:", i)
	}
}
//Parameter

/*线程处理分配管理线程------------>>>>>>多个改写参数线程分配任务*/
func TaskAssignmentParameter(thread int, equ []MultipleParaRewrite, fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string,scrit string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine	
	for w := 1; w <= thread; w++ {
		go workerParameter( jobs, results, equ,fa, ff, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn,scrit)
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

/*分配工作*/
func workerParameter(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite,fa string, ff string, PH string, FH string, PL string, FL string, dH string, Fd string, CH string, FC string, CA string, nE string, FP string, HS string, Hb string, HE string, HL string, HU string, Fn string,scrit string) {

	for i := range jobs {
		BasendCommandParameter(equ[i].ApplianceId,fa, ff, PH, FH, PL, FL, dH, Fd, CH, FC, CA, nE, FP, HS, Hb, HE, HL, HU, Fn,scrit)

		results <- i
		fmt.Println("线程i:", i)
	}
}


/*线程处理分配管理线程------------>>>>>>改写改写非调试模式参数第一组接口线程分配任务*/
func TaskAssignmentFirst(thread int, equ []MultipleParaRewrite, reWaterFlow,windPressureSensor string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine	
	for w := 1; w <= thread; w++ {
		go workerFirst( jobs, results, equ,reWaterFlow,windPressureSensor)
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

/*分配工作*/
func workerFirst(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite,reWaterFlow,windPressureSensor string) {

	for i := range jobs {
		BasendNoDebugFirst(equ[i].ApplianceId,reWaterFlow,windPressureSensor)

		results <- i
		fmt.Println("线程i:", i)
	}
}


/*线程处理分配管理线程------------>>>>>>改写改写非调试模式参数第二组接口线程分配任务*/
func TaskAssignmentNoDebugSecond(thread int, equ []MultipleParaRewrite, MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine	
	for w := 1; w <= thread; w++ {
		go workerNoDebugSecond( jobs, results, equ,MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff)
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

/*分配工作*/
func workerNoDebugSecond(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite,MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff string) {

	for i := range jobs {
		BasendNoDebugSecond(equ[i].ApplianceId,MaxCurrCoeff, MinCurrCoeff, MaxDutyCycCoeff, MinDutyCycCoeff)

		results <- i
		fmt.Println("线程i:", i)
	}
}


/*线程处理分配管理线程------------>>>>>>改写总参数中单个修改参数协议接口线程分配任务*/
func TaskAssignmentsum(thread int, equ []MultipleParaRewrite, index int, value, indexcanshu , Backups string) {

	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)

	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)

	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go workersum(jobs, results, equ, index, value, indexcanshu,Backups)
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
func workersum(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite, index int, value, indexcanshu ,Backups string) {
	swit := true
	for i := range jobs {

		QueryParaSettingsum(equ[i].ApplianceId, index, value, indexcanshu, swit,Backups)

		results <- i
		fmt.Println("线程i:", i)
	}
}

/*线程处理分配管理线程------------>>>>>>改写总参数中单个修改参数协议接口线程分配任务*/
func TaskAssignmentsumStarent(thread int, equ []MultipleParaRewrite, index int, value, indexcanshu , Backups string) {

	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 10000)
	results := make(chan int, 10000)

	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)

	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go workersumStarent(jobs, results, equ, index, value, indexcanshu,Backups)
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
func workersumStarent(jobs <-chan int, results chan<- int, equ []MultipleParaRewrite, index int, value, indexcanshu , Backups string) {
	swit := true
	for i := range jobs {
		FastartingdownDatathread(equ[i].ApplianceId, 12)
		QueryParaSettingsum(equ[i].ApplianceId, index, value, indexcanshu, swit,Backups)
		FastartingdownDatathread(equ[i].ApplianceId, 13)
		results <- i
		fmt.Println("线程i:", i)
	}
}


func NoNetwork(applianceId string){
	db := common.GetDB()
	fmt.Println("设备无联网")
	db.Table("multiple_para_rewrite").Where("appliance_id = ?",applianceId ).
		Updates(map[string]interface{}{
			"succeedflag":        "0",
			"updatetime":         time.Now().Format("2006-01-02 15:04:05"),
			"handleflag":        "1",
		})
}

func Network(applianceId string){
	db := common.GetDB()
	fmt.Println("成功")
	db.Table("multiple_para_rewrite").Where("appliance_id = ?",applianceId ).
		Updates(map[string]interface{}{
			"succeedflag":        "1",
			"updatetime":         time.Now().Format("2006-01-02 15:04:05"),
			"handleflag":        "1",
		})
}