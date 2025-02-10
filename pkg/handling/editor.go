package handling

import "fyne.io/fyne/v2/widget"

// Editor wraps the text editor widget.
type Editor struct {
	Widget *widget.Entry
}

// NewEditor creates a new markdown editor.
func NewEditor() *Editor {
	return &Editor{Widget: widget.NewMultiLineEntry()}
}
