package youi

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/kataras/go-errors"

	"github.com/hamcha/youi/components"
)

type yumlElement struct {
	Name       xml.Name
	Attributes yumlAttributes
	Children   []yumlChild
	Content    []byte
}

type yumlChild struct {
	Element  *yumlElement
	Settings []xml.Attr
}

type yumlAttributes []xml.Attr

// YUML errors
var (
	ErrIncompleteYuml      = errors.New("Incomplete YUML tree")
	ErrCouldNotMakeElement = errors.New("Could not make element \"%s\"")
)

func parseYUML(reader io.Reader) (*yumlElement, error) {
	var scope []*yumlElement
	var current *yumlElement

	decoder := xml.NewDecoder(reader)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch v := token.(type) {
		case xml.StartElement:
			//TODO Filter attributes and settings apart
			if current != nil {
				scope = append(scope, current)
			}
			current = &yumlElement{
				Name:       v.Name,
				Attributes: yumlAttributes(v.Attr),
			}
			if len(scope) > 0 {
				parent := scope[len(scope)-1]
				parent.Children = append(parent.Children, yumlChild{
					Element: current,
				})
			}
		case xml.EndElement:
			if len(scope) == 0 {
				return current, nil
			}
			current, scope = scope[len(scope)-1], scope[:len(scope)-1]
		case xml.CharData:
			if current != nil {
				current.Content = v
			}
		case xml.Comment:
			// Ignore comments, for now
		case xml.ProcInst:
			// Ignore procedures, for now
		case xml.Directive:
			// Ignore directives, for now

		}
	}
	return nil, ErrIncompleteYuml
}

func makeYUMLcomponentTree(element *yumlElement) (components.Component, error) {
	elem, err := makeComponent(element.Name.Space, element.Name.Local, element.Attributes.AttributeList())
	if err != nil {
		return nil, err
	}

	// Check for children
	for _, child := range element.Children {
		childelem, err := makeYUMLcomponentTree(child.Element)
		if err != nil {
			return nil, ErrCouldNotMakeElement.Format(element.Name.Local).AppendErr(err)
		}
		elem.AppendChild(childelem)
	}

	return elem, nil
}

func (y yumlElement) String() string {
	args := []string{}
	for _, arg := range y.Attributes {
		args = append(args, fmt.Sprintf("%s=%s", arg.Name.Local, arg.Value))
	}
	out := fmt.Sprintf("%s:%s (%s)\n", y.Name.Space, y.Name.Local, strings.Join(args, ", "))
	for i, child := range y.Children {
		symbol := '˫'
		if i == len(y.Children)-1 {
			symbol = '˪'
		}
		childindent := strings.Replace(child.Element.String(), "\n", "\n  ", -1)
		out += fmt.Sprintf("  %c %s\n", symbol, childindent)
	}
	return strings.TrimSpace(out)
}

func (y yumlAttributes) AttributeList() components.AttributeList {
	out := make(components.AttributeList)
	for _, attr := range y {
		out[attr.Name.Local] = components.Attribute(attr.Value)
	}
	return out
}