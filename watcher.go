package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jetking/proxy-watcher/entities"
	"github.com/jetking/proxy-watcher/instance"
	"go.uber.org/zap"
)

const Interval = 10 // seconds
const Concurrency = 1
const Dest = "https://wang.mx?source=proxy-watcher"

type workerItem struct {
	Node entities.Node
	Port int
}

func WatcherStart() {
	ch := make(chan workerItem, 32)
	for i := 0; i < Concurrency; i++ {
		go func(n int) {
			instance.Logger().Info("channel started", zap.Int("id", n))
			for item := range ch {
				worker(item.Node, item.Port)
			}
		}(i)
	}
	nodes := Nodes()
	t := time.NewTicker(Interval * time.Second)
	for range t.C {
		for k := range nodes {
			for p := nodes[k].PortRange[0]; p <= nodes[k].PortRange[1]; p++ {
				ch <- workerItem{Node: nodes[k], Port: p}
			}
		}
	}
}

func worker(node entities.Node, port int) {
	instance.Logger().Sugar().Infof("start worker [%s:%d]", node.Host, port)
	proxyURL, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", node.User, node.Pass, node.Host, port))
	proxyURL.User = url.UserPassword(node.User, node.Pass)
	cli := http.DefaultClient
	cli.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	cli.Timeout = 3 * time.Second

	request, err := http.NewRequest("GET", Dest, nil)
	if err != nil {
		instance.Logger().Error("Error creating HTTP request", zap.Error(err))
		return
	}
	t1 := time.Now()
	response, err := cli.Do(request)
	if err != nil {
		instance.Logger().Error("Error sending HTTP request", zap.Error(err))
		return
	}
	times := time.Since(t1).Milliseconds()
	instance.Logger().Sugar().Infof("[%s:%d] done, time: %dms, httpStatus:%s", node.Host, port, times, response.Status)
}
