package util

import (
	"flag"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParametrizationUtils", func() {
	Context("when parsing parameters", func() {
		It("with correct input", func() {
			const failed = "test failed"

			var (
				test  string
				test2 string
			)

			ParameterizedTestCase("testcase", func(flagSet *flag.FlagSet) {
				flagSet.StringVar(&test, "test", failed, "string test.")
				flagSet.StringVar(&test2, "test2", failed, "string test2.")
			})

			fmt.Println(test)
			fmt.Println(test2)
			Expect(test).ShouldNot(Equal(failed))
			Expect(test).ShouldNot(Equal(failed))
		})
	})
})
