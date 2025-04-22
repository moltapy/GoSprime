package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {

	logrus.SetFormatter(&Formatter{})

	cmdHandler.Authors = []any{"moltapy"}

	cmdHandler.Flags = append(cmdHandler.Flags, sampleflag)
	cmdHandler.Flags = append(cmdHandler.Flags, outgroupflag)
	cmdHandler.Flags = append(cmdHandler.Flags, sampleoutflag)
	cmdHandler.Flags = append(cmdHandler.Flags, outgroupoutflag)
	cmdHandler.Flags = append(cmdHandler.Flags, limitsflag)
	cmdHandler.Flags = append(cmdHandler.Flags, headerflag)
	cmdHandler.Flags = append(cmdHandler.Flags, formatflag)

	cmdHandler.Action = actions
}

func main() {
	if err := cmdHandler.Run(context.Background(), os.Args); err != nil {
		if len(os.Args) >= 2 {
			logrus.Fatal(err)
		}
	}
}
