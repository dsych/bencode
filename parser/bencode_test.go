package parser

import (
	"bytes"
	"strings"
	"testing"
)

func TestStringParserPositive(t *testing.T) {
	inputStr := ":foo"
	input := []byte(inputStr)

	output, err := parseString(bytes.NewReader(input), '3')

	if err != nil {
		t.Errorf("Unable to parse correct string, %v", err)
	} else if output.State != BnString {
		t.Errorf("Incorrect flags are set")
	}

	v, err := output.getString()

	if err != nil {
		t.Errorf("Unable to parse correct string. %v", err)
	} else if expected := strings.SplitAfter(inputStr, ":")[1]; v != expected {
		t.Errorf("Incorrectly parsed the string, expected %s got %s", expected, v)
	}

}

func TestIntParsePositive(t *testing.T) {
	inputStr := "-123456e"
	input := []byte(inputStr)

	output, err := parseInt(bytes.NewReader(input), 'i')

	if err != nil {
		t.Errorf("Unable to parse correct int, %v", err)
	} else if output.State != BnInt {
		t.Errorf("Incorrect flags are set")
	}

	v, err := output.getInt()
	expected := int(-123456)
	if err != nil {
		t.Errorf("Unable to parse correct int. %v", err)
	} else if v != expected {
		t.Errorf("Incorrectly parsed the string, expected %d got %d", expected, v)
	}
}

func TestParseListPositive(t *testing.T) {
	inputStr := "i42e3:fooe"
	input := []byte(inputStr)

	output, err := parseList(bytes.NewReader(input), 'l')

	if err != nil {
		t.Errorf("Unable to parse correct list, %v", err)
	} else if output.State != BnList {
		t.Errorf("Incorrect flags are set")
	}

	v, err := output.getList()
	if err != nil || len(v) != 2 {
		t.Errorf("Unable to parse correct list. Got %v, with error %v", v, err)
	}

}

func TestDictParsePositive(t *testing.T) {
	input := []byte("5:helloi-3e4:spam3:foo3:zooli42e3:fooee")
	output, err := parseDict(bytes.NewReader(input), 'd')

	if err != nil {
		t.Errorf("Unable to parse correct dict, %v", err)
	} else if output.State != BnDict {
		t.Errorf("Incorrect flags are set")
	}

	v, err := output.getDict()
	if err != nil || len(v) != 3 {
		t.Errorf("Unable to parse correct dict. %v, Error: %v", v, err)
	}
}
