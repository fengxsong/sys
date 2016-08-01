package sys

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func IsFile(fn string) bool {
	f, err := os.Stat(fn)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func IsExist(fn string) bool {
	_, err := os.Stat(fn)
	return err == nil || os.IsExist(err)
}

func Read2Bytes(fn string) ([]byte, error) {
	return ioutil.ReadFile(fn)
}

func Read2String(fn string) (string, error) {
	b, err := Read2Bytes(fn)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Read2TrimString(fn string) (string, error) {
	s, err := Read2String(fn)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func Read2Uint64(fn string) (uint64, error) {
	s, err := Read2TrimString(fn)
	if err != nil {
		return 0, err
	}
	var u uint64
	if u, err = strconv.ParseUint(s, 10, 64); err != nil {
		return 0, err
	}
	return u, nil
}

func WalkOnlyDirs(dirname string) ([]string, error) {
	if !IsExist(dirname) {
		return nil, fmt.Errorf("No such file or directory: '%s'", dirname)
	}
	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	list := make([]string, 0, sz)
	for _, f := range fs {
		if f.IsDir() {
			basename := f.Name()
			if basename != "." && basename != ".." {
				list = append(list, basename)
			}
		}
	}
	return list, nil
}

func OsWalk(dirname string, recursive, dirOnly bool) ([]string, error) {
	if !IsExist(dirname) {
		return nil, fmt.Errorf("No such file or directory: '%s'", dirname)
	}
	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	list := make([]string, 0, sz)
	for _, f := range fs {
		relpath := filepath.Join(dirname, f.Name())
		if f.IsDir() {
			list = append(list, relpath)
			if recursive {
				subdirList, _ := OsWalk(relpath, recursive, dirOnly)
				list = append(list, subdirList...)
			}
		} else {
			if !dirOnly {
				list = append(list, relpath)
			}
		}
	}
	return list, nil
}
