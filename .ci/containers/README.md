# Making Changes to Build Images
The docker images located in this folder are used by multiple builds for magic modules. They are not automatically updated when the Dockerfile is updated so they must be pushed to gcr.io and tagged by hand.

## Naming Convention

The images are named according to their use. We have a small number of images that get reused in multiple places, based around sets of requirements shared by different parts of the build pipeline. The images are:

- `gcr.io/graphite-docker-images/bash-plus`
- `gcr.io/graphite-docker-images/build-environment`
- `gcr.io/graphite-docker-images/go-plus`

## Updating a docker image

Before you begin, set up Docker (including configuring it to [authenticate with gcloud](https://cloud.google.com/container-registry/docs/advanced-authentication#gcloud-helper)).

1. Make changes to the Dockerfile
2. Build & push the image with the `testing` tag:
   ```bash
   gcloud builds submit . \
    --tag us.gcr.io/graphite-docker-images/go-plus:testing \
    --project graphite-docker-images
   ```
3. Update cloudbuild yaml files to reference the image you just pushed by adding the `:testing` suffix
4. Update files that will cause the cloudbuild yaml changes (and therefore your changes) to be exercised
   - Tip: Modifying `mmv1/third_party/terraform/services/compute/metadata.go.tmpl` will trigger builds for TPG, TPGB, and TGC.
5. Create a PR with these changes.
6. Verify that the cloudbuild steps that should use your testing image _are_ using your testing image (in the Execution Details tab for the step.)
