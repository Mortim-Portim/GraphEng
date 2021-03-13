package GE

import ebiten "github.com/hajimehoshi/ebiten/v2"

func GetNewParticleSystem(layer uint8, pf *ParticleFactory) (ps *ParticleSystem) {
	ps = &ParticleSystem{
		Particles: make([]*Particle, 0),
		layer:     layer,
		DrawBox:   GetRectangle(0, 0, 1, 1),
		pf:        pf,
	}
	return
}

type ParticleSystem struct {
	Particles []*Particle
	layer     uint8
	DrawBox   *Rectangle
	x, y      float64
	pf        *ParticleFactory
}

func (ps *ParticleSystem) Spawn(count, frame int, mass float64, dir *Vector, deltaA float64, X, Y, W, H float64) {
	if ps.pf == nil {
		return
	}

	angle := -deltaA
	for i := 0; i < count; i++ {
		nDir := dir.Copy().RotateZ(angle)
		angle += deltaA * 2 / float64(count)
		ps.Add(ps.pf.GetNew(frame, mass, nDir, X, Y, W, H))
	}
}
func (ps *ParticleSystem) SpawnRandom(count, frame int, mass, X, Y, W, H float64) {
	if ps.pf == nil {
		return
	}
	for i := 0; i < count; i++ {
		ps.Add(ps.pf.GetNewRandom(frame, mass, X, Y, W, H))
	}
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
