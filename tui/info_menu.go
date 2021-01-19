package tui

import (
	"fmt"

	"github.com/maahsome/tview"
	"github.com/maahsome/vault-view/common"
	"github.com/maahsome/vault-view/resource"
	"github.com/nathan-fiscaletti/consolesize-go"
)

type infoMenu struct {
	*tview.TextView
	Items map[string][]menuItem
	lang  *resource.Lang
}

type menuItem struct {
	KeyList     string
	Description string
}

func newInfoMenu(t *Tui) *infoMenu {

	i := &infoMenu{
		TextView: tview.NewTextView().SetDynamicColors(true),
		Items:    newMenuItems(t),
		lang:     t.lang,
	}

	i.display("folders")

	return i
}

func newMenuItems(t *Tui) map[string][]menuItem {

	menu := make(map[string][]menuItem)

	menu["folders"] = []menuItem{
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_s"),      // "s"
			Description: t.lang.GetText("ui", "m_fld_s_desc"), // "show/hide values"
		},
		// menuItem{
		// 	KeyList:     t.lang.GetText("ui", "m_fld_bB"),      // "b/B"
		// 	Description: t.lang.GetText("ui", "m_fld_bB_desc"), // "generate Bash script"
		// },
		// menuItem{
		// 	KeyList:     t.lang.GetText("ui", "m_fld_jJ"),      // "j/J"
		// 	Description: t.lang.GetText("ui", "m_fld_jJ_desc"), // "generate Jenkinsfile"
		// },
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_cC"),      // "c/C"
			Description: t.lang.GetText("ui", "m_fld_cC_desc"), // "KEY/CMD -> clipboard"
		},
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_/"),      // "/"
			Description: t.lang.GetText("ui", "m_fld_/_desc"), // "filter"
		},
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_?"),      // "/"
			Description: t.lang.GetText("ui", "m_fld_?_desc"), // "filter"
		},
	}

	menu["datas"] = []menuItem{
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_s"),      // "s"
			Description: t.lang.GetText("ui", "m_fld_s_desc"), // "show/hide values"
		},
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_cC"),      // "c/C"
			Description: t.lang.GetText("ui", "m_fld_cC_desc"), // "KEY/CMD -> clipboard"
		},
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_v"),      // "v"
			Description: t.lang.GetText("ui", "m_fld_v_desc"), // "VALUE -> clipboard"
		},
		menuItem{
			KeyList:     t.lang.GetText("ui", "m_fld_?"),      // "/"
			Description: t.lang.GetText("ui", "m_fld_?_desc"), // "filter"
		},
	}

	return menu
}

func (i *infoMenu) display(menu string) {

	w, _ := consolesize.GetConsoleSize()

	cols := 1
	if w > 80 && w < 160 {
		cols = 2
	}
	if w >= 160 {
		cols = 3
	}
	common.Logger.Info(fmt.Sprintf("Cols: %d, actual COLS: %d", w, cols))
	output := ""
	for c, m := range i.Items[menu] {
		separate := "  "
		if ((c + 1) % cols) == 0 {
			separate = "\n"
		}
		output += fmt.Sprintf("[purple]<%s>:[grey] %s%s", m.KeyList, m.Description, separate)
	}
	i.SetText(output)

}
