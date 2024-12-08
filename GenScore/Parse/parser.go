package parse

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Args struct {
	WorkPath     *string
	Threads      *int
	SprimeTool   *string
	PopList      *string
	MapFile      *string
	OutGroupFile *string
	GenoPath     *string
	OutFileName  *string
}

func (args *Args) Parse() {

	directory, err := os.Getwd()
	if err != nil {
		log.Printf("Fail to get current working path,Please check!\nError: %s\n", err)
	}
	args.WorkPath = flag.String("w", directory, "Path for SAVING OUTPUT SCORE FILES, should be the parent directory of subpopulation directories")
	args.Threads = flag.Int("t", 1, "Num of PARALLLEL GROUPS, the program run analysis on 22 chromosome at same time, cautiously consider total parallel threads")
	args.SprimeTool = flag.String("j", "", "Path of SPRIME JAR file")
	args.PopList = flag.String("p", "", "Path of SUBPPULATION LIST file, should contain one subpopulation name per line")
	args.MapFile = flag.String("m", "", "Path of GENETIC MAP file")
	args.OutGroupFile = flag.String("f", "", "Path of OUTGROUP SAMPLE LIST file, should contain one sample ID per line")
	args.GenoPath = flag.String("g", "", "Path of BCFTOOLS CONCATED VCF file, should contain all variants in 22 chromosomes and start from the subdirectory of the subpopulation directory")
	args.OutFileName = flag.String("o", "SprimeOut_chr{chrom}", "Name of OUTPUT FILE PREFIX, should use {chrom} to take place for chromosome index")
	flag.Parse()

	if len(os.Args) > 2 {
		if *args.WorkPath == "" {
			log.Fatal("The prama WORKPATH equals a nil value, Please check!")
		}
		if *args.Threads == 0 {
			log.Fatal("The prama THREAD NUMBBERS equals a nil value, Please check!")
		}
		if *args.SprimeTool == "" {
			log.Fatal("The prama SPRIME JAR equals a nil value, Please check!")
		}
		if *args.PopList == "" {
			log.Fatal("The prama POPULATION LIST FILE equals a nil value, Please check!")
		}
		if *args.MapFile == "" {
			log.Fatal("The prama GENETIC MAP FILE equals a nil value, Please check!")
		}
		if *args.OutGroupFile == "" {
			log.Fatal("The prama OUTGROUP SAMPLE LIST FILE equals a nil value, Please check!")
		}
		if *args.GenoPath == "" {
			log.Fatal("The prama BCFTOOLS CONCATED VCF FILE PATH equals a nil value, Please check!")
		}
		if *args.OutFileName == "" {
			log.Fatal("The prama OUTPUT FILE NAME PREFIX equals a nil value, Please check!")
		} else if !strings.Contains(*args.OutFileName, "{chrom}") {
			log.Fatal("The prama OUTPUT FILE NAME PREFIX donnot contain '{chrom}', Please check!")
		}
	} else if len(os.Args) == 1 {
		programPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get the program path,Please check!\n ERROR: %s\n", err)
		}
		fmt.Printf("Usage of %s:\n", programPath)
		flag.PrintDefaults()
		os.Exit(1)
	}

}
