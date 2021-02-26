package GE

func GetTimedParticle() {

}

type Particle struct {
	*WObj
	Mass     float32
	Velocity *Vector
	Forces   []*Force

	TimedOut      bool
	TimedOutFrame int
}

func (p *Particle) Update(frame int) {
	if frame >= p.TimedOutFrame {
		p.TimedOut = true
	}
	p.WObj.Update(frame)
	for _, f := range p.Forces {
		p.Velocity.Add(f.GetAcc(float64(p.Mass)))
	}
	p.MoveBy(p.Velocity.X, p.Velocity.Y)
}
func (p *Particle) IsFinished() bool {
	return p.TimedOut
}
