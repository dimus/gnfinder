package nlp

import (
	"strconv"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/token"
)

// BayesF implements bayes.Featurer
type BayesF struct {
	name  string
	value string
}

// FeatureSet splits features into Uninomial, Species, Ifraspecies groups
type FeatureSet struct {
	Uninomial []BayesF
	Speces    []BayesF
	InfraSp   []BayesF
}

func (fs *FeatureSet) Flatten() []bayes.Featurer {
	l := len(fs.Uninomial) + len(fs.Speces) + len(fs.InfraSp)
	res := make([]bayes.Featurer, 0, l)
	res = append(res, features(fs.Uninomial)...)
	res = append(res, features(fs.Speces)...)
	res = append(res, features(fs.InfraSp)...)
	return res
}

// Name is required by bayes.Featurer
func (b BayesF) Name() bayes.FeatureName { return bayes.FeatureName(b.name) }

// Value is required by bayes.Featurer
func (b BayesF) Value() bayes.FeatureValue {
	return bayes.FeatureValue(b.value)
}

// BayesFeatures creates slices of features for a token that might represent
// genus or other uninomial
func BayesFeatures(ts []token.Token) FeatureSet {
	var fs FeatureSet
	var u, sp, isp *token.Token
	u = &ts[0]

	if !u.Capitalized {
		return fs
	}

	if i := u.Indices.Species; i > 0 {
		sp = &ts[i]
	}

	if i := u.Indices.Infraspecies; i > 0 {
		isp = &ts[i]
	}

	return convertFeatures(u, sp, isp)
}

func convertFeatures(u *token.Token, sp *token.Token,
	isp *token.Token) FeatureSet {
	var fs FeatureSet
	fs.Uninomial = []BayesF{
		{"uniAbbr", strconv.FormatBool(u.Abbr)},
	}
	if !u.Abbr {
		fs.Uninomial = append(fs.Uninomial,
			BayesF{"uniLen", strconv.Itoa(len(u.Cleaned))},
			BayesF{"uniDict", u.UninomialDict.String()},
		)
	}
	if w3 := wordEnd(u); !u.Abbr && w3 != "" {
		fs.Uninomial = append(fs.Uninomial, BayesF{"uniEnd3", w3})
	}
	if u.Indices.Species > 0 {
		fs.Speces = append(fs.Speces,
			BayesF{"spLen", strconv.Itoa(len(sp.Cleaned))},
			BayesF{"spDict", sp.SpeciesDict.String()},
		)
		if w3 := wordEnd(sp); w3 != "" {
			fs.Speces = append(fs.Speces, BayesF{"spEnd3", w3})
		}
	}
	if u.Indices.Rank > 0 {
		fs.InfraSp = []BayesF{
			{"ispRank", "true"},
		}
	}
	if u.Indices.Infraspecies > 0 {
		fs.InfraSp = append(fs.InfraSp,
			BayesF{"ispLen", strconv.Itoa(len(isp.Cleaned))},
			BayesF{"ispDict", isp.SpeciesDict.String()},
		)
	}
	return fs
}

func wordEnd(t *token.Token) string {
	name := []rune(t.Cleaned)
	l := len(name)
	if l < 4 {
		return ""
	}
	w3 := string(name[l-3 : l])
	return w3
}
