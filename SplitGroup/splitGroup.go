package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func splitGroup(reader *bufio.Reader, outgroup *string) {
	var sampleWriter *os.File
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Problem occurred when reading individuals of sample list file: %v", err)
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		// avoid the final blank lines
		if lineStrip != "" {
			line := strings.Split(lineStrip, *args.SepChar)

			if !set.exists(line[1]) && line[1] != *outgroup {
				path := *args.WorkPath + "/" + line[1]
				err = os.MkdirAll(path, 0777)
				if err != nil {
					log.Fatalf("Problem occurred when making directory of subgroups: %v", err)
				}
				sampleList := path + "/sample.txt"
				sampleWriter, err = os.Create(sampleList)
				if err != nil {
					log.Fatalf("Problem occurred when create sample file for each group: %v", err)
				}
				defer sampleWriter.Close()
				writer := bufio.NewWriter(sampleWriter)
				_, err = writer.WriteString(line[0] + "\n")
				if err != nil {
					log.Fatalf("Problem occurred when writing sample ID into buffer: %v", err)
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
					log.Fatalf("Problem occurred when opening %v for writing: %v", sampleList, err)
				}
				defer sampleWriter.Close()
				writer := bufio.NewWriter(sampleWriter)
				_, err = writer.WriteString(line[0] + "\n")
				if err != nil {
					log.Fatalf("Problem occurred when writing %v into buffer: %v", line[0], err)
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
			log.Fatalf("Problem occurred when opening subgroup sample file %v: %v", subGroupList, err)
		}
		writer := bufio.NewWriter(subGroupFile)
		waitGroup.Add(1)
		go pasteOutGroup(writer, outGroup)
	}
	waitGroup.Wait()

	writeOutGroup(outGroup)
}

func pasteOutGroup(writer *bufio.Writer, outgroup []string) {
	defer waitGroup.Done()
	_, err := writer.WriteString(strings.Join(outgroup, "\n"))
	if err != nil {
		log.Fatalf("Problem occurred when pasting outgroup samples: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Problem occurred when flushing buffer into file: %v", err)
	}
}

func writeOutGroup(outgroup []string) {
	path := *args.WorkPath + "/outgroup.txt"
	outGroupWriter, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Problem occurred when creating outgroup sample list file %v: %v", path, err)
	}
	writer := bufio.NewWriter(outGroupWriter)

	_, err = writer.WriteString(strings.Join(outgroup, "\n"))
	if err != nil {
		log.Fatalf("Problem occurred when writing outgroup samples into buffer: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Problem occurred when flushing buffer into file: %v", err)
	}
}
