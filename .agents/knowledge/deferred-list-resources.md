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
