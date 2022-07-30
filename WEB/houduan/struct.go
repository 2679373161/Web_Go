package houduan

//输入数据结构体
type user struct {
	Applianceid string `json:"applianceid" gorm:"type:varchar(50);not null"`
	Datatime    string `json:"datetime" gorm:"type:varchar(50);not null"`
	Flame       int    `json:"flame" gorm:"type:int;not null"`
	Flow        int    `json:"flow" gorm:"type:int;not null"`
	Outtemp     int    `json:"out_temp" gorm:"type:int;not null"`
	Settemp     int    `json:"set_temp" gorm:"type:int;not null"`
	Model       int    `json:"model" gorm:"type:int;not null"`
	ZoneID      int    `json:"zone_id" gorm:"type:int;not null"`
	TempPattern int    `json:"temp_pattern" gorm:"type:int;not null"`
	TabName     string `gorm:"-"`
}

// Multimodality 设备指标
type Multimodality struct {
	ProvinceCode       string  `json:"province_code" gorm:"type:varchar(6);not null"`
	CityCode           string  `json:"city_code" gorm:"type:varchar(6);not null"`
	DevType            string  `json:"dev_type" gorm:"type:varchar(10);not null"`
	DevId              int     `json:"dev_id" gorm:"type:varchar(20);not null"`
	TimeDate           string  `json:"total_time" gorm:"type:varchar(20);not null"`
	ValidTime          string  `json:"valid_time" gorm:"type:varchar(20);not null"`
	PatternNum         int     `json:"pattern_num" gorm:"type:int;not null"`
	AverageTime        string  `json:"average_time" gorm:"type:varchar(10);not null"`
	UnStableProportion float64 `json:"fluctuation_proportion" gorm:"type:float;not null"`
	StableTime         string  `json:"stable_time" gorm:"type:varchar(20);not null"`
	UnStableTime       string  `json:"fluctuation_time" gorm:"type:varchar(20);not null"`
	MaximumTime        string  `json:"maximum_time" gorm:"type:varchar(10);not null"`
	MinimumTime        string  `json:"minimum_time" gorm:"type:varchar(10);not null"`
	UnStableBehavior   int     `json:"un_stable_behavior" gorm:"type:int;not null"`
	AbnormalFlag       int     `json:"abnormal_flag" gorm:"type:int;not null"`
	FlowMultipleScore  int     `json:"flow_multiple_score" gorm:"type:int;not null"`
	Three              int     `json:"three" gorm:"type:int;not null"`
	Five               int     `json:"five" gorm:"type:int;not null"`
	Ten                int     `json:"ten" gorm:"type:int;not null"`
	TabName            string  `gorm:"-"`
}

// OutTemperature 出水温度指标输出
type OutTemperature struct {
	DevId                    string `json:"dev_id" gorm:"type:varchar(50);not null"`
	LocCode                  string `json:"loc_code" gorm:"type:varchar(50);not null"`
	TimeDate                 string `json:"time_date" gorm:"type:varchar(50);not null"`
	StableDuration           string `json:"stable_duration" gorm:"type:varchar(50);not null"`
	HeatingUpDuration        string `json:"heating_up_duration" gorm:"type:varchar(50);not null"`
	OverTemperatureDuration  string `json:"over_temperature_duration" gorm:"type:varchar(50);not null"`
	ShutdownDuration         string `json:"shutdown_duration" gorm:"type:varchar(50);not null"`
	UnderTemperatureDuration string `json:"under_temperature_duration" gorm:"type:varchar(50);not null"`
	Index                    int    `json:"index" gorm:"type:varchar(50);not null"`
	TabName                  string `gorm:"-"`
}

// ModeFragment 模式片段指标输出
type ModeFragment struct {
	DevId                string  `json:"dev_id" gorm:"type:varchar(20);not null"`
	Pattern              int     `json:"pattern" gorm:"type:int;not null"`
	StartTime            string  `json:"start_time" gorm:"type:varchar(20);not null"`
	EndTime              string  `json:"end_time" gorm:"type:varchar(20);not null"`
	DurationTime         string  `json:"duration_time" gorm:"type:varchar(20);not null"`
	Extreme              int     `json:"extreme" gorm:"type:int;not null"`
	MaxChange            float64 `json:"max_change" gorm:"type:float;not null"`
	Average              float64 `json:"average" gorm:"type:float;not null"`
	Deviation            float64 `json:"deviation" gorm:"type:float;not null"`
	UpNumber             int     `json:"up_number" gorm:"type:int;not null"`
	DownNumber           int     `json:"down_number" gorm:"type:int;not null"`
	MultipleScore        int     `json:"multiple_score" gorm:"int;not null"`
	HeatDuration         string  `json:"heat-duration" gorm:"type:varchar(20);not null"`
	UnStableTempDuration string  `json:"un_stable_temp_duration" gorm:"type:varchar(20);not null"`
	UnStableTempPercent  float64 `json:"un_stable_temp_percent" gorm:"type:float;not null"`
	UnHeatDev            float64 `json:"un_heat_dev" gorm:"type:float;not null"`
	HeatDev              float64 `json:"heat_dev" gorm:"type:float;not null"`
	TempNum              int     `json:"temp_num" gorm:"type:int;not null"`
	TabName              string  `gorm:"-"`
}

