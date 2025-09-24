package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

// iServer 接口实现，定义一个Server服务类
type Server struct {
	Name       string            // 服务器名称
	IPVersion  string            // tcp4 or other
	IP         string            // 服务绑定的IP地址
	Port       int               // 服务绑定的端口
	MsgHandler ziface.IMsgHanlde // 当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
}

// 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	utils.ServerConfig.Reload()

	s := &Server{
		Name:       utils.ServerConfig.Name,
		IPVersion:  "tcp4",
		IP:         utils.ServerConfig.Host,
		Port:       utils.ServerConfig.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}

// =============================== 实现 ziface.IServer 里的全部接口方法 ===============================

// 启动服务方法
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s, listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.ServerConfig.Version,
		utils.ServerConfig.MaxConn,
		utils.ServerConfig.MaxPacketSize)

	// 开启一个go去做服务端Listen业务
	go func() {
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 2 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		// 已经监听成功
		fmt.Println("start Zinx server", s.Name, "succ, now listenning...")

		// 3 启动server网络连接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			// 3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(conn, utils.GetUid(), s.MsgHandler)

			// 3.4 启动连接的处理业务
			go dealConn.Start()
		}
	}()
}

// 停止服务方法
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

// 启动业务服务方法
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// 路由功能
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("add router succ!")
}
