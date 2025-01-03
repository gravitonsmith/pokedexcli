package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "   hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "BuLbAsUaR CHARMANDER Squirtle",
			expected: []string{"bulbasuar", "charmander", "squirtle"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("got %v, want %v", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("got %v, want %v", word, expectedWord)
			}
		}
	}
}
