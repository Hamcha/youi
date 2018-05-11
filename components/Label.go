package components

// Label is a drawable text label
type Label struct {
	componentPosition
	componentText

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
}
