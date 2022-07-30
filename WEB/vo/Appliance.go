package vo

// 用于表达验证
type ApplianceSelectRequest struct {
	//CategoryId uint `json:"category_id" binding:"required"`//数据分类
	DevId    string `json:"dev_id" binding:"required"`
	CityCode string `json:"city_code" binding:"required"`
	Model    string `json:"model" binding:"required"`
}
