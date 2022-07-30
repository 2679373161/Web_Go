package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"
	"os"
	"path"
	"strings"

	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"strconv"

	"time"

	"reflect"

	"errors"

	"encoding/csv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ITimeMonitorController interface {
	RestController
	PageList(ctx *gin.Context)
	EquipmentIndexMonitoring(ctx *gin.Context)
	EquipmentHistoryMonitoring(ctx *gin.Context)
	EquipmentRewriteFlag(ctx *gin.Context)
	MultipleEquipmentInfo(ctx *gin.Context)
	EquipmentRewriteFlagTempClear(ctx *gin.Context)
	EquipmentRewriteFlagTempModify(ctx *gin.Context)
	GetUpdata(ctx *gin.Context)
	EquipmentInfoCreat(ctx *gin.Context)
	EquipmentInfoDelete(ctx *gin.Context)
	MultipleEquipmentInfoSingel(ctx *gin.Context)
	GetDelete(ctx *gin.Context)
}

type TimeMonitorController struct {
	DB *gorm.DB
}

func (t TimeMonitorController) Create(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestTableStoreDate); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	//创建tablestoredate
	tableStoreDate := model.TableDate{
		//CategoryId: requestTableStoreDate.CategoryId,
		Datatime: requestTableStoreDate.Datatime,
		Flame:    requestTableStoreDate.Flame,
		Outtemp:  requestTableStoreDate.Outtemp,
		Settemp:  requestTableStoreDate.Settemp,
		//Flow: requestTableStoreDate.Flow,
		//Model: requestTableStoreDate.Model,
	}
	if err := t.DB.Create(&tableStoreDate).Error; err != nil {
		panic(err)
		return
	}
	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "创建成功")

}

func (t TimeMonitorController) Update(ctx *gin.Context) {
	var requestTableStoreDate vo.CreateTableStoreRequest
	//数据验证
	if err := ctx.ShouldBind(&requestTableStoreDate); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	//获取path中的id
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	//更新文章
	if err := t.DB.Model(&tableStoreDate).Update(requestTableStoreDate).Error; err != nil {
		//panic(err)
		fmt.Println(err)
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "更新成功")

}

func (t TimeMonitorController) Show(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	//通过preload加载外键
	if t.DB.Preload("Category").Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"tableStoreDate": tableStoreDate}, "读取成功")
}

func (t TimeMonitorController) Delete(ctx *gin.Context) {
	tableStoreDateId := ctx.Params.ByName("id") //从上下文中解析

	var tableStoreDate model.TableDate
	if t.DB.Where("id=?", tableStoreDateId).First(&tableStoreDate).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	t.DB.Delete(&tableStoreDate)

	response.Fail(ctx, gin.H{"tableStoreDate": tableStoreDate}, "删除成功")

}

/*=========================================================
* 函数名称： EquipmentRewriteFlagTempModify
* 功能描述: 多参数表内的零时标志位修改
 para: 	num_falg   			"pageall" or "select_all" or "single"
      	id          		equipment_ID for row and "single"
       	fun        			"rewriteflag_temp" or "para_out"
        rewriteflag_temp    "1" or "0"
        value_para 			"para_num"
        index      			"indexvalue"
=========================================================*/
func (t TimeMonitorController) EquipmentRewriteFlagTempModify(ctx *gin.Context) {
	num_falg := ctx.DefaultQuery("num_falg", "0")
	rewriteflag_temp := ctx.DefaultQuery("rewriteflag_temp", "10")

	fun := ctx.DefaultQuery("fun", "0")
	value_para := ctx.DefaultQuery("value_para", "0")
	index := ctx.DefaultQuery("index", "0")
	var tableselect []model.Multiple_Equipment
	if fun == "rewriteflag_temp" {
		if num_falg == "single" {
			Appliance_id := ctx.DefaultQuery("applianceid", "0")
			common.DB.Table("multiple_para_rewrite").Model(&tableselect).Where("appliance_id=?", Appliance_id).Update("rewriteflag_temp", rewriteflag_temp)
		} else if num_falg == "pageall" {
			Appliance_id := ctx.QueryArray("applianceid[]")
			common.DB.Table("multiple_para_rewrite").Model(&tableselect).Where("appliance_id in (?)", Appliance_id).Update("rewriteflag_temp", rewriteflag_temp)
			response.Success(ctx, gin.H{"page_change_flag": "1"}, "成功")
			return
		} else if num_falg == "select_all" {
			common.DB.Table("multiple_para_rewrite").Where("rewriteflag_temp != ?", rewriteflag_temp).
				Updates(map[string]interface{}{
					"rewriteflag_temp": rewriteflag_temp,
				})
			response.Success(ctx, gin.H{"all_change_flag": "1"}, "成功")
			return
		}
	} else if fun == "para_out" {

		common.DB.Table("multiple_para_rewrite").Where("rewriteflag_temp = ?", "1").
			Updates(map[string]interface{}{
				"handleflag":  "0",
				"succeedflag": "0",
				"updatetime":  time.Now().Format("2006-01-02 15:04:05"),
			})
		common.DB.Table("multiple_para_rewrite").Where("rewriteflag_temp = ?", "0").
			Updates(map[string]interface{}{
				"handleflag":  "0",
				"succeedflag": "0",
			})
		common.DB.Exec("update multiple_para_rewrite set rewriteflag=rewriteflag_temp")
		response.Success(ctx, gin.H{"para_out_flag": "1"}, "成功")
		return
	}
	response.Success(ctx, gin.H{"outcmd": value_para, "bb": index, "para_out_flag": "0", "page_change_flag": "0"}, "成功")
}
func (t TimeMonitorController) EquipmentRewriteFlagTempClear(ctx *gin.Context) {
	rewriteflag_temp := ctx.DefaultQuery("rewriteflag_temp", "10")
	Appliance_id := ctx.DefaultQuery("applianceid", "0")

	outcmd := UpdateMultiplePara(rewriteflag_temp, Appliance_id)
	response.Success(ctx, gin.H{"outcmd": outcmd}, "成功")
}

/*=========================================================
 * 函数名称： UpdateFindZone
 * 功能描述: 多参数表内的零时标志位清空
 =========================================================*/
func UpdateMultiplePara(rewriteflag_temp string, Appliance_id string) int64 {
	db := common.GetDB()
	var k int64
	if Appliance_id != "0" {
		if rewriteflag_temp == "0" {
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", Appliance_id).
				Updates(map[string]interface{}{
					"rewriteflag_temp": "0",
				})
			fmt.Println("更新成功")
			k = 1
		} else if rewriteflag_temp == "1" {
			db.Table("multiple_para_rewrite").Where("appliance_id = ?", Appliance_id).
				Updates(map[string]interface{}{
					"rewriteflag_temp": "1",
				})
			fmt.Println("更新成功")
			k = 1
		} else {
			k = 0
		}
	} else {
		if rewriteflag_temp == "0" {
			//保证临时状态与初始状态一致
			common.DB.Exec("update multiple_para_rewrite set rewriteflag_temp=rewriteflag")
			fmt.Println("更新成功")
			k = 1
		} else if rewriteflag_temp == "1" {
			db.Table("multiple_para_rewrite").Where("rewriteflag_temp = ? ", "0").
				Updates(map[string]interface{}{
					"rewriteflag_temp": "1",
				})
			fmt.Println("更新成功")
			k = 1
		} else {
			k = 0
		}
	}

	return k
}
func (t TimeMonitorController) EquipmentRewriteFlag(ctx *gin.Context) {

	Appliance_id := ctx.DefaultQuery("applianceid", "188016488514318")

	flag := ctx.DefaultQuery("flag", "0")

	fmt.Println(Appliance_id)

	var Parameter_final []model.ParameterOutput

	common.DB.Table("parameter_final_settings").Where(" appliance_id = ? ", Appliance_id).Find(&Parameter_final)

	temp_num := 0
	for i := 0; i < len(Parameter_final); i++ {
		if flag == "0" {
			if Parameter_final[i].RewriteSuccessFlag == 1 {
				response.Success(ctx, gin.H{"RewriteSuccessFlag": 1}, "成功")

				return
			}
		} else if flag == "1" {

			if Parameter_final[i].RewriteSuccessFlag == 1 {
				temp_num++
				if temp_num >= 62 {
					response.Success(ctx, gin.H{"RewriteSuccessFlag": 2}, "成功")

					return
				}
			}
		}

	}
	response.Success(ctx, gin.H{"RewriteSuccessFlag": 0}, "成功")
}
func (t TimeMonitorController) MultipleEquipmentInfo(ctx *gin.Context) {
	var tableData_multiple_equipment []model.Multiple_Equipment
	var tableData_multiple_para_line []model.ParamenBatch

	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	loading_flag := ctx.DefaultQuery("loading_flag", "-1")
	refresh_flage := ctx.DefaultQuery("refresh_flage", "-1")

	var count_equipment int

	var count_equipment_remain int

	count_equipment_remain = 0
	//state_flag =1 busy //等待修改 001
	//state_flag =0 normal //空闲设备 000
	//state_flag =3 unsuccess //失败设备 101
	//state_flag =4 success //成功设备 111
	//state_flag =5 abnormal //状态异常设备 else 010 011 110 100
	var count int
	var count_success_e int
	count_success_e = 0
	var count_fail_e int
	count_fail_e = 0

	var count_all_rewrite int
	var count_all_rewrite_ok int
	var count_all_rewrite_sucess int
	common.DB.Table("multiple_para_rewrite").Where("handleflag=?", "1").Scan(&tableData_multiple_equipment).Count(&count_all_rewrite_ok)
	common.DB.Table("multiple_para_rewrite").Where("succeedflag=?", "1").Scan(&tableData_multiple_equipment).Count(&count_all_rewrite_sucess)
	common.DB.Table("multiple_para_rewrite").Where("rewriteflag=?", "1").Scan(&tableData_multiple_equipment).Count(&count_all_rewrite)

	if loading_flag != "-1" {

		common.DB.Table("multiple_para_rewrite").Scan(&tableData_multiple_equipment)

	} else if refresh_flage == "0" || refresh_flage == "-1" {

		common.DB.Table("multiple_para_rewrite").Scan(&tableData_multiple_equipment).Count(&count)
		common.DB.Table("multiple_para_rewrite").Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&tableData_multiple_equipment)

	} else if refresh_flage == "1" {
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "0", "0", "0").Scan(&tableData_multiple_equipment).Count(&count)
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "0", "0", "0").Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&tableData_multiple_equipment)
	} else if refresh_flage == "2" {
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "0", "0", "1").Scan(&tableData_multiple_equipment).Count(&count)
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "0", "0", "1").Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&tableData_multiple_equipment)
	} else if refresh_flage == "3" {
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "1", "0", "1").Scan(&tableData_multiple_equipment).Count(&count)
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "1", "0", "1").Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&tableData_multiple_equipment)
	} else if refresh_flage == "4" {
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "1", "1", "1").Scan(&tableData_multiple_equipment).Count(&count)
		common.DB.Table("multiple_para_rewrite").Where("handleflag=? AND succeedflag=? AND rewriteflag=?", "1", "1", "1").Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&tableData_multiple_equipment)
	} else if refresh_flage == "5" {
		common.DB.Table("multiple_para_rewrite").Where("(handleflag != rewriteflag AND succeedflag=?) OR(handleflag=? AND succeedflag=? AND rewriteflag=?) OR(handleflag=? AND succeedflag=? AND rewriteflag=?)", "1", "0", "1", "0", "1", "0", "0").Scan(&tableData_multiple_equipment).Count(&count)
		common.DB.Table("multiple_para_rewrite").Where("(handleflag != rewriteflag AND succeedflag=?) OR(handleflag=? AND succeedflag=? AND rewriteflag=?) OR(handleflag=? AND succeedflag=? AND rewriteflag=?)", "1", "0", "1", "0", "1", "0", "0").Offset(perPage * (currentPage - 1)).Limit(perPage).Scan(&tableData_multiple_equipment)
	}

	for temp_count, tableDate := range tableData_multiple_equipment {
		if refresh_flage == "-1" {
			if tableDate.Rewriteflag == "1" {
				count_equipment = count_equipment + 1
				tableData_multiple_equipment[temp_count].Isrewirted = "true"
			} else {
				tableData_multiple_equipment[temp_count].Isrewirted = "false"
			}
		} else {
			if tableDate.Rewriteflag == "1" {
				count_equipment = count_equipment + 1
			}
			if tableDate.Rewriteflag_temp == "1" {
				tableData_multiple_equipment[temp_count].Isrewirted = "true"
			} else {
				tableData_multiple_equipment[temp_count].Isrewirted = "false"
			}
		}
		if (tableDate.Handleflag == "0") && (tableDate.Rewriteflag == "1") && (tableDate.Succeedflag == "0") {
			count_equipment_remain = count_equipment_remain + 1
			tableData_multiple_equipment[temp_count].State_flag = "1"
			tableData_multiple_equipment[temp_count].Equipment_State = "等待修改"
		} else {
			if tableDate.Handleflag == "1" && tableDate.Succeedflag == "1" {
				if tableDate.Rewriteflag == "1" {
					tableData_multiple_equipment[temp_count].State_flag = "4"
					tableData_multiple_equipment[temp_count].Equipment_State = "修改成功"
					count_success_e = count_success_e + 1
				} else {
					tableData_multiple_equipment[temp_count].State_flag = "5"
					tableData_multiple_equipment[temp_count].Equipment_State = "状态异常"
				}
			} else {
				if tableDate.Rewriteflag == "1" {
					if tableDate.Succeedflag == "1" {
						tableData_multiple_equipment[temp_count].State_flag = "5"
						tableData_multiple_equipment[temp_count].Equipment_State = "状态异常"
					} else {
						tableData_multiple_equipment[temp_count].State_flag = "3"
						tableData_multiple_equipment[temp_count].Equipment_State = "修改失败"
						count_fail_e = count_fail_e + 1
					}
				} else {
					if (tableDate.Succeedflag == "1") || (tableDate.Handleflag == "1") {
						tableData_multiple_equipment[temp_count].State_flag = "5"
						tableData_multiple_equipment[temp_count].Equipment_State = "状态异常"
					} else {
						tableData_multiple_equipment[temp_count].State_flag = "0"
						tableData_multiple_equipment[temp_count].Equipment_State = "未修改"
					}
				}
			}
		}
		// 参数信息处理
		common.DB.Table("parameter_batch_sets").Where("appliance_id=? ", tableDate.Appliance_id).Order("check_alter ASC").Scan(&tableData_multiple_para_line)

		if len(tableData_multiple_para_line) >= 2 {
			P_N := Struct2Map(tableData_multiple_para_line[0])
			P_O := Struct2Map(tableData_multiple_para_line[1])
			if tableData_multiple_para_line[0].CheckAlter == "1" {
				tableData_multiple_equipment[temp_count].Equipment_Para_T = tableData_multiple_para_line[0].AlterCode
				tableData_multiple_equipment[temp_count].Equipment_Para_N = Strval(P_N[tableData_multiple_equipment[temp_count].Equipment_Para_T])
				tableData_multiple_equipment[temp_count].Equipment_Para_O = Strval(P_O[tableData_multiple_equipment[temp_count].Equipment_Para_T])

			} else {
				tableData_multiple_equipment[temp_count].Equipment_Para_T = tableData_multiple_para_line[0].AlterCode
				tableData_multiple_equipment[temp_count].Equipment_Para_O = Strval(P_N[tableData_multiple_equipment[temp_count].Equipment_Para_T])
				tableData_multiple_equipment[temp_count].Equipment_Para_N = Strval(P_O[tableData_multiple_equipment[temp_count].Equipment_Para_T])
			}

		} else {
			tableData_multiple_equipment[temp_count].Equipment_Para_T = "无记录"
			tableData_multiple_equipment[temp_count].Equipment_Para_N = "-"
			tableData_multiple_equipment[temp_count].Equipment_Para_O = "-"
		}

	}
	count_equipment_remain = count_all_rewrite - count_all_rewrite_ok
	if count_equipment_remain != 0 && refresh_flage == "-1" {
		// temp_num_e := count_all_rewrite - count_all_rewrite_ok
		count_equipment_remain = 100 - (count_equipment_remain * 100 / count_all_rewrite)
		response.Success(ctx, gin.H{"busy_flag": "1", "loading": count_equipment_remain, "count_equipment": count_all_rewrite, "count_equipment_ok": count_all_rewrite_ok, "tableData_equipment_info": tableData_multiple_equipment, "count": count}, "成功")
	} else {
		temp_num_e := count_all_rewrite - count_all_rewrite_sucess
		response.Success(ctx, gin.H{"busy_flag": "0", "loading": "100", "tableData_equipment_info": tableData_multiple_equipment, "count_equipment": count_all_rewrite, "count_equipment_ok": count_all_rewrite_ok, "count_success_e": count_all_rewrite_sucess, "count_fail_e": temp_num_e, "count": count}, "成功")
	}
}

