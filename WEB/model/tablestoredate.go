package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)



type Daily_monitorings_error struct {
	Dev_Id        string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type      string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Province_code string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code     string `json:"city_code" gorm:"type:varchar(255);not null"`
	Temp_score    string `json:"temp_score" gorm:"type:varchar(255);not null"`
	Time_date     string `json:"time_date" gorm:"type:varchar(255);not null"`
	E1            string `json:"e1" gorm:"int;not null"`
	C4            string `json:"c4" gorm:"int;not null"`
	Zhendangflag         string  `json:"zhendangflag" gorm:"type:varchar(255);not null"`
    K                 string  `json:"k" gorm:"type:varchar(255);not null"`
	Worst_start       string `json:"worst_start" gorm:"type:varchar(255);not null"`
	Worst_end         string `json:"worst_end" gorm:"type:varchar(255);not null"`
	Worst_score       string `json:"worst_score" gorm:"int;not null"`
	Worst_temppattern string `json:"worst_temppattern" gorm:"type:int;not null"`
	E1fragementflag   string `json:"e1fragementflag" gorm:"type:int;not null"`
}


type PerchangeNums struct {
    DevId string `json:"dev_id" `
    Pername string `json:"pername" gorm:"type:varchar(255)"`
    FaultType string `json:"fault_type" gorm:"type:varchar(255)"`
    UpValue  string `json:"up_value" gorm:"type:varchar(255)"`
    CurreValue string `json:"curre_value" gorm:"type:varchar(255)"`
    UpdataTime string `json:"updata_time" gorm:"type:varchar(255)"`
    SuccessFlag string `json:"success_flag" gorm:"type:varchar(255)"`
	}
	
type Failed_devices struct{
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Data_time string `json:"data_time" gorm:"type:varchar(255);not null"`
	Fault_type   string `json:"fault_type" gorm:"type:varchar(255);not null"`
	Revise_flag string `json:"revise_flag" gorm:"type:varchar(255);not null"`

}


// 多参数结构体
type ParamenBatch struct {
	Appliance_id string `gorm:"primaryKey" json:"appliance_id"`
	FA           string `gorm:"column:030012;type:varchar(255);default:00" json:"FA"`
	FF           string `gorm:"column:030013;type:varchar(255);default:00" json:"FF"`
	PH           string `gorm:"column:030014;type:varchar(255);default:00" json:"PH"`
	FH           string `gorm:"column:030015;type:varchar(255);default:00" json:"FH"`
	PL           string `gorm:"column:030016;type:varchar(255);default:00" json:"PL"`
	FL           string `gorm:"column:030017;type:varchar(255);default:00" json:"FL"`
	DH           string `gorm:"column:030018;type:varchar(255);default:00" json:"dH"`
	Fd           string `gorm:"column:030019;type:varchar(255);default:00" json:"Fd"`
	CH           string `gorm:"column:030020;type:varchar(255);default:00" json:"CH"`
	FC           string `gorm:"column:030021;type:varchar(255);default:00" json:"FC"`
	NE           string `gorm:"column:030022;type:varchar(255);default:00" json:"NE"`
	CA           string `gorm:"column:030023;type:varchar(255);default:00" json:"CA"`
	FP           string `gorm:"column:030024;type:varchar(255);default:00" json:"FP"`
	LF           string `gorm:"column:030025;type:varchar(255);default:00" json:"LF"`
	HS           string `gorm:"column:030026;type:varchar(255);default:00" json:"HS"`
	Hb           string `gorm:"column:030027;type:varchar(255);default:00" json:"Hb"`
	HE           string `gorm:"column:030028;type:varchar(255);default:00" json:"HE"`
	HL           string `gorm:"column:030029;type:varchar(255);default:00" json:"HL"`
	HU           string `gorm:"column:030030;type:varchar(255);default:00" json:"HU"`
	UA           string `gorm:"column:030031;type:varchar(255);default:00" json:"UA"`
	Ub           string `gorm:"column:030032;type:varchar(255);default:00" json:"Ub"`
	Fn           string `gorm:"column:030033;type:varchar(255);default:00" json:"Fn"`
	AlterCode    string `gorm:"column:alter_code;type:varchar(255);default:00" json:"alter_code"`
	CheckAlter   string `gorm:"column:check_alter;type:varchar(255);default:00" json:"check_alter"`
}

type MultipleParaRewrite struct {
	ApplianceId     string `json:"appliance_id" gorm:"column:appliance_id"`
	HandleFlag      string ` json:"handleflag" gorm:"column:handleflag"`
	SucceedFlag     string ` json:"succeedflag" gorm:"column:succeedflag"`
	RewriteFlag     string ` json:"rewriteflag" gorm:"column:rewriteflag`
	Updatetime      string ` json:"updatetime" gorm:"column:updatetime`
	RewriteflagTemp string ` json:"rewriteflag_temp" gorm:"column:rewriteflag_temp`
}


type Multiple_Equipment_info_input struct {
	Appliance_id string `json:"equipment_ID" gorm:"type:varchar(255);not null"`
   Handleflag string `json:"handleflag" gorm:"type:varchar(255);not null"`
   Succeedflag string `json:"succeedflag" gorm:"type:varchar(255);not null"`
   Rewriteflag string `json:"rewriteflag" gorm:"type:varchar(255);not null"`
}


type Fault_overtemp struct{
	Dev_province   string `json:"dev_province" gorm:"type:varchar(255);not null"`
	Dev_city   string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Start_time string `json:"start_time" gorm:"type:varchar(255);not null"`
	Time_date   string `json:"time_date" gorm:"type:varchar(255);not null"`
	Fault_code   string `json:"fault_code" gorm:"type:varchar(255);not null"`
	Total_num string `json:"total_num" gorm:"type:varchar(255);not null"`
	Over_temp_num  string `json:"over_temp_num" gorm:"type:varchar(255);not null"`
	Abnormal_num  string `json:"abnormal_num" gorm:"type:varchar(255);not null"`
	Flow_num  string `json:"flow_num" gorm:"type:varchar(255);not null"`
	Intemp_num   string `json:"intemp_num" gorm:"type:varchar(255);not null"`
	Min_load_num   string `json:"min_load_num" gorm:"type:varchar(255);not null"`
	Avg_over_temp  string `json:"avg_over_temp" gorm:"type:varchar(255);not null"`
 }



type Heatfaultne2 struct{
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Start_time   string `json:"start_time" gorm:"type:varchar(255);not null"`
	Time_date   string `json:"time_date" gorm:"type:varchar(255);not null"`
	Fault_code   string `json:"fault_code" gorm:"type:varchar(255);not null"`
	Current_level_num   string `json:"current_level_num" gorm:"type:varchar(255);not null"`
	Total_num      string `json:"total_num" gorm:"type:varchar(255);not null"`
	Overtime_num   string `json:"overtime_num" gorm:"type:varchar(255);not null"`
	Otp            string `json:"otp" gorm:"type:varchar(255);not null"`
	Avg_overtime   string `json:"avg_overtime" gorm:"type:varchar(255);not null"`
	Avg_heat_flow   string `json:"avg_heat_flow" gorm:"type:varchar(255);not null"`
	Avg_temp_diff   string `json:"avg_temp_diff" gorm:"type:varchar(255);not null"`
	Overshoot_num    string `json:"overshoot_num" gorm:"type:varchar(255);not null"`
	Osp              string `json:"osp" gorm:"type:varchar(255);not null"`
	Avg_overshoot   string `json:"avg_overshoot" gorm:"type:varchar(255);not null"`

}


