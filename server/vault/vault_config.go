package vault

import (
	"github.com/hashicorp/vault/api"
	"ritchie-server/server"
)

type Config struct {
	vaultConfig *api.Config
}

func NewConfig() server.VaultConfig {
	return Config{vaultConfig: api.DefaultConfig()}
}

func (vc Config) Start() (*api.Client, error) {
	_ = vc.vaultConfig.ReadEnvironment()
	client, err := api.NewClient(vc.vaultConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
