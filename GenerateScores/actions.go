package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/codeskyblue/go-sh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var actions = func(ctx context.Context, c *cli.Command) error {

	chroms := make([]string, 0)
	if len(chromlist) > 0 {
		chroms = chromlist
	} else {
		for i := 1; i <= 22; i++ {
			chroms = append(chroms, strconv.Itoa(i))
		}
	}

	var groups []string
	if len(targetgroup) > 0 {
		groups = targetgroup
	} else {
		groups = fullgroups
	}

	if overwrite {
		logrus.Warningf("Overwrite mode on, rewrite path: %s", scorepath)
	}

	fulltasks := len(groups) * len(chroms)
	logrus.Infof("Total tasks: %d, threads: %d, cores: %d, running ... ", fulltasks, threads, incores)
	tunnel = make(chan struct{}, threads)

	for _, group := range groups {

		logrus.Infof("Start processing group %s", group)

		dirgroup := filepath.Join(scorepath, group)
		logrus.Infof("Creating directory of group: %s, path: %s", group, dirgroup)
		if err := os.MkdirAll(dirgroup, os.ModePerm); err != nil {
			return fmt.Errorf("Creating directory: %s error, reason: %v", dirgroup, err)
		}
		logrus.Infof("Complete creating directory for group: %s, result: %s", group, dirgroup)

		var groupWaitGroup sync.WaitGroup
		groupWaitGroup.Add(len(chroms))

		for _, chrom := range chroms {
			output := filepath.Join(dirgroup, group+"_"+outname+"_"+chrom)
			input := strings.ReplaceAll(genopath, "{group}", group)
			waitGroup.Add(1)

			tunnel <- struct{}{}
			go func(toolpath, groupname, inputpath, otgpath, routemap, outputpath, chromosome string) {
				defer func() {
					<-tunnel
					waitGroup.Done()
					groupWaitGroup.Done()
				}()

				SPrimeRun(toolpath, groupname, inputpath, otgpath, routemap, outputpath, chromosome)
			}(jarpath, group, input, outgrouproute, maproute, output, chrom)
		}
		go func(sample string) {
			groupWaitGroup.Wait()
			if !isError {
				logrus.Infof("Complete processing all target chromosomes of group: %s", sample)
			}
		}(group)
	}
	waitGroup.Wait()

	if !isError {
		logrus.Infof("Complete scoring %d groups, result: %s", len(groups), scorepath)
	} else {
		logrus.Errorf("Scoring failed, reason: %v", cmderr)
	}

	return nil
}

func SPrimeRun(toolpath, group, genofile, outgroupfile, mapfile, outfile, chrom string) error {

	parameters := []string{"-Djava.util.concurrent.ForkJoinPool.common.parallelism=12", "-jar", toolpath,
		"gt=" + genofile, "outgroup=" + outgroupfile, "map=" + mapfile, "out=" + outfile, "chrom=" + chrom, "minscore=" + strconv.Itoa(MINSCORE)}

	if !isError {
		err := sh.Command("java", parameters).WriteStdout(os.DevNull)

		if err != nil {
			isError = true
			cmderr = err
			return fmt.Errorf("SPrime process failed, reason: %v", err)
		}

		logrus.Infof("Complete scoring samples in %s, chromosome: %s", group, chrom)
	}

	return nil
}
