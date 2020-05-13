package main

import (
	"github.com/yankooo/go-chassis"
	"github.com/yankooo/go-chassis/core/lager"
	example "github.com/yankooo/go-chassis/examples/fileupload/server/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/fileupload/server/
func main() {
	chassis.RegisterSchema("rest", &example.RestFulUpload{})

	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
