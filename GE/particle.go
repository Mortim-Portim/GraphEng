package GE

type Particle struct {
	*WObj
	Mass     float32
	Velocity *Vector
	Forces   []*Force
}

func (p *Particle) Update() {
	for _, f := range p.Forces {
		p.Velocity.Add(f.GetAcc(float64(p.Mass)))
	}
	p.MoveBy(p.Velocity.X, p.Velocity.Y)
}
