package tui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"

	// "github.com/maahsome/vault-view/clipboard"
	"github.com/atotto/clipboard"
	"github.com/maahsome/vault-view/common"
	"github.com/maahsome/vault-view/resource"
	"github.com/maahsome/vault-view/vault"
	"github.com/sirupsen/logrus"
)

type folder struct {
	Type     string
	Path     string
	Parent   string
	FullPath string
	Version  int
}

type folders struct {
	*tview.Table
	filterWord string
	showTypes  int
	shownPath  string
	parent     string
	lang       *resource.Lang
}

func newFolders(t *Tui) *folders {
	folders := &folders{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		lang:  t.lang,
	}

	folders.SetTitle(fmt.Sprintf(" [[ %s ]] ", t.lang.GetText("ui", vaultFolder))).SetTitleAlign(tview.AlignLeft)
	folders.SetBorder(true)
	folders.SetBorderColor(tcell.ColorDeepSkyBlue)
	folders.setEntries(t, enterFolder)
	folders.setKeybinding(t)
	t.state.resources.markedPathFolders = make(map[string]bool, 0)
	t.state.resources.markedPathDatas = make(map[string]bool, 0)
	return folders
}

func (i *folders) name() string {
	return "folders"
}

func (i *folders) setTitle() {

	itemsShown := ""
	if i.showTypes != allItems {
		switch i.showTypes {
		case dataItems:
			itemsShown = fmt.Sprintf("[green](%s)[white]", i.lang.GetText("ui", vaultData))
		case folderItems:
			itemsShown = fmt.Sprintf("[green](%s)[white]", i.lang.GetText("ui", vaultFolder))
		}
	}
	if len(i.filterWord) > 0 {
		i.SetTitle(fmt.Sprintf(" [[ %s %s ]] - /%s/ ", i.lang.GetText("ui", vaultFolder), itemsShown, i.filterWord)).SetTitleAlign(tview.AlignLeft)
	} else {
		i.SetTitle(fmt.Sprintf(" [[ %s %s]] ", i.lang.GetText("ui", vaultFolder), itemsShown)).SetTitleAlign(tview.AlignLeft)
	}
}

func (i *folders) setKeybinding(t *Tui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		selectedFolder := t.selectedFolder()
		switch event.Key() {
		case tcell.KeyEnter:
			if selectedFolder != nil {
				if selectedFolder.Type == vaultFolder {
					i.filterWord = ""
					i.setTitle()
					i.setEntries(t, enterFolder)
				} else {
					if len(t.state.resources.datas) > 0 {
						t.dataPanel().setEntries(t, enterFolder)
						t.switchPanel("datas")
					}
				}
			}
		// KeyDown/KeyUp happen BEFORE the actual selection changes
		// So we temporarily increase by one, then decrease and allow the
		// default behavior to increase it after we display the path datas.
		case tcell.KeyCtrlSpace:
			if selectedFolder != nil {
				i.toggleSelected(t, selectedFolder)
			}
		case tcell.KeyDown:
			row, _ := t.folderPanel().GetSelection()

			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("down folders/row: %d/%d", len(t.state.resources.folders), row))

			if row < len(t.state.resources.folders) {
				tempRow := row + 1
				i.Select(tempRow, 0)
				t.dataPanel().setEntries(t, enterFolder)
				i.Select(row, 0)
			}
		case tcell.KeyUp:
			row, _ := t.folderPanel().GetSelection()
			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("up folders/row: %d/%d", len(t.state.resources.folders), row))
			if row > 0 {
				tempRow := row - 1
				i.Select(tempRow, 0)
				t.dataPanel().setEntries(t, enterFolder)
				i.Select(row, 0)
			}
		case tcell.KeyRight:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "keystrokes",
			}).Debug("KeyRight")

			if selectedFolder != nil {
				row, _ := t.folderPanel().GetSelection()
				common.Logger.WithFields(logrus.Fields{
					"unit":        "folders",
					"function":    "keystrokes",
					"row":         row,
					"parent":      selectedFolder.Parent,
					"folder_type": selectedFolder.Type,
				}).Info("RIGHT: Remember Where I was")
				t.state.resources.rowTracker[selectedFolder.Parent] = row
				if selectedFolder.Type == vaultFolder {
					i.filterWord = ""
					i.setTitle()
					i.setEntries(t, enterFolder)
				} else {
					t.dataPanel().setEntries(t, enterFolder)
				}
			}
		case tcell.KeyLeft, tcell.KeyEscape:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "keystrokes",
			}).Info("KeyLeft")
			row, _ := t.folderPanel().GetSelection()
			if selectedFolder != nil {
				common.Logger.WithFields(logrus.Fields{
					"unit":     "folders",
					"function": "keystrokes",
					"row":      row,
					"parent":   selectedFolder.Parent,
				}).Info("LEFT: Remember Where I was")
				t.state.resources.rowTracker[selectedFolder.Parent] = row
			}
			i.filterWord = ""
			i.setTitle()
			i.setEntries(t, enterParent)
		}

		var showValuesRune rune
		var runeerr error
		showValuesRune, runeerr = i.lang.GetRune("kbd", "s")
		if runeerr != nil {
			common.Logger.WithError(runeerr).WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "keystrokes",
				"rune":     "s",
			}).Error("Rune undefined")
		}
		switch event.Rune() {
		case showValuesRune:
			obscured = !obscured
			selectedFolder := t.selectedFolder()
			if selectedFolder.Type == vaultData {
				t.dataPanel().setEntries(t, enterFolder)
			}
		case ' ':
			i.toggleSelected(t, selectedFolder)
		case 'c':
			if selectedFolder.Type == vaultData {
				i.selectedCmdToClipboard(t, selectedFolder)
			} else {
				origText := t.state.info.Status.GetText(false)
				t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Only Valid for Data Type")))
				go t.ClearIndicator(1)
			}
		case 'C':
			i.markedCmdToClipboard(t)
		}

		return event
	})
}

