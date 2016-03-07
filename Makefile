all: install

GOPATH:=$(CURDIR)
export GOPATH

dep:
	go get "github.com/toolkits/nux"
	go get "github.com/go-sql-driver/mysql"

install:dep
	go install monitor_agent
	go install monitor_server
