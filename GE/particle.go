package GE

type Particle struct {
	*WObj
	Mass   float64
	Forces []*Force
}
