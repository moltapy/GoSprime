package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"

	"github.com/codeskyblue/go-sh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var actions = func(ctx context.Context, c *cli.Command) error {

	if overwrite {
		logrus.Warningf("Overwrite mode on, rewrite path: %s", outputdir)
	}

	var groups []string
	if len(groupslimit) > 0 {
		groups = groupslimit
	} else {
		groups = fullgroups
	}
	logrus.Infof("Concat ranges: %v", groups)

	for _, group := range groups {
		grouproute := filepath.Join(groupsdir, group)
		filelists := make([]string, 0)
		err := filepath.Walk(grouproute, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Unaccessible path: %s, reason: %v", grouproute, err)
			}
			if !info.IsDir() {
				filelists = append(filelists, path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Walking path: %s error, reason: %v", grouproute, err)
		}
		slices.SortFunc(filelists, func(a, b string) int {
			findIndex := func(path string) int {
				filename := filepath.Base(path)
				re := regexp.MustCompile(`\D*(\d+)\D*`)
				matches := re.FindStringSubmatch(filename)
				if len(matches) < 2 {
					return -1
				}

				num, _ := strconv.Atoi(matches[1])
				return num
			}

			firstNum := findIndex(a)
			secondNum := findIndex(b)

			switch {
			case firstNum < secondNum:
				return -1
			case firstNum > secondNum:
				return 1
			default:
				return 0
			}
		})

		output := filepath.Join(outputdir, group+"_all.vcf.gz")

		if err := Concat(bcftoolpath, output, group, filelists); err != nil {
			return fmt.Errorf("Concat group: %s error, reason: %v", group, err)
		}
	}
	logrus.Infof("Complete concating %d groups, results: %s", len(groups), outputdir)
	return nil
}

func Concat(bcftoolpath, outputpath, group string, sortfiles []string) error {

	parameters := []string{"concat"}
	parameters = append(parameters, sortfiles...)
	parameters = append(parameters, "--naive-force", "--output-type", "z", "--output", outputpath)

	logrus.Infof("Start bcftools concating in group %s", group)

	err := sh.Command(bcftoolpath, parameters).Run()
	if err != nil {
		return fmt.Errorf("Processing bcftools concat error, reason: %v", err)
	}

	logrus.Infof("Complete concating group %s's VCF files", group)

	return nil
}
