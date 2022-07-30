package model

type Total struct {
	DevId        		string  `json:"dev_id" gorm:"type:varchar(20);not null"`
	DevCity             string  `json:"dev_city" gorm:"type:varchar(20);not null"`
	DevType             string  `json:"dev_type" gorm:"type:varchar(20);not null"`
	StartTime           string  `json:"start_time" gorm:"type:varchar(20);not null"`
	EndTime             string  `json:"end_time" gorm:"type:varchar(20);not null"`
	TotalNum            int     `json:"total_num" gorm:"type:varchar(20);not null"`
}

type TempScoreDownDevs struct {
	DevId        		 string `json:"dev_id" gorm:"type:varchar(20);not null"`
	City				 string `json:"dev_city" gorm:"type:varchar(20);not null"`
	CityCode            string  `json:"city_code" gorm:"type:varchar(20);not null"`
	DevType              string  `json:"dev_type" gorm:"type:varchar(20);not null"`
	StartTime            string  `json:"start_time" gorm:"type:varchar(20);not null"`
	EndTime              string  `json:"end_time" gorm:"type:varchar(20);not null"`
	TempScoreDownPercent float64 `json:"temp_score_down_percent" gorm:"type:float;not null"`
	TempScoreAvg         float64 `json:"temp_score_avg" gorm:"type:float;not null"`
	UpperLimit           float64 `json:"upper_limit" gorm:"type:float;not null"`
	LowerLimit           float64 `json:"lower_limit" gorm:"type:float;not null"`
	TempscoreDeviation   float64 `json:"tempscore_deviation" gorm:"type:float;not null"`
}
