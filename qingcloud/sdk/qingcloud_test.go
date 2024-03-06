package sdk

import (
	"fmt"
	"testing"

	"github.com/yunify/qingcloud-sdk-go/config"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

var qcService *qc.QingCloudService

func init() {
	configuration, _ := config.NewDefault()

	configuration.AccessKeyID = "UFTUUCNNFKSKKLFAKYLN"
	configuration.SecretAccessKey = "Zi9zBZ8rLdyj1N68Fgjkul80NVurKQ7x8WkE32VD"
	configuration.Protocol = "http"
	configuration.Host = "36.110.217.162"
	configuration.Port = 7777

	configuration.URI = "/iaas"

	// configuration.LoadUserConfig()
	qcService, _ = qc.Init(configuration)
}

func TestMain(t *testing.T) {
	// 使用重写方法
	user := "admin@scdemo.com" //"luke@yunify.com"
	passwd := "zhu88jie"       //"Luke_123456"

	userServer, _ := NewUser(qcService, "demo")

	iOutput, _ := userServer.CreateSession(&UserInput{
		User:   &user,
		Passwd: &passwd,
	})
	fmt.Println(iOutput.SK)
	// 使用原生SDK
	limit := 10
	offset := 0
	appServer, _ := qcService.App("demo")
	appOutput, _ := appServer.DescribeApps(&qc.DescribeAppsInput{
		Limit:  &limit,
		Offset: &offset,
	})
	fmt.Println(appOutput.AppSet)
}
