package model

import (
	"time"
)

type SignleApplicance struct {
	Applianceid    string `json:"applianceid"form:"applianceid"gorm:"type:varchar(50);not null;primary_key"`
	StartTime      string `json:"starttime"form:"starttime"gorm:"type:varchar(50);not null;primary_key"`
	EndTime        string `json:"endtime"form:"endtime"gorm:"type:varchar(50);not null;primary_key"`
	AdjustTaskFlag string  `json:"adjusttaskflag"form:"adjusttaskflag"gorm:"type:varchar(50);default:0;not null"`
	MoveTaskFlag   string  `json:"movetaskflag"form:"movetaskflag"gorm:"type:varchar(50);default:0;not null"`
	MiningTaskFlag string  `json:"miningtaskflag"form:"miningtaskflag"gorm:"type:varchar(50);default:0;not null"`
	CreatedAt      time.Time `gorm:"not null"`
}

