package GE

import (
	"math/rand"
	"time"
	"fmt"
)

const STANDARD_FADE_TIME = 1.5

type SoundTrack struct {
	Tracks map[string]*Sounds
	current string
	waitingForFading chan bool
	waitingCount, maximumWaitingLength int
	OnRepeat func() bool
	NextTrack string
}

func LoadSoundTrack(path string, maximumWaitingLength int) (*SoundTrack, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	files, err := OSReadDir(path)
	if err != nil {return nil, err}
	s := &SoundTrack{maximumWaitingLength:maximumWaitingLength}; s.Tracks = make(map[string]*Sounds)
	for _,f := range(files) {
		s.Tracks[f], err = LoadSounds(path+f)
		if err != nil {return nil, err}
		s.Tracks[f].OnRepeat = s.onTrackRepeat
	}
	s.waitingForFading = make(chan bool)
	go func(){
		s.waitingForFading <- true
	}()
	return s,nil
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
func (t *SoundTrack) onTrackRepeat() (repeat bool) {
	fmt.Println("Repeating track?")
	repeat = true
	if t.OnRepeat != nil {
		repeat = t.OnRepeat()
	}
	_, ok := t.Tracks[t.NextTrack]
	if ok && t.NextTrack != t.current {
		t.Tracks[t.NextTrack].FadeToR(time.Now().UnixNano(), STANDARD_FADE_TIME, nil)
		t.NextTrack = ""
		return false
	}
	return repeat
}
func (t *SoundTrack) Play(name string) {
	if name == t.current {
		return
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
		if _,ok := t.Tracks[name]; ok {
			if _,ok := t.Tracks[t.current]; ok {
				t.Tracks[t.current].FadeToSoundR(t.Tracks[name], time.Now().UnixNano(), t.waitingForFading)
				t.current = name
			}else{
				t.current = name
				t.Tracks[t.current].FadeToR(time.Now().UnixNano(), STANDARD_FADE_TIME, t.waitingForFading)
			}
		}
	}()
}
func (t *SoundTrack) SetVolume(volume float64) {
	for _,s := range t.Tracks {
		s.SetVolume(volume)
	}
}
func (t *SoundTrack) Pause() {
	t.Tracks[t.current].PauseAll()
}
func (t *SoundTrack) Resume() {
	t.Tracks[t.current].Resume()
}
func (t *SoundTrack) FadeOut() {
	go func(){
		<-t.waitingForFading
		t.Tracks[t.current].FadeOut(STANDARD_FADE_TIME, t.waitingForFading)
	}()
}