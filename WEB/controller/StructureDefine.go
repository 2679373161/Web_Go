package controller

import (
	"sync"
	"time"
)

//输入数据结构体
type user struct {
	Citycode     string  `json:"citycode" gorm:"type:varchar(50);not null"`
	ProdTime     string  `json:"prod_time" gorm:"type:varchar(50);not null"`
	Applianceid  string  `json:"applianceid" gorm:"type:varchar(50);not null;index:applianceid_idx; "`
	Datatime     string  `json:"datatime" gorm:"type:varchar(50);not null;index:datatime_idx; "`
	Flame        int     `json:"flame" gorm:"type:int;not null"`
	Flow         int     `json:"flow" gorm:"type:int;not null"`
	Outtemp      int     `json:"out_temp" gorm:"type:int;not null"`
	Settemp      int     `json:"set_temp" gorm:"type:int;not null"`
	Intemp       int     `json:"in_temp" gorm:"type:int;not null"`
	EffectMark   float64 `json:"effect_mark" gorm:"type:int;not null"`
	WaterPattern int     `json:"water_pattern" gorm:"type:int;"`
	TempPattern  int     `json:"temp_pattern" gorm:"type:int;"`
	FlowUpstep   bool    `gorm:"-"` //水流向上阶跃点
	FlowDownstep bool    `gorm:"-"` //水流向下阶跃点
	Variablerise float64 `json:"variablerise" gorm:"type:int;not null"`
	Liter        string  `json:"liter" gorm:"type:varchar(50);not null"`
	Eerror       string  `json:"eerror" gorm:"type:varchar(50);not null"`
	Cerror       string  `json:"cerror" gorm:"type:varchar(50);not null"`
	NeedLiters   float64 `json:"need_liters" gorm:"type:float;not null"`
	TrueLiters   float64 `json:"true_liters" gorm:"type:float;not null"`
	BehaviorID   int     `json:"behavior_id" gorm:"type:int;not null"`
	Zerowater    string  `json:"zero_water" gorm:"type:varchar(50);not null"`
	//WorkLoad     int     `json:"work_load" gorm:"type:int;not null"`
	MinWorkLoad    bool   `json:"min_work_load" gorm:"type:bool;not null"`
	ZoneID         int    `json:"zone_id" gorm:"type:int;not null"`
	Byte38         string `gorm:"column:byte38;type:varchar(20);not null"`
	Byte14         string `gorm:"column:byte14;type:varchar(20);not null"`
	AbnormalFlag   int    `gorm:"-"` //不写入数据库
	TabName        string `gorm:"-"`
	InsertData     int    `gorm:"column:insert_data;type:int;not null"`
	InvalidReason  int    `gorm:"-"`
	InvalidReasons string `gorm:"-"`
	FlameoutCount  int    `gorm:"-"`
	Point          int    `gorm:"-"` //判断是否为波动的拐点
}

type user04 struct {
	Citycode       string `json:"citycode" gorm:"type:varchar(50);not null;index:citycode_idx;"`
	Applianceid    string `json:"applianceid" gorm:"type:varchar(50);not null;index:applianceid_idx; "`
	Datatime       string `json:"datatime" gorm:"type:varchar(50);not null;index:datatime_idx; "`
	Fantype        string `gorm:"type:varchar(50)"` //风机类型
	Comloadsegment int    `gorm:"type:varchar(50)"` //燃烧符合段
	Theorypwm      string `gorm:"type:varchar(50)"` //需求PWM
	Actualpwm      int    `gorm:"type:int"`         //实际PWM
	Fullmaxload    string `gorm:"type:varchar(50)"`
	Fullminload    string `gorm:"type:varchar(50)"`
	Windspeed      string `gorm:"type:varchar(50)"`
	Fanpwm         string `gorm:"type:varchar(50)"` //风扇PWM
	Byte22         string `gorm:"type:varchar(50)"` //分段阀信息
	Windblock      string `gorm:"type:varchar(50)"` //风堵等级
	Windpressf     string `gorm:"type:varchar(50)"` //风机转速？
}

type uau04 struct {
	Applianceid    string `json:"applianceid" gorm:"type:varchar(50);not null;index:applianceid_idx; "`
	Datatime       string `json:"datatime" gorm:"type:varchar(50);not null;index:datatime_idx; "`
	TempPattern    int    `json:"temp_pattern" gorm:"type:int;"`
	Comloadsegment int    `gorm:"type:varchar(50)"` //燃烧符合段
	Actualpwm      int    `gorm:"type:int"`         //燃烧符合段
}

// Multimodality 设备指标
type Multimodality struct {
	ProvinceCode          string  `json:"province_code" gorm:"type:varchar(6);not null"`
	CityCode              string  `json:"city_code" gorm:"type:varchar(6);not null"`
	DevType               string  `json:"dev_type" gorm:"type:varchar(10);not null"`
	DevId                 int     `json:"dev_id" gorm:"type:varchar(20);not null; "`
	TimeDate              string  `json:"total_time" gorm:"type:varchar(20);not null; "`
	WaterValidTime        string  `json:"water_valid_time" gorm:"type:varchar(20);not null"`
	WaterNum              int     `json:"water_num" gorm:"type:int;not null"`
	AverageTime           string  `json:"average_time" gorm:"type:varchar(10);not null"`
	UnStableProportion    float64 `json:"fluctuation_proportion" gorm:"type:float;not null"`
	StableTime            string  `json:"stable_time" gorm:"type:varchar(20);not null"`
	UnStableTime          string  `json:"fluctuation_time" gorm:"type:varchar(20);not null"`
	MaximumTime           string  `json:"maximum_time" gorm:"type:varchar(10);not null"`
	UnStableBehavior      int     `json:"un_stable_behavior" gorm:"type:int;not null"`
	AbnormalFlag          int     `json:"abnormal_flag" gorm:"type:int;not null"`
	WaterScore            int     `json:"water_score" gorm:"type:int;not null"`
	TempValidTime         string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum               int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration       string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev       float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	TempScore             int     `json:"temp_score" gorm:"type:float;not null"`
	NewTempScore          int     `json:"new_temp_score" gorm:"type:int;not null"`
	HeatTempScore         int     `json:"heat_temp_score" gorm:"type:float;not null"`
	StableTempScore       int     `json:"stable_temp_score" gorm:"type:float;not null"`
	NormalAllPro          float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro       float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro        float64 `json:"abnormal_all" gorm:"type:float;not null"`
	TabName               string  `gorm:"-"`
}

