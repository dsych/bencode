package parser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStringFlattenPositive(t *testing.T) {
	targetVal := "foo"
	obj := BnCode{IsString: true, Value: targetVal}
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
	obj := BnCode{IsInt: true, Value: targetVal}
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
	obj := BnCode{IsList: true, Value: []BnCode{{IsInt: true, Value: 4}, {IsString: true, Value: "abc"}}}
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
	obj := BnCode{IsDict: true, Value: map[string]BnCode{"a": {IsInt: true, Value: 4}, "b": {IsString: true, Value: "abc"}, "list": {IsList: true, Value: []BnCode{{IsInt: true, Value: 5}}}}}

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
