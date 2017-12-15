package token

import (
	"strings"
	"unicode"

	"github.com/gnames/gnfinder/dict"
)

// Features keep properties of a token as a possible candidate for a name part.
type Features struct {
	// Candidate to be a start of a uninomial or binomial.
	NameStartCandidate bool
	// The name looks like a possible genus name.
	PotentialBinomialGenus bool
	// The token has necessary qualities to be a start of a binomial.
	PotentialBinomialSpecies bool
	// Capitalized feature of the first alphabetic character.
	Capitalized bool
	// CapitalizedSpecies -- the first species lphabetic character is capitalized.
	CapitalizedSpecies bool
	// ParensEnd feature: token starts with parentheses.
	ParensStart bool
	// ParensEnd feature: token ends with parentheses.
	ParensEnd bool
	// ParensEndSpecies feature: species token ends with parentheses.
	ParensEndSpecies bool
	// Abbr feature: token ends with a period.
	Abbr bool
	// UninomialDict defines which Genera or Uninomials dictionary (if any)
	// contained the token.
	UninomialDict dict.DictionaryType
	// SpeciesDict defines which Species dictionary (if any) contained the token.
	SpeciesDict dict.DictionaryType
}

func (f *Features) setParensStart(firstRune rune) {
	f.ParensStart = firstRune == rune('(')
}

func (f *Features) setParensEnd(lastRune rune) {
	f.ParensEnd = lastRune == rune(')')
}

func (f *Features) setCapitalized(firstAlphabetRune rune) {
	f.Capitalized = unicode.IsUpper(firstAlphabetRune)
}

func (f *Features) setAbbr(raw []rune, startEnd *[2]int) {
	l := len(raw)
	lenClean := startEnd[1] - startEnd[0] + 1
	if lenClean < 4 && l > 1 && unicode.IsLetter(raw[l-2]) &&
		raw[l-1] == rune('.') {
		f.Abbr = true
	}
}

func (f *Features) setPotentialBinomialGenus(startEnd *[2]int, raw []rune) {
	lenRaw := len(raw)
	lenClean := startEnd[1] - startEnd[0] + 1
	cleanEnd := lenRaw == startEnd[1]+1
	switch lenClean {
	case 0:
		f.PotentialBinomialGenus = false
	case 1:
		f.PotentialBinomialGenus = f.Abbr
	case 2, 3:
		f.PotentialBinomialGenus = f.Abbr || cleanEnd
	default:
		f.PotentialBinomialGenus = cleanEnd
	}
}

func (f *Features) setPotentialBinomialSpecies(startEnd *[2]int) {
	lenClean := startEnd[1] - startEnd[0] + 1
	if lenClean >= 2 && startEnd[0] == 0 {
		f.PotentialBinomialSpecies = true
	}
}

func (f *Features) SetUninomialDict(t *Token, d *dict.Dictionary) {
	name := t.Cleaned
	in := func(dict map[string]struct{}) bool { _, ok := dict[name]; return ok }
	inlow := func(dict map[string]struct{}) bool {
		_, ok := dict[strings.ToLower(name)]
		return ok
	}

	switch {
	case in(d.WhiteGenera):
		f.UninomialDict = dict.WhiteGenus
	case in(d.GreyGenera):
		f.UninomialDict = dict.GreyGenus
	case in(d.WhiteUninomials):
		f.UninomialDict = dict.WhiteUninomial
	case in(d.GreyUninomials):
		f.UninomialDict = dict.GreyUninomial
	case inlow(d.BlackUninomials):
		f.UninomialDict = dict.BlackUninomial
	}
}

func (f *Features) SetSpeciesDict(t *Token, d *dict.Dictionary) {
	name := strings.ToLower(t.Cleaned)
	in := func(dict map[string]struct{}) bool { _, ok := dict[name]; return ok }

	switch {
	case in(d.WhiteSpecies):
		f.SpeciesDict = dict.WhiteSpecies
	case in(d.GreySpecies):
		f.SpeciesDict = dict.GreySpecies
	case in(d.BlackSpecies):
		f.SpeciesDict = dict.BlackSpecies
	}
}