// ModeFragment 模式片段指标输出
type ModeFragment struct {
	DevId        string `json:"dev_id" gorm:"type:varchar(50);not null;index:dev_id_idx; "`         //设备id
	StartTime    string `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime      string `json:"end_time" gorm:"type:varchar(20);not null; "`                        //片段结束时间
	DurationTime string `json:"duration_time" gorm:"type:varchar(20);not null"`                     //片段总时长
	WaterPattern int    `json:"water_pattern" gorm:"type:int;not null"`                             //片段温度模式
	//Flowrange     int      `json:"flowrange" gorm:"type:int;not null"`
	SmallWater int `json:"small_water" gorm:"type:int;not null"` //片段小水流标志位
	//Extreme              int     `json:"extreme" gorm:"type:int;not null"`
	//MaxChange            float64 `json:"max_change" gorm:"type:float;not null"`
	//Average              float64 `json:"average" gorm:"type:float;not null"`
	Deviation            float64 `json:"deviation" gorm:"type:float;not null"`                     //片段水流变异系数
	UpNumber             int     `json:"up_number" gorm:"type:int;not null"`                       //片段用水先向上阶跃次数
	DownNumber           int     `json:"down_number" gorm:"type:int;not null"`                     //片段用水向下阶跃次数
	WaterScore           int     `json:"water_score" gorm:"int;not null"`                          //片段水流评分
	HeatDuration         string  `json:"heat-duration" gorm:"type:varchar(20);not null"`           //片段升温时间
	UnStableTempDuration string  `json:"un_stable_temp_duration" gorm:"type:varchar(20);not null"` //片段不恒温时间
	UnStableTempPercent  float64 `json:"un_stable_temp_percent" gorm:"type:float;not null"`        //片段不恒温占比
	UnHeatDev            float64 `json:"un_heat_dev" gorm:"type:float;not null"`                   //不恒温标准差
	OldUnHeatDev         float64 `json:"old_un_heat_dev" gorm:"type:float;not null"`               //旧不恒温标准差
	TempPattern          int     `json:"temp_pattern" gorm:"type:int;not null"`                    //温度模式
	OvershootValue       int     `json:"overshoot_value" gorm:"type:int;"`                         //升温超调
	StateAccuracy        int     `json:"state_accuracy" gorm:"type:int;"`                          //温度稳态精度
	TempScore            int     `json:"temp_score" gorm:"type:int;not null"`                      //旧温度评分
	NewTempScore         int     `json:"new_temp_score" gorm:"type:int;not null"`                  //新温度评分
	HeatTempScore        int     `json:"heat_temp_score" gorm:"type:int;not null"`                 //升温段评分
	StableTempScore      int     `json:"stable_temp_score" gorm:"type:int;not null"`               //恒温段评分
	OldStableTempScore   int     `json:"old_stable_temp_score" gorm:"type:int;not null"`           //旧恒温段评分
	WaterFlag            int     `json:"water_flag" gorm:"type:int;not null"`                      //水流有效标志位
	TempFlag             int     `json:"temp_flag" gorm:"type:int;not null"`                       //温度有效标志位
	EffectFlag           bool    `json:"effect_flag" gorm:"type:int;not null"`                     //片段水流对温度有影响标志位
	AbnormalState        int     `json:"abnormal_state" gorm:"type:int;not null"`                  //水流稳定温度不稳定标志位
	TabName              string  `gorm:"-"`                                                        //表名
	UnHeatDevMark        int     `json:"un_heat_dev_mark" gorm:"type:int;not null"`                //恒温段相对目标温度标准差评分
	UnStableMark         int     `json:"un_stable_mark" gorm:"type:int;not null"`                  //恒温段不恒温占比评分
	OverShootMark        int     `json:"over_shoot_mark" gorm:"type:int;not null"`                 //升温超调评分
	HeatMark             int     `json:"heat_mark" gorm:"type:int;not null"`                       //升温段评分
	TempRange            int     `json:"temp_range" gorm:"type:int;not null" `                     //温度极差
	OldunStableMark      int     `json:"old_un_stable_mark" gorm:"type:int;not null"`              //旧恒温段不恒温占比评分
	OldunHeatDevMark     int     `json:"old_un_heat_dev_mark" gorm:"type:int;not null"`            //旧恒温段相对目标温度标准差评分
	AbnormalFlag         int     `json:"abnormal_flag" gorm:"type:int;not null"`                   //异常标志位
	FlameoutCount        int     `json:"shoot_down_count" gorm:"type:int;not null"`                //熄火次数                   //再升温次数
	FluctuateCount       int     `json:"fluctuate_count" gorm:"type:int;not null"`
	//测试用，改为float64,方便sql计算
	HeatDuration_f float64 `json:"heat-duration_f" gorm:"type:float;not null"` //升温时间
	//AvgHeatDuration_f	float64  `json:"avg_heat-duration" gorm:"type:float;not null"`
	UnStableTempDuration_f float64 `json:"un_stable_temp_duration_f" gorm:"type:float;not null"` //不恒温时间
	OverTempDuration_f     float64 ` json:"over_temp_duration_f" gorm:"type:float;not null"`     //超温时间
	UnderTempDuration_f    float64 `json:"under_temp_duration_f" gorm:"type:float;not null"`     //低温时间
	ShockTempDuration_f    float64 `json:"shock_temp_duration_f" gorm:"type:float;not null"`     //震荡时间
	Settemp                int     `json:"settemp" gorm:"type:int;"`                             //片段设定温度
	Intemp                 int     `json:"Intemp" gorm:"type:int;"`                              //片段初始温度
	Outtemp                int     `json:"Outtemp" gorm:"type:int;"`                             //片段开始的输出温度
	DurationTime_f         float64 `json:"DurationTime_F" gorm:"type:int;"`                      //用水时长
	LowTempAvgFlow         float64 `json:"low_temp_avg_flow" gorm:"type:float;not null"`         //恒温段平均水流量
	LowTempAvgTemp         float64 `json:"low_temp_avg_temp" gorm:"type:float;not null"`         //恒温段平均水温，调试使用
	LiterDiff              float64 `json:"liter_diff" gorm:"type:float;not null"`                //（提供负荷 - 实际负荷）的平均差值
	MideaFaultCode         string  `json:"midea_fault_code" gorm:"type:varchar(20);not null"`
	MideaFaultCodeC        string  `json:"midea_fault_code_c" gorm:"type:varchar(20);not null"`
	Byte38                 string  `gorm:"column:byte38;type:varchar(20);not null"`
	Byte14                 string  `gorm:"column:byte14;type:varchar(20);not null"`
	InvalidReason          int     `json:"invalid_reason" gorm:"type:int;"`
	InvalidReasons         string  `json:"invalid_reasons" gorm:"type:varchar(20);not null"`
	AvgHeatFlow            float32 `gorm:"type:float;not null"`                                 //升温段平均值
	FaultCode              string  `gorm:"type:varchar(20);not null"`                           //故障code
	HeatFlameoutDuration   string  `json:"heat_flameout_duration" gorm:"type:varchar(20)"`      //升温段一开始点火时间（带s)
	HeatFinalTempDiff      int     `json:"heat_final_tempdiff" gorm:"type:int;not null"`        //升温段末尾温度差
	HeatFlameoutDuration_f float64 `json:"heat_flameout_duration_f" gorm:"type:float;not null"` //升温段一开始点火时间(不带s)
	IntempDiff             int     `json:"intemp_diff" gorm:"type:int;not null"`
	HeatFluctuateCount     int     `gorm:"type:tinyint;"` //升温段波动次数
	////TotalDuration_f float64 `json:"total_heat-duration" gorm:"type:float;not null"`
	HeatFlowFrequency    float64 `json:"heat_flow_frequency" gorm:"type:float;not null"` //升温段水流波动频率
	OverTempPercent      float64 `json:"over_temp_percent" gorm:"type:float;not null"`
	NeedLiter            float64 `json:"need_liter" gorm:"type:float;not null"`
	TrueLiter            float64 `json:"true_liter" gorm:"type:float;not null"`
	OverTemp             float64 `json:"over_temp" gorm:"type:float;not null"`
	FlowAvgChange        float64 `json:"flow_avg_change" gorm:"type:float;not null"`
	OverTempCode         string  `gorm:"type:varchar(20);not null"`
	AbnormalIgnitionFlag bool    `json:"abnormal_ignition_flag" gorm:"type:int;not null"` //点火时长过长且对升温时长有影响的标志位
	IntempFluctuate      int     `json:"Intemp_fluctuate" gorm:"type:int;"`
	IntempMax            int     `json:"Intemp_max" gorm:"type:int;"`
	IntempMin            int     `json:"Intemp_min" gorm:"type:int;"`
	UnHeatDev1           float64 `json:"un_heat_dev1" gorm:"type:float;not null"`       //不恒温标准差
	UnStablePercent1     float64 `json:"un_stable_percent1" gorm:"type:float;not null"` //片段不恒温占比
	AvgTempDiff1         float64 `json:"avg_temp_diff1" gorm:"type:float;not null"`     //不恒温处平均温差
	InsertPercent        float64 `json:"insert_percent" gorm:"type:float;not null"`     //不恒温处插值比例
	OverShootDuration_f  float64 `gorm:"type:int;not null"`                             //超调时间
	RiseTime             float64 `gorm:"type:int;not null"`                             //上升时间
	HeatLowTempFlag      bool    `gorm:"type:tinyint;not null"`                         //升温段低温标志位
	DevType              string  `gorm:"type:varchar(25);not null"`
	FlowGrade            string  `json:"flow_grade" gorm:"type:varchar(20);not null"`
	IntempGrade          string  `json:"intemp_grade" gorm:"type:varchar(20);not null"`
	ScoreByWeight        int     `json:"score_by_weight" gorm:"int;not null"`
	HeatWeight           float64 `json:"heat_weight" gorm:"float;not null"`
	HeatPercent          float64 `json:"heat_percent" gorm:"float;not null"`
	LiterGrade           string  `gorm:"type:varchar(20);"`
	FragmentGrade        int     `gorm:"type:int;"`
}

// DaysSummary 设备日指标汇总
type DaysSummary struct {
	ProvinceCode            string  `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode                string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType                 string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId                   string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx; "`
	TimeDate                string  `json:"time_date" gorm:"type:varchar(20);not null;index:time_date_idx; "`
	WaterValidTime          string  `json:"water_valid_time" gorm:"type:varchar(20);not null"`
	WaterNum                int     `json:"water_num" gorm:"type:int;not null"`
	AverageTime             string  `json:"average_time" gorm:"type:varchar(10);not null"`
	UnStableProportion      float64 `json:"fluctuation_proportion" gorm:"type:float;not null"`
	UnStableTime            string  `json:"unStableTime" gorm:"type:varchar(20);not null"`
	MaximumTime             string  `json:"maximum_time" gorm:"type:varchar(10);not null"`
	UnStableBehavior        int     `json:"un_stable_behavior" gorm:"type:int;not null"`
	WaterScore              int     `json:"water_score" gorm:"type:int;not null"`
	TempValidTime           string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum                 int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration         string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration      string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion   float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev         float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	OldTempScore            int     `json:"temp_score" gorm:"type:int;not null"`
	TempScore               int     `json:"new_temp_score" gorm:"type:int;not null"` //新温度评分
	HeatTempScore           int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore         int     `json:"stable_temp_score" gorm:"type:int;not null"`
	UnStableMark            int     `json:"un_stable_mark" gorm:"type:int;not null"`
	UnHeatDevMark           int     `json:"un_heat_dev_mark" gorm:"type:int;not null"`
	OverShootMark           int     `json:"over_shoot_mark" gorm:"type:int;not null"`
	HeatMark                int     `json:"heat_mark" gorm:"type:int;not null"`
	TempAllNormalNum        int     `json:"temp_all_normal_num" gorm:"type:int;not null`
	ConstantTempAbnormalNum int     `json:"constant_temp_abnormal_num" gorm:"type:int;not null`
	ElevateTempAbnormalNum  int     `json:"elevate_temp_abnormal_num" gorm:"type:int;not null`
	TempAllAbnormalNum      int     `json:"temp_all_abnormal_num" gorm:"type:int;not null`
	NormalAllPro            float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro         float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro   float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro          float64 `json:"abnormal_all" gorm:"type:float;not null"`
	AbnormalCount           int     `json:"abnormal_count" gorm:"type:int;not null"`
	//AbnormalProportion    float64 `json:"abnormal_proportion" gorm:"type:float;not null"`
	AbnormalFlag   int    `json:"abnormal_flag" gorm:"type:int;not null"`
	TempRange      int    `json:"temp_range" gorm:"type:int;not null" `
	OvershootValue int    `json:"overshoot_value" gorm:"type:int;"` //升温超调
	Num1100        int    `json:"num1100" gorm:"type:int;"`         //设备自身问题性能不佳
	Num1000        int    `json:"num1000" gorm:"type:int;"`         //设备自身问题，不恒温占比大、极差小
	Num0100        int    `json:"num0100" gorm:"type:int;"`         //设备自身问题，不恒温占比小、极差大
	NumFluctuate   int    `json:"num_fluctuate" gorm:"type:int;"`   //水流波动原因
	TabName        string `gorm:"-"`
}
type MonthSummary struct {
	ProvinceCode          string  `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode              string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType               string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId                 string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	TimeDate              string  `json:"total_time" gorm:"type:varchar(20);not null;index:time_date_idx"`
	WaterValidTime        string  `json:"water_valid_time" gorm:"type:varchar(20);not null"`
	WaterNum              int     `json:"water_num" gorm:"type:int;not null"`
	EffectiveDays         int     `json:"effective_days" gorm:"type:int;not null"`
	AverageTime           string  `json:"average_time" gorm:"type:varchar(10);not null"`
	UnStableProportion    float64 `json:"fluctuation_proportion" gorm:"type:float;not null"`
	UnStableTime          string  `json:"fluctuation_time" gorm:"type:varchar(20);not null"`
	MaximumTime           string  `json:"maximum_time" gorm:"type:varchar(10);not null"`
	UnStableBehavior      int     `json:"un_stable_behavior" gorm:"type:int;not null"`
	WaterScore            int     `json:"water_score" gorm:"type:int;not null"`
	TempValidTime         string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum               int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration       string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev       float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	TempScore             int     `json:"temp_score" gorm:"type:float;not null"`
	HeatTempScore         int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore       int     `json:"stable_temp_score" gorm:"type:int;not null"`
	NormalAllPro          float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro       float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro        float64 `json:"abnormal_all" gorm:"type:float;not null"`
	AbnormalCount         int     `json:"abnormal_count" gorm:"type:int;not null"`
	AbnormalProportion    float64 `json:"abnormal_proportion" gorm:"type:float;not null"`
	TabName               string  `gorm:"-"`
	TempScoreByDay        int     `json:"temp_score_by_day" gorm:"type:int;not null"`
	Num1100               int     `json:"num1100" gorm:"type:int;"`       //设备自身问题性能不佳
	Num1000               int     `json:"num1000" gorm:"type:int;"`       //设备自身问题，不恒温占比大、极差小
	Num0100               int     `json:"num0100" gorm:"type:int;"`       //设备自身问题，不恒温占比小、极差大
	NumFluctuate          int     `json:"num_fluctuate" gorm:"type:int;"` //水流波动原因
}
type BadFragment struct {
	CityCode         string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType          string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId            string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime        string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime          string  `json:"end_time" gorm:"type:varchar(20);not null; "`                        //片段结束时间
	UnHeatDev1       float64 `json:"un_heat_dev1" gorm:"type:float;not null"`                            //不恒温标准差
	UnStablePercent1 float64 `json:"un_stable_percent1" gorm:"type:float;not null"`                      //片段不恒温占比
	AvgTempDiff1     float64 `json:"avg_temp_diff1" gorm:"type:float;not null"`                          //不恒温处平均温差
	InsertPercent    float64 `json:"insert_percent" gorm:"type:float;not null"`                          //不恒温处插值比例
	FaultCode        string  `json:"fault_code" gorm:"type:varchar(20);not null; "`
	OverTempCode     string  `gorm:"type:varchar(20);not null"`
	InvalidReasons   string  `json:"invalid_reasons" gorm:"type:varchar(20);not null"`
}
type OverTemp struct {
	CityCode       string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType        string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId          string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime      string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime        string  `json:"end_time" gorm:"type:varchar(20);not null; "`                        //片段结束时间
	OverTemp       float64 `json:"over_temp" gorm:"type:float;not null"`
	OverPercent    float64 `json:"over_percent" gorm:"type:float;not null"`
	Deviation      float64 `json:"deviation" gorm:"type:float;not null"` //片段水流变异系数
	FluctuateCount int     `json:"fluctuate_count" gorm:"type:int;not null"`
	FlowAvgChange  float64 `json:"flow_avg_change" gorm:"type:float;not null"`
	OverCode       string  `json:"over_code" gorm:"type:varchar(20);not null; "`
}
type LowTemp struct {
	CityCode       string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType        string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId          string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime      string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime        string  `json:"end_time" gorm:"type:varchar(20);not null; "`                        //片段结束时间
	Settemp        int     `json:"settemp" gorm:"type:int;"`                                           //片段设定温度
	LowTempAvgTemp float64 `json:"low_temp_avg_temp" gorm:"type:float;not null"`                       //恒温段平均水温，调试使用
	LowTempAvgFlow float64 `json:"low_temp_avg_flow" gorm:"type:float;not null"`                       //恒温段平均水流量
	LiterDiff      float64 `json:"liter_diff" gorm:"type:float;not null"`                              //（提供负荷 - 实际负荷）的平均差值
}

type FaultE0 struct {
	ProvinceCode     string `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode         string `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType          string `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId            string `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime        string `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "`
	TimeDate         string `json:"time_date" gorm:"type:varchar(20);not null;index:time_date_idx; "` //统计截止日期
	FaultCode        string `json:"fault_code" gorm:"type:varchar(20);not null"`                      //故障编码
	TotalNum         int    `json:"total" gorm:"type:int;not null"`                                   //参与故障评定的总片段数
	FlameoutNum      int    `json:"flame_out_num" gorm:"type:int;not null"`                           //熄火片段个数
	MaxFlameoutCount int    `json:"max_flame_out_count" gorm:"type:int;not null"`                     //片段最大熄火次数
}
type FaultE1 struct {
	ProvinceCode   string  `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode       string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType        string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId          string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime      string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //统计开始日期
	TimeDate       string  `json:"time_date" gorm:"type:varchar(20);not null;index:time_date_idx; "`   //统计截止日期
	FaultCode      string  `json:"fault_code" gorm:"type:varchar(20);not null"`                        //故障编码
	TotalNum       int     `json:"total" gorm:"type:int;not null"`                                     //参与故障评定的总片段数
	LowTempNum     int     `json:"low_temp_num" gorm:"type:int;not null"`                              //无法恒温片段个数
	SetTemp        int     `json:"set_temp" gorm:"type:int;not null"`                                  //无法恒温时的设定温度（取最小值）
	LowTempMinFlow int     `json:"low_temp_flow" gorm:"type:int;not null"`                             //在设定温度下，无法恒温的最小水流
	TempDiff       int     `json:"temp_diff" gorm:"type:int;not null"`                                 //最小水流处对应的温差(设定温度-平均水温)
	LiterDiff      float64 `json:"liter_diff" gorm:"type:float;not null"`                              //（提供负荷 - 实际负荷）的平均差值
}
type FaultE4 struct {
	ProvinceCode string  `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode     string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType      string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId        string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime    string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //统计开始日期
	TimeDate     string  `json:"time_date" gorm:"type:varchar(20);not null;index:time_date_idx; "`   //统计截止日期
	FaultCode    string  `json:"fault_code" gorm:"type:varchar(20);not null"`                        //故障编码
	TotalNum     int     `json:"total" gorm:"type:int;not null"`                                     //参与故障评定的总片段数
	OverTempNum  int     `json:"over_temp_num" gorm:"type:int;not null"`                             //设备当天总超温片段
	AbnormalNum  int     `json:"abnormal_num" gorm:"type:int;not null"`                              //无外因，正常运行环境还超温个数
	FlowNum      int     `json:"flow_num" gorm:"type:int;not null"`                                  //水流原因超温个数
	IntempNum    int     `json:"intemp_num" gorm:"type:int;not null"`                                //进水温度异常原因超温个数
	MinLoadNum   int     `json:"min_load_num" gorm:"type:int;not null"`                              //最小工况超温个数
	AvgOverTemp  float64 `json:"avg_over_temp" gorm:"type:float;not null"`                           //平均超温值

}

