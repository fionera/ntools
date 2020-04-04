package bwav

import (
	"os"
)

type Decoder func(f *os.File, header *FileHeader, info *ChannelInfo, loopCount int) []int
