package config

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	goprovider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type KeycloakConfig struct {
	Url      string  `pulumi:"url"`
	Username string  `pulumi:"username"`
	Password string  `pulumi:"password" provider:"secret"`
	Insecure *bool   `pulumi:"insecure,optional"`
	Realm    *string `pulumi:"realm,optional"`
	BasePath *string `pulumi:"basePath,optional"`
	Client   *gocloak.GoCloak
}

func (c *KeycloakConfig) Annotate(a infer.Annotator) {
	a.Describe(&c.Url, "URL of the Keycloak Server.")
	a.Describe(&c.Username, "Username for the Keycloak Account.")
	a.Describe(&c.Password, "Password for the Keycloak Account.")
	a.Describe(&c.Insecure, "Allow insecure HTTPS client. Defaults to true.")
	a.Describe(&c.Realm, "Keycloak realm. Defaults to 'master'.")
	a.Describe(&c.BasePath, "Base path for the Keycloak API. Defaults to '/'.")
	a.SetDefault(&c.Url, "localhost:8080")
	a.SetDefault(&c.Username, "admin")
	a.SetDefault(&c.Password, "")
	a.SetDefault(&c.Insecure, true)
	a.SetDefault(&c.Realm, "master")
	a.SetDefault(&c.BasePath, "/")
}

func (c *KeycloakConfig) Configure(ctx context.Context) error {
	goprovider.GetLogger(ctx).Info("Configuring Keycloak Client Provider")

	if c.Url == "" {
		return fmt.Errorf("Keycloak 'url' must be provided in the provider configuration")
	}
	if c.Username == "" {
		return fmt.Errorf("Keycloak 'username' must be provided in the provider configuration")
	}
	if c.Password == "" {
		return fmt.Errorf("Keycloak 'password' must be provided in the provider configuration")
	}

	// Create Keycloak client
	client := gocloak.NewClient(c.Url)
	ctx = context.Background()
	_, err := client.LoginAdmin(ctx, c.Username, c.Password, *c.Realm)
	if err != nil {
		return fmt.Errorf("failed to login to Keycloak: %v", err)
	}

	c.Client = client
	goprovider.GetLogger(ctx).Info("Keycloak Client configured successfully")
	return nil
}