// DaysSummary 设备日指标汇总
type DaysSummary struct {
	ProvinceCode       string  `json:"province_code" gorm:"type:varchar(6);not null"`
	CityCode           string  `json:"city_code" gorm:"type:varchar(6);not null"`
	DevType            string  `json:"dev_type" gorm:"type:varchar(10);not null"`
	DevId              string  `json:"dev_id" gorm:"type:varchar(20);not null"`
	TimeDate           string  `json:"total_time" gorm:"type:varchar(20);not null"`
	ValidTime          string  `json:"valid_time" gorm:"type:varchar(20);not null"`
	PatternNum         int     `json:"pattern_num" gorm:"type:int;not null"`
	AverageTime        string  `json:"average_time" gorm:"type:varchar(10);not null"`
	UnStableProportion float64 `json:"fluctuation_proportion" gorm:"type:float;not null"`
	StableTime         string  `json:"stable_time" gorm:"type:varchar(20);not null"`
	UnStableTime       string  `json:"fluctuation_time" gorm:"type:varchar(20);not null"`
	MaximumTime        string  `json:"maximum_time" gorm:"type:varchar(10);not null"`
	MinimumTime        string  `json:"minimum_time" gorm:"type:varchar(10);not null"`
	UnStableBehavior   int     `json:"un_stable_behavior" gorm:"type:int;not null"`
	AbnormalFlag       int     `json:"abnormal_flag" gorm:"type:int;not null"`
	FlowMultipleScore  int     `json:"flow_multiple_score" gorm:"type:int;not null"`
	Three              int     `json:"three" gorm:"type:int;not null"`
	Five               int     `json:"five" gorm:"type:int;not null"`
	Ten                int     `json:"ten" gorm:"type:int;not null"`
	TabName            string  `gorm:"-"`


	
}
type MonthSummary struct {
	ProvinceCode       string  `json:"province_code" gorm:"type:varchar(6);not null"`
	CityCode           string  `json:"city_code" gorm:"type:varchar(6);not null"`
	DevType            string  `json:"dev_type" gorm:"type:varchar(10);not null"`
	DevId              string  `json:"dev_id" gorm:"type:varchar(20);not null"`
	TimeDate           string  `json:"total_time" gorm:"type:varchar(20);not null"`
	ValidTime          string  `json:"valid_time" gorm:"type:varchar(20);not null"`
	PatternNum         int     `json:"pattern_num" gorm:"type:int;not null"`
	AverageTime        string  `json:"average_time" gorm:"type:varchar(10);not null"`
	UnStableProportion float64 `json:"fluctuation_proportion" gorm:"type:float;not null"`
	StableTime         string  `json:"stable_time" gorm:"type:varchar(20);not null"`
	UnStableTime       string  `json:"fluctuation_time" gorm:"type:varchar(20);not null"`
	MaximumTime        string  `json:"maximum_time" gorm:"type:varchar(10);not null"`
	MinimumTime        string  `json:"minimum_time" gorm:"type:varchar(10);not null"`
	UnStableBehavior   int     `json:"un_stable_behavior" gorm:"type:int;not null"`
	AbnormalFlag       int     `json:"abnormal_flag" gorm:"type:int;not null"`
	FlowMultipleScore  int     `json:"flow_multiple_score" gorm:"type:int;not null"`
	Three              int     `json:"three" gorm:"type:int;not null"`
	Five               int     `json:"five" gorm:"type:int;not null"`
	Ten                int     `json:"ten" gorm:"type:int;not null"`
	TabName            string  `gorm:"-"`
}

