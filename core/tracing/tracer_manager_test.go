package tracing_test

import (
	"errors"
	"github.com/yankooo/go-chassis/core/config"
	"github.com/yankooo/go-chassis/core/config/model"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/core/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func test(o map[string]string) (opentracing.Tracer, error) {
	return nil, nil
}
func fake(o map[string]string) (opentracing.Tracer, error) {
	return nil, errors.New("123")
}
func TestTracerManager(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{}
	tracing.InstallTracer("test", test)
	tracing.InstallTracer("fake", fake)
	err := tracing.Init()
	assert.NoError(t, err)
	config.GlobalDefinition.Tracing = model.TracingStruct{
		Tracer: "test",
	}
	err = tracing.Init()
	assert.NoError(t, err)

	config.GlobalDefinition.Tracing = model.TracingStruct{
		Tracer: "fake",
	}
	err = tracing.Init()
	assert.Error(t, err)
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
