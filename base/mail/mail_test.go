package mail

 import (
     "testing"
     "fmt"
     "gopkg.in/gomail.v2"
 )



 func TestMain(t *testing.T) {
     fmt.Println("开始")
     m := gomail.NewMessage()

     m.SetAddressHeader("From", "luke@yunify.com" /*"发件人地址"*/, "发件人") // 发件人

     m.SetHeader("To",
         m.FormatAddress("franklinhe@yunify.com", "收件人")) // 收件人
     //  m.SetHeader("Cc",
     //     m.FormatAddress("xxxx@foxmail.com", "收件人")) //抄送
     // m.SetHeader("Bcc",
     //     m.FormatAddress("xxxx@gmail.com", "收件人")) /暗送

     m.SetHeader("Subject", "liic测试")     // 主题

     body := `
 		<html>
 		<body>
 		<h3>
 		"Test send to email"
 		</h3>
 		</body>
 		</html>
 		`
     m.SetBody("text/html",body) // 可以放html..还有其他的
     // m.SetBody("我是正文") // 正文

     // m.Attach("我是附件")  //添加附件

     d := gomail.NewPlainDialer("smtp.yunify.com", 465, "luke@yunify.com", "Luke_1234567") // 发送邮件服务器、端口、发件人账号、发件人密码
     if err := d.DialAndSend(m); err != nil {
         fmt.Println("发送失败", err)
         return
     }

     fmt.Println("done.发送成功")

     fmt.Println("结束")
 }
