package main

import (
	"bufio"
	parse "gosprime/SplitGroup/Parse"
	"io"
	"log"
	"os"
	"sync"
)

type Set map[string]bool

var args parse.Args
var set Set
var outGroup []string
var waitGroup, waitSpGroup, waitBcfGroup sync.WaitGroup

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

	var threadNum = 0
	for subgroup := range set {
		waitSpGroup.Add(1)
		threadNum += 1
		go splitVcfFile(subgroup)
		if threadNum%*args.ParaNum == 0 {
			log.Printf("Start Processing bcftools concat: %s with %s", subgroup, *args.OutGroup)
			waitSpGroup.Wait()
		}
	}
	waitSpGroup.Wait()
}
