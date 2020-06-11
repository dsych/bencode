package parser

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

func parseInt(reader io.ByteReader, firstChar byte) (BnCode, error) {
	rc := BnCode{IsInt: true}
	var buffer []byte

	// check if the stream starts with the correct delimiter for char
	if firstChar != 'i' {
		return rc, fmt.Errorf("Unexpected character encountered. Expected %c but got %c", 'i', firstChar)
	}

	isNegative, hasZero := int(1), false

	var err error
	var b byte

readLoop:
	for {
		if b, err = reader.ReadByte(); err != nil {
			return rc, err
		}

		switch b {
		case '-':
			isNegative = -1
		case 'e':
			// terminate the outter loop, we found the termination delimiter
			break readLoop
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			buffer = append(buffer, b)
		case '0':
			if hasZero {
				return rc, fmt.Errorf("Leading zeros are not allowed")
			}
			buffer = append(buffer, b)
			hasZero = true
		default:
			return rc, fmt.Errorf("Unexpected character encountered. Expected a digit, sign or e, got %c", b)
		}

	}

	tmp, err := strconv.Atoi(string(buffer))
	if err != nil {
		return rc, err
	}

	if tmp == 0 && isNegative == -1 {
		return rc, fmt.Errorf("Negative zeros are not allowed")
	}

	// done processing, make sure to apply the sign correctly
	rc.Value = int(tmp) * isNegative
	return rc, nil
}

func parseString(reader io.ByteReader, firstChar byte) (BnCode, error) {
	rc := BnCode{IsString: true}
	var buffer []byte = []byte{firstChar}

	// attemp to read the length of a string
readLoop:
	for {
		if b, err := reader.ReadByte(); err == io.EOF {
			break
		} else if err != nil {
			return rc, err
		} else {
			switch b {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				buffer = append(buffer, b)
			case ':': // delimiter character that divides the length of the string and the actual value
				break readLoop
			default:
				// unexpected character seen, stop
				return rc, fmt.Errorf("Unexpected character '%c' encountered while parse string", b)
			}
		}
	}
	// get the length of the string
	length, err := strconv.Atoi(string(buffer))
	if err != nil {
		return rc, err
	}
	// clear the buffer
	buffer = nil

	// iterate over the entire string. throw if the length is less than the state length
	for i := 0; i < length; i++ {
		if b, err := reader.ReadByte(); err == io.EOF {
			return rc, fmt.Errorf("String length is shorter than expected. Expected %d got %d", length, i)
		} else if err != nil {
			return rc, err
		} else {
			buffer = append(buffer, b)
		}
	}

	// done parsing, record the string value
	rc.Value = string(buffer)
	return rc, nil
}

func parseList(reader io.ByteReader, firstChar byte) (BnCode, error) {

	rc := BnCode{IsList: true}
	// check if the stream starts with the correct delimiter for list
	if firstChar != 'l' {
		return rc, fmt.Errorf("Unexpected character encountered. Expected %c but got %c", 'l', firstChar)
	}

	var tmpList []BnCode
	var err error
	var b byte

readLoop:
	for {
		if b, err = reader.ReadByte(); err != nil {
			return rc, err
		} // we want to return the read byte, to allow parsing methods to perform a full string scan
		switch b {
		case 'e':
			break readLoop
		default:
			if t, err := decode(reader, b); err == nil {
				tmpList = append(tmpList, t)
			}
		}
	}
	rc.Value = tmpList
	return rc, nil
}

func parseDict(reader io.ByteReader, firstChar byte) (BnCode, error) {
	var keys []string
	cache := make(map[string]BnCode)

	rc := BnCode{IsDict: true}
	// check if the stream starts with the correct delimiter for dict
	if firstChar != 'd' {
		return rc, fmt.Errorf("Unexpected character encountered. Expected %c but got %c", 'd', firstChar)
	}

	var b byte
	var err error

readLoop:
	for {
		if b, err = reader.ReadByte(); err != nil {
			return rc, err
		} // we want to return the read byte, to allow parsing methods to perform a full string scan
		switch b {
		case 'e':
			break readLoop
		default:
			// read the key first, it is always expected to be a string
			key, err := parseString(reader, b)
			if err != nil {
				return rc, err
			}

			// get the actual value that could be anything
			val, err := Decode(reader)
			if err != nil {
				return rc, err
			}

			// save the key for later
			keyStr, err := key.getString()
			if err != nil {
				return rc, fmt.Errorf("Unable to convert to string key, %v", err)
			}

			keys = append(keys, keyStr)
			cache[keyStr] = val
		}
	}

	orgKeys := make([]string, len(keys))
	copy(orgKeys, keys)
	sort.Strings(keys)

	// make sure that the original key order is sorted
	for i, k := range orgKeys {
		if k != keys[i] {
			return rc, fmt.Errorf("Dictionary keys are not in lexicographical order")
		}
	}

	rc.Value = cache
	return rc, nil
}

func Decode(reader io.ByteReader) (BnCode, error) {
	if b, err := reader.ReadByte(); err != nil {
		return BnCode{}, err
	} else {
		return decode(reader, b)
	}
}

func decode(reader io.ByteReader, firstChar byte) (BnCode, error) {
	var rc BnCode
	b := firstChar

	switch b {
	// found an int
	case 'i':
		if obj, err := parseInt(reader, b); err == nil {
			// append the result of the int parsing to the original slice
			rc = obj
		} else {
			return BnCode{}, err
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if obj, err := parseString(reader, b); err == nil {
			// append the result of the int parsing to the original slice
			rc = obj
		} else {
			return BnCode{}, err
		}
	case 'l':
		if obj, err := parseList(reader, b); err == nil {
			// append the result of the int parsing to the original slice
			rc = obj
		} else {
			return BnCode{}, err
		}
	case 'd':
		if obj, err := parseDict(reader, b); err == nil {
			// append the result of the int parsing to the original slice
			rc = obj
		} else {
			return BnCode{}, err
		}
	}

	return rc, nil
}
