package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
)

type infoIndicator struct {
	*tview.TextView
}

func newInfoIndicator(t *Tui) *infoIndicator {

	i := &infoIndicator{
		TextView: tview.NewTextView(),
	}

	i.display(t.lang.GetText("ui", "Successfully Connected"))

	return i
}

func (i *infoIndicator) display(status string) {

	i.SetDynamicColors(true).SetTextColor(tcell.ColorDarkOrange)

	i.SetText(status)
}
