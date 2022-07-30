package main

import (
	"ginEssential/controller"
	"ginEssential/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	//预警模块下评分分析模块
	scoresummaryRoutes := r.Group("/api/scoresummary")
	ScoreSummaryController := controller.NewscoresummaryController()
	scoresummaryRoutes.PUT("/:id", ScoreSummaryController.Update)
	scoresummaryRoutes.POST("page/search", ScoreSummaryController.Search)
	scoresummaryRoutes.POST("page/scorechart", ScoreSummaryController.Scorechart)
	scoresummaryRoutes.POST("page/intempchart", ScoreSummaryController.Intempchart)
	scoresummaryRoutes.POST("page/residualheatchart", ScoreSummaryController.Residualheatchart)
	scoresummaryRoutes.POST("page/funnelchart", ScoreSummaryController.Funnelchart)
	scoresummaryRoutes.POST("page/gettype", ScoreSummaryController.Gettype)
	scoresummaryRoutes.POST("page/getequipment", ScoreSummaryController.Getequipment)
	//预警模块
	earlywarningRoutes := r.Group("/api/earlywarning")
	earlywarningController := controller.NewEarlyWarningController()
	earlywarningRoutes.POST("page/search", earlywarningController.Search)
	earlywarningRoutes.POST("page/scorechart", earlywarningController.ScoreChart)


	MonitorRoutes := r.Group("/api/monitor")
	MonitorController := controller.NewMonitorController()
	MonitorRoutes.PUT("/:id", MonitorController.Update)
	MonitorRoutes.POST("page/initialize", MonitorController.Datasave)
	MonitorRoutes.POST("page/getequipment", MonitorController.Getequipment)
	MonitorRoutes.POST("page/gettype", MonitorController.Gettype)
	MonitorRoutes.POST("page/Datachart", MonitorController.Datachart)
	MonitorRoutes.POST("page/ceshi", MonitorController.Ceshi)
	MonitorRoutes.POST("page/DataError", MonitorController.Data_Error)

	EqudayRoutes := r.Group("/api/equday")
	EqudayController := controller.NewEqudayController()
	EqudayRoutes.PUT("/:id", EqudayController.Update)
	EqudayRoutes.POST("page/initialize", EqudayController.Datasave)
	EqudayRoutes.POST("page/getequipment", EqudayController.Getequipment)
	EqudayRoutes.POST("page/gettype", EqudayController.Gettype)
	EqudayRoutes.POST("page/Datachart", EqudayController.Datachart)
	EqudayRoutes.POST("page/ceshi", EqudayController.Ceshi)





	faultdiagnosisne01Routes := r.Group("/api/fault")
	faultdiagnosisne01Controller := controller.Newfaultdiagnosisne01Controller()
    faultdiagnosisne01Routes.POST("page/faultdiagnosis", faultdiagnosisne01Controller.Search)
	faultdiagnosisne01Routes.POST("page/faultsummaries", faultdiagnosisne01Controller.Summaries)
	faultdiagnosisne01Routes.POST("page/getequipment", faultdiagnosisne01Controller.Getequipment)
	faultdiagnosisne01Routes.POST("page/gettype", faultdiagnosisne01Controller.Gettype)
	faultdiagnosisne01Routes.POST("page/ceshi", faultdiagnosisne01Controller.Ceshi)
	faultdiagnosisne01Routes.POST("page/getday", faultdiagnosisne01Controller.Getday)
	faultdiagnosisne01Routes.POST("page/overtemp", faultdiagnosisne01Controller.Overtemp)   //超温诊断近期
	faultdiagnosisne01Routes.POST("page/overtempgetdata", faultdiagnosisne01Controller.Overtempgetdata) //界面刷新时选择框的型号设备获取近期
	faultdiagnosisne01Routes.POST("page/overtempgettype", faultdiagnosisne01Controller.Overtempgettype) //选择框的型号获取近期
	faultdiagnosisne01Routes.POST("page/overtempgetid", faultdiagnosisne01Controller.Overtempgetid)  //选择框的设备获取近期
	faultdiagnosisne01Routes.POST("page/overtempday", faultdiagnosisne01Controller.Overtempday) //超温每日详情
	faultdiagnosisne01Routes.POST("page/overtempdaygetdata", faultdiagnosisne01Controller.Overtempdaygetdata) //界面刷新时选择框的型号设备获取每日
	faultdiagnosisne01Routes.POST("page/overtempdaygettype", faultdiagnosisne01Controller.Overtempdaygettype) //选择框的型号获取每日
	faultdiagnosisne01Routes.POST("page/overtempdaygetid", faultdiagnosisne01Controller.Overtempdaygetid)  //选择框的设备获取每日




	faultdiagnosisne2Routes := r.Group("/api/fault2")
	faultdiagnosisne2Controller := controller.Newfaultdiagnosisne2Controller()
    faultdiagnosisne2Routes.POST("page/faultdiagnosis2", faultdiagnosisne2Controller.Search)
	faultdiagnosisne2Routes.POST("page/faultsummaries2", faultdiagnosisne2Controller.Summaries)
	faultdiagnosisne2Routes.POST("page/getequipment2", faultdiagnosisne2Controller.Getequipment)
	faultdiagnosisne2Routes.POST("page/gettype2", faultdiagnosisne2Controller.Gettype)
	faultdiagnosisne2Routes.POST("page/ceshi2", faultdiagnosisne2Controller.Ceshi)
	faultdiagnosisne2Routes.POST("page/getday2", faultdiagnosisne2Controller.Getday)



	basicdatacontrollerRoutes := r.Group("/api/basicdatacontroller")
	basicdatacontrollerController := controller.NewBasicdatacontrollerController()
	basicdatacontrollerRoutes.PUT("/:id", basicdatacontrollerController.Update)
	basicdatacontrollerRoutes.POST("page/typedata", basicdatacontrollerController.Datatype)
	basicdatacontrollerRoutes.POST("page/initialize", basicdatacontrollerController.Datasave)
	basicdatacontrollerRoutes.POST("page/equipmentdata", basicdatacontrollerController.Dataequipment)
	basicdatacontrollerRoutes.POST("page/modifydata", basicdatacontrollerController.Modifydata)
	basicdatacontrollerRoutes.POST("page/modifytype", basicdatacontrollerController.Modifytype)
	basicdatacontrollerRoutes.POST("page/modifyequipment", basicdatacontrollerController.Modifyequipment)

	dataRoutes := r.Group("/api/data")
	dataController := controller.NewDatasaveController()
	dataRoutes.PUT("/:id", dataController.Update)
	dataRoutes.POST("page/initialize", dataController.Datasave)

	faultstatisticsRoutes := r.Group("/api/faultstatistics")
	faultstatisticsController := controller.NewfaultstatisticsController()
	faultstatisticsRoutes.PUT("/:id", faultstatisticsController.Update)
	faultstatisticsRoutes.POST("page/initialize", faultstatisticsController.Search)
    faultstatisticsRoutes.POST("page/gettype", faultstatisticsController.Gettype)
	faultstatisticsRoutes.POST("page/idfault", faultstatisticsController.Getid)
	faultstatisticsRoutes.POST("page/getequ", faultstatisticsController.Getequ)
	faultstatisticsRoutes.POST("page/typesearch", faultstatisticsController.GetTypeErrorInfo)
	faultstatisticsRoutes.POST("page/typesearchsum", faultstatisticsController.GetTypeErrorInfosum)
	faultstatisticsRoutes.POST("page/prodmonthfaultsum", faultstatisticsController.Prodmonthfaultsum)
	faultstatisticsRoutes.POST("page/equfaultnum", faultstatisticsController.Equnum)//型号的设备故障数量占比每日详情
	faultstatisticsRoutes.POST("page/equfaultnumsum", faultstatisticsController.Equnumsum)//型号的设备故障数量占比时间汇总


	recordingRoutes := r.Group("/api/recording")
	recordingController := controller.NewrecordingController()
	recordingRoutes.PUT("/:id", recordingController.Update)

    recordingRoutes.POST("page/gettype", recordingController.Gettype)
	recordingRoutes.POST("page/idfault", recordingController.Getid)
	recordingRoutes.POST("page/getequ", recordingController.Getequ)
	
	collectionRoutes := r.Group("/api/collection")
	collectionController := controller.NewcollectionController()
	collectionRoutes.PUT("/:id", collectionController.Update)
    collectionRoutes.POST("page/gettype", collectionController.Gettype)
	collectionRoutes.POST("page/idfault", collectionController.Getid)
	collectionRoutes.POST("page/getequ", collectionController.Getequ)




	indextempRoutes := r.Group("/api/indextemp")
	indextempController := controller.NewIndextempController()
	indextempRoutes.PUT("/:id", indextempController.Update)
	indextempRoutes.POST("page/initialize", indextempController.Datasave)
	indextempRoutes.POST("page/menu", indextempController.Menu)
	indextempRoutes.POST("page/tempequipment", indextempController.Tempequipment)
	indextempRoutes.POST("page/temporder", indextempController.Temporder)
	indextempRoutes.POST("page/parameter_change", indextempController.Parameter_change)
    indextempRoutes.POST("page/equipmentsearch", indextempController.Equipmentsearch)
	indextempRoutes.POST("page/provincecodesearch", indextempController.Provincecodesearch)



	menuRoutes := r.Group("/api/indextemp1")
	menuController := controller.NewmenuController()
    menuRoutes.POST("page/menu2", menuController.Trend)

	statisticsRoutes := r.Group("/api/statistics")
	statisticsController := controller.NewstatisticsController()
	statisticsRoutes.POST("page/statistics", statisticsController.Search)
	statisticsRoutes.POST("page/getequipment", statisticsController.Getequipment)
	statisticsRoutes.POST("page/gettype", statisticsController.Gettype)

	statisticsRoutes.POST("page/ceshi", statisticsController.Ceshi)
	statisticsRoutes.POST("page/getday", statisticsController.Getday)
	statisticsRoutes.POST("page/getidnumber", statisticsController.Getidnumber)





	datasavemineRoutes := r.Group("/api/datasavemine")
	datasavemineController := controller.NewDatasaveandmineController()
	datasavemineRoutes.PUT("/:id", datasavemineController.Update)
	datasavemineRoutes.POST("page/initialize", datasavemineController.Datasave)

	categoryRoutes := r.Group("/api/categories")
	categoryController := controller.NewCategoryController()
	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id", categoryController.Update)
	categoryRoutes.GET("/:id", categoryController.Show)
	categoryRoutes.DELETE("/:id", categoryController.Delete)
	//tablestore数据访问路由

	tableStoreDateRoutes := r.Group("/api/tableStoreDates")
	//tableStoreDateRoutes.Use(middleware.AuthMiddleware())//路由中添加中间件
	tableStoreDateController := controller.NewTableStoreController()
	tableStoreDateRoutes.POST("", tableStoreDateController.Create)
	tableStoreDateRoutes.PUT("", tableStoreDateController.Update)
	tableStoreDateRoutes.GET("", tableStoreDateController.Show)
	tableStoreDateRoutes.DELETE("", tableStoreDateController.Delete)
	tableStoreDateRoutes.POST("page/list", tableStoreDateController.PageList)

	TimemonitorRoutes := r.Group("/api/timemonitor")
	//tableStoreDateRoutes.Use(middleware.AuthMiddleware())//路由中添加中间件
	TimemonitorController := controller.NewTimeMonitorController()
	TimemonitorRoutes.POST("", TimemonitorController.Create)
	TimemonitorRoutes.PUT("", TimemonitorController.Update)
	TimemonitorRoutes.GET("", TimemonitorController.Show)
	TimemonitorRoutes.DELETE("", TimemonitorController.Delete)
	TimemonitorRoutes.POST("page/list", TimemonitorController.PageList)
	TimemonitorRoutes.POST("EquipmentIndexMonitoring", TimemonitorController.EquipmentIndexMonitoring)
	TimemonitorRoutes.POST("EquipmentHistoryMonitoring", TimemonitorController.EquipmentHistoryMonitoring)
	TimemonitorRoutes.POST("EquipmentRewriteFlag", TimemonitorController.EquipmentRewriteFlag)
	TimemonitorRoutes.POST("MultipleEquipmentInfo", TimemonitorController.MultipleEquipmentInfo)
	TimemonitorRoutes.POST("EquipmentRewriteFlagTempClear", TimemonitorController.EquipmentRewriteFlagTempClear)
	TimemonitorRoutes.POST("EquipmentRewriteFlagTempModify", TimemonitorController.EquipmentRewriteFlagTempModify)
	TimemonitorRoutes.POST("GetUpdata", TimemonitorController.GetUpdata)
    TimemonitorRoutes.POST("EquipmentInfoCreat", TimemonitorController.EquipmentInfoCreat)
	TimemonitorRoutes.POST("EquipmentInfoDelete", TimemonitorController.EquipmentInfoDelete)


	cleanTaskRoutes := r.Group("/api/cleantasks") //建立清洗任务请求路由组
	cleanTaskController := controller.NewCleanTaskController()
	cleansingleTaskController := controller.NewsingleCleanTaskController()
	cleanTaskRoutes.POST("", cleanTaskController.Create)
	cleanTaskRoutes.PUT("", cleanTaskController.Update)
	cleanTaskRoutes.GET("", cleanTaskController.Show)
	cleanTaskRoutes.DELETE("", cleanTaskController.Delete)
	cleanTaskRoutes.POST("page/list", cleanTaskController.PageList)
	cleanTaskRoutes.POST("/sigleapplicance", cleansingleTaskController.SingleCreate)

	cleanDateExistRoutes := r.Group("/api/cleandatesexist") //存在的数据请求路由组
	cleanDateExistController := controller.NewCleanDateExistController()
	cleanDateExistRoutes.POST("", cleanDateExistController.Create)
	cleanDateExistRoutes.PUT("", cleanDateExistController.Update)
	cleanDateExistRoutes.GET("", cleanDateExistController.Show)
	cleanDateExistRoutes.DELETE("", cleanDateExistController.Delete)
	cleanDateExistRoutes.POST("page/list", cleanDateExistController.PageList)

	applianceSelectRoutes := r.Group("/api/applianceSelect") //存在的数据请求路由组
	applianceSelectController := controller.NewApplianceController()
	applianceSelectRoutes.POST("", applianceSelectController.Create)
	applianceSelectRoutes.PUT("", applianceSelectController.Update)
	applianceSelectRoutes.GET("", applianceSelectController.Show)
	applianceSelectRoutes.DELETE("", applianceSelectController.Delete)
	applianceSelectRoutes.GET("page/list", applianceSelectController.PageList)

	BasicdataRoutes := r.Group("/api/Basicdata") //存在的数据请求路由组
	NewBasicdataController := controller.NewBasicdataController()
	BasicdataRoutes.POST("", NewBasicdataController.Create)
	BasicdataRoutes.PUT("", NewBasicdataController.Update)
	BasicdataRoutes.GET("", NewBasicdataController.Show)
	BasicdataRoutes.DELETE("", NewBasicdataController.Delete)
	BasicdataRoutes.POST("page/list", NewBasicdataController.PageList)
	BasicdataRoutes.POST("page/modify", NewBasicdataController.Modifyequipment)
	BasicdataRoutes.POST("page/getequipment", NewBasicdataController.Getequipment)
	BasicdataRoutes.POST("page/gettype", NewBasicdataController.Gettype)
    BasicdataRoutes.POST("modifyfaulty", controller.Modifyfaultydeviceparameters)
    BasicdataRoutes.POST("hxcqueryrewrint", controller.HxcqueryRewrint)

	SentCmdRoutes := r.Group("/api/sentcmd", middleware.IdCheck()) //热水器参数查询、改写路由组
	{
		SentCmdRoutes.POST("", controller.QueryParameter)
		SentCmdRoutes.POST("/defaultparainformation", controller.QueryParaSettingCmd)
		SentCmdRoutes.POST("/querynoparasettingcmd", controller.QueryNoParaSettingCmd)
		SentCmdRoutes.POST("/updatenobedugparafirstcmd", middleware.NoparasettingfirstCheck(), controller.RewriteNoDebugFirstCmd)
		SentCmdRoutes.POST("/rewritenodebugsecondcmd", middleware.NoparasettingsecondCheck(), controller.RewriteNoDebugSecondCmd)
		SentCmdRoutes.POST("/queryparasettingcmd", controller.QuerySuanFaDataCmd)
		SentCmdRoutes.POST("/rewriteQueryParaSettingCmd", middleware.ParainformationCheck(), controller.RewriteQueryParaSettingCmd)
		SentCmdRoutes.POST("/rewritesingleparacmd", middleware.SingleparaCheck(), controller.RewriteSingleParaCmd)
		SentCmdRoutes.POST("/rewriteFind", controller.RewriteFind)
		SentCmdRoutes.POST("/rewriteparametersettingcmd", middleware.ParaSettingCheck(), controller.RewriteParameterSettingCmd) //改写参数设置
		SentCmdRoutes.POST("/querysetparameter", controller.QuerySetParameter, controller.QueryParameter)
		SentCmdRoutes.POST("/delechangeparameter", controller.DeleChangeParameter)
		SentCmdRoutes.POST("/onekeyrescontempparacmd", controller.OneKeyResConTempParaCmd)
		SentCmdRoutes.POST("/batchonerewrite", middleware.SingleparaCheck(), middleware.BatchRewriteSingleParaCmd)
		SentCmdRoutes.POST("/downdiv", controller.Downdiv)//关机
		SentCmdRoutes.POST("/straupdiv",  controller.Straupdiv)//开机
		SentCmdRoutes.POST("/sumnoparasettingfirstcheck", controller.SumNoparasettingfirstCheck()) //检查接口
		SentCmdRoutes.POST("/batchsumpara", middleware.SingleparaCheck(), controller.BatchRewriteSumParaCmd)
		TimemonitorRoutes.POST("MultipleEquipmentInfoSingel", TimemonitorController.MultipleEquipmentInfoSingel)
		TimemonitorRoutes.POST("GetDelete", TimemonitorController.GetDelete)
	}




	CollectionRoutes := r.Group("/api/collection") //存在的数据请求路由组
	CollectionController := controller.NewCollectionController()
	CollectionRoutes.POST("", CollectionController.Create)
	CollectionRoutes.PUT("", CollectionController.Update)
	CollectionRoutes.GET("", CollectionController.Show)
	CollectionRoutes.DELETE("", CollectionController.Delete)
	CollectionRoutes.POST("page/create", CollectionController.Create1)
	CollectionRoutes.POST("page/list", CollectionController.Tempequipment)
	CollectionRoutes.POST("page/getinfo", CollectionController.Getinfo)
	CollectionRoutes.POST("page/getparachangeinfo", CollectionController.Get_Change_Para_Info)
	CollectionRoutes.POST("page/getparaRowinfo", CollectionController.Get_Row_Para_Info)

	return r
}
