package ginkgo_helper

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("test all collection assertions", func() {
	Context("AllPassed", func() {
		array := []int{1, 2, 4, 6}
		slice := array[:]

		It("should passed when used with array, as all elements passed", func() {
			r, err := AllPassed(BeNumerically(">=", 1)).Match(array)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeTrue())
		})
		It("should passed when used with slice, as all elements passed", func() {
			r, err := AllPassed(BeNumerically(">=", 1)).Match(slice)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeTrue())
		})
		It("should failed when any element failed", func() {
			assertion := AllPassed(BeNumerically(">=", 3))
			r, err := assertion.Match(array)
			failureMessage := assertion.FailureMessage(array)

			fmt.Println(failureMessage)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeFalse())
			Expect(failureMessage).ShouldNot(BeEmpty())
		})

		It("should panic when not call Match first", func() {
			Expect(func() {
				AnyFailed(BeNumerically(">=", 3)).FailureMessage(array)
			}).Should(PanicWithOutput("on AllPassed, Match should be called first"))
		})
		It("should panic when chained with ShouldNot", func() {
			Expect(func() {
				AllPassed(BeNumerically(">=", 3)).NegatedFailureMessage([]int{1, 2, 4, 6})
			}).Should(PanicWithOutput("AllPassed chained with ShouldNot"))
		})
	})

	Context("AnyFailed", func() {
		array := []int{1, 2, 4, 6}
		slice := array[:]

		It("should passed when used with array, and some element failed", func() {
			r, err := AnyFailed(BeNumerically(">=", 3)).Match(array)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeTrue())
		})
		It("should passed when used with slice, and some element failed", func() {
			r, err := AnyFailed(BeNumerically(">=", 3)).Match(slice)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeTrue())
		})
		It("should failed when all elements passed", func() {
			assertion := AnyFailed(BeNumerically(">=", 1))
			r, err := assertion.Match(array)
			failureMessage := assertion.FailureMessage(array)

			fmt.Println(failureMessage)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeFalse())
			Expect(failureMessage).ShouldNot(BeEmpty())
		})

		It("should panic when not call Match first", func() {
			Expect(func() {
				AnyFailed(BeNumerically(">=", 3)).FailureMessage(array)
			}).Should(PanicWithOutput("on AnyFailed, Match should be called first"))
		})
		It("should panic when chained with ShouldNot", func() {
			Expect(func() {
				AnyFailed(BeNumerically(">=", 3)).NegatedFailureMessage(array)
			}).Should(PanicWithOutput("AnyFailed chained with ShouldNot"))
		})
	})
})
