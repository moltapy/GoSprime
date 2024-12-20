package main

import (
	"bufio"
	"bytes"
	parse "gosprime/ConcatAuto/Parse"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var args = parse.Args{}
var subGroups []string
var waitGroup sync.WaitGroup

func init() {
	args.Parse()
}

func main() {
	subGroups = readPops(*args.PopList)
	for _, subPop := range subGroups {
		waitGroup.Add(1)
		go concatAutos(subPop)
	}
	waitGroup.Wait()
}

func readPops(path string) []string {
	var popGroups []string
	popList, err := os.Open(path)
	if err != nil {
		log.Fatalf("Problem occurred when opening population list file: %v", err)
	}
	reader := bufio.NewReader(popList)
	for {
		lineBytes, err := reader.ReadBytes('\n')
		line := strings.TrimRight(string(lineBytes), "\r\n")
		if err != nil && err != io.EOF {
			log.Fatalf("Problem occurred when reading from buffer: %v", err)
		}
		if line != "" {
			popGroups = append(popGroups, line)
		}
		if err == io.EOF {
			break
		}
	}
	return popGroups
}

func concatAutos(pop string) {
	defer waitGroup.Done()

	var pramaList = []string{"concat"}
	pathPrefix := *args.WorkPath + "/" + pop + "/"
	for chrom := 1; chrom <= 22; chrom++ {
		pramaList = append(pramaList, pathPrefix+strings.Replace(*args.VcfName, "{chrom}", strconv.Itoa(chrom), 1))
	}
	pramaList = append(pramaList, "--naive-force", "--output-type", "z", "--output", pathPrefix+*args.ConcatedFile)
	concatCommand := exec.Command(*args.BcfTool, pramaList...)

	var stderr bytes.Buffer
	concatCommand.Stderr = &stderr

	log.Printf("Bcftools concat VCF files in %s start!\n", pop)
	if err := concatCommand.Start(); err != nil {
		log.Fatalf("Problem occurred when starting concat command: %v", err)
	}

	if err := concatCommand.Wait(); err != nil {
		log.Printf("Problem occurred when executing concat command: %v,stderr: %s", err, stderr.String())
	} else {
		log.Printf("Bcftools concat success in %s!\n", pop)
	}

}
