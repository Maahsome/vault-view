package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
	"github.com/nathan-fiscaletti/consolesize-go"
)

var inputWidth = 70

func (t *Tui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'h':
		t.prevPanel()
	case 'l':
		t.nextPanel()
	case 'q':
		t.Stop()
	case '/':
		t.searchPrompt()
	case ':':
		t.commandPrompt()
	case '?':
		t.helpPage()
	}

	switch event.Key() {
	case tcell.KeyTab:
		t.nextPanel()
	case tcell.KeyBacktab:
		t.prevPanel()
	}
}

func (t *Tui) helpPage() {

	w, h := consolesize.GetConsoleSize()

	viewName := "help"
	helpView := tview.NewTextView().SetDynamicColors(true)

	helpText := `
      [white]GENERAL                                                 NAVIGATION
      [red]<s> [brightwhite]      Toggle Showing Vault Values                   [red]<Left Arrow/ESC> [brightwhite]    Back (Data to Folder Frame as well)
      [red]<c> [brightwhite]      Copy Bash Command of Selected Item            [red]<Right Arrow/ENTER> [brightwhite] Open Selected Item
      [red]<C> [brightwhite]      Copy Bash Command for Marked Items            [red]<:quit> [brightwhite]             Quit
      [red]<v> [brightwhite]      Copy the selected Value                       [red]<TAB> [brightwhite]               Switch Frames
      [red]</text> [brightwhite]  Filter the list on "text"

      [red]<c> [brightwhite]      In Data pane, copy FIELD name for selection\n
      [red]<C> [brightwhite]      In Data pane, copy Bash Command for selection

      [red]<:folder>              Filter to list only Folder types
      [red]<:data>                Filter to list only Data types
      [red]<:all>                 Remove type filter


      [red]<q/ESC> [brightwhite]  Exit Help
`
	helpView.SetText(helpText)
	helpView.SetBorder(true)
	helpView.SetBorderColor(tcell.ColorDeepSkyBlue)

	closeHelpView := func() {
		t.closeAndSwitchPanel(viewName, t.state.panels.panel[t.state.panels.currentPanel].name())
	}

	helpView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyLeft {
			closeHelpView()
		}
		if event.Rune() == 'q' {
			closeHelpView()
		}
		return event
	})

	t.pages.AddAndSwitchToPage(viewName, t.modal(helpView, w, h), true).ShowPage("main")

}

func (t *Tui) searchPrompt() {
	t.state.grid.ResizeItem(t.state.command, 3, 1)
	t.state.command.focus(t)
	t.state.command.SetText("/")
	t.state.command.Prompt = "/"
	t.state.command.Cmd = ""
}

func (t *Tui) commandPrompt() {
	t.state.grid.ResizeItem(t.state.command, 3, 1)
	t.state.command.focus(t)
	t.state.command.SetText(":")
	t.state.command.Prompt = ":"
	t.state.command.Cmd = ""
}

func (t *Tui) nextPanel() {
	idx := (t.state.panels.currentPanel + 1) % len(t.state.panels.panel)
	t.switchPanel(t.state.panels.panel[idx].name())
}

func (t *Tui) prevPanel() {
	t.state.panels.currentPanel--

	if t.state.panels.currentPanel < 0 {
		t.state.panels.currentPanel = len(t.state.panels.panel) - 1
	}

	idx := (t.state.panels.currentPanel) % len(t.state.panels.panel)
	t.switchPanel(t.state.panels.panel[idx].name())
}
