package controller

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)
//类中的常量
const ( //常量的定义
	MaxDelay = float64(10)
	g = 300
	ORDER = uint64(2)
	Size = int(30)
	CODEL = 3  //可以不是常量
	F = float64(0.95)
	CR = float64(0.6)
)
//整个函数的类
type GeneticAlgorithm struct {
	DelegateFlag string
	Ts float64
	Uk []float64
	Yk []float64
	VariableDelayState []float64
	MaxCounts uint64
	CODEL int
	X [][]float64
	MaxX []float64
	MinX []float64
	Delta float64
	ampMax uint64
	ampMin uint64
	Index  int
	timeCounter int
	InputTime []float64//狗杂函数里初始化
	InputAmp []float64//构造函数里初始化
	randX rand.Rand
}

func NewGeneticAlgorithm(f func( G *GeneticAlgorithm))*GeneticAlgorithm  {
	G := new(GeneticAlgorithm)
	f(G)//重载构造函数1
	return G
}
func Start(codeLength int,MaxPar []float64,MinPar []float64,delegateFlag string)func(G *GeneticAlgorithm){//第一个构造函数
	return func(G *GeneticAlgorithm) {
		G.CODEL =codeLength
		G.MaxX = MaxPar
		G.MinX = MinPar
		G.DelegateFlag = delegateFlag
		G.Ts = 0.1
		G.InputTime = []float64{0, 12,20, 26, 32, 36, 43, 55,68 }
		G.InputAmp = []float64{0.618,0.25,0.618, 0.824, 0.618, 0.463, 0.927, 0.223,0.618}
		G.Index = 0
		G.timeCounter = 0
	}
}
/*=========================================================
 * 函数名称：Jvalue()
 * 功能描述: 代价函数----------值传递
 =========================================================*/
func (G GeneticAlgorithm)Jvalue(y []float64,yp []float64)float64{
	J := 0.0
	for i := 0; i < len(y); i++ {
		if i <19{
			J += 0.5*(y[i]-yp[i])*(y[i]-yp[i])
		}else{
			J += 5*(y[i]-yp[i])*(y[i]-yp[i])
		}

	}
	return J * 10000
}
/*=========================================================
* 函数名称： Lism()
* 功能描述:二阶传函加滞后模型，结果向量以引用的方式缓存在实参output中
=========================================================*/
func (G GeneticAlgorithm) Lism(sPar []float64,ut []float64,outPut []float64,length uint64)float64{
	a1 := sPar[0]*G.Ts*G.Ts
	a2 := 2*sPar[0]*G.Ts*G.Ts
	a3 := sPar[0]*G.Ts*G.Ts
	b1 := (2*sPar[1]+G.Ts)*(2*sPar[2]+G.Ts)
	b2 := (2*sPar[1]+G.Ts)*(G.Ts-2*sPar[2])+(2*sPar[2]+G.Ts)*(G.Ts-2*sPar[1])
	b3 := (G.Ts-2*sPar[1])*(G.Ts-2*sPar[2])
	y := 0.0
	ydelay := 0.0
	G.Uk[0] = 0
	G.Uk[1] = 0
	G.Yk[0] = 0
	G.Yk[1] = 0
	G.ClearVDS(0)
	for i := 0; i < int(length); i++ {
		y = (a1*float64(ut[i])+a2*G.Uk[1]+a3*G.Uk[0]-b2*G.Yk[1]-b3*G.Yk[0]) / b1
		ydelay = G.DelayModel(y,uint64(sPar[3]*10))
		outPut[i] = ydelay
		G.Yk[1] = y
		G.Uk[0] = G.Uk[1]
		G.Uk[1] = ut[i]
	}
	return y
}
/*=========================================================
 * 函数名称： DelayModel()
 * 功能描述:模拟系统滞后输出的模型
 =========================================================*/
func (G GeneticAlgorithm)DelayModel(inputSignal float64,delay uint64)float64{
	var idxDelay uint64
	var tmp uint64
	var outputSignal float64
	if delay <= uint64(0){
		outputSignal = inputSignal
	}else {
		if delay > G.MaxCounts{
			tmp = G.MaxCounts
		}else {
			tmp = delay
		}
		outputSignal = G.VariableDelayState[G.MaxCounts - tmp]
		for idxDelay = 0; idxDelay <G.MaxCounts - 1; idxDelay++{
			G.VariableDelayState[idxDelay] = G.VariableDelayState[idxDelay + 1]
		}
		G.VariableDelayState[G.MaxCounts - 1] = inputSignal
	}
	return outputSignal
}
/*=========================================================
 * 函数名称：FitFunc()
 * 功能描述:二阶传函加滞后模型，结果向量以引用的方式缓存在实参output中
 =========================================================*/
