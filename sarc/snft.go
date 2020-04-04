package sarc

import (
	"bufio"
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

const SFNTMagicBytes = "SFNT"

type SFNT struct {
	FileNames []string
}

func ReadSFNT(r io.Reader, header *SARCHeader, sfat *SFATHeader) *SFNT {
	if !CheckMagicBytes(r, []byte(SFNTMagicBytes)) {
		logrus.Fatal("Not a SFNT entry")
	}

	_ = ReadU16(r, header.ByteOrder) // Length

	_ = ReadU16(r, header.ByteOrder)

	var fileNames []string

	reader := bufio.NewReader(r)
	for range sfat.Nodes {
		s, err := reader.ReadString(0x00)
		if err != nil {
			logrus.Fatalf("error reading string: %v", err)
		}
		fileNames = append(fileNames, s[:len(s)-1])
	}

	return &SFNT{
		FileNames: fileNames,
	}
}
