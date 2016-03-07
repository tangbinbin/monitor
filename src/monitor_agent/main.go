package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/toolkits/nux"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
	"util"
)

var (
	addr = flag.String("h", "127.0.0.1:8090", "server address")
)

func init() {
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conn, err := net.Dial("tcp", *addr)
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	if err != nil {
		log.Fatalf("open conn failed:%v", err)
	}
	iobuf := bufio.NewWriter(conn)
	for range time.NewTicker(time.Second).C {
		l, _ := nux.LoadAvg()
		info := util.Info{
			Ip:       ip,
			Avg1min:  l.Avg1min,
			Avg5min:  l.Avg5min,
			Avg15min: l.Avg15min,
		}
		b, _ := json.Marshal(info)
		iobuf.Write(b)
		iobuf.Write([]byte("\r\n"))
		iobuf.Flush()
	}
}
