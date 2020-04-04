package sarc

import (
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

const SARCMagicBytes = "SARC"

type SARCHeader struct {
	ByteOrder  binary.ByteOrder
	DataBeginn int
}

func ReadSARCHeader(r io.Reader) *SARCHeader {
	var bo binary.ByteOrder = binary.LittleEndian

	if !CheckMagicBytes(r, []byte(SARCMagicBytes)) {
		logrus.Fatal("Not a SARC entry")
	}

	_ = ReadU16(r, bo) // Length

	byteOrder := ReadU16(r, bo)
	if byteOrder == 0xFFFE {
		logrus.Info("Found LittleEndian file")
		bo = binary.LittleEndian
	} else {
		logrus.Info("Found BigEndian file")
	}

	fileSize := ReadU32(r, bo)
	logrus.Infof("Size: %d", fileSize)

	dataBegin := ReadU32(r, bo) //Begin of Data

	_ = ReadU16(r, bo) // Version

	_ = ReadU16(r, bo)

	return &SARCHeader{
		ByteOrder:  bo,
		DataBeginn: int(dataBegin),
	}
}
