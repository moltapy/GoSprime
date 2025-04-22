package main

import (
	"sync"

	"github.com/urfave/cli/v3"
)

var header bool

var samplelist, outgroup, sampleout, outgroupout, separators string

var groups map[string][]string = make(map[string][]string)

var outgroups []string = make([]string, 0)

var waitGroup sync.WaitGroup

var samplelimits []string

var cmdHandler = &cli.Command{
	Name:  "ExtractSamples",
	Usage: "Quickly Split Samples from Samplelist.txt, if outgroup then add outgroup to each group",
}
