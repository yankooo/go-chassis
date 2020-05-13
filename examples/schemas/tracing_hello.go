package schemas

import (
	"net/http"

	"log"

	"github.com/yankooo/go-chassis/client/rest"
	"github.com/yankooo/go-chassis/core"
	"github.com/yankooo/go-chassis/pkg/util/httputil"
	rf "github.com/yankooo/go-chassis/server/restful"
)

//TracingHello is a struct
type TracingHello struct {
}

//Trace is a method
func (r *TracingHello) Trace(b *rf.Context) {
	log.Println("tracing===", b.Ctx)
	req, err := rest.NewRequest("GET", "http://RESTServerB/sayhello/world", nil)
	if err != nil {
		b.WriteError(500, err)
		return
	}

	resp, err := core.NewRestInvoker().ContextDo(b.Ctx, req)
	if err != nil {
		b.WriteError(500, err)
		return
	}
	defer resp.Body.Close()
	b.Write(httputil.ReadBody(resp))
}

//URLPatterns helps to respond for corresponding API calls
func (r *TracingHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/trace", ResourceFunc: r.Trace},
	}
}
