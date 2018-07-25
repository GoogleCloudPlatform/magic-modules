FROM nmckinley/installed-magic-modules:latest

RUN export GOPATH="/go" && \
    mkdir -p "${GOPATH}/src/github.com/terraform-providers"
ADD build/terraform/ /go/src/github.com/terraform-providers/terraform-provider-google/
RUN cd "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google" && \
    go get && \
    rm -rf "${GOPATH}/src/github.com/terraform-providers/terraform-provider-google"
