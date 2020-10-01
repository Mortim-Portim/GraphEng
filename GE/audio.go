package GE

import (
	"log"
	"io/ioutil"
    "errors"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

const (
	sampleRate = 48000
)

var (
	audioContext *audio.Context
)

// AudioPlayer represents the current audio state.
type AudioPlayer struct {
	
	audioPlayer  *audio.Player
	seBytes      []byte
	seCh         chan []byte
}

func InitAudioContext() {
	AudioContext, err := audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}
	audioContext = AudioContext
}

func NewPlayer(filename string) (*AudioPlayer, error) {
    b, fErr := ioutil.ReadFile(filename)
    if fErr != nil {
        return nil, fErr
    }
	player := &AudioPlayer{
		seCh:         make(chan []byte),
	}
    go func() {
		s, err := wav.Decode(audioContext, audio.BytesReadSeekCloser(b))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := ioutil.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		player.seCh <- b
	}()
	return player, nil
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}

func (p *AudioPlayer) load() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}

	if p.seBytes == nil {
		return errors.New("Wrong sample rate")
	}
    return nil
}

func (p *AudioPlayer) Play() error {
    err := p.load()
    if err != nil {
        return err
    }
	sePlayer, _ := audio.NewPlayerFromBytes(audioContext, p.seBytes)
	sePlayer.Play()
	return nil
}