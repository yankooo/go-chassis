package cipher_test

import (
	security2 "github.com/yankooo/foundation/security"
	"github.com/yankooo/go-chassis/core/lager"
	"github.com/yankooo/go-chassis/security/cipher"
	_ "github.com/yankooo/go-chassis/security/cipher/plugins/aes"
	"github.com/stretchr/testify/assert"
	"testing"
)

//DefaultCipher is a struct
type DefaultCipher struct {
}

func init() {
	cipher.InstallCipherPlugin("default", new)
}
func new() security2.Cipher {

	return &DefaultCipher{}
}

//Encrypt is method used for encryption
func (c *DefaultCipher) Encrypt(src string) (string, error) {
	return src, nil
}

//Decrypt is method used for decryption
func (c *DefaultCipher) Decrypt(src string) (string, error) {
	return src, nil
}

func TestInstallCipherPlugin(t *testing.T) {

	cipher.InstallCipherPlugin("test", new)
	f, err := cipher.GetCipherNewFunc("test")
	assert.NoError(t, err)
	c := f()
	r, _ := c.Encrypt("test")
	assert.Equal(t, "test", r)
	_, err = cipher.GetCipherNewFunc("asd")
	assert.Error(t, err)

	_, err = cipher.GetCipherNewFunc("aes")
	assert.NoError(t, err)
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
