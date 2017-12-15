package dict

import (
	"encoding/csv"
	"net/http"

	// _ is needed for virtual file system
	_ "github.com/gnames/gnfinder/statik"
	"github.com/gnames/gnfinder/util"
	"github.com/rakyll/statik/fs"
)

// DictionaryType describes available dictionaries
type DictionaryType int

func (d DictionaryType) String() string {
	types := [...]string{"NotInDictionary", "WhiteGenus", "GreyGenus",
		"WhiteUninomial", "GreyUninomial", "BlackUninomial", "WhiteSpecies",
		"GreySpecies", "BlackSpecies"}
	return types[d]
}

// DictionaryType dictionaries
const (
	NotInDictionary DictionaryType = iota
	WhiteGenus
	GreyGenus
	WhiteUninomial
	GreyUninomial
	BlackUninomial
	WhiteSpecies
	GreySpecies
	BlackSpecies
)

// Dictionary contains dictionaries used for detecting scientific names
type Dictionary struct {
	BlackUninomials map[string]struct{}
	BlackSpecies    map[string]struct{}
	GreyGenera      map[string]struct{}
	GreyGeneraSp    map[string]struct{}
	GreySpecies     map[string]struct{}
	GreyUninomials  map[string]struct{}
	WhiteGenera     map[string]struct{}
	WhiteSpecies    map[string]struct{}
	WhiteUninomials map[string]struct{}
}

// LoadDictionary contain most popular words in European languages.
func LoadDictionary() Dictionary {
	statikFS, err := fs.New()
	util.Check(err)

	d := Dictionary{}
	d.BlackUninomials = readData(statikFS, "/black/uninomials.csv")
	d.BlackSpecies = readData(statikFS, "/black/species.csv")
	d.GreyGenera = readData(statikFS, "/grey/genera.csv")
	d.GreyGeneraSp = readData(statikFS, "/grey/genera_species.csv")
	d.GreySpecies = readData(statikFS, "/grey/species.csv")
	d.GreyUninomials = readData(statikFS, "/grey/uninomials.csv")
	d.WhiteGenera = readData(statikFS, "/white/genera.csv")
	d.WhiteSpecies = readData(statikFS, "/white/species.csv")
	d.WhiteUninomials = readData(statikFS, "/white/uninomials.csv")
	return d
}

func readData(fs http.FileSystem, path string) map[string]struct{} {
	res := make(map[string]struct{})
	f, err := fs.Open(path)
	var empty struct{}
	util.Check(err)

	defer func() {
		err := f.Close()
		util.Check(err)
	}()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	util.Check(err)

	for _, v := range records {
		res[v[0]] = empty
	}
	return res
}
