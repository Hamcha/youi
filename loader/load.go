package loader

import (
	"fmt"
	"image"
	"io"
	"os"

	resources "gopkg.in/cookieo9/resources-go.v2"
)

// This is just a wrapper for gopkg.in/cookieo9/resources-go.v2, in case someone else
// wants to use his own bundles

// BundleSequence is the bundle sequence that's used for opening/loading data
// If you need to use your own bundles, just override it
var BundleSequence resources.BundleSequence

func init() {
	BundleSequence = resources.DefaultBundle
}

// Open wraps resources.Open so it uses our custom bundle instead of the default one
func Open(path string) (io.ReadCloser, error) {
	fmt.Printf("Looking for \"%s\"\n", path)
	err := resources.CheckPath(path)
	if err != nil {
		return nil, err
	}
	//return BundleSequence.Open(path)

	// BundleSequence.Open is broken, use a slightly different version
	for _, bundle := range BundleSequence {
		if bundle == nil {
			continue
		}
		reader, err := bundle.Open(path)
		if err == nil {
			return reader, nil
		} else if err != resources.ErrNotFound && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return nil, resources.ErrNotFound
}

// Bytes opens a file from the bundle sequence and tries to read all of its contents immediately
func Bytes(path string) ([]byte, error) {
	res, err := Load(path)
	if err != nil {
		return nil, err
	}

	return res.Bytes()
}

// String calls Bytes(..) but returns the result as a string
func String(path string) (string, error) {
	res, err := Bytes(path)
	return string(res), err
}

// Image loads an image and returns it as a image.RGBA
func Image(path string) (*image.RGBA, error) {
	read, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer read.Close()

	return readImage(read)
}
