package main

import (
	"bufio"
	"bytes"
	parse "gosprime/GenScore/Parse"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const MINSCORE = 15000

var subpops []string
var args parse.Args
var sprimeGroup, chromGroup sync.WaitGroup

func init() {
	args = parse.Args{}
	args.Parse()
}

func main() {

	subpops = readSubPops(*args.PopList)

	for index, subpop := range subpops {
		sprimeGroup.Add(1)
		go runSprime(subpop)
		if index+1%*args.Threads == 0 {
			sprimeGroup.Wait()
		}
	}
	sprimeGroup.Wait()

}

func readSubPops(subPopFilePath string) []string {

	var subPopList []string

	subPopFile, err := os.Open(subPopFilePath)
	if err != nil {
		log.Fatalf("Problem occurred when opening %s: %v", subPopFilePath, err)
	}
	reader := bufio.NewReader(subPopFile)

	log.Printf("Start reading subpopulation infos from %s", subPopFilePath)

	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Problem occurred when reading from buffer: %v", err)
		}
		if strings.TrimRight(string(lineBytes), "\r\n") != "" {
			subPopList = append(subPopList, strings.TrimRight(string(lineBytes), "\r\n"))
		}
		if err == io.EOF {
			log.Printf("End reading subpopulation infos, totally read %d subpopulations", len(subPopList))
			log.Printf("Subpopulations Read: %s", strings.Join(subPopList, " "))
			break
		}
	}
	return subPopList
}

func runSprime(pop string) {
	defer sprimeGroup.Done()

	concatedGenoFile := *args.WorkPath + "/" + pop + *args.GenoPath

	log.Printf("Start sprime process in %s!", pop)
	for chrom := 1; chrom <= 22; chrom++ {
		chromGroup.Add(1)
		go runChromosome(*args.SprimeTool, concatedGenoFile, *args.OutGroupFile, *args.MapFile, *args.OutFileName, pop, chrom)
	}
	chromGroup.Wait()
	log.Printf("Finished sprime calculate in %s!", pop)
}

func runChromosome(sprimeJar, genoFile, outgroupFile, mapFile, outputPrefix, pop string, chrom int) {
	defer chromGroup.Done()

	var stderr bytes.Buffer

	sprimeCommand := exec.Command("java", "-jar", sprimeJar, "gt="+genoFile, "outgroup="+outgroupFile, "map="+mapFile, "out="+outputPrefix, "chrom="+strconv.Itoa(chrom), "minscore="+strconv.Itoa(MINSCORE))

	sprimeCommand.Stderr = &stderr

	log.Printf("Sprime calculate at %d in %s start!", chrom, pop)
	if err := sprimeCommand.Start(); err != nil {
		log.Fatalf("Problem occurred when starting sprime command: %v", err)
	}

	if err := sprimeCommand.Wait(); err != nil {
		log.Printf("Problem occurred when executing sprime calculate: %v, stderr: %s", err, stderr.String())
	}

	log.Printf("Sprime calculate success at chromosome %s in %s!\n", chrom, pop)

}
