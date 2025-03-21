package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	type Input struct {
		input    string
		expected []string
	}

	var tests []Input
	tests = append(tests, Input{input: "  Hello  world  ", expected: []string{"Hello", "world"}})

	for _, test := range tests {
		actual := cleanInput(test.input)
		if len(actual) != len(test.expected) {
			t.Errorf("\nLengths not equal. Expected: %v, Actual: %v", test.expected, actual)
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := test.expected[i]
			if word != expectedWord {
				t.Errorf("\nExpected: %v\nActual: %v", expectedWord, word)
			}
		}
	}

}
