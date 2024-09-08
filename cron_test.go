package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assertOutput(expected, actual [][]string) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d rows but got %d", len(expected), len(actual))
	}

	for i := range actual {
		if !reflect.DeepEqual(expected[i], actual[i]) {
			return fmt.Errorf("expected line %v actual %v", expected[i], actual[i])
		}
	}

	return nil
}

func TestValidExpression(t *testing.T) {
	arg := "*/15 0 1,15 * 1-3 /usr/bin/find"

	actual, err := execute(arg)

	if err != nil {
		t.Errorf("Expected nil error but got %v", err)
		return
	}

	expected := [][]string{
		{"minute", "0 15 30 45"},
		{"hour", "0"},
		{"day of month", "1 15"},
		{"month", "1 2 3 4 5 6 7 8 9 10 11 12"},
		{"day of week", "1 2 3"},
		{"command", "/usr/bin/find"},
	}

	err = assertOutput(expected, actual)
	if err != nil {
		t.Error(err)
	}
}

func TestErrorWhenExtraToken(t *testing.T) {
	arg := "*/15 0 1,15 * 1-3 5-8 /usr/bin/find"

	_, err := execute(arg)

	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
}

func TestErrorWhenMissingToken(t *testing.T) {
	arg := "*/15 0 1,15 * /usr/bin/find"

	_, err := execute(arg)

	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
}

func TestErrorWhenInvalidRange(t *testing.T) {
	arg := "*/15 0 7,3 * /usr/bin/find"

	_, err := execute(arg)

	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
}

func TestErrorWhenOutOfBounds(t *testing.T) {
	arg := "*/15 100 1,15 * /usr/bin/find"

	_, err := execute(arg)

	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
}

func TestErrorRangeOutOfBounds(t *testing.T) {
	arg := "*/15 100 0,6 * /usr/bin/find"

	_, err := execute(arg)

	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
}
