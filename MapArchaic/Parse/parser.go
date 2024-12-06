package parse

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Args struct {
	BedMode   *string
	SepChar   *string
	MskFile   *string
	VcfFile   *string
	ScoreFile *string
	RefTag    *string
	OutFile   *string
	ReadDepth *string
}

func (args *Args) Parse() {
	args.BedMode = flag.String("bed", "", "Tag for indicating way to use the prama MASK FILE, use 'include' for including points in mask file,'exclude' for excluding points in mask file(default:Include all points in score file)")
	args.SepChar = flag.String("sep", "\t", "Separator of the output file")
	args.MskFile = flag.String("msk", "", "Mask file, only one allowed as the input")
	args.ScoreFile = flag.String("score", "", "Score file from Sprime")
	args.RefTag = flag.String("tag", "", "Tag for the added column")
	args.OutFile = flag.String("out", "", "Mapped score file path")
	args.ReadDepth = flag.String("depth", "false", "Show read depth in result file(optional), bool, true for showing, false for not showing")
	flag.Parse()
	if len(os.Args) > 2 {
		if *args.BedMode != "" && *args.MskFile == "" {
			log.Fatal("If you use the prama BED, The prama MASK FILE cannot be nil, Please check!")
		} else if *args.MskFile != "" && *args.BedMode == "" {
			log.Println("If you use the prama MASK FILE, and may expect that it works, You may indicate include/exclude the mask file in the prama BED")
		} else if *args.BedMode != "include" && *args.BedMode != "exclude" && *args.BedMode != "" {
			log.Fatal("The prama BED only has three options, include/exclude or nil(default),Please check!")
		} else if *args.ScoreFile == "" {
			log.Fatal("The prama SCORE FILE cannot be nil, Please check!")
		} else if *args.RefTag == "" {
			log.Fatal("The prama REFERENCE TAG cannot be nil, Please check!")
		} else if *args.ReadDepth != "true" && *args.ReadDepth != "false" {
			log.Fatal("The depth received is just true/false, Please check!")
		} else if *args.OutFile == "" {
			log.Fatal("The Prama OUTPUT FILE cannot be nil, Please check!")
		} else {
			log.Default()
		}
	} else if len(os.Args) == 1 {
		programPath, err := os.Executable()
		if err != nil {
			log.Fatal("Problem occurred when getting abspath, Please report this bug!")
		}
		fmt.Printf("Usage of %s:\n", programPath)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
