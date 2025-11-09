package provider

import (
	"strings"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"

	"github.com/pulumi/pulumi-go-provider/integration"
	qeycloak "github.com/raushan606/pulumi-qeycloak/provider/pkg/provider/qeycloak"
)

func NewProvider() p.Provider {
	return infer.Provider(infer.Options{
		// This is the metadata for the provider
		Metadata: schema.Metadata{
			DisplayName: "Qeycloak",
			Description: "The Pulumi Qube Keycloak provider provides resources to interact with a Realm on a Keycloak server.",
			Keywords: []string{
				"pulumi",
				"keycloak",
				"realm",
				"qeycloak",
			},
			Homepage:          "https://pulumi.com",
			License:           "Apache-2.0",
			Repository:        "https://github.com/raushan606/pulumi-qeycloak",
			PluginDownloadURL: "github://api.github.com/raushan606",
			Publisher:         "Raushan Kumar",
			// This contains language specific details for generating the provider's SDKs
			LanguageMap: map[string]any{
				"nodejs": map[string]any{
					"packageName": "@netascode/aci",
					"dependencies": map[string]string{
						"@pulumi/pulumi": "^3.0.0",
					},
				},
				"java": map[string]any{
					"buildFiles":                      "gradle",
					"gradleNexusPublishPluginVersion": "1.1.0",
					"dependencies": map[string]any{
						"com.pulumi:pulumi":               "1.16.3",
						"com.google.code.gson:gson":       "2.13.2",
						"com.google.code.findbugs:jsr305": "3.0.2",
					},
				},
			},
		},
		// A list of `infer.Resource` that are provided by the provider.
		Resources: []infer.InferredResource{
			infer.Resource[*qeycloak.Rest, qeycloak.RestInputs, qeycloak.RestOutputs](),
		},
		Config: infer.Config[*qeycloak.Config](),
	})
}

func Schema(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	s, err := integration.NewServer("aci", semver.MustParse(version), NewProvider()).
		GetSchema(p.GetSchemaRequest{})
	return s.Schema, err
}
