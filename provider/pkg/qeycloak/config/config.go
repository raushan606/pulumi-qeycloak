package config

import (
	"context"
	"fmt"

	goprovider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type KeycloakConfig struct {
	Url      string  `pulumi:"url"`
	Username string  `pulumi:"username"`
	Password string  `pulumi:"password" provider:"secret"`
	Realm    *string `pulumi:"realm,optional"`
}

func (c *KeycloakConfig) Annotate(a infer.Annotator) {
	a.Describe(&c.Url, "URL of the Keycloak Server.")
	a.Describe(&c.Username, "Username for the Keycloak Account.")
	a.Describe(&c.Password, "Password for the Keycloak Account.")
	a.Describe(&c.Realm, "Keycloak realm. Defaults to 'master'.")
}

func (c *KeycloakConfig) Configure(ctx context.Context) error {
	goprovider.GetLogger(ctx).Info("Configuring Keycloak Config Provider")

	if c.Url == "" {
		return fmt.Errorf("keycloak 'url' must be provided in the provider configuration")
	}
	if c.Username == "" {
		return fmt.Errorf("keycloak 'username' must be provided in the provider configuration")
	}
	if c.Password == "" {
		return fmt.Errorf("keycloak 'password' must be provided in the provider configuration")
	}
	return nil
}
