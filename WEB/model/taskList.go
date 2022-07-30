package model

import "github.com/jinzhu/gorm"

type TaskList struct {
	Citycode  string `json:"citycode" gorm:"type:varchar(50);not null"`
	StartTime string `json:"starttime" gorm:"type:varchar(50);not null"`
	EndTime   string `json:"endtime" gorm:"type:varchar(50);not null"`
	gorm.Model
}
