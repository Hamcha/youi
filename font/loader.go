package font

import (
	"encoding/gob"
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/hamcha/youi/loader"
)

// Font holds the necessary data for using a font in youi
type Font struct {
	Texture *image.RGBA
	Atlas   Atlas
}

// ErrFontNotFound means any valid font formats (ttf, sdf+atlas) could not be found
var ErrFontNotFound = errors.New("could not find font files")

type fontCache map[string]*Font

var fonts fontCache

func init() {
	ClearCache()
}

// LoadFont loads either the pregenerated SDF+Atlas or TTF font from disk or binary and returns a usable youi font
func LoadFont(fontName string) (*Font, error) {
	// Check in local cache
	if fnt, ok := fonts[fontName]; ok {
		return fnt, nil
	}

	// Check for a prebuilt SDF texture
	texstream, errPNG := loader.Load("fonts/" + fontName + ".png")
	atlstream, errAtlas := loader.Load("fonts/" + fontName + ".atlas")
	if errPNG == nil && errAtlas == nil {
		defer texstream.Close()
		defer atlstream.Close()

		// Load PNG texture
		teximg, err := texstream.Image()
		if err != nil {
			return nil, err
		}

		// Load Atlas
		var atlas Atlas
		err = gob.NewDecoder(atlstream).Decode(&atlas)
		if err != nil {
			return nil, err
		}

		// Save to cache
		fonts[fontName] = &Font{
			Texture: teximg,
			Atlas:   atlas,
		}
		return fonts[fontName], nil
	}

	// Check for TTF file
	ttffile, err := loader.Bytes("fonts/" + fontName + ".ttf")
	if err == nil {
		// Parse truetype font
		ttf, err := truetype.Parse(ttffile)
		if err != nil {
			return nil, err
		}

		// Generate SDF texture
		fonts[fontName], err = MakeFont(ttf, DefaultTextureFontSize)
		return fonts[fontName], err
	}

	return nil, ErrFontNotFound
}

// Export saves the font texture and atlas to disk
func (f *Font) Export(directory, name string) error {
	// Save font texture (to PNG)
	texfile, err := os.Create(filepath.Join(directory, name+".png"))
	if err != nil {
		return err
	}

	err = png.Encode(texfile, f.Texture)
	if err != nil {
		return err
	}

	texfile.Close()

	// Save font atlas (to gob)
	atlfile, err := os.Create(filepath.Join(directory, name+".atlas"))
	if err != nil {
		return err
	}

	err = gob.NewEncoder(atlfile).Encode(f.Atlas)
	if err != nil {
		return err
	}

	// All done!
	return nil
}

// DefaultFont loads (and generate, if necessary) the default font for youi (goregular)
func DefaultFont() *Font {
	// Check if it's bundled or cached
	fnt, err := LoadFont("default")
	if err == nil {
		return fnt
	}

	// Nothing? Guess we have to generate it
	ttf, _ := truetype.Parse(goregular.TTF)

	fonts["default"], err = MakeFont(ttf, DefaultTextureFontSize)
	if err != nil {
		// This is pretty bad
		panic(err)
	}

	return fonts["default"]
}

// ClearCache removes all entries from the font cache
func ClearCache() {
	fonts = make(fontCache)
}
