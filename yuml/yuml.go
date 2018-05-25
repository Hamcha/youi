package yuml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/kataras/go-errors"
)

// Element contains all the data about a YUML element
type Element struct {
	Name       xml.Name
	Attributes Attributes
	Children   []Child
	Content    []byte
}

// Child contains a YUML element and its parent-related attributes
type Child struct {
	Element  *Element
	Settings []xml.Attr
}

// Attributes are YUML attributes (basically xml.Attr with some extra sugar)
type Attributes []xml.Attr

// YUML errors
var (
	ErrIncompleteYuml = errors.New("Incomplete YUML tree")
)

// ParseYUML tries to read YUML code and parse it as such, returning the root element
func ParseYUML(reader io.Reader) (*Element, error) {
	var scope []*Element
	var current *Element

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
			current = &Element{
				Name:       v.Name,
				Attributes: Attributes(v.Attr),
			}
			if len(scope) > 0 {
				parent := scope[len(scope)-1]
				parent.Children = append(parent.Children, Child{
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

func (y Element) String() string {
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
