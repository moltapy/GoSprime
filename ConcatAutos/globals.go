package main

import "github.com/urfave/cli/v3"

var groupsdir, outputdir, bcftoolpath string

var fullgroups, groupslimit []string

var overwrite bool = false

var handler = &cli.Command{
	Name:  "ConcatAutos",
	Usage: "Quickly Concat VCF files of a group in different chromosomes, bcftools options: --naive-force, --output-type z",
}
