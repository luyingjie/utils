package manage

// IConnection : 定义链接接口
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start() error
	//停止连接，结束当前连接状态
	Stop() error

	//从当前连接获取原始的socket TCPConn
	GetConnection() interface{}
	//获取当前连接ID
	GetConnID() string

	// 发送消息
	SendMsg(msgId int, data []byte) error
}
