package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Memory struct {
	Buffers   uint64
	Cached    uint64
	MemTotal  uint64
	MemFree   uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
}

func (m *Memory) String() string {
	return fmt.Sprintf("<MemTotal:%dkB, MemFree:%dkB, SwapUsedPercent:%.2f >", m.MemTotal, m.MemFree, float64(m.SwapUsed)/float64(m.SwapTotal)*100.0)
}

var M = map[string]bool{
	"MemTotal:":  true,
	"MemFree:":   true,
	"Buffers:":   true,
	"Cached:":    true,
	"SwapTotal:": true,
	"SwapFree:":  true,
}

func MemInfo() (*Memory, error) {
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	m := &Memory{}
	rd := bufio.NewReader(bytes.NewBuffer(b))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		fields := strings.Fields(line)
		fieldName := fields[0]
		if _, ok := M[fieldName]; ok && len(fields) == 3 {
			v, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				continue
			}
			switch fieldName {
			case "MemTotal:":
				m.MemTotal = v
			case "MemFree:":
				m.MemFree = v
			case "Buffers:":
				m.Buffers = v
			case "Cached:":
				m.Cached = v
			case "SwapTotal:":
				m.SwapTotal = v
			case "SwapFree:":
				m.SwapFree = v
			}
		}
	}
	m.SwapUsed = m.SwapTotal - m.SwapFree
	return m, nil
}
