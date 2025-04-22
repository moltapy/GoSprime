package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
}

func (formatter *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var byteBuffer *bytes.Buffer
	if entry.Buffer != nil {
		byteBuffer = entry.Buffer
	} else {
		byteBuffer = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	format := "%s | %s" + strings.Repeat(" ", 9-len(entry.Level.String())) + "| GO_EXECUTABLE - %s\n"
	log := fmt.Sprintf(format, timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	byteBuffer.WriteString(log)
	return byteBuffer.Bytes(), nil
}
