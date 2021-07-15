package examples_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HelloWorld", func() {
	It("should be always pass", func() {
		err := error(nil)
		Expect(err).ShouldNot(HaveOccurred())

		println("My test case 'hello world' is ok!")
	})
})
