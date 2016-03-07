package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/toolkits/nux"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	"util"
)

var (
	addrs = flag.String("h", "127.0.0.1:8090,127.0.0.1:8091", "server address")
)

type TcpAgent struct {
	remoteAddr string
	servers    []string
	localIp    string
	conn       net.Conn
	ioBuf      *bufio.Writer
}

func init() {
	flag.Parse()
}

func main() {
	log.Println("start")
	runtime.GOMAXPROCS(runtime.NumCPU())
	ta := NewTcpAgent()
	go ta.run()
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
	log.Println("stop")
}

func (ta *TcpAgent) run() {
	for range time.NewTicker(2 * time.Second).C {
		l, _ := nux.LoadAvg()
		m, _ := nux.MemInfo()
		mountPoints, _ := nux.ListMountPoint()
		diskInfo := []util.DfInfo{}
		for _, dev := range mountPoints {
			du, _ := nux.BuildDeviceUsage(dev[0], dev[1], dev[2])
			diskInfo = append(diskInfo,
				util.DfInfo{
					Mount:     du.FsFile,
					InodeUsed: du.InodesUsedPercent,
					DiskUsed:  du.BlocksUsedPercent,
				})
		}

		info := util.Info{
			Ip:        ta.localIp,
			Avg1min:   l.Avg1min,
			Avg5min:   l.Avg5min,
			Avg15min:  l.Avg15min,
			Buffers:   m.Buffers,
			Cached:    m.Cached,
			MemTotal:  m.MemTotal,
			MemFree:   m.MemFree,
			SwapTotal: m.SwapTotal,
			SwapUsed:  m.SwapUsed,
			SwapFree:  m.SwapFree,
			DiskInfo:  diskInfo,
			Timestamp: time.Now().Unix(),
		}
		b, _ := json.Marshal(info)
		err := ta.writeLine(b)
		if err != nil {
			ta.ReConnect()
			ta.writeLine(b)
		}
	}
}

func NewTcpAgent() (ta *TcpAgent) {
	ta = new(TcpAgent)
	ta.servers = strings.Split(*addrs, ",")
	for _, server := range ta.servers {
		conn, err := net.Dial("tcp", server)
		if err == nil {
			ta.remoteAddr = server
			ta.conn = conn
			ta.localIp = strings.Split(conn.LocalAddr().String(), ":")[0]
			ta.ioBuf = bufio.NewWriter(ta.conn)
			return
		}
	}
	return
}

func (ta *TcpAgent) writeLine(b []byte) error {
	ta.ioBuf.Write(b)
	ta.ioBuf.Write([]byte("\r\n"))
	return ta.ioBuf.Flush()
}

func (ta *TcpAgent) ReConnect() {
	for _, server := range ta.servers {
		conn, err := net.Dial("tcp", server)
		if err == nil {
			ta.remoteAddr = server
			ta.conn = conn
			ta.localIp = strings.Split(conn.LocalAddr().String(), ":")[0]
			ta.ioBuf = bufio.NewWriter(ta.conn)
			return
		}
	}
}
