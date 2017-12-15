package gnfinder_test

import (
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("NewGnfinder()", func() {
		It("returns new Config object", func() {
			conf := util.NewConfig()
			Expect(conf.Language).To(Equal(lang.NotSet))
			Expect(conf.Bayes).To(BeFalse())
			Expect(conf.BayesOddsThreshold).To(Equal(100.0))
			Expect(conf.URL).To(Equal("http://index-api.globalnames.org/api/graphql"))
			Expect(conf.BatchSize).To(Equal(500))
		})

		It("takes language", func() {
			conf := util.NewConfig(util.WithLanguage(lang.English))
			Expect(conf.Language).To(Equal(lang.English))
		})

		It("sets bayes", func() {
			conf := util.NewConfig(util.WithBayes(true))
			Expect(conf.Bayes).To(BeTrue())
		})

		It("sets bayes' threshold", func() {
			conf := util.NewConfig(util.WithBayesThreshold(200))
			Expect(conf.BayesOddsThreshold).To(Equal(200.0))
		})

		It("sets a url for resolver", func() {
			conf := util.NewConfig(util.WithResolverURL("http://example.org"))
			Expect(conf.URL).To(Equal("http://example.org"))
		})

		It("sets batch size for resolver", func() {
			conf := util.NewConfig(util.WithResolverBatch(333))
			Expect(conf.BatchSize).To(Equal(333))
		})

		It("sets workers' number  for resolver", func() {
			conf := util.NewConfig(util.WithResolverWorkers(1))
			Expect(conf.Workers).To(Equal(1))
		})

		It("sets several options", func() {
			conf := util.NewConfig(util.WithResolverWorkers(10),
				util.WithBayes(true), util.WithLanguage(lang.English))
			Expect(conf.Workers).To(Equal(10))
			Expect(conf.Language).To(Equal(lang.English))
			Expect(conf.Bayes).To(BeTrue())
		})
	})
})
