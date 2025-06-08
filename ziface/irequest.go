package ziface

/*
IRequest 接口：
客户端的所有请求数据都包装到 Request 里
*/
type IRequest interface {
	GetConnection() IConnection // 获取请求连接信息
	GetData() []byte            // 获取请求数据
}
