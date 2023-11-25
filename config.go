package main

import (
	"encoding/json"
	"os"

	"github.com/jetking/proxy-watcher/entities"
	"github.com/jetking/proxy-watcher/instance"
	"go.uber.org/zap"
)

var _nodes []entities.Node

func Nodes() []entities.Node {
	if _nodes != nil {
		return _nodes
	}
	data, err := os.ReadFile("./cfg/nodes.json")
	if err != nil {
		panic(err)
	}
	_nodes = []entities.Node{}
	if err := json.Unmarshal(data, &_nodes); err != nil {
		instance.Logger().Warn("failed to load config", zap.Error(err))
		panic(err)
	}
	return _nodes
}
