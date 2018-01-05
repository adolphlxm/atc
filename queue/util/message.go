package util

import (
	"github.com/golang/protobuf/proto"
)

func MustMessageBody(m proto.Message) []byte {
	buf, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	return buf
}

func FromMessageBody(buf []byte, dst proto.Message) error {
    return proto.Unmarshal(buf, dst)
}