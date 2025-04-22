package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var actions = func(ctx context.Context, c *cli.Command) error {

	if header {
		logrus.Infof("Opening file: %s using header mode ... ", samplelist)
	} else {
		logrus.Infof("Opening file: %s using headerless mode ... ", samplelist)
	}

	samplefile, err := os.Open(samplelist)
	if err != nil {
		logrus.Fatalf("Opening file: %s error, reason: %v", samplelist, err)
	}
	defer samplefile.Close()
	scanner := bufio.NewScanner(samplefile)

	count := 0

	for scanner.Scan() {
		if header && count < 1 {
			count++
			continue
		}

		var line []string
		if separators != "" {
			line = FormatSplit(separators, scanner.Text())
			if len(line) < 2 {
				return fmt.Errorf("Line %d cannot be splited by format string %s", count+1, separators)
			}
		} else {
			line = strings.Split(scanner.Text(), "\t")
			if len(line) < 2 {
				logrus.Infof("Line %d cannot be splited by \\t, retry using space...", count+1)
				line = strings.Split(scanner.Text(), " ")
				if len(line) < 2 {
					return fmt.Errorf("Line %d cannot be splited by just space or \\t", count+1)
				}
			}
		}

		if outgroup != "" && line[1] == outgroup {
			outgroups = append(outgroups, line[0])
		} else {
			if len(samplelimits) > 0 && slices.Contains(samplelimits, line[1]) {
				groups[line[1]] = append(groups[line[1]], line[0])
			} else if len(samplelimits) == 0 {
				groups[line[1]] = append(groups[line[1]], line[0])
			}
		}
		count++
	}

	if len(outgroups) > 0 {
		for i := range groups {
			groups[i] = append(groups[i], outgroups...)
		}
		WriteList(filepath.Join(outgroupout, "outgroup.txt"), outgroups, false)
	}

	for name, samples := range groups {
		waitGroup.Add(1)
		go WriteList(filepath.Join(sampleout, name+".txt"), samples, true)
	}
	waitGroup.Wait()

	if outgroup != "" {
		logrus.Infof("Complete spliting %d groups merge with %s outgroup, subgroups: %s, outgroup: %s",
			len(groups), outgroup, sampleout, outgroupout)
	} else {
		logrus.Infof("Complete spliting %d groups, subgroups: %s", len(groups), sampleout)
	}

	return nil
}

func WriteList(path string, samples []string, parallel bool) error {
	if parallel {
		defer waitGroup.Done()
	}
	file, err := os.Create(path)
	if err != nil {
		logrus.Fatalf("Creating file: %s error, reason: %v", path, err)
	}
	defer file.Close()
	for _, sample := range samples {
		file.WriteString(sample + "\n")
	}
	return nil
}

func FormatSplit(separators, line string) []string {

	escaped := regexp.QuoteMeta(separators)
	pattern := strings.ReplaceAll(escaped, "|", "([^\\s]*)")

	reg := regexp.MustCompile(pattern)
	reglist := reg.FindStringSubmatch(line)
	return reglist[1:]
}
