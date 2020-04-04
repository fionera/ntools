package main

import (
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/fionera/ntools/msbt"
)

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
		msbtHeader := msbt.ReadHeader(f)

		_ = msbt.ReadLBL1(f, msbtHeader)


		atr1Header := msbt.ReadATR1(f, msbtHeader)
		for i:=0; i < int(atr1Header.EntryCount); i++ {
			msbt.ReadATR1Entry(f, msbtHeader)
		}


		f.Close()
	}
}
