package sarc

import (
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"
)

type Node struct {
	FileNameHash   uint32
	FileAttributes uint32
	DataBeginn     uint32
	DataEnd        uint32
}

func ReadNode(r io.Reader, header *SARCHeader) *Node {
	n := &Node{}
	err := binary.Read(r, header.ByteOrder, n)
	if err != nil {
		logrus.Fatal(err)
	}
	return n
}

