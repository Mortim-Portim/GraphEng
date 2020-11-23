package GE

import (
	"time"
)

const STANDARD_FADE_TIME = 1.5

type SoundTrack struct {
	Tracks map[string]*Sounds
	current string
	waitingForFading chan bool
	waitingCount, maximumWaitingLength int
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
	}
	s.waitingForFading = make(chan bool)
	go func(){
		s.waitingForFading <- true
	}()
	return s,nil
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