// MultimodalityCity 城市汇总指标
type MultimodalityCity struct {
	ProvinceCode       string `json:"province_code" gorm:"type:varchar(50);not null"`
	CityCode           string `json:"city_code" gorm:"type:varchar(50);not null"`
	EquipmentNum       int    `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate           string `json:"time_date" gorm:"type:varchar(50);not null"`
	ValidTime          string `json:"valid_time" gorm:"type:varchar(50);not null"`
	PatternNum         int    `json:"pattern_num" gorm:"type:varchar(50);not null"`
	AverageTime        string `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion string `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	StableTime         string `json:"stable_time" gorm:"type:varchar(50);not null"`
	UnStableTime       string `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime        string `json:"maximum_time" gorm:"type:varchar(50);not null"`
	MinimumTime        string `json:"minimum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior   int    `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	AbnormalFlag       int    `json:"abnormal_flag" gorm:"type:varchar(50);not null"`
	FlowMultipleScore  int    `json:"flow_multiple_score" gorm:"type:int(50);not null"`
	TabName            string `gorm:"-"`
}

// MultimodalityProvince 省汇总指标
type MultimodalityProvince struct {
	ProvinceCode       string `json:"province_code" gorm:"type:varchar(50);not null"`
	EquipmentNum       int    `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate           string `json:"time_date" gorm:"type:varchar(50);not null"`
	ValidTime          string `json:"valid_time" gorm:"type:varchar(50);not null"`
	PatternNum         int    `json:"pattern_num" gorm:"type:varchar(50);not null"`
	AverageTime        string `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion string `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	StableTime         string `json:"stable_time" gorm:"type:varchar(50);not null"`
	UnStableTime       string `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime        string `json:"maximum_time" gorm:"type:varchar(50);not null"`
	MinimumTime        string `json:"minimum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior   int    `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	AbnormalFlag       int    `json:"abnormal_flag" gorm:"type:varchar(50);not null"`
	FlowMultipleScore  int    `json:"flow_multiple_score" gorm:"type:int(50);not null"`
	Three              int    `json:"three" gorm:"type:int(50);not null"`
	Five               int    `json:"five" gorm:"type:int(50);not null"`
	Ten                int    `json:"ten" gorm:"type:int(50);not null"`
	TabName            string `gorm:"-"`
}

// MultimodalityType 型号汇总指标
type MultimodalityType struct {
	DevType            string `json:"dev_type" gorm:"type:varchar(50);not null"`
	EquipmentNum       int    `json:"equipment_num" gorm:"type:varchar(50);not null"`
	TimeDate           string `json:"time_date" gorm:"type:varchar(50);not null"`
	ValidTime          string `json:"valid_time" gorm:"type:varchar(50);not null"`
	PatternNum         int    `json:"pattern_num" gorm:"type:varchar(50);not null"`
	AverageTime        string `json:"average_time" gorm:"type:varchar(50);not null"`
	UnStableProportion string `json:"un_stable_proportion" gorm:"type:varchar(50);not null"`
	StableTime         string `json:"stable_time" gorm:"type:varchar(50);not null"`
	UnStableTime       string `json:"un_stable_time" gorm:"type:varchar(50);not null"`
	MaximumTime        string `json:"maximum_time" gorm:"type:varchar(50);not null"`
	MinimumTime        string `json:"minimum_time" gorm:"type:varchar(50);not null"`
	UnStableBehavior   int    `json:"un_stable_behavior" gorm:"type:varchar(50);not null"`
	AbnormalFlag       int    `json:"abnormal_flag" gorm:"type:varchar(50);not null"`
	FlowMultipleScore  int    `json:"flow_multiple_score" gorm:"type:int(50);not null"`
	TabName            string `gorm:"-"`
}

// Equipment 设备id表
type Equipment struct {
	Identifier  int    `json:"identifier" gorm:"type:varchar(20);not null"`
	DevId       int    `json:"dev_id" gorm:"type:varchar(20);not null"`
	DevType     string `json:"dev_type" gorm:"type:varchar(20);not null"`
	CityCode    string `json:"city_code" gorm:"type:varchar(10);not null"`
	InsertTime  string `json:"insert_time" gorm:"type:varchar(20);not null"`
	RefreshTime string `json:"refresh_time" gorm:"type:varchar(20);not null"`
}

