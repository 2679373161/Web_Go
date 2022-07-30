package model

type CleanDateExist struct {
	Citycode string `json:"citycode" gorm:"type:varchar(50);not null;unique;primary_key"`
	Dataday  string `json:"dataday" gorm:"type:varchar(50);not null;unique;primary_key"`
}
