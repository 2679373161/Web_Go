package vo

// 用于表达验证
type CreateTableStoreRequest struct{
	//CategoryId uint `json:"category_id" binding:"required"`//数据分类
	Datatime string `json:"datatime" binding:"required"`
	Flame string `json:"flame" binding:"required"`
	Outtemp string `json:"outtemp" binding:"required"`
	Settemp string `json:"settemp" binding:"required"`
	Flow int `json:"flow" binding:"required"`
	//Model string `json:"model" binding:"required"`
}
type CreateDatasaveRequest struct{
	//CategoryId uint `json:"category_id" binding:"required"`//数据分类
	Id  string `json:"id" binding:"required"`
	Label string `json:"label" binding:"required"`
	Value   string `json:"value" binding:"required"`
	//Bgcolor string `json:"bgcolor" binding:"required"`
	//Txcolor string `json:"txcolor" binding:"required"`
	//Acttxcolor string `json:"acttxcolor" binding:"required"`

	//Model string `json:"model" binding:"required"`
}
//新加入存的表
type CollectionRequest struct {
	Applianceid           string   `json:"applianceid" binding:"required"`
	Start_time    string   `json:"start_time" binding:"required"`
	End_time      string `json:"end_time" binding:"required"`
	Move_task_flag    string `json:"move_task_flag" binding:"required"`
	Mining_task_flag   string `json:"mining_task_flag" binding:"required"`
}
