package sys

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

func Netstat(ext string) (ret map[string]uint64, err error) {
	ret = make(map[string]uint64)
	var b []byte
	b, err = ioutil.ReadFile("/proc/net/netstat")
	if err != nil {
		return
	}

	rd := bufio.NewReader(bytes.NewBuffer(b))
	for {
		var line string
		line, err = rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}

		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}

		subExt := strings.TrimSpace(line[:idx])
		if subExt == ext {
			keys := strings.Fields(strings.TrimSpace(line[idx+1:]))
			var next string
			next, err = rd.ReadString('\n')
			if err != nil {
				return
			}
			valFields := strings.Fields(strings.TrimSpace(next[idx+1:]))
			lenVal := len(valFields)
			for i := 0; i < lenVal; i++ {
				ret[keys[i]], err = strconv.ParseUint(valFields[i], 10, 64)
				if err != nil {
					return
				}
			}
		}
	}
	return
}