// ModeFragment 新的模式片段指标输出
type ModeFragment1 struct {
	DevId        string  `json:"dev_id" gorm:"type:varchar(50);not null;index:dev_id_idx"`
	StartTime    string  `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx"`
	EndTime      string  `json:"end_time" gorm:"type:varchar(20);not null"`
	DurationTime string  `json:"duration_time" gorm:"type:varchar(20);not null"`
	WaterPattern int     `json:"water_pattern" gorm:"type:int;not null"`
	FlowAvg      float64 `json:"flow_avg" gorm:"type:float;not null"`
	SmallWater   float64 `json:"small_water" gorm:"type:float;not null"`
	//Extreme              int     `json:"extreme" gorm:"type:int;not null"`
	//MaxChange            float64 `json:"max_change" gorm:"type:float;not null"`
	//Average              float64 `json:"average" gorm:"type:float;not null"`
	Deviation            float64 `json:"deviation" gorm:"type:float;not null"`
	UpNumber             int     `json:"up_number" gorm:"type:int;not null"`
	DownNumber           int     `json:"down_number" gorm:"type:int;not null"`
	WaterScore           int     `json:"water_score" gorm:"int;not null"`
	HeatDuration         string  `json:"heat-duration" gorm:"type:varchar(20);not null"`
	UnStableTempDuration string  `json:"un_stable_temp_duration" gorm:"type:varchar(20);not null"`
	UnStableTempPercent  float64 `json:"un_stable_temp_percent" gorm:"type:float;not null"`
	UnHeatDev            float64 `json:"un_heat_dev" gorm:"type:float;not null"`
	TempPattern          int     `json:"temp_pattern" gorm:"type:int;not null"`
	OvershootValue       int     `json:"overshoot_value" gorm:"type:int;"`
	StateAccuracy        int     `json:"state_accuracy" gorm:"type:int;"`
	TempScore            int     `json:"temp_score" gorm:"type:int;not null"`
	NewTempScore         int     `json:"new_temp_score" gorm:"type:int;not null"`
	HeatTempScore        int     `json:"heat_temp_score" gorm:"type:int;not null"`
	StableTempScore      int     `json:"stable_temp_score" gorm:"type:int;not null"`
	TempJudgeFlag        bool    `json:"temp_judge_flag" gorm:"type:bool;not null"`
	WaterFlag            int     `json:"water_flag" gorm:"type:int;not null"`
	TempFlag             int     `json:"temp_flag" gorm:"type:int;not null"`
	AbnormalState        int     `json:"abnormal_state" gorm:"type:int;not null"`
	TabName              string  `gorm:"-"`
}

//实验设备id表
//type InstantEquipment struct {
//	Applianceid     string `json:"applianceid" gorm:"type:varchar(50);not null"`
//	StartTime        string `json:"start_time" gorm:"type:varchar(20);not null;index:start_time_idx"`
//	EndTime          string `json:"end_time" gorm:"type:varchar(20);not null;"`
//	Move_task_flag   int    `json:"move_task_flag" gorm:"type:tinyint;not null;default:0"`
//	Mining_task_flag int    `json:"mining_task_flag" gorm:"type:tinyint;not null;default:0"`
//}






func NewUser1(tableName string) *Multimodality {
	return &Multimodality{TabName: tableName}
}

func (u Multimodality) TableName() string {
	return u.TabName
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

// ProvinceCodeName 省份名
func ProvinceCodeName(provinceName string) string {
	value := "0000"
	provinceName = provinceName[0:2]
	finalValue := provinceName + value
	return finalValue
}
func ProvinceCodeName1(provinceName string) string {
	value := "0000"
	provinceName = provinceName[0:2]
	finalValue := provinceName + value
	return finalValue
}

//省编号
func ProvinceCode(provinceName string) string {
	value := "0000"
	provinceName = provinceName[0:2]
	finalValue := provinceName + value
	return finalValue
}

//省份编码
var provinceNumber = [...]string{"130000", "140000", "440000"}

//城市编码
var cityNumber = [...]string{"130100", "130200", "130300", "130400", "130500", "130600", "130700", "130800", "130900", "131000", "131100",
	"140100", "140200", "140300", "140400", "140500", "140600", "140700", "140800", "140900", "141000", "141100", "440100", "440300", "440400",
	"440500", "440600", "440700", "440800", "440900", "441200", "441300", "441400", "441500", "441600", "441700", "441800", "441900", "442000",
	"445100", "445200", "445300"}

//型号编码
var typeNumber = [...]string{"00012HES", "00016HGS", "0012WH6A", "0012WS5A", "0012WS5B", "51012HES", "511000L6", "511000Y5", "511000Y9", "51100DL3",
	"51100GT9", "51100HS4", "51100GTS", "51100HSN", "51100HT5", "51100HXN", "51100JM1", "51100JM3", "51100TD1", "51100WD7", "51100Y8S", "51100ZC3",
	"51100HT7", "51100HWA"}
