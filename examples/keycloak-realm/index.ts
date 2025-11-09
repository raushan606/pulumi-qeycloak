import * as pulumi from "@pulumi/pulumi";
import { Provider, Realm } from "@raushan606/qeycloak";

// Configuration / assumptions:
// We assume a local Keycloak instance listening at http://localhost:8080
// For demo purposes we use arbitrary credentials; adjust to your local setup.
// If you have different admin credentials, either change them here or set via stack config.

// Allow overriding via Pulumi config (optional convenience):
const config = new pulumi.Config();
const keycloakUrl = config.get("keycloakUrl") || "http://localhost:8080"; // include protocol for clarity
const username = config.get("keycloakUsername") || "demo-admin";
const password = config.getSecret("keycloakPassword") || pulumi.secret("demo-password-ChangeMe!");

// Instantiate the provider. Defaults (realm=master, insecure=true, basePath=/) are applied by the provider.
const keycloakProvider = new Provider("localKeycloak", {
    url: keycloakUrl,
    username: username,
    password: password,
    // realm: "master", // optional override
    // insecure: true,  // default true
});

// Create a new realm using the provider.
const demoRealm = new Realm("demoRealm", {
    name: "pulumi-demo-realm",
    displayName: "Pulumi Demo Realm",
    displayNameHtml: "<strong>Pulumi Demo Realm</strong>",
    enabled: true,
    loginTheme: "keycloak",    // assuming default themes
}, { provider: keycloakProvider });

// Export realm id and name
export const realmId = demoRealm.realmId;
export const realmName = demoRealm.name;

// Helpful note output
export const info = pulumi.interpolate`Realm ${demoRealm.name} provisioned with id ${demoRealm.realmId}`;