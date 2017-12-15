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

	conf := util.NewConfig(opts...)
	if conf.Language == lang.NotSet {
		conf.Language = lang.DetectLanguage(text)
		if conf.Language != lang.UnknownLanguage {
			conf.Bayes = true
		}
	}

	heuristic.TagTokens(tokens, d, text, conf)
	if conf.Bayes {
		nlp.TagTokens(tokens, d, conf)
	}

	return collectOutput(tokens, text, conf)
}

func collectOutput(ts []token.Token, text []rune,
	conf *util.Config) Output {
	return Output{}
}
