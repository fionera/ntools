package util

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"
)

func ReadU16(r io.Reader, bo binary.ByteOrder) uint16 {
	b := make([]byte, 2)

	_, err := r.Read(b)
	if err != nil {
		logrus.Fatal(err)
	}

	return bo.Uint16(b)
}

func ReadU32(r io.Reader, bo binary.ByteOrder) uint32 {
	b := make([]byte, 4)

	_, err := r.Read(b)
	if err != nil {
		logrus.Fatal(err)
	}

	return bo.Uint32(b)
}

func CheckMagicBytes(r io.Reader, expectedBytes []byte) bool {
	readBytes := ReadBytes(r, len(expectedBytes))
	logrus.Info("rb: ", string(readBytes))
	return bytes.Equal(readBytes, expectedBytes)
}

func ReadBytes(r io.Reader, i int) []byte {
	b := make([]byte, i)
	_, err := r.Read(b)
	if err != nil {
		logrus.Fatalf("error reading bytes: %v", err)
	}
	
	return b
}