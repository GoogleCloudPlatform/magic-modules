[
  {
    "name": "cloudresourcemanager.googleapis.com/projects/{{.Provider.project}}",
    "asset_type": "cloudresourcemanager.googleapis.com/Project",
    "ancestry_path": "{{.Ancestry}}/project/{{.Provider.project}}",
    "v2_org_policies": [
      {
        "name": "policies/gcp.resourceLocations",
        "spec": {
          "update_time": "{{.Time.RFC3339Nano}}",
          "rules": [
            {
              "values": {
                "allowed_values": [
                  "projects/allowed-project"
                ],
                "denied_values": [
                  "projects/denied-project"
                ]
              },
              "expression": {
                "expression": "resource.matchLabels('labelKeys/123', 'labelValues/345')",
                "title": "sample-condition",
                "description": "A sample condition for the policy",
                "location": "sample-location.log"
              }
            },
            {
              "allow_all": true
            }
          ]
        }
      }
    ],
  }
]
