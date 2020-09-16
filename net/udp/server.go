package udp

import (
	"net"
	// "bytes"
	// "log"
	// "fmt"
	"time"
	"utils/convert"
)

func clientHandle(conn *net.UDPConn, size int, f func(data string)) error {
	//这里如果关闭连接后面循环会一直连接不上。应该是不关闭连接一直去监听连接拿数据。
	// defer conn.Close();
	// time.Sleep(100 * time.Microsecond)
	// buf := make([]byte, 256);
	buf := make([]byte, size)
	if len(buf) == 0 {
		time.Sleep(100 * time.Microsecond)
	}
	//读取数据
	//注意这里返回三个参数
	//第二个是udpaddr
	//下面向客户端写入数据时会用到
	// _, udpaddr, err := conn.ReadFromUDP(buf);
	_, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	//输出接收到的值
	// fmt.Println(string(buf));
	// buf = bytes.TrimRight(buf, "\u0000") //\x00
	if f != nil {
		// 如果后面的操作需要考虑循序，就不明使用并发。
		go f(convert.ByteToString(buf))
		// f(convert.ByteToString(buf))
	}
	//向客户端发送数据，  测试用。
	// conn.WriteToUDP([]byte("hello,client \r\n"), udpaddr);
	return nil
}

//RunServer 运行udp服务端,发布到服务器上可改为0.0.0.0:xxxx
//f 参数是回掉，这里不返回值，直接调用后面的方法并传值，形成一个管道。后面这里改成那配置去反射方法。
//需要错误判断和停止监听的容错处理。
func RunServer(udpType, udpURL string, size int, f func(data string)) error {
	udpaddr, err := net.ResolveUDPAddr(udpType, udpURL)
	if err != nil {
		return err
	}
	//监听端口
	udpconn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return err
	}
	// defer udpconn.Close()
	//udp没有对客户端连接的Accept函数
	for {
		clientHandle(udpconn, size, f)
	}
	return nil
}
