package oom_test

import (
	"testing"

	"github.com/joho/godotenv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	infra_core "github.com/pingcap/test-infra/sdk/core"
)

var (
	suiteTestCtx infra_core.TestContext
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Load dotenv failed.")
	}
	suiteTestCtx, err = infra_core.BuildContext()
	if err != nil {
		panic("Init test context failed.")
	}
}

func TestOom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oom Suite")
}
