package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn      // 连接的socket套接字
	ConnID       uint32            // 连接id
	isClosed     bool              // 连接是否关闭
	Router       ziface.IRouter    // 连接的处理方法router
	MsgHandler   ziface.IMsgHanlde // 消息管理MsgId和对应处理方法的消息管理模块
	ExitBuffChan chan bool         // 告知连接状态的channel
}

// 创建连接
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHanlde) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
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
		// 创建拆包解包对象
		dp := NewDataPack()

		// 读取数据Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConn(), headData); err != nil {
			fmt.Println("read msg head err ", err)
			break
		}

		// 拆包，得到msgId 和 dataLen 放在 msg 中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		// 根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConn(), data); err != nil {
				fmt.Println("read msg data error ", err)
				continue
			}
		}
		msg.SetData(data)

		// 得到客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go c.MsgHandler.DoMsgHandler(&req)
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

// 发送数据
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	// 将data封包并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}

	return nil
}
