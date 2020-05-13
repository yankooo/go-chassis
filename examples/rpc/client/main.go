package main

import (
	"context"

	"github.com/yankooo/go-chassis"
	_ "github.com/yankooo/go-chassis/bootstrap"
	"github.com/yankooo/go-chassis/core"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/examples/schemas/helloworld"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rpc/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	//declare reply struct
	reply := &helloworld.HelloReply{}
	//Invoke with microservice name, schema ID and operation ID
	if err := core.NewRPCInvoker().Invoke(context.Background(), "RPCServer", "HelloService", "SayHello",
		&helloworld.HelloRequest{Name: "Peter"}, reply); err != nil {
		lager.Logger.Error("error" + err.Error())
	}
	lager.Logger.Info(reply.Message)
}
