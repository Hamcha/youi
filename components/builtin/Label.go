package builtin

import "github.com/hamcha/youi/components"

// Label is a drawable text label
type Label struct {
	components.ComponentDrawable
	components.ComponentText

	content      string
	dirtyContent bool
}

// SetText changes the text content of the label
func (l *Label) SetText(str string) {
	l.content = str
	l.dirtyContent = true
}

// Draw draws the label on screen
func (l *Label) Draw() {
	//TODO

	l.ComponentDrawable.Draw()
	l.ComponentText.Draw()
	l.ClearFlags()
}

func (l *Label) ClearFlags() {
	l.dirtyContent = false
}

// ShouldDraw returns whether the label needs to be re-drawn
func (l *Label) ShouldDraw() bool {
	return l.dirtyContent || l.ComponentText.ShouldDraw()
}

func (l *Label) String() string {
	//TODO Add attributes
	return "<Label />"
}

func makeLabel(list components.AttributeList) (components.Component, error) {
	//TODO Parse attributes
	return &Label{}, nil
}
