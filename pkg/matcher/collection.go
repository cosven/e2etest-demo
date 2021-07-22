package matcher

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"github.com/pingcap/errors"
)

type collectionMatched struct {
	successElements  []reflect.Value
	successPositions []int
	failedElements   []reflect.Value
	failedPositions  []int
}

func match(assertion types.GomegaMatcher, actual interface{}) (*collectionMatched, *reflect.Value, error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return nil, nil, errors.New(fmt.Sprintf("unsupported type for AllPassed: %T, only type Array and Slice are supportted.", t))
	}

	v := reflect.ValueOf(actual)
	var result collectionMatched

	for i := 0; i < v.Len(); i++ {
		eRaw := v.Index(i)
		e := eRaw.Interface()

		r, err := assertion.Match(e)
		if err != nil {
			return nil, nil, err
		}

		if r {
			result.successElements = append(result.successElements, eRaw)
			result.successPositions = append(result.successPositions, i)
		} else {
			result.failedElements = append(result.failedElements, eRaw)
			result.failedPositions = append(result.failedPositions, i)
		}
	}

	return &result, &v, nil
}

func (matched *collectionMatched) expectedAllPassedFailure(assertion types.GomegaMatcher, actual interface{}) string {
	var diagnose string
	for i, eRaw := range matched.failedElements {
		e := eRaw.Interface()
		diagnose += "\n"
		diagnose += fmt.Sprintf("At position %d: \n", matched.failedPositions[i])
		diagnose += format.IndentString(assertion.FailureMessage(e), 1)
	}

	return fmt.Sprintf("Expected ALL element of \n%s\nmatches the assertion, yet:%s", format.Object(actual, 1), format.IndentString(diagnose, 1))
}

func (matched *collectionMatched) expectedAnyFailedFailure(assertion types.GomegaMatcher, actual interface{}) string {
	var diagnose string
	for i, eRaw := range matched.successElements {
		e := eRaw.Interface()
		diagnose += "\n"
		diagnose += fmt.Sprintf("At position %d: \n", matched.successPositions[i])
		diagnose += format.IndentString(assertion.NegatedFailureMessage(e), 1)
	}

	return fmt.Sprintf("Expected ANY element of \n%s\nNOT matches the assertion, yet at these indexes, elements matches:%s", format.Object(actual, 1), format.IndentString(diagnose, 1))
}

type allPassedMatcher struct {
	assertion types.GomegaMatcher
	evaluated bool

	matched *collectionMatched
}

// AllPassed asserts that the assertion passed in should pass for ALL elements of the array or slice.
// if not, it will print the failed element and its position as helpful diagnose message.
func AllPassed(assertion types.GomegaMatcher) *allPassedMatcher {
	return &allPassedMatcher{
		evaluated: false,
		assertion: assertion,
	}
}

func (matcher *allPassedMatcher) Match(actual interface{}) (success bool, err error) {
	matched, v, err := match(matcher.assertion, actual)
	if err != nil {
		return false, err
	}

	matcher.evaluated = true
	matcher.matched = matched

	return len(matched.successElements) == v.Len(), err
}

func (matcher *allPassedMatcher) FailureMessage(actual interface{}) (message string) {
	if !matcher.evaluated {
		panic("you should call Match to execute matching first.")
	}
	return matcher.matched.expectedAllPassedFailure(matcher.assertion, actual)
}

func (matcher *allPassedMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	if !matcher.evaluated {
		panic("you should call Match to execute matching on first.")
	}
	return matcher.matched.expectedAnyFailedFailure(matcher.assertion, actual) +
		"\n [ ShouldNot + AllPassed == Should + AnyFailed, ( \\not (\\exists x. A(x)) <=> \\forall x. !A(x) ) ]"
}

type anyFailedMatcher struct {
	assertion types.GomegaMatcher
	evaluated bool

	matched *collectionMatched
}

// AnyFailed asserts that the assertion passed in should fail for AT LEAST ONE element of the array or slice.
// if not, it will print all passed elements and its position as helpful diagnose message.
func AnyFailed(assertion types.GomegaMatcher) *anyFailedMatcher {
	return &anyFailedMatcher{
		evaluated: false,
		assertion: assertion,
	}
}

func (matcher *anyFailedMatcher) Match(actual interface{}) (success bool, err error) {
	matched, _, err := match(matcher.assertion, actual)
	if err != nil {
		return false, err
	}

	matcher.evaluated = true
	matcher.matched = matched

	return len(matched.failedElements) > 0, err
}

func (matcher *anyFailedMatcher) FailureMessage(actual interface{}) (message string) {
	if !matcher.evaluated {
		panic("you should call Match to execute matching first.")
	}
	return matcher.matched.expectedAnyFailedFailure(matcher.assertion, actual)
}

func (matcher *anyFailedMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	if !matcher.evaluated {
		panic("you should call Match to execute matching first.")
	}
	return matcher.matched.expectedAllPassedFailure(matcher.assertion, actual) +
		"\n [ ShouldNot + AnyFailed == Should + AllPassed ( \\not (\\exists x. \\not A(x)) <=> \\not (\\not forall x. A(x)) <=> \\forall x. F(x) ) ]"
}
