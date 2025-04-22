package main

import (
	"sync"

	"github.com/urfave/cli/v3"
)

var samplepath, modernout, modernpath, name, bcftoolpath string

var samples []string = make([]string, 0)

var chromolimits, grouplimits []string

var cores, threads, incores int64

var default_cores int

var waitGroup sync.WaitGroup

var tunnel chan struct{}

var isError bool = false

var overwrite bool = false

var cmdError error

var threadsflag, incoresflag *cli.IntFlag

var cmdHandler = &cli.Command{
	Name:  "ExtractAutos",
	Usage: "Quickly Split Sample variants from VCF file, if you have outgroup, will automatically merge splited groups with outgroup",
}
