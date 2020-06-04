package parser

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

func parseInt(reader io.ByteScanner) (BnCode, error) {
	rc := BnCode{IsInt: true}
	var buffer []byte

	// check if the stream starts with the correct delimiter for char
	b, err := reader.ReadByte()
	if err != nil {
		return rc, err
	} else if b != 'i' {
		return rc, fmt.Errorf("Unexpected character encountered. Expected %c but got %c", 'i', b)
	}

	isNegative, hasZero := int64(1), false

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
	rc.Value = int64(tmp) * isNegative
	return rc, nil
}

func parseString(reader io.ByteScanner) (BnCode, error) {
	rc := BnCode{IsString: true}
	var buffer []byte

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

func parseList(reader io.ByteScanner) (BnCode, error) {

	rc := BnCode{IsList: true}
	// check if the stream starts with the correct delimiter for list
	b, err := reader.ReadByte()
	if err != nil {
		return rc, err
	} else if b != 'l' {
		return rc, fmt.Errorf("Unexpected character encountered. Expected %c but got %c", 'l', b)
	}

	var tmpList []BnCode

readLoop:
	for {
		if b, err = reader.ReadByte(); err != nil {
			return rc, err
		} // we want to return the read byte, to allow parsing methods to perform a full string scan
		if err := reader.UnreadByte(); err != nil {
			return BnCode{}, err
		}
		switch b {
		case 'e':
			break readLoop
		default:
			if t, err := Decode(reader); err == nil {
				tmpList = append(tmpList, t)
			}
		}
	}
	rc.Value = tmpList
	return rc, nil
}

func parseDict(reader io.ByteScanner) (BnCode, error) {
	var keys []string
	cache := make(map[string]BnCode)

	rc := BnCode{IsDict: true}
	// check if the stream starts with the correct delimiter for dict
	b, err := reader.ReadByte()
	if err != nil {
		return rc, err
	} else if b != 'd' {
		return rc, fmt.Errorf("Unexpected character encountered. Expected %c but got %c", 'd', b)
	}

readLoop:
	for {
		if b, err = reader.ReadByte(); err != nil {
			return rc, err
		} // we want to return the read byte, to allow parsing methods to perform a full string scan
		if err := reader.UnreadByte(); err != nil {
			return BnCode{}, err
		}
		switch b {
		case 'e':
			break readLoop
		default:
			// read the key first, it is always expected to be a string
			key, err := parseString(reader)
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

func Decode(reader io.ByteScanner) (BnCode, error) {
	var rc BnCode
	if b, err := reader.ReadByte(); err != nil {
		return BnCode{}, err
	} else {
		// we want to return the read byte, to allow parsing methods to perform a full string scan
		if err := reader.UnreadByte(); err != nil {
			return BnCode{}, err
		}

		switch b {
		// found an int
		case 'i':
			if obj, err := parseInt(reader); err == nil {
				// append the result of the int parsing to the original slice
				rc = obj
			} else {
				return BnCode{}, err
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if obj, err := parseString(reader); err == nil {
				// append the result of the int parsing to the original slice
				rc = obj
			} else {
				return BnCode{}, err
			}
		case 'l':
			if obj, err := parseList(reader); err == nil {
				// append the result of the int parsing to the original slice
				rc = obj
			} else {
				return BnCode{}, err
			}
		case 'd':
			if obj, err := parseDict(reader); err == nil {
				// append the result of the int parsing to the original slice
				rc = obj
			} else {
				return BnCode{}, err
			}
		}
	}

	return rc, nil
}
