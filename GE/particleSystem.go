package GE

import "github.com/hajimehoshi/ebiten"

func GetNewParticleSystem(layer uint8) (ps *ParticleSystem) {
	ps = &ParticleSystem{
		Particles: make([]*Particle, 0),
		layer:     layer,
		DrawBox:   GetRectangle(0, 0, 1, 1),
	}
	return
}

type ParticleSystem struct {
	Particles []*Particle
	layer     uint8
	DrawBox   *Rectangle
	x, y      float64
}

func (ps *ParticleSystem) Add(p *Particle) {
	ps.Particles = append(ps.Particles, p)
}
func (ps *ParticleSystem) Update(frame int) {
	for _, p := range ps.Particles {
		p.Update(frame)
		if !p.IsFinished() {
			ps.DrawBox = ps.DrawBox.BoundingRect(p.Drawbox)
		}
	}
	mdl := ps.DrawBox.GetMiddle()
	ps.x = mdl.X
	ps.y = mdl.Y
	ps.RemoveFinishedParticles()
}
func (ps *ParticleSystem) GetDrawBox() *Rectangle {
	return ps.DrawBox
}
func (ps *ParticleSystem) GetPos() (x, y float64, layer int8) {
	return ps.x, ps.y, int8(ps.layer)
}
func (ps *ParticleSystem) Draw(screen *ebiten.Image, _ int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	for _, p := range ps.Particles {
		p.Draw(screen, xStart-leftTopX*sqSize, yStart-leftTopY*sqSize, sqSize)
	}
}
func (ps *ParticleSystem) DrawOnPlainScreen(screen *ebiten.Image, xStart, yStart, scale float64) {
	for _, p := range ps.Particles {
		p.Draw(screen, xStart, yStart, scale)
	}
}

func (ps *ParticleSystem) RemoveFinishedParticles() {
	rems := 0
	for idx := range ps.Particles {
		if ps.Particles[idx-rems].IsFinished() {
			ps.RemoveParticleByIdx(idx - rems)
			rems++
		}
	}
}
func (ps *ParticleSystem) RemoveParticleByIdx(i int) {
	ps.Particles[i] = ps.Particles[len(ps.Particles)-1]
	ps.Particles = ps.Particles[:len(ps.Particles)-1]
}
