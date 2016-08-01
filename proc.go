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

type PidStatus struct {
	Pid   int
	Name  string
	State string
}

type Proc struct {
	Ps      *PidStatus
	Cmdline string
}

func (p *Proc) String() string {
	return fmt.Sprintf("<Pid:%d, Name:%s, State:%s, Cmdline:%s>", p.Ps.Pid, p.Ps.Name, p.Ps.State, p.Cmdline)
}

func getPidStatus(fn string) (*PidStatus, error) {
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ps := &PidStatus{}
	rd := bufio.NewReader(bytes.NewBuffer(f))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		colonIndex := strings.Index(line, ":")
		switch title := strings.TrimSpace(line[0:colonIndex]); title {
		case "Name":
			ps.Name = strings.TrimSpace(line[colonIndex+1:])
		case "State":
			ps.State = strings.TrimSpace(line[colonIndex+1:])
		}
	}
	return ps, nil
}

func AllProcs() ([]*Proc, error) {
	dirs, err := WalkOnlyDirs("/proc")
	if err != nil {
		return nil, err
	}

	sz := len(dirs)
	if sz == 0 {
		return nil, nil
	}

	procs := make([]*Proc, 0)

	for i := 0; i < sz; i++ {
		pid, e := strconv.Atoi(dirs[i])
		if e != nil {
			continue
		}
		statusFile := fmt.Sprintf("/proc/%d/status", pid)
		cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
		if !IsExist(statusFile) || !IsExist(cmdlineFile) {
			continue
		}
		ps, err := getPidStatus(statusFile)
		if err != nil {
			continue
		}
		ps.Pid = pid
		cmdline, err := Read2String(cmdlineFile)
		if err != nil {
			continue
		}
		procs = append(procs, &Proc{Ps: ps, Cmdline: cmdline})
	}
	return procs, nil
}
