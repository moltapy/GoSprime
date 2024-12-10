package parse

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Args struct {
	WorkPath    *string
	BedMode     *string
	SepChar     *string
	MskFile     *string
	ArchaicFile *string
	ScoreFile   *string
	RefTag      *string
	OutFile     *string
	ReadDepth   *string
}

// TODO: fill the left MAPARCHAIC AND BUILD WORKFLOW FOR IT

func (args *Args) Parse() {

	directory, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory :%v", err)
	}
	args.WorkPath = flag.String("w", directory, "Path of the subpopulation directories, should be the ")
	args.BedMode = flag.String("b", "", "Tag for indicating how to use the prama MASK FILE,options: 'include'/'exclude'; 'include' for including points in mask file, 'exclude' for excluding points in mask file(default: include all points in score file)")
	args.SepChar = flag.String("c", "\t", "Separator of the score file")
	args.MskFile = flag.String("m", "", "Path of the MASK FILE, should be a bed file for a specific chromosome")
	args.ArchaicFile = flag.String("a", "", "Path of the ARCHAIC VCF FILE, should be a gzip VCF file for a specific chromosome")
	args.ScoreFile = flag.String("s", "", "Path of SPRIME GENERATED SCORE FILE")
	args.RefTag = flag.String("t", "", "Tag for ADDED COLUMN")
	args.OutFile = flag.String("o", "", "Path of OUTPUT SCORE FILE WITH MATCHING TAG COLUMN")
	args.ReadDepth = flag.String("d", "false", "Tag for containing READ DEPTH INFO in result file, should be 'true' for showing, 'false' for not showing")
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
