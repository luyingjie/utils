package udp

import (
	"net"
	"testing"
	// "strconv"
	// . "github.com/smartystreets/goconvey/convey"
)

func runClient(udpType, udpURL, data string) {
	udpaddr, _ := net.ResolveUDPAddr(udpType, udpURL)

	udpconn, _ := net.DialUDP("udp", nil, udpaddr)
	udpconn.Write([]byte(data))
}

func TestSendLog(t *testing.T) {
	// for a := 0; a < 10; a++ {
	// runClient("udp4","0.0.0.0:7777","[cloud name] [proxy1] has alerts %7C@@%7C<h3>=== nf_server.log.wf ===</h3>%7C@@%7C2017-11-26 03:49:32,434 CRITICAL -140506499643136- call zucp api failed after 3 retries (/pitrix/lib/pitrix-common/sms/__init__.py:70)")
	// runClient("udp4","0.0.0.0:7777",strconv.Itoa(a))
	runClient("udp4", "0.0.0.0:7777", "[zone_id] [xxx2r01n01] has alerts %7C@@%7C<h3>=== supervisord.log ===</h3>%7C@@%7C2019-06-25 12:32:57,377 INFO exited: storage_server (exit status 255; not expected)")
	// }
	// runClient("udp4", "0.0.0.0:7777", "[cloud name] [webservice0] has alerts %7C@@%7C<h3>=== supervisor ===</h3>%7C@@%7Cws_server                        STOPPED    Oct 12 11:46 AM")
}
