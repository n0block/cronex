package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type expressionType string

const (
	WILDCARD    expressionType = "wildcard"
	RANGE       expressionType = "range"
	STEP        expressionType = "step"
	VALUE       expressionType = "value"
	ENUMERATION expressionType = "enumeration"
)

var parsers = map[expressionType]Parser{
	WILDCARD:    WildcardExpressionParser{},
	RANGE:       NewRangeExpressionParser(),
	STEP:        NewStepExpressionParser(),
	VALUE:       NewValueExpressionParser(),
	ENUMERATION: EnumerationExpressionParser{},
}

type Parser interface {
	Parse(string, int, int) ([]int, error)
}

type RangeExpressionParser struct {
	pattern *regexp.Regexp
}

func NewRangeExpressionParser() RangeExpressionParser {
	return RangeExpressionParser{regexp.MustCompile(`^(\d+)-(\d+)$`)}
}

func (rp RangeExpressionParser) Parse(token string, begin, end int) ([]int, error) {
	match := rp.pattern.FindStringSubmatch(token)

	if len(match) == 0 {
		return nil, fmt.Errorf("unable to parse: %s", token)
	}

	providedBegin, err := strconv.Atoi(match[1])
	if err != nil || providedBegin < begin || end < providedBegin {
		return nil, fmt.Errorf("%d is not within [%d, %d]", providedBegin, begin, end)
	}

	providedEnd, err := strconv.Atoi(match[2])
	if err != nil || end < providedEnd {
		return nil, fmt.Errorf("%d is not within [%d, %d]", providedEnd, begin, end)
	}

	if providedEnd < providedBegin {
		return nil, fmt.Errorf("begin=%d is greater than end=%d", providedEnd, providedBegin)
	}

	return generateValues(providedBegin, providedEnd, 1), nil
}

type StepExpressionParser struct {
	pattern *regexp.Regexp
}

func NewStepExpressionParser() StepExpressionParser {
	return StepExpressionParser{regexp.MustCompile(`^\*/(\d+)$`)}
}

func (se StepExpressionParser) Parse(token string, begin, end int) ([]int, error) {
	match := se.pattern.FindStringSubmatch(token)

	if len(match) == 0 {
		return nil, fmt.Errorf("unable to parse: %s", token)
	}

	step, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, err
	}

	return generateValues(begin, end, step), nil
}

type ValueExpressionParser struct {
	pattern *regexp.Regexp
}

func NewValueExpressionParser() ValueExpressionParser {
	return ValueExpressionParser{regexp.MustCompile(`^(\d+)$`)}
}

func (v ValueExpressionParser) Parse(token string, begin, end int) ([]int, error) {
	match := v.pattern.FindStringSubmatch(token)

	if len(match) == 0 {
		return nil, fmt.Errorf("unable to parse: %s", token)
	}

	value, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, err
	}

	return generateValues(value, value, 1), nil
}

type WildcardExpressionParser struct {
}

func (w WildcardExpressionParser) Parse(token string, begin, end int) ([]int, error) {
	if token != "*" {
		return nil, fmt.Errorf("unable to parse: %s", token)
	}

	return generateValues(begin, end, 1), nil
}

type EnumerationExpressionParser struct {
}

func (e EnumerationExpressionParser) Parse(token string, begin, end int) ([]int, error) {
	strValues := strings.Split(token, ",")

	if len(strValues) < 2 {
		return nil, fmt.Errorf("unable to parse: %s", token)
	}

	values := make([]int, 0, len(strValues))
	for _, s := range strValues {
		value, err := strconv.Atoi(s)
		if err != nil || value < begin || end < value {
			return nil, fmt.Errorf("unable to parse: %s", token)
		}
		values = append(values, value)
	}

	return values, nil
}
