package tui

import (
	"encoding/json"
	"fmt"
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

type dataRecord struct {
	Field        string
	Value        string
	DisplayValue string
	Version      int
}

type dataMark struct {
	Field string
}

type datas struct {
	*tview.Table
	filterWord string
	showTypes  int
	shownPath  string
	parent     string
	lang       *resource.Lang
}

var obscured bool

func newDatas(t *Tui) *datas {
	datas := &datas{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		lang:  t.lang,
	}

	datas.SetTitle(fmt.Sprintf(" [[ %s ]] ", t.lang.GetText("ui", vaultData))).SetTitleAlign(tview.AlignLeft)
	datas.SetBorder(true)
	datas.SetBorderColor(tcell.ColorDeepSkyBlue)
	datas.setEntries(t, enterFolder)
	datas.setKeybinding(t)
	return datas
}

func (i *datas) name() string {
	return "datas"
}

func (i *datas) setTitle() {
	if len(i.filterWord) > 0 {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] - /%s/", i.lang.GetText("ui", vaultData), i.filterWord)).SetTitleAlign(tview.AlignLeft)
	} else {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] ", i.lang.GetText("ui", vaultData))).SetTitleAlign(tview.AlignLeft)
	}
}

func (i *datas) setKeybinding(t *Tui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)

		selectedFolder := t.selectedFolder()
		selectedData := t.selectedData()

		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyEscape:
			t.prevPanel()
		}
		switch event.Rune() {
		case 'c':
			// Copy KEY name to clipboard
			// clipboard .Set(selectedData.Field)
			clipboard.WriteAll(selectedData.Field)
			origText := t.state.info.Status.GetText(false)
			t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Copied to clipboard")))
			go t.ClearIndicator(1)
		case 'C':
			// Copy Vault Commmand to clipboard
			vCommand := fmt.Sprintf("vault kv get --format table --field %s secret%s", selectedData.Field, selectedFolder.FullPath)
			// clipboard.Set(vCommand)
			clipboard.WriteAll(vCommand)
			origText := t.state.info.Status.GetText(false)
			t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Copied to clipboard")))
			go t.ClearIndicator(1)
		case 'v':
			// Copy the VALUE to clipboard
			// clipboard.Set(selectedData.Value)
			clipboard.WriteAll(selectedData.Value)
			origText := t.state.info.Status.GetText(false)
			t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Copied to clipboard")))
			go t.ClearIndicator(1)
		case 's':
			obscured = !obscured
			t.dataPanel().setEntries(t, enterFolder)
		}

		return event
	})
}

