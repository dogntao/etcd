package main

import (
	dis "etcd/discovery"
	"log"
	"time"
)

// 启动方法
func main() {
	serverName := "dt-service"
	ipPort := "127.0.0.1:1"

	server, err := dis.NewServer(serverName, ipPort)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		time.Sleep(time.Second * 10)
		server.Stop()
	}()
	server.Run()
}
