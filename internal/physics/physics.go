package physics

import (
	"errors"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/cmplx"
)

func C1Coef(wavelength float64) complex128 {
	return complex(0, -1*240*math.Pow(math.Pi, 2)/wavelength)
}

func C2Coef(wavelength float64) complex128 {
	return complex(math.Pow(2*math.Pi/wavelength, 2), 0)
}

func KCoef(wavelength, h float64) complex128 {
	return complex(math.Pow(2*math.Pi*(h)/(wavelength), 2)-2, 0)
}

func RCoef(wavelength, h float64) complex128 {
	return complex(240*math.Pow(math.Pi*(h), 2)/wavelength, 0)
}

func C(wavelength, ds, h float64, k, v int) complex128 {
	r := complex(rkv(ds, h, k, v), 0)
	x := complex(0, 2*math.Pi/wavelength)
	return cmplx.Exp(x*r) / r * complex(ds, 0)
}

func C2(wavelength, ds, h float64, ik, jk, iv, jv int) complex128 {
	r := complex(rkv2(ik, jk, iv, jv, h, ds), 0)
	x := complex(0, 2*math.Pi/wavelength)
	return cmplx.Exp(x*r) / r * complex(ds, 0)
}

func rkv(ds, h float64, k, v int) float64 {
	if v == k {
		return math.Sqrt(ds / (4 * math.Pi))
	} else {
		return h * math.Abs(float64(k-v))
	}
}

func rkv2(ik, jk, iv, jv int, h, ds float64) float64 {
	if ik == iv && jk == jv {
		return math.Sqrt(ds / (4 * math.Pi))
	}
	return math.Sqrt(math.Pow(float64(ik-iv)*h, 2) + math.Pow(float64(jk-jv)*h, 2))
}

func MakeA(size int, wavelength, ds, h float64) *mat.Dense {
	data := make([]float64, size*size*4)
	kCoef := KCoef(wavelength, h)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			k := i + 1
			v := j + 1
			val := C(wavelength, ds, h, k-1, v) + kCoef*C(wavelength, ds, h, k, v) + C(wavelength, ds, h, k+1, v)
			data[2*size*i+j] = real(val)
			data[2*size*(size+i)+size+j] = real(val)
			data[2*size*i+size+j] = -imag(val)
			data[2*size*(size+i)+j] = imag(val)
		}
	}
	return mat.NewDense(2*size, 2*size, data)
}

func MakeB(size int, wavelength, h float64, current []float64) (*mat.Dense, error) {
	if size != len(current) {
		return nil, errors.New("wrong size")
	}
	data := make([]float64, size*2)
	for i := 0; i < size; i++ {
		val := complex(0, -1) * RCoef(wavelength, h) * complex(current[i], 0)
		data[i] = real(val)
		data[size+i] = imag(val)
	}
	return mat.NewDense(2*size, 1, data), nil
}

func GetAMX(a_size_i, a_size_j, i_offset, j_offset, size int, jmx []complex128, wavelength, ds, h float64) [][]complex128 {
	res := make([][]complex128, a_size_i)
	for i := 0; i < a_size_i; i++ {
		res[i] = make([]complex128, a_size_j)
		for j := 0; j < a_size_j; j++ {
			for k := 0; k < size; k++ {
				res[i][j] += jmx[k] * C2(wavelength, ds, h, i, j, i_offset, k+j_offset)
			}
			res[i][j] /= 4 * math.Pi
		}
	}
	return res
}

func ddAMXdxdy(AMX [][]complex128, h float64, i, j int) complex128 {
	return (AMX[i+1][j+1] - AMX[i-1][j+1] - AMX[i+1][j-1] + AMX[i-1][j-1]) / complex(4*math.Pow(h, 2), 0)
}

func ddAMXddx(AMX [][]complex128, h float64, i, j int) complex128 {
	return (AMX[i+1][j] - 2*AMX[i][j] + AMX[i-1][j]) / complex(math.Pow(h, 2), 0)
}

func GetMagnetVec(AMX [][]complex128, h, wavelength float64) ([][]complex128, [][]complex128) {
	l := len(AMX)
	l1 := len(AMX[0])
	x := make([][]complex128, l-2)
	y := make([][]complex128, l-2)
	for i := 0; i < l-2; i++ {
		x[i] = make([]complex128, l1-2)
		y[i] = make([]complex128, l1-2)
	}
	for i1 := 0; i1 < l-2; i1++ {
		for j1 := 0; j1 < l1-2; j1++ {
			i := i1 + 1
			j := j1 + 1
			x[i1][j1] = (C2Coef(wavelength)/C1Coef(wavelength))*AMX[i][j] + ddAMXddx(AMX, h, i, j)/C1Coef(wavelength)
			y[i1][j1] = ddAMXdxdy(AMX, h, i, j) / C1Coef(wavelength)
		}
	}
	return x, y
}

func GetElectricVec(xM, yM [][]complex128) ([][]complex128, [][]complex128) {
	l := len(xM)
	l1 := len(xM[0])
	xE := make([][]complex128, l)
	yE := make([][]complex128, l)
	for i, row := range xM {
		xE[i] = make([]complex128, l1)
		yE[i] = make([]complex128, l1)
		for j, val := range row {
			yE[i][j] = -val
			xE[i][j] = yM[i][j]
		}
	}
	return xE, yE
}