func (t TimeMonitorController) EquipmentIndexMonitoring(ctx *gin.Context) {

	Appliance_id := ctx.DefaultQuery("dev_id", "179220395415410")
	multiple_falg := ctx.DefaultQuery("multiple_falg", "10086")
	fmt.Println(Appliance_id)
	var Parameter_final []model.ParameterOutput
	var Parameter_defaults []model.ParameterDefaultsOutput
	var Parameter_defaults_one model.ParameterDefaultsOutput
	var Parameter_final_one model.ParameterOutput
	// var Parameter_Out model.NonParameterOutput

	var tableData_para_temp_out []model.TableData_para_temp_struct
	var tableData_para_temp model.TableData_para_temp_struct
	var tableData_para_un_1_out []model.TableData_para_temp_struct
	var tableData_para_un_1 model.TableData_para_temp_struct
	var tableData_para_un_2_out []model.TableData_para_temp_struct
	var tableData_para_un_2 model.TableData_para_temp_struct
	var tableData_para_single_out []model.TableData_para_single_struct
	var tableData_para_single model.TableData_para_single_struct
	var code_para []model.Code_para_struct
	var code_Serial_para model.Code_para_Serial_struct
	common.DB.Table("parameter_final_settings").Where(" appliance_id = ? ", Appliance_id).Find(&Parameter_final)
	common.DB.Table("parameter_defaults").Where(" appliance_id = ? ", Appliance_id).Find(&Parameter_defaults)
	common.DB.Table("parameter_codes").Scan(&code_para)
	if multiple_falg == "1" {
		for i := 62; i < 84; i++ {
			temp_min, _ := strconv.Atoi(code_para[i].Min)
			temp_max, _ := strconv.Atoi(code_para[i].Max)
			tableData_para_single.Para_name = code_para[i].Parameter
			tableData_para_single.Code = code_para[i].Code
			tableData_para_single.Min_limit = strconv.FormatInt(int64(temp_min), 16)
			tableData_para_single.Max_limit = strconv.FormatInt(int64(temp_max), 16)

			common.DB.Table("parameter_serials").Where("  parameter = ? ", tableData_para_single.Para_name).Find(&code_Serial_para)
			tableData_para_single.Serial_number = code_Serial_para.Serial_number
			tableData_para_single.IsEdit = false
			tableData_para_single_out = append(tableData_para_single_out, tableData_para_single)
		}
		response.Success(ctx, gin.H{"tableData_para_single_out": tableData_para_single_out}, "成功")
		return
	}
	// for i := 0; i < len(Parameter_final); i++ {
	// 	code := Parameter_final[i].Code
	// 	switch {
	// 	case code == "000013":
	// 		Parameter_Out.PH = Parameter_final[i].CurrentValue
	// 	case code == "000014":
	// 		Parameter_Out.FH = Parameter_final[i].CurrentValue
	// 	case code == "000015":
	// 		Parameter_Out.PL = Parameter_final[i].CurrentValue
	// 	case code == "000016":
	// 		Parameter_Out.FL = Parameter_final[i].CurrentValue
	// 	case code == "000017":
	// 		Parameter_Out.DH = Parameter_final[i].CurrentValue
	// 	case code == "000018":
	// 		Parameter_Out.Fd = Parameter_final[i].CurrentValue
	// 	case code == "000019":
	// 		Parameter_Out.CH = Parameter_final[i].CurrentValue
	// 	case code == "000020":
	// 		Parameter_Out.FC = Parameter_final[i].CurrentValue
	// 	case code == "010113":
	// 		Parameter_Out.MaximumLoadFanCurrentDeviationCoefficient = Parameter_final[i].CurrentValue
	// 	case code == "010114":
	// 		Parameter_Out.MinimumLoadFanCurrentDeviationCoefficient = Parameter_final[i].CurrentValue
	// 	case code == "010115":
	// 		Parameter_Out.MaximumLoadFanDutyCycleDeviationCoefficient = Parameter_final[i].CurrentValue
	// 	case code == "010116":
	// 		Parameter_Out.MinimumLoadFanDutyCycleDeviationCoefficient = Parameter_final[i].CurrentValue
	// 	case code == "010117":
	// 		Parameter_Out.BackwaterFlowValue = Parameter_final[i].CurrentValue
	// 	case code == "010118":
	// 		Parameter_Out.FrequencyCompensationValueOfWindPressureSensorAlarmPoint = Parameter_final[i].CurrentValue
	// 	case code == "020214":
	// 		Parameter_Out.KA0 = Parameter_final[i].CurrentValue
	// 	case code == "020314":
	// 		Parameter_Out.KA1 = Parameter_final[i].CurrentValue
	// 	case code == "020414":
	// 		Parameter_Out.KA2 = Parameter_final[i].CurrentValue
	// 	case code == "020514":
	// 		Parameter_Out.KA3 = Parameter_final[i].CurrentValue
	// 	case code == "020215":
	// 		Parameter_Out.KB0 = Parameter_final[i].CurrentValue
	// 	case code == "020315":
	// 		Parameter_Out.KB1 = Parameter_final[i].CurrentValue
	// 	case code == "020415":
	// 		Parameter_Out.KB2 = Parameter_final[i].CurrentValue
	// 	case code == "020515":
	// 		Parameter_Out.KB3 = Parameter_final[i].CurrentValue
	// 	case code == "020216":
	// 		Parameter_Out.KC0 = Parameter_final[i].CurrentValue
	// 	case code == "020316":
	// 		Parameter_Out.KC1 = Parameter_final[i].CurrentValue
	// 	case code == "020416":
	// 		Parameter_Out.KC2 = Parameter_final[i].CurrentValue
	// 	case code == "020516":
	// 		Parameter_Out.KC3 = Parameter_final[i].CurrentValue
	// 	case code == "020217":
	// 		Parameter_Out.KF0 = Parameter_final[i].CurrentValue
	// 	case code == "020317":
	// 		Parameter_Out.KF1 = Parameter_final[i].CurrentValue
	// 	case code == "020417":
	// 		Parameter_Out.KF2 = Parameter_final[i].CurrentValue
	// 	case code == "020517":
	// 		Parameter_Out.KF3 = Parameter_final[i].CurrentValue
	// 	case code == "020218":
	// 		Parameter_Out.T1A0 = Parameter_final[i].CurrentValue
	// 	case code == "020318":
	// 		Parameter_Out.T1A1 = Parameter_final[i].CurrentValue
	// 	case code == "020418":
	// 		Parameter_Out.T1A2 = Parameter_final[i].CurrentValue
	// 	case code == "020518":
	// 		Parameter_Out.T1A3 = Parameter_final[i].CurrentValue
	// 	case code == "020219":
	// 		Parameter_Out.T1C0 = Parameter_final[i].CurrentValue
	// 	case code == "020319":
	// 		Parameter_Out.T1C1 = Parameter_final[i].CurrentValue
	// 	case code == "020419":
	// 		Parameter_Out.T1C2 = Parameter_final[i].CurrentValue
	// 	case code == "020519":
	// 		Parameter_Out.T1C3 = Parameter_final[i].CurrentValue
	// 	case code == "020220":
	// 		Parameter_Out.T2A0 = Parameter_final[i].CurrentValue
	// 	case code == "020320":
	// 		Parameter_Out.T2A1 = Parameter_final[i].CurrentValue
	// 	case code == "020420":
	// 		Parameter_Out.T2A2 = Parameter_final[i].CurrentValue
	// 	case code == "020520":
	// 		Parameter_Out.T2A3 = Parameter_final[i].CurrentValue
	// 	case code == "020221":
	// 		Parameter_Out.T2C0 = Parameter_final[i].CurrentValue
	// 	case code == "020321":
	// 		Parameter_Out.T2C1 = Parameter_final[i].CurrentValue
	// 	case code == "020421":
	// 		Parameter_Out.T2C2 = Parameter_final[i].CurrentValue
	// 	case code == "020521":
	// 		Parameter_Out.T2C3 = Parameter_final[i].CurrentValue
	// 	case code == "020222":
	// 		Parameter_Out.TDA0 = Parameter_final[i].CurrentValue
	// 	case code == "020322":
	// 		Parameter_Out.TDA1 = Parameter_final[i].CurrentValue
	// 	case code == "020422":
	// 		Parameter_Out.TDA2 = Parameter_final[i].CurrentValue
	// 	case code == "020522":
	// 		Parameter_Out.TDA2 = Parameter_final[i].CurrentValue
	// 	case code == "020223":
	// 		Parameter_Out.TDC0 = Parameter_final[i].CurrentValue
	// 	case code == "020323":
	// 		Parameter_Out.TDC1 = Parameter_final[i].CurrentValue
	// 	case code == "020423":
	// 		Parameter_Out.TDC2 = Parameter_final[i].CurrentValue
	// 	case code == "020523":
	// 		Parameter_Out.TDC3 = Parameter_final[i].CurrentValue
	// 	case code == "020224":
	// 		Parameter_Out.WC0 = Parameter_final[i].CurrentValue
	// 	case code == "020324":
	// 		Parameter_Out.WC1 = Parameter_final[i].CurrentValue
	// 	case code == "020424":
	// 		Parameter_Out.WC2 = Parameter_final[i].CurrentValue
	// 	case code == "020524":
	// 		Parameter_Out.WC3 = Parameter_final[i].CurrentValue
	// 	case code == "020225":
	// 		Parameter_Out.WO0 = Parameter_final[i].CurrentValue
	// 	case code == "020325":
	// 		Parameter_Out.WO1 = Parameter_final[i].CurrentValue
	// 	case code == "020425":
	// 		Parameter_Out.WO2 = Parameter_final[i].CurrentValue
	// 	case code == "020525":
	// 		Parameter_Out.WO3 = Parameter_final[i].CurrentValue
	// 	default:
	// 		fmt.Println("切片为空!")
	// 	}
	// }
	fmt.Println(code_para)
	for i := 62; i < 84; i++ {
		temp_min, _ := strconv.Atoi(code_para[i].Min)
		temp_max, _ := strconv.Atoi(code_para[i].Max)
		tableData_para_single.Para_name = code_para[i].Parameter
		tableData_para_single.Code = code_para[i].Code
		tableData_para_single.Min_limit = strconv.FormatInt(int64(temp_min), 16)
		tableData_para_single.Max_limit = strconv.FormatInt(int64(temp_max), 16)
		common.DB.Table("parameter_defaults").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_defaults_one)
		tableData_para_single.Default_value = strconv.FormatInt(int64(Parameter_defaults_one.Default_value), 16)
		common.DB.Table("parameter_final_settings").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_final_one)
		tableData_para_single.Recent_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_single.Modify_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		common.DB.Table("parameter_serials").Where("  parameter = ? ", tableData_para_single.Para_name).Find(&code_Serial_para)
		tableData_para_single.Serial_number = code_Serial_para.Serial_number
		tableData_para_single.IsEdit = false
		tableData_para_single_out = append(tableData_para_single_out, tableData_para_single)
	}
	for i := 8; i < 12; i++ {
		temp_min, _ := strconv.Atoi(code_para[i].Min)
		temp_max, _ := strconv.Atoi(code_para[i].Max)
		tableData_para_un_2.Para_name = code_para[i].Parameter
		tableData_para_un_2.Code = code_para[i].Code
		tableData_para_un_2.Min_limit = strconv.FormatInt(int64(temp_min), 16)
		tableData_para_un_2.Max_limit = strconv.FormatInt(int64(temp_max), 16)
		common.DB.Table("parameter_defaults").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_defaults_one)
		tableData_para_un_2.Default_value = strconv.FormatInt(int64(Parameter_defaults_one.Default_value), 16)
		common.DB.Table("parameter_final_settings").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_final_one)
		tableData_para_un_2.Recent_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_un_2.Modify_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_un_2.IsEdit = false
		tableData_para_un_2_out = append(tableData_para_un_2_out, tableData_para_un_2)
	}
	for i := 12; i < 14; i++ {
		temp_min, _ := strconv.Atoi(code_para[i].Min)
		temp_max, _ := strconv.Atoi(code_para[i].Max)
		tableData_para_un_1.Para_name = code_para[i].Parameter
		tableData_para_un_1.Code = code_para[i].Code
		tableData_para_un_1.Min_limit = strconv.FormatInt(int64(temp_min), 16)
		tableData_para_un_1.Max_limit = strconv.FormatInt(int64(temp_max), 16)
		common.DB.Table("parameter_defaults").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_defaults_one)
		tableData_para_un_1.Default_value = strconv.FormatInt(int64(Parameter_defaults_one.Default_value), 16)
		common.DB.Table("parameter_final_settings").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_final_one)
		tableData_para_un_1.Recent_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_un_1.Modify_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_un_1.IsEdit = false
		tableData_para_un_1_out = append(tableData_para_un_1_out, tableData_para_un_1)
	}
	for i := 14; i < 62; i++ {
		temp_min, _ := strconv.Atoi(code_para[i].Min)
		temp_max, _ := strconv.Atoi(code_para[i].Max)
		tableData_para_temp.Para_name = code_para[i].Parameter
		tableData_para_temp.Code = code_para[i].Code
		tableData_para_temp.Min_limit = strconv.FormatInt(int64(temp_min), 16)
		tableData_para_temp.Max_limit = strconv.FormatInt(int64(temp_max), 16)
		common.DB.Table("parameter_defaults").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_defaults_one)
		tableData_para_temp.Default_value = strconv.FormatInt(int64(Parameter_defaults_one.Default_value), 16)
		common.DB.Table("parameter_final_settings").Where("  appliance_id = ? AND code = ? ", Appliance_id, code_para[i].Code).Find(&Parameter_final_one)
		tableData_para_temp.Recent_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_temp.Modify_date = strconv.FormatInt(int64(Parameter_final_one.CurrentValue), 16)
		tableData_para_temp.IsEdit = false
		tableData_para_temp_out = append(tableData_para_temp_out, tableData_para_temp)
	}
	fmt.Println(tableData_para_temp_out)

	response.Success(ctx, gin.H{"tableData_para_temp_out": tableData_para_temp_out, "tableData_para_un_2_out": tableData_para_un_2_out, "tableData_para_un_1_out": tableData_para_un_1_out, "tableData_para_single_out": tableData_para_single_out}, "成功")
}

