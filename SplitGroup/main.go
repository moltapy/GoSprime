package main

import (
	"bufio"
	parse "gosprime/SplitGroup/Parse"
	"io"
	"log"
	"os"
	"sync"
)

var args parse.Args
var set Set
var outGroup []string
var waitGroup, waitSpGroup, waitBcfGroup sync.WaitGroup

func init() {
	args = parse.Args{}
	args.Parse()
	set = make(Set)
}

func main() {
	sampleFile, err := os.Open(*args.SampleFile)
	if err != nil && err != io.EOF {
		log.Fatalf("Problem occurred when opening sample list file, Please check!\nERROR:%s\n", err)
	}
	defer sampleFile.Close()
	reader := bufio.NewReader(sampleFile)
	_, err = reader.ReadBytes('\n')
	if err != nil {
		log.Fatalf("Problem occurred when reading header line of sample list file, Please check!\nERROR:%s\n", err)
	}
	splitGroup(reader, args.OutGroup)

	threadNum := 0
	for subgroup := range set {
		waitSpGroup.Add(1)
		threadNum += 1
		go splitVcfFile(subgroup)
		if threadNum%*args.ParaNum == 0 {
			waitSpGroup.Wait()
		}
	}
	waitSpGroup.Wait()
}
