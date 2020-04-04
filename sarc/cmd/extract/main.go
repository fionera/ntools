package main

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/fionera/ntools/sarc"
)

func init() {
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

		logrus.Infof("Opening File: %s", fileName)

		sarcHeader := sarc.ReadSARCHeader(f)
		sfatHeader := sarc.ReadSFATHeader(f, sarcHeader)
		sfntHeader := sarc.ReadSFNT(f, sarcHeader, sfatHeader)

		for i, node := range sfatHeader.Nodes {
			nodeName := sfntHeader.FileNames[i]
			
			if len(nodeName) == 0 {
				nodeName = fmt.Sprintf("%s_NO_NAME_%d", path.Base(fileName), i)
			}

			nodeName = path.Join(path.Dir(f.Name()), nodeName)

			err := os.MkdirAll(path.Dir(nodeName), 0755)
			if err != nil {
				logrus.Fatalf("error creating folders: %v", err)
			}

			logrus.Infof("Extracting: %s", nodeName)
			nodeFile, err := os.Create(nodeName)
			if err != nil {
				logrus.Fatalf("error creating file: %v", err)
			}

			buf := make([]byte, node.DataEnd-node.DataBeginn)
			_, err = f.ReadAt(buf, int64(sarcHeader.DataBeginn+int(node.DataBeginn)))
			if err != nil {
				logrus.Fatalf("error reading data: %v", err)
			}

			_, err = nodeFile.Write(buf)
			if err != nil {
				logrus.Fatalf("error writing data: %v", err)
			}
		}

		f.Close()
	}
}
