package GE

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"time"
)

const STANDARDVOLUME = 1.0

/**
SoundEffect represents a collection of short audios, that should sound similar
you can play a specific sound or just a random one (using Set() and SetR())
**/
func fadeInFunc(percent float64, ap Audio) float64 {
	return percent * ap.GetStandardVolume()
}
func fadeOutFunc(percent float64, ap Audio) float64 {
	return ap.GetStandardVolume() - percent*ap.GetStandardVolume()
}

type SoundEffect struct {
	*AudioPlayer
	sounds []*AudioPlayer
}

func (s *SoundEffect) SetVolume(v float64) {
	for _, snd := range s.sounds {
		snd.SetStandardVolume(v)
		snd.SetVolume(v)
	}
}

//Sets a specific audio file
func (s *SoundEffect) Set(idx int) {
	if idx < 0 || idx >= len(s.sounds) {
		return
	}
	s.Pause()
	s.Rewind()
	s.AudioPlayer = s.sounds[idx]
}
func (s *SoundEffect) SetR(seed int64) {
	rand.Seed(seed)
	s.Set(rand.Intn(len(s.sounds)))
}
func (s *SoundEffect) onFinished() {
	s.Pause()
	s.Rewind()
}

//loads all audio files in a folder
func LoadSounds(folder string) (*SoundEffect, error) {
	s := &SoundEffect{}
	s.sounds = make([]*AudioPlayer, 0)
	if folder[len(folder)-1:] != "/" {
		folder += "/"
	}
	files, err1 := ioutil.ReadDir(folder)
	if err1 != nil {
		return nil, err1
	}
	for _, f := range files {
		name := f.Name()
		file := folder + name
		//fmt.Println("Loading: ", file)
		player, err := NewPlayer(file, STANDARDVOLUME, s.onFinished, time.Duration(0))
		if err != nil {
			break
		}
		s.sounds = append(s.sounds, player)
	}
	if len(s.sounds) == 0 {
		return nil, errors.New("Could not load SoundEffects")
	}
	return s, nil
}

//Fades out
func FadeOut(s Audio, seconds float64, done chan bool) {
	millis := seconds * 1000.0
	FadePlayer(s, nil, int(millis), int(millis/2.0), nil, fadeOutFunc, false, done)
	go func() {
		time.Sleep(time.Duration(float64(time.Millisecond) * millis))
		s.Pause()
	}()
}

//Fades to another sound
func FadeTo(s, s2 Audio, seconds float64, done chan bool) {
	millis := seconds * 1000.0
	FadePlayer(s, s2, int(millis), int(millis/2.0), fadeInFunc, fadeOutFunc, true, done)
}
func FadePlayer(oldP, newP Audio, milliseconds, iterations int, volumefaderNew, volumefaderOld func(percent float64, ap Audio) (volume float64), changeToNew bool, done chan bool) {
	delay := float64(milliseconds) / float64(iterations)
	if newP != nil {
		newP.PlayFromBeginning()
		newP.SetVolume(0.0)
	}
	if oldP != nil {
		oldP.SetVolume(volumefaderOld(0, oldP))
	}
	go func() {
		for i := 0; i < iterations; i++ {
			percent := float64(i+1) / float64(iterations)
			if oldP != nil {
				oldP.SetVolume(volumefaderOld(percent, oldP))
			}
			if newP != nil {
				newP.SetVolume(volumefaderNew(percent, newP))
			}
			time.Sleep(time.Duration(int(float64(time.Millisecond) * delay)))
		}
		if changeToNew {
			if oldP != nil {
				oldP.Pause()
			}
		} else if newP != nil {
			newP.Pause()
		}
		if newP != nil {
			newP.SetVolume(volumefaderNew(1.0, newP))
		}
		if done != nil {
			done <- true
		}
	}()
}
