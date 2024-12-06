package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func splitGroup(reader *bufio.Reader, outgroup *string) {
	var sampleWriter *os.File
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal("Problem occurred when reading individuals of sample list file,Please check!", err)
		}
		lineStrip := strings.TrimRight(string(lineBytes), "\r\n")
		// avoid the final blank lines
		if lineStrip != "" {
			line := strings.Split(lineStrip, *args.SepChar)

			if !set.exists(line[1]) && line[1] != *outgroup {
				path := *args.WorkPath + "/" + line[1]
				err = os.MkdirAll(path, 0777)
				if err != nil {
					log.Fatal("Problem occurred when making directory of subgroups,Please check!", err)
				}
				sampleList := path + "/sample.txt"
				sampleWriter, err = os.Create(sampleList)
				if err != nil {
					log.Fatal("Problem occurred when create sample file for each group, Please check!", err)
				}
				defer sampleWriter.Close()
				writer := bufio.NewWriter(sampleWriter)
				_, err = writer.WriteString(line[0] + "\n")
				if err != nil {
					log.Fatal("Problem occurred when writing into buffer, Please check!", err)
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
			log.Fatal("Problem occurred when open subgroup sample file, Please check!", err)
		}
		writer := bufio.NewWriter(subGroupFile)
		waitGroup.Add(1)
		go pasteOutGroup(writer, outGroup)
		waitGroup.Wait()
	}

	writeOutGroup(outGroup)
}

func pasteOutGroup(writer *bufio.Writer, outgroup []string) {
	defer waitGroup.Done()
	_, err := writer.WriteString(strings.Join(outgroup, "\n"))
	if err != nil {
		log.Fatal("Problem occurred when paste outgroup,Please check!", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("Problem occurred when flush buffer,Please check!", err)
	}
}

func writeOutGroup(outgroup []string) {
	path := *args.WorkPath + "/outgroup.txt"
	outGroupWriter, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Problem occurred when create outgroup sample list, Please check!", err)
	}
	writer := bufio.NewWriter(outGroupWriter)

	_, err = writer.WriteString(strings.Join(outgroup, "\n"))
	if err != nil {
		log.Fatal("Problem occurred when write outgroup,Please check!", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("Problem occurred when flush buffer,Please check!", err)
	}
}

func splitVcfFile(subPop string) {

	filePrefix := *args.WorkPath + "/" + subPop
	defer waitSpGroup.Done()
	for chrom := 1; chrom <= 22; chrom++ {
		vcfFile := strings.Replace(*args.ModernFile, "{chrom}", strconv.Itoa(chrom), 1)
		outFile := filePrefix + "/" + strings.Replace("chr{chrom}.vcf.gz", "{chrom}", strconv.Itoa(chrom), 1)

		sampleFile := filePrefix + "/sample.txt"
		waitBcfGroup.Add(1)
		go bcftoolExec(*args.BcfTool, subPop, vcfFile, outFile, sampleFile)
	}
	waitBcfGroup.Wait()
}

func bcftoolExec(tool, subPop, vcfFile, outFile, sampleFile string) {
	defer waitBcfGroup.Done()

	viewSamples := exec.Command(tool, "view", "--samples-file", sampleFile, vcfFile)
	cmdReader, err := viewSamples.StdoutPipe()
	if err != nil {
		log.Fatal("Problem occurred when creating StdoutPipe for sample view", err)
	}

	viewSnps := exec.Command(tool, "view", "-c1", "-m2", "-M2", "-v", "snps")
	viewSnps.Stdin = cmdReader
	if err != nil {
		log.Fatal("Problem occurred when setting Stdin for snp view,Please check!", err)
	}
	cmdReader, err = viewSnps.StdoutPipe()
	if err != nil {
		log.Fatal("Problem occurred when creating StdoutPipe for snp view,Please check!", err)
	}
	annotateVcf := exec.Command(tool, "annotate", "-x", "INFO,^FORMAT/GT", "-Oz")

	annotateVcf.Stdin = cmdReader
	if err != nil {
		log.Fatal("Problem occurred when setting Stdin for bcftools annotate,Please check!", err)
	}

	annotatedFile, err := os.Create(outFile)
	if err != nil {
		log.Fatal("Problem occurred when creating output file", err)
	}
	defer annotatedFile.Close()

	annotateVcf.Stdout = annotatedFile

	var viewSamplesStderr, viewSnpsStderr, annotateVcfStderr bytes.Buffer
	viewSamples.Stderr = &viewSamplesStderr
	viewSnps.Stderr = &viewSnpsStderr
	annotateVcf.Stderr = &annotateVcfStderr

	if err := viewSamples.Start(); err != nil {
		log.Fatal("Problem occurred when starting view samples,Please check!", err)
	}
	if err := viewSnps.Start(); err != nil {
		log.Fatal("Problem occurred when starting view snps,Please check!", err)
	}
	if err := annotateVcf.Start(); err != nil {
		log.Fatal("Problem occurred when starting annotate vcfs,Please check!", err)
	}

	if err := viewSamples.Wait(); err != nil {
		log.Fatalf("Problem occurred when processing view samples,err =%v,stderr:%s", err, viewSamplesStderr.String())
	}
	if err := viewSnps.Wait(); err != nil {
		log.Fatalf("Problem occurred when processing view snps, err=%v,stderr:%s", err, viewSnpsStderr.String())
	}
	if err := annotateVcf.Wait(); err != nil {
		log.Fatalf("Problem occurred when processing bcftools annotate, err=%v,stderr:%s", err, annotateVcfStderr.String())
	}
	// take place
	log.Printf("Success Split %s into %s,Extracted %s merge with %s", vcfFile, outFile, subPop, *args.OutGroup)
}
