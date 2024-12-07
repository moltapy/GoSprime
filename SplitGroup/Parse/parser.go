package parse

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Args struct {
	WorkPath   *string
	SampleFile *string
	OutGroup   *string
	ModernFile *string
	SepChar    *string
	BcfTool    *string
	ParaNum    *int
}

func (args *Args) Parse() {

	directory, err := os.Getwd()
	if err != nil {
		log.Printf("Fail to get current working path, Please check! \nError: %s\n", err)
	}

	args.WorkPath = flag.String("w", directory, "Path for saving VCF files per subgroup")
	args.SampleFile = flag.String("s", "", "Path of the text file saving all individuals and their populations")
	args.OutGroup = flag.String("o", "", "Name of the outgroup, selected outgroup will merge into other subgroups and not be splited separately")
	args.ModernFile = flag.String("m", "", "Path of VCF files containing gene information of modern humans, chromosome index should be taken place by '{chrom}' and should use '.vcf.gz' format")
	args.SepChar = flag.String("c", "\t", "Separator of the sample list file")
	args.BcfTool = flag.String("b", "bcftools", "Path of your bcftools, should be the absolute path to your BCFtools except you use conda and have activated bcftools environment")
	args.ParaNum = flag.Int("p", 4, "Number of the parallel threads when dealing with subpopulations, you will open p * 22(the number of chromosomes) threads in total at one time")
	flag.Parse()

	if len(os.Args) > 2 {
		if *args.WorkPath == "" {
			log.Fatal("The prama WORKPATH equals a nil value, Please check!")
		}
		if *args.SampleFile == "" {
			log.Fatal("The prama SAMPLE FILE equals a nil value, Please check!")
		}
		if *args.ModernFile == "" {
			log.Fatal("The prama MODERN HUMAN VCF FILE equals a nil value, Please check!")
		} else if !strings.Contains(*args.ModernFile, "{chrom}") {
			log.Fatal("The prama MODERN HUMAN VCF FILE donnot contain '{chrom}', Please check!")
		}
	} else if len(os.Args) == 1 {
		programPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get the program path, Please check!\n ERROR: %s\n", err)
		}
		fmt.Printf("Usage of %s:\n", programPath)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
