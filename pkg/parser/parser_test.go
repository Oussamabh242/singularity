package parser

import (
	"reflect"
	"testing"
)

// type Metadata struct {
// 	Queue       string `json:"queue"`
// 	Topic       string `json:"Topic"`
// 	ContentType string `json:"content-type"`
// }
//
// type Packet struct {
// 	PacketType  int
// 	RLenght     uint
// 	MetadataLen uint
// 	Metadata    Metadata
// 	PayloadLen  uint
// 	Payload     string
// }

func TestParse(t *testing.T) {
	t.Run("bytes array into a uint16", func(t *testing.T) {
		res := Intify[uint16]([]byte{0x12, 0x34})
		if res != 4660 {
			t.Errorf("expected 4660 , got %d", res)
		}
	})
  t.Run("bytes array into a uint32", func(t *testing.T) {
		res := Intify[uint32]([]byte{0x12, 0x34,0x12, 0x34})
		if res != 305402420 {
			t.Errorf("expected 4660 , got %d", res)
		}
	})
	t.Run("parsing a full packet", func(t *testing.T) {
		byteArray := []byte{1,00,00, 18, 52, 0, 17, 123, 34, 113, 117, 101, 117, 101, 34, 58, 32, 34, 116, 104, 97, 116, 34, 125, 0, 12, 104, 101, 108, 108, 111, 32, 116, 104, 101, 114, 101, 48}

		got := Parse(byteArray)
		want := Packet{
			PacketType:  1,
			RLenght:     uint(4660),
			MetadataLen: uint(17),
			PayloadLen:  uint(12),
			Metadata: Metadata{
				Queue:       "that",
				Topic:       "",
				ContentType: "",
			},
			Payload: []byte{104, 101, 108, 108, 111, 32, 116, 104, 101, 114, 101, 48},
		}
		if reflect.DeepEqual(got, want) == false {
			t.Errorf(" %v \n %v ", got, want)
		}
  })
  t.Run("parsing a Metadata only packet ", func(t *testing.T) {
    byteArray := []byte{5,0 ,0, 18, 52, 0, 17, 123, 34, 113, 117, 101, 117, 101, 34, 58, 32, 34, 116, 104, 97, 116, 34, 125}
    got := Parse(byteArray)
    want := Packet{
      PacketType:  5,
      RLenght:     uint(4660),
      MetadataLen: uint(17),
      PayloadLen:  0,
      Metadata: Metadata{
        Queue:       "that",
        Topic:       "",
        ContentType: "",
      },
    }
    if reflect.DeepEqual(got, want) == false {
      t.Errorf("want %v \n got %v", got, want)
    }
  })
  t.Run("parsing Empty packet ", func(t *testing.T) {
    byteArray := []byte{1, 0, 0,0,0}
    got := Parse(byteArray)
    want := Packet{
      PacketType: 1,
      RLenght:    0,
    }
    if reflect.DeepEqual(got, want) == false {
      t.Errorf("want %v \n got %v", got, want)
    }
  })


}
