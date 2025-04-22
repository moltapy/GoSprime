package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {

	logrus.SetFormatter(&Formatter{})

	handler.Authors = []any{"moltapy"}

	handler.Flags = append(handler.Flags, bedpathflag)
	handler.Flags = append(handler.Flags, archvariantsflag)
	handler.Flags = append(handler.Flags, scorepathflag)
	handler.Flags = append(handler.Flags, arraynameflag)
	handler.Flags = append(handler.Flags, separatorflag)
	handler.Flags = append(handler.Flags, reverseflag)
	handler.Flags = append(handler.Flags, isdepthflag)

	handler.Action = actions
}

func main() {
	if err := handler.Run(context.Background(), os.Args); err != nil {
		if len(os.Args) > 1 {
			logrus.Fatal(err)
		}
	}
}
