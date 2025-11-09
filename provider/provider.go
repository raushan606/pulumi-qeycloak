package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/raushan606/pulumi-qeycloak/provider/pkg/qeycloak/config"
	"github.com/raushan606/pulumi-qeycloak/provider/pkg/qeycloak/realm"
)

var Version string = "0.0.1"

const Name string = "qeycloak"

func Provider() p.Provider {
	return infer.Provider(infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource(&realm.Realm{}),
		},
		Components: []infer.InferredComponent{},
		Config:     infer.Config(&config.KeycloakConfig{}),
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
		// This is the metadata for the provider
		Metadata: schema.Metadata{
			DisplayName: "Qeycloak",
			Description: "The Pulumi Qube Keycloak provider provides resources to interact with a Realm on a Keycloak server.",
			Keywords: []string{
				"pulumi",
				"keycloak",
				"realm",
				"qeycloak",
				"kind/native",
			},
			Homepage:   "https://pulumi.com",
			License:    "Apache-2.0",
			Repository: "https://github.com/raushan606/pulumi-qeycloak",
			Publisher:  "Raushan Kumar",
			// This contains language specific details for generating the provider's SDKs
			LanguageMap: map[string]any{
				"nodejs": map[string]any{
					"packageName": "@raushan606/qeycloak",
					"dependencies": map[string]string{
						"@pulumi/pulumi": "^3.0.0",
					},
				},
				"java": map[string]any{
					"buildFiles": "maven",
					"dependencies": map[string]any{
						"com.pulumi:pulumi":               "1.16.3",
						"com.google.code.gson:gson":       "2.13.2",
						"com.google.code.findbugs:jsr305": "3.0.2",
					},
				},
			},
		},
	})
}
