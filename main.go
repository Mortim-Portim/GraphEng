package main

import (
	//"os"
	//"image"
	//_ "image/jpeg"
	//"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	"marvin/GraphEng/res"
	//"github.com/hajimehoshi/ebiten/text"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	//"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	//"github.com/golang/freetype/truetype"
	//"github.com/hajimehoshi/ebiten/inpututil"
	//"golang.org/x/image/font"
	//"github.com/nfnt/resize"
	"fmt"
	"image/color"
	"log"
	"time"
	//"math"
)

const (
	screenWidth  = 1600
	screenHeight = 900
	FPS = 30
)

const TestText = "Licht ist eine Form der elektromagnetischen Strahlung. Im engeren Sinne sind vom gesamten elektromagnetischen Spektrum nur die Anteile gemeint, die für das menschliche Auge sichtbar sind. Im weiteren Sinne werden auch elektromagnetische Wellen kürzerer Wellenlänge (Ultraviolett) und größerer Wellenlänge (Infrarot) dazu gezählt.\nDie physikalischen Eigenschaften des Lichts werden durch verschiedene Modelle beschrieben: In der Strahlenoptik wird die geradlinige Ausbreitung des Lichts durch „Lichtstrahlen“ veranschaulicht; in der Wellenoptik wird die Wellennatur des Lichts betont, wodurch auch Beugungs- und Interferenzerscheinungen erklärt werden können. In der Quantenphysik schließlich wird das Licht als ein Strom von Quantenobjekten, den Photonen (veranschaulichend auch „Lichtteilchen“ genannt), beschrieben. Eine vollständige Beschreibung des Lichts bietet die Quantenelektrodynamik. Im Vakuum breitet sich Licht mit der konstanten Lichtgeschwindigkeit von 299.792.458 m/s aus. Trifft Licht auf Materie, so kann es gestreut, reflektiert, gebrochen und verlangsamt oder absorbiert werden.\nLicht ist der für das menschliche Auge adäquate Sinnesreiz. Dabei wird die Intensität des Lichts als Helligkeit wahrgenommen, die spektrale Zusammensetzung als Farbe."

