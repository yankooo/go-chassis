package handler_test

import (
	"github.com/yankooo/go-chassis/core/common"
	"github.com/yankooo/go-chassis/core/config"
	"github.com/yankooo/go-chassis/core/config/model"
	"github.com/yankooo/go-chassis/core/handler"
	"github.com/yankooo/go-chassis/core/invocation"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func TestCreateChain(t *testing.T) {
	t.Log("testing creation of chain with various service type,chain name and handlers")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	e := handler.RegisterHandler("fake", newProviderHandler)
	assert.NoError(t, e)
	c, err := handler.CreateChain("abc", "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = handler.CreateChain(common.Consumer, "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = handler.CreateChain(common.Provider, "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = handler.CreateChain(common.Consumer, "fake")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	chopt := handler.WithChainName("chainName")
	var ch *handler.ChainOptions = new(handler.ChainOptions)
	chopt(ch)
	assert.Equal(t, "chainName", ch.Name)
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func BenchmarkChain_Next(b *testing.B) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	config.GlobalDefinition = &model.GlobalCfg{}
	config.Init()
	iv := &invocation.Invocation{}
	handler.RegisterHandler("f1", createBizkeeperFakeHandler)
	handler.RegisterHandler("f2", createBizkeeperFakeHandler)
	handler.RegisterHandler("f3", createBizkeeperFakeHandler)
	if err := handler.CreateChains(common.Consumer, map[string]string{
		"default": "f1,f2,f3,f1,f2,f3,f1,f2,f3,f1,f2,f3,f1,f2",
	}); err != nil {
		b.Fatal(err)
	}

	c, err := handler.GetChain(common.Consumer, "default")
	if err != nil {
		b.Fatal(err)
	}
	log.Println("----------------------------------------------------")
	log.Println(c)

	for i := 0; i < b.N; i++ {
		c, _ = handler.GetChain(common.Consumer, "default")
		c.Next(iv, func(r *invocation.Response) error {
			return r.Err
		})
		iv.HandlerIndex = 0
	}
}
