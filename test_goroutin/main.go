package main

import (
	"flag"
	"fmt"
	"github.com/fengxsong/sys"
	"runtime"
	"strconv"
	"time"
)

func Test() ([]*sys.PidStatus, error) {
	dirs, err := sys.WalkOnlyDirs("/tmp/test_data")
	if err != nil {
		return nil, err
	}

	sz := len(dirs)
	if sz == 0 {
		return nil, nil
	}

	procs := make([]*sys.PidStatus, 0)

	for i := 0; i < sz; i++ {
		pid, e := strconv.Atoi(dirs[i])
		if e != nil {
			continue
		}
		statusFile := fmt.Sprintf("/tmp/test_data/%d/status", pid)
		if !sys.IsExist(statusFile) {
			continue
		}
		ps, err := sys.GetPidStatus(statusFile)
		if err != nil {
			continue
		}

		ps.Pid = pid
		procs = append(procs, ps)
	}
	return procs, nil
}

func TestGoroutin() ([]*sys.PidStatus, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	dirs, err := sys.WalkOnlyDirs("/tmp/test_data")
	if err != nil {
		return nil, err
	}

	sz := len(dirs)
	if sz == 0 {
		return nil, nil
	}

	procs := make([]*sys.PidStatus, 0)
	pCh := make(chan *sys.PidStatus, sz)

	for i := 0; i < sz; i++ {
		go func(index int) {
			pid, e := strconv.Atoi(dirs[index])
			if e != nil {
				return
			}
			statusFile := fmt.Sprintf("/tmp/test_data/%d/status", pid)
			if !sys.IsExist(statusFile) {
				return
			}
			ps, err := sys.GetPidStatus(statusFile)
			if err != nil {
				return
			}

			ps.Pid = pid
			pCh <- ps
		}(i)
	}
	for i := 0; i < sz; i++ {
		procs = append(procs, <-pCh)
	}
	return procs, nil

}

func main() {
	vv := flag.Bool("vv", false, "if true, use goroutin")
	flag.Parse()
	var test01 []*sys.PidStatus
	var err error
	startTime := time.Now()
	if *vv {
		test01, err = TestGoroutin()
	} else {
		test01, err = Test()
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now().Sub(startTime))
	fmt.Println(len(test01))
}
