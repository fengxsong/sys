package sys

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SysUptime struct {
	Uptime   float64 //the total number of seconds the system has been up
	IdleTime float64 //the sum of how much time each core has spent idle, in seconds
	Days     int64
	Hours    int64
	Mins     int64
}

func (s *SysUptime) String() string {
	loadavg, _ := CheckLoadavg()
	return fmt.Sprintf("%s up %d days, %d:%d, load average: %.2f, %.2f, %.2f",
		time.Now().Format("15:04:05"),
		s.Days,
		s.Hours,
		s.Mins,
		loadavg.Avg1min,
		loadavg.Avg5min,
		loadavg.Avg15min)
}

func Uptime() (*SysUptime, error) {
	s, err := Read2String("/proc/uptime")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(s)
	if len(fields) < 2 {
		return nil, fmt.Errorf("/proc/uptime format not supported")
	}

	su := &SysUptime{}
	su.Uptime, _ = strconv.ParseFloat(fields[0], 64)
	su.IdleTime, _ = strconv.ParseFloat(fields[1], 64)
	minsTotal := su.Uptime / 60.0
	hoursTotal := minsTotal / 60.0
	su.Days = int64(hoursTotal / 24.0)
	su.Hours = int64(hoursTotal) - su.Days*24
	su.Mins = int64(minsTotal) - su.Days*24*60 - su.Hours*60
	return su, nil
}
