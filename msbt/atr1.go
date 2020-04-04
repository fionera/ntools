package msbt

import (
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

type ATR1 struct {
	SectionSize uint32
	EntryCount  uint32
}

const ATR1MagicBytes = "ATR1"

func ReadATR1(r io.Reader, header *FileHeader) *ATR1 {
	if !CheckMagicBytes(r, []byte(ATR1MagicBytes)) {
		logrus.Fatal("Not a ATR1 Section")
	}

	sectionSize := ReadU32(r, header.ByteOrder)
	logrus.Infof("Section size: %d", sectionSize)

	_ = ReadBytes(r, 8)

	entryCount := ReadU32(r, header.ByteOrder)
	logrus.Infof("Entries: %d", entryCount)

	return &ATR1{
		SectionSize: sectionSize,
		EntryCount:  entryCount,
	}
}

func ReadATR1Entry(r io.Reader, header *FileHeader) {
	_ = ReadU32(r, header.ByteOrder)

	_ = ReadU32(r, header.ByteOrder) // Index into TXT2 strings
}
