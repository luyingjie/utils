package udp

import (
	"net"
)

//RunClient 运行udp客户端
//增加出错的处理
func RunClient(udpType, udpURL, data string) error {
	//获取udpaddr
	udpaddr, err := net.ResolveUDPAddr(udpType, udpURL)
	if err != nil {
		return err
	}
	//连接，返回udpconn
	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return err
	}
	//写入数据
	_, err = udpconn.Write([]byte(data))
	if err != nil {
		return err
	}
	//udp 貌似不等待返回结果会直接关闭连接，如果等不到返回结果会阻塞。
	// defer udpconn.Close()
	// buf := make([]byte, 256);
	//读取服务端发送的数据
	// _, err4 := udpconn.Read(buf);
	// clientError(err4);
	// fmt.Println(string(buf));
	return nil
}
