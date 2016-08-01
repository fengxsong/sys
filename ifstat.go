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

type NetIf struct {
	Interface          string
	ReceiveBytes       int64
	ReceivePackets     int64
	ReceiveErrs        int64
	ReceiveDrop        int64
	ReceiveFifo        int64
	ReceiveFrame       int64
	ReceiveCompressed  int64
	ReceiveMulticast   int64
	TransmitBytes      int64
	TransmitPackets    int64
	TransmitErrs       int64
	TransmitDrop       int64
	TransmitFifo       int64
	TransmitColls      int64
	TransmitCarrier    int64
	TransmitCompressed int64
	TotalBytes         int64
	TotalPackets       int64
	TotalErrs          int64
	TotalDrop          int64
}

func (n *NetIf) String() string {
	return fmt.Sprintf("<Interface:%s, TotalBytes:%d, TotalPackets:%d, TotalErrs:%d, TotalDrop:%d.>",
		n.Interface,
		n.TotalBytes,
		n.TotalPackets,
		n.TotalErrs,
		n.TotalDrop)
}

func IfStat() ([]*NetIf, error) {
	fn := "/proc/net/dev"
	fb, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ret := []*NetIf{}
	rd := bufio.NewReader(bytes.NewBuffer(fb))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}

		netIf := NetIf{}
		deviceName := strings.TrimSpace(line[0:idx])
		fields := strings.Fields(line[idx+1:])
		if len(fields) != 16 {
			continue
		}
		netIf.Interface = deviceName

		netIf.ReceiveBytes, _ = strconv.ParseInt(fields[0], 10, 64)
		netIf.ReceivePackets, _ = strconv.ParseInt(fields[1], 10, 64)
		netIf.ReceiveErrs, _ = strconv.ParseInt(fields[2], 10, 64)
		netIf.ReceiveDrop, _ = strconv.ParseInt(fields[3], 10, 64)
		netIf.ReceiveFifo, _ = strconv.ParseInt(fields[4], 10, 64)
		netIf.ReceiveFrame, _ = strconv.ParseInt(fields[5], 10, 64)
		netIf.ReceiveCompressed, _ = strconv.ParseInt(fields[6], 10, 64)
		netIf.ReceiveMulticast, _ = strconv.ParseInt(fields[7], 10, 64)

		netIf.TransmitBytes, _ = strconv.ParseInt(fields[8], 10, 64)
		netIf.TransmitPackets, _ = strconv.ParseInt(fields[9], 10, 64)
		netIf.TransmitErrs, _ = strconv.ParseInt(fields[10], 10, 64)
		netIf.TransmitDrop, _ = strconv.ParseInt(fields[11], 10, 64)
		netIf.TransmitFifo, _ = strconv.ParseInt(fields[12], 10, 64)
		netIf.TransmitColls, _ = strconv.ParseInt(fields[13], 10, 64)
		netIf.TransmitCarrier, _ = strconv.ParseInt(fields[14], 10, 64)
		netIf.TransmitCompressed, _ = strconv.ParseInt(fields[15], 10, 64)

		netIf.TotalBytes = netIf.ReceiveBytes + netIf.TransmitBytes
		netIf.TotalDrop = netIf.ReceiveDrop + netIf.TransmitDrop
		netIf.TotalErrs = netIf.ReceiveErrs + netIf.TransmitErrs
		netIf.TotalPackets = netIf.ReceivePackets + netIf.TransmitPackets

		ret = append(ret, &netIf)
	}
	return ret, nil
}
