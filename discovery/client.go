package discovery

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type client struct {
	serverName string // 服务名
	cli        *clientv3.Client
	ServerList map[string]string
}

func NewClient(serverName string) (*client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:12379", "http://127.0.0.1:22379", "http://127.0.0.1:32379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return nil, err
	}

	client := &client{
		serverName: serverName,
		cli:        cli,
		ServerList: make(map[string]string),
	}

	go client.watch()
	return client, err
}

// 监听方法
func (c *client) watch() {
	key := "services/" + c.serverName
	wcReps := c.cli.Watch(context.TODO(), key)
	for wcResp := range wcReps {
		for _, v := range wcResp.Events {
			switch v.Type {
			case clientv3.EventTypePut:
				c.ServerList[string(v.Kv.Key)] = string(v.Kv.Value)
			case clientv3.EventTypeDelete:
				delete(c.ServerList, string(v.Kv.Key))
			}
		}
	}
}
