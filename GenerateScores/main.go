package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func init() {
	logrus.SetFormatter(&Formatter{})

	handler.Authors = []any{"moltapy"}

	handler.ExitErrHandler = func(ctx context.Context, c *cli.Command, err error) {
		fmt.Fprintf(c.Writer, "ERROR: %v\n", err)
		cli.OsExiter(1)
	}

	handler.OnUsageError = func(ctx context.Context, c *cli.Command, err error, isSubcommand bool) error {
		fmt.Fprintf(c.Writer, "USAGE ERROR: %v\n", err)
		cli.OsExiter(1)
		return nil
	}

	handler.Flags = append(handler.Flags, genopathflag)
	handler.Flags = append(handler.Flags, targetgroupflag)
	handler.Flags = append(handler.Flags, mappathflag)
	handler.Flags = append(handler.Flags, jarpathflag)
	handler.Flags = append(handler.Flags, overwriteflag)
	handler.Flags = append(handler.Flags, outgroupflag)
	handler.Flags = append(handler.Flags, otgnameflag)
	handler.Flags = append(handler.Flags, chromlistflag)
	handler.Flags = append(handler.Flags, scorepathflag)
	handler.Flags = append(handler.Flags, threadsflag)
	handler.Flags = append(handler.Flags, incoresflag)

	handler.Action = actions
}

func main() {
	if err := handler.Run(context.Background(), os.Args); err != nil {
		if len(os.Args) >= 2 {
			logrus.Fatal(err)
		}
	}
}
