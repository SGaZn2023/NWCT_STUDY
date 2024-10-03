package main

import (
	"fmt"
	"io"
	"net"
	"nwct_st/common"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/xtaci/smux"
)

type Client struct {
	// 连接ID
	clientID string
	// 服务器地址
	serverAddr string
}

// 创建连接
func NewClient(clientID string, serverAddr string) *Client {
	return &Client{
		clientID:   clientID,
		serverAddr: serverAddr,
	}
}

// 运行连接
func (c *Client) Run() {
	for {
		err := c.run()
		if err != nil {
			logs.Error("client run failed: %v", err)
			time.Sleep(3 * time.Second)
		}
		logs.Warn("reconnect %s", c.serverAddr)
	}
}

// 连接核心代码
func (c *Client) run() error {
	// 与服务器建立连接
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return err
	}
	// 函数返回前关闭连接
	defer conn.Close()

	// 发送handshake包
	handshakeReq := &common.HandshakeReq{
		ClientID: c.clientID,
	}
	buf, err := handshakeReq.Encode()
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write(buf)
	conn.SetWriteDeadline(time.Time{})

	if err != nil {
		return err
	}

	// 创建 mux session
	mux, err := smux.Client(conn, nil)

	if err != nil {
		return err
	}

	defer mux.Close()

	// 等待 mux stream
	for {
		stream, err := mux.AcceptStream()
		if err != nil {
			return err
		}

		go c.handleStream(stream)
	}

	// 处理 mux stream
}

func (c *Client) handleStream(stream net.Conn) {
	defer stream.Close()

	// pp 解码
	pp := &common.ProxyProtocol{}
	err := pp.Decode(stream)
	if err != nil {
		logs.Error("decode pp failed: %v", err)
		return
	}
	logs.Debug("pp: %+v", pp)

	// 与本地建立连接
	var localConn net.Conn
	switch pp.InternalProtocol {
	case "tcp":
		localConn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", pp.InternalIP, pp.InternalPort))
		if err != nil {
			logs.Error("connect to local failed: %v", err)
			return
		}
		defer localConn.Close()
	default:
		logs.Warn("unsupported internal protocol: %s", pp.InternalProtocol)
	}

	go func() {
		defer localConn.Close()
		defer stream.Close()
		io.Copy(localConn, stream)
	}()
	io.Copy(stream, localConn)

}
