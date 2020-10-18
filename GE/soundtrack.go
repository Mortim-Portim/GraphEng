package GE

import (
	"time"
)

const STANDARD_FADE_TIME = 1.5

type SoundTrack struct {
	Tracks map[string]*Sounds
	current string
}

func LoadSoundTrack(path string) (*SoundTrack, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	files, err := OSReadDir(path)
	if err != nil {return nil, err}
	s := &SoundTrack{}; s.Tracks = make(map[string]*Sounds)
	for _,f := range(files) {
		s.Tracks[f], err = LoadSounds(path+f)
		if err != nil {return nil, err}
	}
	return s,nil
}

func (t *SoundTrack) Play(name string) {
	if _,ok := t.Tracks[name]; ok {
		if _,ok := t.Tracks[t.current]; ok {
			t.Tracks[t.current].FadeToSoundR(t.Tracks[name], time.Now().UnixNano())
			t.current = name
		}else{
			t.current = name
			t.Tracks[t.current].FadeToR(time.Now().UnixNano(), STANDARD_FADE_TIME)
		}
	}
}

func (t *SoundTrack) Pause() {
	t.Tracks[t.current].PauseAll()
}
func (t *SoundTrack) Resume() {
	t.Tracks[t.current].Resume()
}
func (t *SoundTrack) FadeOut() {
	t.Tracks[t.current].FadeOut(STANDARD_FADE_TIME)
}