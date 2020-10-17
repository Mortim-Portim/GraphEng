package GE

import (
	"fmt"
	"math/rand"
	"io/ioutil"
	"errors"
	"time"
	//"github.com/hajimehoshi/ebiten/audio"
)

func StandardFader(percent float64)(float64){return percent}

type Sounds struct {
	sounds map[string]*AudioPlayer
	currentPlayer string
}
func (s *Sounds) ChangeTo(milliseconds, iterations int, new_file string, volumefader func(percent float64)(volume float64)) {
	delay := float64(milliseconds)/float64(iterations)
	old_file := s.currentPlayer; s.currentPlayer = new_file
	s.sounds[new_file].PlayFromBeginning(0.0)
	s.sounds[old_file].SetVolume(1.0)
	go func() {
		for i := 0; i < iterations; i++ {
			percent := float64(i+1)/float64(iterations)
			volume := volumefader(percent)
			
			s.sounds[old_file].SetVolume(1.0-volume)
			s.sounds[new_file].SetVolume(volume)
			time.Sleep(time.Duration(int(float64(time.Millisecond)*delay)))
		}
		s.sounds[old_file].Pause()
	}()
}

func (s *Sounds) FadeTo(new_file string, seconds float64) {
	millis := seconds*1000.0
	s.ChangeTo(int(millis), int(millis/2.0), new_file, StandardFader)
}
func (s *Sounds) FadeToR(seed int, seconds float64) {
	millis := seconds*1000.0
	s.ChangeTo(int(millis), int(millis/2.0), s.GetRandomSound(seed), StandardFader)
}

//loads all audio files in a folder
func LoadSounds(folder string) (*Sounds, error) {
	s := &Sounds{}
	s.sounds = make(map[string]*AudioPlayer)
	if folder[len(folder)-1:] != "/" {
		folder += "/"
	}
	files, err1 := ioutil.ReadDir(folder)
    if err1 != nil {return nil, err1}
	
	for _, f := range files {
		name := f.Name()
		file := folder+name
		fmt.Println("Loading: ", file)
		player, err := NewPlayer(file)
		if err != nil {break}
		s.sounds[name[:len(name)-4]] = player
	}
	if len(s.sounds) == 0 {
		return nil, errors.New("Could not load SoundEffects")
	}
	return s, nil
}
//Plays a specific audio file
func (s *Sounds) PS(file string) {
	if p, ok := s.sounds[s.currentPlayer]; ok && p.IsPlaying() {
		s.sounds[s.currentPlayer].Pause()
	}
	s.currentPlayer = file
	if !s.sounds[file].IsPlaying() {
		s.sounds[file].PlayFromBeginning(1.0)
	}
}
//Plays a random audio file
func (s *Sounds) PR(seed int) {
	s.PS(s.GetRandomSound(seed))
}

func (s *Sounds) GetRandomSound(seed int) string {
	rand.Seed(int64(seed))
	idx := rand.Intn(len(s.sounds))
	counter := 0
	for f,_ := range(s.sounds) {
		if counter == idx {
			return f
		}
		counter ++
	}
	return ""
}