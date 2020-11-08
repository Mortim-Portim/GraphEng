package main

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	
	"marvin/GraphEng/res"
	"fmt"
	"time"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	
	//"image/color"
	
	//cmp "marvin/GraphEng/Compression"
)
const (
	screenWidth  = 1920
	screenHeight = 1080
	FPS = 30
	TestText = "Licht ist eine Form der elektromagnetischen Strahlung. Im engeren Sinne sind vom gesamten elektromagnetischen Spektrum nur die Anteile gemeint, die für das menschliche Auge sichtbar sind. Im weiteren Sinne werden auch elektromagnetische Wellen kürzerer Wellenlänge (Ultraviolett) und größerer Wellenlänge (Infrarot) dazu gezählt.\nDie physikalischen Eigenschaften des Lichts werden durch verschiedene Modelle beschrieben: In der Strahlenoptik wird die geradlinige Ausbreitung des Lichts durch „Lichtstrahlen“ veranschaulicht; in der Wellenoptik wird die Wellennatur des Lichts betont, wodurch auch Beugungs- und Interferenzerscheinungen erklärt werden können. In der Quantenphysik schließlich wird das Licht als ein Strom von Quantenobjekten, den Photonen (veranschaulichend auch „Lichtteilchen“ genannt), beschrieben. Eine vollständige Beschreibung des Lichts bietet die Quantenelektrodynamik. Im Vakuum breitet sich Licht mit der konstanten Lichtgeschwindigkeit von 299.792.458 m/s aus. Trifft Licht auf Materie, so kann es gestreut, reflektiert, gebrochen und verlangsamt oder absorbiert werden.\nLicht ist der für das menschliche Auge adäquate Sinnesreiz. Dabei wird die Intensität des Lichts als Helligkeit wahrgenommen, die spektrale Zusammensetzung als Farbe."
)
func StartGame(g ebiten.Game) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GraphEng Test")
	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(FPS)
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
	GE.CloseLogFile()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


/**
//  _    _                 _____       _             __                
// | |  | |               |_   _|     | |           / _|               
// | |  | |___  ___ _ __    | |  _ __ | |_ ___ _ __| |_ __ _  ___ ___  
// | |  | / __|/ _ \ '__|   | | | '_ \| __/ _ \ '__|  _/ _` |/ __/ _ \ 
// | |__| \__ \  __/ |     _| |_| | | | ||  __/ |  | || (_| | (_|  __/ 
//  \____/|___/\___|_|    |_____|_| |_|\__\___|_|  |_| \__,_|\___\___|


//TestGame implements ebiten.Game, USE FOR TESTING ONLY
type TestGame struct {
	Tbv *GE.TabView
	Tbv_U GE.UpdateFunc
	Tbv_D GE.DrawFunc
	frame int
}
func (g *TestGame) Init(screen *ebiten.Image) {
	g.Tbv_U, g.Tbv_D = g.Tbv.Init(screen, nil)
}
func (g *TestGame) Update(screen *ebiten.Image) error {
	if g.frame == 0 {
		g.Init(screen)
	}else{
		g.Tbv.Update(g.frame)
		g.Tbv.Screens.Member[3].(*GE.TabView).Screens.Member[0].(*GE.Animation).UpdatePeriod = g.Tbv.Screens.Member[3].(*GE.TabView).Screens.Member[2].(*GE.ScrollBar).Current()
		g.Tbv.Draw(screen)
	}
	g.frame ++
	return nil
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	//Initializes the Graphics Engine
	GE.Init("")
	//Sets the StandardFont of the Graphics Engine
	GE.StandardFont = GE.ParseFontFromBytes(res.MONO_TTF)

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates an EditText with standard text "Fett"
	editText := GE.GetEditText("Fett", 10, 100, 90, 20, GE.StandardFont, &color.RGBA{255, 120, 20, 255}, GE.EditText_Selected_Col)

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a blue Button with red Text on it
	button := GE.GetTextButton("Edit", "Edasdit", GE.StandardFont, 10, 100, 90, &color.RGBA{255, 0, 0, 255}, &color.RGBA{0, 0, 255, 255})
	//If Button is clicked select editText
	button.RegisterOnEvent(func(b *GE.Button) {
		if !b.LPressed {
			editText.IsSelected = !editText.IsSelected
		}
	})

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//formats the wikipedia-text to lines with a maximum size of 21 runes
	formatedTestText := GE.FormatTextToWidth(TestText, 21, true)
	//Creates a Scrollable TextView displaying 3 lines at the time
	textView := GE.GetTextView(formatedTestText, 0, 300, 40, 3, GE.StandardFont, &color.RGBA{255, 255, 255, 255}, &color.RGBA{255, 0, 0, 255})

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a ScrollBar with min=0, max=120, initvalue=3
	scrollBar := GE.GetStandardScrollbar(700, 500, 600, 60, 0, 120, 3, GE.StandardFont)
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Loads an ebiten image from bytes
	spellEimg,_ := GE.LoadEbitenImgFromBytes(res.SPELL_ANIM)
	//creates an animation from the spriteSheet
	animation := GE.GetAnimation(1000, 300, 160, 240, 28, 3, spellEimg)

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a slice of UpdateAbles to use for the inner TabView
	innerTabViewUpdateAble := make([]GE.UpdateAble, 3)
	innerTabViewUpdateAble[0] = animation
	innerTabViewUpdateAble[1] = textView
	innerTabViewUpdateAble[2] = scrollBar
	//Loads the inner TabViews tab Buttons
	eTab1,_ := GE.LoadEbitenImgFromBytes(res.TAB1)
	eTab2,_ := GE.LoadEbitenImgFromBytes(res.TAB2)
	eTab3,_ := GE.LoadEbitenImgFromBytes(res.TAB3)
	//Creates the inner TabView
	params_inner := &GE.TabViewParams{Imgs: []*ebiten.Image{eTab1, eTab2, eTab3}, Scrs: innerTabViewUpdateAble, Y: 200, W: screenWidth, H: screenHeight}
	tabView_inner := GE.GetTabView(params_inner)

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a slice of UpdateAbles to use for the outer TabView
	outerTabViewUpdateAble := make([]GE.UpdateAble, 4)
	outerTabViewUpdateAble[0] = editText
	outerTabViewUpdateAble[1] = button
	outerTabViewUpdateAble[2] = textView
	outerTabViewUpdateAble[3] = tabView_inner
	
	params_outer := &GE.TabViewParams{Nms: []string{"Fett", "Sack", "Fettsack", "LOL"}, Scrs: outerTabViewUpdateAble, W: screenWidth, H: screenHeight}
	tabView_outer := GE.GetTabView(params_outer)
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	g := &TestGame{tabView_outer, nil, nil, 0}

	StartGame(g)
}
**/


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


