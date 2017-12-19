package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/util"
)

func main() {
	dir := filepath.Join("..", "..", "scripts")
	data := gnfinder.LoadTrainingData(filepath.Join(dir, "training"))
	output := filepath.Join(dir, "data", "nlp")
	d := dict.LoadDictionary()
	for lang, v := range *data {
		path := filepath.Join(output, lang.String(), "bayes.json")
		nb := gnfinder.Train(path, v, &d)
		err := ioutil.WriteFile(path, nb.Dump(), 0644)
		util.Check(err)
	}
}
