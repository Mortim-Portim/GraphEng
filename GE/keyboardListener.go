package GE

import (
	"github.com/hajimehoshi/ebiten"
	"fmt"
	//"log"
	//"math"
	//"errors"
	"encoding/json"
	"io/ioutil"
)
/**
KeyboardListener represents a struct, that listens for KeyboardEvents

The Key is maped, so that a user can specifiy a Key of his choice to use for an action

Example:
W is maped to index 0, which means Run forward
if SetKeyState(0) is called the next Key the user presses is mapped to i
for Example:
SetKeyState(0) is called
User presses L_Shift
L_Shift is maped to index 0, which means Run forward

when SetKeyState(0) was called the all Keys should be updated(Update), otherwise
only update used keys, increasing performance (UpdateMapped)
**/

//KeyboardListener
type KeyLi struct {
	mapper map[int]int
	
	keyStates map[int]bool
	JustChanged []int
	
	SettingKey int
}
//Resets the configurations
func (l *KeyLi) Reset() {
	l.mapper = make(map[int]int)
	l.keyStates = make(map[int]bool)
	l.JustChanged = make([]int, 0)
	l.SettingKey = -1
}
//Update only the Keys that are used
func (l *KeyLi) UpdateMapped(screen *ebiten.Image, frame int) error {
	l.JustChanged = make([]int, 0)
	for _,ID := range(l.mapper) {
		l.UpdateKeyState(ID)
	}
	return nil
}
//Update all Keys
func (l *KeyLi) Update() {
	l.JustChanged = make([]int, 0)
	for ID,_ := range(AllKeys) {
		l.UpdateKeyState(ID)
	}
}
//Update the state of a specific Key
func (l *KeyLi) UpdateKeyState(KeyID int) {
	lastKeyState, ok := l.keyStates[KeyID]
	if !ok {
		lastKeyState = false
	}
	l.keyStates[KeyID] = false
	if ebiten.IsKeyPressed(AllKeys[KeyID]) {
		l.keyStates[KeyID] = true
	}
	if lastKeyState != l.keyStates[KeyID] && !contains(l.JustChanged, KeyID) {
		l.JustChanged = append(l.JustChanged, KeyID)
		if l.SettingKey >= 0 {
			l.mapper[l.SettingKey] = KeyID
			l.SettingKey = -1
		}
	}
}
//Sets the Key to be reassigned
func (l *KeyLi) SetKeyState(KeyID int) {
	l.SettingKey = KeyID
}

//Returns all KeyIDs that just changed (use GetRawKeyState(ID))
func (l *KeyLi) GetJustChangedKeys() (IDs []int) {
	return l.JustChanged
}
//Returns the state and weather it just changed based on the Keys ID
func (l *KeyLi) GetRawKeyState(KeyID int) (state, change bool) {
	state = l.keyStates[KeyID]
	change = contains(l.JustChanged, KeyID)
	return
}
//Returns the state and weather it just changed based on the Keys mapped ID
func (l *KeyLi) GetMappedKeyState(ID int) (state, change bool) {
	KeyID, ok := l.mapper[ID]
	if !ok {
		KeyID = ID
	}
	return l.GetRawKeyState(KeyID)
}
//Saves the Keyboardmapper to a file
func (l *KeyLi) SaveConfig(path string) {
	SaveMapper(fmt.Sprintf("%s/Keyboardmapper.txt", path), l.mapper)
}
//Loads the Keyboardmapper from a file
func (l *KeyLi) LoadConfig(path string) {
	mapper := LoadMapper(fmt.Sprintf("%s/Keyboardmapper.txt", path))
	if mapper != nil {
		l.mapper = mapper
	}
}
//Loads a map[int]int from a file
func LoadMapper(path string) map[int]int {
	dat, err := ioutil.ReadFile(path)
   	if err != nil {
	   	return nil
   	}
	var newMapper map[int]int
	err2 := json.Unmarshal(dat, &newMapper)
	if err2 != nil {
	   	return nil
   	}
	return newMapper
}
//Saves a map[int]int to a file
func SaveMapper(path string, mapper map[int]int) {
	bytes, err := json.Marshal(mapper)
	CheckErr(err)
	err2 := ioutil.WriteFile(path, bytes, 0644)
    CheckErr(err2)
}