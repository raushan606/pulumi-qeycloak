from pulumi_aci import *

tenant_name = "TEN_PUL1"
tenant = apic.Rest(
    tenant_name,
    dn="uni/tn-"+tenant_name,
    class_name="fvTenant",
    content={
        "name": tenant_name,
        "descr": "My desc new",
    },
)