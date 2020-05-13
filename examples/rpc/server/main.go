package main

import (
	"github.com/yankooo/go-chassis"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/core/server"
	"github.com/yankooo/go-chassis/examples/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rpc/server/

func main() {
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, server.WithSchemaID("HelloService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.")
		return
	}
	chassis.Run()
}
