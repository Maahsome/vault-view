package tui

type panel interface {
	name() string
	setTitle()
	buildPanelData(*Tui, int)
	setEntries(*Tui, int)
	updateEntries(*Tui)
	setKeybinding(*Tui)
	focus(*Tui)
	unfocus()
	setFilterWord(string)
	setFilterType(int)
	setShownPath(string)
	getShownPath() string
	setParent(string)
	getParent() string
}
