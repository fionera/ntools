package bwav

import (
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

type FileHeader struct {
	ByteOrder    binary.ByteOrder
	MajorVersion uint16
	MinorVersion uint16
	CRC          uint32
	Channels     uint16
}

const BWAVMagicBytes = "BWAV"

func ReadHeader(r io.Reader) *FileHeader {
	var bo binary.ByteOrder = binary.BigEndian

	if !CheckMagicBytes(r, []byte(BWAVMagicBytes)) {
		logrus.Fatal("Not a BWAV File")
	}

	byteOrder := ReadU16(r, bo)
	if byteOrder == 0xFFFE {
		logrus.Info("Found LittleEndian file")
		bo = binary.LittleEndian
	} else {
		logrus.Info("Found BigEndian file")
	}

	version := ReadU16(r, bo)
	major := version & 0xFF00 >> 8
	minor := version & 0xFF

	logrus.Infof("Version %d-%d", major, minor)

	crc := ReadU32(r, bo) // CRC32
	_ = ReadU16(r, bo) // Padding

	channelCount := ReadU16(r, bo)
	logrus.Infof("Channels: %d", channelCount)

	return &FileHeader{
		ByteOrder:    bo,
		MajorVersion: major,
		MinorVersion: minor,
		CRC:          crc,
		Channels:     channelCount,
	}
}