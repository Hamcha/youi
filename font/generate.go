package font

import (
	"errors"
	"image"

	"github.com/adinfinit/texpack/maxrect"
	"github.com/adinfinit/texpack/pack"
	"github.com/adinfinit/texpack/sdf"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// MaxTextureSize is the maximum width/height each font texture can have
const MaxTextureSize = 2048

// SDFRadius is the search radius for each pixel during the distance fields generation
const SDFRadius = 4

// DefaultTextureFontSize is the default size used for SDF texture generation
const DefaultTextureFontSize = 32

var (
	// ErrFontDoesntFit means the font glyphs don't fit within the maximum texture limit
	ErrFontDoesntFit = errors.New("font texture doesn't fit within texture size limits")
)

// Atlas is a dictionary that maps each font glyph to its position on a texture atlas
type Atlas map[rune]image.Rectangle

// MakeFont creates an SDF font atlas from a ttf font
func MakeFont(fnt *truetype.Font, fontSize int) (*Font, error) {
	// Create new typeface with specified size
	face := truetype.NewFace(fnt, &truetype.Options{
		Size:    float64(fontSize),
		Hinting: font.HintingFull,
	})

	// Create textpack font struct
	font := &pack.Font{
		Font:      fnt,
		Face:      face,
		MaxBounds: fnt.Bounds(fixed.I(fontSize)),
		Glyphs:    make(map[rune]*pack.Glyph, 256),
		Padding:   SDFRadius / 2,
		Kern:      make(map[[2]rune]fixed.Int26_6, 256),
	}

	// Include extra characters
	//TODO Should include more stuff, maybe user-defined?
	font.IncludeExtendedAscii()

	// Get boxes for all characters
	boxes := font.Boxes()

	// Get best texture size for all boxes
	size, ok := pack.PlaceBoxes(image.Point{MaxTextureSize, MaxTextureSize}, maxrect.Automatic, boxes)
	if !ok {
		return nil, ErrFontDoesntFit
	}

	// Make output image for font map and draw all glyphs
	dst := image.NewRGBA(image.Rectangle{image.ZP, size})
	font.Draw(dst)

	// Use signed distance fields on the texture
	sdf.ApplyRGBA_Alpha(dst, SDFRadius)

	// Generate fontmap
	fontmap := make(Atlas)
	for r, box := range font.Glyphs {
		fontmap[r] = box.Place
	}

	return &Font{
		Texture:  dst,
		Atlas:    fontmap,
		FontSize: fontSize,
	}, nil
}
