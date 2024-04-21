package math_3d

import (
	"fmt"
	"testing"
)

func TestMath(t *testing.T) {
	fmt.Println(VectorDot(&Vector{X: 1, Y: 1, Z: 0}, &Vector{X: 2, Y: 0, Z: 1}))
	fmt.Println(VectorAngle(&Vector{X: 1, Y: 1, Z: 0}, &Vector{X: 1, Y: -1, Z: 0}))
	fmt.Println(VectorCross(&Vector{X: 1, Y: 1, Z: 0}, &Vector{X: 2, Y: 0, Z: 1}))

	m1 := &Matrix{}
	m1.Init([]float64{1, 0, 1}, []float64{2, 0, 2})
	m2 := &Matrix{}
	m2.Init([]float64{1, 0, 1}, []float64{2, 0, 2})
	m3 := &Matrix{}
	m3.Init([]float64{1, 0, 0, 0}, []float64{0, 1, 0, 0}, []float64{0, 0, 1, 0}, []float64{1, 2, 3, 1})
	m4 := &Matrix{}
	m4.Init([]float64{1, 0, 0}, []float64{0, 1, 0}, []float64{0, 0, 1})
	fmt.Println(MatrixAdd(m1, m2))
	fmt.Println(MatrixMulNum(m1, 3))
	//fmt.Println(MatrixMulMatrix(m1, m2)) // panic 这两个矩阵不能乘
	//fmt.Println(MatrixMulMatrix(m1, m3))
	fmt.Println(m3.Inverse())
	fmt.Println(m4.IsOrthogonal())

	v := &Vector{3, 3, 4}
	fmt.Println(MatrixMulMatrix(v.ToMatrix(), TranslationMatrix(&Vector{1, 2, 3}))) // 4,5,7
}
