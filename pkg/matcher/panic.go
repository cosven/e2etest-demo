package matcher

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type panicWithOutput struct {
	name     []string
	delegate *matchers.PanicMatcher
}

// PanicWithOutput asserts that the actual function should panicked. this is implemented as a pure delegate of gomega.Panic.
// the only difference is that this matcher will not only asserts the panic, but also print the carried value to stdout.
func PanicWithOutput(name ...string) types.GomegaMatcher {
	return &panicWithOutput{
		name:     name,
		delegate: &matchers.PanicMatcher{},
	}
}

func (p *panicWithOutput) Match(actual interface{}) (success bool, err error) {
	delegated1, delegated2 := p.delegate.Match(actual)

	v := reflect.TypeOf(*p.delegate)
	y, found := v.FieldByName("object")
	if !found {
		panic("object not found by reflect in <PanicMatcher> panicWithOutput.delegate")
	}

	ptrToDelegate := unsafe.Pointer(p.delegate)
	ptrToObject := unsafe.Pointer(uintptr(ptrToDelegate) + y.Offset)
	object := (*interface{})(ptrToObject)

	fmt.Printf("[%s] panicked with: \n%s\n", strings.Join(p.name, "; "), format.IndentString(fmt.Sprintf("%s", *object), 1))

	return delegated1, delegated2
}

func (p *panicWithOutput) FailureMessage(actual interface{}) (message string) {
	return p.delegate.FailureMessage(actual)
}

func (p *panicWithOutput) NegatedFailureMessage(actual interface{}) (message string) {
	return p.delegate.NegatedFailureMessage(actual)
}
