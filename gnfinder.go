//go:generate statik -f -src=./data/dictionaries
package gnfinder

import (
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/heuristic"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/nlp"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
)

// FindNamesJSON takes a text and returns scientific names found in the text,
// as well as tokens
func FindNamesJSON(data []byte, dict *dict.Dictionary,
	opts ...util.Opt) []byte {
	output := FindNames([]rune(string(data)), dict, opts...)
	return output.ToJSON()
}

// FindNames traverses a text and finds scientific names in it.
func FindNames(text []rune, d *dict.Dictionary, opts ...util.Opt) Output {
	tokens := token.Tokenize(text)

	m := util.NewModel(opts...)
	if m.Language == lang.NotSet {
		m.Language = lang.DetectLanguage(text)
		if m.Language != lang.UnknownLanguage {
			m.Bayes = true
		}
	}

	heuristic.TagTokens(tokens, d, m)
	if m.Bayes {
		nlp.TagTokens(tokens, d, m)
	}

	return collectOutput(tokens, text, m)
}

func collectOutput(ts []token.Token, text []rune,
	m *util.Model) Output {
	var names []Name
	for i := range tokens {
		u := &tokens[i]
		if u.Kind == token.NotName {
			continue
		}
		s := &tokens[i+u.SpeciesIndexOffset]
		name := TokensToName(u, s, t)
		names = append(names, name)
	}
	output := NewOutput(names, tokens, gnf)
	return output
}