func (t TimeMonitorController) EquipmentInfoDelete(ctx *gin.Context) {
	Appliance_id := ctx.QueryArray("applianceid[]")
	var notid int64
	var haveid int64
	var MultipleParaRewrite []model.MultipleParaRewrite
	var ParamenBatch []model.ParamenBatch
	fmt.Println("删除设备：", Appliance_id)
	for _, singel_id := range Appliance_id {
		Result := common.DB.Table("multiple_para_rewrite").Where("appliance_id = ? ", singel_id).First(&model.MultipleParaRewrite{}).Error
		if errors.Is(Result, gorm.ErrRecordNotFound) {
			fmt.Println("没有找到设备")
			notid++
		} else {
			common.DB.Table("multiple_para_rewrite").Where("appliance_id = ?", singel_id).Delete(&MultipleParaRewrite)
			common.DB.Table("parameter_batch_sets").Where("appliance_id = ?", singel_id).Delete(&ParamenBatch)
			haveid++
		}
	}
	response.Success(ctx, gin.H{"upload_ok": haveid, "upload_not": notid}, "成功")
}
func (t TimeMonitorController) GetDelete(ctx *gin.Context) { //GetDelete
	file, err := ctx.FormFile("file_D")
	if err != nil {
		fmt.Println(err)
		response.Success(ctx, gin.H{"success": "0"}, "失败")
		return
	}
	fmt.Println(file.Filename)
	fileSuffix := path.Ext(file.Filename)
	if fileSuffix == ".xls" || fileSuffix == ".xlsx" {
		dst := fmt.Sprintf("/root/aliyun/%s", file.Filename)
		fmt.Println(1233)
		ctx.SaveUploadedFile(file, dst)
		fmt.Println(dst)
		readExceldeldect(dst)
		fmt.Println(1233)
		// 上传文件到指定的目录

		response.Success(ctx, gin.H{"success": "1"}, "成功")
		return
	} else if fileSuffix == ".csv" {
		nowstar := time.Now()
		var Eup []string
		if err := ctx.SaveUploadedFile(file, file.Filename); err != nil {
			response.Success(ctx, gin.H{"error": "0"}, "失败")
			os.Remove(file.Filename)
			return
		}
		clientsFile, err := os.OpenFile(file.Filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			response.Success(ctx, gin.H{"error": "0"}, "失败")
			os.Remove(file.Filename)
			return
		}
		//关闭流
		defer clientsFile.Close()
		//解析csv文件到结构体，clients为自己定义的结构体

		r := csv.NewReader(clientsFile)
		r.FieldsPerRecord = -1
		rows, err := r.ReadAll()
		if err != nil {
			response.Success(ctx, gin.H{"error": "0"}, "失败")
			os.Remove(file.Filename)
			return
		}
		// 处理
		// writeDatas := [][]string{
		// 	[]string{"columnID(栏目ID)","ChannelID(专辑ID)","columnSourceId(栏目数据源表ID)","columnName(栏目名)","URL"},
		// }
		for index, item := range rows {
			// 跳过第一行
			if index < 1 {
				Stamp := strings.TrimSpace(item[0])
				ID := strings.TrimSpace(item[1])
				if ID != "ID" || Stamp != "Stamp" {
					os.Remove(file.Filename)
					return
				}
			} else {
				// Stamp := strings.TrimSpace(item[0])
				ID := strings.TrimSpace(item[1])
				if ID != "null" {
					Eup = append(Eup, ID)
				}
			}
		}
		TaskAssignmendelect(7, Eup)
		nowend := time.Now()
		timeLength := nowend.Sub(nowstar) //两个时间相减，计算时间差
		fmt.Println(timeLength.Seconds(), "s")
		fmt.Print("\n")
		os.Remove(file.Filename)

	}

}