func (i *folders) buildCmdForPath(t *Tui, path string) string {
	vCommand := ""
	var data vault.DataRecord
	var derr error

	if t.vaultCache.CacheDataExist(path, expireMinutes) {
		data = t.vaultCache.CacheDatas[path].Data
	} else {
		data, derr = t.vault.GetData(path)
		if derr != nil {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "buildcmd",
			}).WithError(derr).Error("failed to get data")
		}
	}

	dataFields := make([]string, 0, len(data.Data.Data))

	var folderData map[string]interface{}
	folderData = data.Data.Data
	for k := range folderData {
		dataFields = append(dataFields, k)
		common.Logger.WithFields(logrus.Fields{
			"unit":     "datas",
			"function": "data",
		}).Debug(fmt.Sprintf("Data Field: %s", k))
	}
	sort.Strings(dataFields)

	for _, sortedField := range dataFields {
		envField := strings.ToUpper(sortedField)
		vCommand += fmt.Sprintf("%s=$(vault kv get --format table --field %s secret%s); echo ${%s}\n", envField, sortedField, path, envField)
	}

	return vCommand
}

func (i *folders) selectedCmdToClipboard(t *Tui, s *folder) {

	vCommand := i.buildCmdForPath(t, s.FullPath)
	// clipboard.Set(vCommand)
	clipboard.WriteAll(vCommand)
	origText := t.state.info.Status.GetText(false)
	t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Copied to clipboard")))
	go t.ClearIndicator(1)
}

func (i *folders) markedCmdToClipboard(t *Tui) {
	vCommand := ""
	for k, v := range t.state.resources.markedPathDatas {
		if v {
			vCommand += fmt.Sprintf("%s\n", i.buildCmdForPath(t, k))
		}
	}

	// clipboard.Set(vCommand)
	clipboard.WriteAll(vCommand)
	origText := t.state.info.Status.GetText(false)
	t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Copied to clipboard")))
	go t.ClearIndicator(1)
}

func (i *folders) toggleSelected(t *Tui, selectedFolder *folder) {
	common.Logger.WithFields(logrus.Fields{
		"unit":             "folders",
		"function":         "marking",
		"tuipath":          selectedFolder.FullPath,
		"markedPathFolder": t.state.resources.markedPathFolders[selectedFolder.FullPath],
		"markedPathData":   t.state.resources.markedPathDatas[selectedFolder.FullPath],
	}).Debug(fmt.Sprintf("MarkThis: %s", selectedFolder.FullPath))

	if selectedFolder.Type == vaultFolder {
		t.state.resources.markedPathFolders[selectedFolder.FullPath] = !t.state.resources.markedPathFolders[selectedFolder.FullPath]
	} else {
		t.state.resources.markedPathDatas[selectedFolder.FullPath] = !t.state.resources.markedPathDatas[selectedFolder.FullPath]
	}
	row, _ := t.folderPanel().GetSelection()
	rowColor := tcell.ColorLightBlue
	if selectedFolder.Type == vaultFolder {
		if val, ok := t.state.resources.markedPathFolders[selectedFolder.FullPath]; ok {
			if val {
				rowColor = tcell.ColorMediumSeaGreen
			}
		}
	} else {
		if val, ok := t.state.resources.markedPathDatas[selectedFolder.FullPath]; ok {
			if val {
				rowColor = tcell.ColorMediumSeaGreen
			}
		}
	}
	for col := 0; col <= t.folderPanel().GetColumnCount(); col++ {
		t.folderPanel().GetCell(row, col).SetTextColor(rowColor)
	}
}

