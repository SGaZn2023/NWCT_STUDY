# NWCT APP FOR STUDY

// 客户端配置
type ClientInfo struct {
    ClientID string
    PublicIP string
    PublicPort uint16
    PublicProtocol string
    InternalIP string
    InternalPort uint16
    InternalProtocol string
}

// session
type Session struct {
    ClientID string
    Connection net.Conn
}

// 私有协议（proxyProtocol）
type ProxyProtocol struct {
    ClientID string
    PublicIP string
    PublicPort uint16
    PublicProtocol string
    InternalIP string
    InternalPort uint16
    InternalProtocol string
}

// 双向长连接

- smux
- yamux
- optw (github.com/ICKelin/optw (smux/kcp/quic))
