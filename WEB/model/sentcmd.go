package model

type DailyMonitoring struct {
	ProvinceCode                                             string `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode                                                 string `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType                                                  string `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId                                                    string `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx; "`
	TimeDate                                                 string `json:"time_date" gorm:"type:varchar(20);not null;index:time_date_idx; "`
	TempScore                                                int    `json:"temp_score" gorm:"type:int;not null"`
	HeatTempScore                                            int    `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore                                          int    `json:"stable_temp_score" gorm:"type:int;not null"`
	UnStableMark                                             int     `json:"un_stable_mark" gorm:"type:int;not null"`
	UnHeatDevMark                                            int     `json:"un_heat_dev_mark" gorm:"type:int;not null"`
	OverShootMark                                            int     `json:"over_shoot_mark" gorm:"type:int;not null"`
	HeatMark                                                 int     `json:"heat_mark gorm":"type:int;not null"`
	E1                                                      int       `json:"e1" gorm:"type:int;not null"`
	C4                                                      int       `json:"c4" gorm:"type:int;not null"`
	Worst_start                                             string `json:"worst_start" gorm:"type:varchar(20);not null;`
	Worst_end                                               string `json:"worst_end" gorm:"type:varchar(20);not null;`
	Worst_temppattern								        int `json:"worst_temppattern" gorm:"type:int;not null;`
	Worst_score                                               int `json:"worst_score" gorm:"type:int;not null;`
	Ut                                                      int       `json:"ut" gorm:"type:int;not null"`
	F                                                      float64       `json:"f" gorm:"type:float;"not null"`
	K                                                      float64     `json:"k" gorm:"type:float;"not null"`
	Zhendangflag                                         int       `json:"zhendangflag" gorm:"type:int;not null"`
}
type Bo struct {
	Identifier string `json:"identifier" gorm:"column:Identifier"`
	ProvinceCode string `json:"province_code" gorm:"type:varchar(255)"`
	CityCode string `json:"city_code" gorm:"type:varchar(255)"`
	DevId	string `json:"dev_id"`
	DevType string `json:"dev_type" gorm:"type:varchar(255)"`
	HandleFlag string `json:"handle_flag" gorm:"type:varchar(255)"`
	MonitoringFlag string `json:"monitoring_flag" gorm:"type:varchar(255)"`
	SentcmdFlag string `json:"sentcmd_flag" gorm:"type:varchar(255)"`
}



// NonParamenterSettings 非参数设置可调参数表
type ParameterSerials struct {
        Parameter string `json:"parameter"`
        SerialNumber string `json:"serial_number"`
	}
type NonParamenterSetting struct {
	ApplianceId  		string `gorm:"primaryKey" json:"appliance_id"`
	MaxCurrCoeff 		string `gorm:"column:maximum_load_fan_current_deviation_coefficient;type:varchar(255);default:00" json:"maximum_load_fan_current_deviation_coefficient"`
	MinCurrCoeff 		string `gorm:"column:minimum_load_fan_current_deviation_coefficient;type:varchar(255);default:00" json:"minimum_load_fan_current_deviation_coefficient"`
	MaxDutyCycCoeff 	string `gorm:"column:maximum_load_fan_duty_cycle_deviation_coefficient;type:varchar(255);default:00" json:"maximum_load_fan_duty_cycle_deviation_coefficient"`
	MinDutyCycCoeff 	string `gorm:"column:minimum_load_fan_duty_cycle_deviation_coefficient;type:varchar(255);default:00" json:"minimum_load_fan_duty_cycle_deviation_coefficient"`
	BackwaterFlow		string `gorm:"column:backwater_flow_value;type:varchar(255);default:00" json:"backwater_flow_value"`
	FreqWind			string `gorm:"column:frequency_compensation_value_of_wind_pressure_sensor_alarm_point;type:varchar(255);default:00" json:"frequency_compensation_value_of_wind_pressure_sensor_alarm_point"`
	Updatetime 			string `gorm:"type:varchar(255);" json:"updatetime"`
}
type TowFlyParamenterSetting struct {
	ApplianceId   string `gorm:"primaryKey" json:"appliance_id"`
	BackwaterFlow string `gorm:"column:backwater_flow_value;type:varchar(255);default:00" json:"backwater_flow_value"`
	FreqWind      string `gorm:"column:frequency_compensation_value_of_wind_pressure_sensor_alarm_point;type:varchar(255);default:00" json:"frequency_compensation_value_of_wind_pressure_sensor_alarm_point"`
	Updatetime    string `gorm:"type:varchar(255);" json:"updatetime"`
}

