package model

type ApplianceSelect struct {
	DevId    string `json:"dev_id" gorm:"type:varchar(50);not null;primary_key"`
	CityCode string `json:"city_code" gorm:"type:varchar(50);not null;primary_key"`
	Model    string `json:"model" gorm:"type:varchar(50);not null;primary_key"`
	Select   string `json:"select" gorm:"type:varchar(50);not null"`
}
