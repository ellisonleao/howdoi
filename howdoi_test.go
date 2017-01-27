package main

import "testing"

func TestSanitizeQuestion(t *testing.T) {
	h := new(Howdoi)

	// case 1 - empty input
	badInput := []string{}
	err := h.sanitizeQuestion(badInput)
	if err == nil {
		t.Error("failed on sanitize empty input")
	}

	// case 2 - only spaces on input
	badInput = []string{"", " "}
	err = h.sanitizeQuestion(badInput)
	if err == nil {
		t.Error("failed on sanitize empty input")
	}

	// case 3 valid input
	input := []string{"open", "file", "in", "python"}
	err = h.sanitizeQuestion(input)
	if err != nil {
		t.Error("returning error on valid input")
	}

	if h.Question != "open+file+in+python" {
		t.Error("Wrong return string on sanitizeQuestion input")
	}
}
