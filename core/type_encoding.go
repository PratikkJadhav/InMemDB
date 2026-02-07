package core

import "errors"

func get_type(te uint8) uint8 {
	return (te >> 4) << 4
}

func get_encoding(te uint8) uint8 {
	return (te << 4) >> 4
}

func assertType(te uint8, t uint8) error {
	if get_type(te) != t {
		return errors.New("the operation is not permitted on this type")
	}

	return nil
}

func assertEncoding(te uint8, t uint8) error {
	if get_encoding(te) != t {
		return errors.New("the operation is not permitted on this encoding")
	}

	return nil
}
