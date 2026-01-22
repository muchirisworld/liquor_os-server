package clerk

import (
	"github.com/All-Things-Muchiri/server/internal/config"
	clerk "github.com/clerk/clerk-sdk-go/v2"
	client "github.com/clerk/clerk-sdk-go/v2/client"
)

type ClerkProvider struct {
	client       client.Client
	clientConfig *clerk.ClientConfig
	authConfig   config.AuthConfig
}

func NewProvider(authConfig config.AuthConfig) *ClerkProvider {
	clientCfg := clerk.ClientConfig{}
	clientCfg.Key = &authConfig.SecretKey

	return &ClerkProvider{
		client:       *client.NewClient(&clientCfg),
		clientConfig: &clientCfg,
	}
}
