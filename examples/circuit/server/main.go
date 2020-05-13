package main

import (
	"github.com/yankooo/go-chassis"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/core/server"
	"github.com/yankooo/go-chassis/examples/circuit/server/resource"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &resource.RestFulMessage{}, server.WithSchemaID("RestHelloService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
