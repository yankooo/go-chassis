package client_test

import (
	"github.com/yankooo/go-chassis/pkg/scclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithGlobal(t *testing.T) {
	o := client.WithGlobal()
	opts := &client.CallOptions{}
	o(opts)
	assert.True(t, opts.WithGlobal)
}