// __          __        _     _    _____ _                   _                  
// \ \        / /       | |   | |  / ____| |                 | |                 
//  \ \  /\  / /__  _ __| | __| | | (___ | |_ _ __ _   _  ___| |_ _   _ _ __ ___ 
//   \ \/  \/ / _ \| '__| |/ _` |  \___ \| __| '__| | | |/ __| __| | | | '__/ _ \
//    \  /\  / (_) | |  | | (_| |  ____) | |_| |  | |_| | (__| |_| |_| | | |  __/
//     \/  \/ \___/|_|  |_|\__,_| |_____/ \__|_|   \__,_|\___|\__|\__,_|_|  \___|


type TestGame struct {
	wrld *GE.WorldStructure
	rec  *GE.Recorder
	player *GE.WObj
	frame int
}
func (g *TestGame) Init(screen *ebiten.Image) {
	g.wrld.Add_Drawables = g.wrld.Add_Drawables.Add(g.player)
}

var timeTaken int64
var inputT, objT, lightUpT, lightDT, worldDT time.Duration
func (g *TestGame) Update(screen *ebiten.Image) error {
	startTime := time.Now()
	if g.frame%1 == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.wrld.MoveSmooth(-0.2,0,false,false)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.wrld.MoveSmooth(0.2,0, false,false)
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.wrld.MoveSmooth(0,-0.2, false,false)
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.wrld.MoveSmooth(0,0.2, false,false)
		}
	}
	g.wrld.MoveSmooth(-0.01,0,false,false)
	_,dy := ebiten.Wheel()
	if dy != 0 {
		g.wrld.Lights[0].SetMaximumIntesity(g.wrld.Lights[0].GetMaximumIntesity()+int16(dy*10))
	}
	
	x,y := g.wrld.SmoothMiddle()
	g.player.Update(g.frame)
	g.player.SetToXY(float64(x),float64(y))
	g.wrld.UpdateObjDrawables()
	
	if g.frame%100 < 50 {
		for _,strct := range(g.wrld.Structures) {
			strct.IsUnderstood = true
		}
	}else{
		for _,strct := range(g.wrld.Structures) {
			strct.IsUnderstood = false
		}
	}
	
	g.wrld.UpdateLightLevel(1)
	u_lights := g.wrld.UpdateAllLightsIfNecassary()
	
	//Around 8ms
	g.wrld.Draw(screen)
	
	if g.frame == FPS*6 {
		go g.rec.Save("./res/out.mp4")
	}
	g.rec.NextFrame(screen)
	
	g.frame ++
	timeTaken = time.Now().Sub(startTime).Milliseconds()
	fps := ebiten.CurrentTPS()
	msg := fmt.Sprintf(`TPS: %0.2f, Updating took: %v at frame %v, ul:%v`, fps, timeTaken, g.frame-1, u_lights)
	ebitenutil.DebugPrint(screen, msg)
	GE.LogToFile(msg+"\n")
	fmt.Println(msg)
	return nil
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	GE.Init("")
	GE.StandardFont = GE.ParseFontFromBytes(res.MONO_TTF)
	GE.SetLogFile("./res/log.txt")
	time.Sleep(time.Second)
	
	pathFMat := GE.GetMatrix(10,10,0)
	for x := 2; x < 7; x++ {
		pathFMat.Set(x,3,10)
	}
	pathFMat.Set(2,4,10)
	pathFMat.Set(2,5,10)
	fmt.Println(pathFMat.Print())
	
	nanos := 0
	for i := 0; i < 100; i++ {
		stComp := time.Now()
		GE.FindPathMat(pathFMat, [2]int{1,1}, [2]int{8,8}, true)
		nanos += int(time.Now().Sub(stComp).Nanoseconds())
	}
	fmt.Println("Computing A* took: ", time.Duration(nanos/100))
	
	FP := GE.FindPathMat(pathFMat, [2]int{1,1}, [2]int{8,8}, false)
	for _,p := range(FP) {
		pathFMat.Set(p[0],p[1],888)
		fmt.Println(p)
	}
	fmt.Println(pathFMat.Print())
	
	XT := 50; YT := 50
	
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates the index Matrix
	wmatI := GE.GetMatrix(300, 300, 0)
	//Saves the matrix in a compressed form to the file system
	err1 := wmatI.Save("./res/wmatI.txt")
	if err1 != nil {
		panic(err1)
	}
	//Loads the matrix from the file system
	err2 := wmatI.Load("./res/wmatI.txt")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("wmatI width: ", wmatI.WAbs())

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a WorldStructure object
	wrld := GE.GetWorldStructure(0, 0, screenWidth, screenHeight, XT, YT, 32, 18)
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	errtl := wrld.LoadTiles("./res/tiles/")
	if errtl != nil {panic(errtl)}
	errol := wrld.LoadStructureObjs("./res/structObjs/")
	if errol != nil {panic(errol)}
	
	wrld.AddNamedStructureObj("house1", 	10, 2)
	wrld.AddNamedStructureObj("tree2", 		5, 5)
	wrld.AddNamedStructureObj("tree2big", 	14, 14)
	
	wrld.UpdateObjMat()
	
	obs := wrld.ObjectsToBytes()
	fmt.Printf("Objects are %v bytes\n", len(obs))
	wrld.BytesToObjects(obs)
	
	player,err := GE.GetWObjFromPath("./res/anims/test")
	if err != nil {panic(err)}
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Add a light source to the world
	light1 := GE.GetLightSource(&GE.Point{10,10}, &GE.Vector{0,-1,0}, 360, 400, 0.01, false)
	light2 := GE.GetLightSource(&GE.Point{15,15}, &GE.Vector{0,-1,0}, 360, 400, 0.01, false)
	
	wrld.Lights = append(wrld.Lights, light1, light2)//, light2, light3, light4, light5, light6, light7, light8, light9, light10)
	wrld.UpdateLIdxMat()
	wrld.UpdateLightValue(wrld.Lights, true)
	
	lbs := wrld.LightsToBytes()
	fmt.Printf("Lights are %v bytes\n", len(lbs))
	wrld.BytesToLights(lbs)
	
	//Sets the start point
	wrld.SetMiddle(14,14,true)
	wrld.SetLightStats(10,255, 0.3)
	wrld.SetLightLevel(15)
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Saves the compressed world
	startComp := time.Now()
	errS := wrld.Save("./res/wrld.map")
	if errS != nil {
		panic(errS)
	}
	//Calculates how long it took to save the world
	endComp := time.Now()
	fmt.Println("Saving wrld took: ", endComp.Sub(startComp))
	
	//loads the compressed world
	startDeComp := time.Now()
	newWrld, errL := GE.LoadWorldStructure(0,0,screenWidth,screenHeight, "./res/TestMap4.map", "./res/tiles/", "./res/structObjs/")
	if errL != nil {
		GE.ShitImDying(errL)
	}
	//Calculates how long it took to load the world
	endDeComp := time.Now()
	fmt.Println("Loading wrld took: ", endDeComp.Sub(startDeComp))
	
	//Sets the start point
	//newWrld.SetMiddle(14,14,true)
	newWrld.SetLightStats(10,255, 0.3)
	newWrld.SetLightLevel(15)
	//Creates a raster
	//newWrld.GetFrame(2, 90)
	newWrld.SetDisplayWH(32,18)
	
	g := &TestGame{newWrld, GE.GetNewRecorder(FPS*5, 192, 108, FPS), player, 0}	
	g.Init(nil)

	StartGame(g)
}


/**
type TestGame struct {}
func (g *TestGame) Update(screen *ebiten.Image) error {return nil}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 600, 400
}

func main() {
	GE.Init("")
	
	go func(){
		st, err := GE.LoadSoundTrack("./res/audio/Soundtrack")
		GE.ShitImDying(err)
		st.Play("main")
		fmt.Println("Playing Main")
		
		time.Sleep(time.Second*3)
		st.Play("ork")
		fmt.Println("Playing Ork")
		
		time.Sleep(time.Second*3)
		st.Pause()
		fmt.Println("Pausing")
		
		time.Sleep(time.Second*3)
		st.Resume()
		fmt.Println("Resuming")
		
		time.Sleep(time.Second*3)
		st.FadeOut()
		fmt.Println("Fading out")
	}()
	StartGame(&TestGame{})
}
**/