func (i *datas) buildPanelData(t *Tui, operation int) {

	selectedFolder := t.selectedFolder()
	if selectedFolder != nil {
		if selectedFolder.Type == vaultData {
			common.Logger.WithFields(logrus.Fields{
				"unit":            "datas",
				"function":        "data",
				"selected_parent": selectedFolder.Parent,
				"selected_path":   selectedFolder.Path,
			}).Debug(fmt.Sprintf("SET ShownPath"))
			i.setShownPath(selectedFolder.FullPath)

			var data vault.DataRecord
			var derr error
			if t.vaultCache.CacheDataExist(selectedFolder.FullPath, expireMinutes) {
				// Add a check here to see how old our data is
				common.Logger.WithFields(logrus.Fields{
					"unit":     "datas",
					"function": "data",
				}).Debug("Loading FIELD/VALUES from Cache... wooo!")
				data = t.vaultCache.CacheDatas[selectedFolder.FullPath].Data
			} else {
				data, derr = t.vault.GetData(selectedFolder.FullPath)
				if derr != nil {
					common.Logger.WithFields(logrus.Fields{
						"unit":     "datas",
						"function": "data",
					}).WithError(derr).Error("failed to get data")
				}
				if ok, err := t.vaultCache.UpdateCacheData(selectedFolder.FullPath, data); !ok {
					common.Logger.WithFields(logrus.Fields{
						"unit":     "datas",
						"function": "data",
					}).WithError(err).Error("Bad UpdateCacheData")
				}
			}

			common.Logger.WithFields(logrus.Fields{
				"unit":     "datas",
				"function": "data",
			}).Trace(fmt.Sprintf("Cached DATA: %#v", data))
			t.state.resources.datas = make(map[string]*dataRecord, 0)
			t.state.resources.dataRows = make(map[int]string, 0)
			dataFields := make([]string, 0, len(data.Data.Data))

			var folderData map[string]interface{}

			folderData = data.Data.Data
			common.Logger.WithFields(logrus.Fields{
				"unit":     "datas",
				"function": "data",
			}).Trace(fmt.Sprintf("JSON/DATA data.Data: %#v", data.Data.Data))
			for k := range folderData {
				dataFields = append(dataFields, k)
				common.Logger.WithFields(logrus.Fields{
					"unit":     "datas",
					"function": "data",
				}).Debug(fmt.Sprintf("Data Field: %s", k))
			}

			common.Logger.WithFields(logrus.Fields{
				"unit":     "datas",
				"function": "data",
			}).Debug(fmt.Sprintf("JSON[version]: %d", data.Data.Metadata.Version))
			sort.Strings(dataFields)
			rowCount := 0

			for _, sortedField := range dataFields {
				dataValue := folderData[sortedField]
				common.Logger.WithFields(logrus.Fields{
					"unit":     "datas",
					"function": "data",
				}).Trace(fmt.Sprintf("JSON/DATA: Field: %s, Value: %#v", sortedField, dataValue))
				if strings.Index(sortedField, i.filterWord) == -1 {
					continue
				}
				displayValue := ""
				actualValue := ""
				if dv, ok := dataValue.(string); !ok {

					jv, err := json.Marshal(dataValue)
					if err != nil {
						displayValue = "JSON"
						actualValue = "JSON"
					} else {
						displayValue = string(jv[:])
						actualValue = displayValue
					}
				} else {
					displayValue = dv
					actualValue = dv
				}
				common.Logger.WithFields(logrus.Fields{
					"unit":     "datas",
					"function": "data",
				}).Trace("Adding: ", fmt.Sprintf("Field: %s, Value: %s", sortedField, dataValue))
				if obscured {
					displayValue = "**********"
				}

				t.state.resources.datas[sortedField] = &dataRecord{
					Field:        sortedField,
					Value:        actualValue,
					DisplayValue: displayValue,
					Version:      data.Data.Metadata.Version,
				}
				t.state.resources.dataRows[rowCount] = sortedField
				rowCount++
			}
		} else {
			t.state.resources.datas = make(map[string]*dataRecord, 0)
		}
	} else {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "datas",
			"function": "data",
		}).Warn("failed to get path data")
		t.state.resources.datas = make(map[string]*dataRecord, 0)
	}
}

func (i *datas) setEntries(t *Tui, operation int) {
	i.buildPanelData(t, operation)
	table := i.Clear()

	headers := []string{
		i.lang.GetText("ui", "FIELD"),
		i.lang.GetText("ui", "VALUE"),
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorLightYellow,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	dataFields := make([]string, 0, len(t.state.resources.datas))
	for k := range t.state.resources.datas {
		dataFields = append(dataFields, k)
	}
	sort.Strings(dataFields)

	c := 0
	for _, sortedField := range dataFields {
		dataRec := t.state.resources.datas[sortedField]
		table.SetCell(c+1, 0, tview.NewTableCell(dataRec.Field).
			SetTextColor(tcell.ColorLightBlue).
			SetMaxWidth(40).
			SetExpansion(0))

		table.SetCell(c+1, 1, tview.NewTableCell(dataRec.DisplayValue).
			SetTextColor(tcell.ColorLightBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		c++
	}

	table.Select(0, 0)
	i.ScrollToBeginning()
}

func (i *datas) updateEntries(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		i.setEntries(t, enterFolder)
	})
}

func (i *datas) focus(t *Tui) {
	i.SetSelectable(true, false)
	t.app.SetFocus(i)
}

func (i *datas) unfocus() {
	i.SetSelectable(false, false)
}

func (i *datas) setFilterWord(word string) {
	i.filterWord = word
}

func (i *datas) setFilterType(which int) {
	i.showTypes = which
}

func (i *datas) setShownPath(path string) {
	i.shownPath = path
}

func (i *datas) getShownPath() string {
	return i.shownPath
}

func (i *datas) setParent(path string) {
	i.parent = path
}

func (i *datas) getParent() string {
	return i.parent
}
