package tui

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
	"github.com/maahsome/vault-view/vault"
)

type infoLogo struct {
	*tview.TextView
	Vault  *vaultInfo
	Host   *hostInfo
	AppVer *appVer
}

type appVer struct {
	Name    string
	Version string
}

type vaultInfo struct {
	Addr    string
	Version string
}

type hostInfo struct {
	OSType       string
	Architecture string
}

func newInfoLogo(semVer string) *infoLogo {

	i := &infoLogo{
		TextView: tview.NewTextView(),
		Host:     newHostInfo(),
		AppVer:   newAppVerInfo(semVer),
	}

	i.display()

	return i
}

func newAppVerInfo(semVer string) *appVer {
	return &appVer{
		Name:    "vault-view",
		Version: semVer,
	}
}

func newHostInfo() *hostInfo {
	return &hostInfo{
		OSType:       runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
}

func newVaultInfo(vaultCli vault.Client) *vaultInfo {

	vv, _ := vaultCli.GetVersion()

	return &vaultInfo{
		Addr:    os.Getenv("VAULT_ADDR"),
		Version: vv,
	}
}

func (i *infoLogo) display() {

	appVersion := fmt.Sprintf("%s", i.AppVer.Version)

	i.SetDynamicColors(true).SetTextColor(tcell.ColorDarkOrange)

	logo := ` ___   __
 __ | / / %s
 __ |/ /___   __
 _____/ __ | / /
        __ |/ /
        _____/
`

	i.SetText(fmt.Sprintf(logo, appVersion))
}