//超调
type FaultE3 struct {
	ProvinceCode    string  `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode        string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType         string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId           string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime       string  ` gorm:"type:varchar(20);not null;index:start_time_idx; "` //统计开始日期
	TimeDate        string  `json:"time_date" gorm:"type:varchar(20);not null; "`      //统计截止日期
	FaultCode       string  `json:"fault_code" gorm:"type:varchar(20);not null"`       //故障编码
	CurrentLevelNum int16   ` gorm:"type:smallint;not null"`                           //该设备当前场景片段数
	TotalNum        int16   `json:"total" gorm:"type:smallint;not null"`               //该设备总片段数
	OvershootNum    int16   `json:"low_temp_num" gorm:"type:smallint;not null"`        //当前场景超调个数
	Osp             int8    `gorm:"type:smallint;not null"`                            //超调占比=超调个数/当前场景个数   OvershootNum/CurrentLevelNum
	AvgOvershoot    float32 `gorm:"type:float;not null"`                               //平均超调两
	AvgHeatFlow     float32 `json:"set_temp" gorm:"type:float;not null"`               //平均水流量
	AvgTempDiff     float32 `json:"temp_diff" gorm:"type:float;not null"`              //平均温差(设定温度-平均水温)
}

//升温时间长
type FaultE2 struct {
	ProvinceCode    string  `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode        string  `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType         string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId           string  `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx"`
	StartTime       string  ` gorm:"type:varchar(20);not null;index:start_time_idx; "` //统计开始日期
	TimeDate        string  `json:"time_date" gorm:"type:varchar(20);not null; "`      //统计截止日期
	FaultCode       string  `json:"fault_code" gorm:"type:varchar(20);not null"`       //故障编码
	CurrentLevelNum int16   ` gorm:"type:smallint;not null"`                           //该设备当前场景片段数
	TotalNum        int16   `json:"total" gorm:"type:smallint;not null"`               //该设备总片段数
	OvertimeNum     int16   `json:"low_temp_num" gorm:"type:smallint;not null"`        //该设备升温超市个数
	Otp             int8    `gorm:"type:smallint;not null"`                            //升温超长占比 Otp=OverTimeNum/CurrentLevelNum
	AvgOvertime     float32 `gorm:"type:float;not null"`                               //平均升温时长
	AvgHeatFlow     float32 `json:"set_temp" gorm:"type:float;not null"`               //平均水流量
	AvgTempDiff     float32 `json:"temp_diff" gorm:"type:float;not null"`              //平均温差(设定温度-平均水温)
}

type DevTypeFaultStruct struct {
	DevType   string `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	FaultCode string `json:"fault_code" gorm:"type:varchar(20);not null"`  //故障编码
	Count     int    `gorm:"type:tinyint"`                                 //机型个数
	TimeDate  string `json:"time_date" gorm:"type:varchar(20);not null; "` //统计截止日期
}

// MultimodalityCity 城市汇总指标
type MultimodalityCity struct {
	ProvinceCode          string  `json:"province_code" gorm:"type:varchar(50);not null"`
	CityCode              string  `json:"city_code" gorm:"type:varchar(50);not null"`
	EquipmentNum          int     `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate              string  `json:"time_date" gorm:"type:varchar(50);not null"`
	WaterValidTime        string  `json:"water_valid_time" gorm:"type:varchar(50);not null"`
	WaterNum              int     `json:"water_num" gorm:"type:varchar(50);not null"`
	AverageTime           string  `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion    string  `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	UnStableTime          string  `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime           string  `json:"maximum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior      int     `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	WaterScore            int     `json:"water_score" gorm:"type:int(50);not null"`
	TempValidTime         string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum               int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration       string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev       float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	TempScore             int     `json:"temp_score" gorm:"type:float;not null"`
	HeatTempScore         int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore       int     `json:"stable_temp_score" gorm:"type:int;not null"`
	NormalAllPro          float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro       float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro        float64 `json:"abnormal_all" gorm:"type:float;not null"`
	TabName               string  `gorm:"-"`
}

// MultimodalityProvince 省汇总指标
type MultimodalityProvince struct {
	ProvinceCode          string  `json:"province_code" gorm:"type:varchar(50);not null"`
	EquipmentNum          int     `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate              string  `json:"time_date" gorm:"type:varchar(50);not null"`
	WaterValidTime        string  `json:"water_valid_time" gorm:"type:varchar(50);not null"`
	WaterNum              int     `json:"water_num" gorm:"type:varchar(50);not null"`
	AverageTime           string  `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion    string  `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	UnStableTime          string  `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime           string  `json:"maximum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior      int     `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	WaterScore            int     `json:"water_score" gorm:"type:int(50);not null"`
	TempValidTime         string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum               int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration       string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev       float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	TempScore             int     `json:"temp_score" gorm:"type:float;not null"`
	HeatTempScore         int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore       int     `json:"stable_temp_score" gorm:"type:int;not null"`
	NormalAllPro          float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro       float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro        float64 `json:"abnormal_all" gorm:"type:float;not null"`
	TabName               string  `gorm:"-"`
}
type MultimodalityRegion struct {
	RegionCode            string  `json:"province_code" gorm:"type:varchar(50);not null"`
	EquipmentNum          int     `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate              string  `json:"time_date" gorm:"type:varchar(50);not null"`
	WaterValidTime        string  `json:"water_valid_time" gorm:"type:varchar(50);not null"`
	WaterNum              int     `json:"water_num" gorm:"type:varchar(50);not null"`
	AverageTime           string  `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion    string  `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	UnStableTime          string  `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime           string  `json:"maximum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior      int     `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	WaterScore            int     `json:"water_score" gorm:"type:int(50);not null"`
	TempValidTime         string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum               int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration       string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev       float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	TempScore             int     `json:"temp_score" gorm:"type:float;not null"`
	HeatTempScore         int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore       int     `json:"stable_temp_score" gorm:"type:int;not null"`
	NormalAllPro          float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro       float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro        float64 `json:"abnormal_all" gorm:"type:float;not null"`
	TabName               string  `gorm:"-"`
}

// MultimodalityType 型号汇总指标
type MultimodalityType struct {
	DevType               string  `json:"dev_type" gorm:"type:varchar(50);not null"`
	EquipmentNum          int     `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate              string  `json:"time_date" gorm:"type:varchar(50);not null"`
	WaterValidTime        string  `json:"water_valid_time" gorm:"type:varchar(50);not null"`
	WaterNum              int     `json:"water_num" gorm:"type:varchar(50);not null"`
	AverageTime           string  `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion    string  `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	UnStableTime          string  `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime           string  `json:"maximum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior      int     `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	WaterScore            int     `json:"water_score" gorm:"type:int(50);not null"`
	TempValidTime         string  `json:"temp_valid_time" gorm:"type:varchar(20);not null"`
	TempNum               int     `json:"temp_num" gorm:"type:int;not null"`
	AveHeatDuration       string  `json:"ave_heat_duration" gorm:"type:varchar(20);not null"`
	AveUnSableDuration    string  `json:"ave_un_sable_duration" gorm:"type:varchar(20);not null"`
	AveUnstableProportion float64 `json:"ave_unstable_proportion" gorm:"type:float;not null"`
	UnStableTempDev       float64 `json:"un_stable_temp_dev" gorm:"type:float;not null"`
	TempScore             int     `json:"temp_score" gorm:"type:float;not null"`
	HeatTempScore         int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore       int     `json:"stable_temp_score" gorm:"type:int;not null"`
	NormalAllPro          float64 `json:"normal_all" gorm:"type:float;not null"`
	AbnormalHeatPro       float64 `json:"abnormal_heat" gorm:"type:float;not null"`
	AbnormalStableTempPro float64 `json:"abnormal_stable_temp" gorm:"type:float;not null"`
	AbnormalAllPro        float64 `json:"abnormal_all" gorm:"type:float;not null"`
	TabName               string  `gorm:"-"`
}

// Equipment 设备id表
type Equipment struct {
	Identifier     string `json:"identifier" gorm:"type:varchar(20);not null"`
	ProvinceCode   string `json:"province_code" gorm:"type:varchar(20);not null"`
	CityCode       string `json:"city_code" gorm:"type:varchar(10);not null"`
	DevId          string `json:"dev_id" gorm:"type:varchar(20);not null"`
	DevType        string `json:"dev_type" gorm:"type:varchar(20);not null"`
	HandleFlag     string `json:"handle_flag" gorm:"type:varchar(20);not null"`
	MonitoringFlag int    `json:"monitoring_flag" gorm:"type:int;not null"`
	Liter          string `json:"liter" gorm:"type:varchar(20);not null;"`
}

