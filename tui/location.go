package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
)

type location struct {
	*tview.TextView
}

func newLocation() *location {
	return &location{
		TextView: tview.NewTextView().SetDynamicColors(true).SetTextColor(tcell.ColorDeepSkyBlue),
	}
}

func (l *location) update(path string) {
	l.SetText(path)
}
