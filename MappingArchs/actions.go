package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var actions = func(ctx context.Context, c *cli.Command) error {

	// genoInfos用于替换dataArray和deptharray,[0,1]或者[1,1]初始化
	genoInfos := make([][2]int, maxpos+1)

	if isreverse {
		logrus.Infof("Reverse mode on, include mask: %s", maskpath)
		if maxpos < 1000000 {
			for i := 0; i <= maxpos; i++ {
				genoInfos[i] = [2]int{Inactive, Defaultdepth}
			}
		} else {
			for i := 0; i <= maxpos; i++ {
				go initGenoInfo(genoInfos, i, Inactive, Defaultdepth)
			}
		}

	} else {
		logrus.Infof("Reverse mode off, exclude mask: %s", maskpath)
		if maxpos < 1000000 {
			for i := 0; i <= maxpos; i++ {
				genoInfos[i] = [2]int{Active, Defaultdepth}
			}
		} else {
			for i := 0; i <= maxpos; i++ {
				go initGenoInfo(genoInfos, i, Active, Defaultdepth)
			}
		}

	}

	if isdepth {
		logrus.Infof("Depth mode on, output will contain read-depth columns")
	}

	sep = strings.Replace(sep, "\\t", "\t", -1)

	if maskpath != "" {
		maskhandler, err := os.Open(maskpath)
		if err != nil {
			return fmt.Errorf("Opening file: %s error, reason: %v", maskpath, err)
		}
		defer maskhandler.Close()

		maskreader, err := gzip.NewReader(maskhandler)
		if err != nil {
			return fmt.Errorf("Reading file: %s error, reason: %v", maskpath, err)
		}
		defer maskreader.Close()

		bufreader := bufio.NewReader(maskreader)

		for {
			line, _, err := bufreader.ReadLine()
			if err != nil && err != io.EOF {
				return fmt.Errorf("Reading line from buffer error, reason: %v", err)
			}

			if err == io.EOF {
				logrus.Infof("Complete reading file: %s", maskpath)
				break
			}

			contents := strings.Split(string(line), "\t")
			start, err := strconv.Atoi(contents[1])
			if err != nil {
				return fmt.Errorf("Start pos %s is non-int, reason: %v", contents[1], err)
			}
			end, err := strconv.Atoi(contents[2])
			if err != nil {
				return fmt.Errorf("End pos %s is non-int, reason: %v", contents[2], err)
			}
			for i := start + 1; i <= end && i <= maxpos; i++ {
				genoInfos[i][0] = (genoInfos[i][0] ^ Active) << MaskSite
			}
		}
	}

	archaichandler, err := os.Open(archvcfpath)
	if err != nil {
		return fmt.Errorf("Opening file: %s error, reason: %v", archvcfpath, err)
	}
	defer archaichandler.Close()

	archaicreader, err := gzip.NewReader(archaichandler)
	if err != nil {
		return fmt.Errorf("Reading file: %s error, reason: %v", archvcfpath, err)
	}
	defer archaicreader.Close()

	archbufreader := bufio.NewReader(archaicreader)

	for {
		line, _, err := archbufreader.ReadLine()
		if err != nil && err != io.EOF {
			return fmt.Errorf("Reading line from buffer error, reason: %v", err)
		}
		if err == io.EOF {
			return fmt.Errorf("File: %s contains no informative lines", archvcfpath)
		}

		if !strings.HasPrefix(string(line), "##") {
			break
		}
	}

	for {
		line, _, err := archbufreader.ReadLine()
		if err != nil && err != io.EOF {
			return fmt.Errorf("Reading line from buffer error, reason: %v", err)
		}
		if err == io.EOF {
			logrus.Infof("Complete reading file: %s", archvcfpath)
			break
		}

		lines := strings.Split(string(line), "\t")
		position, err := strconv.Atoi(lines[1])
		if err != nil {
			return fmt.Errorf("Second column in file: %s contains non-int value: %s", archvcfpath, lines[1])
		}

		if position > maxpos {
			break
		}

		refAllele, altAllele, depth, genotype := lines[3], lines[4], lines[7], lines[9]

		if len(refAllele) < 2 && len(altAllele) < 2 {

			if int(genotype[0]-48) == 0 {
				genoInfos[position][0] = (genoInfos[position][0] | int(refAllele[0]-48)<<LeftSite)
			} else if int(genotype[0]-48) == 1 {
				genoInfos[position][0] = (genoInfos[position][0] | int(altAllele[0]-48)<<LeftSite)
			}

			if int(genotype[2]-48) == 0 {
				genoInfos[position][0] = (genoInfos[position][0] | int(refAllele[0]-48)<<RightSite)
			} else if int(genotype[2]-48) == 1 {
				genoInfos[position][0] = (genoInfos[position][0] | int(altAllele[0]-48)<<RightSite)
			}

			re := regexp.MustCompile(`\d+`)
			depthval := re.FindString(string(depth))
			if depthval == "" {
				if !depthTag {
					logrus.Warningf("Depth value not found in %s, continue with 1, stop and check if needed, hint: program use first int value of INFO column as depth", string(depth))
					depthTag = true
				}
			} else {
				depthInt, _ := strconv.Atoi(depthval)
				genoInfos[position][1] = depthInt
			}
		}
	}

	atomicRewrite(scorepath, genoInfos)

	logrus.Infof("Complete mapping archaic, column: %s, score file: %s", arrayname, scorepath)

	return nil
}

func atomicRewrite(filename string, genoInfos [][2]int) error {
	tmpFile := filename + ".tmp"

	input, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer output.Close()
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	count := 0
	for scanner.Scan() {
		count++
		line := scanner.Text()

		var modified string

		if count == 1 {
			modified = line + sep + arrayname
			if isdepth {
				modified += sep + arrayname + "_depth"
			}
		} else {
			contents := strings.Split(line, "\t")

			pos, err := strconv.Atoi(contents[1])
			if err != nil {
				return fmt.Errorf("Second column of score file contains non-int value: %s", contents[1])
			}

			k, err := strconv.Atoi(contents[6])
			if err != nil {
				return fmt.Errorf("Seventh column of score file contains non-int value: %s", contents[6])
			}
			var snp int
			switch k {
			case 0:
				snp = int([]byte(contents[3])[0] - 48)
			case 1:
				snp = int([]byte(contents[4])[0] - 48)
			default:
				return fmt.Errorf("Seventh column value: %s of score file cannot match snp", contents[6])
			}

			modified = line + processPosition(pos, snp, isdepth, genoInfos)

		}

		if _, err := writer.WriteString(modified + "\n"); err != nil {
			os.Remove(tmpFile)
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		os.Remove(tmpFile)
		return err
	}
	if err := writer.Flush(); err != nil {
		os.Remove(tmpFile)
		return err
	}

	return os.Rename(tmpFile, filename)
}

func processPosition(pos, snp int, depthOptions bool, genoInfos [][2]int) string {

	resStr := ""

	if genoInfos[pos][0]&(Active<<MaskSite) == 0 || genoInfos[pos][1] < 0 {
		resStr += sep + "notcomp"
	} else {
		if genoInfos[pos][0]&(snp<<LeftSite) == 0 || genoInfos[pos][0]&(snp<<RightSite) == 0 {
			resStr += sep + "match"
		} else {
			resStr += sep + "mismatch"
		}
	}
	if depthOptions {
		resStr += sep + strconv.Itoa(genoInfos[pos][1])
	}
	return resStr
}

func initGenoInfo(array [][2]int, index, status, depth int) {
	array[index] = [2]int{status, depth}
}
