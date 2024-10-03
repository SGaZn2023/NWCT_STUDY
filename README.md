# NWCT APP FOR STUDY

## 项目介绍

此项目为一个内网穿透软件（暂仅支持TCP穿透），仅用于学习交流

## 使用方法

将编译好的 gateway(.exe) 与配置文件 gateway.yaml、proxy.json 放入服务器的同一个文件夹下
在该文件夹下输入命令

```bash
./gateway -c gateway.yaml
```

将编译好的 client(.exe) 放在需要被内网穿透的计算机的同一个文件夹下
在该文件夹下输入命令

```bash
./client -client_id="该ID名任取" -server_addr="服务器地址（带端口）"
```

## 配置文件

gateway.yaml

```yaml
listener_file: ./proxy.json    // 读取proxy.json文件
```

proxy.json

```json
[
  {
    "client_id": "此连接ID，需与上述命令中的ID相同", 
    "public_protocol": "tcp",
    "public_ip": "服务器IP，最好写服务器的内网IP",
    "public_port": 20000, 
    "internal_protocol": "tcp",
    "internal_ip": "需要内网穿透的计算机的IP地址，一般为 127.0.0.1",
    "internal_port": 5101
  }
]
```

## 项目配置细节

### 客户端配置

```Go
type ClientInfo struct {
    ClientID string
    PublicIP string
    PublicPort uint16
    PublicProtocol string
    InternalIP string
    InternalPort uint16
    InternalProtocol string
}
```

### session

```Go
type Session struct {
    ClientID string
    Connection net.Conn
}
```

### 私有协议（proxyProtocol）

```Go
type ProxyProtocol struct {
    ClientID string
    PublicIP string
    PublicPort uint16
    PublicProtocol string
    InternalIP string
    InternalPort uint16
    InternalProtocol string
}
```
### 双向长连接

- smux 暂时使用这个
- yamux
- optw (github.com/ICKelin/optw (smux/kcp/quic))