type MideaNumFault struct {
	DevType  string `json:"dev_type" gorm:"type:varchar(50);not null"`
	TimeDate string `json:"time_date" gorm:"type:varchar(50);not null"`
	Total     int `json:"total" gorm:"type:varchar(50);not null"`
	E0num       int    `json:"e0num" gorm:"type:int;not null"`
	E1num        int    `json:"e1num" gorm:"type:int;not null"`
	E2num        int    `json:"e2num" gorm:"type:int;not null"`
	E3num        int    `json:"e3num" gorm:"type:int;not null"`
	E4num       int    `json:"e4num" gorm:"type:int;not null"`
	E5num        int    `json:"e5num" gorm:"type:int;not null"`
	E6num        int    `json:"e6num" gorm:"type:int;not null"`
	E8num        int    `json:"e8num" gorm:"type:int;not null"`
	EANUM       int    `json:"eanum" gorm:"type:int;not null"`
	EENUM       int    `json:"eenum" gorm:"type:int;not null"`
	F2num        int    `json:"f2num" gorm:"type:int;not null"`
	C0num        int    `json:"c0num" gorm:"type:int;not null"`
	C1num        int    `json:"c1num" gorm:"type:int;not null"`
	C2num        int    `json:"c2num" gorm:"type:int;not null"`
	C3num        int    `json:"c3num" gorm:"type:int;not null"`
	C4num        int    `json:"c4num" gorm:"type:int;not null"`
	C5num        int    `json:"c5num" gorm:"type:int;not null"`
	C6num        int    `json:"c6num" gorm:"type:int;not null"`
	C7num        int    `json:"c7num" gorm:"type:int;not null"`
	C8num        int    `json:"c8num" gorm:"type:int;not null"`
	EHNUM        int    `json:"ehnum" gorm:"type:int;not null"`
	EFNUM        int    `json:"efnum" gorm:"type:int;not null"`
}



type MideaTypeFault struct {
	DevType         string `json:"dev_type" gorm:"type:varchar(50);not null"`
	TimeDate        string `json:"time_date" gorm:"type:varchar(50);not null"`
	E0              int    `json:"e0" gorm:"type:int;not null"`
	E1              int    `json:"e1" gorm:"type:int;not null"`
	E2              int    `json:"e2" gorm:"type:int;not null"`
	E3              int    `json:"e3" gorm:"type:int;not null"`
	E4              int    `json:"e4" gorm:"type:int;not null"`
	E5              int    `json:"e5" gorm:"type:int;not null"`
	E6              int    `json:"e6" gorm:"type:int;not null"`
	E8              int    `json:"e8" gorm:"type:int;not null"`
	EA              int    `json:"ea" gorm:"type:int;not null"`
	EE              int    `json:"ee" gorm:"type:int;not null"`
	F2              int    `json:"f2" gorm:"type:int;not null"`
	C0              int    `json:"c0" gorm:"type:int;not null"`
	C1              int    `json:"c1" gorm:"type:int;not null"`
	C2              int    `json:"c2" gorm:"type:int;not null"`
	C3              int    `json:"c3" gorm:"type:int;not null"`
	C4              int    `json:"c4" gorm:"type:int;not null"`
	C5              int    `json:"c5" gorm:"type:int;not null"`
	C6              int    `json:"c6" gorm:"type:int;not null"`
	C7              int    `json:"c7" gorm:"type:int;not null"`
	C8              int    `json:"c8" gorm:"type:int;not null"`
	EH              int    `json:"eH" gorm:"type:int;not null"`
	EF              int    `json:"eF" gorm:"type:int;not null"`
	MAX_Error       string `json:"max_error" gorm:"type:varchar(50);not null"`
	Max_Error_count int    `json:"max_error_count" gorm:"type:int;not null"`
	Max_Other_count int    `json:"max_other_count" gorm:"type:int;not null"`
	All_Error_count int    `json:"all_error_count" gorm:"type:int;not null"`
	// 修改该结构体一定要小心 该数据处理过程可能越界
}



type MideaFault struct {
	DevID       string  `json:"dev_id" gorm:"type:varchar(50);not null"`
	DevType     string  `json:"dev_type" gorm:"type:varchar(50);not null"`
	TimeDate    string  `json:"time_date" gorm:"type:varchar(50);not null"`
	ProdTime    string  `json:"prod_time" gorm:"type:varchar(50);not null"`
	ProdMonth    string  `json:"prod_month" gorm:"type:varchar(50);not null"`
	E0          int      `json:"e0" gorm:"type:varchar(255);not null"`
	E1          int      `json:"e1" gorm:"type:varchar(255);not null"`
	E2          int      `json:"e2" gorm:"type:varchar(255);not null"`
	E3          int      `json:"e3" gorm:"type:varchar(255);not null"`
	E4          int      `json:"e4" gorm:"type:varchar(255);not null"`
	E5          int      `json:"e5" gorm:"type:varchar(255);not null"`
	E6          int      `json:"e6" gorm:"type:varchar(255);not null"`
	E8          int      `json:"e8" gorm:"type:varchar(255);not null"`
	EA          int      `json:"ea" gorm:"type:varchar(255);not null"`
	EE          int      `json:"ee" gorm:"type:varchar(255);not null"`
	F2          int      `json:"f2" gorm:"type:varchar(255);not null"`
	C0          int      `json:"c0" gorm:"type:varchar(255);not null"`
	C1          int      `json:"c1" gorm:"type:varchar(255);not null"`
	C2          int      `json:"c2" gorm:"type:varchar(255);not null"`
	C3          int      `json:"c3" gorm:"type:varchar(255);not null"`
	C4          int      `json:"c4" gorm:"type:varchar(255);not null"`
	C5          int      `json:"c5" gorm:"type:varchar(255);not null"`
	C6          int      `json:"c6" gorm:"type:varchar(255);not null"`
	C7          int      `json:"c7" gorm:"type:varchar(255);not null"`
	C8          int      `json:"c8" gorm:"type:varchar(255);not null"`
	EH          int      `json:"eh" gorm:"type:varchar(255);not null"`
	EF          int      `json:"ef" gorm:"type:varchar(255);not null"`
 }


type Statistics2 struct{
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
     Avg_score   string `json:"avg_score" gorm:"type:varchar(255);not null"`
}
type Statistics1 struct{
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
	High80Day     int `json:"high80day" gorm:"type:varchar(255);not null"`
	Sixto80Day     int `json:"sixto80day" gorm:"type:varchar(255);not null"`
	Low60Day      int `json:"low60day" gorm:"type:varchar(255);not null"`
	High80Device     int `json:"high80device" gorm:"type:varchar(255);not null"`
	Sixto80Device     int `json:"sixto80device" gorm:"type:varchar(255);not null"`
	Low60Device      int `json:"low60device" gorm:"type:varchar(255);not null"`

}

type Statistics struct{
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
	High80     int `json:"high80" gorm:"type:varchar(255);not null"`
	Sixto80     int `json:"sixto80" gorm:"type:varchar(255);not null"`
	Low60      int `json:"low60" gorm:"type:varchar(255);not null"`
	Temp_score  int `json:"temp_score" gorm:"type:varchar(255);not null"`
	//Lowchaotiao int `json:"lowchaotiao" gorm:"type:varchar(255);not null"`
	//Lowjiare int `json:"lowjiare" gorm:"type:varchar(255);not null"`
	//Lowbuhengwen int `json:"lowbuhengwen" gorm:"type:varchar(255);not null"`
    //Lowjicha     int `json:"lowjicha" gorm:"type:varchar(255);not null"`
	Un_stable_mark60    int `json:"un_stable_mark60" gorm:"type:varchar(255);not null"`
	Un_heat_dev_mark60 int `json:"un_heat_dev_mark60" gorm:"type:varchar(255);not null"`
	Over_shoot_mark60    int `json:"over_shoot_mark60" gorm:"type:varchar(255);not null"`
	 Heat_mark60  int `json:"heat_mark60" gorm:"type:varchar(255);not null"`
	Un_stable_mark    int `json:"un_stable_mark" gorm:"type:varchar(255);not null"`
	Un_heat_dev_mark int `json:"un_heat_dev_mark" gorm:"type:varchar(255);not null"`
	Over_shoot_mark    int `json:"over_shoot_mark" gorm:"type:varchar(255);not null"`
	Heat_mark  int `json:"heat_mark" gorm:"type:varchar(255);not null"`
	High80device     int `json:"high80device" gorm:"type:varchar(255);not null"`
	Sixto80device     int `json:"sixto80device" gorm:"type:varchar(255);not null"`
	Low60device      int `json:"low60device" gorm:"type:varchar(255);not null"`

}



