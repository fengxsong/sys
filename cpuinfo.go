package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
)

func NumCPU() int {
	return runtime.NumCPU()
}

func CpuMHz() (mhz string, err error) {
	fn := "/proc/cpuinfo"
	var fb []byte
	fb, err = ioutil.ReadFile(fn)
	if err != nil {
		return
	}

	rd := bufio.NewReader(bytes.NewBuffer(fb))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if !strings.Contains(line, "MHz") {
			continue
		}
		arr := strings.Split(line, ":")
		return strings.TrimSpace(arr[1]), nil
	}
	return "", fmt.Errorf("no cpu MHz in %s", fn)
}
