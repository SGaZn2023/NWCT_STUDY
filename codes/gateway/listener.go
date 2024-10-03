package main

import (
	"fmt"
	"io"
	"net"
	"nwct_st/common"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	writeTimeout = time.Second * 1
)

type Listener struct {
	pp          *common.ProxyProtocol
	sessionMgr  *SessionManager
	closeOnce   sync.Once
	close       chan struct{}
	tcpListener net.Listener
}

func NewListener(pp *common.ProxyProtocol, sessionMgr *SessionManager) *Listener {
	return &Listener{
		pp:         pp,
		close:      make(chan struct{}),
		sessionMgr: sessionMgr,
	}
}

func (l *Listener) ListenAndServer() error {
	switch l.pp.PublicProtocol {
	case "tcp":
		return l.listenAndServeTCP()
	default:
		return fmt.Errorf("unsupported protocol: %s", l.pp.PublicProtocol)
	}
}

func (l *Listener) listenAndServeTCP() error {
	listenerAddr := fmt.Sprintf("%s:%d", l.pp.PublicIP, l.pp.PublicPort)
	listener, err := net.Listen("tcp", listenerAddr)
	if err != nil {
		return err
	}

	defer listener.Close()
	l.tcpListener = listener

	for {
		conn, err := listener.Accept()

		if err != nil {
			return err
		}

		go l.handleConn(conn)
	}
}

func (l *Listener) handleConn(conn net.Conn) {
	defer conn.Close()

	// 查询seesion
	tunnelConn, err := l.sessionMgr.GetSessionByClientID(l.pp.ClientID)
	if err != nil {
		logs.Warn("get session for client %s failed", l.pp.ClientID)
		return
	}

	defer tunnelConn.Close()

	// 封装proxyprotocol
	ppdoby, err := l.pp.Encode()

	if err != nil {
		logs.Warn("encode pp failed: %v", err)
		return
	}

	tunnelConn.SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err = tunnelConn.Write(ppdoby)
	tunnelConn.SetWriteDeadline(time.Time{})
	if err != nil {
		logs.Warn("write pp body failed: %v", err)
		return
	}

	// 双向数据拷贝
	go func() {
		defer tunnelConn.Close()
		defer conn.Close()
		io.Copy(tunnelConn, conn)
	}()
	io.Copy(conn, tunnelConn)

}

func (l *Listener) Close() {
	l.closeOnce.Do(func() {
		close(l.close)
		if l.tcpListener != nil {
			l.tcpListener.Close()
		}
	})
}
