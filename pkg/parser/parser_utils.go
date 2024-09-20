package parser

import (
	"strconv"
)

func byteToBinary(integer byte) string {
	b := strconv.FormatInt(int64(integer), 2)
	ns := ""
	for i := 0; i < (7 - len(b)); i++ {
		ns += "0"
	}
	return ns + b
}

func binaryToUint(binary string) uint {
	num, _ := strconv.ParseUint(binary, 2, 0)

	return uint(num)
}