func (i *folders) buildPanelData(t *Tui, operation int) {

	fetchPath := "/"
	selectedFolder := t.selectedFolder()
	common.Logger.WithFields(logrus.Fields{
		"unit":     "",
		"function": "data",
	}).Trace(fmt.Sprintf("Selected Folder: %#v", selectedFolder))

	// Determine the folder we will build data for
	switch operation {
	case enterFolder:
		if selectedFolder != nil {
			fetchPath = selectedFolder.FullPath
		}
	case enterParent:
		common.Logger.WithFields(logrus.Fields{
			"unit":     "folders",
			"function": "data",
		}).Debug(fmt.Sprintf("NAV[folders]: Selected Parent from Panel: %s", i.getParent()))
		fetchPath = fmt.Sprintf("%s/", filepath.Dir(strings.TrimSuffix(i.getParent(), "/")))
		common.Logger.WithFields(logrus.Fields{
			"unit":     "folders",
			"function": "data",
		}).Debug(fmt.Sprintf("NAV[folders]: Initial FetchPath: %s", fetchPath))
		if fetchPath == "//" || fetchPath == "./" {
			fetchPath = "/"
		}
	case applyFilter:
		fetchPath = i.getParent()
	}

	i.setShownPath(fetchPath)
	common.Logger.WithFields(logrus.Fields{
		"unit":     "folders",
		"function": "data",
	}).Info(fmt.Sprintf("SET ShownPath: %s", fetchPath))

	var folders map[string]vault.Paths
	var serr error

	if t.vaultCache.CachePathExists(fetchPath) {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "folders",
			"function": "cache",
		}).Info("Loading PATHS from Cache... wooo!")
		folders = t.vaultCache.CachePaths[fetchPath].Paths
	} else {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "folders",
			"function": "data",
		}).Debug(fmt.Sprintf("Fetching Path Datas: %s", fetchPath))
		folders, serr = t.vault.GetPaths(fetchPath)
		if serr != nil {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "data",
			}).Error("failed to get paths")
		}
		if ok, err := t.vaultCache.UpdateCachePath(fetchPath, folders); !ok {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "folders",
				"function": "cache",
			}).WithError(err).Error("Bad UpdateCachePath")
		}
	}

	t.state.resources.folders = make(map[string]*folder, 0)
	t.state.resources.folderRows = make(map[int]string, 0)

	folderPaths := make([]string, 0, len(folders))
	for k := range folders {
		folderPaths = append(folderPaths, k)
	}
	sort.Strings(folderPaths)

	if len(folders) > 0 {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "folders",
			"function": "data",
		}).Debug(fmt.Sprintf("NAV[folders] SET Parent: %s", folders[folderPaths[0]].Parent))
		i.setParent(folders[folderPaths[0]].Parent)
		if t.state.location != nil {
			t.state.location.update(fmt.Sprintf("\n [white]%s", folders[folderPaths[0]].Parent))
		}
	}

	rowCount := 0
	for _, sortedPath := range folderPaths {
		folderInfo := folders[sortedPath]
		common.Logger.WithFields(logrus.Fields{
			"unit":          "folders",
			"function":      "data",
			"folder_type":   folderInfo.Type,
			"folder_path":   folderInfo.Path,
			"folder_parent": folderInfo.Parent,
			"filter_word":   i.filterWord,
			"show_types":    i.showTypes,
		}).Debug(fmt.Sprintf("SortedKEY: [%s]", sortedPath))
		if strings.Index(folderInfo.Path, i.filterWord) == -1 {
			continue
		}
		if i.showTypes == dataItems && folderInfo.Type == vaultFolder {
			continue
		}
		if i.showTypes == folderItems && folderInfo.Type == vaultData {
			continue
		}

		var folderData vault.DataRecord
		// CACHE: Load the datas as we are building the KEY list.
		if folderInfo.Type == vaultData {
			// Check Cache
			if !t.vaultCache.CachePathExists(folderInfo.FullPath) {
				vaultPaths := make(map[string]vault.Paths)
				vaultPaths[folderInfo.FullPath] = vault.Paths{
					Type:    folderInfo.Type,
					Path:    folderInfo.Path,
					Parent:  folderInfo.Parent,
					Version: folderData.Data.Metadata.Version,
				}
				if len(vaultPaths) > 0 {
					// TODO: PERFORMANCE, this used to be called as a go routine
					// It makes the interface load MUCH faster, though the VERSION
					// for the "Data" types lags in the interface, need some way
					// to go back and populate them once the cache is complete.
					// There is a framework here for time based reloading that
					// is probably the deal...
					// go t.vaultCache.PreloadPaths(vaultPaths)
					t.vaultCache.PreloadPaths(vaultPaths)
				}
			}
			// Attempt to find VERSION in the Data cache, this will lag in the display
			if t.vaultCache.CacheDataExist(folderInfo.FullPath, expireMinutes) {
				folderData = t.vaultCache.GetCacheData(folderInfo.FullPath)
			}
		} else {
			// Folder, pass this along to pre-load KEYs in the folders of this path
			go t.vaultCache.PreloadFolderPaths(folderInfo.FullPath)
		}

		t.state.resources.folders[folderInfo.FullPath] = &folder{
			Type:     folderInfo.Type,
			Path:     folderInfo.Path,
			Parent:   folderInfo.Parent,
			FullPath: folderInfo.FullPath,
			Version:  folderData.Data.Metadata.Version,
		}
		t.state.resources.folderRows[rowCount] = folderInfo.FullPath

		// TODO: I was originally just changing them all to FALSE upon load, this may still be the desired state
		// I supposed I could loop through the array as my "work" book and generate scripts from ALL marked items
		// rather than just marked items on the page.  Could be interesting.
		// if t.state.resources.markedFolders != nil {
		// 	t.state.resources.markedFolders[fmt.Sprintf("%s%s", folderInfo.Parent, folderInfo.Path)] = false
		// }
		rowCount++
	}
}

