package tui

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
)

type command struct {
	Cmd string
}

type commands struct {
	*tview.TextView
	Prompt string
	Cmd    string
}

func newCommand(t *Tui) *commands {
	commands := &commands{
		TextView: tview.NewTextView(),
	}

	commands.SetWordWrap(true)
	commands.SetWrap(true)
	commands.SetDynamicColors(true)
	commands.SetBorder(true)
	commands.SetBorderColor(tcell.ColorDeepSkyBlue)
	commands.setKeybinding(t)
	return commands
}

func (i *commands) setKeybinding(t *Tui) {
	i.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyBackspace2, tcell.KeyBackspace, tcell.KeyDelete:
			if len(i.Cmd) > 0 {
				i.Cmd = i.Cmd[:len(i.Cmd)-1]
				i.SetText(fmt.Sprintf("%s %s", i.Prompt, i.Cmd))
			}

		case tcell.KeyEnter:
			t.state.grid.ResizeItem(i, 0, 0)
			currentPanel := t.state.panels.panel[t.state.panels.currentPanel]
			t.switchPanel(currentPanel.name())
			if i.Prompt == ":" {
				switch i.Cmd {
				case "quit", "q", "qu", "qui":
					t.Stop()
				case "data", "d", "da", "dat":
					currentPanel.setFilterType(dataItems)
					currentPanel.setTitle()
					currentPanel.setEntries(t, applyFilter)
				case "folders", "f", "fo", "fol", "fold", "folde", "folder":
					currentPanel.setFilterType(folderItems)
					currentPanel.setTitle()
					currentPanel.setEntries(t, applyFilter)
				case "all", "a", "al":
					currentPanel.setFilterType(allItems)
					currentPanel.setTitle()
					currentPanel.setEntries(t, applyFilter)
				}
			}
			if i.Prompt == "/" {
				// set our search filter
				currentPanel.setFilterWord(i.Cmd)
				currentPanel.setTitle()
				currentPanel.setEntries(t, applyFilter)
			}
		case tcell.KeyEscape:
			i.SetText("")
			i.Cmd = ""
			t.state.grid.ResizeItem(i, 0, 0)
			currentPanel := t.state.panels.panel[t.state.panels.currentPanel]
			t.switchPanel(currentPanel.name())
		}

		if isPrint(evt.Rune()) {
			i.Cmd += string(evt.Rune())
			i.SetText(fmt.Sprintf("%s %s", i.Prompt, i.Cmd))
		}

		return evt
	})

}

func isPrint(r rune) bool {
	if unicode.IsPrint(r) {
		return true
	}
	return false
}

func (i *commands) name() string {
	return "commands"
}

func (i *commands) focus(t *Tui) {
	t.app.SetFocus(i)
}
