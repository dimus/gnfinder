package heuristic

import (
	"fmt"
	"strings"

	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
)

func TagTokens(ts []token.Token, d *dict.Dictionary, m *util.Model) {
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
		exploreNameCandidate(ts[i:upperIndex], d, m)
	}
}

func exploreNameCandidate(ts []token.Token, d *dict.Dictionary,
	m *util.Model) bool {

	u := &ts[0]

	if u.Features.UninomialDict == dict.WhiteUninomial {
		u.Decision = token.Uninomial
		return true
	} else if u.InParentheses() &&
		u.Features.UninomialDict == dict.WhiteGenus {
		u.Decision = token.Uninomial
		return true
	}

	i := speciesTokenIndex(ts, d)

	if i == 0 && u.Features.UninomialDict == dict.WhiteGenus {
		ts[0].Decision = token.Uninomial
		return true
	}
	if _, blk := d.BlackUninomials[strings.ToLower(u.Cleaned)]; blk {
		return false
	}

	u.Indices.Species = i

	if ok := checkAsGenusSpecies(ts, d, m); !ok {
		return false
	}

	return checkInfraspecies(ts, d, m)
}

func speciesTokenIndex(ts []token.Token, d *dict.Dictionary) (i int) {
	i = 2
	if len(ts) > i && ts[1].InParentheses() {
		ts[i].SetSpeciesDict(&ts[i], d)
		if checkAsSpecies(&ts[i], d) {
			return
		}
	}

	i = 1
	if len(ts) > i {
		ts[i].SetSpeciesDict(&ts[i], d)
		if checkAsSpecies(&ts[i], d) {
			return i
		}
	}
	return 0
}

func checkAsSpecies(t *token.Token, d *dict.Dictionary) bool {
	f := &t.Features
	if !f.PotentialBinomialSpecies ||
		!(f.SpeciesDict == dict.WhiteSpecies || f.SpeciesDict == dict.GreySpecies) {
		return false
	}
	return true
}

func checkAsGenusSpecies(ts []token.Token, d *dict.Dictionary,
	m *util.Model) bool {
	g := &ts[0]
	s := &ts[g.Indices.Species]
	if g.UninomialDict == dict.WhiteGenus {
		g.Decision = token.Binomial
		return true
	}

	if checkGreyGeneraSp(g, s, d) {
		g.Decision = token.Binomial
		return true
	}

	if s.Features.SpeciesDict == dict.WhiteSpecies {
		g.Decision = token.PossibleBinomial
		return true
	}
	return false
}

func checkGreyGeneraSp(g *token.Token, s *token.Token,
	d *dict.Dictionary) bool {
	sp := fmt.Sprintf("%s %s", g.Cleaned, s.Cleaned)
	if _, ok := d.GreyGeneraSp[sp]; ok {
		return true
	}
	return false
}

func checkInfraspecies(ts []token.Token, d *dict.Dictionary,
	m *util.Model) bool {
	return false
}
