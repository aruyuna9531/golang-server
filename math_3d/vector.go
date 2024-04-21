package math_3d

import "math"

// TODO 考虑向量长度为0的问题

type Vector struct {
	X, Y, Z float64
}

func (v *Vector) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector) Normalize() *Vector {
	l := v.Len()
	return &Vector{
		X: v.X / l,
		Y: v.Y / l,
		Z: v.Z / l,
	}
}

func (v *Vector) ToMatrix() *Matrix {
	ret := &Matrix{}
	ret.Init([]float64{v.X, v.Y, v.Z, 1}) // w = 1
	return ret
}

// 向量点积
func VectorDot(a, b *Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// 向量夹角
func VectorAngle(a, b *Vector) float64 {
	return math.Acos(VectorDot(a, b) / (a.Len() * b.Len()))
}

// 向量叉积
func VectorCross(a, b *Vector) *Vector {
	return &Vector{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}
