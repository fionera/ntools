package main

import (
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/fionera/ntools/bwav"
)

var loops *int

func init() {
	loops = flag.Int("loops", 0, "how many times should the song loop")
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		logrus.Fatal("Please supply a file")
	}

	for _, fileName := range flag.Args() {
		f, err := os.Open(fileName)
		if err != nil {
			logrus.Fatalf("error opening file: %v", err)
		}

		if _, err := os.Stat(WavFileName(f)); flag.NArg() > 1 && err == nil {
			logrus.Infof("Skipping File: %s", fileName)
			continue
		}

		logrus.Infof("Opening File: %s", fileName)
		header := bwav.ReadHeader(f)
		channelInfos := make([]*bwav.ChannelInfo, int(header.Channels))

		for i := range channelInfos {
			channelInfos[i] = bwav.ReadChannelInfo(f, header)
		}

		EncodeToWav(f, header, channelInfos)
		f.Close()
	}
}
