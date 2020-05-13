package aes

import (
	"github.com/yankooo/go-chassis/security/cipher"
	"os"

	security2 "github.com/yankooo/foundation/security"
	"github.com/yankooo/go-chassis/pkg/goplugin"
	"github.com/go-mesh/openlogging"
)

const cipherPlugin = "cipher_plugin.so"

//Cipher interface declares Init(), Encrypyt(), Decrypyt() methods
type Cipher interface {
	Init()
	Encrypt(src string) (string, error)
	Decrypt(src string) (string, error)
}

// HWAESCipher is a cipher used in huawei
type HWAESCipher struct {
	gcryptoEngine Cipher
}

func init() {
	if v, exist := os.LookupEnv("CIPHER_ROOT"); exist {
		err := os.Setenv("PAAS_CRYPTO_PATH", v)
		if err != nil {
			openlogging.Warn("can not set env for cipher: " + err.Error())
		}
	}
	cipher.InstallCipherPlugin("aes", new)
}

func new() security2.Cipher {
	cipher, err := goplugin.LookUpSymbolFromPlugin(cipherPlugin, "Cipher")
	if err != nil {
		if os.IsNotExist(err) {
			openlogging.GetLogger().Errorf("%s not found", cipherPlugin)
		} else {
			openlogging.GetLogger().Errorf("Load %s failed, err [%s]", cipherPlugin, err.Error())
		}
		return nil
	}
	cipherInstance, ok := cipher.(Cipher)
	if !ok {
		openlogging.GetLogger().Infof("E: Expecting Cipher interface, but got something else.")
		return nil
	}
	cipherInstance.Init()
	return &HWAESCipher{
		gcryptoEngine: cipherInstance,
	}
}

//Encrypt is method used for encryption
func (ac *HWAESCipher) Encrypt(src string) (string, error) {
	return ac.gcryptoEngine.Encrypt(src)
}

//Decrypt is method used for decryption
func (ac *HWAESCipher) Decrypt(src string) (string, error) {
	return ac.gcryptoEngine.Decrypt(src)
}