type Aliyun_info struct{
      Citycode string `json:"citycode" gorm:"type:varchar(255);not null"`
      Applianceid   string `json:"applianceid" gorm:"type:varchar(50);not null"`
	  Datatime      string `json:"datatime" gorm:"type:varchar(50);not null"`
      Comloadsegment string `json:"comloadsegment" gorm:"type:varchar(50);not null"`
	  Actualpwm    string `json:"actualpwm" gorm:"type:varchar(50);not null"`
	  Ut   string `json:"ut" gorm:"type:varchar(50);not null"`
}



type Fault_summaries struct{
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_type   string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Time_date   string `json:"time_date" gorm:"type:varchar(255);not null"`
	Fault_code   string `json:"fault_code" gorm:"type:varchar(255);not null"`
	Total_num string `json:"total_num" gorm:"type:varchar(255);not null"`
	Temp_num  string `json:"temp_num" gorm:"type:varchar(255);not null"`
	Flameout_num  string `json:"flameout_num" gorm:"type:varchar(255);not null"`
	Max_flameout_count  string `json:"max_flameout_count" gorm:"type:varchar(255);not null"`
	Set_temp   string `json:"set_temp" gorm:"type:varchar(255);not null"`
	Low_temp_min_flow   string `json:"low_temp_min_flow" gorm:"type:varchar(255);not null"`
	Temp_diff  string `json:"temp_diff" gorm:"type:varchar(255);not null"`
	Low_temp_num  string `json:"low_temp_num" gorm:"type:varchar(255);not null"`
    Liter_diff string `json:"liter_diff" gorm:"type:varchar(255);not null"`
}


type Migration_information_record struct{
	Date   string `json:"date" gorm:"type:varchar(255);not null"`
	Total_time    string `json:"total_time" gorm:"type:varchar(255);not null"`
}

type Dayinformation struct{
	//Update_time      string   `json:"update_time" gorm:"type:varchar(50);not null"`
	Abnormal_count int `json:"abnormal_count" gorm:"type:varchar(50);not null"`
	Temp_all_normal int `json:"temp_all_normal" gorm:"type:varchar(50);not null"`
	Constant_temp_abnormal   int `json:"constant_temp_abnormal" gorm:"type:varchar(50);not null"`
	Elevate_temp_abnormal    int `json:"elevate_temp_abnormal" gorm:"type:varchar(50);not null"`
	Temp_all_abnormal       int  `json:"temp_all_abnormal" gorm:"type:varchar(50);not null"`
}

type Monitor_fragement struct
{
	Dev_id   string `json:"dev_id" gorm:"type:varchar(255);not null"`
    City_code   string `json:"city_code" gorm:"type:varchar(255);not null"`
    Start_time        string `json:"start_time" gorm:"type:varchar(255);not null"`
	End_time        string `json:"end_time" gorm:"type:varchar(255);not null"`
	Temp_pattern   string `json:"temp_pattern" gorm:"type:varchar(255);not null"`
}


type Score_equipment  struct
{
	Date            string `json:"date" gorm:"type:varchar(255);not null"`
	Province   string `json:"province" gorm:"type:varchar(255);not null"`
	First_equ    string `json:"first_equ" gorm:"type:varchar(255);not null"`
	Second_equ string `json:"second_equ" gorm:"type:varchar(255);not null"`
	Third_equ string `json:"third_equ" gorm:"type:varchar(255);not null"`
	
	Fourth_equ string `json:"fourth_equ" gorm:"type:varchar(255);not null"`
	Fifth_equ string `json:"fifth_equ" gorm:"type:varchar(255);not null"`
	Sixth_equ string `json:"sixth_equ" gorm:"type:varchar(255);not null"`
}


type Score_province  struct
{
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	Temp_score    string `json:"temp_score" gorm:"type:varchar(255);not null"`
	
}

type Type_oneday_avgscores struct{
	Dev_type    string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Time_date   string `json:"time_date" gorm:"type:varchar(50);not null"`
	Avgscore    string `json:"avgscore" gorm:"type:int;not null"`
}
type Region_oneday_average_scores struct{
	Region_code   string `json:"region" gorm:"type:varchar(255);not null"`
	Dev_region    string `json:"dev_region" gorm:"type:varchar(255);not null"`
	Time_date     string `json:"time_date" gorm:"type:varchar(50);not null"`
	Avg_score      string `json:"avgscore" gorm:"type:int;not null"`
}
type City_oneday_average_scores struct{
	City_code   string `json:"region" gorm:"type:varchar(255);not null"`
	Dev_city    string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Time_date     string `json:"time_date" gorm:"type:varchar(50);not null"`
	Avg_score      string `json:"avgscore" gorm:"type:int;not null"`
}


type Multiple_Equipment struct {
	 Appliance_id  string `json:"equipment_ID" gorm:"type:varchar(255);not null"`
	 Handleflag string `json:"handleflag" gorm:"type:varchar(255);not null"`
    Succeedflag   string `json:"succeedflag" gorm:"type:varchar(255);not null"`
    Updatetime   string `json:"equipment_Time" gorm:"type:varchar(255);not null"`
   Rewriteflag string `json:"rewriteflag" gorm:"type:varchar(255);not null"`
    Rewriteflag_temp string `json:"rewriteflag_temp" gorm:"type:varchar(255);not null"`
   Index string `json:"index" gorm:"type:varchar(255);not null"`
   Value string `json:"value" gorm:"type:varchar(255);not null"`
  Isrewirted string `json:"Isrewirted" gorm:"type:varchar(255);not null"`
   State_flag string `json:"State_flag" gorm:"type:varchar(255);not null"`
   Equipment_State string `json:"equipment_State" gorm:"type:varchar(255);not null"`
    Equipment_Para_O string `json:"Equipment_Para_O" gorm:"type:varchar(255);not null"`
Equipment_Para_N string `json:"Equipment_Para_N" gorm:"type:varchar(255);not null"`
 Equipment_Para_T string `json:"Equipment_Para_T" gorm:"type:varchar(255);not null"`
	}


type Effective_statistics struct{
	Dev_id    string   `json:"dev_id" gorm:"type:varchar(50);not null"`
	Effective_day     int   `json:"effective_day" gorm:"type:varchar(50);not null"`
}


type Equipment_search struct{
	Dev_Id    string   `json:"dev_id" gorm:"type:varchar(50);not null"`
	Dev_type  string   `json:"dev_type" gorm:"type:varchar(50);not null"`
	Time_date  string   `json:"time_date" gorm:"type:varchar(50);not null"`
	Temp_score  int   `json:"temp_score" gorm:"type:varchar(50);not null"`
	Ave_unstable_proportion  float32   `json:"ave_unstable_proportion" gorm:"type:varchar(50);not null"`
	Temp_valid_time string   `json:"temp_valid_time" gorm:"type:varchar(50);not null"`
    Ave_heat_duration string   `json:"ave_heat_duration" gorm:"type:varchar(50);not null"`
    Ave_un_sable_duration  string   `json:"ave_un_sable_duration" gorm:"type:varchar(50);not null"`
    Temp_num int   `json:"temp_num" gorm:"type:varchar(50);not null"`
	City_code string   `json:"city_code" gorm:"type:varchar(50);not null"`
	//Dev_province    string `json:"dev_province" gorm:"type:varchar(255);not null"`
	//Dev_city        string `json:"dev_city" gorm:"type:varchar(255);not null"`

}

