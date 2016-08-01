package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"syscall"
)

var FSSPEC_IGNORE = map[string]struct{}{
	"none":  struct{}{},
	"nodev": struct{}{},
	"tmpfs": struct{}{},
}

var FSTYPE_IGNORE = map[string]struct{}{
	"cgroup":     struct{}{},
	"debugfs":    struct{}{},
	"devtmpfs":   struct{}{},
	"rpc_pipefs": struct{}{},
	"rootfs":     struct{}{},
}

var FS_PREFIX_IGNORE = []string{"/dev", "/sys", "/net", "/misc", "/proc", "lib"}

func IsIgnoreFs(fs_file string) bool {
	for _, prefix := range FS_PREFIX_IGNORE {
		if strings.HasPrefix(fs_file, prefix) {
			return true
		}
	}
	return false
}

type DeviceUsage struct {
	FsSpec            string
	FsFile            string
	FsVfstype         string
	BlocksAll         uint64
	BlocksUsed        uint64
	BlocksFree        uint64
	BlocksUsedPercent float64
	BlocksFreePercent float64
	InodesAll         uint64
	InodesUsed        uint64
	InodesFree        uint64
	InodesUsedPercent float64
	InodesFreePercent float64
}

func (d *DeviceUsage) String() string {
	return fmt.Sprintf("<FsSpec:%s, FsFile:%s, FsVfstype:%s, BlocksFreePercent:%f, InodesFreePercent:%f..>",
		d.FsSpec,
		d.FsFile,
		d.FsVfstype,
		d.BlocksFreePercent,
		d.InodesFreePercent)
}

func ListMounts() ([][3]string, error) {
	fn := "/proc/mounts"
	fb, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	m := make([][3]string, 0)
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
		// Docs come from the fstab(5)
		// fs_spec     # Mounted block special device or remote filesystem e.g. /dev/sda1
		// fs_file     # Mount point e.g. /data
		// fs_vfstype  # File system type e.g. ext4
		// fs_mntops   # Mount options
		// fs_freq     # Dump(8) utility flags
		// fs_passno   # Order in which filesystem checks are done at reboot time
		fs_spec, fs_file, fs_vfstype := fields[0], fields[1], fields[2]
		if _, exist := FSSPEC_IGNORE[fs_spec]; exist {
			continue
		}
		if _, exist := FSTYPE_IGNORE[fs_vfstype]; exist {
			continue
		}
		if IsIgnoreFs(fs_file) {
			continue
		}
		if strings.HasPrefix(fs_spec, "/dev") {
			deviceFound := false
			for idx := range m {
				if m[idx][0] == fs_spec {
					deviceFound = true
					if len(fs_file) < len(m[idx][1]) {
						m[idx][1] = fs_file
					}
					break
				}
			}
			if !deviceFound {
				m = append(m, [3]string{fs_spec, fs_file, fs_vfstype})
			}
		} else {
			m = append(m, [3]string{fs_spec, fs_file, fs_vfstype})
		}
	}
	return m, nil
}

func CheckDeviceUsage(_fsSpec, _fsFile, _fsVfstype string) (*DeviceUsage, error) {
	du := &DeviceUsage{FsSpec: _fsFile, FsFile: _fsFile, FsVfstype: _fsVfstype}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(_fsFile, &fs)
	if err != nil {
		return nil, err
	}

	//blocks
	used := fs.Blocks - fs.Bfree
	du.BlocksAll = uint64(fs.Frsize) * fs.Blocks
	du.BlocksUsed = uint64(fs.Frsize) * used
	du.BlocksFree = uint64(fs.Frsize) * fs.Bavail
	if fs.Blocks == 0 {
		du.BlocksUsedPercent = 100.0
	} else {
		du.BlocksUsedPercent = float64(used) * 100.0 / float64(used+fs.Bavail)
	}
	du.BlocksFreePercent = 100.0 - du.BlocksUsedPercent

	//inodes
	du.InodesAll = fs.Files
	du.InodesFree = fs.Ffree
	du.InodesUsed = fs.Files - fs.Ffree
	if fs.Files == 0 {
		du.InodesUsedPercent = 100.0
	} else {
		du.InodesUsedPercent = float64(du.InodesUsed) * 100.0 / float64(du.InodesAll)
	}
	du.InodesFreePercent = 100.0 - du.InodesUsedPercent
	return du, nil
}
