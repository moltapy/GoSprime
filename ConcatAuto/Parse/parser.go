package parse

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Args struct {
	WorkPath     *string
	PopList      *string
	ConcatedFile *string
	VcfName      *string
	BcfTool      *string
}

func (args *Args) Parse() {
	directory, err := os.Getwd()
	if err != nil {
		log.Printf("Fail to get current working path: %v", err)
	}
	args.WorkPath = flag.String("w", directory, "Path of the working directory, should be the PARENT DIRECTORY of the subpopulation dirs")
	args.PopList = flag.String("p", "", "Path of the text file with names of subpopulation, per name at each line")
	args.ConcatedFile = flag.String("o", "Summary.vcf.gz", "Path of the concated VCF file, should end with '.vcf.gz' and locate in the dir of each subpopulation in population list file")
	args.VcfName = flag.String("v", "", "Name of the source VCF file, chromosome index should be taken place by '{chrom}', for example: chr{chrom}.vcf.gz")
	args.BcfTool = flag.String("b", "bcftools", "Path of your bcftools, should be the absolute path to your BCFtools except you use conda and have activated bcftools environment")
	flag.Parse()

	if len(os.Args) > 2 {
		if *args.WorkPath == "" {
			log.Fatal("The prama WORKPATH equals a nil value, Please check!")
		}
		if *args.PopList == "" {
			log.Fatal("The prama POPULATION LIST equals a nil value, Please check!")
		}
		if *args.VcfName == "" {
			log.Fatal("The prama VCFFILE NAME equals a nil value, Please check!")
		}
	} else if len(os.Args) == 1 {
		programPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get the program path: %v", err)
		}
		fmt.Printf("Usage of %s:\n", programPath)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
