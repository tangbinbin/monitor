package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"util"
)

var (
	port     = flag.Int("P", 8090, "server port")
	addr     = flag.String("h", "127.0.0.1:3306", "mysql addr")
	user     = flag.String("u", "monitor", "mysql user")
	password = flag.String("p", "monitor", "mysql password")
	db       *sql.DB
)

type TcpServer struct {
	remoteAddr string
	remoteIp   string
	conn       net.Conn
	ioBuf      *bufio.Reader
}

func init() {
	flag.Parse()
	initDb()
}

func initDb() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=100ms",
		*user, *password, *addr, "monitor")
	var err error
	db, err = sql.Open("mysql", connStr)
	db.SetMaxOpenConns(10)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println("start server")
	runtime.GOMAXPROCS(runtime.NumCPU())
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("open listen port error:%v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("accept error:%v", err)
			}
			tc := newTcpServer(conn)
			go tc.run()
		}
	}()
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
	log.Println("stop")
}

func newTcpServer(conn net.Conn) *TcpServer {
	return &TcpServer{
		remoteAddr: conn.RemoteAddr().String(),
		remoteIp:   strings.Split(conn.RemoteAddr().String(), ":")[0],
		conn:       conn,
		ioBuf:      bufio.NewReaderSize(conn, 8192),
	}
}

func (tc *TcpServer) run() {
	for {
		line, _, err := tc.ioBuf.ReadLine()
		if err == io.EOF {
			log.Printf("conn closed %s", tc.remoteAddr)
			break
		}
		if err != nil {
			log.Printf("readline err:%v", err)
			break
		}
		info := util.Info{}
		json.Unmarshal(line, &info)
		saveInfo(info)
	}
}

func saveInfo(info util.Info) {
	statsfields := []string{
		"host",
		"load1",
		"load5",
		"load15",
		"buffers",
		"cached",
		"memtotal",
		"memfree",
		"swaptotal",
		"swapused",
		"swapfree",
		"created_time"}
	statsvalues := []string{}
	for i := 0; i < len(statsfields); i++ {
		statsvalues = append(statsvalues, "?")
	}
	statssql := fmt.Sprintf("insert into stats (%s) values (%s)",
		strings.Join(statsfields, ","), strings.Join(statsvalues, ","))
	db.Exec(statssql,
		info.Ip,
		info.Avg1min,
		info.Avg5min,
		info.Avg15min,
		info.Buffers,
		info.Cached,
		info.MemTotal,
		info.MemFree,
		info.SwapTotal,
		info.SwapUsed,
		info.SwapFree,
		info.Timestamp)
	dinfofields := []string{
		"host",
		"mount",
		"inodeused",
		"diskused",
		"created_time"}
	dinfovalues := []string{}
	for i := 0; i < len(dinfofields); i++ {
		dinfovalues = append(dinfovalues, "?")
	}
	dinfosql := fmt.Sprintf("insert into diskinfo (%s) values (%s)",
		strings.Join(dinfofields, ","), strings.Join(dinfovalues, ","))
	for _, dinfo := range info.DiskInfo {
		db.Exec(dinfosql,
			info.Ip,
			dinfo.Mount,
			dinfo.InodeUsed,
			dinfo.DiskUsed,
			info.Timestamp)
	}
}