type ModeBehavior struct {
	ProvinceCode string `json:"province_code" gorm:"type:varchar(50);not null"`
	CityCode     string `json:"city_code" gorm:"type:varchar(50);not null"`
	DevId        string `json:"dev_id" gorm:"type:varchar(20);not null; "`
	StartTime    string `json:"start_time" gorm:"type:varchar(20);not null; "`
	EndTime      string `json:"end_time" gorm:"type:varchar(20);not null; "`
	DurationTime string `json:"duration_time" gorm:"type:varchar(20);not null"`
	FragmentNum  int    `json:"fragment_num" gorm:"type:varchar(20);not null"`
	EffectFlag   int    `json:"effect" gorm:"type:int;not null"`
}

//用水行为
type BehaviorSummary struct {
	ProvinceCode string `json:"province_code" gorm:"type:varchar(50);not null"`
	CityCode     string `json:"city_code" gorm:"type:varchar(50);not null"`
	DevId        string `json:"dev_id" gorm:"type:varchar(50);not null; "`
	DataTime     string `json:"start_time" gorm:"type:varchar(50);not null; "`
	Sec0p        int    `json:"sec0P" gorm:"type:int;not null"`
	Sec30p       int    `json:"sec30P" gorm:"type:int;not null"`
	Min3p        int    `json:"min3P" gorm:"type:int;not null"`
	Min10p       int    `json:"min10P" gorm:"type:int;not null"`
}

//数据特征
type DataFeature struct {
	StartTime           string  `json:"start_time" gorm:"type:varchar(50);not null"`
	UpdateTime          string  `json:"update_time" gorm:"type:varchar(50);not null"`
	ProvinceCount       int     `json:"province_count" gorm:"type:int;not null"`
	CityCount           int     `json:"city_count" gorm:"type:int;not null"`
	TypeCount           int     `json:"type_count" gorm:"type:int;not null"`
	EquCount            int     `json:"equ_count" gorm:"type:int;not null"`
	DayDataSum          int64   `json:"day_data_sum" gorm:"type:bigint;not null"`
	AvgDataSum          int64   `json:"avg_data_sum" gorm:"type:bigint;not null"`
	TotalDataSum        int64   `json:"total_data_sum" gorm:"type:bigint;not null"`
	HandleDuration      string  `json:"handle_duration" gorm:"type:varchar(50);not null"`
	FreeSpace           float64 `json:"free_space" gorm:"type:float;"`
	DataDays            int     `json:"data_days" gorm:"type:int;not null"`
	AbnormalCount       int     `json:"abnormal_count" gorm:"type:int;not null"`
	DayAllDevAvgScore   int     `json:"day_all_dev_avg_score" gorm:"type:int;not null"`
	UnderEightFive      int64   `json:"under_eight_five" gorm:"type:bigint;not null"`    //高于60低于85分的
	TopEightFive        int64   `json:"top_eight_five" gorm:"type:bigint;not null"`      //高于85的
	UnderSixty          int64   `json:"under_sixty" gorm:"type:bigint;not null"`         //低于60分的
	OneHundred          int64   `json:"one_hundred" gorm:"type:bigint;not null"`         //100分
	DevAvgHeatScore     int     `json:"dev_avg_heat_score" gorm:"type:int;not null"`     //平均升温段评分
	DevAvgStableScore   int     `json:"dev_avg_stable_score" gorm:"type:int;not null"`   //平均恒温温段评分
	AvgWaterNum         int64   `json:"avg_water_num" gorm:"type:bigint;not null"`       //平均用水次数
	AvgWaterTime        string  `json:"avg_water_time" gorm:"type:varchar(50);not null"` //平均用水时长
	AvgHeatTime         string  `json:"avg_heat_time" gorm:"type:varchar(50);not null"`  //平均升温时长
	TotalProcessingTime string  `json:"total_processing_time" gorm:"type:varchar(50);not null"`
}

//日设备监控
type DailyMonitoring struct {
	ProvinceCode                                             string `json:"province_code" gorm:"type:varchar(6);not null;index:province_code_idx"`
	CityCode                                                 string `json:"city_code" gorm:"type:varchar(6);not null;index:city_code_idx"`
	DevType                                                  string `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId                                                    string `json:"dev_id" gorm:"type:varchar(20);not null;index:dev_id_idx; "`
	TimeDate                                                 string `json:"time_date" gorm:"type:varchar(20);not null;index:time_date_idx; "`
	TempScore                                                int    `json:"temp_score" gorm:"type:int;not null"`
	HeatTempScore                                            int    `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore                                          int    `json:"stable_temp_score" gorm:"type:int;not null"`
	AbnormalFlag                                             int    `json:"abnormal_flag" gorm:"type:int;not null"`
	PH                                                       int    `json:"PH" gorm:"type:int;not null;default:0"`
	PL                                                       int    `json:"PL" gorm:"type:int;not null;default:0"`
	FH                                                       int    `json:"FH" gorm:"type:int;not null;default:0"`
	DH                                                       int    `json:"dH" gorm:"type:int;not null;default:0"`
	Fd                                                       int    `json:"Fd" gorm:"type:int;not null;default:0"`
	CH                                                       int    `json:"CH" gorm:"type:int;not null;default:0"`
	FC                                                       int    `json:"FC" gorm:"type:int;not null;default:0"`
	FL                                                       int    `json:"FL" gorm:"type:int;not null;default:0"`
	MaximumLoadFanCurrentDeviationCoefficient                int    `json:"maximum_load_fan_current_deviation_coefficient" gorm:"type:int;not null;default:0"`
	MinimumLoadFanCurrentDeviationCoefficient                int    `json:"minimum_load_fan_current_deviation_coefficient" gorm:"type:int;not null;default:0"`
	MaximumLoadFanDutyCycleDeviationCoefficient              int    `json:"maximum_load_fan_duty_cycle_deviation_coefficient" gorm:"type:int;not null;default:0"`
	MinimumLoadFanDutyCycleDeviationCoefficient              int    `json:"minimum_load_fan_duty_cycle_deviation_coefficient" gorm:"type:int;not null;default:0"`
	BackwaterFlowValue                                       int    `json:"backwater_flow_value" gorm:"type:int;not null;default:0"`
	FrequencyCompensationValueOfWindPressureSensorAlarmPoint int    `json:"frequency_compensation_value_of_wind_pressure_sensor_alarm_point" gorm:"type:int;not null;default:0"`
	KA0                                                      int    `json:"ka_0" gorm:"type:int;not null;default:0"`
	KA1                                                      int    `json:"ka_1" gorm:"type:int;not null;default:0"`
	KA2                                                      int    `json:"ka_2" gorm:"type:int;not null;default:0"`
	KA3                                                      int    `json:"ka_3" gorm:"type:int;not null;default:0"`
	KB0                                                      int    `json:"kb_0" gorm:"type:int;not null;default:0"`
	KB1                                                      int    `json:"kb_1" gorm:"type:int;not null;default:0"`
	KB2                                                      int    `json:"kb_2" gorm:"type:int;not null;default:0"`
	KB3                                                      int    `json:"kb_3" gorm:"type:int;not null;default:0"`
	KC0                                                      int    `json:"kc_0" gorm:"type:int;not null;default:0"`
	KC1                                                      int    `json:"kc_1" gorm:"type:int;not null;default:0"`
	KC2                                                      int    `json:"kc_2" gorm:"type:int;not null;default:0"`
	KC3                                                      int    `json:"kc_3" gorm:"type:int;not null;default:0"`
	KF0                                                      int    `json:"kf_0" gorm:"type:int;not null;default:0"`
	KF1                                                      int    `json:"kf_1" gorm:"type:int;not null;default:0"`
	KF2                                                      int    `json:"kf_2" gorm:"type:int;not null;default:0"`
	KF3                                                      int    `json:"kf_3" gorm:"type:int;not null;default:0"`
	T1A0                                                     int    `json:"t_1_a_0" gorm:"type:int;not null;default:0"`
	T1A1                                                     int    `json:"t_1_a_1" gorm:"type:int;not null;default:0"`
	T1A2                                                     int    `json:"tia_2" gorm:"type:int;not null;default:0"`
	T1A3                                                     int    `json:"t_1_a_3" gorm:"type:int;not null;default:0"`
	T1C0                                                     int    `json:"t_1_c_0" gorm:"type:int;not null;default:0"`
	T1C1                                                     int    `json:"t_1_c_1" gorm:"type:int;not null;default:0"`
	T1C2                                                     int    `json:"t_1_c_2" gorm:"type:int;not null;default:0"`
	T1C3                                                     int    `json:"t_1_c_3" gorm:"type:int;not null;default:0"`
	T2A0                                                     int    `json:"t_2_a_0" gorm:"type:int;not null;default:0"`
	T2A1                                                     int    `json:"t_2_a_1" gorm:"type:int;not null;default:0"`
	T2A2                                                     int    `json:"t_2_a_2" gorm:"type:int;not null;default:0"`
	T2A3                                                     int    `json:"t_2_a_3" gorm:"type:int;not null;default:0"`
	T2C0                                                     int    `json:"t_2_c_0" gorm:"type:int;not null;default:0"`
	T2C1                                                     int    `json:"t_2_c_1" gorm:"type:int;not null;default:0"`
	T2C2                                                     int    `json:"t_2_c_2" gorm:"type:int;not null;default:0"`
	T2C3                                                     int    `json:"t_2_c_3" gorm:"type:int;not null;default:0"`
	TDA0                                                     int    `json:"tda_0" gorm:"type:int;not null;default:0"`
	TDA1                                                     int    `json:"tda_1" gorm:"type:int;not null;default:0"`
	TDA2                                                     int    `json:"tda_2" gorm:"type:int;not null;default:0"`
	TDA3                                                     int    `json:"tda_3" gorm:"type:int;not null;default:0"`
	TDC0                                                     int    `json:"tdc_0" gorm:"type:int;not null;default:0"`
	TDC1                                                     int    `json:"tdc_1" gorm:"type:int;not null;default:0"`
	TDC2                                                     int    `json:"tdc_2" gorm:"type:int;not null;default:0"`
	TDC3                                                     int    `json:"tdc_3" gorm:"type:int;not null;default:0"`
	WC0                                                      int    `json:"wc_0" gorm:"type:int;not null;default:0"`
	WC1                                                      int    `json:"wc_1" gorm:"type:int;not null;default:0"`
	WC2                                                      int    `json:"wc_2" gorm:"type:int;not null;default:0"`
	WC3                                                      int    `json:"wc_3" gorm:"type:int;not null;default:0"`
	WO0                                                      int    `json:"wo_0" gorm:"type:int;not null;default:0"`
	WO1                                                      int    `json:"wo_1" gorm:"type:int;not null;default:0"`
	WO2                                                      int    `json:"wo_2" gorm:"type:int;not null;default:0"`
	WO3                                                      int    `json:"wo_3" gorm:"type:int;not null;default:0"`
	E1                                                       int    `json:"e1" gorm:"type:int;not null;default:0"`
	C4                                                       int    `json:"c4" gorm:"type:int;not null;default:0"`
}

type ParameterOutput struct {
	Appliance_id       int    `json:"appliance_id" gorm:"type:int;not null"`
	Code               string `json:"code" gorm:"type:string;not null"`
	CurrentValue       int    `json:"current_value" gorm:"type:int;not null"`
	RewriteSuccessFlag int    `json:"rewrite_success_flag" gorm:"type:int;not null"`
	Updatetime         string `json:"updatetime" gorm:"type:varchar(50);not null"`
}
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

//及时监控设备表
type InstantEquipment struct {
	Applianceid      string `json:"applianceid" gorm:"type:varchar(50);not null"`
	StartTime        string `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx"`
	EndTime          string `json:"end_time" gorm:"type:varchar(20);not null;"`
	Move_task_flag   int    `json:"move_task_flag" gorm:"type:tinyint;not null;default:0"`
	Mining_task_flag int    `json:"mining_task_flag" gorm:"type:tinyint;not null;default:0"`
	AdjustTaskFlag   int    `json:"adjust_task_flag" gorm:"type:tinyint;not null;default:0"`
}

