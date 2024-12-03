package main

import (
	"bufio"
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
}
