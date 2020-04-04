package main

import (
	"os"
	"path"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/sirupsen/logrus"

	"github.com/fionera/ntools/bwav"
	"github.com/fionera/ntools/bwav/codecs/dsp_adpcm"
	"github.com/fionera/ntools/bwav/codecs/pcm"
)

var decoders = map[uint16]bwav.Decoder {
	0: pcm.Decode,
	1: dsp_adpcm.Decode,
}

func combineChannels(in [][]int) (out []int) {
	i := 0
	for {
		added := false
		for _, v := range in {
			if i > len(v)-1 {
				break
			}

			out = append(out, v[i])
			added = true
		}

		if !added {
			break
		}

		i++
	}

	return out
}

func EncodeToWav(f *os.File, header *bwav.FileHeader, channels []*bwav.ChannelInfo) {
	data := make([][]int, len(channels))
	for i, info := range channels {
		logrus.Infof("Decoding Channel %d", i)

		decoder := decoders[info.Codec]
		if decoder == nil {
			logrus.Fatal("Unknown Codec")
			return
		}
		
		data[i] = decoder(f, header, info, *loops)
	}

	fileName := WavFileName(f)
	logrus.Infof("Output: %s", fileName)
	out, err := os.Create(fileName)
	if err != nil {
		logrus.Fatal(err)
	}

	var buffer *audio.IntBuffer

	buffer = &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: len(channels),
			SampleRate:  int(channels[0].SampleRate),
		},
		Data:           combineChannels(data),
		SourceBitDepth: 16,
	}

	e := wav.NewEncoder(out,
		buffer.Format.SampleRate,
		buffer.SourceBitDepth,
		buffer.Format.NumChannels,
		1)
	if err = e.Write(buffer); err != nil {
		logrus.Fatal(err)
	}
	// close the encoder to make sure the headers are properly
	// set and the data is flushed.
	if err = e.Close(); err != nil {
		logrus.Fatal(err)
	}
	out.Close()
}

func WavFileName(f *os.File) string {
	name := path.Base(f.Name())
	return name[:len(name)-len(path.Ext(name))] + ".wav"
}
