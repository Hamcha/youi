package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/hamcha/youi/font"
)

func main() {
	in := flag.String("in", "", "Font TTF file")
	name := flag.String("name", "", "Font name")
	out := flag.String("out", "fonts", "Directory where to save the generated textures")
	size := flag.Int("size", font.DefaultTextureFontSize, "Font size")
	flag.Parse()

	if *in == "" {
		fmt.Fprintln(os.Stderr, "-in flag required but missing")
		flag.Usage()
		os.Exit(1)
	}

	if *name == "" {
		fmt.Fprintln(os.Stderr, "-name flag required but missing")
		flag.Usage()
		os.Exit(1)
	}

	// Create output directory if missing
	if _, err := os.Stat(*out); os.IsNotExist(err) {
		checkErr(os.MkdirAll(*out, 0755), "Could not create output directory")
	}

	// Load file from disk
	ttffile, err := ioutil.ReadFile(*in)
	checkErr(err, "Could not read input file")

	// Parse truetype font
	ttf, err := truetype.Parse(ttffile)
	checkErr(err, "Could not parse file as TTF file")

	// Generate SDF texture and atlas
	font, err := font.MakeFont(ttf, *size)
	checkErr(err, "Could not generate texture or atlas")

	// Save generated files to disk
	checkErr(font.Export(*out, *name), "Could not save files to disk")
}

func checkErr(err error, msg string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: "+msg+":\n", args...)
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
