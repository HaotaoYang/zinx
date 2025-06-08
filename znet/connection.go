package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn   // 连接的socket套接字
	ConnID       uint32         // 连接id
	isClosed     bool           // 连接是否关闭
	Router       ziface.IRouter // 连接的处理方法router
	ExitBuffChan chan bool      // 告知连接状态的channel
}

// 创建连接
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// 读数据
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		// 读取最大的数据到buf中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}
		// 得到客户端请求的Request数据
		req := Request{
			conn: c,
			data: buf,
		}
		// 从路由 Routers 中找到对应的Handle
		go func(request ziface.IRequest) {
			// 执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// 启动连接
func (c *Connection) Start() {
	// 读取数据
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息, 不在阻塞
			return
		}
	}
}

// 停止连接
func (c *Connection) Stop() {
	// 1 连接已关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// TODO Connection Stop()

	// 关闭socket连接
	c.Conn.Close()

	// 通知连接关闭
	c.ExitBuffChan <- true

	// 关闭全部管道
	close(c.ExitBuffChan)
}

// 获取socket TCPConn
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

// 获取连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取客户端地址
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
