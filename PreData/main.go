package main

import (
	"bufio"
	parse "gosprime/PreData/Parse"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Set map[string]bool

var args parse.Args
var set Set
var outGroup []string
var waitGroup, waitSpGroup sync.WaitGroup

func init() {
	args = parse.Args{}
	args.Parse()
	set = make(Set)
}

func (set Set) exists(val string) bool {
	_, ok := set[val]
	if ok {
		return true
	} else {
		return false
	}
}

func (set Set) add(val string) {
	set[val] = true
}

func main() {
	sampleFile, err := os.Open(*args.SampleFile)
	if err != nil && err != io.EOF {
		log.Fatal("Problem occurred when opening sample list file, Please check!", err)
	}
	defer sampleFile.Close()
	reader := bufio.NewReader(sampleFile)
	_, err = reader.ReadBytes('\n')
	if err != nil {
		log.Fatal("Problem occurred when reading header line of sample list file, Please check!", err)
	}
	splitGroup(reader, args.OutGroup)
}

func splitGroup(reader *bufio.Reader, outgroup *string) {
	var sampleWriter *os.File
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal("Problem occurred when reading individuals of sample list file,Please check!", err)
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		// avoid the final blank lines
		if lineStrip != "" {
			line := strings.Split(lineStrip, *args.SepChar)

			if !set.exists(line[1]) && line[1] != *outgroup {
				path := *args.WorkPath + "/" + line[1]
				err = os.MkdirAll(path, 0777)
				if err != nil {
					log.Fatal("Problem occurred when making directory of subgroups,Please check!", err)
				}
				sampleList := path + "/sample.txt"
				sampleWriter, err = os.Create(sampleList)
				if err != nil {
					log.Fatal("Problem occurred when create sample file for each group, Please check!", err)
				}
				defer sampleWriter.Close()
				writer := bufio.NewWriter(sampleWriter)
				_, err = writer.WriteString(line[0] + "\n")
				if err != nil {
					log.Fatal("Problem occurred when writing into buffer, Please check!", err)
				}
				err = writer.Flush()
				if err != nil {
					log.Fatalf("Problem occurred when flushing buffer into file: %v", err)
				}
				set.add(line[1])
			} else if set.exists(line[1]) && line[1] != *outgroup {
				path := *args.WorkPath + "/" + line[1]
				sampleList := path + "/sample.txt"
				sampleWriter, err = os.OpenFile(sampleList, os.O_APPEND|os.O_RDWR, 0666)
				if err != nil {
					panic(err)
				}
				defer sampleWriter.Close()
				writer := bufio.NewWriter(sampleWriter)
				_, err = writer.WriteString(line[0] + "\n")
				if err != nil {
					panic(err)
				}

				err = writer.Flush()
				if err != nil {
					log.Fatalf("Problem occurred when flushing buffer into file: %v", err)
				}
			} else if line[1] == *outgroup {
				outGroup = append(outGroup, line[0])
			}
		}
		if err == io.EOF {
			break
		}
	}
	// use goroutine to quick add all outGroup into those files
	for subGroup := range set {
		subGroupList := *args.WorkPath + "/" + subGroup + "/sample.txt"
		subGroupFile, err := os.OpenFile(subGroupList, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal("Problem occurred when open subgroup sample file, Please check!", err)
		}
		writer := bufio.NewWriter(subGroupFile)
		waitGroup.Add(1)
		go pasteOutGroup(writer, outGroup)
		waitGroup.Wait()
	}

	writeOutGroup(outGroup)
}

func pasteOutGroup(writer *bufio.Writer, outgroup []string) {
	defer waitGroup.Done()
	_, err := writer.WriteString(strings.Join(outgroup, "\n"))
	if err != nil {
		log.Fatal("Problem occurred when paste outgroup,Please check!", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("Problem occurred when flush buffer,Please check!", err)
	}
}

func writeOutGroup(outgroup []string) {
	path := *args.WorkPath + "/outgroup.txt"
	outGroupWriter, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Problem occurred when create outgroup sample list, Please check!", err)
	}
	writer := bufio.NewWriter(outGroupWriter)

	_, err = writer.WriteString(strings.Join(outgroup, "\n"))
	if err != nil {
		log.Fatal("Problem occurred when write outgroup,Please check!", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("Problem occurred when flush buffer,Please check!", err)
	}
}

func splitVcfFile(subPop string) {
	for chrom := 1; chrom <= 22; chrom++ {
		vcfFile := strings.Replace(*args.ModernFile, "{chrom}", string(chrom), 1)
		outFile := strings.Replace("chr{chrom}.vcf.gz", "{chrom}", string(chrom), 1)

		sampleFile := *args.WorkPath + "/" + subPop + "/sample.txt"
		waitSpGroup.Add(1)
		go bcftoolExec(vcfFile, outFile, sampleFile)

	}
}

func bcftoolExec(vcfFile, outFile, sampleFile string) {
	defer waitSpGroup.Done()

}
