package math_3d

import "math"

type Quaternion struct {
	v *Vector // 虚部
	s float64 // 实部
	// q = ix+jy+kz+w, i^2=j^2=k^2=-1, jk=-kj=i, ki=-ik=j, ij=-ji=k
}

// 共轭四元数，q*
func (q *Quaternion) Conjugate() *Quaternion {
	return &Quaternion{
		v: q.v.MulNum(-1),
		s: q.s,
	}
}

// 四元数的模，size
func (q *Quaternion) Len() float64 {
	return math.Sqrt(q.v.X*q.v.X + q.v.Y*q.v.Y + q.v.Z*q.v.Z + q.s*q.s)
}

// 四元数的逆，q^-1
func (q *Quaternion) Inverse() *Quaternion {
	return &Quaternion{
		v: q.Conjugate().v.MulNum(1 / q.Len()),
		s: q.Conjugate().s / q.Len(),
	}
}

func QuaternionsAdd(qs ...*Quaternion) *Quaternion {
	ret := &Quaternion{}
	for _, q := range qs {
		ret.v = VectorAdd(ret.v, q.v)
		ret.s += q.s
	}
	return ret
}

// 四元数乘法（不满足交换律）
func QuaternionsMul(p, q *Quaternion) *Quaternion {
	return &Quaternion{
		v: VectorAdd(q.v.MulNum(p.s), p.v.MulNum(q.s), VectorCross(p.v, q.v)),
		s: p.s*q.s - VectorDot(p.v, q.v),
	}
}

// --------------------

// 3D旋转相关

// 根据旋转绕的轴和旋转角生成四元数
func GenQuaternionByAxisAndRad(axis *Vector, rad float64) *Quaternion {
	// q是单位四元数，q^-1 = q*
	return &Quaternion{
		v: axis.Normalize().MulNum(math.Sin(rad / 2)),
		s: math.Cos(rad / 2),
	}
}

func (v *Vector) Rotation3D(qs ...*Quaternion) *Vector {
	if len(qs) == 0 {
		return v
	}
	res := qs[len(qs)-1]
	for i := len(qs) - 2; i >= 0; i-- {
		res = QuaternionsMul(res, qs[i])
	}
	res = QuaternionsMul(res, &Quaternion{
		v: v,
		s: 0,
	})
	for i := 0; i < len(qs); i++ {
		res = QuaternionsMul(res, qs[i].Inverse())
	}
	return &Vector{
		X: res.v.X,
		Y: res.v.Y,
		Z: res.v.Z,
	}
}
