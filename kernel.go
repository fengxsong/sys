package sys

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func KernelMaxFiles() (uint64, error) {
	return Read2Uint64("/proc/sys/fs/file-max")
}

func KernelAllocatedFiles() (uint64, error) {
	fn := "/proc/sys/fs/file-nr"
	s, err := Read2TrimString(fn)
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(s)
	if len(fields) != 3 {
		return 0, fmt.Errorf("%s format error", fn)
	}
	return strconv.ParseUint(fields[0], 10, 64)
}

func KernelMaxProcs() (uint64, error) {
	return Read2Uint64("/proc/sys/kernel/pid_max")
}

func Hostname() (string, error) {
	return os.Hostname()
}
