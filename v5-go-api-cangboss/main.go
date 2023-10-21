package main

import (
	"warehouse/v5-go-api-cangboss/conf"
	"warehouse/v5-go-api-cangboss/mq"
	"warehouse/v5-go-component/cp_app"

	_ "warehouse/v5-go-api-cangboss/api"
)

func init() {
	cp_app.InitInstance("cangboss", "1.2.5", mq.InitMQ, conf.InitConf)
}

func main() {
	//go mq.RocketMQProducer()
	cp_app.Run()
}
