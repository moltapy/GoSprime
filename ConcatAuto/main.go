package main

import (
	"bufio"
	parse "gosprime/ConcatAuto/Parse"
	"io"
	"log"
	"os"
	"os/exec"
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
		waitGroup.Wait()
	}
}

func readPops(path string) []string {
	var popGroups []string
	popList, err := os.Open(path)
	if err != nil {
		log.Fatal("Problem occurred when open population list file, Please check!", err)
	}
	reader := bufio.NewReader(popList)
	for {
		lineBytes, err := reader.ReadBytes('\n')
		line := strings.TrimRight(string(lineBytes), "\r\n")
		if err != nil && err != io.EOF {
			log.Fatal(err)
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

	var fileList string
	for chrom := 1; chrom <= 22; chrom++ {
		fileList += *args.WorkPath + "/" + pop + "/" + strings.Replace(*args.VcfName, "{chrom}", string(rune(chrom)), 1) + " "
	}

	concatCommand := exec.Command(*args.BcfTool, "concat", "--file-list", strings.TrimRight(fileList, " "),
		"--naive-force", "--output-type", "z", "--output", *args.ConcatedFile)

	if err := concatCommand.Start(); err != nil {
		log.Fatal(err)
	}

	if err := concatCommand.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Println("Success!")

}
