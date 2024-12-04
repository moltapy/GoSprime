package main

import (
	"bufio"
	parse "gosprime/PreData/Parse"
	"io"
	"log"
	"os"
	"strings"
)

type Set map[string]bool

var args parse.Args
var set Set

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
		panic(err)
	}
	defer sampleFile.Close()
	reader := bufio.NewReader(sampleFile)
	_, err = reader.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	splitGroup(reader, args.OutGroup)
}

func splitGroup(reader *bufio.Reader, outgroup *string) {
	var sampleWriter *os.File
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			// to fill
			log.Fatal("Take place")
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		if lineStrip != "" {
			line := strings.Split(lineStrip, *args.SepChar)
			if !set.exists(line[1]) && line[1] != *outgroup {
				path := *args.WorkPath + "/" + line[1]
				err = os.MkdirAll(path, 0777)
				if err != nil {
					log.Fatal(err)
				}
				sampleList := path + "/sample.txt"
				sampleWriter, err = os.Create(sampleList)
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
					log.Fatalf("Error flushing buffer: %v", err)
				}
				set.add(line[1])
			} else {
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
					log.Fatalf("Error flushing buffer: %v", err)
				}
			}
		}

		if err == io.EOF {
			break
		}
	}
}