func (i *folders) setEntries(t *Tui, operation int) {
	i.buildPanelData(t, operation)
	table := i.Clear()

	headers := []string{
		i.lang.GetText("ui", "TYPE"),
		i.lang.GetText("ui", "PATH"),
		i.lang.GetText("ui", "VERSION"),
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorLightYellow,
			BackgroundColor: tcell.ColorDefault,
		})
	}

	folderPaths := make([]string, 0, len(t.state.resources.folders))
	for k := range t.state.resources.folders {
		folderPaths = append(folderPaths, k)
	}
	sort.Strings(folderPaths)

	c := 0
	for _, sortedPath := range folderPaths {
		folder := t.state.resources.folders[sortedPath]

		rowColor := tcell.ColorLightBlue
		if folder.Type == vaultFolder {
			if val, ok := t.state.resources.markedPathFolders[sortedPath]; ok {
				if val {
					rowColor = tcell.ColorMediumSeaGreen
				}
			}
		} else {
			if val, ok := t.state.resources.markedPathDatas[sortedPath]; ok {
				if val {
					rowColor = tcell.ColorMediumSeaGreen
				}
			}
		}
		table.SetCell(c+1, 0, tview.NewTableCell(i.lang.GetText("ui", folder.Type)).
			SetTextColor(rowColor).
			SetMaxWidth(10).
			SetExpansion(0))

		table.SetCell(c+1, 1, tview.NewTableCell(folder.Path).
			SetTextColor(rowColor).
			SetMaxWidth(1).
			SetExpansion(1))

		version := ""
		if folder.Version > 0 {
			version = fmt.Sprintf("%d", folder.Version)
		}
		table.SetCell(c+1, 2, tview.NewTableCell(version).
			SetTextColor(rowColor).
			SetMaxWidth(1).
			SetExpansion(1))
		c++
	}

	lastRow := 0
	if len(folderPaths) > 0 {
		lastRow = t.state.resources.rowTracker[t.state.resources.folders[folderPaths[0]].Parent]
	}
	if lastRow <= c {
		table.Select(lastRow, 0)
	} else {
		table.Select(0, 0)
	}
	i.ScrollToBeginning()
	if t.dataPanel() != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "folders",
			"function": "tuibuild",
		}).Debug("Setting Entries for Datas panel")
		t.dataPanel().setEntries(t, enterFolder)
	}

}

func (i *folders) updateEntries(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		i.setEntries(t, enterFolder)
	})
}

func (i *folders) focus(t *Tui) {
	i.SetSelectable(true, false)
	t.app.SetFocus(i)
}

func (i *folders) unfocus() {
	i.SetSelectable(false, false)
}

func (i *folders) setFilterWord(word string) {
	i.filterWord = word
}

func (i *folders) setFilterType(which int) {
	i.showTypes = which
}

func (i *folders) setShownPath(path string) {
	i.shownPath = path
}

func (i *folders) getShownPath() string {
	return i.shownPath
}

func (i *folders) setParent(path string) {
	i.parent = path
}

func (i *folders) getParent() string {
	return i.parent
}
