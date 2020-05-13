package main

import (
	"github.com/yankooo/go-chassis"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/server/restful"
	"net/http"
)

type DemoResource struct {
}

func (r *DemoResource) Limit(b *restful.Context) {
	b.ReadResponseWriter().WriteHeader(http.StatusOK)
	b.ReadResponseWriter().Write([]byte("ok"))
}

// URLPatterns returns routes
func (r *DemoResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/limit", ResourceFunc: r.Limit},
	}
}

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/{project_root}/

func main() {
	chassis.RegisterSchema("rest", &DemoResource{})
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
