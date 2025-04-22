package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var groupsdirflag = &cli.StringFlag{
	Name:        "dir",
	Aliases:     []string{"d"},
	Usage:       "Path of the parent directory of groups containing VCF files",
	Required:    true,
	Destination: &groupsdir,
	Action: func(ctx context.Context, c *cli.Command, s string) error {
		fileInfo, err := os.Stat(s)

		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Invalid path: %s, reason: %v", s, err)
			}
			return fmt.Errorf("Unaccessable path: %s, reason: %v", s, err)
		}

		if !fileInfo.IsDir() {
			return fmt.Errorf("Invalid path: %s, reason: not a directory", s)
		}

		fullgroups = make([]string, 0)
		err = filepath.Walk(s, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Unaccessible path: %s, reason: %v", s, err)
			}
			if info.IsDir() {
				group := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
				if group != strings.Split(s, "/")[len(strings.Split(s, "/"))-1] {
					fullgroups = append(fullgroups, group)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Walking path: %s error, reason: %v", s, err)
		}
		return nil
	},
}

var groupslimitflag = &cli.StringSliceFlag{
	Name:        "groups",
	Aliases:     []string{"g"},
	Usage:       "Target groups to concat, default will use all in input dir",
	Destination: &groupslimit,
	Action: func(ctx context.Context, c *cli.Command, s []string) error {
		for _, group := range s {
			if !slices.Contains(fullgroups, group) {
				return fmt.Errorf("Unexpected group: %s, groups list: %v", group, fullgroups)
			}
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

var outputdirflag = &cli.StringFlag{
	Name:        "outdir",
	Aliases:     []string{"o"},
	Usage:       "Directory path of concated VCF files per group",
	Destination: &outputdir,
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
			logrus.Infof("Directory: %s not found, creating...", s)
			mkerr := os.MkdirAll(s, os.ModePerm)
			if mkerr != nil {
				logrus.Fatalf("Creating directory: %s error, reason: %v", s, err)
			} else {
				logrus.Infof("Complete making: %s, continue ... ", s)
			}
		} else {
			if isEmpty, _ := IsDirEmpty(s); !isEmpty && !overwrite {
				logrus.Warningf("Directory: %s found and not empty, using '--overwrite' to overwrite", s)
				os.Exit(0)
				return nil
			}
		}
		return nil
	},
}

var bcftoolpathflag = &cli.StringFlag{
	Name:    "bcftools",
	Aliases: []string{"b"},
	Usage: `Path of bcftools, if you use conda to manage envs, just activate env with bcftools, 
	else you should put bcftools path here, miniconda3 path is always : 
	~/miniconda3/envs/bcftools/bin/bcftools`,
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
