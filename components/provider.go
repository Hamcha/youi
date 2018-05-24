package components

import (
	"strconv"
)

// Attribute is a single componment attribute
type Attribute string

// AttributeList is a list of attributes for a component, divided by key
type AttributeList map[string]Attribute

// ComponentProvider is a function that takes a list of attributes and creates a component with the attributes applied
type ComponentProvider func(AttributeList) (Component, error)

//
// Functions to easily convert attributes to their target types
//

// String returns the underlying unparsed string value of the attribute
func (a Attribute) String() string {
	return string(a)
}

// Float32 tries to parse an attribute as a float32 number
func (a Attribute) Float32() (float32, error) {
	ret, err := strconv.ParseFloat(string(a), 32)
	return float32(ret), err
}

// Int tries to parse an attribute as an integer number
func (a Attribute) Int() (int, error) {
	ret, err := strconv.Atoi(string(a))
	return ret, err
}

// Get returns either the requested attribute or a default value
func (a AttributeList) Get(name string, def string) Attribute {
	attr, ok := a[name]
	if !ok {
		return Attribute(def)
	}
	return attr
}
