package parse

import "flag"

type Args struct {
	WorkPath     *string
	PopList      *string
	ConcatedFile *string
	VcfName      *string
	BcfTool      *string
}

func (args *Args) Parse() {
	args.WorkPath = flag.String("w", "", "father dir of the poplations")
	args.PopList = flag.String("p", "", "path of population lists")
	args.ConcatedFile = flag.String("o", "", "path of concated file")
	args.VcfName = flag.String("v", "", "vcf file name")
	args.BcfTool = flag.String("b", "", "bcftool path")

	flag.Parse()
}
