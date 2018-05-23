package components

const Namespace = "https://yuml.ovo.ovh/schema/components/1.0"

var AllComponents = map[string]Component{
	"Canvas": &Canvas{},
	"Image":  &Image{},
	"Label":  &Label{},
}
