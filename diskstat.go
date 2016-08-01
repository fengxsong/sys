package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type DiskStat struct {
	//for more detail refer to https://www.kernel.org/doc/Documentation/iostats.txt
	Major          int
	Minor          int
	Device         string
	ReadCompleted  uint64 //the total number of reads completed successfully.
	ReadMerged     uint64 //Reads and writes which are adjacent to each other may be merged for efficiency.
	ReadSectors    uint64 //the total number of sectors read successfully.
	TimeReading    uint64 //the total number of milliseconds spent by all reads (as measured from __make_request() to end_that_request_last()).
	WriteCompleted uint64 //the total number of writes completed successfully.
	WriteMerged    uint64 //See the description of field `ReadMerged`.
	WriteSectors   uint64 //the total number of sectors written successfully.
	TimeWrite      uint64 //the total number of milliseconds spent by all writes
	IosInProgress  uint64 //I/Os currently in progress
	TimeIos        uint64 //milliseconds spent doing I/Os
	TimeWeighted   uint64
	Ts             time.Time
}

func (ds *DiskStat) String() string {
	return fmt.Sprintf("<Device:%s, Major:%d, Minor:%d...>", ds.Device, ds.Major, ds.Minor)
}

func CheckDiskStats() ([]*DiskStat, error) {
	fn := "/proc/diskstats"
	fb, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ret := make([]*DiskStat, 0)
	rd := bufio.NewReader(bytes.NewBuffer(fb))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		fields := strings.Fields(line)
		if fields[3] == "0" {
			continue
		}
		size := len(fields)
		if size != 14 {
			continue
		}
		ds := &DiskStat{}
		if ds.Major, err = strconv.Atoi(fields[0]); err != nil {
			return nil, err
		}
		if ds.Minor, err = strconv.Atoi(fields[1]); err != nil {
			return nil, err
		}
		ds.Device = fields[2]

		if ds.ReadCompleted, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
			return nil, err
		}
		if ds.ReadMerged, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
			return nil, err
		}
		if ds.ReadSectors, err = strconv.ParseUint(fields[5], 10, 64); err != nil {
			return nil, err
		}
		if ds.TimeReading, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
			return nil, err
		}
		if ds.WriteCompleted, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
			return nil, err
		}
		if ds.WriteMerged, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
			return nil, err
		}
		if ds.WriteSectors, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
			return nil, err
		}
		if ds.TimeWrite, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
			return nil, err
		}
		if ds.IosInProgress, err = strconv.ParseUint(fields[11], 10, 64); err != nil {
			return nil, err
		}
		if ds.TimeIos, err = strconv.ParseUint(fields[12], 10, 64); err != nil {
			return nil, err
		}
		if ds.TimeWeighted, err = strconv.ParseUint(fields[13], 10, 64); err != nil {
			return nil, err
		}
		ds.Ts = time.Now()
		ret = append(ret, ds)
	}
	return ret, nil
}
