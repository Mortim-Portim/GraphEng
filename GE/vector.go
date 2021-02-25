package GE

import (
	"fmt"
	"math"

	cmp "github.com/mortim-portim/GraphEng/compression"
)

func GetVectorFromRot(a float64) (v *Vector) {
	v = &Vector{}
	v.RotateAbsZ(a)
	return
}

//3d Vector
type Vector struct {
	X, Y, Z float64
}

//Returns true if v is filled
func (v *Vector) IsFilled() bool {
	if v.X != 0 || v.Y != 0 || v.Z != 0 {
		return true
	}
	return false
}

//Rotates v around the Z axis
func (v *Vector) RotateZ(angle float64) *Vector {
	return v.Rotate(&Vector{0, 0, 1}, angle)
}

//Rotates v around the axis
func (v *Vector) Rotate(axis *Vector, angle float64) *Vector {
	angle *= (math.Pi / 180)
	if axis == nil || !axis.IsFilled() {
		panic("Cannot Rotate Vector around Zero or nil Axis")
	}
	r := axis.X
	s := axis.Y
	t := axis.Z
	m := r*v.X + s*v.Y + t*v.Z
	xPrime := (float64)(r*m*(1-math.Cos(angle)) + v.X*math.Cos(angle) + (-t*v.Y+s*v.Z)*math.Sin(angle))
	yPrime := (float64)(s*m*(1-math.Cos(angle)) + v.Y*math.Cos(angle) + (t*v.X-r*v.Z)*math.Sin(angle))
	zPrime := (float64)(t*m*(1-math.Cos(angle)) + v.Z*math.Cos(angle) + (-s*v.X+r*v.Y)*math.Sin(angle))
	v.X = xPrime
	v.Y = yPrime
	v.Z = zPrime
	return v
}

//Sets the Rotation of v around the Z axis
func (v *Vector) RotateAbsZ(angle float64) *Vector {
	return v.RotateAbs(&Vector{0, 0, 1}, angle)
}

//Sets the Rotation of v around the axis
func (v *Vector) RotateAbs(axis *Vector, angle float64) *Vector {
	v.X = 0
	v.Y = -1
	v.Z = 0
	v.Rotate(axis, angle)
	return v
}

//Returns the Rotation of v around the Z axis
func (v *Vector) GetRotationZ() (angle float64) {
	angle = 0
	if v.IsFilled() {
		if v.Y < 0 {
			angle = math.Atan(math.Abs(v.X)/-v.Y) / (math.Pi) * 180.0
		} else {
			angle = 90 + math.Atan(v.Y/math.Abs(v.X))/(math.Pi)*180.0
		}
		if v.X < 0 {
			angle = 360 - angle
		}
	}
	return
}

//Returns the dotproduct of two Vectors
func (v *Vector) DotProduct(B *Vector) (product float64) {
	if B == nil {
		panic("Cannot compute the DotProduct of a nil Vector")
	}
	product = 0
	product += v.X * B.X
	product += v.Y * B.Y
	product += v.Z * B.Z
	return product
}

//Returns the crossproduct of two Vectors
func (v *Vector) CrossProduct(B *Vector) *Vector {
	if B == nil {
		panic("Cannot compute the CrossProduct of a nil Vector")
	}
	cross_P := Vector{}
	cross_P.X = (v.Y*B.Z - v.Z*B.Y)
	cross_P.Y = (v.Z*B.X - v.X*B.Z)
	cross_P.Z = (v.X*B.Y - v.Y*B.X)
	return &cross_P
}

//Returns the sum of two Vectors
func (v *Vector) Add(B *Vector) *Vector {
	if B == nil {
		panic("Cannot add a nil Vector")
	}
	return &Vector{v.X + B.X, v.Y + B.Y, v.Z + B.Z}
}

//Subtracts B from v
func (v *Vector) Sub(B *Vector) *Vector {
	if B == nil {
		panic("Cannot subtract a nil Vector")
	}
	return &Vector{v.X - B.X, v.Y - B.Y, v.Z - B.Z}
}

//Returns the product of a Vector an a Skalar
func (v *Vector) Mul(num float64) *Vector {
	prod := Vector{}
	prod.X = (v.X * num)
	prod.Y = (v.Y * num)
	prod.Z = (v.Z * num)
	return &prod
}

//Checks for value equality
func (v *Vector) Equals(B *Vector) bool {
	if B != nil {
		if v.X == B.X && v.Y == B.Y && v.Z == B.Z {
			return true
		}
	}
	return false
}

//Returns the length of v
func (v *Vector) Absolute() float64 {
	if v.IsFilled() {
		return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	}
	return 0
}

//Returns normalized v
func (v *Vector) Normalize() *Vector {
	length := v.Absolute()
	if length != 0 {
		v.X = (v.X / length)
		v.Y = (v.Y / length)
		v.Z = (v.Z / length)
	} else {
		v.X = 0
		v.Y = 0
		v.Z = 0
	}
	return v.Copy()
}

//Returns a Copy of v
func (v *Vector) Copy() *Vector {
	return &Vector{v.X, v.Y, v.Z}
}

//Returns Infos about v
func (v *Vector) GetInfos() string {
	return fmt.Sprintf("[%0.3f, %0.3f, %0.3f]", v.X, v.Y, v.Z)
}

func (v *Vector) ToBytes() []byte {
	return append(append(cmp.Float64ToBytes(v.X), cmp.Float64ToBytes(v.Y)...), cmp.Float64ToBytes(v.Z)...)
}
func (v *Vector) XYToBytes() []byte {
	return append(cmp.Float64ToBytes(v.X), cmp.Float64ToBytes(v.Y)...)
}
func XYVectorFromBytes(bs []byte) (v *Vector) {
	v = &Vector{}
	v.X = cmp.BytesToFloat64(bs[:8])
	v.Y = cmp.BytesToFloat64(bs[8:])
	v.Z = 0
	return
}
func VectorFromBytes(bs []byte) (v *Vector) {
	v = &Vector{}
	v.X = cmp.BytesToFloat64(bs[:8])
	v.Y = cmp.BytesToFloat64(bs[8:16])
	v.Z = cmp.BytesToFloat64(bs[16:])
	return
}
func GetRandomVector() *Vector {
	angle := RandomFloat(0, 360)
	return GetVectorFromRot(angle)
}
