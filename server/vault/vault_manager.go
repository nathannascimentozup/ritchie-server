package vault

import (
	"fmt"
	"os"

	vlogin "github.com/ZupIT/go-vault-session/pkg/login"
	"github.com/ZupIT/go-vault-session/pkg/token"
	"github.com/hashicorp/vault/api"

	"ritchie-server/server"
	"ritchie-server/server/slicer"

	log "github.com/sirupsen/logrus"
)

const vaultPath = "ritchie/credential/%s"

type Manager struct {
	client *api.Client
}

func NewVaultManager(c *api.Client) server.VaultManager {
	return Manager{client: c}
}

func (vm Manager) Write(key string, data map[string]interface{}) error {
	vm.setToken()
	if _, err := vm.client.Logical().Write(pathResolver(key), data); err != nil {
		log.Error("Vault write error ", err)
		return err
	}
	return nil
}

func (vm Manager) Read(key string) (map[string]interface{}, error) {
	vm.setToken()
	res, err := vm.client.Logical().Read(pathResolver(key))

	if err != nil {
		log.Error("Vault read error ", err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return res.Data, nil
}

func (vm Manager) Delete(key string) error {
	vm.setToken()
	_, err := vm.client.Logical().Delete(pathResolver(key))
	if err != nil {
		log.Error("Vault delete error ", err)
		return err
	}

	return nil
}

func (vm Manager) List(key string) ([]interface{}, error) {
	vm.setToken()
	res, err := vm.client.Logical().List(pathResolver(key))
	if err != nil {
		log.Error("Vault list error ", err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	a := res.Data["keys"]

	keys, err := slicer.NewSlicer(a).Interface()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (vm Manager) Start(c *api.Client) {
	l := vlogin.NewHandler(c)
	s := l.Handle()

	ch := make(chan string)
	r := token.NewHandler(c, s, ch)

	r.Handle()
	go func() {
		for {
			msg := <-ch
			log.Info(msg)
		}
	}()

}

func (vm Manager) setToken() {
	vm.client.SetToken(os.Getenv(api.EnvVaultToken))
}

func pathResolver(key string) string {
	return fmt.Sprintf(vaultPath, key)
}
