package ginkgo_helper

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"github.com/pingcap/errors"
)

type allPassedMatcher struct {
	assertion types.GomegaMatcher
	evaluated bool

	actual          interface{}
	failedElements  []reflect.Value
	failedPositions []int
}

func AllPassed(assertion types.GomegaMatcher) *allPassedMatcher {
	return &allPassedMatcher{
		evaluated: false,
		assertion: assertion,
	}
}

func (matcher *allPassedMatcher) Match(actual interface{}) (success bool, err error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return false, errors.New(fmt.Sprintf("unsupported type for AllPassed: %T, only type Array and Slice are supportted.", t))
	}

	matched, err := match(true, matcher.assertion, actual)
	if err != nil {
		return false, err
	}

	matcher.actual = actual
	matcher.failedElements = matched.undesiredElements
	matcher.failedPositions = matched.undesiredPositions
	matcher.evaluated = true

	return !matched.undesired, err
}

type collectionMatched struct {
	undesired          bool
	undesiredElements  []reflect.Value
	undesiredPositions []int
}

func match(desired bool, assertion types.GomegaMatcher, actual interface{}) (*collectionMatched, error) {
	v := reflect.ValueOf(actual)
	var result collectionMatched
	result.undesired = false

	for i := 0; i < v.Len(); i++ {
		eRaw := v.Index(i)
		e := eRaw.Interface()

		r, err := assertion.Match(e)
		if err != nil {
			return nil, err
		}

		if r != desired {
			result.undesiredElements = append(result.undesiredElements, eRaw)
			result.undesiredPositions = append(result.undesiredPositions, i)
			result.undesired = true
		}
	}

	return &result, nil
}

func (matcher *allPassedMatcher) FailureMessage(actual interface{}) (message string) {
	if !matcher.evaluated {
		panic("you should call Match to execute matching first.")
	}

	var diagnose string
	for i, eRaw := range matcher.failedElements {
		e := eRaw.Interface()
		diagnose += "\n"
		diagnose += fmt.Sprintf("At position %d: \n", matcher.failedPositions[i])
		diagnose += format.IndentString(matcher.assertion.FailureMessage(e), 1)
	}

	return fmt.Sprintf("Expected ALL element of \n%s\nmatches the assertion, yet:%s", format.Object(actual, 1), format.IndentString(diagnose, 1))
}

func (matcher *allPassedMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	panic(`AllPassed only supports Should. Use AnyFailed if you want ShouldNot + AllPassed. (ShouldNot + AllPassed == Should + AnyFailed, \not (\exists x. A(x)) <=> \forall x. !A(x) )`)
}

type anyFailedMatcher struct {
	assertion types.GomegaMatcher
	evaluated bool

	actual           interface{}
	successElements  []reflect.Value
	successPositions []int
}

func AnyFailed(assertion types.GomegaMatcher) *anyFailedMatcher {
	return &anyFailedMatcher{
		evaluated: false,
		assertion: assertion,
	}
}

func (matcher *anyFailedMatcher) Match(actual interface{}) (success bool, err error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return false, errors.New(fmt.Sprintf("unsupported type for AllPassed: %T, only type Array and Slice are supportted.", t))
	}

	matched, err := match(false, matcher.assertion, actual)
	if err != nil {
		return false, err
	}

	matcher.actual = actual
	matcher.successElements = matched.undesiredElements
	matcher.successPositions = matched.undesiredPositions
	matcher.evaluated = true

	return len(matcher.successElements) != reflect.ValueOf(actual).Len(), err // TODO: deep comparing?
}

func (matcher *anyFailedMatcher) FailureMessage(actual interface{}) (message string) {
	if !matcher.evaluated {
		panic("you should call Match to execute matching on first.")
	}

	var diagnose string
	for i, eRaw := range matcher.successElements {
		e := eRaw.Interface()
		diagnose += "\n"
		diagnose += fmt.Sprintf("At position %d: \n", matcher.successPositions[i])
		diagnose += format.IndentString(matcher.assertion.NegatedFailureMessage(e), 1)
	}

	return fmt.Sprintf("Expected ANY element of \n%s\nNOT matches the assertion, yet at these indexes, elements matches:%s", format.Object(actual, 1), format.IndentString(diagnose, 1))
}

func (matcher *anyFailedMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	panic(`AnyFailed only supports Should. Use AllPassed if you want to use ShouldNot + AnyFailed. (ShouldNot + AnyFailed == Should + AllPassed, \not (\exists x. \not A(x)) <=> \not (\not forall x. A(x)) <=> \forall x. F(x) )`)
}
