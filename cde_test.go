package main

import (
	"reflect"
	"testing"
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

		if len(argv) != 1 {
			t.Errorf("Expected length of 1, Got %d", len(argv))
		}
	}
}

func TestHelpReformattingWithCommand(t *testing.T) {
	t.Parallel()

	tests := []string{"--help", "-h", "help"}
	expected := "test"
	expectedArgv := []string{"test", "--help"}

	for _, test := range tests {
		actual, argv := parseArgs([]string{test, "test"})

		if actual != expected {
			t.Errorf("Expected %s, Got %s", expected, actual)
		}

		if !reflect.DeepEqual(expectedArgv, argv) {
			t.Errorf("Expected %v, Got %v", expectedArgv, argv)
		}
	}
}
