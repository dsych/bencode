package parser

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

func flattenInt(src BnCode, dest io.ByteWriter) error {
	if src.State != BnInt {
		return fmt.Errorf("Source object does not hold an int value")
	}

	val, err := src.getInt()
	if err != nil {
		return err
	}

	dest.WriteByte('i')
	for _, v := range []byte(strconv.Itoa(val)) {
		dest.WriteByte(v)
	}
	dest.WriteByte('e')

	return nil
}

func flattenString(src BnCode, dest io.ByteWriter) error {
	if src.State != BnString {
		return fmt.Errorf("Source object does not hold a string value")
	}

	val, err := src.getString()
	if err != nil {
		return err
	}
	for _, b := range []byte(strconv.Itoa(len(val))) {
		dest.WriteByte(b)
	}
	dest.WriteByte(':')
	for _, v := range []byte(val) {
		dest.WriteByte(v)
	}

	return nil
}

func flattenList(src BnCode, dest io.ByteWriter) error {
	if src.State != BnList {
		return fmt.Errorf("Source object does not hold a list value")
	}

	val, err := src.getList()
	if err != nil {
		return err
	}

	dest.WriteByte('l')
	for _, v := range val {
		if err := Encode(v, dest); err != nil {
			return err
		}
	}

	dest.WriteByte('e')

	return nil
}

func flattenDict(src BnCode, dest io.ByteWriter) error {
	if src.State != BnDict {
		return fmt.Errorf("Source object does not hold a dictionary")
	}

	val, err := src.getDict()
	if err != nil {
		return err
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

	dest.WriteByte('d')
	for _, key := range keys {
		v := val[key]
		// encode the key by creating a tmp wrapper object
		tmp := BnCode{State: BnString, Value: key}
		flattenString(tmp, dest)
		// insert the actual value
		if err := Encode(v, dest); err != nil {
			return err
		}
	}
	dest.WriteByte('e')

	return nil
}

func Encode(src BnCode, dest io.ByteWriter) error {
	var err error = nil
	switch src.State {
	case BnInt:
		err = flattenInt(src, dest)
	case BnString:
		err = flattenString(src, dest)
	case BnList:
		err = flattenList(src, dest)
	case BnDict:
		err = flattenDict(src, dest)
	default:
		err = fmt.Errorf("Unknown type encountered")
	}
	return err
}
