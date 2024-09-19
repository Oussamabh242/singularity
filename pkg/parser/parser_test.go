package parser

import (
  "testing"
)


func TestParse(t *testing.T)  {
  test := []byte{224, 128, 30, 128, 8, 77, 121, 47, 84, 111, 112, 105, 99, 128, 11, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, }  
got := Parse(test)
  if got.RLenght != 1 {
    t.Errorf("%v" , got)
  }
}