type Day_data struct{
Update_time      string   `json:"update_time" gorm:"type:varchar(50);not null"`
Day_all_dev_avg_score int `json:"day_all_dev_avg_score" gorm:"type:varchar(50);not null"`
//	Constant_temp_abnormal   int `json:"constant_temp_abnormal" gorm:"type:varchar(50);not null"`
//	Elevate_temp_abnormal    int `json:"elevate_temp_abnormal" gorm:"type:varchar(50);not null"`
  //  Temp_all_abnormal       int  `json:"temp_all_abnormal" gorm:"type:varchar(50);not null"`
 
}
type Citycode_migrate struct{
	City_migrate_code      string `json:"city_migrate_code" gorm:"type:varchar(50);not null"`
    Opt                   string `json:"opt" gorm:"type:varchar(50);not null"`
}
type Weekdate struct{
	Dev_Id    string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Time_date string `json:"time_date" gorm:"type:varchar(50);not null"`
	Temp_score string `json:"temp_score" gorm:"type:varchar(50);not null"`
} 
type Score struct{
	Heat_temp_score         int     `json:"heat_temp_score"  gorm:"type:int;not null"`
	Stable_temp_score       int     `json:"stable_temp_score" gorm:"type:int;not null"`
 }

type TableDate struct {
	//Applianceid uint64 `gorm:"primary_key"`
	//Datatime string `json:"data_time" gorm:"type:varchar(50);not null"`
	//Flame string `json:"flame" gorm:"type:varchar(50);not null"`
	//Outtemp string `json:"out_temp" gorm:"type:varchar(50);not null"`
	//Settemp string `json:"set_temp" gorm:"type:varchar(50);not null"`
	//Flow int `json:"flow" gorm:"type:varchar(50);not null"`
	//Temp_pattern int `json:"temp_pattern" gorm:"type:varchar(50);not null"`
	//Water_pattern int `json:"water_pattern" gorm:"type:varchar(50);not null"`
	//Zone_id string `json:"Zone_id" gorm:"type:varchar(50);not null"`
	Citycode      string `json:"citycode" gorm:"type:varchar(50);not null"`
	Applianceid   string `json:"applianceid" gorm:"type:varchar(50);not null"`
	Datatime      string `json:"datatime" gorm:"type:varchar(50);not null"`
	Flame         string `json:"flame" gorm:"type:varchar(50);not null"`
	Flow          int    `json:"flow" gorm:"type:varchar(50);not null"`
	Intemp        string    `json:"in_temp" gorm:"type:varchar(50);not null"`
	Outtemp       string `json:"out_temp" gorm:"type:varchar(50);not null"`
	Settemp       string `json:"set_temp" gorm:"type:varchar(50);not null"`
	Water_pattern int    `json:"water_pattern" gorm:"type:varchar(50);"`
	Temp_pattern  int    `json:"temp_pattern" gorm:"type:varchar(50);"`
	BehaviorID    int    `json:"behavior_id" gorm:"type:varchar(50);not null"`
	Zone_id       string `json:"zone_id" gorm:"type:varchar(50);not null"`
	Effect_mark  string `json:"effect_mark" gorm:"type:varchar(50);not null"`
}
type RunDate struct {
	//Applianceid uint64 `gorm:"primary_key"`
	//Datatime string `json:"data_time" gorm:"type:varchar(50);not null"`
	//Flame string `json:"flame" gorm:"type:varchar(50);not null"`
	//Outtemp string `json:"out_temp" gorm:"type:varchar(50);not null"`
	//Settemp string `json:"set_temp" gorm:"type:varchar(50);not null"`
	//Flow int `json:"flow" gorm:"type:varchar(50);not null"`
	//Temp_pattern int `json:"temp_pattern" gorm:"type:varchar(50);not null"`
	//Water_pattern int `json:"water_pattern" gorm:"type:varchar(50);not null"`
	//Zone_id string `json:"Zone_id" gorm:"type:varchar(50);not null"`
	Citycode      string `json:"citycode" gorm:"type:varchar(50);not null"`
	Applianceid   string `json:"applianceid" gorm:"type:varchar(50);not null"`
	Datatime      string `json:"datatime" gorm:"type:varchar(50);not null"`
	Flame         string `json:"flame" gorm:"type:varchar(50);not null"`
	Flow          int    `json:"flow" gorm:"type:varchar(50);not null"`
	Outtemp       string `json:"out_temp" gorm:"type:varchar(50);not null"`
	Settemp       string `json:"set_temp" gorm:"type:varchar(50);not null"`
	Water_pattern int    `json:"water_pattern" gorm:"type:varchar(50);"`
	Temp_pattern  int    `json:"temp_pattern" gorm:"type:varchar(50);"`
	BehaviorID    int    `json:"behavior_id" gorm:"type:varchar(50);not null"`
	Zone_id       string `json:"zone_id" gorm:"type:varchar(50);not null"`
	//Effect_mark  string `json:"effect_mark" gorm:"type:varchar(50);not null"`
}

