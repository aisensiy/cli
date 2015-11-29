package main

import (
	"testing"
	"reflect"
)

func TestHelpReformatting(t *testing.T) {
	t.Parallel()

	tests := []string{"--help", "-h", "help"}
	expected := "help"

	for _, test := range tests {
		actual, argv := parseArgs([]string{test})

		if actual != expected {
			t.Errorf("Expected %s, Got %s", expected, actual)
		}

		if len(argv) != 0 {
			t.Errorf("Expected length of 0, Got %d", len(argv))
		}
	}
}

func TestHelpReformattingWithCommand(t *testing.T) {
	t.Parallel()

	tests := []string{"--help", "-h", "help"}
	expected := "test"

	for _, test := range tests {
		actual, argv := parseArgs([]string{"test", test})

		if actual != expected {
			t.Errorf("Expected %s, Got %s", expected, actual)
		}


		if !reflect.DeepEqual([]string{test}, argv) {
			t.Errorf("Expected %v, Got %v", []string{test}, argv)
		}
	}
}