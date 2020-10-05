package GE

import (
	"fmt"
	"math/rand"
)

//IMPORTANT: initialize a random seed (rand.Seed()) before calling PR()

type SoundEffect struct {
	sounds []*AudioPlayer
}
//loads a soundeffect consiting out of multiple audio files
func LoadSoundEffect(folder string) (s *SoundEffect) {
	s = &SoundEffect{}
	s.sounds = make([]*AudioPlayer,0)
	for i := 0; i < 20; i++ {
		file := fmt.Sprintf("%s/%v.wav", folder, i+1)
		fmt.Println("Loading: ", file)
		player, err := NewPlayer(file)
		if err != nil {
			break
		}
		s.sounds = append(s.sounds, player)
	}
	if len(s.sounds) == 0 {
		panic("Could not load SoundEffects")
	}
	return
}
//Plays a random audio file
func (s *SoundEffect) PR() {
	idx := rand.Intn(len(s.sounds))
	s.sounds[idx].Play()
}
