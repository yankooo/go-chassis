package main

import (
	"context"
	"log"
	"time"

	"github.com/yankooo/go-chassis"
	_ "github.com/yankooo/go-chassis/bootstrap"
	"github.com/yankooo/go-chassis/client/rest"
	"github.com/yankooo/go-chassis/core"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	restInvoker := core.NewRestInvoker()

	// use the configured chain
	for {
		callRest(restInvoker, 10)
		<-time.After(time.Second)
	}
}

func callRest(invoker *core.RestInvoker, i int) {
	url := "http://istioserver/sayhello/b"
	if i < 10 {
		url = "http://istioserver/sayhello/a"
	}
	req, _ := rest.NewRequest("GET", url, nil)
	//use the invoker like http client.
	resp1, err := invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
		//lager.Logger.Errorf(err, "call request fail (%s) (%d) ", string(resp1.ReadBody()), resp1.GetStatusCode())
		//return
	}
	log.Println(i, "REST SayHello ------------------------------ ", resp1.StatusCode, string(httputil.ReadBody(resp1)))

	resp1.Body.Close()
}