type TableData struct {
	Id string `json:"id" gorm:"type:varchar(255);not null"`
	//Bgcolor string `json:"bgcolor" gorm:"type:varchar(50);not null"`
	//Txcolor string `json:"txcolor" gorm:"type:varchar(50);not null"`
	//Acttxcolor string `json:"acttxcolor" gorm:"type:varchar(50);not null"`
	Label string `json:"label"  gorm:"type:varchar(50);not null"`
	Value string `json:"value" gorm:"type:varchar(50);not null"`
}
type Dataparameter struct {
	Bgcolor    string `json:"bgcolor" gorm:"type:varchar(50);not null"`
	Txcolor    string `json:"txcolor" gorm:"type:varchar(50);not null"`
	Acttxcolor string `json:"acttxcolor" gorm:"type:varchar(50);not null"`
	Menucolor  string `json:"menucolor" gorm:"type:varchar(50);not null"`
	Viewcolor  string `json:"viewcolor" gorm:"type:varchar(50);not null"`
	Timestamp  string `json:"timestamp" gorm:"type:varchar(50);not null"`
}
type Tableoption struct {
	Dev_Id          string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type        string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Devid           string `json:"devid" gorm:"type:varchar(255);not null"`
	Loc_Code        int    `json:"loc_code" gorm:"type:int;not null"`
	Total_Time      string `json:"total_time" gorm:"type:time;not null"`
	Mode_Proportion string `json:"mode_proportion" gorm:"type:varchar(255);not null"`
	Min_Time        string `json:"min_time" gorm:"type:time;not null"`
	Max_Time        string `json:"max_time" gorm:"type:time;not null"`
	Total_Num       int    `json:"total_num" gorm:"type:int;not null"`
	Applianceid     uint64 `gorm:"primary_key"`
}
type Table_all_select struct {
	Dev_Id          string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Model_type      string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Province_code   string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code       string `json:"city_code" gorm:"type:varchar(255);not null"`
	Dev_province    string `json:"dev_province" gorm:"type:varchar(255);not null"`
	Dev_city        string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Opt             string `json:"handle_flag" gorm:"type:int;not null"`
	Monitoring_flag string `json:"monitoring_flag" gorm:"type:int;not null"`
}
type Daily_monitorings struct {
	Dev_Id        string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type      string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Province_code string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code     string `json:"city_code" gorm:"type:varchar(255);not null"`

	Temp_score    string `json:"temp_score" gorm:"type:varchar(255);not null"`
	Time_date     string `json:"time_date" gorm:"type:varchar(255);not null"`
	Abnormal_flag string `json:"abnormal_flag" gorm:"type:varchar(255);not null"`
}
type Tableplace struct {
	Dev_Id                   string  `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type                 string  `json:"dev_type" gorm:"type:varchar(255);not null"`
	Water_valid_time         string  `json:"water_valid_time" gorm:"type:varchar(255);not null"`
	City_code                string  `json:"city_code" gorm:"type:int;not null"`
	Province_code            string  `json:"province_code" gorm:"type:int;not null"`
	Region_code              string  `json:"region_code" gorm:"type:varchar(255);not null"`
	Time_date                string  `json:"time_date" gorm:"type:time;not null"`
	Water_num                int     `json:"water_num" gorm:"type:time;not null"`
	Stable_behavior          string  `json:"stable_behavior" gorm:"type:time;not null"`
	Un_stable_behavior       int     `json:"un_stable_behavior" gorm:"type:time;not null"`
	Average_time             string  `json:"average_time" gorm:"type:varchar(255);not null"`
	Stable_proportion        string  `json:"stable_proportion" gorm:"type:varchar(255);not null"`
	Un_stable_proportion     float32 `json:"un_stable_proportion" gorm:"type:varchar(255);not null"`
	Stable_time              string  `json:"stable_time" gorm:"type:varchar(255);not null"`
	Un_stable_time           string  `json:"un_stable_time" gorm:"type:varchar(255);not null"`
	Maximum_time             string  `json:"maximum_time" gorm:"type:varchar(255);not null"`
	Minimum_time             string  `json:"minimum_time" gorm:"type:varchar(255);not null"`
	Equipment_num            int     `json:"equipment_num" gorm:"type:varchar(255);not null"`
	Water_score              int     `json:"water_score" gorm:"type:varchar(255);not null"`
	Temp_valid_time          string  `json:"temp_valid_time" gorm:"type:varchar(255);not null"`
	Temp_num                 int     `json:"temp_num" gorm:"type:time;not null"`
	Ave_heat_duration        string  `json:"ave_heat_duration" gorm:"type:varchar(255);not null"`
	Ave_un_sable_duration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(255);not null"`
	Ave_unstable_proportion  float32 `json:"ave_unstable_proportion" gorm:"type:varchar(255);not null"`
	Un_stable_temp_dev       float32 `json:"un_stable_temp_dev" gorm:"type:varchar(255);not null"`
	Temp_score               int     `json:"temp_score" gorm:"type:varchar(255);not null"`
	New_teme_score            int     `json:"New_teme_score" gorm:"type:varchar(255);not null"`
	Normal_all_pro           float32 `json:"normal_all_pro" gorm:"type:varchar(255);not null"`
	Abnormal_heat_pro        float32 `json:"abnormal_heat_pro" gorm:"type:varchar(255);not null"`
	Abnormal_stable_temp_pro float32 `json:"abnormal_stable_temp_pro" gorm:"type:varchar(255);not null"`
	Abnormal_all_pro         float32 `json:"abnormal_all_pro" gorm:"type:varchar(255);not null"`
	Temp_score_by_day         int     `json:"temp_score_by_day" gorm:"type:varchar(255);not null"`
    Overshoot_value            int     `json:"Overshoot_value" gorm:"type:varchar(255);not null"`
    Temp_range          int     `json:"temp_range" gorm:"type:varchar(255);not null"`
    Heat_temp_score     int     `json:"heat_temp_score" gorm:"type:varchar(255);not null"`
	Stable_temp_score     int     `json:"stable_temp_score" gorm:"type:varchar(255);not null"`
}

type MonitorTable struct {
	Dev_Id                  string  `json:"dev_id" gorm:"type:varchar(255);not null"`
	Start_time              string  `json:"start_time" gorm:"type:varchar(255);not null"`
	End_time                string  `json:"end_time" gorm:"type:varchar(255);not null"`
	Duration_time           string  `json:"duration_time" gorm:"type:varchar(255);not null"`
	Water_pattern           int     `json:"water_pattern" gorm:"type:int;not null"`
	Flow_avg                float64 `json:"flow_avg" gorm:"type:float;not null"`
	Small_water             float64 `json:"small_water" gorm:"type:float;not null"`
	Deviation               float64 `json:"deviation" gorm:"type:float;not null"`
	Up_number               int     `json:"up_number" gorm:"type:int;not null"`
	Down_number             int     `json:"down_number" gorm:"type:int;not null"`
	Water_score             int     `json:"water_score" gorm:"type:int;not null"`
	Heat_duration           string  `json:"heat_duration" gorm:"type:varchar(255);not null"`
	Un_stable_temp_duration string  `json:"un_stable_temp_duration" gorm:"type:varchar(255);not null"`
	Un_stable_temp_percent  float32 `json:"un_stable_temp_percent" gorm:"type:float;not null"`
	Un_heat_dev             float32 `json:"un_heat_dev" gorm:"type:float;not null"`
	Temp_pattern            int     `json:"temp_pattern"  gorm:"type:int;not null"`
	Overshoot_value         int     `json:"overshoot_value"  gorm:"type:int;not null"`
	State_accuracy          int     `json:"state_accuracy"  gorm:"type:int;not null"`
	Temp_score              int     `json:"temp_score"  gorm:"type:int;not null"`
	New_temp_score          int     `json:"new_temp_score"  gorm:"type:int;not null"`
	Heat_temp_score         int     `json:"heat_temp_score"  gorm:"type:int;not null"`
	Stable_temp_score       int     `json:"stable_temp_score" gorm:"type:int;not null"`
	Temp_score_judge_flag   int     `json:"temp_judge_flag" gorm:"type:tinyint;not null"`
	Water_flag              int     `json:"water_flag" gorm:"type:int;not null"`
	Temp_flag               int     `json:"temp_flag" gorm:"type:int;not null"`
	Abnormal_state          int     `json:"abnormal_state" gorm:"type:int;not null"`
}

