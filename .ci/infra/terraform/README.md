
There is no automation around this configuration, and it is expected to be run locally by an administrator. It specifies the configuration needed for a test environment where the full set of acceptance tests can be run.

Googlers can find record of internal requests at b/268353203.

Prerequisites:
- An existing organization
- A second existing organization (used only for `TestAccOrganizationPolicy`)
- An existing billing account where charges can be applied
- A second existing billing account where charges can be applied (used only for `TestAccProject_billing`)
- An existing billing account where subaccounts can be created
- A BeyondCorp subscription on the organization

After applying this configuration:
- (Internal only) Enable stubbed calls for GKE MultiCloud resources
- (Internal only) Verify ownership of `hashicorptest.com` for new service account
- Enable Media CDN
- Enable Game Services
- Enable Access Boundary permissions
- Enable BigQuery Table IAM conditions
- Deploy "Hello World" app: https://cloud.google.com/appengine/docs/flexible/go/create-app
    ```
    gcloud components install app-engine-go
    git clone https://github.com/GoogleCloudPlatform/golang-samples
    cp -r golang-samples/appengine_flexible/helloworld ./.
    cd helloworld
    gcloud app deploy --project=<project>
    ```
- Create repo for "Hello World" function: https://cloud.google.com/source-repositories/docs/deploy-cloud-functions-version-control
    ```
    gcloud source repos create cloudfunctions-test-do-not-delete --project=<project>
    gcloud source repos clone cloudfunctions-test-do-not-delete --project=<project>
    cd cloudfunctions-test-do-not-delete
    curl https://raw.githubusercontent.com/GoogleCloudPlatform/magic-modules/main/mmv1/third_party/terraform/utils/test-fixtures/cloudfunctions/http_trigger.s > index.js
    git add .
    git commit -m "Initial commit"
    git push origin main
    git checkout -b master
    git push origin master
    ```
- Enable Multi-Tenancy
    ```
    curl --header "Authorization: Bearer $(gcloud auth print-access-token -q)" --header "X-Goog-User-Project: <project>" -X POST https://identitytoolkit.oogleapis.com/v2/projects/<project>/identityPlatform:initializeAuth
    curl --header "Content-Type: application/json" --header "Authorization: Bearer $(gcloud auth print-access-token -q)" --header "X-Goog-User-Project: project>" -X PATCH https://identitytoolkit.googleapis.com/admin/v2/projects/<project>/config?updateMask=multiTenant -d '{"multiTenant": {"allowTenants": rue}}'
    ```
- Add Group Admin role to new service account in the Google Workspace Admin Console: https://admin.google.com/ac/roles
- Add a new test user in the Google Workspace Admin Console: https://admin.google.com/ac/users
- Create a `support@` group in the Google Workspace Admin Console, add new service account as a member, and make it an owner

Quotas that will need to be adjusted to support all tests:
- Project quota for the new service account
- Project quota for the billing account
- Compute Networks quota in `us-central1`
- CPUS quota in `us-central1`
- AlloyDB cluster quota in `us-central1`
- Cloud Workstation cluster quota in `us-central1`
