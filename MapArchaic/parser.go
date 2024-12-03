package maparchaic

import "flag"

type Args struct {
	bedMode   *string
	sepChar   *string
	mskFile   *string
	vcfFile   *string
	scoreFile *string
	refTag    *string
	readDepth *string
}

func (args *Args) Parse() {
	args.bedMode = flag.String("bed", "keep", "Bed file for keep/remove/not use")
	args.sepChar = flag.String("s", "\t", "Define the separator in the output file")
	flag.Parse()
}
