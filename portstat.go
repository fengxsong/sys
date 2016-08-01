package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ListeningPorts() ([]int64, error) {
	return listeningPorts("ss", true, "-t", "-l", "-n")
}

func UdpPorts() ([]int64, error) {
	return listeningPorts("ss", true, "-u", "-a", "-n")
}

func listeningPorts(name string, sort bool, args ...string) ([]int64, error) {
	b, _, err := ExecCmdBytes(name, args...)
	if err != nil {
		return nil, err
	}

	ports := []int64{}
	rd := bufio.NewReader(bytes.NewBuffer(b))
	//ignore the first line
	line, err := rd.ReadString('\n')
	if err != nil {
		return nil, err
	}

	for {
		line, err = rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		fields := strings.Fields(line)
		fieldsLen := len(fields)
		if fieldsLen != 4 && fieldsLen != 5 {
			return nil, fmt.Errorf("output of command %s error", name)
		}

		portColumnIndex := 2
		if fieldsLen == 5 {
			portColumnIndex = 3
		}

		location := strings.LastIndex(fields[portColumnIndex], ":")
		port := fields[portColumnIndex][location+1:]

		if p, e := strconv.ParseInt(port, 10, 64); e != nil {
			return ports, fmt.Errorf("parse port to int64 fail: %s", e.Error())
		} else {
			ports = append(ports, p)
		}
	}
	return Int64Set(ports, sort), nil
}
