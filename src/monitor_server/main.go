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
	ioBuf      *bufio.ReadWriter
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
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error:%v", err)
		}
		tc := newTcpServer(conn)
		go tc.run()
	}
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
		ioBuf:      bufio.NewReadWriter(bufio.NewReaderSize(conn, 1024), bufio.NewWriter(conn)),
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
	ssql := "insert into avgload (host,load1,load5,load15) values (?,?, ?, ?)"
	db.Exec(ssql, info.Ip, info.Avg1min, info.Avg5min, info.Avg15min)
}

func (tc *TcpServer) read(b []byte) (err error) {
	total := 0
	for {
		var n int
		n, err = tc.ioBuf.Read(b[total:])
		if err != nil {
			break
		}
		total += n
		if total == len(b) {
			break
		}
	}
	return
}
