package sys

import (
	"fmt"
	"strconv"
	"strings"
)

type Loadavg struct {
	Avg1min  float64
	Avg5min  float64
	Avg15min float64
}

func (l *Loadavg) String() string {
	return fmt.Sprintf("<1min:%f, 5min:%f, 15min:%f>", l.Avg1min, l.Avg5min, l.Avg15min)
}

func CheckLoadavg() (*Loadavg, error) {
	s, err := Read2TrimString("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	loadAvg := &Loadavg{}
	fields := strings.Fields(s)
	if loadAvg.Avg1min, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg5min, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg15min, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return nil, err
	}
	return loadAvg, nil
}
