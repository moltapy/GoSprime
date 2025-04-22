package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var sampleflag = &cli.StringFlag{
	Name:        "sample",
	Aliases:     []string{"s"},
	Usage:       "Path of samplelist file",
	Required:    true,
	Destination: &samplelist,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessable path: %s, reason: %v", s, err)
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a file", s)
		}
		return nil
	},
}

var outgroupflag = &cli.StringFlag{
	Name:        "outgroup",
	Aliases:     []string{"g"},
	Usage:       "Name of outgroup samples",
	Destination: &outgroup,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		return nil
	},
}

var sampleoutflag = &cli.StringFlag{
	Name:        "samplepath",
	Aliases:     []string{"p"},
	Usage:       "Output path of splited samples",
	Destination: &sampleout,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		_, err := os.Stat(s)
		if os.IsNotExist(err) {
			logrus.Infof("Directory: %s not found, creating ... ", s)
			err := os.MkdirAll(s, os.ModePerm)
			if err != nil {
				logrus.Fatalf("Creating directory: %s error, reason: %v", s, err)
			} else {
				logrus.Infof("Complete making directory: %s", s)
			}
		} else if err != nil {
			return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
		}
		return nil
	},
}

var outgroupoutflag = &cli.StringFlag{
	Name:        "outgrouppath",
	Aliases:     []string{"o"},
	Usage:       "Output path of extracted outgroup samples",
	Destination: &outgroupout,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		_, err := os.Stat(s)
		if os.IsNotExist(err) {
			logrus.Infof("Directory: %s not found, creating ... ", s)
			err := os.MkdirAll(s, os.ModePerm)
			if err != nil {
				logrus.Fatalf("Creating directory %s error, reason: %v", s, err)
			} else {
				logrus.Infof("Complete making directory %s", s)
			}
		} else if err != nil {
			return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
		}
		return nil
	},
}

var limitsflag = &cli.StringSliceFlag{
	Name:        "limits",
	Aliases:     []string{"l"},
	Usage:       "Limits of extracted samples, just extract selected groups, input as 'ASW,ACB/...'",
	Destination: &samplelimits,
	Action: func(ctx context.Context, c *cli.Command, s []string) error {
		return nil
	},
}

var headerflag = &cli.BoolFlag{
	Name:        "title",
	Aliases:     []string{"t"},
	Usage:       "Tag to annotate if samplelist file contains header",
	Value:       true,
	Destination: &header,
	Action: func(ctx context.Context, c *cli.Command, b bool) error {
		return nil
	},
}

var formatflag = &cli.StringFlag{
	Name:        "format",
	Aliases:     []string{"f"},
	Usage:       "Format string to indicate each separator,use | as placehoders for each line,like |\\t| |",
	Destination: &separators,
	Action: func(ctx context.Context, c *cli.Command, s string) error {

		if strings.Count(s, "|") < 2 {
			return fmt.Errorf("Format string error, using '| 'as placeholder twice or more")
		}
		return nil
	},
}
