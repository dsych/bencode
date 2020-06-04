package parser

import (
	"bytes"
	"strings"
	"testing"
)

func TestStringParserPositive(t *testing.T) {
	inputStr := "3:foo"
	input := []byte(inputStr)

	output, err := parseString(bytes.NewReader(input))

	if err != nil {
		t.Errorf("Unable to parse correct string, %v", err)
	} else if !output.IsString || (output.IsInt || output.IsDict || output.IsList) {
		t.Errorf("Incorrect flags are set")
	}

	v, err := output.getString()

	if err != nil {
		t.Errorf("Unable to parse correct string. %v", err)
	} else if v != strings.SplitAfter(inputStr, ":")[1] {
		t.Errorf("Incorrectly parsed the string, expected %s got %s", inputStr, v)
	}

}
