# Deferred List Resources

Resources that were scoped out of a prior PR and need a dedicated follow-up PR.
When the agent runs without `TARGET_PRODUCT` set, entries here take priority over
selecting a new product.

Each row: **product**, **resource** (PascalCase YAML stem), **pattern** (oracle P-NN),
**reason** (one sentence), **follow-up branch** (exact git branch name).

## Table

| Product | Resource | Pattern | Reason | Follow-up branch |
|---------|----------|---------|--------|------------------|
| apigee | TargetServer | P-08 | Bare-array list response — generator needs `list_response_is_array: true` support | add-apigee-list-resources-followup |
| apigee | EnvironmentKeyvaluemaps | P-08 | Bare-array list response — generator needs `list_response_is_array: true` support | add-apigee-list-resources-followup |
| apigee | Organization | P-09 | List items use `"organization"` identity key but resource uses `"name"` — needs custom list decoder | add-apigee-list-resources-followup |
| networkservices | AgentGateway | P-17 | `base_url` defaults to `global` location; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | EdgeCacheKeyset | P-17 | `base_url` hardcodes `/locations/global/`; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | EdgeCacheOrigin | P-17 | `base_url` hardcodes `/locations/global/`; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | EdgeCacheService | P-17 | `base_url` hardcodes `/locations/global/`; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | EndpointPolicy | P-17 | `base_url` hardcodes `/locations/global/`; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | Gateway | P-17 | `base_url` defaults to `global` location; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | GrpcRoute | P-17 | `base_url` defaults to `global` location; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | HttpRoute | P-17 | `base_url` hardcodes `/locations/global/`; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | LbRouteExtension | P-17 | `base_url` uses `{{location}}` but resource creation scope mismatches list scope | add-networkservices-list-resources-followup |
| networkservices | LbTrafficExtension | P-17 | `base_url` uses `{{location}}` but resource creation scope mismatches list scope | add-networkservices-list-resources-followup |
| networkservices | Mesh | P-17 | `base_url` defaults to `global` location; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | ServiceLbPolicies | P-17 | `base_url` defaults to `global` location; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | TcpRoute | P-17 | `base_url` hardcodes `/locations/global/`; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| networkservices | TlsRoute | P-17 | `base_url` defaults to `global` location; list scope mismatches resource creation scope | add-networkservices-list-resources-followup |
| oracledatabase | AutonomousDatabase | P-22 | ListQuery create uses 240m insert timeout; VCR 6h wall clock cannot record it with sibling tests | add-oracledatabase-list-resources-followup |
| oracledatabase | CloudExadataInfrastructure | P-22 | ListQuery create uses 240m insert timeout; VCR 6h wall clock cannot record it | add-oracledatabase-list-resources-followup |
| oracledatabase | CloudVmCluster | P-22 | ListQuery create uses 120m insert timeout; VCR 6h wall clock cannot record it | add-oracledatabase-list-resources-followup |
| oracledatabase | DbSystem | P-22 | ListQuery create uses 120m insert timeout; VCR 6h wall clock cannot record it | add-oracledatabase-list-resources-followup |
| oracledatabase | ExadbVmCluster | P-22 | ListQuery create uses 120m insert timeout; VCR 6h wall clock cannot record it | add-oracledatabase-list-resources-followup |
| oracledatabase | ExascaleDbStorageVault | P-22 | ListQuery create uses 120m insert timeout; VCR 6h wall clock cannot record it | add-oracledatabase-list-resources-followup |
| oracledatabase | GoldengateDeployment | P-22 | ListQuery create ~78m plus sibling tests exceeded VCR 6h recording wall clock | add-oracledatabase-list-resources-followup |
| oracledatabase | GoldengateConnectionAssignment | P-23 | ListQuery reuses FullExample's shared permanent connection+deployment under t.Parallel(); VCR recording failed while FullExample passed | add-oracledatabase-list-resources-followup |
