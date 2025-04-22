package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var samplepathflag = &cli.StringFlag{
	Name:        "samplepath",
	Aliases:     []string{"s"},
	Usage:       "Path of the parent directory of samplelist files per group",
	Required:    true,
	Destination: &samplepath,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessible path: %s, reason: %v", s, err)
		}

		if !fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a directory", s)
		}
		err = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Unaccessible path: %s, reason: %v", path, err)
			}
			if !info.IsDir() {
				group := strings.Split(strings.Split(path, "/")[len(strings.Split(path, "/"))-1], ".")[0]
				samples = append(samples, group)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Walking path: %s error, reason: %v", s, err)
		}
		return nil
	},
}

var modernpathflag = &cli.StringFlag{
	Name:        "modernpath",
	Aliases:     []string{"m"},
	Usage:       "Path of the VCF files for modern human, use '{chr}' as chromosome index placeholder",
	Required:    true,
	Destination: &modernpath,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		if !strings.Contains(s, "{chr}") {
			return fmt.Errorf("Path: %s should contain '{chr}' as chromsome index placeholder", s)
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

var modernoutflag = &cli.StringFlag{
	Name:        "modernout",
	Aliases:     []string{"o"},
	Usage:       "Path of the vcf files for modern human, will be used as the parent dir of splited vcf files per group",
	Required:    true,
	Destination: &modernout,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				logrus.Infof("Directory: %s not found, creating ... ", s)
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
			logrus.Infof("Directory: %s not found, creating...", s)
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

var nameflag = &cli.StringFlag{
	Name:        "name",
	Aliases:     []string{"n"},
	Usage:       "Name of the vcf files for modern human, use '{chr}' as chromosome index placeholder",
	Value:       "chrom_{chr}.vcf.gz",
	Destination: &name,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		if !strings.Contains(s, "{chr}") {
			return fmt.Errorf("Name: %s should contain '{chr}' as chromsome index placeholder", s)
		}
		return nil
	},
}

var chromolimitflag = &cli.StringSliceFlag{
	Name:        "autos",
	Aliases:     []string{"a"},
	Usage:       "Slice of the chromosomes to process, default will be 1-22",
	Destination: &chromolimits,
	Action: func(ctx context.Context, c *cli.Command, s []string) error {
		for _, chromo := range s {

			chromosome, err := strconv.Atoi(chromo)
			if err != nil {
				return fmt.Errorf("Chromosome index range contains non-int value: %s", chromo)
			}
			if chromosome > 22 || chromosome < 1 {
				return fmt.Errorf("Chromosome index value: %s invalid, not in 1 ~ 22", chromo)
			}
		}
		return nil
	},
}

var groupslimitflag = &cli.StringSliceFlag{
	Name:        "groups",
	Aliases:     []string{"g"},
	Usage:       "Slice of the groups to process, default will be all groups in samplepath",
	Destination: &grouplimits,
	Action: func(ctx context.Context, c *cli.Command, s []string) error {
		if len(s) > 0 {
			for _, item := range s {
				if !slices.Contains(samples, item) {
					return fmt.Errorf("Unexpected group: %s, groups list: %v", item, samples)
				}
			}
		}
		return nil
	},
}

var bcftoolsflag = &cli.StringFlag{
	Name:    "bcftools",
	Aliases: []string{"b"},
	Usage: `Path of bcftools, if you use conda to manage envs, just activate env with bcftools, 
	else you should put bcftools path here, miniconda3 path is always : ~/miniconda3/envs/bcftools/bin/bcftools`,
	Value:       "bcftools",
	Destination: &bcftoolpath,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid bcftools path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessible bcftools path: %s, reason: %v", s, err)
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a file", s)
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
				logrus.Infof("Threads: %d, all cores: %d", i, cores)
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
