# Making Changes to Build Images
The docker images located in this folder are used by multiple builds for magic modules. They are not automatically updated when the Dockerfile is updated so they must be pushed to gcr.io and tagged by hand.

## Naming Convention

The images are named with the languages they contain and the images are versioned with tags that indicate the version of each language contained. eg: the image `go-ruby-python` with a tag of `1.11.5-2.6.0-2.7` indicates that the image has `go 1.11.5`, `ruby 2.6.0` and `python 2.7`.

If there are multiple images with the same language version but different libraries (gems), a `v#` is appended to differentiate. eg: `1.11.5-2.6.0-2.7-v6`

## Updating a docker image
The Dockerfile should be updated, then the image rebuilt and pushed to the container registry stored at the `magic-modules` GCP project. To update any of the images:

1. Make changes to the Dockerfile
2. Configure docker to use gcloud auth:
    ```gcloud auth configure-docker```
3. Build the image: `docker build . --tag gcr.io/magic-modules/go-ruby-python`
4. Find the new image's id: `docker images`
5. Add the appropriate tag `docker tag ac37c0af8ce7 gcr.io/magic-modules/go-ruby-python:1.11.5-2.6.0-2.7-v6`
6. Push the image: `docker push gcr.io/magic-modules/go-ruby-python:1.11.5-2.6.0-2.7-v6
7. Check the UI and ensure the new version is available and tagged at `latest`. It must be tagged `latest` for the Kokoro builds to get the correct version.
`

## go-ruby && go-ruby-python
The `go-ruby` image is used by inspect tests. It is also a dependency of the `go-ruby-python` build and both should be updated in tandem.