func (G GeneticAlgorithm)FitFunc(sPar []float64,ut []float64,outPut []float64,length uint64)[]float64{
	y := 0.0
	for i := 0; i < int(length) ; i++ {
		y = sPar[0]*math.Exp(-sPar[1]*ut[i])+sPar[2]
		outPut = append(outPut,y)
	}
	return outPut
}
/*=========================================================
 * 函数名称：ClearVDS()
 * 功能描述:重置模拟系统滞后的数组
 =========================================================*/
func (G GeneticAlgorithm)ClearVDS(resetValue float64)  {
	for i := 0; i <len(G.VariableDelayState) ; i++ {
		G.VariableDelayState[i] = resetValue
	}
}
/*=========================================================
 * 函数名称：ClearSys()
 * 功能描述:重置输入、输出向量缓存
 =========================================================*/
func (G GeneticAlgorithm)ClearSys(resetValue float64)  {
	if G.Uk != nil && G.Yk != nil && len(G.Uk) == int(ORDER) && len(G.Yk)== int(ORDER){
		for i := 0; i < int(ORDER) ; i++ {
			G.Uk[i] = resetValue
			G.Yk[i] = resetValue
		}
	}else {
		fmt.Println("错误，数组未分配空间")
	}
}
///*=========================================================
// * 测试函数
// * 函数名称：SeriesInit()
// =========================================================*/
func (G GeneticAlgorithm)SeriesInit(nbits uint64,ts float64,amMax uint64,amMin uint64)  {
	G.Delta = ts
	G.ampMax = amMax
	G.ampMin = amMin
	G.Index = 0
	G.timeCounter = 0
}
/*=========================================================
 * 测试函数
 * 函数名称：MSeries()
 =========================================================*/
func(G GeneticAlgorithm)MSeries()float64{
	u := 0.0
	t := float64(G.timeCounter) * G.Delta
	if t > G.InputTime[G.Index]{
		u = float64(G.InputTime[G.Index] * float64((G.ampMax - G.ampMin)+G.ampMin))
		if (t >= G.InputTime[G.Index + 1]) && (G.Index < (len(G.InputTime) -2)){
			G.Index +=1
		}else if G.Index >= len(G.InputTime) - 2{
			u = float64(0.2 * float64((G.ampMax - G.ampMin)+G.ampMin))
		}
	}
	if G.timeCounter > 700{
		u = -1
	}
	G.timeCounter++
	return u
}
/*=========================================================
 * 测试函数
 * 函数名称：MSeries()
 =========================================================*/
func (G GeneticAlgorithm)Test(){
	//i := 0//未使用
	//var  sPar [4]float64 = [4]float64{2,1,3,0}//未使用
	//u := make([]float64,1000)//未使用
	//y := make([]float64,1000)//未使用
	ut := make([]float64,640)//c#中浮点型数组
	for i := 0; i < 640 ; i++ {
		ut[i] = G.MSeries()
	}
	output := make([]float64,640)
	var sys []float64 = []float64{0.03,6,5,5}
	G.Lism(sys,ut,output,640)
	//res := []float64//未使用
	//遗传算法
}
/*=========================================================
 * 测试函数
 * y是温升 ut是控制量 length 参数向量的长度
 * 函数名称：genetic()
 =========================================================*/
