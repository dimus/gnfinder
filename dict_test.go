package gnfinder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dictionary", func() {
	Describe("GreyUninomials", func() {
		It("has grey uninomials list", func() {
			l := len(dictionary.GreyUninomials)
			Expect(l).To(Equal(237))
			_, ok := dictionary.GreyUninomials["Gastropoda"]
			Expect(ok).To(Equal(true))
		})
	})
	Describe("WhiteGenera", func() {
		It("has white genus list", func() {
			l := len(dictionary.WhiteGenera)
			Expect(l).To(Equal(415619))
			_, ok := dictionary.WhiteGenera["Plantago"]
			Expect(ok).To(Equal(true))
		})
	})
})