// WaterFeature 水流量特征
type WaterFeatureStruct struct {
	DurationTime       string  `json:"duration_time" gorm:"type:varchar(20);not null"` //片段总时长
	WaterPattern       int     `json:"water_pattern" gorm:"type:int;not null"`         //片段温度模式
	SmallWater         int     `json:"small_water" gorm:"type:int;not null"`           //片段小水流标志位
	Deviation          float64 `json:"deviation" gorm:"type:float;not null"`           //片段水流标准差
	UpNumber           int     `json:"up_number" gorm:"type:int;not null"`             //片段用水先向上阶跃次数
	DownNumber         int     `json:"down_number" gorm:"type:int;not null"`           //片段用水向下阶跃次数
	WaterScore         int     `json:"water_score" gorm:"int;not null"`                //片段水流评分
	AbnormalState      int     `json:"abnormal_state" gorm:"type:int;not null"`        //水流稳定温度不稳定标志位
	WaterFlag          int     `json:"water_flag" gorm:"type:int;not null"`
	FluctuateCount     int     `json:"fluctuate_count" gorm:"type:int;not null"` //突变次数
	Avg                float64 `json:"avg" gorm:"type:float;not null"`
	Flowrange          int     `json:"avg_water" gorm:"type:int;not null"`
	LowTempAvgFlow     float64 `json:"low_temp_avg_flow" gorm:"type:float;not null"` //恒温段平均水流量
	LowTempAvgTemp     float64 `json:"low_temp_avg_temp" gorm:"type:float;not null"` //恒温段平均水温，调试使用
	FaultCode          string  `json:"fault_code" gorm:"type:varchar(20);not null"`  //故障编码
	LiterDiff          float64 `json:"liter_diff" gorm:"type:float;not null"`        //（提供负荷 - 实际负荷）
	Coefficient        float64 //变异系数
	AvgHeatFlow        float32 `gorm:"type:float;not null"` //升温段平均值
	HeatFluctuateCount int     `gorm:"type:tinyint;"`       //升温段波动次数
	OverTempPercent    float64 `json:"over_temp_percent" gorm:"type:float;not null"`
	OverTemp           float64 `json:"over_temp" gorm:"type:float;not null"`
	FlowChange         float64 `json:"flow_change" gorm:"type:float;not null"`
	FlowGrade          string  `json:"flow_grade" gorm:"type:varchar(20);not null"`
	AvgNeedLiter       float64 ` gorm:"type:float;"`
	LiterGrade         string  `gorm:"type:varchar(20);"`
}
type TempFeatureStruct struct {
	DurationTime         string        `json:"duration_time" gorm:"type:varchar(20);not null"`           //片段总时长
	HeatDuration         time.Duration `json:"heat-duration" gorm:"type:varchar(20);not null"`           //片段升温时间
	UnStableTempDuration time.Duration `json:"un_stable_temp_duration" gorm:"type:varchar(20);not null"` //片段不恒温时间
	UnStableTempPercent  float64       `json:"un_stable_temp_percent" gorm:"type:float;not null"`        //片段不恒温占比
	OldStableTempPercent float64
	OldTempRange         int
	FlameoutCount        int     //熄火次数
	UnHeatDev            float64 `json:"un_heat_dev" gorm:"type:float;not null"`         //不恒温标准差
	OldUnHeatDev         float64 `json:"old_un_heat_dev" gorm:"type:float;not null"`     //旧不恒温标准差
	TempPattern          int     `json:"temp_pattern" gorm:"type:int;not null"`          //温度模式
	OvershootValue       int     `json:"overshoot_value" gorm:"type:int;"`               //升温超调
	StateAccuracy        int     `json:"state_accuracy" gorm:"type:int;"`                //温度稳态精度
	TempScore            int     `json:"temp_score" gorm:"type:int;not null"`            //旧温度评分
	NewTempScore         int     `json:"new_temp_score" gorm:"type:int;not null"`        //新温度评分
	HeatTempScore        int     `json:"heat_temp_score" gorm:"type:int;not null"`       //升温段评分
	StableTempScore      int     `json:"stable_temp_score" gorm:"type:int;not null"`     //恒温段评分
	OldStableTempScore   int     `json:"old_stable_temp_score" gorm:"type:int;not null"` //旧恒温段评分
	TempFlag             int     `json:"temp_flag" gorm:"type:int;not null"`             //温度有效标志位
	EffectFlag           bool    `json:"effect_flag" gorm:"type:int;not null"`           //片段水流对温度有影响标志位
	AbnormalState        int     `json:"abnormal_state" gorm:"type:int;not null"`        //水流稳定温度不稳定标志位
	UnHeatDevMark        int     `json:"un_heat_dev_mark" gorm:"type:int;not null"`      //恒温段相对目标温度标准差评分
	UnStableMark         int     `json:"un_stable_mark" gorm:"type:int;not null"`        //恒温段不恒温占比评分
	OverShootMark        int     `json:"over_shoot_mark" gorm:"type:int;not null"`       //升温超调评分
	HeatMark             int     `json:"heat_mark" gorm:"type:int;not null"`             //升温段评分
	TempRange            int     `json:"temp_range" gorm:"type:int;not null" `           //温度极差
	OldunStableMark      int     `json:"old_un_stable_mark" gorm:"type:int;not null"`    //旧恒温段不恒温占比评分
	OldunHeatDevMark     int     `json:"old_un_heat_dev_mark" gorm:"type:int;not null"`  //旧恒温段相对目标温度标准差评分
	//测试用，改为float64,方便sql计算
	HeatDuration_f float64 `json:"heat-duration_f" gorm:"type:float;not null"` //升温段开火时候不过冲时长
	//AvgHeatDuration_f	float64  `json:"avg_heat-duration" gorm:"type:float;not null"`
	UnStableTempDuration_f float64       `json:"un_stable_temp_duration_f" gorm:"type:float;not null"`    //不恒温时间
	OldOverTempDuration_f  float64       `json:"over_temp_duration_f" gorm:"type:float;not null"`         //超温时间
	OldUnderTempDuration_f float64       `json:"under_temp_duration_f" gorm:"type:float;not null"`        //低温时间
	OldShockTempDuration_f float64       `json:"shock_temp_duration_f" gorm:"type:float;not null"`        //震荡时间
	Settemp                int           `json:"settemp" gorm:"type:int;"`                                //片段设定温度
	Intemp                 int           `json:"Intemp" gorm:"type:int;"`                                 //片段初始温度
	Outtemp                int           `json:"Outtemp" gorm:"type:int;"`                                //片段开始的输出温度
	HeatFlameoutDuration   time.Duration `json:"heat_flameout_duration" gorm:"type:varchar(20);not null"` //升温段一开始点火时长
	HeatFinalTempDiff      int           `json:"heat_final_tempdiff" gorm:"type:int;not null"`            //升温段末尾温度差
	HeatFlameoutDuration_f float64       `json:"heat_flameout_duration_f" gorm:"type:float;not null"`

	OverShootDuration_f float64 `gorm:"type:int;not null"`
	RiseTime            float64 `gorm:"type:int;not null"`
	IntempGrade         string  `json:"intemp_grade" gorm:"type:varchar(20);not null"`
	IntempDiff          int     `json:"intemp_diff" gorm:"type:int;not null"`
	IntempFluctuate     int     `json:"Intemp_fluctuate" gorm:"type:int;"`
	IntempMax           int     `json:"Intemp_max" gorm:"type:int;"`
	IntempMin           int     `json:"Intemp_min" gorm:"type:int;"`
	UnHeatDev1          float64 `json:"un_heat_dev1" gorm:"type:float;not null"`       //不恒温标准差
	UnStablePercent1    float64 `json:"un_stable_percent1" gorm:"type:float;not null"` //片段不恒温占比
	AvgTempDiff1        float64 `json:"avg_temp_diff1" gorm:"type:float;not null"`     //不恒温处平均温差
	HeatLowTempFlag     bool    `gorm:"type:tinyint;not null"`                         //升温段低温标志位
	InsertPercent       float64 `json:"insert_percent" gorm:"type:float;not null"`     //不恒温处插值比例
	////TotalDuration_f float64 `json:"total_heat-duration" gorm:"type:float;not null"`
}

//水流量识别中间标记结构体
type PatternRec struct {
	deltaflow float64
	up        int
	down      int
}

//锁的实现
type Locker struct {
	lock sync.Mutex
}

func NewUser(tableName string) *user {
	return &user{TabName: tableName}
}
func (u user) TableName() string {
	return u.TabName
}
func NewUser2(tableName string) *ModeFragment {
	return &ModeFragment{TabName: tableName}
}
func (u ModeFragment) TableName() string {
	return u.TabName
}
func ProvinceName(tableName string) *MultimodalityProvince {
	return &MultimodalityProvince{TabName: tableName}
}
func (u MultimodalityProvince) TableName() string {
	return u.TabName
}

func RegionName(tableName string) *MultimodalityRegion {
	return &MultimodalityRegion{TabName: tableName}
}
func (u MultimodalityRegion) TableName() string {
	return u.TabName
}

func CityName(tableName string) *MultimodalityCity {
	return &MultimodalityCity{TabName: tableName}
}
func (u MultimodalityCity) TableName() string {
	return u.TabName
}
func TypeName(tableName string) *MultimodalityType {
	return &MultimodalityType{TabName: tableName}
}
func (u MultimodalityType) TableName() string {
	return u.TabName
}

// OperationDataSheetName 运行数据表名字
func OperationDataSheetName(CityCode string) string {
	valueStart := "data"
	dataSheetName := CityCode[0:4]
	OperationDataSheet := valueStart + dataSheetName

	return OperationDataSheet
}

// DataOutPutTableName 数据输出表名
func DataOutPutTableName(dataSheetName string) string {
	value := "00"
	finalValue := dataSheetName + value
	return finalValue
}

// PatternFragmentTableName 模式片段表名
func PatternFragmentTableName(patternName string) string {
	value := "fragment"
	finalValue := value + patternName
	return finalValue
}

func ProvinceCodeName1(provinceName string) string {
	value := "0000"
	provinceName = provinceName[0:2]
	finalValue := provinceName + value
	return finalValue
}

// ProvinceCode 省编号
func ProvinceCode(provinceName string) string {
	value := "0000"
	provinceName = provinceName[0:2]
	finalValue := provinceName + value
	return finalValue
}

// DuplicateDataProcess 重复数据处理
func DuplicateDataProcess(allData []string) []string {

	var uniqueData []string
	var repeatFlag bool

	for i, n := 0, len(allData); i < n; i++ {
		for j, m := 0, i; j < m; j++ {
			if allData[i] == allData[j] {
				repeatFlag = true
				break
			}
		}

		if repeatFlag == false {
			uniqueData = append(uniqueData, allData[i])
		}
		repeatFlag = false
	}
	return uniqueData
}

func (l *Locker) Lock(key string) (success bool, err error) {
	l.lock.Lock()

	return true, nil
}

func (l *Locker) Unlock(key string) error {
	l.lock.Unlock()
	return nil
}

/************以下为数据表**********以下为数据表********以下为数据表***********以下为数据表*******************/

