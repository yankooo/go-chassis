package profile

import (
	"github.com/yankooo/go-chassis/core/config"
	"github.com/yankooo/go-chassis/core/registry"
	"github.com/yankooo/go-chassis/core/router"
	_ "github.com/yankooo/go-chassis/core/router/servicecomb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfile(t *testing.T) {
	err := router.BuildRouter("cse")
	assert.NoError(t, err)
	rr := map[string][]*config.RouteRule{"test": {{Precedence: 10}}}
	router.DefaultRouter.SetRouteRule(rr)

	registry.MicroserviceInstanceIndex = registry.NewIndexCache()
	registry.MicroserviceInstanceIndex.Set("test", []*registry.MicroServiceInstance{{InstanceID: "id"}})

	p := newProfile()

	assert.Equal(t, 10, p.RouteRule["test"][0].Precedence)
	assert.Equal(t, "id", p.Discovery["test"][0].InstanceID)
}
