package antenna

import (
	"coursework3/internal/physics"
	"encoding/json"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
)

func getDeg(x, y float64) float64 {
	deg := math.Atan2(y, x)
	if deg < 0 {
		deg = 2*math.Pi + deg
	}
	return deg * 180 / math.Pi
}

type AntennaSet struct {
	antennas []*SlotAntenna
	len      int
}

func MakeSet(antennas ...*SlotAntenna) *AntennaSet {
	set := &AntennaSet{
		len: len(antennas),
	}
	for _, antenna := range antennas {
		antenna.initialize()
		set.antennas = append(set.antennas, antenna)
	}
	return set
}

func (as *AntennaSet) Marshal(ind int) ([]byte, error) {
	ind = ind % as.len
	b, err := json.Marshal(as.antennas[ind])
	if err != nil {
		return nil, err
	}
	return b, nil
}

type SlotAntenna struct {
	SlotSize   int
	Wavelength float64
	Ds         float64
	Size_I     int
	Size_J     int
	H          float64
	I_offset   int
	J_offset   int
	MD         [][]float64
	ED         [][]float64
}

func NewSlotAntenna(size, size_i, size_j int, wavelength, h float64) *SlotAntenna {
	return &SlotAntenna{
		SlotSize:   size,
		Wavelength: wavelength,
		Ds:         h * h,
		H:          h,
		Size_J:     size_j,
		Size_I:     size_i,
	}
}

func computeDegs(xM, yM, xE, yE [][]complex128) ([][]float64, [][]float64) {
	md, ed := make([][]float64, len(xM)), make([][]float64, len(xM))
	for i := 0; i < len(xM); i++ {
		md[i], ed[i] = make([]float64, len(xM[0])), make([]float64, len(xM[0]))
		for j := 0; j < len(xM[0]); j++ {
			//md[i][j], ed[i][j] = getDeg(imag(xM[i][j]), imag(yM[i][j])), getDeg(imag(xE[i][j]), imag(yE[i][j]))
			md[i][j], ed[i][j] = getDeg(real(xM[i][j]), real(yM[i][j])), getDeg(real(xE[i][j]), real(yE[i][j]))
		}
	}
	return md, ed
}

func (an *SlotAntenna) initialize() {
	j_offset := (an.Size_J - an.SlotSize) / 2
	i_offset := (an.Size_I - 1) / 2
	an.I_offset = i_offset
	an.J_offset = j_offset
	middle := an.SlotSize / 2
	A := physics.MakeA(an.SlotSize, an.Wavelength, an.Ds, an.H)
	a := make([]float64, an.SlotSize)
	a[middle] = 1
	B, err := physics.MakeB(an.SlotSize, an.Wavelength, an.H, a)
	if err != nil {
		log.Fatal(err)
	}
	O := &mat.Dense{}
	err = O.Solve(A, B)
	if err != nil {
		log.Fatal(err)
	}
	jmx := make([]complex128, an.SlotSize)
	for i := 0; i < an.SlotSize; i++ {
		jmx[i] = complex(O.At(i, 0), O.At(i+an.SlotSize, 0))
	}
	AMX := physics.GetAMX(an.Size_I, an.Size_J, i_offset, j_offset, an.SlotSize, jmx, an.Wavelength, an.Ds, an.H)
	xM, yM := physics.GetMagnetVec(AMX, an.H, an.Wavelength)
	xE, yE := physics.GetElectricVec(xM, yM)
	md, ed := computeDegs(xM, yM, xE, yE)
	an.MD = md
	an.ED = ed
}
