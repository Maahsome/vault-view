package tui

import (
	"github.com/maahsome/tview"
	"github.com/maahsome/vault-view/common"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/sirupsen/logrus"
)

type info struct {
	*tview.Flex
	Status *infoIndicator
	Menu   *infoMenu
}

func newInfo(t *Tui, semVer string) *info {

	common.Logger.WithFields(logrus.Fields{
		"unit":     "info",
		"function": "tui",
	}).Debug("Creating newInfo")
	w, _ := consolesize.GetConsoleSize()
	logo := newInfoLogo(semVer)
	menu := newInfoMenu(t)
	indicator := newInfoIndicator(t)

	i := &info{
		tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(logo, 20, 1, false).
			AddItem(menu, (((w-20)/4)*3), 3, false).
			AddItem(indicator, ((w - 20) / 4), 2, false),
		indicator,
		menu,
	}

	return i
}
