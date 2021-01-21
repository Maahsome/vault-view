package tui

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

// TestNewInfoLogo - testing newInfoLogo
func TestNewInfoLogo(t *testing.T) {

	tv := newInfoLogo("v0.0.1")

	logo := ` ___   __
 __ | / / %s
 __ |/ /___   __
 _____/ __ | / /
        __ |/ /
        _____/
`
	assert.Equal(t, tv.GetText(true), fmt.Sprintf(logo, "v0.0.1"))

}
