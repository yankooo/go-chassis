package main

import (
	"github.com/go-mesh/openlogging"
	"github.com/yankooo/go-chassis"
	"github.com/yankooo/go-chassis/core/server"
	"github.com/yankooo/go-chassis/examples/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.RestFulHello{}, server.WithSchemaID("RestHelloService"))
	if err := chassis.Init(); err != nil {
		openlogging.Fatal("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
