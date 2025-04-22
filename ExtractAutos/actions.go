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

	var samplegroups []string

	if len(grouplimits) > 0 {
		samplegroups = grouplimits
	} else {
		samplegroups = samples
	}

	chromlists := make([]string, 0)
	if len(chromolimits) > 0 {
		chromlists = chromolimits
	} else {
		for i := 1; i <= 22; i++ {
			chromlists = append(chromlists, strconv.Itoa(i))
		}
	}

	if overwrite {
		logrus.Warningf("Overwrite mode on, rewrite path: %s", modernout)
	}

	tasktotal := len(samplegroups) * len(chromlists)
	logrus.Infof("Total tasks: %d, threads: %d, cores: %d, running ... ", tasktotal, threads, incores)
	tunnel = make(chan struct{}, threads)

	for _, sample := range samplegroups {

		samplelist := filepath.Join(samplepath, sample+".txt")
		sampleout := filepath.Join(modernout, sample)
		logrus.Infof("Creating directory of group: %s, path: %s", sample, sampleout)
		if err := os.MkdirAll(sampleout, os.ModePerm); err != nil {
			return fmt.Errorf("Creating directory: %s error, reason: %v", sampleout, err)
		}
		logrus.Infof("Complete creating directory for group: %s, result: %s", sample, sampleout)

		var sampleWaitGroup sync.WaitGroup
		sampleWaitGroup.Add(len(chromlists))

		for _, chrom := range chromlists {
			output := filepath.Join(sampleout, strings.ReplaceAll(name, "{chr}", chrom))
			input := strings.ReplaceAll(modernpath, "{chr}", chrom)
			waitGroup.Add(1)
			tunnel <- struct{}{}
			go func(samplelist, input, output, chrom string) {
				defer func() {
					<-tunnel
					waitGroup.Done()
					sampleWaitGroup.Done()
				}()

				ExtractAuto(samplelist, input, output, chrom)
			}(samplelist, input, output, chrom)
		}

		go func(sample string) {
			sampleWaitGroup.Wait()
			if !isError {
				logrus.Infof("Complete processing all target chromosomes of group: %s", sample)
			}

		}(sample)
	}
	waitGroup.Wait()
	if !isError {
		logrus.Infof("Complete splitting %d groups, result: %s", len(samplegroups), modernout)
	} else {
		logrus.Errorf("Splitting failed, reason: %v", cmdError)
		os.Exit(-1)
	}

	return nil
}

func ExtractAuto(samplelist, input, output, chrom string) error {

	command1 := []string{"view", "--samples-file", samplelist, input}
	command2 := []string{"view", "-c1", "-m2", "-M2", "-v", "snps"}
	command3 := []string{"annotate", "-x", "INFO,^FORMAT/GT", "-Oz"}

	if !isError {
		err := sh.Command(bcftoolpath, command1).Command(bcftoolpath, command2).Command(bcftoolpath, command3).WriteStdout(output)
		if err != nil {
			isError = true
			cmdError = err
			return fmt.Errorf("Bcftools process failed, reason: %v", err)
		}

		logrus.Infof("Complete splitting samples in %s, chromosome: %s", samplelist, chrom)
	}

	return nil
}
