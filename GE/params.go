package GE

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"errors"
	"image/color"
	"io/ioutil"
	"time"
)
func GetTime() string {
	jned := strings.Join(strings.Split(fmt.Sprintf("%v", time.Now().UTC()), " ")[:2], "_")
	return strings.ReplaceAll(strings.Split(jned, ".")[0], ":", "-")
}
/**
Params represents a struct that can be loaded from a file, that contains paramters

The Parameters of that file can than be accessed with params.Get(string) or params.GetS(string)

The file to load should look like this:

ParameterName1		:		value
ParameterName2		:		value
ParameterName3		:		value
ParameterName4		:		value

floats should be written like this: 3189.89
**/

// TODO Load the following Params from directory
var (
	EditText_Placeholder_Col = &(color.RGBA{255,255,255,100})
	EditText_Selected_Col = &(color.RGBA{180,180,180,255})
	EditText_Back_Col = &(color.RGBA{40,40,40,200})
	
	ReduceColOnButtonDown = 100
	MoveOnButtonDown = 0.0//1.0/20.0
	
	
	TabBack_Col = &(color.RGBA{255,255,255,255})
	TabText_Col = &(color.RGBA{0,0,0,255})
	TabsDistance = 0.0
	TabsHeight = 1.0/15.0
)

//Initializes Params
func InitParams(p *Params) {
	if p == nil {
		return
	}
	ReduceColOnButtonDown = int(p.Get("ReduceColOnButtonDown"))
	MoveOnButtonDown = 		(p.Get("MoveOnButtonDown"))
	TabsDistance = 			(p.Get("TabsDistance"))
	TabsHeight = 			(p.Get("TabsHeight"))
	
	EditText_Placeholder_Col = 	&color.RGBA{uint8(p.Get("EditText_Placeholder_Col_R")),uint8(p.Get("EditText_Placeholder_Col_G")),uint8(p.Get("EditText_Placeholder_Col_B")),uint8(p.Get("EditText_Placeholder_Col_A"))}
	EditText_Selected_Col = 	&color.RGBA{uint8(p.Get("EditText_Selected_Col_R")),uint8(p.Get("EditText_Selected_Col_G")),uint8(p.Get("EditText_Selected_Col_B")),uint8(p.Get("EditText_Selected_Col_A"))}
	EditText_Back_Col = 		&color.RGBA{uint8(p.Get("EditText_Back_Col_R")),uint8(p.Get("EditText_Back_Col_G")),uint8(p.Get("EditText_Back_Col_B")),uint8(p.Get("EditText_Back_Col_A"))}
	TabBack_Col = 				&color.RGBA{uint8(p.Get("TabBack_Col_R")),uint8(p.Get("TabBack_Col_G")),uint8(p.Get("TabBack_Col_B")),uint8(p.Get("TabBack_Col_A"))}
	TabText_Col = 				&color.RGBA{uint8(p.Get("TabText_Col_R")),uint8(p.Get("TabText_Col_G")),uint8(p.Get("TabText_Col_B")),uint8(p.Get("TabText_Col_A"))}
}

type List struct {
	strs []string
}
func (l *List) LoadFromFile(path string) error {
	f, err := os.Open(path)
	CheckErr(err)
	scanner := bufio.NewScanner(f)
	l.strs = make([]string, 0)
    for scanner.Scan() {
    	line := scanner.Text()
    	if len(strings.ReplaceAll(line, " ", "")) > 0 {
	    	l.strs = append(l.strs, line)
    	}
    }
	return nil
}
func (l *List) GetSlice() []string {
	return l.strs
}
func (l *List) Get(idx int) (string, error) {
	if idx >= 0 && idx < len(l.strs) {
		return l.strs[idx], nil
	}
	return "", errors.New(fmt.Sprintf("Index %v out of range with slice of length %v", idx, len(l.strs)))
}
func (l *List) Print() (out string) {
	out = ""
	for i,str := range(l.strs) {
		out += fmt.Sprintf("%v: %s\n", i, str)
	}
	out = out[:len(out)-2]
	return
}

//stores a string an if possible a float64 value for each key
type Params struct {
	p map[string]float64
	strs map[string]string
}

//Loads params from a file
func (p *Params) LoadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {return err}
	scanner := bufio.NewScanner(f)
	p.p = make(map[string]float64)
	p.strs = make(map[string]string)
    for scanner.Scan() {
    	line := scanner.Text()
    	line = strings.ReplaceAll(line, " ", "")
    	ps := strings.Split(line, ":")
    	if len(ps) >= 2 {
    		fl, err2 := strconv.ParseFloat(ps[1], 64)
    		p.strs[ps[0]] = strings.Join(ps[1:], ":")
    		//fmt.Printf("Loading Param '%s': %v = '%s'\n", ps[0], ps[1:], strings.Join(ps[1:], ":"))
    		if err2 == nil {
	    		p.p[ps[0]] = fl
    		}
    	}
    }
	return nil
}
//Saves params to a file
func (p *Params) SaveToFile(path string) error {
	data := ""
	for key,val := range(p.strs) {
		data += fmt.Sprintf("%s:%s\n", key, val)
	}
	if len(data) == 0 {
		return nil
	}
	data = data[:len(data)-1]
	ioutil.WriteFile(path, []byte(data), 0644)
	return nil
}
//returns the string value for the key
func (p *Params) GetS(key string) (string) {
	if val, ok := p.strs[key]; ok {
	    return val
	}
	return ""
}
func (p *Params) SetS(key, val string) {
	p.strs[key] = val
	fl, err2 := strconv.ParseFloat(val, 64)
    if err2 == nil {
	    p.p[key] = fl
    }else{
	    p.p[key] = 0.0
    }
}
func (p *Params) Set(key string, val float64) {
	p.p[key] = val
	p.strs[key] = fmt.Sprintf("%0.8f",val)
}
func (p *Params) SetBool(key string, val bool) {
	if val {
		p.SetS(key, "true")
	}else{
		p.SetS(key, "false")
	}
}
//returns the boolean value
func (p *Params) GetBool(key string, standard bool) bool {
	val := strings.ToLower(p.GetS(key))
	if val == "false" || val == "0" {
		return false
	}
	if val == "true" || val == "1" {
		return true
	}
	return standard
}
//returns the string value for the key
func (p *Params) Get(key string) (float64) {
	if val, ok := p.p[key]; ok {
	    return val
	}
	return 0
}
//Prints the paramters an values
func (p *Params) Print() string {
	out := ""
	for k,v := range(p.strs) {
		out += fmt.Sprintf("%s : %s, %0.4f\n", k, v, p.p[k])
	}
	return out
}