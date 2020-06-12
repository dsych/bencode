package parser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStringFlattenPositive(t *testing.T) {
	targetVal := "foo"
	obj := BnCode{State: BnString, Value: targetVal}
	writer := bytes.NewBuffer([]byte{})

	err := flattenString(obj, writer)
	if err != nil {
		t.Errorf("Unable to flatten a correct string. %v", err)
	}

	res := fmt.Sprintf("%d:%s", len(targetVal), targetVal)
	s := writer.String()
	if err != nil {
		t.Errorf("Unable to flatten a correct string. %v", err)
	} else if s != res {
		t.Errorf("Incorrectly flattened the string, expected %s got %s", res, s)
	}

}

func TestIntFlattenPositive(t *testing.T) {
	targetVal := -42
	obj := BnCode{State: BnInt, Value: targetVal}
	writer := bytes.NewBuffer([]byte{})

	err := flattenInt(obj, writer)
	if err != nil {
		t.Errorf("Unable to flatten a correct int. %v", err)
	}

	res := fmt.Sprintf("i%de", targetVal)
	s := writer.String()
	if err != nil {
		t.Errorf("Unable to flatten a correct int. %v", err)
	} else if s != res {
		t.Errorf("Incorrectly flattened the int, expected %s got %s", res, s)
	}
}

func TestListFlattenPositive(t *testing.T) {
	obj := BnCode{State: BnList, Value: []BnCode{{State: BnInt, Value: 4}, {State: BnString, Value: "abc"}}}
	writer := bytes.NewBuffer([]byte{})

	err := flattenList(obj, writer)
	if err != nil {
		t.Errorf("Unable to flatten a correct list. %v", err)
	}

	res := "li4e3:abce"
	s := writer.String()
	if err != nil {
		t.Errorf("Unable to flatten a correct list. %v", err)
	} else if s != res {
		t.Errorf("Incorrectly flattened the list, expected %s got %s", res, s)
	}
}

func TestDictFlattenPositive(t *testing.T) {
	obj := BnCode{State: BnDict, Value: map[string]BnCode{"a": {State: BnInt, Value: 4}, "b": {State: BnString, Value: "abc"}, "list": {State: BnList, Value: []BnCode{{State: BnInt, Value: 5}}}}}

	writer := bytes.NewBuffer([]byte{})

	err := flattenDict(obj, writer)
	if err != nil {
		t.Errorf("Unable to flatten a correct dict. %v", err)
	}

	res := "d1:ai4e1:b3:abc4:listli5eee"
	s := writer.String()
	if err != nil {
		t.Errorf("Unable to flatten a correct dict. %v", err)
	} else if s != res {
		t.Errorf("Incorrectly flattened the dict, expected %s got %s", res, s)
	}
}
