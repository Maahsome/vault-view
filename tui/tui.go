package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
	"github.com/maahsome/vault-view/common"
	"github.com/maahsome/vault-view/resource"
	"github.com/maahsome/vault-view/vault"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/sirupsen/logrus"
)

type panels struct {
	currentPanel int
	panel        []panel
}

// vault resources
type resources struct {
	folders           map[string]*folder
	folderRows        map[int]string
	markedPathFolders map[string]bool
	markedPathDatas   map[string]bool
	datas             map[string]*dataRecord
	dataRows          map[int]string
	markedDatas       map[string][]*dataMark
	rowTracker        map[string]int
}

type state struct {
	panels    panels
	location  *location
	command   *commands
	grid      *tview.Flex
	info      *info
	resources resources
	stopChans map[string]chan int
}

const (
	enterFolder = iota
	enterParent
	applyFilter
)

const (
	allItems = iota
	folderItems
	dataItems
)

const (
	vaultFolder = "Folder"
	vaultData   = "Data"
)

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

// Tui - Structure that runs the application
type Tui struct {
	app        *tview.Application
	pages      *tview.Pages
	state      *state
	vault      vault.Client
	vaultCache *vault.Cache
	lang       *resource.Lang
	semVer     string
}

// New create new Tui
func New(version string) *Tui {
	vaultCli := vault.NewVault()
	vaultCache := vault.NewCache(vaultCli)

	obscured = true
	return &Tui{
		app:        tview.NewApplication(),
		state:      newState(),
		vault:      vaultCli,
		vaultCache: vaultCache,
		lang:       resource.NewLanguage(),
		semVer:     version,
	}
}

func (t *Tui) folderPanel() *folders {
	for _, panel := range t.state.panels.panel {
		if panel.name() == "folders" {
			return panel.(*folders)
		}
	}
	return nil
}

func (t *Tui) dataPanel() *datas {
	for _, panel := range t.state.panels.panel {
		if panel.name() == "datas" {
			return panel.(*datas)
		}
	}
	return nil
}

func (t *Tui) initPanels() {
	folders := newFolders(t)
	command := newCommand(t)
	info := newInfo(t)
	location := newLocation()
	datas := newDatas(t)

	vaultInfo := newVaultInfo(t.vault)
	vaultVersion := fmt.Sprintf("%s", vaultInfo.Version)
	vaultEndpoint := fmt.Sprintf("%s", vaultInfo.Addr)
	address := tview.NewTextView().SetTextColor(tcell.ColorWhite).
		SetText(fmt.Sprintf(" %s (%s)", vaultEndpoint, vaultVersion))
	location.update("\n [white]/")
	go t.ClearIndicator(10)

	t.state.panels.panel = append(t.state.panels.panel, folders)
	t.state.panels.panel = append(t.state.panels.panel, datas)
	t.state.command = command
	t.state.location = location
	t.state.info = info
	t.state.resources.rowTracker = make(map[string]int)

	_, h := consolesize.GetConsoleSize()

	grid := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(info, 6, 1, false).
		AddItem(address, 2, 1, false).
		AddItem(command, 3, 1, false).
		AddItem(folders, 0, (h-4)/2, true).
		AddItem(datas, 0, (h-4)/2, false).
		AddItem(location, 3, 1, false)

	grid.ResizeItem(t.state.command, 0, 0)

	t.state.grid = grid

	t.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	t.app.SetRoot(t.pages, true)
	t.switchPanel("folders")
}

// Start start application
func (t *Tui) Start() error {
	t.initPanels()
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}

	return nil
}

// ClearIndicator - Clear the status message after s seconds
func (t *Tui) ClearIndicator(s time.Duration) {
	time.Sleep(s * time.Second)
	t.app.QueueUpdateDraw(func() {
		t.state.info.Status.display("")
	})
}

// Stop stop application
func (t *Tui) Stop() error {
	t.app.Stop()
	return nil
}

func (t *Tui) selectedFolder() *folder {
	if len(t.state.resources.folders) == 0 {
		return nil
	}

	if t.folderPanel() != nil {
		row, _ := t.folderPanel().GetSelection()
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selection",
			"row":      row,
			"datamap":  t.state.resources.folderRows[row-1],
		}).Debug("Row for Selection")
		if len(t.state.resources.folders) == 0 {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selection",
			}).Debug("Our folders resources is currenty EMPTY")
			return nil
		}
		if row-1 < 0 {
			// if we are at row 0, then return the folder for item 0
			return t.state.resources.folders[t.state.resources.folderRows[0]]
		}
		return t.state.resources.folders[t.state.resources.folderRows[row-1]]
	}
	return nil
}

func (t *Tui) selectedData() *dataRecord {
	if len(t.state.resources.datas) == 0 {
		return nil
	}

	if t.dataPanel() != nil {
		row, _ := t.dataPanel().GetSelection()
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selection",
			"row":      row,
			"datamap":  t.state.resources.dataRows[row-1],
		}).Debug("Row for Selection: ", fmt.Sprintf("%d", row))
		if len(t.state.resources.datas) == 0 {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selection",
			}).Debug("Our datas resources is currenty EMPTY")
			return nil
		}
		if row-1 < 0 {
			// if we are at row 0, then return the folder for item 0
			return t.state.resources.datas[t.state.resources.dataRows[0]]
		}
		return t.state.resources.datas[t.state.resources.dataRows[row-1]]
	}
	return nil
}

func (t *Tui) switchPanel(panelName string) {
	for i, panel := range t.state.panels.panel {
		if panel.name() == panelName {
			t.state.info.Menu.display(panelName)
			panel.focus(t)
			t.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}

func (t *Tui) closeAndSwitchPanel(removePanel, switchPanel string) {
	t.pages.RemovePage(removePanel).ShowPage("main")
	t.switchPanel(switchPanel)
}

func (t *Tui) modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (t *Tui) currentPanel() panel {
	return t.state.panels.panel[t.state.panels.currentPanel]
}
