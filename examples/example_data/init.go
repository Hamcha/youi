package example_data

import (
	"github.com/hamcha/youi/loader"
	resources "gopkg.in/cookieo9/resources-go.v2"
)

func init() {
	pkg, err := resources.OpenPackage("github.com/hamcha/youi/examples/example_data")
	if err != nil {
		panic(err)
	}
	loader.BundleSequence = append(loader.BundleSequence, pkg)
}
