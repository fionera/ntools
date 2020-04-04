package msbt

import (
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

type FileHeader struct {
	ByteOrder binary.ByteOrder
	Sections  uint16
	FileSize  uint32
}

const MSBTMagicBytes = "MsgStdBn"

func ReadHeader(r io.Reader) *FileHeader {
	var bo binary.ByteOrder = binary.BigEndian

	if !CheckMagicBytes(r, []byte(MSBTMagicBytes)) {
		logrus.Fatal("Not a MSBT file")
	}

	byteOrder := ReadU16(r, bo)
	if byteOrder == 0xFFFE {
		logrus.Info("Found LittleEndian file")
		bo = binary.LittleEndian
	} else {
		logrus.Info("Found BigEndian file")
	}

	_ = ReadU16(r, bo) // Unknown
	_ = ReadU16(r, bo) // Unknown

	sections := ReadU16(r, bo)
	logrus.Infof("File has %d sections", sections)

	_ = ReadU16(r, bo) // Unknown
	fileSize := ReadU32(r, bo)
	logrus.Infof("Filesize: %d", fileSize)

	unknownBytes := make([]byte, 10)
	_, err := r.Read(unknownBytes)
	if err != nil {
		logrus.Fatalf("error reading bytes: %v", err)
	}

	return &FileHeader{
		ByteOrder:    bo,
		Sections: sections,
		FileSize: fileSize,
	}
}