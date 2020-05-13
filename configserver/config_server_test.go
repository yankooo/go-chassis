package configserver_test

import (
	"testing"

	"github.com/yankooo/go-chassis/configserver"
	"github.com/yankooo/go-chassis/core/config"
	"github.com/yankooo/go-chassis/core/config/model"
	_ "github.com/yankooo/go-chassis/core/registry/servicecenter"
	_ "github.com/yankooo/go-chassis/initiator"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigServerEndpoint(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Config: model.Config{
				Client: model.ConfigClient{},
			},
		},
	}
	_, err := configserver.GetConfigServerEndpoint()
	assert.Error(t, err)
}
