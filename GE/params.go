package GE

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"image/color"
)

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
	MoveOnButtonDown = 1.0/20.0
	
	
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

//stores a string an if possible a float64 value for each key
type Params struct {
	p map[string]float64
	strs map[string]string
}

//Loads params from a file
func (p *Params) LoadFromFile(path string) error {
	f, err := os.Open(path)
	CheckErr(err)
	scanner := bufio.NewScanner(f)
	p.p = make(map[string]float64)
	p.strs = make(map[string]string)
    for scanner.Scan() {
    	line := scanner.Text()
    	line = strings.ReplaceAll(line, " ", "")
    	ps := strings.Split(line, ":")
    	if len(ps) >= 2 {
    		fl, err2 := strconv.ParseFloat(ps[1], 64)
    		p.strs[ps[0]] = ps[1]
    		if err2 == nil {
	    		p.p[ps[0]] = fl
    		}
    	}
    }
	return nil
}
//returns the string value for the key
func (p *Params) GetS(key string) (string) {
	if val, ok := p.strs[key]; ok {
	    return val
	}
	return ""
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