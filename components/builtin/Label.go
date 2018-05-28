package builtin

import "github.com/hamcha/youi/components"

// Label is a drawable text label
type Label struct {
	components.Base
	components.Text
}

// Draw draws the label on screen
func (l *Label) Draw() {
	//TODO

	l.Text.Draw()
	l.Text.ClearFlags()
}

// ShouldDraw returns whether the label needs to be re-drawn
func (l *Label) ShouldDraw() bool {
	return l.Text.ShouldDraw()
}

func (l *Label) String() string {
	//TODO Add attributes
	return "<Label />"
}

func makeLabel(list components.AttributeList) (components.Component, error) {
	//TODO Parse attributes
	return &Label{}, nil
}
