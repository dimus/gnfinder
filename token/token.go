// Package token deals with breaking a text into tokens. It cleans names broken
// by new lines, concatenating pieces together. Tokens are connected to
// features. Features are used for heuristic and Bayes' approaches for finding
// names.
package token

import (
	"unicode"

	"github.com/gnames/bayes"
)

// Token represents a word separated by spaces in a text. Words split by new
// lines are concatenated.
type Token struct {
	// Raw is a verbatim presentation of a token as it appears in a text.
	Raw []rune
	// Cleaned is a presentation of a token after normalization.
	Cleaned string
	// Start is the index of the first rune of a token. The first rune
	// does not have to be alpha-numeric.
	Start int
	// End is the index of the last rune of a token. The last rune does not
	// have to be alpha-numeric.
	End int
	// Decision tags the first token of a possible name with a classification
	// decision.
	Decision
	// Indices of semantic elements of a possible name.
	Indices
	// NLP data
	NLP
	// Features is a collection of features associated with the token
	Features
}

// Decision definds possible kinds of name candidates.
type Decision int

// Possible Decisions
const (
	NotName Decision = iota
	Uninomial
	Binomial
	PossibleBinomial
	BayesUninomial
	BayesBinomial
)

var decisionsStrings = [...]string{"NotName", "Uninomial", "Binomial",
	"PossibleBinomial", "Uninomial(nlp)", "Binomial(nlp)"}

// String representation of a Decision
func (h Decision) String() string {
	return decisionsStrings[h]
}

// In returns true if a Decision is included in given constants.
func (h Decision) In(hds ...Decision) bool {
	for _, hd := range hds {
		if h == hd {
			return true
		}
	}
	return false
}

// Indices of the elmements for a name candidate.
type Indices struct {
	Species      int
	Rank         int
	Infraspecies int
}

// NLP collects data received from Bayes' algorithm
type NLP struct {
	// Odds are posterior odds.
	Odds float64
	// OddsDetails are elements from which Odds are calculated.
	OddsDetails bayes.Likelihoods
	// LabelFreq is used to calculate prior odds of names appearing in a
	// document
	LabelFreq bayes.LabelFreq
}

// NewToken constructs a new Token object.
func NewToken(text []rune, start int, end int) Token {
	t := Token{
		Raw:   text[start:end],
		Start: start,
		End:   end,
	}
	t.Clean()
	return t
}

// Clean converts a verbatim (Raw) string of a token into normalized cleaned up
// version.
func (t *Token) Clean() {
	l := len(t.Raw)
	f := &t.Features

	f.setParensStart(t.Raw[0])
	f.setParensEnd(t.Raw[l-1])

	res, startEnd := t.normalize()

	f.setAbbr(t.Raw, startEnd)
	if t.Features.Capitalized {
		res[0] = unicode.ToUpper(res[0])
		f.setPotentialBinomialGenus(startEnd, t.Raw)
	} else {
		// makes it impossible to have capitalized species
		f.setPotentialBinomialSpecies(startEnd)
	}

	if f.Abbr {
		res = append(res, rune('.'))
	}
	t.Cleaned = string(res)
}

func (t *Token) normalize() ([]rune, *[2]int) {
	var res []rune
	firstLetter := true
	var startEnd [2]int
	for i, v := range t.Raw {
		if unicode.IsLetter(v) || v == rune('-') {
			if firstLetter {
				startEnd[0] = i
				t.Features.setCapitalized(v)
				firstLetter = false
			}
			startEnd[1] = i
			res = append(res, unicode.ToLower(v))
		} else {
			res = append(res, rune('�'))
		}
	}
	return res[startEnd[0] : startEnd[1]+1], &startEnd
}

// InParentheses is true if token is surrounded by parentheses.
func (t *Token) InParentheses() bool {
	if t.Features.ParensStart && t.Features.ParensEnd {
		return true
	}
	return false
}