type ParamenSetting struct {
	ApplianceId   string `gorm:"primaryKey" json:"appliance_id"`
	FA            string `gorm:"column:FA;type:varchar(255);default:00" json:"FA"`
	FF            string `gorm:"column:FF;type:varchar(255);default:00" json:"FF"`
	PH            string `gorm:"column:PH;type:varchar(255);default:00" json:"PH"`
	FH            string `gorm:"column:FH;type:varchar(255);default:00" json:"FH"`
	PL            string `gorm:"column:PL;type:varchar(255);default:00" json:"PL"`
	FL            string `gorm:"column:FL;type:varchar(255);default:00" json:"FL"`
	DH            string `gorm:"column:dH;type:varchar(255);default:00" json:"dH"`
	Fd            string `gorm:"column:Fd;type:varchar(255);default:00" json:"Fd"`
	CH            string `gorm:"column:CH;type:varchar(255);default:00" json:"CH"`
	FC            string `gorm:"column:FC;type:varchar(255);default:00" json:"FC"`
	NE            string `gorm:"column:NE;type:varchar(255);default:00" json:"NE"`
	CA            string `gorm:"column:CA;type:varchar(255);default:00" json:"CA"`
	FP            string `gorm:"column:FP;type:varchar(255);default:00" json:"FP"`
	LF            string `gorm:"column:LF;type:varchar(255);default:00" json:"LF"`
	HS            string `gorm:"column:HS;type:varchar(255);default:00" json:"HS"`
	Hb            string `gorm:"column:Hb;type:varchar(255);default:00" json:"Hb"`
	HE            string `gorm:"column:HE;type:varchar(255);default:00" json:"HE"`
	HL            string `gorm:"column:HL;type:varchar(255);default:00" json:"HL"`
	HU            string `gorm:"column:HU;type:varchar(255);default:00" json:"HU"`
	UA            string `gorm:"column:UA;type:varchar(255);default:00" json:"UA"`
	Ub            string `gorm:"column:Ub;type:varchar(255);default:00" json:"Ub"`
	Fn            string `gorm:"column:Fn;type:varchar(255);default:00" json:"Fn"`
	Updatetime    string `gorm:"type:varchar(255);" json:"updatetime"`
}



//ParameterCode 参数编码表

type ParameterCode struct {
	Parameter string
	Code	  string
}

// ParameterChangesSetting 参数变化记录表

type ParameterChangesSetting struct {
	ApplianceId		string
	Code			string
	Value			string
	LastValue	    string
	Updatetime		string
	LatestParameterFlag		string
}

// ParameterFinalSetting 参数最终表

type ParameterFinalSetting struct {
	ApplianceId				string
	Code					string
	CurrentValue			string
	RewriteSuccessFlag		string
	Updatetime				string
}

// ParameterDefaults 默认参数表

type ParameterDefault struct {
	ApplianceId			string
	Code				string
	DefaultValue		string
	Updatetime			string
}

