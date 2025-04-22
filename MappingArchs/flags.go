package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var separatorflag = &cli.StringFlag{
	Name:        "separator",
	Aliases:     []string{"s"},
	Usage:       "Tag for indicating the separator in output file",
	Destination: &sep,
	Value:       "\t",
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		if len(s) < 1 {
			logrus.Warningf("Nil separator, columns will attach to the former")
		}
		return nil
	},
}

var bedpathflag = &cli.StringFlag{
	Name:        "maskpath",
	Aliases:     []string{"m"},
	Usage:       "Path for mask file formatted as BED, if provided default will be exclude mask segments, otherwise will keep all segments, use --reverse to just include mask segments",
	Destination: &maskpath,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessible path: %s, reason: %v", s, err)
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a file", s)
		}
		return nil
	},
}

var archvariantsflag = &cli.StringFlag{
	Name:        "vcfarch",
	Aliases:     []string{"v"},
	Usage:       "Path of archaic VCF file",
	Required:    true,
	Destination: &archvcfpath,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessible path: %s, reason: %v", s, err)
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a file", s)
		}
		return nil
	},
}

var arraynameflag = &cli.StringFlag{
	Name:        "arrayname",
	Aliases:     []string{"n"},
	Usage:       "Tag for the added match result column",
	Required:    true,
	Destination: &arrayname,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		return nil
	},
}

var scorepathflag = &cli.StringFlag{
	Name:        "scorepath",
	Aliases:     []string{"p"},
	Usage:       "Path of generated score file",
	Required:    true,
	Destination: &scorepath,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessible path: %s, reason: %v", s, err)
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a file", s)
		}

		CheckScore(s, arrayname)

		return nil
	},
}

var reverseflag = &cli.BoolFlag{
	Name:  "reverse",
	Usage: "Tag for indicating the usage of mask bed file, add to use mask file as include mask, otherwise exclude",
	Action: func(ctx context.Context, c *cli.Command, b bool) error {
		if isreverse = c.Bool("reverse"); isreverse && maskpath == "" {
			return fmt.Errorf("Reverse mode error, provide mask file first")
		}
		return nil
	},
}

var isdepthflag = &cli.BoolFlag{
	Name:  "depth",
	Usage: "Tag for indicating if showing read depth or not, add to add a column showing read depth",
	Action: func(ctx context.Context, c *cli.Command, b bool) error {
		isdepth = c.Bool("depth")
		return nil
	},
}

func CheckScore(path, arrayname string) {

	logrus.Infof("Checking validation of file: %s ... ", path)

	file, err := os.Open(path)
	if err != nil {
		logrus.Errorf("Opening file: %s error, reason: %v", path, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil && err != io.EOF {
		logrus.Errorf("Reading header from buffer error, reason: %v", err)
		os.Exit(-1)
	}
	if err == io.EOF {
		logrus.Errorf("File: %s contains no informative lines", path)
		os.Exit(-1)
	}
	if slices.Contains(strings.Split(string(line), "\t"), arrayname) {
		logrus.Fatalf("Score file already contains column %s", arrayname)
		os.Exit(-1)
	}

	lastline := ""

	for {
		line, _, err := reader.ReadLine()

		if err != nil && err != io.EOF {
			logrus.Errorf("Reading file: %s error, reason: %v", path, err)
			os.Exit(-1)
		}
		if err == io.EOF {
			if !(len(line) > 0) {
				if lastline != "" {
					maxpos, err = strconv.Atoi(strings.Split(lastline, "\t")[1])
					if err != nil {
						logrus.Fatalf("Last line pos: %s is non-int value", strings.Split(lastline, "\t")[1])
						os.Exit(-1)
					}
				} else {
					logrus.Fatalf("File: %s contains no informative lines", path)
					os.Exit(-1)
				}
			} else {
				maxpos, err = strconv.Atoi(strings.Split(string(line), "\t")[1])
				if err != nil {
					logrus.Fatalf("Last line pos: %s is non-int value", strings.Split(string(line), "\t")[1])
					os.Exit(-1)
				}
			}
			break
		}
		if len(line) > 0 {
			lastline = string(line)
		}
	}

	logrus.Infof("File: %s valid, max length: %d", path, maxpos)

}
