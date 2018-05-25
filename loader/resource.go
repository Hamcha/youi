package loader

import (
	"image"
	"io"
	"io/ioutil"
)

// Resource is a wrapper for a ReadCloser with some extra sugar
type Resource struct {
	io.ReadCloser
}

// Load tries to open a file and returns a Resource when successful
func Load(path string) (Resource, error) {
	res, err := Open(path)
	return Resource{res}, err
}

// Bytes tries to read the file data as bytes
// Warning: once used, the reader underneath will be closed!
func (r Resource) Bytes() ([]byte, error) {
	defer r.Close()
	return ioutil.ReadAll(r)
}

// String tries to read the file data as a string
// Warning: once used, the reader underneath will be closed!
func (r Resource) String() (string, error) {
	byt, err := r.Bytes()
	return string(byt), err
}

// Image tries to decode the file as an image
// Warning: once used, the reader underneath will be closed!
// Image loads an image and returns it as a image.RGBA
func (r Resource) Image() (*image.RGBA, error) {
	defer r.Close()
	return readImage(r)
}
