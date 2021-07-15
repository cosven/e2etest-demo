# E2E Test Framework

## Quick Start

### Add a new test case

1. Create a new directory to save your test cases.

```sh
cd testcase
mkdir -p examples/
cd examples/
```

2. Create a new ginkgo test suite.

```sh
# Install ginkgo first.
go get github.com/onsi/ginkgo/ginkgo

# Create a test suite.
ginkgo bootstrap
```

3. Run the test suite.

```
$ ginkgo

Running Suite: Examples Suite
=============================
Random Seed: 1626320026
Will run 0 of 0 specs


Ran 0 of 0 Specs in 0.001 seconds
SUCCESS! -- 0 Passed | 0 Failed | 0 Pending | 0 Skipped
PASS

Ginkgo ran 1 suite in 1.086402735s
Test Suite Passed
```

4. Write your own test case and run it.

```sh
ginkgo generate hello_world
```

```golang
$ cat hello_world_test.go
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
```

Try to run the test case.

```
$ ginkgo
Running Suite: Examples Suite
=============================
Random Seed: 1626320700
Will run 1 of 1 specs

•My test case 'hello world' is ok!

Ran 1 of 1 Specs in 0.000 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped
PASS

Ginkgo ran 1 suite in 1.174906876s
Test Suite Passed
```

Actually, this E2E test framework use ginkgo to manage all the test cases. 
Check [ginkgo docs](https://onsi.github.io/ginkgo/) for more details.

When you want to write a test case which use [infra sdk](https://github.com/pingcap/test-infra), 
you can check the example `testcase/oom/tikv_test.go`.