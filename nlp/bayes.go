package nlp

import (
	"fmt"
	"io/ioutil"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
	"github.com/rakyll/statik/fs"
)

func TagTokens(ts []token.Token, d *dict.Dictionary, m *util.Model) {
	nb := naiveBayesFromDump(m)
	l := len(ts)
	for i := range ts {
		t := &ts[i]
		if t.Features.Capitalized && t.UninomialDict != dict.BlackUninomial {
			ts2 := ts[i:util.UpperIndex(i, l)]
			fs := BayesFeatures(ts2)
			priorOdds := nameFrequency()
			odds := predictOdds(nb, t, &fs, priorOdds)
			processBayesResults(odds, ts, i, m.BayesOddsThreshold)
		}
	}
}

func processBayesResults(odds []bayes.Posterior, ts []token.Token, i int,
	oddsThreshold float64) {
	uni := &ts[i]
	decideUninomial(odds, uni, oddsThreshold)

	if uni.Indices.Species == 0 || odds[1].MaxLabel != Name {
		return
	}

	sp := &ts[i+uni.Indices.Species]
	decideSpeces(odds, uni, sp, oddsThreshold)
	if uni.Indices.Infraspecies == 0 || odds[2].MaxLabel != Name {
		return
	}
	isp := &ts[i+uni.Indices.Infraspecies]
	decideInfraspeces(odds, uni, isp, oddsThreshold)
}

func decideInfraspeces(odds []bayes.Posterior, uni *token.Token,
	isp *token.Token, oddsThreshold float64) {
	if isp.SpeciesDict == dict.BlackSpecies {
		return
	}
	isp.Odds = odds[2].MaxOdds
	isp.OddsDetails = token.NewOddsDetails(odds[2].Likelihoods)

	if isp.Odds >= oddsThreshold && uni.Odds > 1 &&
		uni.Decision.In(token.NotName, token.PossibleBinomial) {
		uni.Decision = token.BayesTrinomial
	}
}

func decideSpeces(odds []bayes.Posterior, uni *token.Token, sp *token.Token,
	oddsThreshold float64) {
	if sp.SpeciesDict == dict.BlackSpecies {
		return
	}
	sp.Odds = odds[1].MaxOdds
	sp.OddsDetails = token.NewOddsDetails(odds[1].Likelihoods)
	if sp.Odds >= oddsThreshold && uni.Odds > 1 &&
		uni.Decision.In(token.NotName, token.PossibleBinomial) {
		uni.Decision = token.BayesBinomial
	}
}

func decideUninomial(odds []bayes.Posterior, uni *token.Token,
	oddsThreshold float64) {
	if odds[0].MaxLabel == Name {
		uni.Odds = odds[0].MaxOdds
	} else {
		uni.Odds = 1 / odds[0].MaxOdds
	}
	uni.OddsDetails = token.NewOddsDetails(odds[0].Likelihoods)
	uni.LabelFreq = odds[0].LabelFreq

	if odds[0].MaxLabel == Name &&
		odds[0].MaxOdds >= oddsThreshold &&
		uni.Decision == token.NotName &&
		uni.UninomialDict != dict.BlackUninomial &&
		!uni.Abbr {
		uni.Decision = token.BayesUninomial
	}
}

func predictOdds(nb *bayes.NaiveBayes, t *token.Token, fs *FeatureSet,
	odds bayes.LabelFreq) []bayes.Posterior {
	oddsUni, err := nb.Predict(features(fs.Uninomial), bayes.WithPriorOdds(odds))
	util.Check(err)
	if t.Indices.Species == 0 {
		return []bayes.Posterior{oddsUni}
	}

	oddsSp, err := nb.Predict(features(fs.Speces), bayes.WithPriorOdds(odds))
	util.Check(err)
	if t.Indices.Infraspecies == 0 {
		return []bayes.Posterior{oddsUni, oddsSp}
	}
	f := features(fs.InfraSp)
	oddsInfraSp, err := nb.Predict(f, bayes.WithPriorOdds(odds))
	util.Check(err)
	return []bayes.Posterior{oddsUni, oddsSp, oddsInfraSp}
}

func nameFrequency() bayes.LabelFreq {
	return map[bayes.Labeler]float64{
		Name:    1.0,
		NotName: 10.0,
	}
}

func naiveBayesFromDump(m *util.Model) *bayes.NaiveBayes {
	nb := bayes.NewNaiveBayes()
	bayes.RegisterLabel(labelMap)
	staticFS, err := fs.New()
	util.Check(err)

	dir := fmt.Sprintf("/nlp/%s/bayes.json", m.Language.String())
	f, err := staticFS.Open(dir)
	util.Check(err)

	defer func() {
		err := f.Close()
		util.Check(err)
	}()

	json, err := ioutil.ReadAll(f)
	util.Check(err)
	nb.Restore(json)
	return nb
}

func features(bf []BayesF) []bayes.Featurer {
	f := make([]bayes.Featurer, len(bf))
	for i, v := range bf {
		f[i] = v
	}
	return f
}
