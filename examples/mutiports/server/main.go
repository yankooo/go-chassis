package main

import (
	"github.com/yankooo/go-chassis"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/examples/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.Hello{})
	chassis.RegisterSchema("rest-legacy", &schemas.Legacy{})
	chassis.RegisterSchema("rest-admin", &schemas.Admin{})
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
