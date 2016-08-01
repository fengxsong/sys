package main

import (
	"flag"
	"fmt"
	"github.com/fengxsong/sys"
	"os"
)

func main() {
	test := flag.Bool("test", false, "if test is true, only run the sys func")
	dirname := flag.String("dirname", ".", "dirname")
	recursive := flag.Bool("r", true, "recursive")
	dirOnly := flag.Bool("dd", false, "if dirs only")
	flag.Parse()

	if *test {

		fmt.Println(sys.NumCPU())
		fmt.Println(sys.CpuMHz())
		fmt.Println(sys.CurrentProcStat())
		fmt.Println(sys.ListMounts())
		fmt.Println(sys.CheckDeviceUsage("/dev/mapper/centos-root", "/", "xfs"))
		fmt.Println(sys.IfStat())
		fmt.Println(sys.CheckDiskStats())
		fmt.Println(sys.KernelMaxFiles())
		fmt.Println(sys.KernelMaxProcs())
		fmt.Println(sys.Hostname())
		fmt.Println(sys.MemInfo())
		fmt.Println(sys.Netstat("TcpExt"))
		fmt.Println(sys.Netstat("IpExt"))
		fmt.Println(sys.ListeningPorts())
		fmt.Println(sys.UdpPorts())
		fmt.Println(sys.Uptime())
		allProcs, _ := sys.AllProcs()
		for _, i := range allProcs {
			fmt.Println("Cmdline: ", i.Cmdline)
			fmt.Println("Pid: ", i.Ps.Pid)
			fmt.Println("Name: ", i.Ps.Name)
			fmt.Println("State: ", i.Ps.State)
		}
		os.Exit(0)
	}
	list, err := sys.OsWalk(*dirname, *recursive, *dirOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, i := range list {
		fmt.Println(i)
	}
}
