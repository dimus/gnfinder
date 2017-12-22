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
			fs := features(BayesFeatures(ts2))
			priorOdds := nameFrequency()
			odds := predictOdds(nb, t, fs, priorOdds)
			processBayesResults(odds, ts, i, m.BayesOddsThreshold)
		}
	}
}

func predictOdds(nb *bayes.NaiveBayes, t *token.Token, fs []BayesF,
	odds bayes.LabelFreq) []bayes.Posterior {
	oddsUni, err := nb.Predict(features(fUni), bayes.WithPriorOdds(odds))
	util.Check(err)
	if t.SpeciesIndexOffset > 0 {
		oddsSp, err := nb.Predict(features(fSp), bayes.WithPriorOdds(odds))
		util.Check(err)
		return []bayes.Posterior{oddsUni, oddsSp}
	}
	return []bayes.Posterior{oddsUni}
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
