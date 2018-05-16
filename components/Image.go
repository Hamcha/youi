package components

import (
	"image"
)

// Image is a simple box that can contain an image or any sort of drawable surface
type Image struct {
	componentBase

	content      *image.RGBA
	dirtyContent bool
}

func (i *Image) SetImage(img *image.RGBA) {
	i.content = img
	i.dirtyContent = true
}

func (i *Image) ShouldDraw() bool {
	return i.dirtyContent || i.componentBase.isDirty()
}

func (i *Image) Draw() {
	//TODO
}
