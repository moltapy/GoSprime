package parse

import "flag"

type Args struct {
	WorkPath   *string
	SampleFile *string
	//SubgroupFile *string
	OutGroup   *string
	ModernFile *string
	SepChar    *string
}

func (args *Args) Parse() {
	args.WorkPath = flag.String("wkp", "", "WorkPath of the whole project")
	args.SampleFile = flag.String("smp", "", "Sample file path")
	//args.SubgroupFile = flag.String("sbg", "", "Generated subgroup file")
	args.OutGroup = flag.String("oug", "", "Outgroup pop name")
	args.ModernFile = flag.String("mdv", "", "modern vcf file")
	args.SepChar = flag.String("sep", "\t", "separation")
	flag.Parse()
}
