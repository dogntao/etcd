package discovery

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type server struct {
	serverName string // 服务名
	ipPort     string // ip+端口
	cli        *clientv3.Client
	leaseID    clientv3.LeaseID
}

func NewServer(serverName string, ipPort string) (*server, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:12379", "http://127.0.0.1:22379", "http://127.0.0.1:32379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return nil, err
	}

	return &server{
		serverName: serverName,
		ipPort:     ipPort,
		cli:        cli,
	}, nil
}

// 启动方法
func (s *server) Run() error {
	ch, err := s.register()
	if err != nil {
		log.Fatal(err)
		return err
	}

	for {
		select {
		case cv, ok := <-ch:
			if !ok {
				fmt.Print("server stop")
				return nil
			} else {
				fmt.Print(cv.TTL)
			}
		}
	}
}

// 往etcd写入数据，并keep-alive
func (s *server) register() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	key := "services/" + s.serverName
	grant, err := s.cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.cli.Put(context.TODO(), key, s.ipPort, clientv3.WithLease(grant.ID))
	s.leaseID = grant.ID
	return s.cli.KeepAlive(context.TODO(), grant.ID)
}

// 停止方法
func (s *server) Stop() {
	s.cli.Revoke(context.TODO(), s.leaseID)
}
