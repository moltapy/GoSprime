package main

import "github.com/urfave/cli/v3"

const Active = 1

const Inactive = 0

const Defaultdepth = 1

const MaskSite = 2

const LeftSite = 1

const RightSite = 0

var sep, maskpath, archvcfpath, scorepath, arrayname string

var isdepth, isreverse bool = false, false

var depthTag bool = false

var maxpos int

var handler = &cli.Command{
	Name:  "MappingArchs",
	Usage: "Mapping archaic segments to score file",
}
