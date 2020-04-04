package sarc

import (
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

const SFATMagicBytes = "SFAT"

type SFATHeader struct {
	NodeCount int
	Nodes     []*Node
}

func ReadSFATHeader(r io.Reader, header *SARCHeader) *SFATHeader {
	if !CheckMagicBytes(r, []byte(SFATMagicBytes)) {
		logrus.Fatal("Not a SFAT entry")
	}

	_ = ReadU16(r, header.ByteOrder) // Length

	nodeCount := ReadU16(r, header.ByteOrder) // Node Count

	_ = ReadU32(r, header.ByteOrder) // Hash Key

	var nodes []*Node
	for i := 0; i < int(nodeCount); i++ {
		nodes = append(nodes, ReadNode(r, header))
	}

	return &SFATHeader{
		NodeCount: int(nodeCount),
		Nodes:     nodes,
	}
}