//1. WeekDates 周表
type WeekDate struct {
	DevType              string `json:"dev_tyoe" gorm:"varchar(255);not null"`
	DevId                string `json:"dev_id" gorm:"varchar(255);not null"`
	TimeDate             string `json:"time_date" gorm:"varchar(255);not null"`
	TempScore            int    `json:"temp_score" gorm:"int;not null"`
	TempAllNormal        int    `json:"temp_all_normal" gorm:"int;not null"`
	ConstantTempAbnormal int    `json:"constant_temp_abnormal" gorm:"int;not null"`
	ElevateTempAbnormal  int    `json:"elevate_temp_abnormal" gorm:"int;not null"`
	TempAllAbnormal      int    `json:"temp_all_abnormal" gorm:"int;not null"`
	AvgOvershootValue    int    `json:"avg_overshoot_value" gorm:"int;not null"`
	AvgStateAccuracy     int    `json:"avg_state_accuracy" gorm:"int;not null"`
	HeatDurationLong     int    `json:"heat_duration_long" gorm:"int;not null"`
	OvershootValueBig    int    `json:"overshoot_value_big" gorm:"int;not null"`
	LongBothBig          int    `json:"long_both_big" gorm:"int;not null"`
	UnderTimeDuration    int    `json:"under_time_duration" gorm:"int;not null"`
	OverTimeDuration     int    `json:"over_time_duration" gorm:"int;not null"`
	TemprangeBig         int    `json:"temprange_big" gorm:"int;not null"`
	WaterEffectNum       int    `json:"water_effect_num" gorm:"int;not null"`
	FlameEffectNum       int    `json:"flame_effect_num" gorm:"int;not null"`
	WandfEffectNum       int    `json:"wandf_effect_num" gorm:"int;not null"`
}

//2. TempAbnormalMode 设备异常片段
type TempAbnormalMode struct {
	DevId       string `json:"dev_id" gorm:"type:varchar(50);not null"`
	CityCode    string `json:"city_code" gorm:"type:varchar(6);not null"`
	StartTime   string `json:"start_time" gorm:"type:varchar(50);not null"`
	EndTime     string `json:"end_time" gorm:"type:varchar(50);not null"`
	TempPattern int    `json:"temp_pattern" gorm:"type:int;not null"`
}

//3. TypeOnedayAvgscore 型号表
type TypeOnedayAvgscore struct {
	TimeDate                string `json:"time_date" gorm:"primary_key"`
	DevType                 string `json:"dev_type" gorm:"primary_key"`
	Avgscore                string `json:"avgscore" gorm:"varchar(50);not null"`
	AvgOvershootValue       string `json:"avg_overshoot_value" gorm:"varchar(50);not null"` //平均超调
	AvgStateAccuracy        string `json:"avg_state_accuracy" gorm:"varchar(50);not null"`  //平均稳态精度
	AvgHeatScore            string `json:"avg_heat_score" gorm:"varchar(50);not null"`      //平均升温段评分
	AvgStableScore          string `json:"avg_stable_score" gorm:"varchar(50);not null"`    //平均恒温段评分
	TempAllNormalNum        string `json:"temp_all_normal_num" gorm:"varchar(50);not null"`
	ConstantTempAbnormalNum string `json:"constant_temp_abnormal_num" gorm:"varchar(50);not null"`
	ElevateTempAbnormalNum  string `json:"elevate_temp_abnormal_num" gorm:"varchar(50);not null"`
	TempAllAbnormalNum      string `json:"temp_all_abnormal_num" gorm:"varchar(50);not null"`
}

//4.省份表
type ProvinceOnedayAverageScore struct {
	TimeDate     string `json:"time_date" gorm:"type:varchar(20);primary_key"`
	ProvinceCode string `json:"procince_code" gorm:"type:varchar(6);primary_key"`
	AvgScore     int    `json:"avg_score" gorm:"type:int;not null"`
}

//5.城市表
type CityOnedayAverageScore struct {
	TimeDate              string `json:"time_date" gorm:"type:varchar(20);primary_key"`
	CityCode              string `json:"city_code" gorm:"type:varchar(6);primary_key"`
	AvgScore              string `json:"avg_score" gorm:"type:varchar(20);not null"`
	AveUnStableProportion string `json:"ave_un_stable_proportion" gorm:"type:varchar(20);not null"`
}

func (CityOnedayAverageScore) TableName() string {
	return "city_oneday_average_scores"
}

//6. RegionOnedayAverageScores地区表
type RegionOnedayAverageScores struct {
	TimeDate              string `json:"time_date" gorm:"type:varchar(20);primary_key"`
	RegionCode            string `json:"region_code" gorm:"type:varchar(6);primary_key"`
	AvgScore              string `json:"avg_score" gorm:"type:varchar(20);not null"`
	AveUnStableProportion string `json:"ave_un_stable_proportion" gorm:"type:varchar(20);not null"`
	UnderEightFive        string `json:"under_eight_five" gorm:"type:varchar(20);not null"` //高于60低于85分的
	TopEightFive          string `json:"top_eight_five" gorm:"type:varchar(20);not null"`   //高于85的
	UnderSixty            string `json:"under_sixty" gorm:"type:varchar(20);not null"`      //低于60分的
	OneHundred            string `json:"one_hundred" gorm:"type:varchar(20);not null"`      //100分
}

func (RegionOnedayAverageScores) TableName() string {
	return "region_oneday_average_scores"
}

//7.型号、城市、升温段得分、恒温段得分在各个分数段内的数量
//type:01型号，02城市、03恒温段、04升温段
type DistributionTable struct {
	TimeDate         string `json:"time_date" gorm:"varchar(255);primary_key"`
	Type             string `json:"type" gorm:"varchar(255);primary_key"`
	TopEightFive     int    `json:"top_eight_five" `
	SixtyToEightFive int    `json:"sixty_to_eight_five" `
	UnderSixty       int    `json:"under_sixty `
}

//8.地区省份排名表
type NationalRankingTables struct {
	Date               string `json:"date" gorm:"varchar(255);not null"`
	FirstRegion        string `json:"first_region" gorm:"type:varchar(20);not null"`
	FirstRegionCode    string `json:"first_region_code" gorm:"type:varchar(20);not null"`
	SecondRegion       string `json:"second_region" gorm:"type:varchar(20);not null"`
	SecondRegionCode   string `json:"second_region_code" gorm:"type:varchar(20);not null"`
	ThirdRegion        string `json:"third_region" gorm:"type:varchar(20);not null"`
	ThirdRegionCode    string `json:"third_region_code" gorm:"type:varchar(20);not null"`
	FirstProvince      string `json:"first_province" gorm:"type:varchar(20);not null"`
	FirstProvinceCode  string `json:"first_province_code" gorm:"type:varchar(20);not null"`
	SecondProvince     string `json:"second_province" gorm:"type:varchar(20);not null"`
	SecondProvinceCode string `json:"second_province_code" gorm:"type:varchar(20);not null"`
	ThirdProvince      string `json:"third_province" gorm:"type:varchar(20);not null"`
	ThirdProvinceCode  string `json:"third_province_code" gorm:"type:varchar(20);not null"`
	FourthProvince     string `json:"fourth_province" gorm:"type:varchar(20);not null"`
	FourthProvinceCode string `json:"fourth_province_code" gorm:"type:varchar(20);not null"`
	FifthProvince      string `json:"fifth_province" gorm:"type:varchar(20);not null"`
	FifthProvinceCode  string `json:"fifth_province_code" gorm:"type:varchar(20);not null"`
}

//9.低分设备表
type LowScoreEquRanking struct {
	Date      string `json:"date" gorm:"varchar(255);not null"`
	Province  string `json:"Province" gorm:"type:varchar(20);not null"`
	FirstEqu  string `json:"first_equ" gorm:"type:varchar(20);not null"`
	SecondEqu string `json:"second_equ" gorm:"type:varchar(20);not null"`
	ThirdEqu  string `json:"third_equ" gorm:"type:varchar(20);not null"`
	FourthEqu string `json:"fourth_equ" gorm:"type:varchar(20);not null"`
	FifthEqu  string `json:"fifth_equ" gorm:"type:varchar(20);not null"`
	SixthEqu  string `json:"sixth_equ" gorm:"type:varchar(20);not null"`
}

//10.设备型号表
type ModelTypeNum struct {
	ModelType string `json:"model_type" gorm:"varchar(255);primary_key"`
	Num       string `json:"num" gorm:"type:varchar(50);not null"`
	Opt       string `json:"opt" gorm:"type:varchar(50);not null"`
}

//11.全国地区表
type MideaLocCode struct {
	DevProvince  string `json:"dev_province" gorm:"type:varchar(255);not null"`
	ProvinceCode string `json:"province_code" gorm:"type:varchar(255);not null"`
	DevCity      string `json:"dev_city" gorm:"type:varchar(255);not null"`
	CityCode     string `json:"city_code" gorm:"type:varchar(255);not null"`
	RegionCode   string `json:"region_code" gorm:"type:varchar(255);not null"`
	DevRegion    string `json:"dev_region" gorm:"type:varchar(255);not null"`
}

//12.城市选择表
type CitySelectCode struct {
	Code string `json:"code" gorm:"varchar(255);not null"`
	Opt  int    `json:"opt" gorm:"type:int;not null"`
}

/************************************前端查询用表 ******************************************/
type EquipmentsearchTable struct {
	Devid              string `json:"dev_id" gorm:"type:varchar(50);not null"`
	Devtype            string `json:"dev_type" gorm:"type:varchar(50);not null"`
	Timedate           string `json:"time_date" gorm:"type:varchar(50);not null"`
	Tempscore          int    `json:"temp_score" gorm:"type:varchar(50);not null"`
	Unstableproportion string `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	Tempvalidtime      string `json:"temp_valid_time" gorm:"type:varchar(50);not null"`
	Aveheatduration    string `json:"ave_heat_duration" gorm:"type:varchar(50);not null"`
	Aveunsableduration string `json:"ave_un_sable_duration" gorm:"type:varchar(50);not null"`
	Tempnum            string `json:"temp_num" gorm:"type:varchar(50);not null"`
}

func (EquipmentsearchTable) TableName() string {
	return "equipmentsearch"
}

type FragmentMmonitoringsTable struct {
	ProvinceCode string `json:"province_code" gorm:"type:varchar(6);not null"`
	CityCode     string `json:"city_code" gorm:"type:varchar(255);not null"`
	DevType      string `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	DevId        string `json:"dev_id" gorm:"varchar(255);not null"`
	TimeDate     string `json:"time_date" gorm:"varchar(255)"`
	TempScore    int    `json:"temp_score" gorm:"type:int;not null"`
	AbnormalFlag int    `json:"abnormal_flag" gorm:"type:int;not null"`
}

func (FragmentMmonitoringsTable) TableName() string {
	return "fragment_monitorings"
}

type tablefragment struct {
	DevId                string  `json:"dev_id" gorm:"varchar(255);not null"`
	DevType              string  `json:"dev_type" gorm:"type:varchar(10);not null;index:dev_type_idx"`
	WaterPattern         int     `json:"water_pattern" gorm:"type:int;"`
	TempPattern          int     `json:"temp_pattern" gorm:"type:int;"`
	StartTime            string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "`
	EndTime              string  `json:"end_time" gorm:"type:varchar(20);not null; "`
	DurationTime         string  `json:"duration_time" gorm:"type:varchar(20);not null"`
	Extreme              int     `json:"extreme" gorm:"type:int;not null"`
	MaxChange            float64 `json:"max_change" gorm:"type:float;not null"`
	Average              float64 `json:"average" gorm:"type:float;not null"`
	Deviation            float64 `json:"deviation" gorm:"type:float;not null"`                     //片段水流变异系数
	UpNumber             int     `json:"up_number" gorm:"type:int;not null"`                       //片段用水先向上阶跃次数
	DownNumber           int     `json:"down_number" gorm:"type:int;not null"`                     //片段用水向下阶跃次数
	WaterScore           int     `json:"water_score" gorm:"int;not null"`                          //片段水流评分
	HeatDuration         string  `json:"heat-duration" gorm:"type:varchar(20);not null"`           //片段升温时间
	UnStableTempDuration string  `json:"un_stable_temp_duration" gorm:"type:varchar(20);not null"` //片段不恒温时间
	UnStableTempPercent  float64 `json:"un_stable_temp_percent" gorm:"type:float;not null"`        //片段不恒温占比
	UnHeatDev            float64 `json:"un_heat_dev" gorm:"type:float;not null"`                   //不恒温标准差
	StateAccuracy        int     `json:"state_accuracy" gorm:"type:int;"`                          //温度稳态精度
	EffectFlag           bool    `json:"effect_flag" gorm:"type:int;not null"`                     //片段水流对温度有影响标志位
	NewTempScore         int     `json:"new_temp_scoreF" gorm:"type:int;not null"`
}

