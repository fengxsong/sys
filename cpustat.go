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

type CpuUsage struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment (since 2.6.11)
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel. (since 2.6.24)
	Total   uint64 // total of all time fields
}

func (c *CpuUsage) String() string {
	return fmt.Sprintf("<User:%d, Nice:%d, System:%d, Idle:%d, Iowait:%d, Irq:%d, SoftIrq:%d, Steal:%d, Guest:%d, Total:%d",
		c.User,
		c.Nice,
		c.System,
		c.Idle,
		c.Iowait,
		c.Irq,
		c.SoftIrq,
		c.Steal,
		c.Guest,
		c.Total)
}

type ProcStat struct {
	Cpu          *CpuUsage
	Cpus         []*CpuUsage
	Ctxt         uint64
	Processes    uint64
	ProcsRunning uint64
	ProcsBlocked uint64
}

func (p *ProcStat) String() string {
	return fmt.Sprintf("<Cpu:%v, Cpus:%v, Ctx:%d, Processes:%d, ProcsRunning:%d, ProcsBlocked:%d",
		p.Cpu,
		p.Cpus,
		p.Ctxt,
		p.Processes,
		p.ProcsRunning,
		p.ProcsBlocked)
}

func CurrentProcStat() (*ProcStat, error) {
	fn := "/proc/stat"
	fb, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ps := &ProcStat{Cpus: make([]*CpuUsage, NumCPU())}
	rd := bufio.NewReader(bytes.NewBuffer(fb))

	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return ps, err
		}
		parseLine(line, ps)
	}
	return ps, nil
}

func parseLine(line string, ps *ProcStat) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return
	}
	fieldName := fields[0]
	if fieldName == "cpu" {
		ps.Cpu = parseCpuFields(fields)
		return
	}
	if strings.HasPrefix(fieldName, "cpu") {
		idx, err := strconv.Atoi(fieldName[3:])
		if err != nil || idx > len(ps.Cpus) {
			return
		}
		ps.Cpus[idx] = parseCpuFields(fields)
		return
	}
	if fieldName == "ctxt" {
		ps.Ctxt, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
	if fieldName == "processes" {
		ps.Processes, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
	if fieldName == "procs_running" {
		ps.ProcsRunning, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
	if fieldName == "procs_blocked" {
		ps.ProcsBlocked, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
}

func parseCpuFields(fields []string) *CpuUsage {
	c := new(CpuUsage)
	l := len(fields)
	for i := 1; i < l; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}
		c.Total += val
		switch i {
		case 1:
			c.User = val
		case 2:
			c.Nice = val
		case 3:
			c.System = val
		case 4:
			c.Idle = val
		case 5:
			c.Iowait = val
		case 6:
			c.Irq = val
		case 7:
			c.SoftIrq = val
		case 8:
			c.Steal = val
		case 9:
			c.Guest = val
		}
	}
	return c
}
