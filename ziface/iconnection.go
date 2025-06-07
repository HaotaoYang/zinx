package ziface

import "net"

// 定义连接接口
type IConnection interface {
	Start()                   // 启动连接
	Stop()                    // 停止连接
	GetTCPConn() *net.TCPConn // 获取连接socket
	GetConnID() uint32        // 获取连接id
	RemoteAddr() net.Addr     // 获取客户端地址
}

// 定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
