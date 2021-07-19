package oom_test

import (
	"errors"
	"os"
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
	err := error(nil)
	if err = godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic("Load dotenv failed.")
		}
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
