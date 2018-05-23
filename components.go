package youi

import (
	"errors"
	"reflect"

	"github.com/hamcha/youi/components"
)

type componentList map[string]reflect.Type

var namespaces map[string]componentList

// Component errors
var (
	ErrUnknownNamespace = errors.New("YUML namespace not found (have you forgot to register its components?)")
	ErrUnknownComponent = errors.New("YUML component not found in namespace")
)

// RegisterComponent registers a component under a namespace and name so it can be used in YUML code
func RegisterComponent(namespace, name string, component components.Component) {
	ns, ok := namespaces[namespace]
	if !ok {
		namespaces[namespace] = make(componentList)
		ns = namespaces[namespace]
	}
	ns[name] = reflect.TypeOf(component)
}

func makeComponent(namespace, name string) (interface{}, error) {
	ns, ok := namespaces[namespace]
	if !ok {
		return nil, ErrUnknownNamespace
	}

	typ, ok := ns[name]
	if !ok {
		return nil, ErrUnknownComponent
	}

	return reflect.New(typ).Interface(), nil
}

func initBuiltinComponents() {
	for name, elem := range components.AllComponents {
		RegisterComponent(components.Namespace, name, elem)
	}
}
