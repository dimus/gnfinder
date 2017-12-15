//go:generate statik -f -src=./data/dictionaries
package gnfinder

import (
	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/heuristic"
	"github.com/gnames/gnfinder/token"
)

// FindNamesJSON takes a text and returns scientific names found in the text,
// as well as tokens
func FindNamesJSON(data []byte, dict *dict.Dictionary, opts ...Opt) []byte {
	output := FindNames([]rune(string(data)), dict, opts...)
	return output.ToJSON()
}

// FindNames traverses a text and finds scientific names in it.
func FindNames(text []rune, d *dict.Dictionary, opts ...Opt) Output {
	tokens := token.Tokenize(text)

	conf := NewConfig(text, opts...)

	heuristic.TagTokens(tokens, d, text, conf)
	if conf.Bayes {
		bayes.TagTokens(tokens, d, conf)
	}

	return collectOutput(tokens, text, conf)
}

func collectOutput(ts []Token, text []rune, conf *Config) Output {
	return Output{}
}
