package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {

	logrus.SetFormatter(&Formatter{})

	handler.Authors = []any{"moltapy"}

	handler.Flags = append(handler.Flags, groupsdirflag)
	handler.Flags = append(handler.Flags, groupslimitflag)
	handler.Flags = append(handler.Flags, overwriteflag)
	handler.Flags = append(handler.Flags, outputdirflag)
	handler.Flags = append(handler.Flags, bcftoolpathflag)

	handler.Action = actions
}

func main() {

	if err := handler.Run(context.Background(), os.Args); err != nil {
		if len(os.Args) >= 2 {
			logrus.Fatal(err)
		}
	}

}
