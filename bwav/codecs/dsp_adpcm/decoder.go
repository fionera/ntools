package dsp_adpcm

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/fionera/ntools/bwav"
)

var nibbleToInt = []int32{0, 1, 2, 3, 4, 5, 6, 7, -8, -7, -6, -5, -4, -3, -2, -1}

type channel struct {
	hist1 int32
	hist2 int32
}

const BytesPerFrame = 0x08
const SamplesPerFrame = (BytesPerFrame - 0x01) * 2

func Decode(f *os.File, header *bwav.FileHeader, info *bwav.ChannelInfo, loopCount int) []int {
	var output []int

	currentLoop := 0

	shouldLoop := false
	if loopCount > 0 {
		shouldLoop = true
	}

	c := &channel{
		hist1: int32(info.HistorySample1),
		hist2: int32(info.HistorySample2),
	}

	for sample := 0; sample < int(info.SampleCount); sample += SamplesPerFrame {
		if shouldLoop && sample >= int(info.LoopEndSample)-SamplesPerFrame {
			currentLoop++
			if currentLoop < loopCount {
				sample = int(info.LoopStartSample)
			} else {
				sample = int(info.LoopEndSample)
			}
		}

		framesIn := sample / SamplesPerFrame
		frameOffset := int(info.AbsoluteStartOffset) + BytesPerFrame*framesIn

		frame := make([]uint8, BytesPerFrame)
		_, err := f.ReadAt(frame, int64(frameOffset))
		if err != nil && err != io.EOF {
			logrus.Fatal(err)
		}

		samples := decodeFrame(info, c, frame)
		output = append(output, samples...)

		if shouldLoop && currentLoop >= loopCount {
			break
		}
	}

	return output
}

func decodeFrame(info *bwav.ChannelInfo, c *channel, frame []byte) []int {
	var output []int

	var scale int32 = 1 << ((frame[0] >> 0) & 0xF)
	coefIndex := (frame[0] >> 4) & 0xF

	coef1 := int32(info.DSP_ADPCM[coefIndex*2+0])
	coef2 := int32(info.DSP_ADPCM[coefIndex*2+1])

	for i := 0; i < SamplesPerFrame; i++ {
		var sample int32
		nibbles := frame[0x01+i/2]

		if i&1 == 1 {
			sample = nibbleToInt[int(nibbles&0xf)]
		} else {
			sample = nibbleToInt[int(nibbles>>4)]
		}

		sample = (sample * scale) << 11
		sample = (sample + 1024 + coef1*c.hist1 + coef2*c.hist2) >> 11

		if sample > 32767 {
			sample = 32767
		} else if sample < -32768 {
			sample = -32768
		}

		output = append(output, int(sample))

		c.hist2 = c.hist1
		c.hist1 = sample
	}

	return output
}
