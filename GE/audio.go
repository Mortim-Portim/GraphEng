package GE

import (
	"os"
	"time"
	//"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
)

//IMPORTANT: Call InitAudioContext() before creating a new Player
const sampleRate = 48000
var audioContext *audio.Context

func InitAudioContext() {
	audioContext,_ = audio.NewContext(sampleRate)
}


type AudioPlayer struct {
	*audio.Player
	repeating bool
	total time.Duration
}

//Creates a new audio player
func NewPlayer(filename string) (*AudioPlayer, error) {
	f, err := os.Open(filename)
	if err != nil {return nil, err}
	
	d, err := mp3.Decode(audioContext, f)
	if err != nil {return nil, err}

	// Create an audio.Player that has one stream.
	p, err := audio.NewPlayer(audioContext, d)
	if err != nil {return nil, err}
	
	ap := &AudioPlayer{Player:p, total:time.Second * time.Duration(d.Length()) / 4 / sampleRate}
	return ap, nil
}
func (p *AudioPlayer) PlayFromBeginning(volume float64) error {
	p.SetVolume(volume)
	p.Rewind()
	p.Play()
	return nil
}
func (p *AudioPlayer) Pause() {
	p.repeating = false
	p.Player.Pause()
}
func (p *AudioPlayer) Repeat(OnRepeat func() bool) {
	if !p.repeating {
		p.repeating = true
		go func(){
			for p.repeating {
				remaining := p.total - p.Current()
				time.Sleep(remaining+time.Millisecond)
				if p.repeating {
					if !OnRepeat() {
						p.repeating = false
						return
					}
					p.Rewind()
					p.Play()
				}
			}
		}()
	}
}