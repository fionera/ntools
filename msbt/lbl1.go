package msbt

import (
	"io"

	"github.com/sirupsen/logrus"

	. "github.com/fionera/ntools/util"
)

type LBL1 struct {
	GroupCount uint32

	Groups []*Group
}

const LBL1MagicBytes = "LBL1"

func ReadLBL1(r io.Reader, header *FileHeader) *LBL1 {
	if !CheckMagicBytes(r, []byte(LBL1MagicBytes)) {
		logrus.Fatal("Not a LBL1 Section")
	}

	sectionSize := ReadU32(r, header.ByteOrder)
	logrus.Infof("Section Size: %d", sectionSize)
	
	ReadBytes(r, 8) // Padding
	
	groupCount := ReadU32(r, header.ByteOrder)
	logrus.Infof("Group count: %d", groupCount)

	var groups []*Group
	for i := 0; i < int(groupCount); i++ {
		groups = append(groups, ReadLBL1Group(r, header))
	}

	for _, g := range groups {
		for i := 0; i < int(g.LabelCount); i++ {
			g.Labels = append(g.Labels, ReadLBL1Label(r, header))
		}
	}

	return &LBL1{
		GroupCount: groupCount,
		Groups:     groups,
	}
}

type Group struct {
	LabelCount uint32
	Offset     uint32

	Labels []*Label
}

func ReadLBL1Group(r io.Reader, header *FileHeader) *Group {
	return &Group{
		LabelCount: ReadU32(r, header.ByteOrder),
		Offset:     ReadU32(r, header.ByteOrder),
	}
}

type Label struct {
	Name  string
	Index uint32
}

func ReadLBL1Label(r io.Reader, header *FileHeader) *Label {
	name := ReadBytes(r, int(10))
	logrus.Info(string(name))
	//strLen := ReadU32(r, header.ByteOrder)
	//name := ReadBytes(r, int(strLen))
	index := ReadU32(r, header.ByteOrder)

	return &Label{
		Name:  string(name),
		Index: index,
	}
}
