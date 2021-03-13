package GE

import (
	"image"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

/**
Animation struct to load and play animations

all frames are contained in a spriteSheet from the left to the right
  #   |  #   |  #   |  #   |  #   |
 # #  | # #  | # #  | # #  | # #  |
  #   |  #   |# #   |  #   |  #   |
 ###  |####  | #### | ###  | ###  |
# # # |  # # |  #   |# # # |# # # |
  #   |  #   |  # # |  #   |  #   |
 # #  | # ## | # #  | # #  | # #  |
#   # |#     |#     |#   # |#   # |

Animation implements UpdateAble
**/
type Animation struct {
	ImageObj
	sprites, current, spriteWidth, spriteHeight, UpdatePeriod, lastFrame int

	spriteSheet *ebiten.Image
}

func (a *Animation) Clone() *Animation {
	return &Animation{*a.ImageObj.Copy(), a.sprites, a.current, a.spriteWidth, a.spriteHeight, a.UpdatePeriod, a.lastFrame, a.spriteSheet}
}
func (a *Animation) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return a.Update, a.DrawImageObj
}
func (a *Animation) Start(screen *ebiten.Image, data interface{}) {}
func (a *Animation) Stop(screen *ebiten.Image, data interface{})  {}
func (a *Animation) Update(frame int) {
	if a.UpdatePeriod > 0 && a.lastFrame != frame && frame%a.UpdatePeriod == 0 {
		a.lastFrame = frame
		a.current++
		if a.current >= a.sprites || frame == 0 {
			a.current = 0
		}
		a.Img = a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image)
	} else if a.UpdatePeriod == 0 && a.Img == nil {
		a.Img = a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image)
	}
}
func (a *Animation) SetTo(frame int) {
	frame = frame % a.sprites
	if a.current != frame {
		a.current = frame
		a.Img = a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image)
	}
}
