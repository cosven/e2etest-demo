package ginkgo_helper

import (
	. "github.com/onsi/gomega"
)

func Try(xs ...interface{}) interface{} {
	if len(xs) == 0 {
		return nil
	}
	if err, ok := xs[len(xs)-1].(error); ok && err != nil {
		Expect(err).ShouldNot(HaveOccurred())
	}
	return xs[0]
}
