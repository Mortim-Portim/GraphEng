package main

import (
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	"marvin/GraphEng/res"
	"fmt"
	"log"
	"time"
)
const (
	screenWidth  = 1600
	screenHeight = 900
	FPS = 30
	TestText = "Licht ist eine Form der elektromagnetischen Strahlung. Im engeren Sinne sind vom gesamten elektromagnetischen Spektrum nur die Anteile gemeint, die für das menschliche Auge sichtbar sind. Im weiteren Sinne werden auch elektromagnetische Wellen kürzerer Wellenlänge (Ultraviolett) und größerer Wellenlänge (Infrarot) dazu gezählt.\nDie physikalischen Eigenschaften des Lichts werden durch verschiedene Modelle beschrieben: In der Strahlenoptik wird die geradlinige Ausbreitung des Lichts durch „Lichtstrahlen“ veranschaulicht; in der Wellenoptik wird die Wellennatur des Lichts betont, wodurch auch Beugungs- und Interferenzerscheinungen erklärt werden können. In der Quantenphysik schließlich wird das Licht als ein Strom von Quantenobjekten, den Photonen (veranschaulichend auch „Lichtteilchen“ genannt), beschrieben. Eine vollständige Beschreibung des Lichts bietet die Quantenelektrodynamik. Im Vakuum breitet sich Licht mit der konstanten Lichtgeschwindigkeit von 299.792.458 m/s aus. Trifft Licht auf Materie, so kann es gestreut, reflektiert, gebrochen und verlangsamt oder absorbiert werden.\nLicht ist der für das menschliche Auge adäquate Sinnesreiz. Dabei wird die Intensität des Lichts als Helligkeit wahrgenommen, die spektrale Zusammensetzung als Farbe."
)
func StartGame(g ebiten.Game) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GameEngine Test")
	//ebiten.SetFullscreen(true)
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetMaxTPS(FPS)
	if err := ebiten.RunGame(g); err != nil {
		GE.CloseLogFile()
		log.Fatal(err)
	}
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

import "image/color"
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
	frame int
}
func (g *TestGame) Init(screen *ebiten.Image) {}
func (g *TestGame) Update(screen *ebiten.Image) error {
	if g.frame%2 == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.wrld.Move(-1,0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.wrld.Move(1,0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.wrld.Move(0,-1)
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.wrld.Move(0,1)
		}
	}
	x,y := g.wrld.Middle()
	g.wrld.Objects[0].SetToXY(float64(x),float64(y))
	g.wrld.DrawBack(screen)
	g.frame ++
	
	msg := fmt.Sprintf(`TPS: %0.2f`, ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
	GE.LogToFile(msg+"\n")
	
	return nil
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	GE.Init("")
	GE.StandardFont = GE.ParseFontFromBytes(res.MONO_TTF)

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates the index Matrix
	wmatI := GE.GetMatrix(20, 20, 0)
	//Creates the light Matrix
	wmatL :=  GE.GetMatrix(20, 20, 0)
	wmatL.Fill(0,0,19,19, 0)
	wmatL.Fill(0,0,17,17, 40)
	wmatL.Fill(0,0,15,15, 70)
	wmatL.Fill(0,0,13,13, 100)
	wmatL.Fill(0,0,11,11, 130)
	wmatL.Fill(0,0,9,9, 160)
	wmatL.Fill(0,0,7,7, 190)
	wmatL.Fill(0,0,5,5, 220)
	wmatL.Fill(0,0,3,3, 255)
	//Prints the full layer matrix
	fmt.Println("wmatL width: ", wmatL.WAbs())
	//Prints a submatrix of the layer matrix
	fmt.Println(wmatL.SubMatrix(3,2,15,12).Print())
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Saves the matrix in a compressed form to the file system
	err1 := wmatL.Save("./res/wmatL.txt")
	if err1 != nil {
		panic(err1)
	}
	
	//Loads the matrix from the file system
	err2 := wmatL.Load("./res/wmatL.txt")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("wmatL width: ", wmatL.WAbs())

	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a WorldStructure object
	wrld := GE.GetWorldStructure(0, 0, 1600, 900, 16, 9)
	wrld.TileMat = wmatI; wrld.LightMat = wmatL
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//Creates a raster
	wrld.GetFrame(2, 90)
	//Sets the start point
	wrld.SetMiddle(10,10)
	
	//Saves the compressed world
	startComp := time.Now()
	errS := wrld.Save("./res/wrld.txt")
	if errS != nil {
		panic(errS)
	}
	//Calculates how long it took to save the world
	endComp := time.Now()
	fmt.Println("Saving wrld took: ", endComp.Sub(startComp))
	
	//loads the compressed world
	startDeComp := time.Now()
	errL := wrld.Load("./res/wrld.txt")
	if errL != nil {
		panic(errL)
	}
	//Calculates how long it took to load the world
	endDeComp := time.Now()
	fmt.Println("Loading wrld took: ", endDeComp.Sub(startDeComp))
	
	//----------------------------------------------------------------------------------------------------------------------------------------------
	//loads all tiles
	tiles, errT := GE.ReadTiles("./res/tiles/")
	if errT != nil {
		panic(errT)
	}
	for _,t := range(tiles) {
		fmt.Println(t.Name)
	}
	//loads all objs
	objs, errO := GE.ReadStructures("./res/structObjs/")
	if errO != nil {
		panic(errO)
	}
	for _,o := range(objs) {
		fmt.Println(o.Name)
	}
	wrld.Tiles = tiles
	wrld.AddStruct(objs[0])
	
	g := &TestGame{wrld, 0}

	GE.SetLogFile("./res/log.txt")

	StartGame(g)
}