func readExceldeldect(excelPath string) {
	// var Multiple_Equipment_info []model.Multiple_Equipment_info_input
	//defer db.Close()
	// nowstar := time.Now()
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		fmt.Println(err)
		fmt.Println(222)
		return
	}
	var Eup []string
	// 获取工作表中指定单元格的值
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	nowstar := time.Now()
	for key, row := range rows {
		// 去掉标题行
		if key > 0 {
			for _, colCell := range row {
				//   fmt.Print(colCell, "\t")
				Eup = append(Eup, colCell)
				//  fmt.Print(Eup)
			}
		}

	}
	TaskAssignmendelect(7, Eup)
	nowend := time.Now()
	timeLength := nowend.Sub(nowstar) //两个时间相减，计算时间差
	fmt.Println(timeLength.Seconds(), "s")
	fmt.Print("\n")
}
func TaskAssignmendelect(thread int, equ []string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 1100)
	results := make(chan int, 1100)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go workerSingledelct(jobs, results, equ)
	}
	// 每台设备一个任务
	for j := 0; j < len(equ); j++ {
		jobs <- j
		work <- equ[j]
	}
	close(jobs)

	// 输出结果
	for a := 0; a < len(equ); a++ {
		<-results
	}
}

/*分配工作*/
func workerSingledelct(jobs <-chan int, results chan<- int, equ []string) {
	var notid int64
	var haveid int64
	var MultipleParaRewrite []model.MultipleParaRewrite
	var ParamenBatch []model.ParamenBatch
	for i := range jobs {
		Result := common.DB.Table("multiple_para_rewrite").Where("appliance_id = ? ", equ[i]).First(&model.MultipleParaRewrite{}).Error
		if errors.Is(Result, gorm.ErrRecordNotFound) {
			fmt.Println("没有找到设备")
			notid++
		} else {
			common.DB.Table("multiple_para_rewrite").Where("appliance_id = ?", equ[i]).Delete(&MultipleParaRewrite)
			common.DB.Table("parameter_batch_sets").Where("appliance_id = ?", equ[i]).Delete(&ParamenBatch)
			haveid++
		}
		results <- i
	}

}
func (t TimeMonitorController) EquipmentHistoryMonitoring(ctx *gin.Context) {
	Appliance_id := ctx.DefaultQuery("applianceid", "179220395415410")
	fmt.Println(Appliance_id)
	Code_singel := ctx.DefaultQuery("code", "0")
	Single_flag := ctx.DefaultQuery("single_flag", "0")
	starttime := ctx.DefaultQuery("starttime", "0")
	endtime := ctx.DefaultQuery("endtime", "0")
	// data_time := ctx.DefaultQuery("starttime", "")
	fmt.Println(starttime)
	fmt.Println(endtime)
	var para_history []model.Para_history_struct
	var tableData_para_un_2_history []model.Para_history_out_struct
	var tableData_para_un_2_history_line model.Para_history_out_struct
	var tableData_para_single_history []model.Para_history_out_struct
	var tableData_para_single_history_line model.Para_history_out_struct
	var tableData_para_temp_history []model.Para_history_out_struct
	var tableData_para_temp_history_line model.Para_history_out_struct
	var tableData_para_un_1_history []model.Para_history_out_struct
	var tableData_para_un_1_history_line model.Para_history_out_struct
	var tableData_para_one_code []model.Para_history_out_struct
	var tableData_para_one_code_line model.Para_history_out_struct

	// if data_time == "" {

	var code_para []model.Code_para_struct
	common.DB.Table("parameter_codes").Scan(&code_para)
	if Single_flag != "0" {
		var Code_name string
		common.DB.Table("parameter_changes_settings").Where(" appliance_id = ? and code = ? ", Appliance_id, Code_singel).Order("updatetime DESC").Find(&para_history)
		fmt.Println(para_history)
		for _, tableDate := range code_para {
			if tableDate.Code == Code_singel {
				Code_name = tableDate.Parameter
				break
			}
		}
		for _, tableDate := range para_history {
			tableData_para_one_code_line.Code = tableDate.Code
			tableData_para_one_code_line.Para_name = Code_name
			tableData_para_one_code_line.LastValue = strconv.FormatInt(int64(tableDate.LastValue), 16)
			tableData_para_one_code_line.Value = strconv.FormatInt(int64(tableDate.Value), 16)
			tableData_para_one_code_line.Appliance_id = tableDate.Appliance_id
			tableData_para_one_code_line.Updatetime = tableDate.Updatetime
			tableData_para_one_code_line.LatestParameterFlag = tableDate.LatestParameterFlag
			tableData_para_one_code = append(tableData_para_one_code, tableData_para_one_code_line)
		}
		fmt.Println(tableData_para_one_code)
		response.Success(ctx, gin.H{"tableData_para_one_code": tableData_para_one_code}, "成功")
		return

	}

	common.DB.Table("parameter_changes_settings").Where(" appliance_id = ? and updatetime >= ? and updatetime <=?", Appliance_id, starttime, endtime).Order("updatetime DESC").Find(&para_history)
	fmt.Println(para_history)
	for _, tableDate := range para_history {
		for i := 8; i < 12; i++ {
			if tableDate.Code == code_para[i].Code {
				tableData_para_un_2_history_line.Code = code_para[i].Code
				tableData_para_un_2_history_line.Para_name = code_para[i].Parameter
				tableData_para_un_2_history_line.LastValue = strconv.FormatInt(int64(tableDate.LastValue), 16)
				tableData_para_un_2_history_line.Value = strconv.FormatInt(int64(tableDate.Value), 16)
				tableData_para_un_2_history_line.Appliance_id = tableDate.Appliance_id
				tableData_para_un_2_history_line.Updatetime = tableDate.Updatetime
				tableData_para_un_2_history_line.LatestParameterFlag = tableDate.LatestParameterFlag
				tableData_para_un_2_history = append(tableData_para_un_2_history, tableData_para_un_2_history_line)

			}
		}
		for i := 12; i < 14; i++ {
			if tableDate.Code == code_para[i].Code {
				tableData_para_un_1_history_line.Code = code_para[i].Code
				tableData_para_un_1_history_line.Para_name = code_para[i].Parameter
				tableData_para_un_1_history_line.LastValue = strconv.FormatInt(int64(tableDate.LastValue), 16)
				tableData_para_un_1_history_line.Value = strconv.FormatInt(int64(tableDate.Value), 16)
				tableData_para_un_1_history_line.Appliance_id = tableDate.Appliance_id
				tableData_para_un_1_history_line.Updatetime = tableDate.Updatetime
				tableData_para_un_1_history_line.LatestParameterFlag = tableDate.LatestParameterFlag
				tableData_para_un_1_history = append(tableData_para_un_1_history, tableData_para_un_1_history_line)
			}
		}
		for i := 14; i < 62; i++ {
			if tableDate.Code == code_para[i].Code {
				tableData_para_temp_history_line.Code = code_para[i].Code
				tableData_para_temp_history_line.Para_name = code_para[i].Parameter
				tableData_para_temp_history_line.LastValue = strconv.FormatInt(int64(tableDate.LastValue), 16)
				tableData_para_temp_history_line.Value = strconv.FormatInt(int64(tableDate.Value), 16)
				tableData_para_temp_history_line.Appliance_id = tableDate.Appliance_id
				tableData_para_temp_history_line.Updatetime = tableDate.Updatetime
				tableData_para_temp_history_line.LatestParameterFlag = tableDate.LatestParameterFlag
				tableData_para_temp_history = append(tableData_para_temp_history, tableData_para_temp_history_line)
			}
		}
		for i := 62; i < 84; i++ {
			if tableDate.Code == code_para[i].Code {
				tableData_para_single_history_line.Code = code_para[i].Code
				tableData_para_single_history_line.Para_name = code_para[i].Parameter
				tableData_para_single_history_line.LastValue = strconv.FormatInt(int64(tableDate.LastValue), 16)
				tableData_para_single_history_line.Value = strconv.FormatInt(int64(tableDate.Value), 16)
				tableData_para_single_history_line.Appliance_id = tableDate.Appliance_id
				tableData_para_single_history_line.Updatetime = tableDate.Updatetime
				tableData_para_single_history_line.LatestParameterFlag = tableDate.LatestParameterFlag
				tableData_para_single_history = append(tableData_para_single_history, tableData_para_single_history_line)
			}
		}
	}
	fmt.Println(tableData_para_single_history)

	response.Success(ctx, gin.H{"tableData_para_temp_history": tableData_para_temp_history, "tableData_para_un_1_history": tableData_para_un_1_history, "tableData_para_un_2_history": tableData_para_un_2_history, "tableData_para_single_history": tableData_para_single_history}, "成功")
}
func (t TimeMonitorController) GetUpdata(ctx *gin.Context) {
	file, err := ctx.FormFile("file_G")
	if err != nil {
		fmt.Println(err)
		response.Success(ctx, gin.H{"success": "0"}, "失败")
		return
	}
	fmt.Println(file.Filename)
	fileSuffix := path.Ext(file.Filename)
	if fileSuffix == ".xls" || fileSuffix == ".xlsx" {
		dst := fmt.Sprintf("/root/aliyun/%s", file.Filename)
		fmt.Println(1233)
		ctx.SaveUploadedFile(file, dst)
		fmt.Println(dst)
		readExcel(dst)
		fmt.Println(1233)
		// 上传文件到指定的目录

		response.Success(ctx, gin.H{"success": "1"}, "成功")
		return
	} else if fileSuffix == ".csv" {
		nowstar := time.Now()
		var Eup []string
		if err := ctx.SaveUploadedFile(file, file.Filename); err != nil {
			response.Success(ctx, gin.H{"error": "0"}, "失败")
			os.Remove(file.Filename)
			return
		}
		clientsFile, err := os.OpenFile(file.Filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			response.Success(ctx, gin.H{"error": "0"}, "失败")
			os.Remove(file.Filename)
			return
		}
		//关闭流
		defer clientsFile.Close()
		//解析csv文件到结构体，clients为自己定义的结构体

		r := csv.NewReader(clientsFile)
		r.FieldsPerRecord = -1
		rows, err := r.ReadAll()
		if err != nil {
			response.Success(ctx, gin.H{"error": "0"}, "失败")
			os.Remove(file.Filename)
			return
		}
		// 处理
		// writeDatas := [][]string{
		// 	[]string{"columnID(栏目ID)","ChannelID(专辑ID)","columnSourceId(栏目数据源表ID)","columnName(栏目名)","URL"},
		// }
		for index, item := range rows {
			// 跳过第一行
			if index < 1 {
				Stamp := strings.TrimSpace(item[0])
				ID := strings.TrimSpace(item[1])
				if ID != "ID" || Stamp != "Stamp" {
					os.Remove(file.Filename)
					return
				}
			} else {
				// Stamp := strings.TrimSpace(item[0])
				ID := strings.TrimSpace(item[1])
				if ID != "null" {
					Eup = append(Eup, ID)
				}
			}
		}
		TaskAssignment_my(7, Eup)
		nowend := time.Now()
		timeLength := nowend.Sub(nowstar) //两个时间相减，计算时间差
		fmt.Println(timeLength.Seconds(), "s")
		fmt.Print("\n")
		os.Remove(file.Filename)

	}

}

