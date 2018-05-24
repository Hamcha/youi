package youi

import (
	"fmt"

	"github.com/kataras/go-errors"

	"github.com/hamcha/youi/components"
	"github.com/hamcha/youi/components/builtin"
)

type componentList map[string]components.ComponentProvider

var namespaces map[string]componentList

// Component errors
var (
	ErrUnknownNamespace = errors.New("YUML namespace \"%s\" not found (have you forgot to register its components?)")
	ErrUnknownComponent = errors.New("YUML component \"%s\" not found in namespace \"%s\"")
)

// RegisterComponent registers a component under a namespace and name so it can be used in YUML code
func RegisterComponent(namespace, name string, provider components.ComponentProvider) {
	ns, ok := namespaces[namespace]
	if !ok {
		namespaces[namespace] = make(componentList)
		ns = namespaces[namespace]
	}
	ns[name] = provider
}

func makeComponent(namespace, name string, attributes components.AttributeList) (components.Component, error) {
	ns, ok := namespaces[namespace]
	if !ok {
		return nil, ErrUnknownNamespace.Format(namespace)
	}

	provider, ok := ns[name]
	if !ok {
		return nil, ErrUnknownComponent.Format(name, namespace)
	}

	return provider(attributes)
}

func initBuiltinComponents() {
	for name, elem := range builtin.AllComponents {
		RegisterComponent(builtin.Namespace, name, elem)
	}
}

type ComponentMap map[string][]string

// Components returns a map of all registered components
func Components() (out ComponentMap) {
	out = make(map[string][]string)
	for ns, components := range namespaces {
		out[ns] = make([]string, len(components))
		index := 0
		for name := range components {
			out[ns][index] = name
			index++
		}
	}
	return
}

func (c ComponentMap) String() (out string) {
	for ns, components := range c {
		out += fmt.Sprintf("%s\n", ns)
		for _, name := range components {
			out += fmt.Sprintf("  %s\n", name)
		}
	}
	return
}
