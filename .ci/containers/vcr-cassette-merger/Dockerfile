from gcr.io/google.com/cloudsdktool/cloud-sdk:alpine as resource

RUN apk add --no-cache bash
RUN apk add --no-cache curl
RUN apk add --no-cache jq

ADD vcr_merge.sh /vcr_merge.sh
ENTRYPOINT ["/vcr_merge.sh"]
