package pcm

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/fionera/ntools/bwav"
)

func Decode(f *os.File, header *bwav.FileHeader, info *bwav.ChannelInfo, loopCount int) []int {
	var output []int

	data := make([]byte, int(info.SampleCount)*2)
	_, err := f.ReadAt(data, int64(info.AbsoluteStartOffset))
	if err != nil && err != io.EOF {
		logrus.Fatal(err)
	}

	var currentLoop int
	shouldLoop := false
	if loopCount > 0 {
		shouldLoop = true
	}

	for i := 0; i < len(data); i += 2 {
		if shouldLoop && i >= int(info.LoopEndSample) {
			currentLoop++
			if currentLoop < loopCount {
				i = int(info.LoopStartSample)
			} else {
				i = int(info.LoopEndSample)
			}
		}

		u := header.ByteOrder.Uint16(data[i : i+2])
		output = append(output, int(u))

		if shouldLoop && currentLoop >= loopCount {
			break
		}
	}

	return output
}
