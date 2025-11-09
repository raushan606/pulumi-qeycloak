package qeycloak

import (
	"os"
	"strconv"

	"github.com/Nerzal/gocloak/v13"
	"github.com/netascode/go-aci"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Config struct {
	Url      string `pulumi:"url"`
	Username string `pulumi:"username"`
	Password string `pulumi:"password" provider:"secret"`
	Insecure *bool  `pulumi:"insecure,optional"`
	Retries  *int   `pulumi:"retries,optional"`
	Logging  *bool  `pulumi:"logging,optional"`
	Client   *gocloak.Client
}

var _ = (infer.Annotated)((*Config)(nil))

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Url, "URL of the Cisco APIC web interface. This can also be set as the ACI_URL environment variable.")
	a.Describe(&c.Username, "Username for the APIC Account. This can also be set as the ACI_USERNAME environment variable.")
	a.Describe(&c.Password, "Password for the APIC Account. This can also be set as the ACI_PASSWORD environment variable.")
	a.Describe(&c.Insecure, "Allow insecure HTTPS client. This can also be set as the ACI_INSECURE environment variable. Defaults to true.")
	a.Describe(&c.Retries, "Number of retries for REST API calls. This can also be set as the ACI_RETRIES environment variable. Defaults to 3.")
	a.Describe(&c.Logging, "Enable debug logging. This can also be set as the ACI_LOGGING environment variable. Defaults to false.")
	a.SetDefault(&c.Url, "", "ACI_URL")
	a.SetDefault(&c.Username, "", "ACI_USERNAME")
	a.SetDefault(&c.Password, "", "ACI_PASSWORD")

	// not yet supported in pulumi-go-provider
	// a.SetDefault(&c.Insecure, true, "ACI_INSECURE")
	// a.SetDefault(&c.Retries, 3, "ACI_RETRIES")
	// a.SetDefault(&c.Logging, false, "ACI_LOGGING")
}

var _ = (infer.CustomCheck[*Config])((*Config)(nil))

// workaround for https://github.com/pulumi/pulumi-go-provider/issues/110
func (c *Config) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (*Config, []p.CheckFailure, error) {
	c.Url = newInputs["url"].StringValue()
	c.Username = newInputs["username"].StringValue()
	c.Password = newInputs["password"].StringValue()
	if newInputs["insecure"].IsBool() {
		insecure := newInputs["insecure"].BoolValue()
		c.Insecure = &insecure
	}
	if newInputs["retries"].IsNumber() {
		retries := int(newInputs["retries"].NumberValue())
		c.Retries = &retries
	}
	if newInputs["logging"].IsBool() {
		logging := newInputs["logging"].BoolValue()
		c.Logging = &logging
	}

	return c, []p.CheckFailure{}, nil
}

var _ = (infer.CustomConfigure)((*Config)(nil))

func (c *Config) Configure(ctx p.Context) error {
	var insecure bool
	if c.Insecure == nil {
		insecureStr := os.Getenv("ACI_INSECURE")
		if insecureStr == "" {
			insecure = true
		} else {
			insecure, _ = strconv.ParseBool(insecureStr)
		}
	} else {
		insecure = *c.Insecure
	}

	var retries int
	if c.Retries == nil {
		retriesStr := os.Getenv("ACI_RETRIES")
		if retriesStr == "" {
			retries = 3
		} else {
			retries64, _ := strconv.ParseInt(retriesStr, 0, 0)
			retries = int(retries64)
		}
	} else {
		retries = *c.Retries
	}

	var logging bool
	if c.Logging == nil {
		loggingStr := os.Getenv("ACI_LOGGING")
		if loggingStr == "" {
			logging = false
		} else {
			logging, _ = strconv.ParseBool(loggingStr)
		}
	} else {
		logging = *c.Logging
	}

	client, _ := aci.NewClient(c.Url, c.Username, c.Password, aci.Insecure(insecure), aci.MaxRetries(retries), aci.Logging(logging))
	c.Client = &client
	return nil
}