type Datafeature struct {
	Start_time     string `json:"start_time" gorm:"type:varchar(255);not null"`
	Update_time    string `json:"update_time" gorm:"type:varchar(255);not null"`
	City_count     string `json:"city_count" gorm:"type:varchar(255);not null"`
	Type_count     string `json:"type_count" gorm:"type:varchar(255);not null"`
	Equ_count      string `json:"equ_count" gorm:"type:varchar(255);not null"`
	Day_data_sum   string `json:"day_data_sum" gorm:"type:varchar(255);not null"`
	Avg_data_sum   string `json:"avg_data_sum" gorm:"type:varchar(255);not null"`
	Total_data_sum string `json:"total_data_sum" gorm:"type:varchar(255);not null"`
	Data_days      string `json:"data_days" gorm:"type:varchar(255);not null"`
	Free_space     string `json:"free_space" gorm:"type:varchar(255);not null"`
	Under_eight_five string `json:"under_eight_five" gorm:"type:varchar(255);not null"`
    Top_eight_five    string `json:"top_eight_five" gorm:"type:varchar(255);not null"`
	Under_sixty   string   `json:"under_sixty" gorm:"type:varchar(255);not null"`
	One_hundred   string   `json:"one_hundred" gorm:"type:varchar(255);not null"`
   Avg_water_num  string    `json:"avg_water_num" gorm:"type:varchar(255);not null"`
   Avg_water_time  string    `json:"avg_water_time" gorm:"type:varchar(255);not null"`
   Avg_heat_time  string    `json:"avg_heat_time" gorm:"type:varchar(255);not null"`
   Dev_avg_stable_score string    `json:"dev_avg_stable_score" gorm:"type:varchar(255);not null"`
   Dev_avg_heat_score string    `json:"dev_avg_heat_score" gorm:"type:varchar(255);not null"`
   Day_all_dev_avg_score string  `json:"day_all_dev_avg_score" gorm:"type:varchar(255);not null"`
   Total_processing_time  string  `json:"total_processing_time" gorm:"type:varchar(255);not null"`
}
type Regiontempwater struct {
	Region     string
	Tempscore  []int
	Waterscore []int
	Month      []string
}
type Tableselect struct {
	Dev_Id        string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type      string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Province_code string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code     string `json:"city_code" gorm:"type:varchar(255);not null"`
	Region_code   string `json:"region_code" gorm:"type:varchar(255);not null"`
	Dev_province  string `json:"dev_province" gorm:"type:varchar(255);not null"`
	Dev_city      string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Dev_region    string `json:"dev_region" gorm:"type:varchar(255);not null"`
	Disable       string `json:"disable" gorm:"type:int;not null"`
	Distype       string `json:"distype" gorm:"type:int;not null"`
	Disequipment  string `json:"disequipment" gorm:"type:int;not null"`
	Discityselect string `json:"discityselect" gorm:"type:int;not null"`
	Handle_flag   string `json:"handle_flag" gorm:"type:int;not null"`
	Monitor_flag  string `json:"monitoring_flag" gorm:"type:int;not null"`
}
type Typeselection struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type Modeltype_and_nums struct {
	Model_type string `json:"model_type" gorm:"type:varchar(255);not null"`
	Num        string `json:"num" gorm:"type:varchar(255);not null"`
	Opt        string `json:"opt" gorm:"type:varchar(255);not null"`
}
type Fragmentflow struct {
	Equipment_num int `json:"equipment_num" gorm:"type:varchar(255);not null"`
	Three         int `json:"three" gorm:"type:varchar(255);not null"`
	Five          int `json:"five" gorm:"type:varchar(255);not null"`
	Ten           int `json:"ten" gorm:"type:varchar(255);not null"`
}
type Behavior struct {
	City_code     string  `json:"city_code" gorm:"type:int;not null"`
	Province_code string  `json:"province_code" gorm:"type:int;not null"`
	Dev_Id        string  `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type      string  `json:"dev_type" gorm:"type:varchar(255);not null"`
	Start_time    string  `json:"start_time" gorm:"type:varchar(255);not null"`
	Sec0p         float32 `json:"sec0p" gorm:"type:varchar(255);not null"`
	Sec30p        float32 `json:"sec30p" gorm:"type:varchar(255);not null"`
	Min3p         float32 `json:"min3p" gorm:"type:varchar(255);not null"`
	Min10p        float32 `json:"min10p" gorm:"type:varchar(255);not null"`
}
type Tablebehavior struct {
	Dev_id        string
	Dev_type      string
	City          string
	Duration_time []float32
}
type Tabletimerate struct {
	Province      string
	Equipment_num int
	Duration_time map[string]float64
}
type Tablefragment struct {
	Dev_Id                  string  `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type                string  `json:"dev_type" gorm:"type:varchar(255);not null"`
	Water_pattern           string  `json:"water_pattern" gorm:"type:varchar(255);not null"`
	Temp_pattern            string  `json:"temp_pattern" gorm:"type:varchar(255);not null"`
	Start_time              string  `json:"start_time" gorm:"type:varchar(255);not null"`
	End_time                string  `json:"end_time" gorm:"type:varchar(255);not null"`
	Duration_time           string  `json:"duration_time" gorm:"type:varchar(255);not null"`
	Extreme                 string  `json:"extreme" gorm:"type:varchar(255);not null"`
	Max_change              string  `json:"max_change" gorm:"type:varchar(255);not null"`
	Average                 string  `json:"average" gorm:"type:varchar(255);not null"`
	Deviation               float32 `json:"deviation" gorm:"type:varchar(255);not null"`
	Up_number               string  `json:"up_number" gorm:"type:varchar(255);not null"`
	Down_number             string  `json:"down_number" gorm:"type:varchar(255);not null"`
	Water_score             int     `json:"water_score" gorm:"type:varchar(255);not null"`
	Temp_score              int     `json:"temp_score" gorm:"type:varchar(255);not null"`
	New_Temp_score              int     `json:"new_temp_score" gorm:"type:varchar(255);not null"`
	Heat_temp_score            int     `json:"heat_temp_score" gorm:"type:varchar(255);not null"`
	Stable_temp_score       int     `json:"stable_temp_score" gorm:"type:varchar(255);not null"`
	Heat_duration           string  `json:"heat_duration" gorm:"type:varchar(255);not null"`
	Un_stable_temp_percent  float32 `json:"un_stable_temp_percent" gorm:"type:varchar(255);not null"`
	Un_stable_temp_duration string  `json:"un_stable_temp_duration" gorm:"type:varchar(255);not null"`
	Un_heat_dev             string  `json:"un_heat_dev" gorm:"type:varchar(255);not null"`
	Overshoot_value         int     `json:"overshoot_value" gorm:"type:varchar(255);not null"`
	State_accuracy          int     `json:"state_accuracy" gorm:"type:varchar(255);not null"`
	Effect_flag              int     `json:"effect_flag" gorm:"type:varchar(255);not null"`
	Fault_code             string   `json:"fault_code" gorm:"type:varchar(255);not null"`
	Score_by_weight           string   `json:"score_by_weight" gorm:"type:varchar(255);not null"`
	Heat_weight              string `json:"heat_weight" gorm:"type:varchar(255);not null"`
	Heat_percent             string `json:"heat_percent" gorm:"type:varchar(255);not null"`
	E1_num                     int     `json:"e1num" gorm:"type:varchar(255);not null"`
}
type Tablewaterbehavior struct {
	City_code     string `json:"city_code" gorm:"type:int;not null"`
	Province_code string `json:"province_code" gorm:"type:int;not null"`
	Dev_Id        string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Start_time    string `json:"start_time" gorm:"type:varchar(255);not null"`
	End_time      string `json:"end_time" gorm:"type:varchar(255);not null"`
	Duration_time string `json:"duration_time" gorm:"type:varchar(255);not null"`
	Fragment_num  string `json:"fragment_num" gorm:"type:varchar(255);not null"`
	Effect_flag   string `json:"effect_flag" gorm:"type:varchar(255);not null"`
}
type TableDate1 struct {
	Dev_Id          string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Dev_type        string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Devid           string `json:"devid" gorm:"type:varchar(255);not null"`
	Loc_Code        int    `json:"loc_code" gorm:"type:int;not null"`
	Total_Time      string `json:"total_time" gorm:"type:time;not null"`
	Mode_Proportion string `json:"mode_proportion" gorm:"type:varchar(255);not null"`
	Min_Time        string `json:"min_time" gorm:"type:time;not null"`
	Max_Time        string `json:"max_time" gorm:"type:time;not null"`
	Total_Num       int    `json:"total_num" gorm:"type:int;not null"`
	Applianceid     uint64 `gorm:"primary_key"`
}
type TableDate2 struct {
	Dev_type               string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Stable_time            string `json:"stable_time" gorm:"type:varchar(255);not null"`
	Fluctuation_time       string `json:"fluctuation_time" gorm:"type:varchar(255);not null"`
	Down_time              string `json:"down_time" gorm:"type:varchar(255);not null"`
	Up_time                string `json:"up_time" gorm:"type:varchar(255);not null"`
	Osc_time               string `json:"osc_time" gorm:"type:varchar(255);not null"`
	Stable_proportion      string `json:"stable_proportion" gorm:"type:varchar(255);not null"`
	Fluctuation_proportion string `json:"fluctuation_proportion" gorm:"type:varchar(255);not null"`
	Down_proportion        string `json:"down_proportion" gorm:"type:varchar(255);not null"`
	Up_proportion          string `json:"up_proportion" gorm:"type:varchar(255);not null"`
	Osc_proportion         string `json:"osc_proportion" gorm:"type:varchar(255);not null"`
}
type TableDate3 struct {
	Dev_Id                 string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Time_date              string `json:"time_date" gorm:"type:varchar(255);not null"`
	Valid_time             string `json:"valid_time" gorm:"type:varchar(255);not null"`
	Average_time           string `json:"average_time" gorm:"type:varchar(255);not null"`
	Province_code          string `json:"province_code" gorm:"type:varchar(255);not null"`
	City_code              string `json:"city_code" gorm:"type:varchar(255);not null"`
	Region_code            string `json:"region_code" gorm:"type:varchar(255);not null"`
	Dev_province           string `json:"dev_province" gorm:"type:varchar(255);not null"`
	Dev_city               string `json:"dev_city" gorm:"type:varchar(255);not null"`
	Dev_region             string `json:"dev_region" gorm:"type:varchar(255);not null"`
	Loc_code               string `json:"id" gorm:"type:varchar(255);not null"`
	Dev_type               string `json:"dev_type" gorm:"type:varchar(255);not null"`
	Dev_num                string `json:"dev_num" gorm:"type:varchar(255);not null"`
	Validdata_vloume       string `json:"validdata_vloume" gorm:"type:varchar(255);not null"`
	Dailydata_vloume       string `json:"dailydata_vloume" gorm:"type:varchar(255);not null"`
	Stable_time            string `json:"stable_time" gorm:"type:varchar(255);not null"`
	Fluctuation_time       string `json:"fluctuation_time" gorm:"type:varchar(255);not null"`
	Down_time              string `json:"down_time" gorm:"type:varchar(255);not null"`
	Up_time                string `json:"up_time" gorm:"type:varchar(255);not null"`
	Osc_time               string `json:"osc_time" gorm:"type:varchar(255);not null"`
	Stable_proportion      string `json:"stable_proportion" gorm:"type:varchar(255);not null"`
	Fluctuation_proportion string `json:"fluctuation_proportion" gorm:"type:varchar(255);not null"`
	Down_proportion        string `json:"down_proportion" gorm:"type:varchar(255);not null"`
	Up_proportion          string `json:"up_proportion" gorm:"type:varchar(255);not null"`
	Osc_proportion         string `json:"osc_proportion" gorm:"type:varchar(255);not null"`
}

