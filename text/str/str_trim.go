package str

import (
	"strings"
)

var (
	defaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.Trim(str, defaultTrimChars)
	} else {
		return strings.Trim(str, defaultTrimChars+characterMask[0])
	}
}

func TrimStr(str string, cut string) string {
	return TrimLeftStr(TrimRightStr(str, cut), cut)
}

func TrimLeft(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeft(str, defaultTrimChars)
	} else {
		return strings.TrimLeft(str, defaultTrimChars+characterMask[0])
	}
}

func TrimLeftStr(str string, cut string) string {
	var lenCut = len(cut)
	for len(str) >= lenCut && str[0:lenCut] == cut {
		str = str[lenCut:]
	}
	return str
}

func TrimRight(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRight(str, defaultTrimChars)
	} else {
		return strings.TrimRight(str, defaultTrimChars+characterMask[0])
	}
}

func TrimRightStr(str string, cut string) string {
	var lenStr = len(str)
	var lenCut = len(cut)
	for lenStr >= lenCut && str[lenStr-lenCut:lenStr] == cut {
		lenStr = lenStr - lenCut
		str = str[:lenStr]

	}
	return str
}
