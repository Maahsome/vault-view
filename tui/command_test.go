package tui

import (
	"testing"

	"github.com/maahsome/tview"
)

// TestNewCommand - testing tui.newCommand
func TestNewCommand(t *testing.T) {

	tui := &Tui{
		app:   tview.NewApplication(),
		state: newState(),
	}

	tv := newCommand(tui)

	tv.Prompt = ":"
	tv.Cmd = "quit"

}
