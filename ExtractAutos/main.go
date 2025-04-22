package main

import (
	"context"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

func init() {

	logrus.SetFormatter(&Formatter{})

	cmdHandler.Authors = []any{"moltapy"}

	runtime.GOMAXPROCS(int(incores))

	cmdHandler.Flags = append(cmdHandler.Flags, samplepathflag)
	cmdHandler.Flags = append(cmdHandler.Flags, modernpathflag)
	cmdHandler.Flags = append(cmdHandler.Flags, overwriteflag)
	cmdHandler.Flags = append(cmdHandler.Flags, modernoutflag)
	cmdHandler.Flags = append(cmdHandler.Flags, nameflag)
	cmdHandler.Flags = append(cmdHandler.Flags, chromolimitflag)
	cmdHandler.Flags = append(cmdHandler.Flags, groupslimitflag)
	cmdHandler.Flags = append(cmdHandler.Flags, bcftoolsflag)
	cmdHandler.Flags = append(cmdHandler.Flags, threadsflag)
	cmdHandler.Flags = append(cmdHandler.Flags, incoresflag)

	cmdHandler.Action = actions
}

func main() {
	if err := cmdHandler.Run(context.Background(), os.Args); err != nil {
		if len(os.Args) >= 2 {
			logrus.Fatal(err)
		}
	}
}
