package main

import (
	dis "etcd/discovery"
	"fmt"
	"log"
	"time"
)

// 启动方法
func main() {
	serverName := "dt-service"

	client, err := dis.NewClient(serverName)
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Println(client.ServerList)
		time.Sleep(time.Second * 1)
	}
}
