package parse

import "flag"

type Args struct {
	WorkPath   *string
	SampleFile *string
	OutGroup   *string
	ModernFile *string
	SepChar    *string
	BcfTool    *string
}

func (args *Args) Parse() {
	args.WorkPath = flag.String("wkp", "", "WorkPath of the whole project")
	args.SampleFile = flag.String("smp", "", "Sample file path")
	args.OutGroup = flag.String("oug", "", "Outgroup pop name")
	args.ModernFile = flag.String("mdv", "", "modern vcf file")
	args.SepChar = flag.String("sep", "\t", "separation")
	args.BcfTool = flag.String("btl", "", "BCFTOOLS")
	flag.Parse()
}