func readExcel(excelPath string) {
	// var Multiple_Equipment_info []model.Multiple_Equipment_info_input
	//defer db.Close()
	// nowstar := time.Now()
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		fmt.Println(err)
		fmt.Println(222)
		return
	}
	var Eup []string
	// 获取工作表中指定单元格的值
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	nowstar := time.Now()
	for key, row := range rows {
		// 去掉标题行
		if key > 0 {
			for _, colCell := range row {
				//   fmt.Print(colCell, "\t")
				Eup = append(Eup, colCell)
				//  fmt.Print(Eup)
			}
		}

	}
	TaskAssignment_my(7, Eup)
	nowend := time.Now()
	timeLength := nowend.Sub(nowstar) //两个时间相减，计算时间差
	fmt.Println(timeLength.Seconds(), "s")
	fmt.Print("\n")
}

func TaskAssignment_my(thread int, equ []string) {
	//最大处理台数，处理完程序阻塞（值设置大于需处理台数）
	jobs := make(chan int, 100)
	work := make(chan string, 1100)
	results := make(chan int, 1100)
	//最大线程数，通过INI配置文件设置(一般一个核2—4个线程)
	runtime.GOMAXPROCS(thread)
	// 开启goroutine
	for w := 1; w <= thread; w++ {
		go workerSingle_my(jobs, results, equ)
	}
	// 每台设备一个任务
	for j := 0; j < len(equ); j++ {
		jobs <- j
		work <- equ[j]
	}
	close(jobs)

	// 输出结果
	for a := 0; a < len(equ); a++ {
		<-results
	}
}

/*分配工作*/
func workerSingle_my(jobs <-chan int, results chan<- int, equ []string) {
	// db, err := gorm.Open("mysql", "root:582235@(127.0.0.1:3306)/aliyun?charset=utf8mb4")
	// if err != nil {
	// 	panic(err)
	// }
	var Multiple_Equipment_info_row model.Multiple_Equipment_info_input
	for i := range jobs {
		//BasendCommand(equ[i])
		Multiple_Equipment_info_row.Appliance_id = equ[i]
		Multiple_Equipment_info_row.Handleflag = "0"
		Multiple_Equipment_info_row.Succeedflag = "0"
		Multiple_Equipment_info_row.Rewriteflag = "0"
		common.DB.Table("multiple_para_rewrite").Create(&Multiple_Equipment_info_row)
		results <- i
		//fmt.Println("线程i:", i)
	}
}

func (t TimeMonitorController) EquipmentInfoCreat(ctx *gin.Context) {
	Appliance_id := ctx.QueryArray("applianceid[]")
	var Multiple_Equipment_info_row model.Multiple_Equipment_info_input
	for _, singel_id := range Appliance_id {
		Multiple_Equipment_info_row.Appliance_id = singel_id
		Multiple_Equipment_info_row.Handleflag = "0"
		Multiple_Equipment_info_row.Succeedflag = "0"
		Multiple_Equipment_info_row.Rewriteflag = "0"
		common.DB.Table("multiple_para_rewrite").Create(&Multiple_Equipment_info_row)
	}
	response.Success(ctx, gin.H{"upload_ok": Appliance_id}, "成功")
}

func (t TimeMonitorController) MultipleEquipmentInfoSingel(ctx *gin.Context) {
	Appliance_id := ctx.DefaultQuery("applianceid", "179220395415410")
	var para_all []model.ParamenBatch
	var para_all_new model.ParamenBatch
	var para_all_old model.ParamenBatch
	var para_all_save model.ParamenBatch

	common.DB.Table("parameter_batch_sets").Where(" appliance_id = ?", Appliance_id).Find(&para_all)

	for _, para_info := range para_all {
		if para_info.CheckAlter == "1" {
			para_all_new = para_info
		}
		if para_info.CheckAlter == "0" {
			para_all_old = para_info
		}
		if para_info.CheckAlter == "2" {
			para_all_save = para_info
		}
	}
	response.Success(ctx, gin.H{"para_all_new": para_all_new, "para_all_old": para_all_old, "para_all_save": para_all_save}, "成功")

}

func (t TimeMonitorController) PageList(ctx *gin.Context) {
	dev_idshou := ctx.DefaultQuery("dev_id", "109951162778761")
	timeLow := ctx.DefaultQuery("timeLow", "2020-01-01")
	timeHigh := ctx.DefaultQuery("timeHigh", "2022-01-01")
	flag := ctx.DefaultQuery("flag", "0")
	downloadflag := ctx.DefaultQuery("downloadflag", "0")
	currentPage, _ := strconv.Atoi(ctx.DefaultQuery("currentPage", "0"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "0"))

	//averagelow := ctx.DefaultQuery("averagelow", "0")
	//averagehigh := ctx.DefaultQuery("averagehigh", "150")
	//model2:=ctx.DefaultQuery("model2","1")
	//model3 := ctx.DefaultQuery("model", "0")
	//timelength := ctx.DefaultQuery("timelength", "0")
	//month1:=ctx.DefaultQuery("month","1")+"-01"
	//month2:=ctx.DefaultQuery("month1","1")+"-01"
	//month3:=ctx.DefaultQuery("month2","1")
	////month4:=ctx.DefaultQuery("month3","1")
	//pageType:=ctx.DefaultQuery("pagetype","2")
	//categoryType:=ctx.QueryArray("type1[]")
	//city:=ctx.DefaultQuery("city","0")
	//city1:=ctx.DefaultQuery("city1","0")
	//flag:=ctx.DefaultQuery("flag","0")
	//flag1:=ctx.DefaultQuery("flag1","0")
	//scorelow:= ctx.DefaultQuery("scorelow", "0")
	//scorehigh:= ctx.DefaultQuery("scorehigh", "99")
	//citySummaryFlag,_:=strconv.ParseBool(ctx.DefaultQuery("citySummaryFlag", "false"))
	//provinceSummaryFlag,_:= strconv.ParseBool(ctx.DefaultQuery("provinceSummaryFlag", "false"))
	//typeSummaryFlag,_:= strconv.ParseBool(ctx.DefaultQuery("typeSummaryFlag", "false"))
	//yearflag,_:= strconv.Atoi(ctx.DefaultQuery("yearflag", "0"))
	//monthflag,_:= strconv.Atoi(ctx.DefaultQuery("monthflag", "1"))
	//dayflag,_:= strconv.Atoi(ctx.DefaultQuery("dayflag", "0"))
	//start:= ctx.DefaultQuery("start", "2021-01-01")
	//end:= ctx.DefaultQuery("end", "2021-07-01")
	//Flag,_:=strconv.ParseBool(ctx.DefaultQuery("Flag", "false"))
	//timeStamp:= ctx.DefaultQuery("timeStamp", "15:02")
	//timeLowInt:=util.TimeParse(timeLow)
	//timeHighInt:=util.TimeParse(timeHigh)

	//var dev_id []string//时间轴
	//var start_time []string
	//var end_time []string
	//var duration_time []string
	//var water_pattern []int
	//var flow_avg []float64
	//var small_water []float64
	//var deviation []float64
	//var up_number []int
	//var down_number []int

	//var heat_duration []string
	//var un_stable_temp_duration []string
	//var un_stable_temp_percent []float32
	//var un_heat_dev []float32
	//var temp_pattern []int//水流量
	//var overshoot_value []int//火焰反馈
	//var state_accuracy []int//火焰反馈
	//var temp_score []int//火焰反馈
	//var new_temp_score []int//火焰反馈
	//var heat_temp_score []int//设定温度
	//var stable_temp_score []int//输出温度
	//var temp_judge_flag  []int
	//var water_flag [] int
	////特征指标
	//var temp_flag []int
	//var abnormal_state [] int

	//var Max_time [] string
	//var Total_num  [] int
	////单月份指标
	var datatime []string  //时间轴
	var datatime1 []string //时间轴
	var flow []int         //水流量
	var flame []string     //火焰反馈
	var settemp []string   //设定温度
	var outtemp []string   //输出温度
	var model1 []int
	var zone_id []string

	var water_score []int
	var temp_score []int     //火焰反馈
	var new_temp_score []int //火焰反馈

	var tableStoreDates []model.MonitorTable
	var tableStoreDates1 []model.RunDate //运行数据库

	var count int

	var tablename1 string
	tablename1 = "e_fragment" + dev_idshou
	var tablename2 string
	tablename2 = "e_data" + dev_idshou
	if flag == "0" { //所有数据
		common.IndexDB.Table(tablename1).Where(" dev_id=? AND start_time>=? AND end_time<=? ", dev_idshou, timeLow, timeHigh).Find(&tableStoreDates)

		common.IndexDB.Table(tablename1).Where(" dev_id=? AND start_time>=? AND end_time<=? ", dev_idshou, timeLow, timeHigh).Count(&count)
		fmt.Println(count)
		common.RunDB.Table(tablename2).Where(" applianceid=? AND datatime BETWEEN ? AND ? ", dev_idshou, timeLow, timeHigh).Find(&tableStoreDates1)
		for _, tableDate := range tableStoreDates {
			datatime1 = append(datatime1, tableDate.Start_time)
			water_score = append(water_score, tableDate.Water_score) //数据放入切片中
			temp_score = append(temp_score, tableDate.Temp_score)
			new_temp_score = append(new_temp_score, tableDate.New_temp_score)
		}

		for _, tableDate := range tableStoreDates1 {
			datatime = append(datatime, tableDate.Datatime) //数据放入切片中
			flow = append(flow, tableDate.Flow)
			flame = append(flame, tableDate.Flame)
			settemp = append(settemp, tableDate.Settemp)
			outtemp = append(outtemp, tableDate.Outtemp)
			model1 = append(model1, tableDate.Water_pattern)
			zone_id = append(zone_id, tableDate.Zone_id)
		}
		fmt.Println(tableStoreDates)
		//fmt.Println(tableStoreDates1)
		response.Success(ctx, gin.H{"data": tableStoreDates, "date1": tableStoreDates1, "data_time1": datatime1, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp, "model": model1, "zone_id": zone_id, "water_score": water_score, "temp_score": temp_score, "new_temp_score": new_temp_score, "total": count}, "成功")

	} else if flag == "1" { //详细数据 某个片段
		common.RunDB.Table(tablename2).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_idshou, timeLow, timeHigh).Find(&tableStoreDates1)
		for _, tableDate := range tableStoreDates1 {
			datatime = append(datatime, tableDate.Datatime) //数据放入切片中
			flow = append(flow, tableDate.Flow)
			flame = append(flame, tableDate.Flame)
			settemp = append(settemp, tableDate.Settemp)
			outtemp = append(outtemp, tableDate.Outtemp)
			model1 = append(model1, tableDate.Water_pattern)
			zone_id = append(zone_id, tableDate.Zone_id)
		}
		fmt.Println(tableStoreDates1)
		fmt.Println(tableStoreDates)
		response.Success(ctx, gin.H{"data": tableStoreDates1, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp, "model": model1,
			"zone_id": zone_id}, "成功")
	} else if flag == "2" {
		if downloadflag == "0" {
			common.IndexDB.Table(tablename1).Where(" dev_id=? AND start_time>=? AND end_time<=? )", dev_idshou, timeLow, timeHigh).Count(&count)
			fmt.Println(count)
			common.IndexDB.Table(tablename1).Where(" dev_id=? AND start_time>=? AND end_time<=?)", dev_idshou, timeLow, timeHigh).Offset(perPage * (currentPage - 1)).Limit(perPage).Find(&tableStoreDates)
		}
		if downloadflag == "1" {
			common.IndexDB.Table(tablename1).Where(" dev_id=? AND start_time>=? AND end_time<=? )", dev_idshou, timeLow, timeHigh).Count(&count)
			fmt.Println(count)
			common.IndexDB.Table(tablename1).Where(" dev_id=? AND start_time>=? AND end_time<=? )", dev_idshou, timeLow, timeHigh).Find(&tableStoreDates)
		}
		for _, tableDate := range tableStoreDates {
			datatime1 = append(datatime1, tableDate.Start_time)
			water_score = append(water_score, tableDate.Water_score) //数据放入切片中
			temp_score = append(temp_score, tableDate.Temp_score)
			new_temp_score = append(new_temp_score, tableDate.New_temp_score)
		}
		fmt.Println(tableStoreDates)
		response.Success(ctx, gin.H{"data": tableStoreDates, "water_score": water_score, "temp_score": temp_score, "new_temp_score": new_temp_score, "total": count}, "成功")
	}

}

