package genscore

import "flag"

type Args struct {
	WorkPath     *string
	Threads      *int
	SprimeTool   *string
	PopList      *string
	MapFile      *string
	OutGroupFile *string
	OutFileName  *string
}

func (args *Args) Parse() {
	args.WorkPath = flag.String("w", "", "Path for SAVING OUTPUT SCORE FILES, should be the parent directory of subpopulation directories")
	args.Threads = flag.Int("t", 1, "Num of PARALLLEL GROUPS, the program run analysis on 22 chromosome at same time, cautiously consider total parallel threads")
	args.SprimeTool = flag.String("j", "", "Path of SPRIME JAR file")
	args.PopList = flag.String("p", "", "Path of SUBPPULATION LIST file, should contain one subpopulation name per line")
	args.MapFile = flag.String("m", "", "Path of GENETIC MAP file")
	args.OutGroupFile = flag.String("g", "", "Path of OUTGROUP SAMPLE LIST file, should contain one sample ID per line")
	args.OutFileName = flag.String("o", "SprimeOut_chr{chrom}", "Name of OUTPUT FILE PREFIX, should use {chrom} to take place for chromosome index")
	flag.Parse()
}
