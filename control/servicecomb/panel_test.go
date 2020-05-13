package servicecomb_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/yankooo/go-chassis/control"
	_ "github.com/yankooo/go-chassis/control/servicecomb"
	"github.com/yankooo/go-chassis/core/common"
	"github.com/yankooo/go-chassis/core/config"
	"github.com/yankooo/go-chassis/core/invocation"
	"github.com/yankooo/go-chassis/core/loadbalancer"
	_ "github.com/yankooo/go-chassis/initiator"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("cse.loadbalance.strategy.name", loadbalancer.StrategyRandom)
	archaius.Set("cse.loadbalance.strategy.Server.name", loadbalancer.StrategyLatency)
	archaius.Set("cse.loadbalance.strategy.name", loadbalancer.StrategyLatency)
	archaius.Set("cse.flowcontrol.Consumer.qps.limit.Server", 100)
	archaius.Set("cse.isolation.Consumer.maxConcurrentRequests", 100)
	err := config.ReadLBFromArchaius()
	if err != nil {
		panic(err)
	}
	err = config.ReadHystrixFromArchaius()
	if err != nil {
		panic(err)
	}
}
func TestPanel_GetLoadBalancing(t *testing.T) {
	opts := control.Options{
		Infra: "archaius",
	}
	err := control.Init(opts)
	assert.NoError(t, err)

	inv := invocation.Invocation{
		MicroServiceName: "Server",
	}
	c := control.DefaultPanel.GetLoadBalancing(inv)
	assert.Equal(t, loadbalancer.StrategyLatency, c.Strategy)

	inv = invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	c = control.DefaultPanel.GetLoadBalancing(inv)
	assert.Equal(t, loadbalancer.StrategyLatency, c.Strategy)

	inv = invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "fake",
	}
	c = control.DefaultPanel.GetLoadBalancing(inv)
	assert.Equal(t, loadbalancer.StrategyLatency, c.Strategy)

	command, cb := control.DefaultPanel.GetCircuitBreaker(inv, common.Consumer)
	assert.Equal(t, "Consumer.fake", command)
	assert.Equal(t, 100, cb.MaxConcurrentRequests)
	inv.MicroServiceName = "Server"
	rl := control.DefaultPanel.GetRateLimiting(inv, common.Consumer)
	assert.Equal(t, 100, rl.Rate)
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.Server", rl.Key)
	assert.Equal(t, true, rl.Enabled)
	t.Run("get server side rate limiting",
		func(t *testing.T) {
			rl := control.DefaultPanel.GetRateLimiting(inv, common.Provider)
			t.Log(rl)
			assert.Equal(t, "cse.flowcontrol.Provider.qps.global.limit", rl.Key)
		})
}

func BenchmarkPanel_GetLoadBalancing(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "Server",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetLoadBalancing(inv)

	}
}
func BenchmarkPanel_GetLoadBalancing2(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetLoadBalancing(inv)

	}
}
func BenchmarkPanel_GetCircuitBreaker(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetCircuitBreaker(inv, common.Consumer)

	}
}
func BenchmarkPanel_GetRateLimiting(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetRateLimiting(inv, common.Consumer)

	}
}
