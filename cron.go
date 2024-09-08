package main

import (
	"fmt"
	"strings"
)

type FieldDefinition struct {
	name               string
	begin              int
	end                int
	allowedExpressions []expressionType
}

func NewFieldDefinition(
	name string,
	begin,
	end int,
	allowedExpressions []expressionType) FieldDefinition {
	if end < begin {
		panic(fmt.Errorf(
			"field def'%s' cannot be created because begin=%d is greater than end=%d",
			name,
			begin,
			end,
		))
	}

	if len(allowedExpressions) == 0 {
		panic(fmt.Errorf(
			"field def '%s' cannot be created because no allowed expressions provided",
			name,
		))
	}

	return FieldDefinition{name, begin, end, allowedExpressions}
}

type CronField struct {
	name   string
	values []int
}

func (c CronField) String() string {
	var sb strings.Builder

	for i, val := range c.values {
		if i != 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", val))
	}

	return sb.String()
}

func NewCronFieldFromToken(def FieldDefinition, token string) (CronField, error) {
	for _, expressionType := range def.allowedExpressions {
		parser, ok := parsers[expressionType]
		if !ok {
			return CronField{}, fmt.Errorf("%s is an unsupported expression", expressionType)
		}

		values, err := parser.Parse(token, def.begin, def.end)

		if err == nil {
			return CronField{def.name, values}, nil
		}
	}

	return CronField{}, fmt.Errorf("unable to parse cron field '%s'", def.name)
}

type CronExpression struct {
	fields []CronField
}

func (c CronExpression) WriteExpression() [][]string {
	output := make([][]string, 0, len(c.fields))
	for _, field := range c.fields {
		line := []string{field.name, field.String()}
		output = append(output, line)
	}

	return output
}

func NewCronExpressionFromTokens(tokens []string, defs []FieldDefinition) (CronExpression, error) {
	if len(tokens) != len(defs) {
		err := fmt.Errorf(
			"token count must match def count: got %d expected %d",
			len(tokens),
			len(defs),
		)
		return CronExpression{}, err
	}

	fields := make([]CronField, 0, len(defs))
	for i, def := range defs {
		cronField, err := NewCronFieldFromToken(def, tokens[i])

		if err != nil {
			return CronExpression{}, err
		}

		fields = append(fields, cronField)
	}

	return CronExpression{fields}, nil
}