func (tablefragment) TableName() string {
	return "tablefragments"
}

type DefaultInfo struct {
	DateTime        string `json:"date_time" gorm:"column:更新时间;type:varchar(20);not null"`
	SN8             string `json:"sn_8" gorm:"column:SN8码;type:varchar(20);not null"`
	DevID           string `json:"dev_id" gorm:"column:设备号;type:varchar(20);not null"`
	StartTime       string `json:"start_time" gorm:"column:用水开始时间;type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime         string `json:"end_time" gorm:"column:用水结束时间;type:varchar(20);not null; "`
	HeatMark        int    `json:"heat_mark" gorm:"column:升温时长评分;type:int;not null"` //升温段评分
	OverShootMark   int    `json:"over_shoot_mark" gorm:"column:温度过冲评分;type:int;not null"`
	UnHeatDevMark   int    `json:"un_heat_dev_mark" gorm:"column:恒温段波动评分;type:int;not null"`
	UnStableMark    int    `json:"un_stable_mark" gorm:"column:恒温段不恒温占比评分;type:int;not null"`
	StableTempScore int    `json:"stable_temp_score" gorm:"column:恒温总评分;type:int;not null"`
	FlameAbnormal   string `json:"flame_abnormal" gorm:"column:异常类别;type:varchar(20);not null"`
	//DefaultCode     string   `json:"default_code" gorm:"column:故障代码;type:varchar(20);not null"`
}
type DefaultInfo1 struct {
	DateTime  string `json:"date_time" gorm:"column:更新时间;type:varchar(20);not null"`
	SN8       string `json:"sn_8" gorm:"column:SN8码;type:varchar(20);not null"`
	DevID     string `json:"dev_id" gorm:"column:设备号;type:varchar(20);not null"`
	StartTime string `json:"start_time" gorm:"column:用水开始时间;type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime   string `json:"end_time" gorm:"column:用水结束时间;type:varchar(20);not null; "`
}
type DefaultInfo2 struct {
	DateTime     string `json:"date_time" gorm:"column:更新时间;type:varchar(20);not null"`
	SN8          string `json:"sn_8" gorm:"column:SN8码;type:varchar(20);not null"`
	DevID        string `json:"dev_id" gorm:"column:设备号;type:varchar(20);not null"`
	StartTime    string `json:"start_time" gorm:"column:用水开始时间;type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime      string `json:"end_time" gorm:"column:用水结束时间;type:varchar(20);not null; "`
	DefaultCode  string `json:"default_code" gorm:"column:故障代码E;type:varchar(20);not null"`
	DefaultCodeC string `json:"default_code_c" gorm:"column:故障代码C;type:varchar(20);not null"`
	Byte38       string `gorm:"column:byte38;type:varchar(20);not null"`
	Byte14       string `gorm:"column:byte14;type:varchar(20);not null"`
	FlameOut     string `gorm:"column:是否熄火;type:varchar(20);not null"`
}
type Citycode struct {
	Code string `json:"code" gorm:"type:varchar(20);not null"`
	Opt  int    `json:"opt" gorm:"type:int;not null"`
}
type Neufault struct {
	DateTime     string `json:"date_time" gorm:"column:date_time;type:varchar(20);not null"`
	ProvinceCode string `json:"province_code" gorm:"type:varchar(6);not null"`
	CityCode     string `json:"city_code" gorm:"type:varchar(6);not null"`
	DevType      string `json:"dev_type" gorm:"type:varchar(10);not null"`
	DevId        string `json:"dev_id" gorm:"type:varchar(20);not null; "`
}
type Ufault struct {
	SN8       string `gorm:"type:varchar(20);not null; "`
	DevId     string `json:"dev_id" gorm:"type:varchar(50);not null;index:dev_id_idx; "`         //设备id
	StartTime string `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime   string `json:"end_time" gorm:"type:varchar(20);not null; "`
	Errorr    string `json:"errorr" gorm:"type:varchar(20);not null; "`
}

type FaultsE2E3Infostruct struct {
	LowLoc              int8
	Size                int8
	LowLocStr           string
	SizeStr             string
	DayStartTimeStr     string
	StartTimeStr        string
	EndTimeStr          string
	percentageThreshold int
	CountThreshold      int
}
type MideaFault2 struct {
	DevID    string `json:"dev_id" gorm:"type:varchar(50);not null"`
	ProdTime string `json:"prod_time" gorm:"type:varchar(50);not null"`
	DevType  string `json:"dev_type" gorm:"type:varchar(50);not null"`
	TimeDate string `json:"time_date" gorm:"type:varchar(50);not null"`
	E0       int    `json:"e0" gorm:"type:int;not null"`
	E1       int    `json:"e1" gorm:"type:int;not null"`
	E2       int    `json:"e2" gorm:"type:int;not null"`
	E3       int    `json:"e3" gorm:"type:int;not null"`
	E4       int    `json:"e4" gorm:"type:int;not null"`
	E5       int    `json:"e5" gorm:"type:int;not null"`
	E6       int    `json:"e6" gorm:"type:int;not null"`
	E8       int    `json:"e8" gorm:"type:int;not null"`
	EA       int    `json:"ea" gorm:"type:int;not null"`
	EE       int    `json:"ee" gorm:"type:int;not null"`
	F2       int    `json:"f2" gorm:"type:int;not null"`
	C0       int    `json:"c0" gorm:"type:int;not null"`
	C1       int    `json:"c1" gorm:"type:int;not null"`
	C2       int    `json:"c2" gorm:"type:int;not null"`
	C3       int    `json:"c3" gorm:"type:int;not null"`
	C4       int    `json:"c4" gorm:"type:int;not null"`
	C5       int    `json:"c5" gorm:"type:int;not null"`
	C6       int    `json:"c6" gorm:"type:int;not null"`
	C7       int    `json:"c7" gorm:"type:int;not null"`
	C8       int    `json:"c8" gorm:"type:int;not null"`
	EH       int    `json:"eH" gorm:"type:int;not null"`
	EF       int    `json:"eF" gorm:"type:int;not null"`
}
type MideaTypeFault struct {
	DevType  string `json:"dev_type" gorm:"type:varchar(50);not null"`
	TimeDate string `json:"time_date" gorm:"type:varchar(50);not null"`
	E0       int    `json:"e0" gorm:"type:int;not null"`
	E1       int    `json:"e1" gorm:"type:int;not null"`
	E2       int    `json:"e2" gorm:"type:int;not null"`
	E3       int    `json:"e3" gorm:"type:int;not null"`
	E4       int    `json:"e4" gorm:"type:int;not null"`
	E5       int    `json:"e5" gorm:"type:int;not null"`
	E6       int    `json:"e6" gorm:"type:int;not null"`
	E8       int    `json:"e8" gorm:"type:int;not null"`
	EA       int    `json:"ea" gorm:"type:int;not null"`
	EE       int    `json:"ee" gorm:"type:int;not null"`
	F2       int    `json:"f2" gorm:"type:int;not null"`
	C0       int    `json:"c0" gorm:"type:int;not null"`
	C1       int    `json:"c1" gorm:"type:int;not null"`
	C2       int    `json:"c2" gorm:"type:int;not null"`
	C3       int    `json:"c3" gorm:"type:int;not null"`
	C4       int    `json:"c4" gorm:"type:int;not null"`
	C5       int    `json:"c5" gorm:"type:int;not null"`
	C6       int    `json:"c6" gorm:"type:int;not null"`
	C7       int    `json:"c7" gorm:"type:int;not null"`
	C8       int    `json:"c8" gorm:"type:int;not null"`
	EH       int    `json:"eH" gorm:"type:int;not null"`
	EF       int    `json:"eF" gorm:"type:int;not null"`
}
type MideaprodtimeFault struct {
	TimeDate  string `json:"time_date" gorm:"type:varchar(50);not null"`
	ProdMonth string `json:"prod_month" gorm:"type:varchar(50);not null"`
	E0        int    `json:"e0" gorm:"type:int;not null"`
	E1        int    `json:"e1" gorm:"type:int;not null"`
	E2        int    `json:"e2" gorm:"type:int;not null"`
	E3        int    `json:"e3" gorm:"type:int;not null"`
	E4        int    `json:"e4" gorm:"type:int;not null"`
	E5        int    `json:"e5" gorm:"type:int;not null"`
	E6        int    `json:"e6" gorm:"type:int;not null"`
	E8        int    `json:"e8" gorm:"type:int;not null"`
	EA        int    `json:"ea" gorm:"type:int;not null"`
	EE        int    `json:"ee" gorm:"type:int;not null"`
	F2        int    `json:"f2" gorm:"type:int;not null"`
	C0        int    `json:"c0" gorm:"type:int;not null"`
	C1        int    `json:"c1" gorm:"type:int;not null"`
	C2        int    `json:"c2" gorm:"type:int;not null"`
	C3        int    `json:"c3" gorm:"type:int;not null"`
	C4        int    `json:"c4" gorm:"type:int;not null"`
	C5        int    `json:"c5" gorm:"type:int;not null"`
	C6        int    `json:"c6" gorm:"type:int;not null"`
	C7        int    `json:"c7" gorm:"type:int;not null"`
	C8        int    `json:"c8" gorm:"type:int;not null"`
	EH        int    `json:"eH" gorm:"type:int;not null"`
	EF        int    `json:"eF" gorm:"type:int;not null"`
}
type MideaNumFault struct {
	DevType  string `json:"dev_type" gorm:"type:varchar(50);not null"`
	TimeDate string `json:"time_date" gorm:"type:varchar(50);not null"`
	Total    int    `json:"total" gorm:"type:varchar(50);not null"`
	E0num    int    `json:"e0num" gorm:"type:int;not null"`
	E1num    int    `json:"e1num" gorm:"type:int;not null"`
	E2num    int    `json:"e2num" gorm:"type:int;not null"`
	E3num    int    `json:"e3num" gorm:"type:int;not null"`
	E4num    int    `json:"e4num" gorm:"type:int;not null"`
	E5num    int    `json:"e5num" gorm:"type:int;not null"`
	E6num    int    `json:"e6num" gorm:"type:int;not null"`
	E8num    int    `json:"e8num" gorm:"type:int;not null"`
	EANUM    int    `json:"eanum" gorm:"type:int;not null"`
	EENUM    int    `json:"eenum" gorm:"type:int;not null"`
	F2num    int    `json:"f2num" gorm:"type:int;not null"`
	C0num    int    `json:"c0num" gorm:"type:int;not null"`
	C1num    int    `json:"c1num" gorm:"type:int;not null"`
	C2num    int    `json:"c2num" gorm:"type:int;not null"`
	C3num    int    `json:"c3num" gorm:"type:int;not null"`
	C4num    int    `json:"c4num" gorm:"type:int;not null"`
	C5num    int    `json:"c5num" gorm:"type:int;not null"`
	C6num    int    `json:"c6num" gorm:"type:int;not null"`
	C7num    int    `json:"c7num" gorm:"type:int;not null"`
	C8num    int    `json:"c8num" gorm:"type:int;not null"`
	EHNUM    int    `json:"ehnum" gorm:"type:int;not null"`
	EFNUM    int    `json:"efnum" gorm:"type:int;not null"`
}
type CutValve struct {
	DateTime   string `json:"date_time" gorm:"column:更新时间;type:varchar(20);not null"`
	SN8        string `json:"sn_8" gorm:"column:SN8码;type:varchar(20);not null"`
	DevID      string `json:"dev_id" gorm:"column:设备号;type:varchar(20);not null"`
	TimeDate   string `json:"time_date" gorm:"column:切阀时间;type:varchar(20);not null;index:time_date_idx; "` //开始时间
	CutValve   string `json:"cut_valve" gorm:"column:切阀情况;type:varchar(20);not null"`
	OverFlag   string `json:"over_flag" gorm:"column:切阀前状态;type:varchar(20);not null"`
	Waterstate string `json:"water_state" gorm:"column:水流量波动情况;type:varchar(20);not null"`
}
type Userinfo struct {
	Citycode       string `json:"citycode" gorm:"type:varchar(50);not null"`
	Applianceid    string `json:"applianceid" gorm:"type:varchar(50);not null;index:applianceid_idx; "`
	Datatime       string `json:"datatime" gorm:"type:varchar(50);not null;index:datatime_idx; "`
	TempPattern    int    `json:"temp_pattern" gorm:"type:int;"`
	Comloadsegment string `json:"comloadsegment" gorm:"type:string;not null"`
	Flow           int    `json:"flow" gorm:"type:string;not null"`
}

