package ginkgo_helper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinkgoHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinkgoHelper Suite")
}