type TestGame struct {
	Tbv *GE.TabView
	Tbv_U GE.UpdateFunc
	Tbv_D GE.DrawFunc

	wrld *GE.WorldStructure
	
	frame int
}
func (g *TestGame) Init(screen *ebiten.Image) {
	g.Tbv_U, g.Tbv_D = g.Tbv.Init(screen, nil)
}
func (g *TestGame) Update(screen *ebiten.Image) error {
	if g.frame == 0 {
		g.Init(screen)
	}else{
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
		g.wrld.StructureObjs[0].SetToXY(int16(x),int16(y))
		
		g.Tbv.Update(g.frame)
		
		g.Tbv.Screens.Member[3].(*GE.TabView).Screens.Member[0].(*GE.Animation).UpdatePeriod = g.Tbv.Screens.Member[3].(*GE.TabView).Screens.Member[2].(*GE.ScrollBar).Current()
		g.Tbv.Draw(screen)
		
		g.wrld.Draw(screen)
		
		//g.wrld.UpdateCollisionMat()
		//fmt.Println(g.wrld.CollisionMat.Print())
	}
	g.frame ++
	return nil
}
func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	GE.Init("")
	GE.StandardFont = GE.ParseFontFromBytes(res.MONO_TTF)

	formatedTestText := GE.FormatTextToWidth(TestText, 21, true)

	edT := GE.GetEditText("Fett", 10, 100, 90, 20, GE.StandardFont, &color.RGBA{255, 120, 20, 255}, GE.EditText_Selected_Col)

	btn := GE.GetTextButton("Edit", "Edasdit", GE.StandardFont, 10, 100, 90, &color.RGBA{255, 0, 0, 255}, &color.RGBA{0, 0, 255, 255})
	btn.RegisterOnEvent(func(b *GE.Button) {
		if !b.LPressed {
			edT.IsSelected = !edT.IsSelected
		}
	})

	TextView := GE.GetTextView(formatedTestText, 0, 300, 40, 3, GE.StandardFont, &color.RGBA{255, 255, 255, 255}, &color.RGBA{255, 0, 0, 255})

	ScrollBar := GE.GetStandardScrollbar(700, 500, 600, 60, 0, 120, 3, GE.StandardFont)
	
	spellEimg,_ := GE.LoadEbitenImgFromBytes(res.SPELL_ANIM)
	animation := GE.GetAnimation(1000, 300, 160, 240, 28, 3, spellEimg)

	up2data := make([]GE.UpdateAble, 3)
	up2data[0] = animation
	up2data[1] = TextView
	up2data[2] = ScrollBar
	
	eTab1,_ := GE.LoadEbitenImgFromBytes(res.TAB1)
	eTab2,_ := GE.LoadEbitenImgFromBytes(res.TAB2)
	eTab3,_ := GE.LoadEbitenImgFromBytes(res.TAB3)
	params2 := &GE.TabViewParams{Imgs: []*ebiten.Image{eTab1, eTab2, eTab3}, Scrs: up2data, Y: 200, W: screenWidth, H: screenHeight}
	tbv2 := GE.GetTabView(params2)

	updatable := make([]GE.UpdateAble, 4)
	updatable[0] = edT
	updatable[1] = btn
	updatable[2] = TextView
	updatable[3] = tbv2
	params := &GE.TabViewParams{Nms: []string{"Fett", "Sack", "Fettsack", "LOL"}, Scrs: updatable, W: screenWidth, H: screenHeight}
	tbv := GE.GetTabView(params)

	wmatI := GE.GetMatrix(20, 20, 0)
		
	wmatL :=  GE.GetMatrix(20, 20, 0)
	//wmatL.Fill(0,0,19,19, -10)
	//wmatL.Fill(0,0,18,18, -8)
	wmatL.Fill(0,0,16,16, -6)
	wmatL.Fill(0,0,14,14, -4)
	wmatL.Fill(0,0,12,12, -2)
	wmatL.Fill(0,0,10,10, 0)
	wmatL.Fill(0,0,8,8, 2)
	wmatL.Fill(0,0,6,6, 4)
	wmatL.Fill(0,0,4,4, 6)
	wmatL.Fill(0,0,2,2, 8)
	wmatL.Fill(0,0,0,0, 10)
	
	fmt.Println(wmatL.Print())
	fmt.Println(wmatL.SubMatrix(3,2,15,12).Print())
	
	startComp := time.Now()
	err1 := wmatL.Save("./res/wmatL.txt")
	if err1 != nil {
		panic(err1)
	}
	endComp := time.Now()
	fmt.Println("Compressing took: ", endComp.Sub(startComp))
	
	startDeComp := time.Now()
	err2 := wmatL.Load("./res/wmatL.txt")
	if err2 != nil {
		panic(err2)
	}
	endDeComp := time.Now()
	fmt.Println("Decompressing took: ", endDeComp.Sub(startDeComp))

	wrld := GE.GetWorldStructure(900, 300, 500, 500, 9, 9)
	wrld.IdxMat = wmatI; wrld.LayerMat = wmatL
	
	eTile16,_ := GE.LoadEbitenImgFromBytes(res.TILE16)
	wrld.AddTile(&GE.Tile{eTile16})
	
	treeImg,_ := GE.LoadEbitenImg("./res/jump.png")
	tree := GE.GetAnimation(0, 0, 16, 32, 16, 5, treeImg)
	strct := GE.GetStructureObj(tree, GE.GetRectangle(1,1,1,1), 16, true)
	
	wrld.AddStructureObj(strct)
	wrld.GetFrame(2, 90)
	//wrld.SetMiddle(&GE.Point{1,1})
	wrld.SetMiddle(&GE.Point{4,4})
	
	errS := wrld.Save("./res/wrld.txt")
	if errS != nil {
		panic(errS)
	}
	errL := wrld.Load("./res/wrld.txt")
	if errL != nil {
		panic(errL)
	}

	g := &TestGame{tbv, nil, nil, wrld, 0}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GameEngine Test")
	//ebiten.SetFullscreen(true)
	//ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetMaxTPS(FPS)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
