package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection // 客户端连接
	data []byte             // 客户端请求数据
}

// 获取连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求数据
func (r *Request) GetData() []byte {
	return r.data
}
