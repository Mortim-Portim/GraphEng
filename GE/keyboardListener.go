package GE

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
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

var MOVE_A_D = false
var counter = 0
var MOVING_A = false

func CheckForAutoMove(key ebiten.Key) bool {
	if !MOVE_A_D {
		return false
	}
	if (key == ebiten.KeyA && MOVING_A) || (key == ebiten.KeyD && !MOVING_A) {
		counter++
		if counter > 20 {
			counter = 0
			MOVING_A = !MOVING_A
		}
		return true
	}
	return false
}

//KeyboardListener
type KeyLi struct {
	mapper    map[int]int
	keyStates map[int]bool

	mL, sL      sync.Mutex
	JustChanged []int

	SettingKey int

	EventListeners map[int][]func(l *KeyLi, state bool)
}

func (l *KeyLi) MappIDToKey(id int, key ebiten.Key) {
	oid := GetKeyID(key)
	l.mL.Lock()
	l.mapper[id] = oid
	l.mL.Unlock()
}
func (l *KeyLi) MappKey(key ebiten.Key) int {
	id := GetKeyID(key)
	l.mL.Lock()
	l.mapper[id] = id
	l.mL.Unlock()
	return id
}

//Registers a listener for an event of a single key
func (l *KeyLi) RegisterKeyEventListener(keyID int, listener func(*KeyLi, bool)) {
	l.mL.Lock()
	if _, ok := l.mapper[keyID]; !ok {
		l.mapper[keyID] = keyID
	}
	_, ok := l.EventListeners[l.mapper[keyID]]
	if !ok {
		l.EventListeners[l.mapper[keyID]] = make([]func(l *KeyLi, state bool), 0)
	}
	l.EventListeners[l.mapper[keyID]] = append(l.EventListeners[l.mapper[keyID]], listener)
	l.mL.Unlock()
}

//Resets the configurations
func (l *KeyLi) Reset() {
	l.mL.Lock()
	l.mapper = make(map[int]int)
	l.mL.Unlock()
	l.sL.Lock()
	l.keyStates = make(map[int]bool)
	l.sL.Unlock()
	l.JustChanged = make([]int, 0)
	l.SettingKey = -1
	l.EventListeners = make(map[int][]func(l *KeyLi, state bool))
}

//Update only the Keys that are used
func (l *KeyLi) UpdateMapped() error {
	l.JustChanged = make([]int, 0)
	l.mL.Lock()
	for _, ID := range l.mapper {
		l.UpdateKeyState(ID)
	}
	l.mL.Unlock()
	return nil
}

//Update all Keys
func (l *KeyLi) Update() {
	l.JustChanged = make([]int, 0)
	for ID := range AllKeys {
		l.UpdateKeyState(ID)
	}
}

//Update the state of a specific Key
func (l *KeyLi) UpdateKeyState(KeyID int) {
	l.sL.Lock()
	lastKeyState, ok := l.keyStates[KeyID]
	l.sL.Unlock()
	if !ok {
		lastKeyState = false
	}
	l.sL.Lock()
	l.keyStates[KeyID] = false
	if ebiten.IsKeyPressed(AllKeys[KeyID]) || CheckForAutoMove(AllKeys[KeyID]) {
		l.keyStates[KeyID] = true
	}
	l.sL.Unlock()
	if lastKeyState != l.keyStates[KeyID] && !containsI(l.JustChanged, KeyID) {
		l.JustChanged = append(l.JustChanged, KeyID)

		listeners, ok := l.EventListeners[KeyID]
		if ok {
			for _, listener := range listeners {
				listener(l, l.keyStates[KeyID])
			}
		}

		if l.SettingKey >= 0 {
			l.mL.Lock()
			l.mapper[l.SettingKey] = KeyID
			l.mL.Unlock()
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
	l.sL.Lock()
	state = l.keyStates[KeyID]
	l.sL.Unlock()
	change = containsI(l.JustChanged, KeyID)
	return
}

//Returns the state and weather it just changed based on the Keys ID
func GetRawKeyStateFast(KeyID int) (state, change bool) {
	key := AllKeys[KeyID]
	return GetKeyStateFast(key)
}

//Returns the state and weather it just changed based on the Keys ID
func GetKeyStateFast(key ebiten.Key) (state, change bool) {
	state = ebiten.IsKeyPressed(key)
	change = IsKeyJustDown(key)
	return
}

//Returns the state and weather it just changed based on the Keys mapped ID
func (l *KeyLi) GetMappedKeyState(IDs ...int) (state, change bool) {
	state = true
	for _, ID := range IDs {
		l.mL.Lock()
		KeyID, ok := l.mapper[ID]
		l.mL.Unlock()
		if !ok {
			KeyID = ID
		}
		s, c := l.GetRawKeyState(KeyID)
		if !s {
			state = false
		}
		if c {
			change = true
		}
	}
	return
}

//Returns the state and weather it just changed based on the Keys mapped ID
func (l *KeyLi) GetMappedKeyStateFast(ID int) (state, change bool) {
	l.mL.Lock()
	KeyID, ok := l.mapper[ID]
	l.mL.Unlock()
	if !ok {
		KeyID = ID
	}
	return GetRawKeyStateFast(KeyID)
}

//Saves the Keyboardmapper to a file
func (l *KeyLi) SaveConfig(path string) {
	l.mL.Lock()
	SaveMapper(path, l.mapper)
	l.mL.Unlock()
}

//Loads the Keyboardmapper from a file
func (l *KeyLi) LoadConfig(path string) {
	mapper := LoadMapper(path)
	if mapper != nil && len(mapper) > 0 {
		l.mL.Lock()
		for i, k := range mapper {
			l.mapper[i] = k
		}
		l.mL.Unlock()
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

//Returns the KeyID of an specified ebiten key
func GetKeyID(key ebiten.Key) int {
	for i, k := range AllKeys {
		if int(k) == int(key) {
			return i
		}
	}
	return -1
}
