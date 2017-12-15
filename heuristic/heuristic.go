package heuristic

import (
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
)

func TagTokens(ts []token.Token, d *dict.Dictionary, text []rune,
	conf *util.Config) {
	l := len(ts)

	for i := range ts {
		t := &ts[i]

		if !t.Features.Capitalized {
			continue
		}

		t.Features.SetUninomialDict(t, d)
		upperIndex := i + 4
		if l < upperIndex {
			upperIndex = l
		}
		exploreNameCandidate(ts[i:upperIndex], d, text, conf)
	}
}

func exploreNameCandidate(ts []token.Token, d *dict.Dictionary, text []rune,
	conf *util.Config) {
}