type ConstantTempPara struct {
	ID          int    `gorm:"primary_key;auto_increment" json:"id"`
	ApplianceId string `json:"appliance_id" gorm:"type:varchar(255);unique;not null"`
	Ka0         string `json:"ka0" gorm:"type:varchar(255)"`
	Ka1         string `json:"ka1" gorm:"type:varchar(255)"`
	Ka2         string `json:"ka2" gorm:"type:varchar(255)"`
	Ka3         string `json:"ka3" gorm:"type:varchar(255)"`
	Kb0         string `json:"Kb0" gorm:"type:varchar(255)"`
	Kb1         string `json:"kb1" gorm:"type:varchar(255)"`
	Kb2         string `json:"kb2" gorm:"type:varchar(255)"`
	Kb3         string `json:"kb3" gorm:"type:varchar(255)"`
	Kc0         string `json:"kc0" gorm:"type:varchar(255)"`
	Kc1         string `json:"kc1" gorm:"type:varchar(255)"`
	Kc2         string `json:"kc2" gorm:"type:varchar(255)"`
	Kc3         string `json:"kc3" gorm:"type:varchar(255)"`
	Kf0         string `json:"kf0" gorm:"type:varchar(255)"`
	Kf1         string `json:"kf1" gorm:"type:varchar(255)"`
	Kf2         string `json:"kf2" gorm:"type:varchar(255)"`
	Kf3         string `json:"kf3" gorm:"type:varchar(255)"`
	T1a0        string `json:"t1a0" gorm:"type:varchar(255)"`
	T1a1        string `json:"t1a1" gorm:"type:varchar(255)"`
	T1a2        string `json:"t1a2" gorm:"type:varchar(255)"`
	T1a3        string `json:"t1a3" gorm:"type:varchar(255)"`
	T1c0        string `json:"t1c0" gorm:"type:varchar(255)"`
	T1c1        string `json:"t1c1" gorm:"type:varchar(255)"`
	T1c2        string `json:"t1c2" gorm:"type:varchar(255)"`
	T1c3        string `json:"t1c3" gorm:"type:varchar(255)"`
	T2a0        string `json:"t2a0" gorm:"type:varchar(255)"`
	T2a1        string `json:"t2a1" gorm:"type:varchar(255)"`
	T2a2        string `json:"t2a2" gorm:"type:varchar(255)"`
	T2a3        string `json:"t2a3" gorm:"type:varchar(255)"`
	T2c0        string `json:"t2c0" gorm:"type:varchar(255)"`
	T2c1        string `json:"t2c1" gorm:"type:varchar(255)"`
	T2c2        string `json:"t2c2" gorm:"type:varchar(255)"`
	T2c3        string `json:"t2c3" gorm:"type:varchar(255)"`
	Tda0        string `json:"tda0" gorm:"type:varchar(255)"`
	Tda1        string `json:"tda1" gorm:"type:varchar(255)"`
	Tda2        string `json:"tda2" gorm:"type:varchar(255)"`
	Tda3        string `json:"tda3" gorm:"type:varchar(255)"`
	Tdc0        string `json:"tdc0" gorm:"type:varchar(255)"`
	Tdc1        string `json:"tdc1" gorm:"type:varchar(255)"`
	Tdc2        string `json:"tdc2" gorm:"type:varchar(255)"`
	Tdc3        string `json:"tdc3" gorm:"type:varchar(255)"`
	Wc0         string `json:"wc0" gorm:"type:varchar(255)"`
	Wc1         string `json:"wc1" gorm:"type:varchar(255)"`
	Wc2         string `json:"wc2" gorm:"type:varchar(255)"`
	Wc3         string `json:"wc3" gorm:"type:varchar(255)"`
	Wo0         string `json:"wo0" gorm:"type:varchar(255)"`
	Wo1         string `json:"wo1" gorm:"type:varchar(255)"`
	Wo2         string `json:"wo2" gorm:"type:varchar(255)"`
	Wo3         string `json:"wo3" gorm:"type:varchar(255)"`
	Updatetime  string `json:"updatetime" gorm:"type:varchar(255)"`
}

func (ConstantTempPara) TableName() string {
	return "constant_temperature_parameters"
}

type ParameterSettings struct {
	ApplianceId string `json:"appliance_id" gorm:"primary_key"`
	PH0          string `json:"pH" gorm:"type:varchar(255)"`
	FH0          string `json:"fH" gorm:"type:varchar(255)"`
	PL0          string `json:"pL" gorm:"type:varchar(255)"`
	FL0          string `json:"fL" gorm:"type:varchar(255)"`
	DH0          string `json:"dh" gorm:"type:varchar(255)"`
	Fd0          string `json:"fd" gorm:"type:varchar(255)"`
	CH0          string `json:"cH" gorm:"type:varchar(255)"`
	FC0          string `json:"fC" gorm:"type:varchar(255)"`
	Updatetime  string `json:"updatetime" gorm:"type:varchar(255)"`
}

type ParameterCodes struct {
	Parameter string `json:"parameter"`
	Code      string `json:"code"`
}

type ParameterFinalSettings struct {
	ApplianceId        string `json:"appliance_id" gorm:"primary_key"`
	Code               string `json:"code" gorm:"type:varchar(255)"`
	CurrentValue       string `json:"CurrentValue " gorm:"type:varchar(255)"`
	RewriteSuccessFlag string `json:"rewrite_success_flag" gorm:"type:varchar(255)"`
	Updatetime         string `json:"updatetime" gorm:"type:varchar(255)"`
}

type ParameterDefaults struct {
	ApplianceId  string `json:"applianceId" gorm:"primary_key"`
	Code         string `json:"code" gorm:"type:varchar(255)"`
	DefaultValue string `json:"defaultValue" gorm:"type:varchar(255)"`
	Updatetime   string `json:"updatetime" gorm:"type:varchar(255)"`
}

