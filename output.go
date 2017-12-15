package gnfinder

import (
	"bytes"
	"time"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
	jsoniter "github.com/json-iterator/go"
)

// Output type is the result of name-finding.
type Output struct {
	Meta  `json:"metadata"`
	Names []Name `json:"names"`
}

// Meta contains meta-information of name-finding result.
type Meta struct {
	// Date represents time when output was generated.
	Date time.Time `json:"date"`
	// Language of the document
	Language string `json:"language"`
	// TotalTokens is a number of 'normalized' words in the text
	TotalTokens int `json:"total_words"`
	// TotalNameCandidates is a number of words that might be a start of
	// a scientific name
	TotalNameCandidates int `json:"total_candidates"`
	// TotalNames is a number of scientific names found
	TotalNames int `json:"total_names"`
	// CurrentName (optional) is the index of the names array that designates a
	// "position of a cursor". It is used by programs like gntagger that allow
	// to work on the list of found names interactively.
	CurrentName int `json:"current_index,omitempty"`
}

// OddsDatum is a simplified version of a name, that stores boolean decision
// (Name/NotName), and corresponding odds of the name.
type OddsDatum struct {
	Name bool
	Odds float64
}

// Name represents one found name.
type Name struct {
	Type        string            `json:"type"`
	Verbatim    string            `json:"verbatim"`
	Name        string            `json:"name"`
	Odds        float64           `json:"odds,omitempty"`
	Likelihoods bayes.Likelihoods `json:"likelihoods,omitempty"`
	OffsetStart int               `json:"start"`
	OffsetEnd   int               `json:"end"`
	Annotation  string            `json:"annotation"`
}

// ToJSON converts Output to JSON representation.
func (o *Output) ToJSON() []byte {
	res, err := jsoniter.MarshalIndent(o, "", "  ")
	util.Check(err)
	return res
}

// FromJSON converts JSON representation of Outout to Output object.
func (o *Output) FromJSON(data []byte) {
	r := bytes.NewReader(data)
	err := jsoniter.NewDecoder(r).Decode(o)
	util.Check(err)
}

// NewOutput is a constructor for Output type.
func NewOutput(names []Name, tokens []token.Token, gnf *GnFinder) Output {
	m := Meta{
		Date:                time.Now(),
		TotalTokens:         len(tokens),
		TotalNameCandidates: candidatesNum(tokens),
		TotalNames:          len(names),
	}
	o := Output{Meta: m, Names: names}
	return o
}
