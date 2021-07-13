package oom_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oom Suite")
}
