package math_3d

import "fmt"

func floatIsZero(f float64) bool {
	return f < 1e-9 && f > -1e-9
}

func floatEqual(a, b float64) bool {
	return floatIsZero(a - b)
}

type Matrix struct {
	rows    int
	columns int
	m       [][]float64
}

func (m *Matrix) Init(arr ...[]float64) {
	m.rows = len(arr)
	if m.rows == 0 {
		return
	}
	m.columns = len(arr[0])
	if m.columns == 0 {
		return
	}
	for i, v := range arr {
		if len(v) != m.columns {
			panic(fmt.Sprintf("illegal matrix init param of index %d length(%d) not equal to first index(%d)", i, len(v), m.columns))
		}
		m.m = append(m.m, v)
	}
}

// 逆矩阵
func (m *Matrix) Inverse() *Matrix {
	if m.rows != m.columns {
		panic(fmt.Sprintf("Inverse matrix of %d×%d not exist", m.rows, m.columns))
	}
	var N = m.rows
	W := make([][]float64, N)
	for i := 0; i < N; i++ {
		W[i] = make([]float64, 2*N)
	}
	result := make([][]float64, N)
	for i := 0; i < N; i++ {
		result[i] = make([]float64, N)
	}
	var tem_1, tem_2, tem_3 float64

	// 对矩阵右半部分进行扩增
	for i := 0; i < N; i++ {
		for j := 0; j < 2*N; j++ {
			if j < N {
				W[i][j] = m.m[i][j]
			} else {
				if j-N == i {
					W[i][j] = 1
				} else {
					W[i][j] = 0
				}
			}
		}
	}

	for i := 0; i < N; i++ {
		// 判断矩阵第一行第一列的元素是否为0，若为0，继续判断第二行第一列元素，直到不为0，将其加到第一行
		if floatIsZero(W[i][i]) {
			var j int
			for j = i + 1; j < N; j++ {
				if !floatIsZero(W[j][i]) {
					break
				}
			}
			if j == N {
				fmt.Println("这个矩阵不能求逆")
				break
			}
			//将前面为0的行加上后面某一行
			for k := 0; k < 2*N; k++ {
				W[i][k] += W[j][k]
			}
		}

		//将前面行首位元素置1
		tem_1 = W[i][i]
		for j := 0; j < 2*N; j++ {
			W[i][j] = W[i][j] / tem_1
		}

		//将后面所有行首位元素置为0
		for j := i + 1; j < N; j++ {
			tem_2 = W[j][i]
			for k := i; k < 2*N; k++ {
				W[j][k] = W[j][k] - tem_2*W[i][k]
			}
		}
	}

	// 将矩阵前半部分标准化
	for i := N - 1; i >= 0; i-- {
		for j := i - 1; j >= 0; j-- {
			tem_3 = W[j][i]
			for k := i; k < 2*N; k++ {
				W[j][k] = W[j][k] - tem_3*W[i][k]
			}
		}
	}

	//得出逆矩阵
	for i := 0; i < N; i++ {
		for j := N; j < 2*N; j++ {
			result[i][j-N] = W[i][j]
		}
	}

	return &Matrix{
		rows:    N,
		columns: N,
		m:       result,
	}
}

// 转置矩阵
func (m *Matrix) Transpose() *Matrix {
	rm := make([][]float64, m.columns)
	for i := 0; i < m.columns; i++ {
		rm[i] = make([]float64, m.rows)
	}
	for i, x := range m.m {
		for j, v := range x {
			rm[j][i] = v
		}
	}
	return &Matrix{
		rows:    m.columns,
		columns: m.rows,
		m:       rm,
	}
}

// 是正交矩阵
func (m *Matrix) IsOrthogonal() bool {
	return MatrixIsSame(m.Transpose(), m.Inverse())
}

// 矩阵相等
func MatrixIsSame(a, b *Matrix) bool {
	if a.rows != b.rows || a.columns != b.columns {
		return false
	}
	for i := 0; i < a.rows; i++ {
		for j := 0; j < a.columns; j++ {
			if !floatEqual(a.m[i][j], b.m[i][j]) {
				return false
			}
		}
	}
	return true
}

// 矩阵加
func MatrixAdd(a, b *Matrix) *Matrix {
	if a.rows != b.rows || a.columns != b.columns {
		panic(fmt.Sprintf("illegal matrix add of size %d×%d and %d×%d", a.rows, a.columns, b.rows, b.columns))
	}
	cm := make([][]float64, a.rows)
	for i, va := range a.m {
		tmp := make([]float64, len(va))
		for j := 0; j < len(va); j++ {
			tmp[j] = a.m[i][j] + b.m[i][j]
		}
		cm[i] = tmp
	}
	return &Matrix{
		rows:    a.rows,
		columns: a.columns,
		m:       cm,
	}
}

// 矩阵乘（数）
func MatrixMulNum(a *Matrix, num float64) *Matrix {
	cm := make([][]float64, a.rows)
	for i, va := range a.m {
		tmp := make([]float64, len(va))
		for j, vb := range va {
			tmp[j] = vb * num
		}
		cm[i] = tmp
	}
	return &Matrix{
		rows:    a.rows,
		columns: a.columns,
		m:       cm,
	}
}

// 矩阵乘矩阵（a×b）
func MatrixMulMatrix(a, b *Matrix) *Matrix {
	if a.columns != b.rows {
		panic(fmt.Sprintf("illegal matrix mul matrix of size %d×%d and %d×%d", a.rows, a.columns, b.rows, b.columns))
	}
	cm := make([][]float64, a.rows)
	for i, _ := range a.m {
		tmp := make([]float64, b.columns)
		for j := 0; j < b.columns; j++ {
			tmp[j] = 0
			for k := 0; k < a.columns; k++ {
				tmp[j] += a.m[i][k] * b.m[k][j]
			}
		}
		cm[i] = tmp
	}
	return &Matrix{
		rows:    a.rows,
		columns: b.columns,
		m:       cm,
	}
}
