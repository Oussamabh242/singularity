package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// intify a generic function that takes an array of bytes and return
// its encoding into INTS or UNITS
func Intify[T any](x []byte  ) T {
	buf := bytes.NewReader(x)
	var num T

  err := binary.Read(buf, binary.BigEndian, &num)
  if err != nil{
    fmt.Println("error" , err)
  }

	return num
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
