package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func SplitArgByWhitespace(arg string, expectedTokensCount int) ([]string, error) {
	tokenScanner := bufio.NewScanner(strings.NewReader(arg))
	tokenScanner.Split(bufio.ScanWords)

	tokens := make([]string, 0, expectedTokensCount)

	for i := 0; tokenScanner.Scan(); i++ {
		if len(tokens) < expectedTokensCount {
			tokens = append(tokens, tokenScanner.Text())
		} else {
			return nil, fmt.Errorf(
				"invalid argument string: expected %d components separated by whitespace but got more",
				expectedTokensCount,
			)
		}
	}

	if len(tokens) < expectedTokensCount {
		return nil, fmt.Errorf(
			"invalid argument string: expected %d components separated by whitespace but got less",
			expectedTokensCount,
		)
	}

	return tokens, nil
}

func execute(arg string) ([][]string, error) {
	defs := []FieldDefinition{
		NewFieldDefinition(
			"minute",
			0,
			59,
			[]expressionType{
				WILDCARD,
				RANGE,
				STEP,
				VALUE,
				ENUMERATION,
			},
		),
		NewFieldDefinition(
			"hour",
			0,
			23,
			[]expressionType{
				WILDCARD,
				RANGE,
				STEP,
				VALUE,
				ENUMERATION,
			},
		),
		NewFieldDefinition(
			"day of month",
			1,
			31,
			[]expressionType{
				WILDCARD,
				RANGE,
				STEP,
				VALUE,
				ENUMERATION,
			},
		),
		NewFieldDefinition(
			"month",
			1,
			12,
			[]expressionType{
				WILDCARD,
				RANGE,
				STEP,
				VALUE,
				ENUMERATION,
			},
		),
		NewFieldDefinition(
			"day of week",
			1,
			7,
			[]expressionType{
				WILDCARD,
				RANGE,
				STEP,
				VALUE,
				ENUMERATION,
			},
		),
	}
	var expectedTokensCount = len(defs) + 1

	tokens, err := SplitArgByWhitespace(arg, expectedTokensCount)
	if err != nil {
		return nil, err
	}

	cronEx, err := NewCronExpressionFromTokens(tokens[:len(defs)], defs)
	if err != nil {
		return nil, err
	}

	output := cronEx.WriteExpression()

	command := tokens[len(tokens)-1]
	output = append(output, []string{"command", command})

	return output, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(
			"error: must provide an argument string (e.g. \"*/15 0 1,15 * 1-3 /usr/bin/find\")",
		)
		return
	}

	output, err := execute(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, line := range output {
		fmt.Fprintf(w, "%s\t%s\n", line[0], line[1])
	}
	w.Flush()
}
