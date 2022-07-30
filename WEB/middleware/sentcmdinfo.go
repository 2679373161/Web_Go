package middleware

import (
	"ginEssential/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*=========================================================
 * 函数名称： IdCheck
 * 功能描述:  验证id是否正确
 * 创建日期： 2021.12.25
 =========================================================*/

func IdCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		applianceid := ctx.DefaultQuery("applianceid", "188016488514318")
		if len(applianceid) != 15 {
			response.Success(ctx, gin.H{"errflag": "1"}, "设备号错误")
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

/*=========================================================
 * 函数名称： NoparasettingfirstCheck
 * 功能描述:  验证 非参数设置可调参数第一组 的逻辑
 * 创建日期： 2021.12.25
 =========================================================*/

func NoparasettingfirstCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//回水水流值
		ReWaterFlow := ctx.Query("rewaterflow")
		//风压传感器报警点补偿值
		WindPressureSensor := ctx.Query("windpressuresensor")
		_, err := strconv.ParseInt(ReWaterFlow, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(ReWaterFlow) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		_, err = strconv.ParseInt(WindPressureSensor, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(WindPressureSensor) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

/*=========================================================
 * 函数名称： NoparasettingsecondCheck
 * 功能描述:  验证 非参数设置可调参数第二组 的逻辑
 * 创建日期： 2021.12.25
 =========================================================*/

func NoparasettingsecondCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//最大负荷风机电流偏差系数MaxCurrCoeff
		MaxCurrCoeff := ctx.Query("maxcurrcoeff")
		//最小负荷风机电流偏差系数MinCurrCoeff
		MinCurrCoeff := ctx.Query("mincurrcoeff")
		//最大负荷风机占空比偏差系数MaxDutyCycCoeff
		MaxDutyCycCoeff := ctx.Query("maxdutycyccoeff")
		//最小负荷风机占空比偏差系数MinDutyCycCoeff
		MinDutyCycCoeff := ctx.Query("mindutycyccoeff")

		MaxCurrCoeffInt, err := strconv.ParseInt(MaxCurrCoeff, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(MaxCurrCoeff) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if MaxCurrCoeffInt < 80&&MaxCurrCoeffInt!=0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if MaxCurrCoeffInt > 120 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}

		MinCurrCoeffInt, err := strconv.ParseInt(MinCurrCoeff, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(MinCurrCoeff) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if MinCurrCoeffInt < 80&&MinCurrCoeffInt!=0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if MinCurrCoeffInt > 120 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}

		MaxDutyCycCoeffInt, err := strconv.ParseInt(MaxDutyCycCoeff, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(MaxDutyCycCoeff) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if MaxDutyCycCoeffInt < 80&&MaxDutyCycCoeffInt!=0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if MaxDutyCycCoeffInt > 120 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		MinDutyCycCoeffInt, err := strconv.ParseInt(MinDutyCycCoeff, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(MinDutyCycCoeff) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if MinDutyCycCoeffInt < 80&&MinDutyCycCoeffInt!=0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if MinDutyCycCoeffInt > 120 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

/*=========================================================
 * 函数名称： SingleparaCheck
 * 功能描述:  验证 单个改写参数设置参数值 的逻辑
 * 创建日期： 2021.12.25
 =========================================================*/

func SingleparaCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		index := ctx.Query("index") //改写序号
		value := ctx.Query("value") //改写值
		SingleMin := ctx.Query("singlemin")
		SingleMax := ctx.Query("singlemax")
		SingleMinInt, _ := strconv.ParseInt(SingleMin, 16, 64)
		SingleMaxInt, _ := strconv.ParseInt(SingleMax, 16, 64)

		_, err := strconv.ParseInt(index, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(index) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}

		ValueInt, _ := strconv.ParseInt(value, 16, 64)
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(value) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if ValueInt < SingleMinInt {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if ValueInt > SingleMaxInt {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

/*=========================================================
 * 函数名称： ParainformationCheck
 * 功能描述:  验证 恒温算法2.0相关数据设置默认参数信息 的逻辑
 * 创建日期： 2021.12.25
 =========================================================*/

func ParainformationCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		segment := ctx.DefaultQuery("segment", "00") //段序号
		ka := ctx.DefaultQuery("ka", "00")           //模型参数Ka*100
		kb := ctx.DefaultQuery("kb", "00")           //模型参数Kb*100
		kc := ctx.DefaultQuery("kc", "00")           //模型参数Kc*1000
		kf := ctx.DefaultQuery("kf", "00")           //bit0-6:模型参数Kf*10,bit7:正负(0：正1：负)
		T1a := ctx.DefaultQuery("T1a", "00")         //模型参数 T1a*1000
		T1c := ctx.DefaultQuery("T1c", "00")         //模型参数T1c*10
		T2a := ctx.DefaultQuery("T2a", "00")         //模型参数T2a*1000
		T2c := ctx.DefaultQuery("T2c", "00")         //模型参数T2c*100
		Tda := ctx.DefaultQuery("Tda", "00")         //模型参数Tda*1000
		Tdc := ctx.DefaultQuery("Tdc", "00")         //模型参数Tdc*10
		Wc := ctx.DefaultQuery("Wc", "00")           //Wc*10
		Wo := ctx.DefaultQuery("Wo", "00")           //Wo*10

		segmentInt, err := strconv.ParseInt(segment, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if segmentInt != 1 && segmentInt != 2 && segmentInt != 3 && segmentInt != 0 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入不存在的段序号")
			ctx.Abort()
			return
		}
		kaInt, err := strconv.ParseInt(ka, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(ka) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if kaInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		kbInt, err := strconv.ParseInt(kb, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(kb) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if kbInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		kcInt, err := strconv.ParseInt(kc, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(kc) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if kcInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		kfInt, err := strconv.ParseInt(kf, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(kf) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if kfInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		T1aInt, err := strconv.ParseInt(T1a, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(T1a) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if T1aInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		T1cInt, err := strconv.ParseInt(T1c, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(T1c) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if T1cInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		T2aInt, err := strconv.ParseInt(T2a, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(T2a) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if T2aInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		T2cInt, err := strconv.ParseInt(T2c, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(T2c) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if T2cInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		TdaInt, err := strconv.ParseInt(Tda, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(Tda) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if TdaInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		TdcInt, err := strconv.ParseInt(Tdc, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(Tdc) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if TdcInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		WcInt, err := strconv.ParseInt(Wc, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(Wc) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if WcInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		WoInt, err := strconv.ParseInt(Wo, 16, 64) //段序号
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(Wo) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if WoInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

/*=========================================================
 * 函数名称： ParaSettingCheck
 * 功能描述:  验证 多个改写参数设置的参数值 的逻辑
 * 创建日期： 2022.01.06
 =========================================================*/
func ParaSettingCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		FA := ctx.DefaultQuery("FA", "00")
		FF := ctx.DefaultQuery("FF", "00")
		PH := ctx.DefaultQuery("PH", "00")
		FH := ctx.DefaultQuery("FH", "00")
		PL := ctx.DefaultQuery("PL", "00")
		FL := ctx.DefaultQuery("FL", "00")
		dH := ctx.DefaultQuery("dH", "00")
		Fd := ctx.DefaultQuery("Fd", "00")
		CH := ctx.DefaultQuery("CH", "00")
		FC := ctx.DefaultQuery("FC", "00")
		CA := ctx.DefaultQuery("CA", "00")
		nE := ctx.DefaultQuery("nE", "00")
		FP := ctx.DefaultQuery("FP", "00")
		HS := ctx.DefaultQuery("HS", "00")
		Hb := ctx.DefaultQuery("Hb", "00")
		HE := ctx.DefaultQuery("HE", "00")
		HL := ctx.DefaultQuery("HL", "00")
		HU := ctx.DefaultQuery("HU", "00")
		Fn := ctx.DefaultQuery("Fn", "00")
		_, err := strconv.ParseInt(FA, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(FA) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		_, err = strconv.ParseInt(FF, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(FF) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		PHInt, err := strconv.ParseInt(PH, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(PH) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if PHInt < 10 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if PHInt > 200 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		FHInt, err := strconv.ParseInt(FH, 16, 64) //参数值
		if FHInt != 0 {
			if err != nil {
				response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
				ctx.Abort()
				return
			}
			if len(FH) > 2 {
				response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
				ctx.Abort()
				return
			}
			if FHInt < 10 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
				ctx.Abort()
				return
			} else if FHInt > 160 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
				ctx.Abort()
				return
			}
		}
		PLInt, err := strconv.ParseInt(PL, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(PL) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if PLInt < 10 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		if PLInt > 200 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		FLInt, err := strconv.ParseInt(FL, 16, 64) //参数值
		if FLInt != 0 {
			if err != nil {
				response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
				ctx.Abort()
				return
			}
			if len(FL) > 2 {
				response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
				ctx.Abort()
				return
			}
			if FLInt < 10 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
				ctx.Abort()
				return
			} else if FLInt > 160 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
				ctx.Abort()
				return
			}
		}
		dHInt, err := strconv.ParseInt(dH, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(dH) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if dHInt < 10 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if dHInt > 200 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		FdInt, err := strconv.ParseInt(Fd, 16, 64) //参数值
		if FdInt != 0 {
			if err != nil {
				response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
				ctx.Abort()
				return
			}
			if len(Fd) > 2 {
				response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
				ctx.Abort()
				return
			}
			if FdInt < 10 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
				ctx.Abort()
				return
			} else if FdInt > 160 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
				ctx.Abort()
				return
			}
		}
		CHInt, err := strconv.ParseInt(CH, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(CH) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if CHInt < 10 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if CHInt > 200 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		FCInt, err := strconv.ParseInt(FC, 16, 64) //参数值
		if FCInt != 0 {
			if err != nil {
				response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
				ctx.Abort()
				return
			}
			if len(FC) > 2 {
				response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
				ctx.Abort()
				return
			}
			if FCInt < 10 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
				ctx.Abort()
				return
			} else if FCInt > 160 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
				ctx.Abort()
				return
			}
		}
		CAInt, err := strconv.ParseInt(CA, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(CA) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if CAInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if CAInt > 5 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		nEInt, err := strconv.ParseInt(nE, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(nE) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if nEInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if nEInt > 1 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		FPInt, err := strconv.ParseInt(FP, 16, 64) //参数值
		if FPInt != 0 {
			if err != nil {
				response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
				ctx.Abort()
				return
			}
			if len(FP) > 2 {
				response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
				ctx.Abort()
				return
			}
			if FPInt < 4 {
				response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
				ctx.Abort()
				return
			}
		}
		HSInt, err := strconv.ParseInt(HS, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(HS) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if HSInt < 50 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		}
		if HSInt > 100 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		HbInt, err := strconv.ParseInt(Hb, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(Hb) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if HbInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if HbInt > 5 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		HEInt, err := strconv.ParseInt(HE, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(HE) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if HEInt < 20 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if HEInt > 50 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		HLInt, err := strconv.ParseInt(HL, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(HL) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if HLInt < 40 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if HLInt > 60 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		HUInt, err := strconv.ParseInt(HU, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(HU) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if HUInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if HUInt > 99 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		FnInt, err := strconv.ParseInt(Fn, 16, 64) //参数值
		if err != nil {
			response.Success(ctx, gin.H{"errflag": "2"}, "参数格式错误")
			ctx.Abort()
			return
		}
		if len(Fn) > 2 {
			response.Success(ctx, gin.H{"errflag": "2"}, "输入参数超过范围")
			ctx.Abort()
			return
		}
		if FnInt < 0 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出下限")
			ctx.Abort()
			return
		} else if FnInt > 2 {
			response.Success(ctx, gin.H{"errflag": "4"}, "输入参数超出上限")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