//预警设备表
type EarlyWarningDev struct {
	DevId                string
	City                 string
	DevType              string  `json:"dev_type" gorm:"type:varchar(20);not null"`
	StartTime            string  `json:"start_time" gorm:"type:varchar(20);not null"`
	EndTime              string  `json:"end_time" gorm:"type:varchar(20);not null"`
	TempScoreDownPercent float64 `json:"temp_score_down_percent" gorm:"type:float;not null"`
	TempScoreAvg         float64 `json:"temp_score_avg" gorm:"type:float;not null"`
	UpperLimit           float64 `json:"upper_limit" gorm:"type:float;not null"`
	LowerLimit           float64 `json:"lower_limit" gorm:"type:float;not null"`
	TempscoreDeviation   float64 `json:"tempscore_deviation" gorm:"type:float;not null"`
}
type SliceOfPercent struct {
	firstOfSecond  []float64
	secondOfSecond []float64
	thirdOfSecond  []float64
	fourthOfSecond []float64
}
type ParamTable struct {
	SN          string `gorm:"column:SN;type:varchar(20);not null"`
	MachineType string `gorm:"column:machine_type;type:varchar(20);not null"`
	MinLoad     int    `gorm:"column:min_load;type:int;not null"`
	MaxLoad     int    `gorm:"column:max_load;type:int;not null"`
	PH          string `gorm:"column:PH;type:varchar(10);not null"`
	PL          string `gorm:"column:PL;type:varchar(10);not null"`
	MaxSegment  int    `gorm:"column:max_segment;type:int;not null"`
}

//type PerchangeNums struct{
//	 DevId string `json:"dev_id" `
//	 FaultType string `gorm:"column:fault_type;type:varchar(20);not null"`
//	 DataTime   string `gorm:"column:data_time;type:varchar(20);not null"`
//	 ChangeCount int `gorm:"column:change_num;type:int;not null"`
//}
type FragmentInfo struct {
	DateTime  string `json:"date_time" gorm:"column:更新时间;type:varchar(20);not null"`
	SN8       string `json:"sn_8" gorm:"column:SN8码;type:varchar(20);not null"`
	DevID     string `json:"dev_id" gorm:"column:设备号;type:varchar(20);not null"`
	StartTime string `json:"start_time" gorm:"column:用水开始时间;type:varchar(20);not null;index:start_time_idx; "` //开始时间
	EndTime   string `json:"end_time" gorm:"column:用水结束时间;type:varchar(20);not null; "`
	Info      int    `json:"info" gorm:"column:提示信息;type:varchar(20);not null; "`
}

// 参数变化记录
type PerchangeNums struct {
	DevId       string `json:"dev_id" `
	Pername     string `json:"pername" gorm:"type:varchar(255)"`
	UpValue     string `json:"up_value " gorm:"type:varchar(255)"`
	CurreValue  string `json:"curre_value" gorm:"type:varchar(255)"`
	UpdataTime  string `json:"updata_time" gorm:"type:varchar(255)"`
	SuccessFlag string `json:"success_flag" gorm:"type:varchar(255)"`
	FaultType   string `gorm:"type:varchar(255)"`
}

// 参数变化次数记录
type NumPerChanges struct {
	FaultType string `json:"fault_type" gorm:"type:varchar(255)"`
	DevId     string `json:"dev_id" gorm:"type:varchar(255)"`
	ChangeNum string `json:"change_num" gorm:"type:varchar(255)"`
	DataTime  string `gorm:"column:data_time;type:varchar(20);not null"`
}

//需要上门维修的故障设备表
type FailedDevice struct {
	DevId      string
	DevType    string
	FaultType  string
	DataTime   string
	ReviseFlag int
}
type Real_mideaFault struct {
	DevID      string `json:"dev_id" gorm:"type:varchar(50);not null"`
	DevType    string `json:"dev_type" gorm:"type:varchar(10);not null"`
	Start_time string `json:"start_time" gorm:"type:varchar(50);not null"`
	End_time   string `json:"end_time" gorm:"type:varchar(50);not null"`
	E0         int    `json:"e0" gorm:"type:int;not null"`
	E1         int    `json:"e1" gorm:"type:int;not null"`
	E2         int    `json:"e2" gorm:"type:int;not null"`
	E3         int    `json:"e3" gorm:"type:int;not null"`
	E4         int    `json:"e4" gorm:"type:int;not null"`
	E5         int    `json:"e5" gorm:"type:int;not null"`
	E6         int    `json:"e6" gorm:"type:int;not null"`
	E8         int    `json:"e8" gorm:"type:int;not null"`
	EA         int    `json:"ea" gorm:"type:int;not null"`
	EE         int    `json:"ee" gorm:"type:int;not null"`
	F2         int    `json:"f2" gorm:"type:int;not null"`
	C0         int    `json:"c0" gorm:"type:int;not null"`
	C1         int    `json:"c1" gorm:"type:int;not null"`
	C2         int    `json:"c2" gorm:"type:int;not null"`
	C3         int    `json:"c3" gorm:"type:int;not null"`
	C4         int    `json:"c4" gorm:"type:int;not null"`
	C5         int    `json:"c5" gorm:"type:int;not null"`
	C6         int    `json:"c6" gorm:"type:int;not null"`
	C7         int    `json:"c7" gorm:"type:int;not null"`
	C8         int    `json:"c8" gorm:"type:int;not null"`
	EH         int    `json:"eH" gorm:"type:int;not null"`
	EF         int    `json:"eF" gorm:"type:int;not null"`
}
type FragPatternRecognition struct {
	DevID           string `grom:"varchar(20);index:applianceid_idx;"`
	StartTime       string `grom:"varchar(20);index:start_time_idx;"`
	EndTime         string `grom:"varchar(20);"`
	Comloadsegment0 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	Comloadsegment1 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	Comloadsegment2 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	Comloadsegment3 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	PWMLoadStat     int    `grom:"tinyint;"` //工况状态 -1最小 0 正常 1 最大
	WaterLoadStat   int    `grom:"tinyint;"` //产水状态
	TempScore       int    `json:"temp_score" gorm:"type:tinyint;not null"`
	HeatTempScore   int    `json:"heat_temp_score" gorm:"type:float;not null"`
	StableTempScore int    `json:"stable_temp_score" gorm:"type:float;not null"`
}
type Signalinput struct{ 
	DevType      string `json:"dev_type" gorm:"type:varchar(50);not null"`
	Applianceid      string `json:"applianceid" gorm:"type:varchar(50);"`
	Datatime    string `json:"datatime" gorm:"type:varchar(50);"`
	Ut string `json:"ut" gorm:"type:varchar(50);"`
	Temprange int `json:"temprange" gorm:"type:varchar(50);"`
	Outtemp   int `json:"outtemp" gorm:"type:varchar(50);"`
	Intemp   int `json:"intemp" gorm:"type:varchar(50);"`
    Flow     int `json:"flow" gorm:"type:varchar(50);"`
    TempPattern  string `json:"temp_pattern" gorm:"type:varchar(50);"`
}
type Stable_paramter struct{
	DevType      string `json:"dev_type" gorm:"type:varchar(50);not null"`
	Applianceid      string `json:"applianceid" gorm:"type:varchar(50);"`
	Start_time string `json:"start_time" gorm:"type:varchar(50);not null"`
	End_time   string `json:"end_time" gorm:"type:varchar(50);not null"`
	Ut int `json:"ut" gorm:"-"`
	Temprange int `json:"temprange" gorm:"-"`
	F float64 `json:"f" gorm:"-"`
	Comloadsegment0 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常	
  Ut0 int `json:"ut0" gorm:"type:varchar(50);"`
   Temprange0 int `json:"temprange0" gorm:"type:varchar(50);"`
	K0             float64 `json:"k0" gorm:"type:float;"`
	F0             float64 `json:"f0" gorm:"type:float;"`
	Comloadsegment1 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	Ut1 int `json:"ut1" gorm:"type:varchar(50);"`
	Temprange1 int `json:"temprange1" gorm:"type:varchar(50);"`
	F1             float64 `json:"f1" gorm:"type:float;"`
    K1              float64 `json:"k1" gorm:"type:float;"`
	Comloadsegment2 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	Ut2 int `json:"ut2" gorm:"type:varchar(50);"`
	Temprange2 int `json:"temprange2" gorm:"type:varchar(50);"`
	K2               float64 `json:"k2" gorm:"type:float;"`
	F2             float64 `json:"f2" gorm:"type:float;"`
	Comloadsegment3 int    `grom:"tinyint;"` //在该燃烧负荷段上出现异常
	Ut3 int `json:"ut3" gorm:"type:varchar(50);"`
	Temprange3 int `json:"temprange3" gorm:"type:varchar(50);"`
	K3                float64 `json:"k3" gorm:"type:float;"`
	F3             float64 `json:"f3" gorm:"type:float;"`
    Flag   int `json:"flag" gorm:"type:varchar(50);"`
}
type BasicData struct {
	firstSegment
	secondSegment
	thirdSegment
	fouthSegment
}
type firstSegment struct {
	basicVectorT []float64
	basicVectorK []float64
}
type secondSegment struct {
	basicVectorT []float64
	basicVectorK []float64
}
type thirdSegment struct {
	basicVectorT []float64
	basicVectorK []float64
}
type fouthSegment struct {
	basicVectorT []float64
	basicVectorK []float64
}
type  savegood struct {
	jd   float64
	ka    float64
	kb    float64
	kc    float64
}

//修改恒温参数表
type ReviseConstParaK struct{
	DevId             string	`json:"dev_id" gorm:"type:varchar(255)"`
	StartTime         string	`json:"start_time" gorm:"type:varchar(255)"`
	EndTime           string	`json:"end_time" gorm:"type:varchar(255)"`
	Comloadsegment0   int		`json:"comloadsegment0" gorm:"type:int"`
	Ka0				  float64   `json:"ka0" gorm:"type:float"`
	Kb0				  float64	`json:"kb0" gorm:"type:float"`
	Kc0               float64	`json:"kc0" gorm:"type:float"`
	Comloadsegment1   int	    `json:"comloadsegment1" gorm:"type:int"`
	Ka1               float64	`json:"ka1" gorm:"type:float"`
	Kb1               float64	`json:"kb1" gorm:"type:float"`
	Kc1               float64	`json:"kc1" gorm:"type:float"`
	Comloadsegment2   int	    `json:"comloadsegment2" gorm:"type:int"`
	Ka2               float64	`json:"ka2" gorm:"type:float"`
	Kb2               float64	`json:"kb2" gorm:"type:float"`
	Kc2               float64	`json:"kc2" gorm:"type:float"`
	Comloadsegment3   int	   `json:"comloadsegment3" gorm:"type:int"`
	Ka3               float64	`json:"ka3" gorm:"type:float"`
	Kb3               float64	`json:"kb3" gorm:"type:float"`
	Kc3               float64	`json:"kc3" gorm:"type:float"`
}
func (r ReviseConstParaK) TableName() string {
	return "revise_constparak"
}