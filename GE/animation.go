package GE

import (
	"github.com/hajimehoshi/ebiten"
	"image"
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
	sprites, current, spriteWidth, spriteHeight, UpdatePeriod int
	
	spriteSheet *ebiten.Image
}
func (a *Animation) Clone() (*Animation) {
	return &Animation{a.ImageObj, a.sprites, a.current, a.spriteWidth, a.spriteHeight, a.UpdatePeriod, a.spriteSheet}
}
func (a *Animation) Init(screen *ebiten.Image, data interface{}) (UpdateFunc, DrawFunc) {
	return a.Update, a.DrawImageObj
}
func (a *Animation) Start(screen *ebiten.Image, data interface{}) {}
func (a *Animation) Stop(screen *ebiten.Image, data interface{}) {}
func (a *Animation) Update(frame int) {
	if a.UpdatePeriod > 0 && frame%a.UpdatePeriod == 0 {
		a.current ++
		if a.current >= a.sprites || frame == 0 {
			a.current = 0
		}
		a.Img = a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image)
	}else if a.UpdatePeriod == 0 && a.Img == nil {
		a.Img = a.spriteSheet.SubImage(image.Rect(a.spriteWidth*a.current, 0, a.spriteWidth*(a.current+1), a.spriteHeight)).(*ebiten.Image)
	}
}
