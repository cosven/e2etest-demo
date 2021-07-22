package matcher

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("test all collection assertions", func() {
	var allPassedMessage string
	dualAssertion := BeNumerically(">=", 3)

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
			assertion := AllPassed(dualAssertion)
			r, err := assertion.Match(array)
			failureMessage := assertion.FailureMessage(array)

			fmt.Println(failureMessage)
			allPassedMessage = failureMessage

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeFalse())
			Expect(failureMessage).ShouldNot(BeEmpty())
		})
		It("should failed when chained with ShouldNot", func() { // AnyFailed
			assertion := AllPassed(BeNumerically(">=", 1))
			r, err := assertion.Match(array)
			r = !r // such negation is done by the driver of Match, as in gomega@v1.10.1/internal/assertion/assertion.go:63
			failureMessage := assertion.NegatedFailureMessage(array)

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
		It("should as if AllPassed were used when chained with ShouldNot", func() { // AllPassed
			assertion := AnyFailed(dualAssertion)
			r, err := assertion.Match(array)
			r = !r // such negation is done by the driver of Match, as in gomega@v1.10.1/internal/assertion/assertion.go:63P
			failureMessage := assertion.NegatedFailureMessage(array)

			fmt.Println(failureMessage)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(BeFalse())
			Expect(AllButLastLine(failureMessage)).Should(Equal(allPassedMessage))
		})

		It("should panic when not call Match first", func() {
			Expect(func() {
				AnyFailed(BeNumerically(">=", 3)).FailureMessage(array)
			}).Should(PanicWithOutput("on AnyFailed, Match should be called first"))
		})
	})
})

func AllButLastLine(str string) string {
	split := strings.Split(str, "\n")
	return strings.Join(split[:len(split)-1], "\n")
}