func (G GeneticAlgorithm) genetic(y []float64, ut []float64, length uint64) (re []float64 , je []float64) {
	var BestS []float64 //最优个体的指针
	var yp []float64    //每次存储辨识出的系统参数对应的输出向量
	var JiD []float64   //进化过程中的代价函数
	var Ji float64      //最优个体的代价函数
	var res []float64 //辨识结果
	var nihedaijia []float64
	for i := 0; i < Size; i++ { //初始化种群{for (i:=0,j;i<len())//Golang中必须事先指定容器的大小
		var rowSlice []float64
		for j := 0; j < CODEL; j++ {
			var x float64
			rand.Seed(time.Now().UnixNano())
			a := rand.Float64()
			x= G.MinX[j] + (G.MaxX[j]-G.MinX[j])*a
			rowSlice = append(rowSlice,x)
		}
		G.X = append(G.X,rowSlice)
	}

	BestS = G.X[0]//初始化最优个体
	J1 := 0.0
	J2 := 0.0

	for i := 1; i < Size; i++ { //选择初始种群中的最优个体
		if G.DelegateFlag == "1" {
			G.Lism(BestS, ut, yp, length)
			J1 = G.Jvalue(y, yp) //计算最优个体的代价函数
			G.Lism(G.X[i], ut, yp, length)
			J2 = G.Jvalue(y, yp) //计算当前个体的代价函数
		} else if G.DelegateFlag == "0" {
			b:=G.FitFunc(BestS, ut, yp, length)
			J1 = G.Jvalue(y,b) //计算最优个体的代价函数
			c:=G.FitFunc(G.X[i], ut, yp, length)
			J2 = G.Jvalue(y,c) //计算当前个体的代价函数
		}

		if J1 > J2 {
			BestS = G.X[i]
			J1 = J2
			nihedaijia = append(nihedaijia,J2)//寻找最优个体，将过程的代价函数放到切片中
		}
	}
	Ji = J1 //将个体最优代价函数给Ji

	for i := 0; i < g; i++ { //主循环，进化G代
		// float[] vi;
		// float[] hi;
		//var vi []float64            //变异个体
		vi := make([]float64,3)
		var hi [CODEL]float64       //交叉后代=3
		for j := 0; j < Size; j++ { //变异
			var r [CODEL]int64//Golang中必须事先指定容器的大小
			//int[] r = new int[CODEL];
			for k := 0; i < CODEL; i++ {
				r[k] = 0
			}
			res2 := judge(r, int64(j))
			if res2 {//交叉

				for i1 := 0; i1<CODEL; i1++ {
					//rand.Seed(time.Now().UnixNano())
					//timeFlag := rand.Float64()
					r[i1] = int64(float64(Size-1)*rand.Float64())
				}
			}
			for k := 0; k < CODEL; k++ {
				hi[k] = BestS[k] + F*(G.X[r[0]][k]-G.X[r[1]][k])//交叉概率
				if hi[k] > G.MinX[k] {//约束调参范围
					hi[k] = hi[k]
				} else {
					hi[k] = G.MinX[k]
				}

				if hi[k] < G.MaxX[k] {
					hi[k] = hi[k]
				} else {
					hi[k] = G.MaxX[k]
				}
			}
			for k := 0; k < CODEL; k++ {
				if  rand.Float64() < CR {//变异概率0.6
					vi[k] = hi[k]
					//vi = append(vi,hi[k])
				} else {
					vi[k] = G.X[j][k] //不变异
					//vi = append(vi,G.X[j][k])
				}
			}
			if G.DelegateFlag == "1" { //选择淘汰
				G.Lism(vi, ut, yp, length)
				J1 = G.Jvalue(y, yp)
				G.Lism(G.X[j], ut, yp, length)
				J2 = G.Jvalue(y, yp)
			} else if G.DelegateFlag == "0" {
				e:=G.FitFunc(vi, ut, yp, length)
				J1 = G.Jvalue(y, e)
				d:=G.FitFunc(G.X[j], ut, yp, length)
				J2 = G.Jvalue(y, d)
			}
			if J1 < J2 {
				//var rowSlice []float64
				for k := 0; k < CODEL; k++ {
					G.X[j][k] = vi[k]//将变异良好的个体放进去或者说淘汰
					//rowSlice = append(rowSlice,vi[k])
				}
				//G.X = append(G.X,rowSlice)
			}
			var J0 float64
			if G.DelegateFlag == "1" { //选择淘汰
				G.Lism(G.X[j], ut, yp, length)
				J0 = G.Jvalue(y, yp)
			} else if G.DelegateFlag == "0" {
				f:=G.FitFunc(G.X[j], ut, yp, length)
				J0 = G.Jvalue(y, f)
			}
			if J0 < Ji {
				Ji = J0
				BestS = G.X[j]
			}
		}
		JiD = append(JiD,Ji)
	}
	//if JiD[len(JiD)-1] < 2{
	//	p := plot.New()
	//	p.Title.Text = "遗传代数-代价趋势图"
	//	p.X.Label.Text = "遗传代数"
	//	p.Y.Label.Text = "代价"
	//	plotutil.AddLinePoints(p,Plot(JiD))
	//	p.Save(4*vg.Inch, 4*vg.Inch, "points.png")
	//}
	fmt.Println("选择最优个体的每一代的代价：",nihedaijia)
//	fmt.Println("变异每一代的代价：",JiD)
	res = BestS
	return res,JiD
}

func judge(r [3]int64, index int64) (re bool) {
	 var res bool
	 res = false
	 for i := 0; i < len(r); i++ {
		 if r[i] == index {
			 res = true
			 }

		 }
	 for i := 0; i < CODEL-1; i++ {
		 for j := i + 1; j < CODEL; j++ {
			 if r[i] == r[j] {
				 res = true
				 }
			 }
		 }
		 return res
}
//func Plot(data []float64)plotter.XYs{
//	points := make(plotter.XYs,len(data))
//	for i := range points{
//		points[i].X = float64(i)
//		points[i].Y = data[i]
//	}
//	return points
//}




