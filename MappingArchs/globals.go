package main

import "github.com/urfave/cli/v3"

const Active = 1

const Inactive = 0

const Defaultdepth = 1

const MaskSite = 4

const LeftSite = 2

const RightSite = 0

const Mask_A, Mask_T = 0, 1

const Mask_C, Mask_G = 2, 3

var sep, maskpath, archvcfpath, scorepath, arrayname string

var isdepth, isreverse bool = false, false

var depthTag bool = false

var maxpos int

var typeMask map[byte]int = map[byte]int{
	'A': Mask_A,
	'T': Mask_T,
	'G': Mask_G,
	'C': Mask_C,
}

var handler = &cli.Command{
	Name:  "MappingArchs",
	Usage: "Mapping archaic segments to score file",
}
