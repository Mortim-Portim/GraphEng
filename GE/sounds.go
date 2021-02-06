package GE

import (
	"math/rand"
	"io/ioutil"
	"errors"
	"time"
)
const STANDARDVOLUME = 1.0

func (s *Sounds) fadeInFunc(percent float64)(float64){return percent*s.StandardVolume}
func (s *Sounds) fadeOutFunc(percent float64)(float64){return s.StandardVolume-percent*s.StandardVolume}

type Sounds struct {
	sounds map[string]*AudioPlayer
	currentPlayer string
	StandardVolume float64
	OnRepeat func() bool
}
func (s *Sounds) SetVolume(volume float64) {
	s.StandardVolume = volume
	if pl,ok := s.sounds[s.currentPlayer]; ok && pl.IsPlaying() {
		s.sounds[s.currentPlayer].SetVolume(s.StandardVolume)
	}
}
func (s *Sounds) PlayInfinite(OnRepeat func() bool) {
	p, ok := s.sounds[s.currentPlayer]
	if ok {
		p.Repeat(OnRepeat)
	}
}
func (s *Sounds) FadeToSoundR(s2 *Sounds, seed int64, done chan bool) {
	new_file := s2.GetRandomSound(seed)
	s.FadeToSound(s2, new_file, STANDARD_FADE_TIME, done)
}
func (s *Sounds) FadeToSound(s2 *Sounds, new_file string, seconds float64, done chan bool) {
	oldP, ok := s.sounds[s.currentPlayer]
	if !ok {
		oldP = nil
	}
	
	newP, ok := s2.sounds[new_file]
	if !ok || newP == oldP {
		newP = nil
	}else{
		s2.currentPlayer = new_file
	}
	millis := seconds*1000.0
	FadePlayer(oldP, newP, int(millis), int(millis/2.0), s2.fadeInFunc, s.fadeOutFunc, true, done)
	if s.OnRepeat != nil {
		s.PlayInfinite(s.OnRepeat)
	}
}
func (s *Sounds) FadeOut(seconds float64, done chan bool) {
	millis := seconds*1000.0
	s.ChangeTo(int(millis), int(millis/2.0), "", nil, s.fadeOutFunc, false, done)
	go func(){
		time.Sleep(time.Duration(float64(time.Millisecond)*millis))
		s.PauseAll()
	}()
}
func (s *Sounds) FadeTo(new_file string, seconds float64, done chan bool) {
	millis := seconds*1000.0
	s.ChangeTo(int(millis), int(millis/2.0), new_file, s.fadeInFunc, s.fadeOutFunc, true, done)
}
func (s *Sounds) FadeToR(seed int64, seconds float64, done chan bool) {
	millis := seconds*1000.0
	s.ChangeTo(int(millis), int(millis/2.0), s.GetRandomSound(seed), s.fadeInFunc, s.fadeOutFunc, true, done)
}

//loads all audio files in a folder
func LoadSounds(folder string) (*Sounds, error) {
	s := &Sounds{StandardVolume:STANDARDVOLUME}
	s.sounds = make(map[string]*AudioPlayer)
	if folder[len(folder)-1:] != "/" {
		folder += "/"
	}
	files, err1 := ioutil.ReadDir(folder)
    if err1 != nil {return nil, err1}
	
	for _, f := range files {
		name := f.Name()
		file := folder+name
		//fmt.Println("Loading: ", file)
		player, err := NewPlayer(file)
		if err != nil {break}
		s.sounds[name[:len(name)-4]] = player
	}
	if len(s.sounds) == 0 {
		return nil, errors.New("Could not load SoundEffects")
	}
	return s, nil
}
func (s *Sounds) Resume() {
	s.sounds[s.currentPlayer].Play()
}
//Plays a specific audio file
func (s *Sounds) PS(file string) {
	if p, ok := s.sounds[s.currentPlayer]; ok && p.IsPlaying() {
		s.sounds[s.currentPlayer].Pause()
	}
	s.currentPlayer = file
	if !s.sounds[file].IsPlaying() {
		s.sounds[file].PlayFromBeginning(s.StandardVolume)
	}
	if s.OnRepeat != nil {
		s.PlayInfinite(s.OnRepeat)
	}
}
//Plays a random audio file
func (s *Sounds) PR(seed int64) {
	s.PS(s.GetRandomSound(seed))
}
func (s *Sounds) PauseAll() {
	for _,p := range(s.sounds) {
		if p.IsPlaying() {
			p.Pause()
		}
	}
}
func (s *Sounds) GetRandomSound(seed int64) string {
	rand.Seed(seed)
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

func (s *Sounds) ChangeTo(milliseconds, iterations int, new_file string, volumefaderNew, volumefaderOld  func(percent float64)(volume float64), changeToNew bool, done chan bool) {
	oldP, ok := s.sounds[s.currentPlayer]
	if !ok {
		oldP = nil
	}
	newP, ok := s.sounds[new_file]
	if !ok || newP == oldP {
		newP = nil
	}else{
		s.currentPlayer = new_file
	}
	FadePlayer(oldP, newP, milliseconds, iterations, volumefaderNew, volumefaderOld, changeToNew, done)
	if s.OnRepeat != nil {
		s.PlayInfinite(s.OnRepeat)
	}
}

func FadePlayer(oldP, newP *AudioPlayer, milliseconds, iterations int, volumefaderNew, volumefaderOld  func(percent float64)(volume float64), changeToNew bool, done chan bool) {
	delay := float64(milliseconds)/float64(iterations)
	if newP != nil {
		newP.PlayFromBeginning(0.0)
	}
	if oldP != nil {
		oldP.SetVolume(volumefaderOld(0))
	}
	go func() {
		for i := 0; i < iterations; i++ {
			percent := float64(i+1)/float64(iterations)
			if oldP != nil {
				oldP.SetVolume(volumefaderOld(percent))
			}
			if newP != nil {
				newP.SetVolume(volumefaderNew(percent))
			}
			time.Sleep(time.Duration(int(float64(time.Millisecond)*delay)))
		}
		if changeToNew {
			if oldP != nil {
				oldP.Pause()
			}
		}else if newP != nil {
			newP.Pause()
		}
		if newP != nil {
			newP.SetVolume(volumefaderNew(1.0))
		}
		if done != nil {
			done <- true
		}
	}()
}