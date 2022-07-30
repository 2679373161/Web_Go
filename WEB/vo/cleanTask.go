package vo

// 用于表达验证
type CleanTaskRequest struct {
	//CategoryId uint `json:"category_id" binding:"required"`//数据分类
	Citycode  string `json:"citycode" binding:"required"`
	Starttime string `json:"starttime" binding:"required"`
	Endtime   string `json:"endtime" binding:"required"`
}
type SingleCleanTaskRequest struct {
	//CategoryId uint `json:"category_id" binding:"required"`//数据分类
	Applianceid  string `json:"applianceid" `
	StartTime string    `json:"starttime" `
	EndTime   string    `json:"endtime" `
	AdjustTaskFlag string  `json:"adjusttaskflag"`
}