//新加入的本地frag测试表
type Tablefragceshi struct {
	Dev_Id                 string `json:"dev_id" gorm:"type:varchar(255);not null"`
	Start_time             string `json:"start_time" gorm:"type:varchar(255);not null"`
	End_time               string `json:"end_time" gorm:"type:varchar(255);not null"`
	Deviation              string `json:"deviation" gorm:"type:float;not null"`
	Water_score            string `json:"water_score" gorm:"type:varchar(255);not null"`
	Water_pattern          string `json:"water_pattern" gorm:"type:varchar(255);not null"`
	Heat_duration          string `json:"heat_duration" gorm:"type:int;not null"`
	Un_stable_temp_percent string `json:"un_stable_temp_percent" gorm:"type:int;not null"`
	Overshoot_value        string `json:"overshoot_value" gorm:"type:int;not null"`
	State_accuracy         string `json:"state_accuracy" gorm:"type:int;not null"`
	Temp_score             string `json:"temp_score" gorm:"type:int;not null"`
	NewTempScore           string `json:"new_temp_score" gorm:"type:int;not null"`
	Temp_pattern           string `json:"temp_pattern" gorm:"type:int;not null"`
}

//新加入的采集信息存入表
type Collection struct {
	Applianceid      string `json:"applianceid" gorm:"type:varchar(50);not null"`
	StartTime        string `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx"`
	EndTime          string `json:"end_time" gorm:"type:varchar(20);not null;"`
	Move_task_flag   string `json:"move_task_flag" gorm:"type:tinyint;not null;default:0"`
	Mining_task_flag string `json:"mining_task_flag" gorm:"type:tinyint;not null;default:0"`
}

//设备最终参数
type ParameterOutput struct {
	Appliance_id       int    `json:"appliance_id" gorm:"type:int;not null"`
	Code               string `json:"code" gorm:"type:string;not null"`
	CurrentValue       int    `json:"current_value" gorm:"type:int;not null"`
	RewriteSuccessFlag int    `json:"rewrite_success_flag" gorm:"type:int;not null"`
	Updatetime         string `json:"updatetime" gorm:"type:varchar(50);not null"`
}

//设备默认参数
type ParameterDefaultsOutput struct {
	Appliance_id  int    `json:"appliance_id" gorm:"type:int;not null"`
	Code          string `json:"code" gorm:"type:string;not null"`
	Default_value int    `json:"default_value" gorm:"type:int;not null"`
	Updatetime    string `json:"updatetime" gorm:"type:varchar(50);not null"`
}

//数据表结构体
type TableData_para_temp_struct struct {
	Para_name     string `json:"para_name" gorm:"type:string;not null"`
	Code          string `json:"code" gorm:"type:string;not null"`
	Default_value string `json:"default_value" gorm:"type:string;not null"`
	Recent_date   string `json:"recent_date" gorm:"type:string;not null"`
	Advice_date   string `json:"advice_date" gorm:"type:string;not null"`
	Modify_date   string `json:"modify_date" gorm:"type:string;not null"`
	Max_limit     string `json:"Max_limit" gorm:"type:string;not null"`
	Min_limit     string `json:"Min_limit" gorm:"type:string;not null"`
	IsEdit        bool   `json:"isEdit" gorm:"type:bool;not null"`
}
type TableData_para_single_struct struct {
	Para_name     string `json:"para_name" gorm:"type:string;not null"`
	Code          string `json:"code" gorm:"type:string;not null"`
	Default_value string `json:"default_value" gorm:"type:string;not null"`
	Recent_date   string `json:"recent_date" gorm:"type:string;not null"`
	Advice_date   string `json:"advice_date" gorm:"type:string;not null"`
	Modify_date   string `json:"modify_date" gorm:"type:string;not null"`
	Max_limit     string `json:"Max_limit" gorm:"type:string;not null"`
	Min_limit     string `json:"Min_limit" gorm:"type:string;not null"`
	Serial_number string `json:"serial_number" gorm:"type:string;not null"`
	IsEdit        bool   `json:"isEdit" gorm:"type:bool;not null"`
}
type Para_history_struct struct {
	Appliance_id        int    `json:"appliance_id" gorm:"type:int;not null"`
	Updatetime          string `json:"updatetime" gorm:"type:varchar;not null;"`
	LastValue           int    `json:"last_value" gorm:"type:int;not null"`
	Value               int    `json:"value" gorm:"type:int;not null"`
	Code                string `json:"code" gorm:"type:varchar(255)"`
	LatestParameterFlag string `json:"latest_parameter_flag" gorm:"type:varchar(255)"`
}
type Para_history_out_struct struct {
	Appliance_id        int    `json:"appliance_id" gorm:"type:int;not null"`
	Para_name           string `json:"para_name" gorm:"type:string;not null"`
	LastValue           string `json:"last_value" gorm:"type:varchar;not null;"`
	Updatetime          string `json:"updatetime" gorm:"type:varchar;not null;"`
	Value               string `json:"value" gorm:"type:string;not null"`
	Code                string `json:"code" gorm:"type:varchar(255)"`
	LatestParameterFlag string `json:"latest_parameter_flag" gorm:"type:varchar(255)"`
}

//编码结构体
type Code_para_struct struct {
	Parameter string `json:"parameter" gorm:"type:string;not null"`
	Code      string `json:"code" gorm:"type:string;not null"`
	Max       string `json:"max" gorm:"type:string;not null"`
	Min       string `json:"min" gorm:"type:string;not null"`
}

