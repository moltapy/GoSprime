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

func main() {
	// Parser part
	var args = parse.Args{}
	args.Parse()

	var maxLength, validLines int

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

	var data = make([][3]string, maxLength+1)
	var start, end string
	for i := 0; i <= maxLength; i++ {
		if *args.BedMode == "include" {
			data[i][0] = "0"
		} else {
			data[i][0] = "1"
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
					data[i][0] = "1"
				} else if *args.BedMode == "exclude" {
					data[i][0] = "0"
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
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		line := strings.Split(lineStrip, "\t")

		if line[1]
	}
}
