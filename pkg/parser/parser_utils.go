package parser

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func intify(x []byte) uint {
	buf := bytes.NewReader(x)
	var num int16
	binary.Read(buf, binary.BigEndian, &num)

	return uint(num)
}

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
