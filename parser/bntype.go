package parser

import (
	"errors"
)

type BnCode struct {
	IsInt, IsString, IsList, IsDict bool

	Value interface{}
}

func (obj *BnCode) getInt() (int64, error) {
	val, ok := obj.Value.(int64)
	if !ok {
		return val, errors.New("Given Value is not an int")
	}

	return val, nil
}

func (obj *BnCode) getString() (string, error) {
	val, ok := obj.Value.(string)
	if !ok {
		return val, errors.New("Given Value is not an string")
	}

	return val, nil
}

func (obj *BnCode) getDict() (map[string]BnCode, error) {
	val, ok := obj.Value.(map[string]BnCode)
	if !ok {
		return val, errors.New("Given Value is not a dictionary")
	}

	return val, nil
}

func (obj *BnCode) getList() ([]BnCode, error) {
	val, ok := obj.Value.([]BnCode)
	if !ok {
		return val, errors.New("Given Value is not a dictionary")
	}

	return val, nil
}
