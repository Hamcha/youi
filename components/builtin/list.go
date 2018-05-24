package builtin

import "github.com/hamcha/youi/components"

const Namespace = "https://yuml.ovo.ovh/schema/components/1.0"

var AllComponents = map[string]components.ComponentProvider{
	"Page":   makePage,
	"Canvas": makeCanvas,
	"Image":  makeImage,
	"Label":  makeLabel,
}
