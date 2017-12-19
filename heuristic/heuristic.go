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
		exploreNameCandidate(ts[i:util.UpperIndex(i, l)], d, m)
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

	checkInfraspecies(ts, d, m)
	return true
}

func speciesTokenIndex(ts []token.Token, d *dict.Dictionary) (i int) {
	if !ts[0].PotentialBinomialGenus {
		return 0
	}

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

	if s.Features.SpeciesDict == dict.WhiteSpecies && !s.Capitalized {
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
	m *util.Model) {
	if !ts[ts[0].Indices.Species].Features.PotentialTrinomialSpecies {
		return
	}
	i := ts[0].Indices.Species + 2
	if len(ts) > i {
		ts[i].SetSpeciesDict(&ts[i], d)
		if _, ok := d.Ranks[ts[i-1].Cleaned]; ok && checkAsSpecies(&ts[i], d) {
			ts[0].Indices.Rank = i - 1
			setInfraspecies(&ts[0], i)
			return
		}
	}

	i--
	if len(ts) > i {
		ts[i].SetSpeciesDict(&ts[i], d)
		if checkAsSpecies(&ts[i], d) {
			setInfraspecies(&ts[0], i)
			return
		}
	}
	return
}

func setInfraspecies(g *token.Token, i int) bool {
	g.Decision = token.Trinomial
	g.Indices.Infraspecies = i
	return true
}
