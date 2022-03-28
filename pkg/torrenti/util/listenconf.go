package util

import (
	"fmt"
	"net"
)

type ListenConf struct {
	Addr string `env:"ADDR"`
	Port int    `env:"PORT"`
}

func (c ListenConf) GetAddr() string {
	addr := c.Addr
	if addr == "" && c.Port != 0 {
		addr = fmt.Sprintf(":%v", c.Port)
	}
	return addr
}

func (c ListenConf) Serve(ss ServeService) error {
	listener, err := c.Listen()
	if err != nil {
		return err
	}
	return ss.Serve(listener)
}

type ServeService interface {
	Serve(l net.Listener) error
}

func (c ListenConf) Listen() (net.Listener, error) {
	addr := c.GetAddr()
	if addr == "" {
		return nil, fmt.Errorf("no address or port")
	}
	return net.Listen("tcp", addr)
}
