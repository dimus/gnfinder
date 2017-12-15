package lang

import (
	"github.com/abadojack/whatlanggo"
)

// Language represents the language of a text.
type Language int

// Set of languages. Last one is an indicator of the 'edge', as well as
// a default value for GnFinder.Language.
const (
	UnknownLanguage Language = iota
	English
	NotSet
)

func (l Language) String() string {
	languages := [...]string{"other", "eng", ""}
	return languages[l]
}

// LanguagesSet returns a 'set' of languages for more effective
// lookup of a language.
func LanguagesSet() map[Language]struct{} {
	var empty struct{}
	ls := make(map[Language]struct{})
	for i := 0; i < int(NotSet); i++ {
		ls[Language(i)] = empty
	}
	return ls
}

// DetectLanguage finds the most probable language for a text.
func DetectLanguage(text []rune) Language {
	sampleLength := len(text)
	if sampleLength > 20000 {
		sampleLength = 20000
	}
	info := whatlanggo.Detect(string(text[0:sampleLength]))
	code := whatlanggo.LangToString(info.Lang)
	switch code {
	case "eng":
		return English
	default:
		return UnknownLanguage
	}
}