type ParameterChangesSettings struct {
	ApplianceId         string `json:"appliance_id"`
	Code                string `json:"code" gorm:"type:varchar(255)"`
	Value               string `json:"value" gorm:"type:varchar(255)"`
	LastValue	        string `json:"last_value" gorm:"type:varchar(255)"`
	Updatetime          string `json:"updatetime" gorm:"type:varchar(255)"`
	LatestParameterFlag string `json:"latest_parameter_flag" gorm:"type:varchar(255)"`
}

type MideaDeviceAllSelectOnline struct {
	DevId           string `json:"dev_id" gorm:"varchar(255);not null"`
	DevType         string `json:"dev_type" gorm:"varchar(255);not null"`
	Sn              string `json:"sn" gorm:"type:varchar(255);not null"`
	LocationIp      string `json:"location_ip" gorm:"type:varchar(255);not null"`
	Country         string `json:"country" gorm:"type:varchar(255);not null"`
	Province        string `json:"province" gorm:"type:varchar(50);not null"`
	City            string `json:"city" gorm:"type:varchar(50);not null"`
	ActiveTime      string `json:"active_time" gorm:"type:varchar(50);not null"`
	RefreshDatetime string `json:"refresh_datetime" gorm:"type:varchar(50);not null"`
	CityCode        string `json:"city_migrate_code" gorm:"type:varchar(50);not null"`
	ModelType       string `json:"model_type" gorm:"type:varchar(50);not null"`
	ProvinceCode    string `json:"province_code" gorm:"type:varchar(50);not null"`
	CityEntireCode  string `json:"city_code" gorm:"type:varchar(50);not null"`
	MonitoringFlag  string `json:"monitoring_flag" gorm:"type:varchar(50);not null"`
	HandleFlag      string `json:"handle_flag" gorm:"type:varchar(50);not null"`
	Opt             string `json:"opt" gorm:"type:varchar(50);not null"`
	TempMultiple    string `json:"temp_multiple" gorm:"type:varchar(50);not null"`
}
type ParamHistory struct {
	DevId    string `json:"dev_id" gorm:"column:dev_id;type:varchar(20);not null;index:dev_id_idx; "`
	DataTime string `gorm:"column:data_time;type:varchar(20);not null;"`
	Wa       float64 `gorm:"column:wa;type:varchar(20);not null;"`
	Wb       float64 `gorm:"column:wb;type:varchar(20);not null;"`
	Wc       float64 `gorm:"column:wc;type:varchar(20);not null;"`
}

type  Keyparamaters  struct{
	SN                   string `gorm:"column:SN;type:varchar(255);" `
	ModelType            string `gorm:"column:model_type;type:varchar(255);default:00" json:"model_type"`
	Powerboardcode       string `gorm:"column:powerboardcode;type:varchar(255);default:00" json:"powerboardcode"`
	Displayboardcode     string `gorm:"column:displayboardcode;type:varchar(255);default:00" json:"displayboardcode"`
	SegmentInformation   string `gorm:"column:segment_information;type:varchar(255);default:00" json:"segment_information"`
	MaxSegment           string `gorm:"column:max_segment;type:varchar(255);default:00" json:"max_segment"`
	MinLoad              string `gorm:"column:min_load;type:varchar(255);default:00" json:"min_load"`
	MaxLoad              string `gorm:"column:max_load;type:varchar(255);default:00" json:"max_load"`
	FanType              string `gorm:"column:fan_type;type:varchar(255);default:00" json:"fan_type"`
	FA                   string `gorm:"column:FA;type:varchar(255);default:00" json:"FA"`
	FF                   string `gorm:"column:FF;type:varchar(255);default:00" json:"FF"`
	PH                   string `gorm:"column:PH;type:varchar(255);default:00" json:"PH"`
	FH                   string `gorm:"column:FH;type:varchar(255);default:00" json:"FH"`
	PL                   string `gorm:"column:PL;type:varchar(255);default:00" json:"PL"`
	FL                   string `gorm:"column:FL;type:varchar(255);default:00" json:"FL"`
	DH                   string `gorm:"column:DH;type:varchar(255);default:00" json:"DH"`
	FD                   string `gorm:"column:FD;type:varchar(255);default:00" json:"FD"`
	CH                   string `gorm:"column:CH;type:varchar(255);default:00" json:"CH"`
	FC                   string `gorm:"column:FC;type:varchar(255);default:00" json:"FC"`
}


var Setpar3 = []string{"FA", "FF", "PH", "FH", "PL", "FL", "dH", "Fd","CH","FC"}

func (MideaDeviceAllSelectOnline) TableName() string {
	return "midea_device_all_select_onlines"
}