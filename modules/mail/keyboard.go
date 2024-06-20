package mail

import (
	"github.com/gdamore/tcell/v2"
)

func (widget *Widget) initializeKeyboardControls() {
	widget.SetKeyboardChar("n", widget.Next, "Next item in list")
	widget.SetKeyboardChar("p", widget.Prev, "Next item in list")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.Select, "Select")
	widget.SetKeyboardChar("q", widget.HandleQ, "Go Back")
}
