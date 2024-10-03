package main

import (
	"net"
	"nwct_st/common"
	"time"

	"github.com/astaxie/beego/logs"
)

type Gateway struct {
	ListenAddr string
	sessionMgr *SessionManager
}

func NewGateway(listenAddr string, sessionMgr *SessionManager) *Gateway {
	gw := &Gateway{
		ListenAddr: listenAddr,
		sessionMgr: sessionMgr,
	}

	go gw.checkOnlineInterval()
	return gw
}

func (gw *Gateway) ListenAndServe() error {
	listener, err := net.Listen("tcp", gw.ListenAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			return err
		}

		go gw.handleConn(conn)
	}
}

func (gw *Gateway) handleConn(conn net.Conn) {
	// defer conn.Close()

	handshakeReq := &common.HandshakeReq{}
	err := handshakeReq.Decode(conn)

	if err != nil {
		logs.Error("decode handshake request failed: %v", err)
		return
	}

	logs.Debug("handshake from: %v", handshakeReq.ClientID)

	// 创建session

	_, err = gw.sessionMgr.CreateSession(handshakeReq.ClientID, conn)
	if err != nil {
		logs.Error("create session failed: %v", err)
		return
	}
}

func (gw *Gateway) checkOnlineInterval() {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()

	for range tick.C {
		gw.sessionMgr.Range(func(k string, v *Session) bool {
			if v.Connection.IsClosed() {
				logs.Info("session %v is offline", v.ClientID)
				return false
			}
			logs.Debug("session %v is online", v.ClientID)
			return true
		})
	}
}
