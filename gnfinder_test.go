package gnfinder_test

import (
	. "github.com/gnames/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gnfinder", func() {
	Describe("FindNames", func() {
		It("finds names", func() {
			s := "Plantago major and Pardosa moesta are spiders and plants"
			output := FindNames([]rune(s), dictionary)
			Expect(output.Names[0].Name).To(Equal("Plantago major"))
			Expect(len(output.Names)).To(Equal(2))
		})

		It("works with very short/empty texts", func() {
			s := "  \n\t    \v\r\n"
			output := FindNames([]rune(s), dictionary)
			Expect(len(output.Names)).To(Equal(0))
			s = "Pomatomus"
			output = FindNames([]rune(s), dictionary)
			Expect(len(output.Names)).To(Equal(1))
			s = "Pomatomus saltator"
			output = FindNames([]rune(s), dictionary)
			Expect(len(output.Names)).To(Equal(1))
		})

		It("recognizes subgenus", func() {
			s := "Pomatomus (Pomatomus) saltator"
			output := FindNames([]rune(s), dictionary)
			Expect(len(output.Names)).To(Equal(2))
			Expect(output.Names[1].Name).To(Equal("Pomatomus"))
		})

		It("finds names in a book", func() {
			output := FindNames([]rune(string(book)), dictionary)
			Expect(len(output.Names)).To(Equal(4455))
		})

		// 	It("finds names in a book with new BayesOddsThreshold", func() {
		// 		output := FindNames([]rune(string(book)),
		// 			dictionary, WithBayesThreshold(1))
		// 		Expect(len(output.Names)).To(Equal(5049))
		// 	})

		It("recognizes 'impossible', unknown and abbreviated binomials", func() {
			s := [][2]string{
				{"{Pardosa) moesta", "Pardosa"},
				{"Pardosa Moesta", "Pardosa"},
				{"\"Pomatomus, saltator", "Pomatomus"},
				{"Pomatomus 'saltator'", "Pomatomus"},
				{"{P. moesta.", "P. moesta"},
				{"Po. saltator", "Po. saltator"},
				{"Pom. saltator", "Pom. saltator"},
				{"SsssAAAbbb saltator!", "Ssssaaabbb saltator"},
				{"ZZZ saltator!", "Zzz saltator"},
				{"One possible Pomatomus saltator...", "Pomatomus saltator"},
				{"[Different Pomatomus ]saltator...", "Pomatomus"},
			}
			for _, v := range s {
				output := FindNames([]rune(v[0]), dictionary)
				Expect(len(output.Names)).To(Equal(1))
				Expect(output.Names[0].Name).To(Equal(v[1]))
			}
		})
	})

	It("rejects black dictionary genera", func() {
		s := []string{"The moesta", "This saltator"}
		for _, v := range s {
			output := FindNames([]rune(v), dictionary)
			Expect(len(output.Names)).To(Equal(0))
		}
	})

	It("does not recognize one letter genera", func() {
		output := FindNames([]rune("I saltator"), dictionary)
		Expect(len(output.Names)).To(Equal(0))
	})

	Describe("FindNamesJSON()", func() {
		It("finds names and returns json representation", func() {
			s := "Plantago major and Pardosa moesta are spiders and plants"
			output := FindNamesJSON([]byte(s), dictionary)
			Expect(string(output)[0:17]).To(Equal("{\n  \"metadata\": {"))
		})
	})
})
