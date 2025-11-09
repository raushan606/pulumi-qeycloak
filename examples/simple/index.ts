import * as aci from "@netascode/aci";

const tenantName = "TEN-PULUMI1";
const tenant = new aci.apic.Rest(
    tenantName,
    {
        dn: "uni/tn-"+tenantName,
        class_name: "fvTenant",
        content: {
            "name": tenantName,
            "descr": "Tenant created by Pulumi"
        }
    });
