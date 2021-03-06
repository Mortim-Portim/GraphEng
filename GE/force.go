package GE

import (
	"fmt"
	"math/rand"
	"time"
)

//Generates a random number including l and u
func RandomInt(l, u int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(u-l+1) + l
}

func RandomFloat(l, u float64) float64 {
	return rand.Float64()*(u-l+1) + l
}

type ForceUpdater func(f *Force)

func GetRandomUF(min, max, deltaAngle, deltaAmount float64) func(*Force) {
	return func(f *Force) {
		f.direction.RotateZ(RandomFloat(-deltaAngle, deltaAngle))
		f.amount += RandomFloat(-deltaAmount, deltaAmount)
		if f.amount > max || f.amount < min {
			f.amount = (max + min) / 2
		}
	}
}
func fUD(f *Force) {
	if f.Count <= 0 {
		f.amount = 0
		f.inActive = true
	}
	f.Count--
	return
}
func fU(f *Force) {
	if f.amount <= 0 {
		f.inActive = true
	}
	return
}
func fUAlwaysActive(f *Force) {
	return
}

type Force struct {
	direction *Vector
	amount    float64
	inActive  bool
	Count     int
	Update    func(f *Force)
}

func GetResultingForce(fs ...*Force) (f *Force) {
	f = GetNewPlainForceAlwaysActive()
	for _, fa := range fs {
		f.SetForceVec(fa.GetForceVec().Add(f.GetForceVec()))
	}
	return
}
func getStandardForce() *Force {
	f := Force{}
	f.amount = 0
	f.Count = 0
	f.direction = &Vector{0, 0, 0}
	return &f
}
func GetNewPlainForce() *Force {
	f := getStandardForce()
	f.Update = fU
	return f
}
func GetNewPlainForceAlwaysActive() *Force {
	f := getStandardForce()
	f.Update = fUAlwaysActive
	return f
}
func GetNewPlainForceDir(x, y, z float64) *Force {
	f := GetNewPlainForce()
	f.SetForceVec(&Vector{x, y, z})
	return f
}
func GetNewPlainForceOfDuration(duration int) *Force {
	if duration < 0 {
		panic("Cannot Create Force with negative Duration")
	}
	f := getStandardForce()
	f.Update = fUD
	f.Count = duration
	return f
}
func GetNewForce(updater ForceUpdater, v *Vector) *Force {
	if updater == nil || v == nil {
		panic("Cannot Create Force with nil Updater oder nil Vector")
	}
	f := getStandardForce()
	f.SetForceVec(v)
	f.Update = updater
	return f
}
func GetNewForceWithDurationAndDir(duration int, v *Vector) *Force {
	f := GetNewPlainForceOfDuration(duration)
	f.SetForceVec(v)
	return f
}
func GetNewRandomForce(minV, maxV float64) *Force {
	return GetNewForce(fU, GetRandomVector(minV, maxV))
}
func (f *Force) SetForceVec(v *Vector) {
	if v == nil {
		panic("Cannot set nil Vector as Force Vector")
	}
	f.amount = v.Absolute()
	f.direction = v.Normalize()
}
func (f *Force) GetForceVec() *Vector {
	if f.direction == nil {
		f.direction = &Vector{}
	}
	return f.direction.Mul(f.amount)
}
func (f *Force) GetAcc(mass float64) *Vector {
	if mass == 0 {
		panic("Cannot calculate Acceleration for Mass=0")
	}
	return f.GetForceVec().Mul(1.0 / mass)
}
func (f *Force) Copy() *Force {
	return &Force{direction: f.direction.Copy(), amount: f.amount, inActive: f.inActive, Count: f.Count, Update: f.Update}
}

func (f *Force) GetInfos() string {
	if f == nil {
		f = GetNewPlainForce()
	}
	v := f.GetForceVec()
	return fmt.Sprintf("%s, |%0.3f|, T:%v", v.GetInfos(), f.amount, f.Count)
}

func (f *Force) IsActive() bool {
	return !f.inActive
}

func (f *Force) Add(f2 *Force) *Force {
	if f2 == nil {
		panic("Cannot add nil to Force")
	}
	nF := f.Copy()
	nF.SetForceVec(f.GetForceVec().Add(f2.GetForceVec()))
	return nF
}
