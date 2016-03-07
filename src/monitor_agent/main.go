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
	for range time.NewTicker(time.Second).C {
		l, _ := nux.LoadAvg()
		info := util.Info{
			Ip:       ta.localIp,
			Avg1min:  l.Avg1min,
			Avg5min:  l.Avg5min,
			Avg15min: l.Avg15min,
		}
		b, _ := json.Marshal(info)
		ta.ioBuf.Write(b)
		ta.ioBuf.Write([]byte("\r\n"))
		err := ta.ioBuf.Flush()
		if err != nil {
			ta.ReConnect()
			ta.ioBuf.Write(b)
			ta.ioBuf.Write([]byte("\r\n"))
			ta.ioBuf.Flush()
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
