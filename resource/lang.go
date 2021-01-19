package resource

import (
	"errors"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// Lang - Language Interface
type Lang struct {
	lang string
}

// Main UI Text
var ui = map[string]map[string]string{
	"American English": {
		"Folder":                   "Folder",
		"Data":                     "Data",
		"FIELD":                    "FIELD",
		"VALUE":                    "VALUE",
		"TYPE":                     "TYPE",
		"PATH":                     "PATH",
		"VERSION":                  "VERSION",
		"Successfully Connected":   "Successfully Connected",
		"Copied to clipboard":      "Copied to clipboard",
		"Only Valid for Data Type": "Only Valid for Data Type",
		"m_fld_s":                  "s",
		"m_fld_s_desc":             "show/hide values",
		"m_fld_bB":                 "b/B",
		"m_fld_bB_desc":            "generate Bash script",
		"m_fld_jJ":                 "j/J",
		"m_fld_jJ_desc":            "generate Jenkinsfile",
		"m_fld_cC":                 "c/C",
		"m_fld_cC_desc":            "FIELD/CMD -> clipboard",
		"m_fld_/":                  "/",
		"m_fld_/_desc":             "filter",
		"m_fld_v":                  "v",
		"m_fld_v_desc":             "VALUE -> clipboard",
		"m_fld_?":                  "?",
		"m_fld_?_desc":             "help",
	},
	"German": {
		"Folder":                   "Mappe",
		"Data":                     "Daten",
		"FIELD":                    "FELD",
		"VALUE":                    "WERT",
		"TYPE":                     "TYP",
		"PATH":                     "PFAD",
		"VERSION":                  "VERSION",
		"Successfully Connected":   "Verbindung hergestellt",
		"Copied to clipboard":      "In die Zwischenablage kopiert",
		"Only Valid for Data Type": "(en)Only Valid for Data Type",
		"m_fld_s":                  "s",
		"m_fld_s_desc":             "Werte anzeigen/ausblenden",
		"m_fld_bB":                 "b/B",
		"m_fld_bB_desc":            "Bash Skript erstellen",
		"m_fld_jJ":                 "j/J",
		"m_fld_jJ_desc":            "Jenkinsfile erstellen",
		"m_fld_cC":                 "c/C",
		"m_fld_cC_desc":            "Schlüssel/Befehl -> Zwischenablage",
		"m_fld_/":                  "/",
		"m_fld_/_desc":             "filter",
		"m_fld_v":                  "v",
		"m_fld_v_desc":             "WERT -> Zwischenablage",
		"m_fld_?":                  "?",
		"m_fld_?_desc":             "hilfe",
	},
}

// Keep Debugging Translations Separate, they don't really matter.
// Potential future use.
var dbg = map[string]map[string]string{
	"American English": {
		"Folder": "Folder",
	},
	"German": {
		"Folder": "Mappe",
	},
}

// Keyboard Interface, potentially use different keys for different languages
// Potential future use
var kbd = map[string]map[string]rune{
	"American English": {
		"s": 's',
	},
	"German": {
		"s": 's',
	},
}

var appLangs = []language.Tag{
	language.AmericanEnglish, // en-US fallback
	language.German,          // de
}

var matcher = language.NewMatcher(appLangs)

// NewLanguage - Create a new language instance
func NewLanguage() *Lang {

	tag, _, _ := matcher.Match(language.Make(os.Getenv("LANG")))
	return &Lang{lang: display.English.Tags().Name(tag)}

}

// GetRune - Return runes for keybindings
func (l *Lang) GetRune(resource string, key string) (rune, error) {
	switch resource {
	case "kbd":
		if val, ok := kbd[l.lang]; ok {
			if text, ok := val[key]; ok {
				return text, nil
			}
			return '^', errors.New("No rune defined")
		}
		return '^', errors.New("No language defined")
	}
	return '^', errors.New("No resource defined")
}

// GetText - Return text for a language resource
func (l *Lang) GetText(resource string, key string) string {
	switch resource {
	case "ui":
		if val, ok := ui[l.lang]; ok {
			if text, ok := val[key]; ok {
				return text
			}
			return "No Language Key"
		}
		return "No Language"
	case "dbg":
		if val, ok := dbg[l.lang]; ok {
			if text, ok := val[key]; ok {
				return text
			}
			return "No Language Key"
		}
		return "No Language"
	}
	return "No Resource Type"
}
