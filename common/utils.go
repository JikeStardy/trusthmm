package common

import "errors"

func ParseArgs(args []string, params ...*string) error {
	if len(params) > len(args) {
		return errors.New("too less args")
	}
	for idx, ptr := range params {
		*ptr = args[idx]
	}
	return nil
}

func CreateInvokeArgs(invokeFunction string, args ...string) [][]byte {
	ret := make([][]byte, 0, len(args)+1)
	ret = append(ret, []byte(invokeFunction))
	for _, arg := range args {
		ret = append(ret, []byte(arg))
	}
	return ret
}
