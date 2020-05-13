package control_test

import (
	"github.com/yankooo/go-chassis/control"
	_ "github.com/yankooo/go-chassis/control/servicecomb"
	"github.com/yankooo/go-chassis/core/config"
	"github.com/yankooo/go-chassis/core/config/model"
	"github.com/yankooo/go-chassis/core/invocation"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallPlugin(t *testing.T) {
	control.InstallPlugin("test", nil)

}
func TestInit(t *testing.T) {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
	config.GlobalDefinition = &model.GlobalCfg{
		Panel: model.ControlPanel{
			Infra: "",
		},
	}
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	err := control.Init(opts)
	assert.NoError(t, err)
	opts.Infra = "xxx"
	err = control.Init(opts)
	t.Log(err)
	assert.Error(t, err)
}

func TestNewCircuitCmd(t *testing.T) {
	config.HystrixConfig = &model.HystrixConfigWrapper{
		HystrixConfig: &model.HystrixConfig{
			CircuitBreakerProperties: &model.CircuitWrapper{
				Scope: "",
			},
		},
	}
	i := invocation.Invocation{
		MicroServiceName: "mall",
		SchemaID:         "rest",
		OperationID:      "/test",
		Endpoint:         "127.0.0.1:8081",
	}
	cmd := control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.rest./test", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "instance"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.127.0.0.1:8081", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "instance-api"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.127.0.0.1:8081.rest./test", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "api"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.rest./test", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "service"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall", cmd)
}
