FROM alpine

RUN apk add --no-cache bash
RUN apk add --no-cache curl
RUN apk add --no-cache git
RUN apk add --no-cache jq

ADD teamcityparams.xml /teamcityparams.xml
ADD run_vcr_tests.sh /run_vcr_tests.sh
ADD check_membership.sh /check_membership.sh
ADD teamcityparamsrecording.xml /teamcityparamsrecording.xml
ENTRYPOINT ["/check_membership.sh"]
