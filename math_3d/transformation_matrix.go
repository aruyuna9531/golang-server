package math_3d

import (
	"fmt"
	"math"
)

// 以下矩阵都是基于3维向量的

// t的转换矩阵
// 1 0 0 0
// 0 1 0 0
// 0 0 1 0
// X Y Z 1
// 向量X={x,y,z,1}乘以该矩阵相当于X+t
// 该矩阵的逆矩阵就是最后一行前3个值取反其他不变
func TranslationMatrix(t *Vector) *Matrix {
	ret := &Matrix{}
	ret.Init([]float64{1, 0, 0, 0}, []float64{0, 1, 0, 0}, []float64{0, 0, 1, 0}, []float64{t.X, t.Y, t.Z, 1})
	return ret
}

type AxisType int

const (
	AxisTypeX AxisType = 0
	AxisTypeY AxisType = 1
	AxisTypeZ AxisType = 2
)

// 旋转矩阵，一个向量乘以这个矩阵的结果是这个向量绕axisType轴旋转rad弧度的结果
func RotationMatrix(axisType AxisType, rad float64) *Matrix {
	ret := &Matrix{}
	switch axisType {
	case AxisTypeX: // 绕X轴旋转，自Y正向向Z正向
		ret.Init([]float64{1, 0, 0, 0}, []float64{0, math.Cos(rad), math.Sin(rad), 0}, []float64{0, -math.Sin(rad), math.Cos(rad), 0}, []float64{0, 0, 0, 1})
	case AxisTypeY: // 绕Y轴旋转，自X正向向Z正向（时针方向和另外两个不一样）
		ret.Init([]float64{math.Cos(rad), 0, -math.Sin(rad), 0}, []float64{0, 1, 0, 0}, []float64{math.Sin(rad), 0, math.Cos(rad), 0}, []float64{0, 0, 0, 1})
	case AxisTypeZ: // 绕Z轴旋转，自X正向向Y正向
		ret.Init([]float64{math.Cos(rad), math.Sin(rad), 0, 0}, []float64{-math.Sin(rad), math.Cos(rad), 0, 0}, []float64{0, 0, 1, 0}, []float64{0, 0, 0, 1})
	default:
		panic(fmt.Sprintf("illegal axisType %d", axisType))
	}
	return ret
}

// 放缩矩阵，输入向量分别在X，Y，Z轴放大的倍数，向量乘以输出矩阵可得到结果
func ScaleMatrix(Xscale, Yscale, Zscale float64) *Matrix {
	ret := &Matrix{}
	ret.Init([]float64{Xscale, 0, 0, 0}, []float64{0, Yscale, 0, 0}, []float64{0, 0, Zscale, 0}, []float64{0, 0, 0, 1})
	return ret
}

// 转换顺序：Scale→Rotation→Translation（SRT）顺序调整会出现结果不满足期望的情况

// -------------------------

// RotationBy v向量绕任意axisVec向量旋转rad弧度得到的向量（复杂，建议引入四元数（quaternion））
func (v *Vector) RotationBy(axisVec *Vector, rad float64) *Vector {
	axisNor := axisVec.Normalize()
	m1 := &Matrix{}
	m1.Init([]float64{1, 0, 0}, []float64{0, 1, 0}, []float64{0, 0, 1})
	m2 := &Matrix{}
	m2.Init([]float64{0, -axisNor.Z, axisNor.Y}, []float64{axisNor.Z, 0, -axisNor.X}, []float64{-axisNor.Y, axisNor.X, 0})
	m3 := &Matrix{}
	m3.Init(
		[]float64{axisNor.X * axisNor.X, axisNor.X * axisNor.Y, axisNor.X * axisNor.Z},
		[]float64{axisNor.Y * axisNor.X, axisNor.Y * axisNor.Y, axisNor.Y * axisNor.Z},
		[]float64{axisNor.Z * axisNor.X, axisNor.Z * axisNor.Y, axisNor.Z * axisNor.Z},
	)
	rotationMatrix := MatrixAdd(MatrixMulNum(m1, math.Cos(rad)), MatrixMulNum(m2, math.Sin(rad)), MatrixMulNum(m3, 1-math.Cos(rad)))
	vv := &Matrix{}
	vv.Init([]float64{v.X, v.Y, v.Z})
	final := MatrixMulMatrix(vv, rotationMatrix)
	return &Vector{X: final.m[0][0], Y: final.m[0][1], Z: final.m[0][2]}
}
