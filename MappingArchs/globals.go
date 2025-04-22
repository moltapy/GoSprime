package main

import "github.com/urfave/cli/v3"

var sep, maskpath, archvcfpath, scorepath, arrayname string

var isdepth, isreverse bool = false, false

var depthTag bool = false

var maxpos int

var handler = &cli.Command{
	Name:  "MappingArchs",
	Usage: "Mapping archaic segments to score file",
}
