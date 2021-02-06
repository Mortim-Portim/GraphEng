package GE

import (
	"math/rand"
	"strings"
	"time"
)

const STANDARD_FADE_TIME = 1.5

type SoundTrack struct {
	Tracks map[string]*AudioPlayer
	current string
	waitingForFading chan bool
	waitingCount, maximumWaitingLength int
	NextTrack string
	OnFinished func()
}

func LoadSoundTrack(path string, maximumWaitingLength int) (*SoundTrack, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	files, err := OSReadDir(path)
	if err != nil {return nil, err}
	s := &SoundTrack{maximumWaitingLength:maximumWaitingLength}; s.Tracks = make(map[string]*AudioPlayer)
	for _,f := range(files) {
		s.Tracks[strings.Split(f, ".")[0]], err = NewPlayer(path+f, STANDARDVOLUME, s.onTrackAlmostFinished, STANDARD_FADE_TIME*1000*time.Millisecond)
		if err != nil {return nil, err}
		//s.Tracks[f].OnRepeat = s.onTrackRepeat
	}
	s.waitingForFading = make(chan bool)
	go func(){
		s.waitingForFading <- true
	}()
	return s,nil
}
func (t *SoundTrack) onTrackAlmostFinished() {
	if t.OnFinished != nil {t.OnFinished()}
	if t.NextTrack == t.current {t.NextTrack = ""}
	curr,_ := t.Tracks[t.current]
	_, ok := t.Tracks[t.NextTrack]
	if ok {
		t.Play(t.NextTrack)
		t.NextTrack = ""
	}else{
		curr.Repeat()
	}
}
func (t *SoundTrack) Play(name string) {
	if name == t.current {
		if t.Tracks[t.current].IsPlaying() {
			return
		}
		t.current = ""
	}
	for t.waitingCount > t.maximumWaitingLength {
		t.waitingForFading <- false
		t.waitingCount --
	}
	t.waitingCount ++
	go func(){
		done := <-t.waitingForFading
		if !done {
			return
		}
		t.waitingCount --
		if next,ok := t.Tracks[name]; ok {
			if curr,ok := t.Tracks[t.current]; ok {
				FadeTo(curr, next, STANDARD_FADE_TIME, t.waitingForFading)
			}else{
				FadeTo(nil, next, STANDARD_FADE_TIME, t.waitingForFading)
			}
			t.current = name
		}else{
			t.waitingForFading <- true
		}
	}()
}
func (t *SoundTrack) SetVolume(volume float64) {
	for _,s := range t.Tracks {
		s.SetStandardVolume(volume)
		s.SetVolume(volume)
	}
}
func (t *SoundTrack) Pause() {
	for _,s := range t.Tracks {
		s.Pause()
	}
}
func (t *SoundTrack) Resume() {
	tr, ok := t.Tracks[t.current]
	if ok && !tr.IsPlaying() {
		tr.Resume()
	}
}
func (t *SoundTrack) FadeOut() {
	go func(){
		<-t.waitingForFading
		tr, ok := t.Tracks[t.current]
		if ok {
			FadeOut(tr, STANDARD_FADE_TIME, t.waitingForFading)
		}
	}()
}


func (t *SoundTrack) GetRandomTrack(seed int64) string {
	rand.Seed(seed)
	idx := rand.Intn(len(t.Tracks))
	counter := 0
	for st,_ := range(t.Tracks) {
		if counter == idx {
			return st
		}
		counter ++
	}
	return ""
}