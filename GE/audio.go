package GE

import (
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
)

//IMPORTANT: Call InitAudioContext() before creating a new Player
const sampleRate = 48000

var audioContext *audio.Context

func InitAudioContext() {
	audioContext = audio.NewContext(sampleRate)
}

type Audio interface {
	SetStandardVolume(v float64)
	GetStandardVolume() float64
	SetVolume(volume float64)
	Volume() float64
	Seek(offset time.Duration) error
	Current() time.Duration
	Rewind() error
	IsPlaying() bool
	Pause()
	Play()
	Resume()
	PlayFromBeginning()
	StopRepeating()
	Repeat()
}

type AudioPlayer struct {
	Ap                  *audio.Player
	repeating           bool
	total, finishedTime time.Duration

	StandardVolume float64
	OnFinished     func()
}

func (p *AudioPlayer) SetStandardVolume(v float64) {
	p.StandardVolume = v
}
func (p *AudioPlayer) GetStandardVolume() float64 {
	return p.StandardVolume
}
func (p *AudioPlayer) SetVolume(volume float64) {
	p.Ap.SetVolume(volume)
}
func (p *AudioPlayer) Volume() float64 {
	return p.Ap.Volume()
}
func (p *AudioPlayer) Seek(offset time.Duration) error {
	return p.Ap.Seek(offset)
}
func (p *AudioPlayer) Current() time.Duration {
	return p.Ap.Current()
}
func (p *AudioPlayer) Rewind() error {
	return p.Ap.Rewind()
}
func (p *AudioPlayer) IsPlaying() bool {
	return p.Ap.IsPlaying()
}
func (p *AudioPlayer) Pause() {
	p.StopRepeating()
	p.Ap.Pause()
}
func (p *AudioPlayer) Resume() {
	p.Ap.Play()
}
func (p *AudioPlayer) Play() {
	p.Ap.Play()
	go func() {
		time.Sleep(p.finishedTime)
		if p.OnFinished != nil && p.IsPlaying() {
			p.OnFinished()
		}
	}()
	return
}
func (p *AudioPlayer) PlayFromBeginning() {
	p.SetVolume(p.StandardVolume)
	p.Rewind()
	p.Play()
}
func (p *AudioPlayer) StopRepeating() {
	p.repeating = false
}
func (p *AudioPlayer) Repeat() {
	if !p.repeating {
		p.repeating = true
		go func() {
			for p.repeating {
				remaining := p.total - p.Current()
				time.Sleep(remaining + time.Millisecond)
				if p.repeating {
					p.PlayFromBeginning()
				}
			}
		}()
	}
}

//Creates a new audio player
func NewPlayer(filename string, StandardVolume float64, OnFinished func(), fT time.Duration) (*AudioPlayer, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	d, err := mp3.Decode(audioContext, f)
	if err != nil {
		return nil, err
	}

	// Create an audio.Player that has one stream.
	p, err := audio.NewPlayer(audioContext, d)
	if err != nil {
		return nil, err
	}
	total := time.Second * time.Duration(d.Length()) / 4 / sampleRate
	ap := &AudioPlayer{
		Ap:             p,
		total:          total,
		StandardVolume: StandardVolume,
		OnFinished:     OnFinished,
		finishedTime:   total - fT,
	}
	return ap, nil
}
