package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func splitVcfFile(subPop string) {

	log.Printf("Start Processing bcftools concat: %s with %s", subPop, *args.OutGroup)

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
