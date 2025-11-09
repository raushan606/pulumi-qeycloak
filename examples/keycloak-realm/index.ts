import * as pulumi from "@pulumi/pulumi";
import { Provider } from "@raushan606/qeycloak";
import { Realm } from "@raushan606/qeycloak/realm";


const config = new pulumi.Config();
const keycloakUrl = config.get("keycloakUrl") || "http://localhost:8080";
const username = config.get("keycloakUsername") || "admin";
const password = config.getSecret("keycloakPassword") || pulumi.secret("admin");

const keycloakProvider = new Provider("localKeycloak", {
    url: keycloakUrl,
    username: username,
    password: password.apply(pwd => pwd),
    realm: "master",
    insecure: true,
    basePath: "/",
});

const demoRealm = new Realm("qube-realm", {
    name: "payara-qube",
    enabled: true,
    displayName: "Payara Qube",
    displayNameHtml: "<div class=\"kc-logo-text\"><span>Payara Qube</span></div>",
    loginTheme: "payara",
    accountTheme: "payara",
    adminTheme: "payara",
    emailTheme: "payara",
}, { provider: keycloakProvider });

// Export realm id and name
export const realmId = demoRealm.realmId;
export const realmName = demoRealm.name;
export const realm = demoRealm;

// Helpful note output
export const info = pulumi.interpolate`Realm ${demoRealm.name} provisioned with id ${demoRealm.realmId}`;