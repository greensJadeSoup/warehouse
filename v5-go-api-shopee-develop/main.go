package main

import (
	"warehouse/v5-go-api-shopee/conf"
	"warehouse/v5-go-api-shopee/mq"
	"warehouse/v5-go-component/cp_app"

	_ "warehouse/v5-go-api-shopee/api"
)

func init() {
	cp_app.InitInstance("shopee", "1.2.2", conf.InitConf, mq.InitMQ)
}

func main() {
	//go mq.ProduceRocketMQ()
	cp_app.Run()
}
