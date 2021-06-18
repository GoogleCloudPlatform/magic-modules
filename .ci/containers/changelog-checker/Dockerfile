FROM hashicorpdev/go-changelog

ENV GITHUB_REPO=magic-modules
ENV GITHUB_OWNER=GoogleCloudPlatform

RUN apt-get update
RUN apt-get install -y jq

ADD check_changelog.sh /check_changelog.sh
ENTRYPOINT ["/check_changelog.sh"]
