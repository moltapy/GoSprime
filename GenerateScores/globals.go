package main

import (
	"sync"

	"github.com/urfave/cli/v3"
)

const MINSCORE = 150000

var genopath, maproute, jarpath, outgrouproute, scorepath, outname string

var fullgroups, targetgroup, chromlist []string

var incores, cores, threads int64

var default_cores int

var overwrite, isError bool = false, false

var cmderr error

var tunnel chan struct{}

var threadsflag, incoresflag *cli.IntFlag

var waitGroup sync.WaitGroup

var handler = &cli.Command{
	Name:  "GenerateScores",
	Usage: "Execute SPrime scripts to generate scores per group",
}
