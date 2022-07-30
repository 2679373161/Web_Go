package vo

// 用于表达验证
type CleanDateExistRequest struct {
	//CategoryId uint `json:"category_id" binding:"required"`//数据分类
	Citycode string `json:"citycode" binding:"required"`
	Dataday  string `json:"dataday" binding:"required"`
}
