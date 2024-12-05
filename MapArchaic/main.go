package main

import (
	"bufio"
	"compress/gzip"
	parse "gosprime/MapArchaic/Parse"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const NOTCOMP = "notcomp"
const MATCH = "match"
const MISMATCH = "mismatch"

var args parse.Args
var maxLength, validLines int
var start, end string

func init() {
	args = parse.Args{}
	args.Parse()
}

func main() {
	// Open score file
	scoreFile, err := os.Open(*args.ScoreFile)
	if err != nil {
		panic(err)
	}
	defer scoreFile.Close()

	// Read the score file and line by line process
	reader := bufio.NewReader(scoreFile)
	_, err = reader.ReadBytes('\n')
	if err != nil {
		log.Fatal("Exception occured when reading the first line in score file,Please check!")
	}

	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal("Exception occured when reading score file, Please check!", err)
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		line := strings.Split(lineStrip, *args.SepChar)

		if len(line) >= 2 {
			pos, err := strconv.Atoi(line[1])
			if err != nil {
				log.Fatal("The second column of score file cannot be transferred to Interger, Please check!", err)
			}
			validLines += 1
			if maxLength <= pos {
				maxLength = pos
			} else {
				log.Printf("The line which position is %d is smaller than the upper lines, Please make sure your score file is right!", pos)
			}

			if maxLength == 0 {
				log.Fatal("The score file chromosome Length equals 0, Please check your sep or file!")
			}
		} else if len(line) == 1 && line[0] != "" {
			log.Fatal("There is a line donnot split correctly because it's length is just 1!", line)
		} else {
			log.Default()
		}

		if err == io.EOF {
			log.Printf("Finished reading score file: Max Length: %d, Valid lines: %d\n", maxLength, validLines)
			break
		}
	}

	var data = make([][3]byte, maxLength+1)
	var depth = make([]int, maxLength+1)
	for i := range depth {
		depth[i] = -1
	}
	for i := 0; i <= maxLength; i++ {
		if *args.BedMode == "include" {
			data[i][0] = '0'
		} else {
			data[i][0] = '1'
		}
	}

	if *args.BedMode != "" {

		mskFile, err := os.Open(*args.MskFile)
		if err != nil {
			log.Fatalf("There is a problem in opening the file %s,Please check!", *args.MskFile)
		}
		defer mskFile.Close()

		mskLineBytes, err := gzip.NewReader(mskFile)
		if err != nil {
			log.Fatal("There is a problem in read gzip MASK FILE,Please check!", err)
		}
		defer mskLineBytes.Close()

		reader = bufio.NewReader(mskLineBytes)

		for {
			lineBytes, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal("There is a problem in reading the gzip MASK FILE line by line, Please check!", err)
			}
			lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
			line := strings.Split(lineStrip, "\t")

			_, start, end = line[0], line[1], line[2]
			start, err := strconv.Atoi(start)
			if err != nil {
				log.Fatalf("Problem occured when transfer the column 2 in MASK FILE from String to Interger, start position is %d\n", start)
			}
			end, err := strconv.Atoi(end)
			if err != nil {
				log.Fatalf("Problem occured when transfer the column 2 in MASK FILE from String to Interger, start position is %d\n", end)
			}

			for i := start + 1; i <= end; i++ {
				if i > maxLength {
					break
				}
				if *args.BedMode == "include" {
					data[i][0] = '1'
				} else if *args.BedMode == "exclude" {
					data[i][0] = '0'
				}
			}
		}
	}

	archFile, err := os.Open(*args.VcfFile)
	if err != nil {
		log.Fatal("Problem occured when opening the archaic genotype file, Please check!", err)
	}
	defer archFile.Close()

	archLineBytes, err := gzip.NewReader(archFile)
	if err != nil {
		log.Fatal("Problem occured when reading the archaic genotype file, Please check!", err)
	}
	defer archLineBytes.Close()

	reader = bufio.NewReader(archLineBytes)

	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal("Problem occurred when reading the archaic genotype file line by line,Please check!", err)
		}

		if string(lineBytes[1]) != "#" {
			break
		}
	}

	_, err = reader.ReadBytes('\n')
	if err != nil {
		log.Fatal("Problem occurred when skipping header in archaic genotype file, Please check!", err)
	}

	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal("Problem occurred when reading the archaic genotype data line by line,Please check!", err)
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		line := strings.Split(lineStrip, "\t")
		pos, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal("Problem occurred when transfering the second column(position) into Interger, Please check!", err)
		}

		if pos < maxLength {
			if len(line[3]) < 2 && len(line[4]) < 2 {
				if line[9][0] == '0' {
					data[pos][1] = line[3][0]
				}
				if line[9][0] == '1' {
					data[pos][1] = line[4][0]
				}
				if line[9][2] == '0' {
					data[pos][2] = line[3][0]
				}
				if line[9][2] == '1' {
					data[pos][2] = line[4][0]
				}
				depthIndex := strings.Index(line[7], "DP=")
				if depthIndex != -1 {
					depthIndex += 3
					depth[pos], err = strconv.Atoi(line[7][depthIndex:])
					if err != nil {
						log.Fatal("Problem occurred when transfering read depth into Interger,Please check!")
					}
				} else {
					depth[pos] = 1
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	scoreFile, err = os.Open(*args.ScoreFile)
	if err != nil {
		panic(err)
	}
	defer scoreFile.Close()

	outFile, err := os.Open(*args.OutFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	reader = bufio.NewReader(scoreFile)
	writer := bufio.NewWriter(outFile)

	header, err := reader.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	var newHeader string
	if *args.ReadDepth == "true" {
		newHeader = string(header) + *args.SepChar + *args.RefTag + *args.SepChar + *args.RefTag + "_DP" + "\n"
	} else {
		newHeader = string(header) + *args.SepChar + *args.RefTag + "\n"
	}
	_, err = writer.WriteString(newHeader)
	if err != nil {
		panic(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("Problem occurred when flushing the header in buffer into file,Please check!", err)
	}

	var snps = make([]string, 2)
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal("Problem occurred when reading the archaic genotype data line by line,Please check!", err)
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		line := strings.Split(lineStrip, "\t")
		pos, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal("Problem occurred when transfering the second column(position) into Interger, Please check!", err)
		}
		snps[0], snps[1] = line[3], line[4]

		allele, err := strconv.Atoi(line[6])
		if err != nil {
			log.Fatal("Problem occurred when transfer the seventh column(allele) into Interger,Please check!", err)
		}

		var writeLine = string(lineBytes)
		if data[pos][0] == '0' || depth[pos] < 0 {
			writeLine += *args.SepChar + NOTCOMP
			if *args.ReadDepth == "true" {
				writeLine += *args.SepChar + string(rune(depth[pos]))
			}
			writeLine += "\n"
		} else {
			if snps[allele] == string(data[pos][1]) || snps[allele] == string(data[pos][2]) {
				writeLine += *args.SepChar + MATCH
			} else {
				writeLine += *args.SepChar + MISMATCH
			}
			if *args.ReadDepth == "true" {
				writeLine += *args.SepChar + string(rune(depth[pos]))
			}
			writeLine += "\n"
		}
		_, err = writer.WriteString(writeLine)
		if err != nil {
			log.Fatal("Problem occurred when write lines into buffer,Please check!", err)
		}

		if err == io.EOF {
			log.Printf("Mapping %s Finished", *args.RefTag)
			break
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal("Problem occurred when flushing lines from buffer into file,Please check!", err)
	}
}
