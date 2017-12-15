package util

import (
	"runtime"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder/lang"
)

// Config keeps configuration variables
type Config struct {
	// Language of the text
	Language lang.Language
	// Bayes flag forces to run Bayes name-finding on unknown languages
	Bayes bool
	// BayesOddsThreshold sets the limit of posterior odds. Everything bigger
	// that this limit will go to the names output.
	BayesOddsThreshold float64
	// TextOdds captures "concentration" of names as it is found for the whole
	// text by heuristic name-finding. It should be close enough for real
	// number of names in text. We use it when we do not have local conentration
	// of names in a region of text.
	TextOdds bayes.LabelFreq
	// NameDistribution keeps data about position of names candidates and
	// their value according to heuristic and Bayes name-finding algorithms.
	// NameDistribution
	// ResolverConf
	Resolver
}

// Resolver contains configuration of Resolver data
type Resolver struct {
	URL       string
	BatchSize int
	Workers   int
}

// NewConfig creates Config object with default data, or with data coming
// from opts.
func NewConfig(opts ...Opt) *Config {
	conf := &Config{
		Language: lang.NotSet,
		TextOdds: bayes.LabelFreq{
			bayes.Label("Name"):    0.0,
			bayes.Label("NotName"): 0.0,
		},
		BayesOddsThreshold: 100.0,
		// NameDistribution: NameDistribution{
		//   Index: make(map[int]int),
		// },
		Resolver: Resolver{
			URL:       "http://index-api.globalnames.org/api/graphql",
			BatchSize: 500,
			Workers:   runtime.NumCPU(),
		},
	}
	for _, o := range opts {
		err := o(conf)
		Check(err)
	}

	return conf
}

// Opt are options for gnfinder
type Opt func(g *Config) error

// WithLanguage option forces a specific language to be associated with a text.
func WithLanguage(l lang.Language) func(*Config) error {
	return func(conf *Config) error {
		conf.Language = l
		return nil
	}
}

// WithBayes is an option that forces running bayes name-finding even when
// the language is not supported by training sets.
func WithBayes(b bool) func(*Config) error {
	return func(conf *Config) error {
		conf.Bayes = b
		return nil
	}
}

// WithBayesThreshold is an option for name finding, that sets new threshold
// for results from the Bayes name-finding. All the name candidates that have a
// higher threshold will appear in the resulting names output.
func WithBayesThreshold(odds float64) func(*Config) error {
	return func(conf *Config) error {
		conf.BayesOddsThreshold = odds
		return nil
	}
}

// WithResolverURL option sets a new url for name resolution service.
func WithResolverURL(url string) func(*Config) error {
	return func(conf *Config) error {
		conf.URL = url
		return nil
	}
}

// WithResolverBatch option sets the batch size of name-strings to send to the
// resolution service.
func WithResolverBatch(n int) func(*Config) error {
	return func(conf *Config) error {
		conf.BatchSize = n
		return nil
	}
}

// WithResolverWorkers option sets the number of workers to process
// name-resolution jobs.
func WithResolverWorkers(n int) func(*Config) error {
	return func(conf *Config) error {
		conf.Workers = n
		return nil
	}
}
