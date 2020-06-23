package bencode

import (
	"fmt"
	"sort"
	"strconv"
)

func flattenInt(src BnCode) ([]byte, error) {
	rc := []byte{}
	if src.State != BnInt {
		return rc, fmt.Errorf("Source object does not hold an int value")
	}

	val, err := src.GetInt()
	if err != nil {
		return rc, err
	}

	rc = append(rc, 'i')
	for _, v := range []byte(strconv.Itoa(val)) {
		rc = append(rc, v)
	}
	rc = append(rc, 'e')

	return rc, nil
}

func flattenString(src BnCode) ([]byte, error) {
	rc := []byte{}
	if src.State != BnString {
		return rc, fmt.Errorf("Source object does not hold a string value")
	}

	val, err := src.GetString()
	if err != nil {
		return rc, err
	}
	rc = []byte(strconv.Itoa(len(val)))

	rc = append(rc, ':')
	rc = append(rc, []byte(val)...)

	return rc, nil
}

func flattenList(src BnCode) ([]byte, error) {
	rc, tmp := []byte{}, []byte{}
	if src.State != BnList {
		return tmp, fmt.Errorf("Source object does not hold a list value")
	}

	val, err := src.GetList()
	if err != nil {
		return tmp, err
	}

	rc = append(rc, 'l')
	var enc []byte
	for _, v := range val {
		if enc, err = Encode(v); err != nil {
			return tmp, err
		}
		rc = append(rc, enc...)
	}

	rc = append(rc, 'e')

	return rc, nil
}

func flattenDict(src BnCode) ([]byte, error) {
	rc, emptyRc := []byte{}, []byte{}
	if src.State != BnDict {
		return emptyRc, fmt.Errorf("Source object does not hold a dictionary")
	}

	val, err := src.GetDict()
	if err != nil {
		return emptyRc, err
	}

	// we need to insert the keys in the sorted order, hence
	// we collect and sort the keys first and then iterate over
	// the sorted keys encoding them in the correct order
	keys := make([]string, len(val))
	i := 0
	for key := range val {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	rc = append(rc, 'd')
	for _, key := range keys {
		v := val[key]
		// encode the key by creating a tmp wrapper object
		tmp := BnCode{State: BnString, Value: key}
		enc, err := flattenString(tmp)
		rc = append(rc, enc...)

		// insert the actual value
		if enc, err = Encode(v); err != nil {
			return emptyRc, err
		}
		rc = append(rc, enc...)
	}
	rc = append(rc, 'e')

	return rc, nil
}

// Encode attempts to flatten the src BnCode object into dest stream.
//
// Follows rules described here: https://en.wikipedia.org/wiki/Bencode
func Encode(src BnCode) ([]byte, error) {
	switch src.State {
	case BnInt:
		return flattenInt(src)
	case BnString:
		return flattenString(src)
	case BnList:
		return flattenList(src)
	case BnDict:
		return flattenDict(src)
	default:
		return []byte{}, fmt.Errorf("Unknown type encountered")
	}
}
