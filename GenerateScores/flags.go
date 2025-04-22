package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var genopathflag = &cli.StringFlag{
	Name:        "genopath",
	Aliases:     []string{"g"},
	Usage:       "Path pattern of concated genotype files,use {group} as placeholder for sample group,path separators recommend '/'",
	Destination: &genopath,
	Required:    true,
	Action: func(ctx context.Context, c *cli.Command, s string) error {

		cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {}

		if !strings.Contains(s, "{group}") {
			return fmt.Errorf("Path: %s should contain '{group}' as subgroup placeholder", s)
		}

		parent := filepath.Dir(s)

		fileInfo, err := os.Stat(parent)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessible path: %s, reason: %v", s, err)
		}

		if !fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a directory", s)
		}

		fullgroups = make([]string, 0)
		pattern := strings.ReplaceAll(filepath.Base(s), "{group}", `(\S+)`)
		logrus.Infof("start finding groups in %s", parent)
		err = filepath.Walk(parent, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Unaccessible path: %s, reason: %v", path, err)
			}
			if !info.IsDir() {
				group := MatchGroup(pattern, path)
				fullgroups = append(fullgroups, group)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Walking path: %s error, reason: %v", s, err)
		}
		if len(fullgroups) == 0 {
			return fmt.Errorf("Group matching in %s error, result is empty", parent)
		} else {
			logrus.Infof("Complete group matching in %s, groups %d found", parent, len(fullgroups))
		}
		return nil
	},
}

var targetgroupflag = &cli.StringSliceFlag{
	Name:        "range",
	Aliases:     []string{"r"},
	Usage:       "Slice of target groups, default will use all groups in genopath",
	Destination: &targetgroup,
	Action: func(ctx context.Context, c *cli.Command, s []string) error {
		for _, group := range s {
			if !slices.Contains(fullgroups, group) {
				return fmt.Errorf("Unexpected group: %s, groups list: %v", group, fullgroups)
			}
		}
		logrus.Infof("Target group list: %v", targetgroup)
		return nil
	},
}

var mappathflag = &cli.StringFlag{
	Name:        "maproute",
	Aliases:     []string{"m"},
	Usage:       "Path of reference genome map",
	Destination: &maproute,
	Required:    true,
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

var jarpathflag = &cli.StringFlag{
	Name:        "jarpath",
	Aliases:     []string{"j"},
	Usage:       "Path of sprime.jar",
	Destination: &jarpath,
	Required:    true,
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
	Aliases:     []string{"u"},
	Usage:       "Path of outgroup samples list",
	Destination: &outgrouproute,
	Required:    true,
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

var otgnameflag = &cli.StringFlag{
	Name:        "outname",
	Aliases:     []string{"n"},
	Usage:       "Name of the outgroup",
	Destination: &outname,
	Required:    true,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		return nil
	},
}

var chromlistflag = &cli.StringSliceFlag{
	Name:        "chromlist",
	Aliases:     []string{"l"},
	Usage:       "Slice of chromosomes to execute SPrime,default will be all autosomes",
	Destination: &chromlist,
	Action: func(ctx context.Context, c *cli.Command, s []string) error {
		for _, chromo := range s {

			chrom, err := strconv.Atoi(chromo)
			if err != nil {
				return fmt.Errorf("Chromosome index range contains non-int value: %s", chromo)
			}
			if chrom > 22 || chrom < 1 {
				return fmt.Errorf("Chromosome index value: %s invalid, not in 1 ~ 22", chromo)
			}
		}

		if len(chromlist) > 0 {
			logrus.Infof("Target chromosome list: %v", chromlist)
		}

		return nil
	},
}

var overwriteflag = &cli.BoolFlag{
	Name:  "overwrite",
	Usage: "Tag for overwrite existing output directory or not",
	Action: func(ctx context.Context, c *cli.Command, b bool) error {
		overwrite = c.Bool("overwrite")
		return nil
	},
}

var scorepathflag = &cli.StringFlag{
	Name:        "output",
	Aliases:     []string{"o"},
	Usage:       "Path of output score files, will be the parent directory of each group's directory",
	Destination: &scorepath,
	Required:    true,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				logrus.Infof("Directory %s not found, creating ... ", s)
				mkerr := os.MkdirAll(s, os.ModePerm)
				if mkerr != nil {
					logrus.Fatalf("Creating directory: %s error, reason: %v", s, err)
				} else {
					logrus.Infof("Complete making: %s, continue ... ", s)
					fileInfo, err = os.Stat(s)
					if err != nil {
						logrus.Fatalf("Verification directory: %s error, reason: %v", s, err)
					}
				}
			} else {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}

		}

		if !fileInfo.IsDir() {
			logrus.Infof("Directory %s not found, creating ... ", s)
			mkerr := os.MkdirAll(s, os.ModePerm)
			if mkerr != nil {
				logrus.Fatalf("Creating directory: %s error, reason: %v", s, err)
			} else {
				logrus.Infof("Complete making: %s, continue ... ", s)
			}
		} else {
			if isEmpty, _ := IsDirEmpty(s); !isEmpty && !overwrite {
				logrus.Warningf("Directory: %s found and not empty, using '--overwrite' to overwrite, skipping ... ", s)
				os.Exit(0)
				return nil
			}
		}
		return nil
	},
}

func init() {

	cores = int64(runtime.NumCPU())
	switch cores <= 16 {
	case true:
		if cores < 8 {
			default_cores = int(cores)
		} else {
			default_cores = 8
		}

	case false:
		default_cores = int(cores / 4)
	}

	threadsflag = &cli.IntFlag{
		Name:        "threads",
		Aliases:     []string{"t"},
		Usage:       " Int value of max goroutines parallel",
		Value:       int64(default_cores / 2),
		Destination: &threads,
		Action: func(ctx context.Context, c *cli.Command, i int64) error {
			if i > cores/2 && cores > 16 {
				logrus.Warningf("Threads: %d, all cores: %d", i, cores)
			}
			return nil
		},
	}

	incoresflag = &cli.IntFlag{
		Name:        "cores",
		Aliases:     []string{"c"},
		Usage:       "Int value to define cores using for parallel",
		Value:       int64(default_cores),
		Destination: &incores,
		Action: func(ctx context.Context, c *cli.Command, i int64) error {
			if incores <= 0 || incores > cores {
				return fmt.Errorf("Cores: %d invalid, max cores: %d", incores, cores)
			}
			if incores >= cores/3 && cores > 16 {
				logrus.Warningf("Cores: %d, more than 1/3 all cores: %d", incores, cores)
			}
			return nil
		},
	}
}

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	names, err := f.Readdirnames(1)
	if err != nil && err != io.EOF {
		return false, err
	}
	return len(names) == 0, nil
}

func MatchGroup(pattern, path string) string {

	filename := filepath.Base(path)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(filename)
	bestmatch := ""

	if len(matches) < 1 {
		logrus.Warningf("Fail to match group in %s, ensure it's not necessary", path)
	} else if len(matches) == 2 {
		bestmatch = matches[1]
		logrus.Infof("Group matched in %s, matched group: %s", path, bestmatch)
	} else {
		logrus.Warningf("Path: %s have multiple matches, match list: %v, continue with %s", path, matches[1:], bestmatch)
	}
	return bestmatch
}
