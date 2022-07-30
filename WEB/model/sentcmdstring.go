package model

var SingleParaIndex = map[string]string{
	"01": "FA",
	"02": "FF",
	"03": "PH",
	"04": "FH",
	"05": "PL",
	"06": "FL",
	"07": "dH",
	"08": "Fd",
	"09": "CH",
	"0a": "FC",
	"0b": "Fn",
	"0c": "cA",
	"0d": "nE",
	"0e": "FP",
	"0f": "HS",
	"10": "Hb",
	"11": "HE",
	"12": "HL",
	"13": "HU",
	"14": "nH",
}

// NonParaNmae 非参数设置可调参数 名字
var NonParaNmae = []string{
	"maximum_load_fan_current_deviation_coefficient",//最大负荷风机电流偏差系数
	"minimum_load_fan_current_deviation_coefficien",//最小负荷风机电流偏差系数
	"maximum_load_fan_duty_cycle_deviation_coefficien",//最大负荷风机占空比偏差系数
	"minimum_load_fan_duty_cycle_deviation_coefficient",//最小负荷风机占空比偏差系数
	"backwater_flow_value",//回水水流值
	"frequency_compensation_value_of_wind_pressure_sensor_alarm_point",//风压传感器报警点频率补偿值
}

var FlyNonParaNmae = []string{
	"backwater_flow_value",
	"frequency_compensation_value_of_wind_pressure_sensor_alarm_point",
}

var ParaSetting = []string{
     "FA",
     "FF",
     "PH",
     "FH",
     "PL",
     "FL",
     "dH",
     "Fd",
     "CH",
     "FC",
     "nE",
     "cA",
     "FP",
     "LF",
     "HS",
     "Hb",
     "HE",
     "HL",
     "HU",
     "UA",
     "Ub",
     "Fn",

}


var SetCode = []string{
	"000013", //0 PH0
	"000014", //1 FH0
	"000015", //2 PL0
	"000016", //3 FL0
	"000017", //4 dH0
	"000018", //5 Fd0
	"000019", //6 CH0
	"000020", //7 FC0
}

var Setpar = []string{"PH0", "FH0", "PL0", "FL0", "dH0", "Fd0", "CH0", "FC0"}

// Paranames 恒温算法参数 名字
var Paranames = []string{
	"Ka0", "Kb0", "Kc0", "Kf0", "T1a0", "T1c0", "T2a0", "T2c0", "Tda0", "Tdc0", "Wc0", "Wo0",
	"Ka1", "Kb1", "Kc1", "Kf1", "T1a1", "T1c1", "T2a1", "T2c1", "Tda1", "Tdc1", "Wc1", "Wo1",
	"Ka2", "Kb2", "Kc2", "Kf2", "T1a2", "T1c2", "T2a2", "T2c2", "Tda2", "Tdc2", "Wc2", "Wo2",
	"Ka3", "Kb3", "Kc3", "Kf3", "T1a3", "T1c3", "T2a3", "T2c3", "Tda3", "Tdc3", "Wc3", "Wo3",
}

var Setpar2 = []string{"FA", "FF", "PH", "FH", "PL", "FL", "dH", "Fd","CH","FC","nE","CA","FP",
	"LF","HS","Hb","HE","HL","HU","UA","Ub","Fn"}