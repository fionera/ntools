package bwav

import (
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"
)

type ChannelInfo struct {
	Codec                uint16    //	0 = PCM, 1 = NGC_DSP
	Pan                  uint16    //	Channel Pan. 0 for left, 1 for right, 2 for middle
	SampleRate           uint32    //	Sample Rate
	SampleCount          uint32    //	Number of samples
	SampleCount2         uint32    //	Number of samples again?
	DSP_ADPCM            [16]int16 //	DSP-ADPCM Coefficients
	AbsoluteStartOffset  uint32    //	Absolute start offset of the sample data
	AbsoluteStartOffset2 uint32    //	Absolute start offset of the sample data again?
	Loop                 uint32    //	Is 1 if the channel loops
	LoopEndSample        uint32    //	Loop End Sample
	LoopStartSample      uint32    //	Loop Start Sample
	PredictorScale       uint16    //	Predictor Scale?
	HistorySample1       uint16    //	History Sample 1?
	HistorySample2       uint16    //	History Sample 2?
	_                    uint16    //	Padding?
}

func ReadChannelInfo(r io.Reader, h *FileHeader) *ChannelInfo {
	c := &ChannelInfo{}

	err := binary.Read(r, h.ByteOrder, c)
	if err != nil {
		logrus.Fatalf("error reading channel info: %v", err)
	}

	return c
}
