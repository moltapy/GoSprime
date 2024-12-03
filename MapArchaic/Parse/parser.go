package parse

import (
	"flag"
	"log"
	"strconv"
)

type Args struct {
	BedMode   *string
	SepChar   *string
	MskFile   *string
	VcfFile   *string
	ScoreFile *string
	RefTag    *string
	ReadDepth int
}

func (args *Args) Parse() {

	var rdDepth *string

	args.BedMode = flag.String("bed", "keep", "Bed file for keep/remove/not use")
	args.SepChar = flag.String("sep", "\t", "Define the separator in the output file")
	args.MskFile = flag.String("msk", "", "Mask file, only one allowed as the input")
	args.ScoreFile = flag.String("score", "", "Score file from Sprime")
	args.RefTag = flag.String("tag", "", "Tag for the added column")
	rdDepth = flag.String("depth", "", "Add read depth for match(optional)")
	flag.Parse()

	if *args.MskFile == "" {
		log.Fatal("The prama MASK FILE cannot be nil, Please check!")
	} else if *args.ScoreFile == "" {
		log.Fatal("The prama SCORE FILE cannot be nil, Please check!")
	} else if *args.RefTag == "" {
		log.Fatal("The prama REFERENCE TAG cannot be nil, Please check!")
	} else if *rdDepth != "" {
		var err error
		args.ReadDepth, err = strconv.Atoi(*rdDepth)
		if err != nil {
			log.Fatal("The depth received is not an Interger, Please check!")
		}
	} else {
		log.Default()
	}
}
