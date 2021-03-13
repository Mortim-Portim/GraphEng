package GE

import ebiten "github.com/hajimehoshi/ebiten/v2"

func GetParticle(frame, frameTime, fps int, mass float64, dir *Vector, anim *Animation, X, Y, W, H float64) (p *Particle) {
	p = &Particle{
		Anim:          anim,
		FPS:           fps,
		W:             W,
		H:             H,
		Mass:          mass,
		Velocity:      dir,
		Forces:        make([]*Force, 0),
		TimedOutFrame: frame + frameTime,
	}
	p.SetMiddle(X, Y)
	return
}
func GetParticleWithRandomDirAndForce(frame, frameTime, fps int, mass float64, anim *Animation, X, Y, W, H float64) (p *Particle) {
	p = GetParticle(frame, frameTime, fps, mass, GetRandomVector(0, 0.0000001), anim, X, Y, W, H)
	p.Forces = append(p.Forces, GetNewRandomForce(0, 0))
	return
}

func GetNewParticleFactory(frameTime, fps int, anim *Animation) *ParticleFactory {
	return &ParticleFactory{
		frameTime: frameTime,
		FPS:       fps,
		anim:      anim,
	}
}

type ParticleFactory struct {
	frameTime, FPS int
	anim           *Animation
}

func (pf *ParticleFactory) GetNewRandom(frame int, mass float64, X, Y, W, H float64) *Particle {
	return GetParticleWithRandomDirAndForce(frame, pf.frameTime, pf.FPS, mass, pf.anim.Clone(), X, Y, W, H)
}
func (pf *ParticleFactory) GetNew(frame int, mass float64, dir *Vector, X, Y, W, H float64) *Particle {
	return GetParticle(frame, pf.frameTime, pf.FPS, mass, dir, pf.anim.Clone(), X, Y, W, H)
}

type Particle struct {
	Anim    *Animation
	FPS     int
	W, H    float64
	Drawbox *Rectangle

	Mass     float64
	Velocity *Vector
	Forces   []*Force

	TimedOut      bool
	TimedOutFrame int
}

func (p *Particle) Update(frame int) {
	if p.TimedOutFrame != 0 && frame >= p.TimedOutFrame {
		p.TimedOut = true
	}
	p.Anim.Update(frame)
	p.Velocity = p.Velocity.Add(GetResultingForce(p.Forces...).GetAcc(p.Mass).Mul(1 / float64(p.FPS)))
	p.MoveBy(p.Velocity.X, p.Velocity.Y)
	p.Anim.Angle = p.Velocity.GetRotationZ()
}
func (p *Particle) IsFinished() bool {
	return p.TimedOut
}
func (p *Particle) MoveBy(dx, dy float64) {
	p.Drawbox.MoveBy(dx, dy)
}
func (p *Particle) SetTopLeft(x, y float64) {
	p.Drawbox = GetRectangle(x, y, x+p.W, y+p.H)
}
func (p *Particle) SetMiddle(x, y float64) {
	p.SetTopLeft(x-p.W/2, y-p.H/2)
}
func (p *Particle) GetPos() (float64, float64) {
	pnt := p.Drawbox.GetMiddle()
	return pnt.X, pnt.Y
}
func (p *Particle) Draw(screen *ebiten.Image, xStart, yStart, scale float64) {
	p.Anim.SetXYWH(xStart+p.Drawbox.min.X*scale, yStart+p.Drawbox.min.Y*scale, p.Drawbox.Bounds().X*scale, p.Drawbox.Bounds().Y*scale)
	p.Anim.Draw(screen)
}
