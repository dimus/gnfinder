package nlp

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/token"
	"github.com/gnames/gnfinder/util"
	jsoniter "github.com/json-iterator/go"
)

// NameString represents name-string of a name.
type NameString string

// TrainingNames keeps meta-information associated with NameString.
type TrainingNames map[NameString]*TrainingNameMeta

// TrainingData associates a Language with training data
type TrainingData map[lang.Language]*TrainingLanguageData

// TrainingLanguageData keeps data from text that does not contain names
// and data with names. This will be obsolete when we have better way to do
// training.
type TrainingLanguageData struct {
	NoNamesText []rune
	*TrainingNames
}

// TrainingNameMeta keeps meta-information about a NameString
type TrainingNameMeta struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Type  string `json:"type"`
}

// LoadTrainingData loads TrainingData from a file.
func LoadTrainingData(dir string) *TrainingData {
	td := make(TrainingData)
	for i := 1; i < int(lang.NotSet); i++ {
		lang := lang.Language(i)
		path := filepath.Join(dir, lang.String())
		textPtr := loadText(path)
		namesPtr := loadNames(path)
		td[lang] = &TrainingLanguageData{textPtr, namesPtr}
	}
	return &td
}

// Train runs training process
func Train(path string, data *TrainingLanguageData,
	d *dict.Dictionary) *bayes.NaiveBayes {
	lfs := processNames(data.TrainingNames, d)
	lfs = append(lfs, processNoNames(data.NoNamesText, d)...)
	nb := bayes.TrainNB(lfs)
	return nb
}

func processNames(ns *TrainingNames,
	d *dict.Dictionary) []bayes.LabeledFeatures {
	var lfs []bayes.LabeledFeatures
	label := Name
	for _, v := range *ns {
		text := []rune(v.Text)
		ts := token.Tokenize(text)
		l := len(ts)
		for i := range ts {
			t := &ts[i]
			if (t.Start <= v.Start && t.End > v.Start) || len(v.Text) == 0 {
				ts2 := ts[i:util.UpperIndex(i, l)]
				token.SetIndices(ts2, d)
				fs := BayesFeatures(ts2)
				lfs = append(lfs, bayes.LabeledFeatures{
					Features: fs.Flatten(),
					Label:    label,
				})
			}
		}
	}
	return lfs
}

func processNoNames(t []rune, d *dict.Dictionary) []bayes.LabeledFeatures {
	var lfs []bayes.LabeledFeatures
	label := NotName
	ts := token.Tokenize(t)
	l := len(ts)
	for i := range ts {
		t := &ts[i]
		if t.Features.Capitalized {
			ts2 := ts[i:util.UpperIndex(i, l)]
			token.SetIndices(ts2, d)
			fs := BayesFeatures(ts2)
			lfs = append(lfs, bayes.LabeledFeatures{
				Features: fs.Flatten(),
				Label:    label,
			})
		}
	}
	return lfs
}

func loadText(path string) []rune {
	path = filepath.Join(path, "no_names.txt")
	bytes, err := ioutil.ReadFile(path)
	text := []rune(string(bytes))
	util.Check(err)
	return text
}

func loadNames(path string) *TrainingNames {
	tn := make(TrainingNames)
	path = filepath.Join(path, "names.json")
	text, err := ioutil.ReadFile(path)
	util.Check(err)
	r := bytes.NewReader(text)
	err = jsoniter.NewDecoder(r).Decode(&tn)
	util.Check(err)
	return &tn
}
