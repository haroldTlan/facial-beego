package main

import (
	"flag"

	_ "DemoForService/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	_ "github.com/go-sql-driver/mysql"
	"verifyconsumer/controllers"
	_ "verifyconsumer/routers"
)

var (
	nsq_ip      = beego.AppConfig.String("nsq") + ":" + beego.AppConfig.String("nsq_pub_port")
	nsqdAddr    = flag.String("nsqds", nsq_ip, "nsqd http address")
	maxInFlight = flag.Int("max-in-flight", 200, "Maximum amount of messages in flight to consume")
	publish     = controllers.Publish
	//apiKey      = flag.String("api_key", "5bf9daae737e43d0a05c37fc4ca6c3fa", "API key of linkface")
	//apiSecret   = flag.String("api_secret", "784eebe6a8a243ce92091597dc8bff78", "API secret of linkface")
	apiKey    = flag.String("api_key", "pPcCQavGfltqRq6m8vIFgALMpmaS3BhI", "API key of face++")
	apiSecret = flag.String("api_secret", "dbwoZSx8BUYRIwQudeSxZiab7WnQEhST", "API secret of face++")
)

func init() {
	err := orm.RegisterDataBase("default", "mysql", "root:passwd@tcp(127.0.0.1:3306)/facial?charset=utf8")
	if err != nil {
		panic(err)
	}

	logs.SetLogger(logs.AdapterFile, `{"filename":"/var/log/facial.log","daily":false,"maxdays":365,"level":3}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
}
func main() {
	flag.Parse()

	go runConsumer(*maxInFlight, *nsqdAddr)

	// WebSocket.
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")

	beego.Run()
}