//编码结构体
type Code_para_Serial_struct struct {
	Parameter     string `json:"parameter" gorm:"type:string;not null"`
	Serial_number string `json:"serial_number" gorm:"type:string;not null"`
}

//参数输出结构体
type NonParameterOutput struct {
	PH                                                       int `json:"PH" gorm:"type:int;not null"`
	PL                                                       int `json:"PL" gorm:"type:int;not null"`
	FH                                                       int `json:"FH" gorm:"type:int;not null"`
	DH                                                       int `json:"dH" gorm:"type:int;not null"`
	Fd                                                       int `json:"Fd" gorm:"type:int;not null"`
	CH                                                       int `json:"CH" gorm:"type:int;not null"`
	FC                                                       int `json:"FC" gorm:"type:int;not null"`
	FL                                                       int `json:"FL" gorm:"type:int;not null"`
	MaximumLoadFanCurrentDeviationCoefficient                int `json:"maximum_load_fan_current_deviation_coefficient" gorm:"type:int;not null"`
	MinimumLoadFanCurrentDeviationCoefficient                int `json:"minimum_load_fan_current_deviation_coefficient" gorm:"type:int;not null"`
	MaximumLoadFanDutyCycleDeviationCoefficient              int `json:"maximum_load_fan_duty_cycle_deviation_coefficient" gorm:"type:int;not null"`
	MinimumLoadFanDutyCycleDeviationCoefficient              int `json:"minimum_load_fan_duty_cycle_deviation_coefficient" gorm:"type:int;not null"`
	BackwaterFlowValue                                       int `json:"backwater_flow_value" gorm:"type:int;not null"`
	FrequencyCompensationValueOfWindPressureSensorAlarmPoint int `json:"frequency_compensation_value_of_wind_pressure_sensor_alarm_point" gorm:"type:int;not null"`
	KA0                                                      int `json:"ka_0" gorm:"type:int;not null"`
	KA1                                                      int `json:"ka_1" gorm:"type:int;not null"`
	KA2                                                      int `json:"ka_2" gorm:"type:int;not null"`
	KA3                                                      int `json:"ka_3" gorm:"type:int;not null"`
	KB0                                                      int `json:"kb_0" gorm:"type:int;not null"`
	KB1                                                      int `json:"kb_1" gorm:"type:int;not null"`
	KB2                                                      int `json:"kb_2" gorm:"type:int;not null"`
	KB3                                                      int `json:"kb_3" gorm:"type:int;not null"`
	KC0                                                      int `json:"kc_0" gorm:"type:int;not null"`
	KC1                                                      int `json:"kc_1" gorm:"type:int;not null"`
	KC2                                                      int `json:"kc_2" gorm:"type:int;not null"`
	KC3                                                      int `json:"kc_3" gorm:"type:int;not null"`
	KF0                                                      int `json:"kf_0" gorm:"type:int;not null"`
	KF1                                                      int `json:"kf_1" gorm:"type:int;not null"`
	KF2                                                      int `json:"kf_2" gorm:"type:int;not null"`
	KF3                                                      int `json:"kf_3" gorm:"type:int;not null"`
	T1A0                                                     int `json:"t_1_a_0" gorm:"type:int;not null"`
	T1A1                                                     int `json:"t_1_a_1" gorm:"type:int;not null"`
	T1A2                                                     int `json:"tia_2" gorm:"type:int;not null"`
	T1A3                                                     int `json:"t_1_a_3" gorm:"type:int;not null"`
	T1C0                                                     int `json:"t_1_c_0" gorm:"type:int;not null"`
	T1C1                                                     int `json:"t_1_c_1" gorm:"type:int;not null"`
	T1C2                                                     int `json:"t_1_c_2" gorm:"type:int;not null"`
	T1C3                                                     int `json:"t_1_c_3" gorm:"type:int;not null"`
	T2A0                                                     int `json:"t_2_a_0" gorm:"type:int;not null"`
	T2A1                                                     int `json:"t_2_a_1" gorm:"type:int;not null"`
	T2A2                                                     int `json:"t_2_a_2" gorm:"type:int;not null"`
	T2A3                                                     int `json:"t_2_a_3" gorm:"type:int;not null"`
	T2C0                                                     int `json:"t_2_c_0" gorm:"type:int;not null"`
	T2C1                                                     int `json:"t_2_c_1" gorm:"type:int;not null"`
	T2C2                                                     int `json:"t_2_c_2" gorm:"type:int;not null"`
	T2C3                                                     int `json:"t_2_c_3" gorm:"type:int;not null"`
	TDA0                                                     int `json:"tda_0" gorm:"type:int;not null"`
	TDA1                                                     int `json:"tda_1" gorm:"type:int;not null"`
	TDA2                                                     int `json:"tda_2" gorm:"type:int;not null"`
	TDA3                                                     int `json:"tda_3" gorm:"type:int;not null"`
	TDC0                                                     int `json:"tdc_0" gorm:"type:int;not null"`
	TDC1                                                     int `json:"tdc_1" gorm:"type:int;not null"`
	TDC2                                                     int `json:"tdc_2" gorm:"type:int;not null"`
	TDC3                                                     int `json:"tdc_3" gorm:"type:int;not null"`
	WC0                                                      int `json:"wc_0" gorm:"type:int;not null"`
	WC1                                                      int `json:"wc_1" gorm:"type:int;not null"`
	WC2                                                      int `json:"wc_2" gorm:"type:int;not null"`
	WC3                                                      int `json:"wc_3" gorm:"type:int;not null"`
	WO0                                                      int `json:"wo_0" gorm:"type:int;not null"`
	WO1                                                      int `json:"wo_1" gorm:"type:int;not null"`
	WO2                                                      int `json:"wo_2" gorm:"type:int;not null"`
	WO3                                                      int `json:"wo_3" gorm:"type:int;not null"`
}
type Real_mideaFault struct{
	DevID    string `json:"dev_id" gorm:"type:varchar(50);not null"`
	Dev_type    string `json:"dev_type" gorm:"type:varchar(50);not null"`
	Start_time string `json:"start_time" gorm:"type:varchar(50);not null"`
	End_time   string `json:"end_time" gorm:"type:varchar(50);not null"`
	E0              int    `json:"e0" gorm:"type:int;not null"`
	E1              int    `json:"e1" gorm:"type:int;not null"`
	E2              int    `json:"e2" gorm:"type:int;not null"`
	E3              int    `json:"e3" gorm:"type:int;not null"`
	E4              int    `json:"e4" gorm:"type:int;not null"`
	E5              int    `json:"e5" gorm:"type:int;not null"`
	E6              int    `json:"e6" gorm:"type:int;not null"`
	E8              int    `json:"e8" gorm:"type:int;not null"`
	EA              int    `json:"ea" gorm:"type:int;not null"`
	EE              int    `json:"ee" gorm:"type:int;not null"`
	F2              int    `json:"f2" gorm:"type:int;not null"`
	C0              int    `json:"c0" gorm:"type:int;not null"`
	C1              int    `json:"c1" gorm:"type:int;not null"`
	C2              int    `json:"c2" gorm:"type:int;not null"`
	C3              int    `json:"c3" gorm:"type:int;not null"`
	C4              int    `json:"c4" gorm:"type:int;not null"`
	C5              int    `json:"c5" gorm:"type:int;not null"`
	C6              int    `json:"c6" gorm:"type:int;not null"`
	C7              int    `json:"c7" gorm:"type:int;not null"`
	C8              int    `json:"c8" gorm:"type:int;not null"`
	EH              int    `json:"eH" gorm:"type:int;not null"`
	EF              int    `json:"eF" gorm:"type:int;not null"`
}



func (tableStoreDate *TableDate) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
