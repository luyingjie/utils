package udp

import (
	"net"
	// "log"
	// "fmt"
	"utils/error"
)

//RunClient 运行udp客户端
//增加出错的处理
func RunClient(udpType, udpURL, data string) {
	//获取udpaddr
	udpaddr, err := net.ResolveUDPAddr(udpType, udpURL)
	error.Try(2000, 3, "utils/net/udp/client/RunClient/ResolveUDPAddr", err)
	//连接，返回udpconn
	udpconn, err2 := net.DialUDP("udp", nil, udpaddr)
	error.Try(2000, 3, "utils/net/udp/client/RunClient/DialUDP", err2)
	//写入数据
	_, err3 := udpconn.Write([]byte(data))
	error.Try(5000, 3, "utils/net/udp/client/RunClient/Write", err3)
	//udp 貌似不等待返回结果会直接关闭连接，如果等不到返回结果会阻塞。
	// defer udpconn.Close()
	// buf := make([]byte, 256);
	//读取服务端发送的数据
	// _, err4 := udpconn.Read(buf);
	// clientError(err4);
	// fmt.Println(string(buf));
}
