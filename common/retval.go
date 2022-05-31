package common

import "fmt"

var (
	OK        = []byte("OK")
	FAIL      = []byte("FAIL")
	EXIST     = []byte("EXIST")
	NOT_EXIST = []byte("NOT EXIST")

	EMPTY_JSON_BYTE = []byte("")
)

func NoSuchFunction(contract, function string) string {
	return fmt.Sprintf("Contract: %s has no function: %s", contract, function)
}