//		for _, tableDate := range tableStoreDates {
//
//			datatime = append(datatime, tableDate.Datatime) //数据放入切片中
//			flow = append(flow, tableDate.Flow)
//			flame = append(flame, tableDate.Flame)
//			settemp = append(settemp, tableDate.Settemp)
//			outtemp = append(outtemp, tableDate.Outtemp)
//			model1 = append(model1, tableDate.Water_pattern)
//			zone_id = append(zone_id, tableDate.Zone_id)
//
//		}
//
//		//	response.Success(ctx,gin.H{"data":tableStoreDates,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province},"成功")
//		response.Success(ctx, gin.H{"data": tableStoreDates, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp, "model": model1,
//			"zone_id": zone_id, "data1": tableplace, "total_time": Total_time, "min_time": Min_time, "max_time": Max_time, "total_num": Total_num,
//			"datanum": tableStoreDates3}, "成功")
//
//	} else if flag == "2" {
//		timeLow=timeLow+" 00:00:00"
//		timeHigh=timeHigh+" 23:59:59"
//		var tablefragment []model.Tablefragment
//		var city1 = "fragment" + city
//		//fmt.Println(averagelow)
//		//fmt.Println(averagehigh)
//		if model3 == "0" {
//			if timelength == "1" {
//
//				common.IndexDB.Table(city1).Where("dev_id=?   AND start_time BETWEEN  ? AND ? AND " +
//					"water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id=?  AND start_time BETWEEN  ? AND ? " +
//					"water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id=? AND start_time BETWEEN  ? AND ?" +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id=?  AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//				common.IndexDB.Table(city1).Where("dev_id=? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//
//		}else  if model3=="5"{
//			if timelength == "1" {
//				common.IndexDB.Table(city1).Where("dev_id=?   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
//					" water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//
//
//		}else  if model3=="1"{
//			if timelength == "1" {
//				common.IndexDB.Table(city1).Where("dev_id=?   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
//					"  water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id=?  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//
//
//		} else {
//			if timelength == "1" {
//
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? " +
//					"AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//
//				common.IndexDB.Table(city1).Where("dev_id=? AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//		}
//
//
//		//	common.IndexDB.Table(city1).Where("dev_id=? AND pattern=? AND start_time BETWEEN  ? AND ?",dev_id,model3,timeLow,timeHigh).Find(&tablefragment)
//		//for _, tableDate := range tablefragment {
//		//
//		//	datatime = append(datatime, tableDate.Multiple_score) //数据放入切片中
//		//
//		//
//		//}
//		//fmt.Println(datatime)
//		//	response.Success(ctx,gin.H{"data":tablefragment},"成功")})=
//		response.Success(ctx, gin.H{"data": tablefragment}, "成功")
//
//	}else if flag == "3" {
//		dev_type:=ctx.DefaultQuery("dev_type","0000")
//		timeLow=timeLow+" 00:00:00"
//		timeHigh=timeHigh+" 23:59:59"
//		var tablefragment []model.Tablefragment
//		var city1 = "fragment" + city
//		var Dev_id [] string
//		//fmt.Println(averagelow)
//		//fmt.Println(averagehigh)
//
//		t.DB.Table("bo").Where("dev_type = ? AND city_code = ? ",dev_type,city).Find(&tableStoreDates4)
//		for _,tableDate :=range tableStoreDates4{
//			Dev_id=append(Dev_id,tableDate.Dev_Id)
//		}
//		if model3 == "0" {
//			if timelength == "1" {
//
//				common.IndexDB.Table(city1).Where("dev_id in (?)   AND start_time BETWEEN  ? AND ? AND " +
//					"water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)  AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND start_time BETWEEN  ? AND ?" +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)  AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//
//		}else  if model3=="5"{
//			if timelength == "1" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
//					" water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, "7","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//
//
//		}else  if model3=="1"{
//			if timelength == "1" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)   AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? AND" +
//					" water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ?" +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id in (?)  AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern BETWEEN  ? AND ? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, "6","8",timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//
//
//		} else {
//			if timelength == "1" {
//
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			} else if timelength == "2" {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "3" {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? " +
//					"AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else if timelength == "4" {
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					" AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//				//fmt.Println(tablefragment)
//			} else {
//
//				common.IndexDB.Table(city1).Where("dev_id in (?) AND water_pattern=? AND start_time BETWEEN  ? AND ? " +
//					"AND water_score  BETWEEN  ? AND ?", Dev_id, model3, timeLow, timeHigh,scorelow,scorehigh).Find(&tablefragment)
//			}
//		}
//
//
//		//	common.IndexDB.Table(city1).Where("dev_id=? AND pattern=? AND start_time BETWEEN  ? AND ?",dev_id,model3,timeLow,timeHigh).Find(&tablefragment)
//		//for _, tableDate := range tablefragment {
//		//
//		//	datatime = append(datatime, tableDate.Multiple_score) //数据放入切片中
//		//
//		//
//		//}
//		//fmt.Println(datatime)
//		//	response.Success(ctx,gin.H{"data":tablefragment},"成功")})=
//		if len(tablefragment)>1000&&flag1=="0"{
//			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"数据量超过1000条(建议具体查询条件)"})
//			return
//		}else{
//			response.Success(ctx, gin.H{"data": tablefragment}, "成功")
//		}
//
//	}}else
//if pageType == "indexbehavior" {
//
//	//fmt.Println(scorehigh)
//	if flag == "1" {
//		var tablename string
//		tablename = "data" + city+"_"+timeLow[0:4]+timeLow[5:7]+timeLow[8:10]
//
//		common.RunDB.Table(tablename).Where(" applianceid=? AND datatime BETWEEN ? AND ?", dev_id,timeLow,timeHigh).Find(&tableStoreDates)
//		//	fmt.Println(tableStoreDates)
//
//		for _, tableDate := range tableStoreDates {
//
//			datatime = append(datatime, tableDate.Datatime) //数据放入切片中
//			flow = append(flow, tableDate.Flow)
//			flame = append(flame, tableDate.Flame)
//			settemp = append(settemp, tableDate.Settemp)
//			outtemp = append(outtemp, tableDate.Outtemp)
//			model1 = append(model1, tableDate.Water_pattern)
//			zone_id = append(zone_id, tableDate.Zone_id)
//
//		}
//		//	response.Success(ctx,gin.H{"data":tableStoreDates,"stable_proportion":Stable_proportion,"un_stable_proportion":Un_stable_proportion,"province":province},"成功")
//		// fmt.Println(len(datatime))
//		response.Success(ctx, gin.H{"data": tableStoreDates, "data_time": datatime, "flow": flow, "flame": flame, "set_temp": settemp, "out_temp": outtemp, "model": model1,
//			"zone_id": zone_id, "data1": tableplace, "total_time": Total_time, "min_time": Min_time, "max_time": Max_time, "total_num": Total_num,
//			"datanum": tableStoreDates3}, "成功")
//
//	} else if flag == "2" {
//		timeLow=timeLow+" 00:00:00"
//		timeHigh=timeHigh+" 23:59:59"
//		var tablefragment []model.Tablewaterbehavior
//		//fmt.Println(averagelow)
//		//fmt.Println(averagehigh)
//		if model3 == "0" {
//			common.IndexDB.Table("mode_behaviors").Where("dev_id=?   AND start_time BETWEEN  ? AND ? ", dev_id, timeLow, timeHigh).Find(&tablefragment)
//		}else  if model3=="1"{
//			common.IndexDB.Table("mode_behaviors").Where("dev_id=?  AND start_time BETWEEN  ? AND ? AND effect_flag = ?", dev_id, timeLow, timeHigh,"1").Find(&tablefragment)
//		}else  if model3=="2"{
//			common.IndexDB.Table("mode_behaviors").Where("dev_id=?  AND start_time BETWEEN  ? AND ? AND effect_flag = ?", dev_id, timeLow, timeHigh,"0").Find(&tablefragment)
//		}
//		//	common.IndexDB.Table(city1).Where("dev_id=? AND pattern=? AND start_time BETWEEN  ? AND ?",dev_id,model3,timeLow,timeHigh).Find(&tablefragment)
//		//for _, tableDate := range tablefragment {
//		//
//		//	datatime = append(datatime, tableDate.Multiple_score) //数据放入切片中
//		//
//		//
//		//}
//		//	response.Success(ctx,gin.H{"data":tablefragment},"成功")})=
//		response.Success(ctx, gin.H{"data": tablefragment}, "成功")
//
//	}} else
//if pageType == "indexmodeltype" {
//	var filter[]string
//	time := make(map[string] []string)
//	xiaoshu := make(map[string] []float32)
//	zhengshu := make(map[string] []int)
//	if flag == "1" {
//		if category=="0000"{
//			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? ", timeLow,timeHigh, city).Find(&tableplace)
//		}else{
//			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow,timeHigh, city,category).Find(&tableplace)
//		}
//
//	} else {
//		if category=="0000"{
//			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND dev_type=? ",timeLow,timeHigh, dev_id).Find(&tableplace)
//		}else{
//			common.IndexDB.Table("days_summaries").Where("time_date BETWEEN  ? AND ? AND city_code=? AND dev_type=?", timeLow,timeHigh, category,dev_id).Find(&tableplace)
//		}
//
//	}
//	//fmt.Println(tableplace)
//	for _, tableDate := range tableplace {
//		var flag bool=false
//		for  i:=0;i< len(filter);i++{
//			if(tableDate.Dev_Id==filter[i]){
//				flag=true
//			}
//		}
//		if flag==false{
//			filter = append(filter, tableDate.Dev_Id)
//			Dev_type = append(Dev_type, tableDate.Dev_type)
//			city_code = append(city_code, tableDate.City_code)
//			t.DB.Table("midea_loc_code").Where("city_code=?", city_code[len(city_code)-1]).Find(&tableStoreDates4)
//			for _, tableDate1 := range tableStoreDates4 {
//				province = append(province, tableDate1.Dev_city)
//			}
//			time[tableDate.Dev_Id]=[]string{"0","0","0"}
//			xiaoshu[tableDate.Dev_Id]=[]float32{0.0,0.0}
//			zhengshu[tableDate.Dev_Id]=[]int{0,0,0}
//		}
//		time0,_:=strconv.Atoi(time[tableDate.Dev_Id][0])
//		time1,_:=strconv.Atoi(time[tableDate.Dev_Id][1])
//		time2,_:=strconv.Atoi(time[tableDate.Dev_Id][2])
//		time[tableDate.Dev_Id][0]=strconv.Itoa(time2sec(tableDate.Water_valid_time)+time0)
//		time[tableDate.Dev_Id][1]=strconv.Itoa(time2sec(tableDate.Average_time)+time1)
//		time[tableDate.Dev_Id][2]=strconv.Itoa(time2sec(tableDate.Maximum_time)+time2)
//		xiaoshu[tableDate.Dev_Id][0]=tableDate.Un_stable_proportion+xiaoshu[tableDate.Dev_Id][0]
//		xiaoshu[tableDate.Dev_Id][1]=xiaoshu[tableDate.Dev_Id][1]+1
//		zhengshu[tableDate.Dev_Id][0]=tableDate.Water_score+zhengshu[tableDate.Dev_Id][0]
//		zhengshu[tableDate.Dev_Id][1]=tableDate.Water_num+zhengshu[tableDate.Dev_Id][1]
//		zhengshu[tableDate.Dev_Id][2]=zhengshu[tableDate.Dev_Id][2]+1
//		//Dev_type = append(Dev_type, tableDate.Dev_type)                            //数据放入切片中
//		//Stable_proportion = append(Stable_proportion, tableDate.Stable_proportion) //数据放入切片中
//		//Un_stable_proportion = append(Un_stable_proportion, tableDate.Un_stable_proportion)
//		//city_code = append(city_code, tableDate.City_code)
//		//t.DB.Table("midea_loc_code").Where("city_code=?", city_code[len(city_code)-1]).Find(&tableStoreDates4)
//		//for _, tableDate1 := range tableStoreDates4 {
//		//	province = append(province, tableDate1.Dev_city)
//		//}
//
//	}
//	for key := range time {
//		Time0,_:=strconv.Atoi(time[key][0])
//		Time1,_:=strconv.Atoi(time[key][1])
//		Time2,_:=strconv.Atoi(time[key][2])
//		time[key][0]=strconv.Itoa(Time0/zhengshu[key][2])
//		time[key][1]=strconv.Itoa(Time1/zhengshu[key][2])
//		time[key][2]=strconv.Itoa(Time2/zhengshu[key][2])
//		xiaoshu[key][0]=xiaoshu[key][0]/xiaoshu[key][1]
//		zhengshu[key][0]=zhengshu[key][0]/zhengshu[key][2]
//		zhengshu[key][1]=zhengshu[key][1]/zhengshu[key][2]
//	}
//
//	fmt.Println(Dev_type)
//	fmt.Println(filter)
//	fmt.Println(xiaoshu)
//	response.Success(ctx, gin.H{"data": tableplace, "stable_proportion": Stable_proportion, "un_stable_proportion": Un_stable_proportion, "province": province, "dev_type": Dev_type,"filter":filter,"time":time,"xiaoshu":xiaoshu,"zhengshu":zhengshu}, "成功")
//
//}else
//if pageType =="menuhome"{
//	var average_time []int
//	var dev_id []string
//	if flag =="1"{
//		common.IndexDB.Table("month_summaries").Where("time_date =? AND province_code =?",month1,city).Find(&tableplace)
//		for _,tableDate1 :=range tableplace{
//			average_time=append(average_time,time2sec(tableDate1.Average_time))
//			dev_id=append(dev_id,tableDate1.Dev_Id)
//		}
//	}else{
//		//	fmt.Println(timeLow)
//		common.IndexDB.Table("province_months").Where("time_date =?",month1).Find(&tableplace)
//		//    fmt.Println(tableplace)
//		for _,tableDate1 :=range tableplace{
//			province_code=append(province_code,tableDate1.Province_code)
//			t.DB.Table("midea_loc_code").Where("province_code =?",province_code[len(province_code)-1]).Find(&tableStoreDates4)
//			province=append(province,tableStoreDates4[0].Dev_province)
//		}
//	}
//	response.Success(ctx,gin.H{"data":tableplace,"province":province,"average_time":average_time,"dev_id":dev_id},"成功")
//} else
//if pageType == "menu" {
//	var province1 []string
//
//
//	var tabletimerate1 [] model.Tabletimerate
//	var fragmentflow    []model.Fragmentflow
//	t.DB.Table("midea_loc_code").Find(&tableStoreDates4)
//
//	//循环省份表
//	for _, tableDate := range tableStoreDates4 {
//		city_code=append(city_code,tableDate.City_code)
//
//		province=append(province,tableDate.Dev_province)
//
//		var flag_1 bool=false
//		for  i:=0;i< len(province1);i++{
//			if(province[len(province)-1]==province1[i]){
//				flag_1=true
//			}
//		}
//		//若省份改变，进行该省份水流量片段分布查询
//		if flag_1==false{
//			province1 = append(province1, province[len(province)-1])
//			province_code=append(province_code,tableDate.Province_code)
//			common.IndexDB.Table("province_days").Where("province_code= ? AND time_date BETWEEN ? AND ?", province_code[len(province_code)-1], timeLow, timeHigh).Find(&fragmentflow)
//
//			var three int = 0
//			var five int = 0
//			var ten int = 0
//			var equipment_num int=0
//			var tabletimerate model.Tabletimerate
//			tabletimerate.Province=province1[len(province1)-1]
//			for _, tableDate1 := range fragmentflow {
//				three=three+tableDate1.Three
//				five=five+tableDate1.Five
//				ten=ten+tableDate1.Ten
//				//fmt.Print(ten)
//				equipment_num=equipment_num+tableDate1.Equipment_num
//
//
//			}
//			var three1 float64 = 0.0
//			var five1 float64 = 0.0
//			var ten1 float64 = 0.0
//			three1=float64(three)/float64(equipment_num)
//			five1=float64(five)/float64(equipment_num)
//			ten1=float64(ten)/float64(equipment_num)
//			//if equipment_num!=0{
//			//	three=three/equipment_num
//			//	five=five/equipment_num
//			//	ten=ten/equipment_num
//			//
//			//}
//
//
//			if three!=0&&five!=0&&ten!=0{
//				dataparameter1 := make(map[string]float64)
//				dataparameter1["three"]=three1
//				dataparameter1["five"]=five1
//				dataparameter1["ten"]=ten1
//				tabletimerate.Duration_time=dataparameter1
//				//fmt.Print(tabletimerate)
//				tabletimerate.Equipment_num=equipment_num
//				tabletimerate1=append(tabletimerate1,tabletimerate)
//			}
//
//
//		}
//
//	}
//	//查询全国，汇总计算
//	if len(tabletimerate1)!=0{
//		var three float64 = 0
//		var five float64 = 0
//		var ten float64 = 0
//		var equipment_num int=0
//		for i:=0;i<len(tabletimerate1);i++{
//			three=three+tabletimerate1[i].Duration_time["three"]
//			five=five+tabletimerate1[i].Duration_time["five"]
//			ten=ten+tabletimerate1[i].Duration_time["ten"]
//			equipment_num=equipment_num+tabletimerate1[i].Equipment_num
//		}
//		var tabletimerate model.Tabletimerate
//		tabletimerate.Province="全国"
//		dataparameter1 := make(map[string]float64)
//		dataparameter1["three"]=three
//		dataparameter1["five"]=five
//		dataparameter1["ten"]=ten
//		tabletimerate.Duration_time=dataparameter1
//		tabletimerate.Equipment_num=equipment_num
//		tabletimerate1=append(tabletimerate1,tabletimerate)
//	}
//
//
//
//	response.Success(ctx, gin.H{"data1": tabletimerate1}, "成功")
//}else if pageType == "newmenu" {
//	var behavior [] model.Behavior
//	var Tablebehavior1 [] model.Tablebehavior
//	var id [] string
//	SheBei := make(map[string] []float32)
//	common.IndexDB.Table("behavior_summaries").Where("data_time BETWEEN ? AND ?",timeLow,timeHigh).Find(&behavior)
//	//SheBei[behavior[0].Dev_Id]=[]float32{0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0}
//	//id[0]=behavior[0].Dev_Id
//	for _, tableDate := range behavior {
//		var Tablebehavior2 model.Tablebehavior
//
//		var flag_1 bool=false
//		for i:=0;i< len(id);i++{
//			if id[i]==tableDate.Dev_Id{
//				SheBei[tableDate.Dev_Id][0]=SheBei[tableDate.Dev_Id][0]+tableDate.Sec0p
//				SheBei[tableDate.Dev_Id][1]=SheBei[tableDate.Dev_Id][1]+tableDate.Sec30p
//				SheBei[tableDate.Dev_Id][2]=SheBei[tableDate.Dev_Id][2]+tableDate.Min3p
//				SheBei[tableDate.Dev_Id][3]=SheBei[tableDate.Dev_Id][3]+tableDate.Min10p
//				flag_1=true
//			}
//		}
//		if flag_1==false{
//			id=append(id,tableDate.Dev_Id)
//			SheBei[tableDate.Dev_Id]=[]float32{0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0}
//			SheBei[tableDate.Dev_Id][0]=tableDate.Sec0p
//			SheBei[tableDate.Dev_Id][1]=tableDate.Sec30p
//			SheBei[tableDate.Dev_Id][2]=tableDate.Min3p
//			SheBei[tableDate.Dev_Id][3]=tableDate.Min10p
//		}
//		if tableDate.Sec0p==0&&tableDate.Sec30p==0&&tableDate.Min3p==0&&tableDate.Min10p==0{
//		}else{
//			SheBei[tableDate.Dev_Id][4]=SheBei[tableDate.Dev_Id][4]+1
//		}
//		SheBei[tableDate.Dev_Id][5]=SheBei[tableDate.Dev_Id][0]/SheBei[tableDate.Dev_Id][4]
//		SheBei[tableDate.Dev_Id][6]=SheBei[tableDate.Dev_Id][1]/SheBei[tableDate.Dev_Id][4]
//		SheBei[tableDate.Dev_Id][7]=SheBei[tableDate.Dev_Id][2]/SheBei[tableDate.Dev_Id][4]
//		SheBei[tableDate.Dev_Id][8]=SheBei[tableDate.Dev_Id][3]/SheBei[tableDate.Dev_Id][4]
//		if flag_1==false{
//			Tablebehavior2.Dev_id=tableDate.Dev_Id
//			t.DB.Table("bo").Where("dev_id = ? ",tableDate.Dev_Id).Find(&tableStoreDates4)
//			Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
//			t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
//			Tablebehavior2.City=tableStoreDates4[0].Dev_city
//			Tablebehavior2.Duration_time=SheBei[tableDate.Dev_Id]
//			Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
//		}else{
//			for i:=0;i< len(Tablebehavior1);i++{
//				if Tablebehavior1[i].Dev_id==tableDate.Dev_Id{
//					Tablebehavior1[i].Duration_time=SheBei[tableDate.Dev_Id]
//				}
//			}
//		}
//
//		//if id!=tableDate.Dev_Id{
//		//	var Tablebehavior2 model.Tablebehavior
//		//	Tablebehavior2.Dev_id=id
//		//	t.DB.Table("id1").Where("dev_id = ? ",id).Find(&tableStoreDates4)
//		//	Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
//		//	t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
//		//	Tablebehavior2.City=tableStoreDates4[0].Dev_city
//		//	SheBei[id][0]=SheBei[id][0]/SheBei[id][4]
//		//	SheBei[id][1]=SheBei[id][1]/SheBei[id][4]
//		//	SheBei[id][2]=SheBei[id][2]/SheBei[id][4]
//		//	SheBei[id][3]=SheBei[id][3]/SheBei[id][4]
//		//	Tablebehavior2.Duration_time=SheBei[id]
//		//	Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
//		//	id=tableDate.Dev_Id
//		//	SheBei[tableDate.Dev_Id]=[]float32{0.0,0.0,0.0,0.0,0.0}
//		//	i++
//		//}
//		//SheBei[tableDate.Dev_Id][0]=SheBei[tableDate.Dev_Id][0]+tableDate.Sec0p
//		//SheBei[tableDate.Dev_Id][1]=SheBei[tableDate.Dev_Id][1]+tableDate.Sec30p
//		//SheBei[tableDate.Dev_Id][2]=SheBei[tableDate.Dev_Id][2]+tableDate.Min3p
//		//SheBei[tableDate.Dev_Id][3]=SheBei[tableDate.Dev_Id][3]+tableDate.Min10p
//		//if tableDate.Sec0p==0&&tableDate.Sec30p==0&&tableDate.Min3p==0&&tableDate.Min10p==0{
//		//}else{
//		//	SheBei[tableDate.Dev_Id][4]=SheBei[tableDate.Dev_Id][4]+1
//		//}
//		//fmt.Print(Tablebehavior1)
//	}
//	fmt.Print(Tablebehavior1)
//	//for key,content:= range SheBei{
//	//	fmt.Print(key,content)
//	//	var Tablebehavior2 model.Tablebehavior
//	//	Tablebehavior2.Dev_id=key
//	//	t.DB.Table("id1").Where("dev_id = ? ",key).Find(&tableStoreDates4)
//	//	Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
//	//	t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
//	//	Tablebehavior2.City=tableStoreDates4[0].Dev_city
//	//	Tablebehavior2.Duration_time=content
//	//	Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
//	//}
//	//var Tablebehavior2 model.Tablebehavior
//	//Tablebehavior2.Dev_id=id
//	//t.DB.Table("id1").Where("dev_id = ? ",id).Find(&tableStoreDates4)
//	//Tablebehavior2.Dev_type=tableStoreDates4[0].Dev_type
//	//t.DB.Table("midea_loc_code").Where("city_code =?",tableStoreDates4[0].City_code).Find(&tableStoreDates4)
//	//Tablebehavior2.City=tableStoreDates4[0].Dev_city
//	//SheBei[id][0]=SheBei[id][0]/SheBei[id][4]
//	//SheBei[id][1]=SheBei[id][1]/SheBei[id][4]
//	//SheBei[id][2]=SheBei[id][2]/SheBei[id][4]
//	//SheBei[id][3]=SheBei[id][3]/SheBei[id][4]
//	//Tablebehavior2.Duration_time=SheBei[id]
//	//Tablebehavior1=append(Tablebehavior1,Tablebehavior2)
//	//fmt.Print(Tablebehavior1)
//	response.Success(ctx, gin.H{"SheBei": Tablebehavior1}, "成功")
//}

func NewTimeMonitorController() ITimeMonitorController {
	db := common.GetDB()
	db.AutoMigrate(model.TableDate{})
	return TimeMonitorController{DB: db}
}

//func time2sec1(time string) (second int) {
//	//var timeArr = strings.Split(time,":")
//	//var hour = timeArr[0]
//	//var minute = timeArr[1]
//	//var sec = timeArr[2]
//	//hour1, _ := strconv.Atoi(hour)
//	//minute1,_:=strconv.Atoi(minute)
//	//sec1 ,_ :=strconv.Atoi(sec)
//	//return hour1 * 3600 + minute1 * 60 + sec1
//	hour:="";
//	min:="";
//	sec:="";
//	j:=0;
//	k:=0;
//	// hour=time.split('h')[0];
//	// min=time.split('h')[1].split('m')[0];
//	// sec=time.split('h')[1].split('m')[1].split('s')[0];
//	for i:=0;i<len(time);i++{
//		if time[i:i+1]=="h"{
//			hour=time[0:i]
//			j=i+1;
//		}
//		if time[i:i+1]=="m"{
//			if j!=0{
//				min=time[j:i]
//			}else{
//				min=time[0:i]
//			}
//			k=i+1;
//		}
//		if time[i:i+1]=="s"{
//			if k!=0{
//				sec=time[k:i]
//			}else{
//				sec=time[0:i]
//			}
//			k=i+1;
//		}
//	}
//	hour1,_:=strconv.Atoi(hour)
//	min1,_:=strconv.Atoi(min)
//	sec1,_:=strconv.Atoi(sec)
//	second=hour1*3600+min1*60+sec1;
//	// console.log(sec)
//	return
//}

//指标明细计算//个数,总时长,最大时间,最短时间,最大阶跃幅值,时间明细,水流量极差明细,水流量幅值明细,水流量均值明细,水流量标准差明细
func feature1(u []model.TableDate, model int, gapTime time.Duration, zoneTime time.Duration) (number int, sum time.Duration, max time.Duration, min time.Duration, maxChange float64, timeLen []time.Duration, flowExtreme []int, flowMaxChange []float64, flowAvg []float64, flowDeviation []float64) {
	var length int
	var addT time.Duration
	var datatime []string //时间轴
	for _, u := range u {
		datatime = append(datatime, u.Datatime) //数据放入切片中
	}
	for i, n := 0, len(u)-1; i < n; i++ {
		length++
		t2, _ := time.Parse("2006-01-02 15:04:05", datatime[i+1])
		t1, _ := time.Parse("2006-01-02 15:04:05", datatime[i])
		deltaT := t2.Sub(t1)
		if deltaT > gapTime {
			if addT < zoneTime {
				//过滤无效区间
			} else {
				//fmt.Printf("有效区间--长度:%v,结束id:%v,节点时间差值:%v，区间总时间:%v\n", length, i+1, deltaT, addT)
				timeLen = append(timeLen, featureTime(u[i-length+1:i+1], model)...) //有效区间进入时间特征统计函数,左闭右开;将结果存放进一个切片u1中
				extreme, maxChange, avg, deviation := featureFlow(u[i-length+1:i+1], model)
				flowExtreme = append(flowExtreme, extreme...)
				flowMaxChange = append(flowMaxChange, maxChange...)
				flowAvg = append(flowAvg, avg...)
				flowDeviation = append(flowDeviation, deviation...)
			}
			length = 0
			addT = 0
		} else if i == n-1 { //最后一个区间
			addT += deltaT //加上最后一个时间差
			if addT < zoneTime {
				//过滤无效区间
			} else {
				//fmt.Printf("最后有效区间--长度:%v,结束id:%v,节点时间差值:%v，区间总时间:%v\n", length+1, i+2, deltaT, addT)
				timeLen = append(timeLen, featureTime(u[i-length+1:i+2], model)...) //有效区间进入时间特征统计函数,左闭右开;将结果存放进一个切片u1中
				extreme, maxChange, avg, deviation := featureFlow(u[i-length+1:i+2], model)
				flowExtreme = append(flowExtreme, extreme...)
				flowMaxChange = append(flowMaxChange, maxChange...)
				flowAvg = append(flowAvg, avg...)
				flowDeviation = append(flowDeviation, deviation...)
			}
		} else {
			addT += deltaT
		}
	}
	if len(timeLen) != 0 { //防止空数组（无该模式的情况）
		number, sum, max, min = featureTimeCalculate(timeLen)
		maxChange = featureFlowCalculate(flowMaxChange)
		//fmt.Println(FlowMaxChange)
	}
	return
}

// featureTime 时间特征明细
func featureTime1(u []model.TableDate, model int) []time.Duration {
	var number = 0
	var timeLen []time.Duration
	for i, n := 0, len(u); i < n; i++ {
		if u[i].Water_pattern == model {
			number++
			if i == n-1 { //区间末尾模式片段
				t2, _ := time.Parse("2006-01-02 15:04:05", u[i].Datatime)
				t1, _ := time.Parse("2006-01-02 15:04:05", u[i-number+1].Datatime)
				timeLen = append(timeLen, t2.Sub(t1))
			} else {
			} //不做处理
		} else {
			if number != 0 {
				t2, _ := time.Parse("2006-01-02 15:04:05", u[i-1].Datatime)
				t1, _ := time.Parse("2006-01-02 15:04:05", u[i-number].Datatime)
				timeLen = append(timeLen, t2.Sub(t1))
			}
			number = 0
		}
	}
	return timeLen
}

// featureTimeCalculate 时间特征指标计算
func featureTimeCalculate1(timeLen []time.Duration) (number int, sum time.Duration, max time.Duration, min time.Duration) {
	max = timeLen[0]
	min = timeLen[0]
	for i, n := 0, len(timeLen); i < n; i++ {
		sum += timeLen[i]
		if timeLen[i] > max {
			max = timeLen[i]
		} else if timeLen[i] < min {
			min = timeLen[i]
		}
	}
	number = len(timeLen)
	return
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// featureFlow 水流量特征明细（模式片段极差、最大幅值（有符号）、标准差、平均值）
//func featureFlow1(u []model.TableDate, model int) (flowExtreme []int, flowMaxChange []float64, flowAvg []float64, flowDeviation []float64) {
//	var number = 0
//	var deltaflow []int
//	for i, n := 0, len(u)-1; i < n; i++ {
//		deltaflow = append(deltaflow, u[i+1].Flow-u[i].Flow)
//	}
//	for i, n := 0, len(u); i < n; i++ {
//		if u[i].Water_pattern == model {
//			number++
//			if i == n-1 { //区间末尾模式片段
//				max := u[i-number+1].Flow
//				min := u[i-number+1].Flow
//				var sumFlow = 0
//				var sumFF = 0.0
//				maxChange := math.Abs(float64(deltaflow[i-number+1]))
//				for i1, n1 := i-number+1, i+1; i1 < n1; i1++ {
//					sumFlow += u[i1].Flow
//					if u[i1].Flow > max {
//						max = u[i1].Flow
//					} else if u[i1].Flow < min {
//						min = u[i1].Flow
//					}
//				}
//				for i1, n1 := i-number+1, i+1; i1 < n1; i1++ {
//					sumFF += (float64(sumFlow)/float64(number) - float64(u[i].Flow)) * (float64(sumFlow)/float64(number) - float64(u[i].Flow))
//				}
//				flowExtreme = append(flowExtreme, max-min)
//				avg, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sumFlow)/float64(number)), 64)
//				flowAvg = append(flowAvg, avg)
//				deviation, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Sqrt(sumFF/float64(number))), 64)
//				flowDeviation = append(flowDeviation, deviation)
//				for i1, n1 := i-number+1, i; i1 < n1; i1++ { //水流量差值
//					if math.Abs(float64(deltaflow[i1])) > math.Abs(maxChange) {
//						maxChange = float64(deltaflow[i1])
//					}
//				}
//				flowMaxChange = append(flowMaxChange, maxChange)
//			} else {
//			}
//		} else {
//			if number != 0 {
//				max := u[i-number].Flow
//				min := u[i-number].Flow
//				var sumFlow int
//				var sumFF float64
//				maxChange := math.Abs(float64(deltaflow[i-number]))
//				for i1, n1 := i-number, i; i1 < n1; i1++ { //水流量
//					sumFlow += u[i1].Flow
//					if u[i1].Flow > max {
//						max = u[i1].Flow
//					} else if u[i1].Flow < min {
//						min = u[i1].Flow
//					}
//				}
//				for i1, n1 := i-number, i; i1 < n1; i1++ {
//					sumFF += (float64(sumFlow)/float64(number) - float64(u[i].Flow)) * (float64(sumFlow)/float64(number) - float64(u[i].Flow))
//				}
//				flowExtreme = append(flowExtreme, max-min)
//				avg, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(sumFlow)/float64(number)), 64)
//				flowAvg = append(flowAvg, avg)
//				deviation, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Sqrt(sumFF/float64(number))), 64)
//				flowDeviation = append(flowDeviation, deviation)
//
//				for i1, n1 := i-number, i-1; i1 < n1; i1++ { //水流量差值
//					if math.Abs(float64(deltaflow[i1])) > math.Abs(maxChange) {
//						maxChange = float64(deltaflow[i1])
//					}
//				}
//				flowMaxChange = append(flowMaxChange, maxChange)
//			}
//			number = 0
//		}
//	}
//
//	return
//}
//
//// featureFlowCalculate 求水流量幅值最大模（有符号）
//func featureFlowCalculate1(u []float64) float64 {
//
//	var max = math.Abs(u[0])
//	for i, n := 0, len(u); i < n; i++ { //水流量差值
//		if math.Abs(u[i]) >= math.Abs(max) {
//			max = u[i]
//		}
//	}
//	return max
//
//}
