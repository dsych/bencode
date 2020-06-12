package parser

import (
	"errors"
)

const (
	// BnInt enum that indicates the state of the Value of BnCode
	BnInt = iota
	// BnString enum that indicates the state of the Value of BnCode
	BnString = iota
	// BnList enum that indicates the state of the Value of BnCode
	BnList = iota
	// BnDict enum that indicates the state of the Value of BnCode
	BnDict = iota
)

// BnCode is structure that wraps main Bencode types :
//
// 1. int - constant BnInt
//
// 2. string - constant BnString
//
// 3. list - constatnt BnList
//
// 4. dictionary - constant BnDict
//
// Each of the types have corresponding code that will show the current Value state,
// hence you can only call geInt method on the BnCode, which State is set to BnInt.
type BnCode struct {
	State int
	Value interface{}
}

// GetInt tries converting Value to int.
//
// Returns error if unable to cast to int or State is not BnInt.
func (obj *BnCode) GetInt() (int, error) {
	if val, ok := obj.Value.(int); obj.State != BnInt || !ok {
		return val, errors.New("Given Value is not an int")
	} else {
		return val, nil
	}
}

// GetString tries converting Value to string.
//
// Returns error if unable to cast to string or State is not BnString.
func (obj *BnCode) GetString() (string, error) {
	if val, ok := obj.Value.(string); obj.State != BnString || !ok {
		return val, errors.New("Given Value is not an string")
	} else {
		return val, nil
	}
}

// GetDict tries converting Value to dictionary
//
// Returns error if unable to cast to dictionary or State is not BnDict
func (obj *BnCode) GetDict() (map[string]BnCode, error) {
	if val, ok := obj.Value.(map[string]BnCode); obj.State != BnDict || !ok {
		return val, errors.New("Given Value is not a dictionary")
	} else {
		return val, nil
	}
}

// GetList tries converting Value to list
//
// Returns error if unable to cast to list or State is not BnList
func (obj *BnCode) GetList() ([]BnCode, error) {

	if val, ok := obj.Value.([]BnCode); obj.State != BnList || !ok {
		return val, errors.New("Given Value is not a list")
	} else {
		return val, nil
	}
}
