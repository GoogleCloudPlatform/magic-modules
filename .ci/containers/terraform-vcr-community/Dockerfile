FROM alpine

RUN apk add --no-cache bash
RUN apk add --no-cache curl
RUN apk add --no-cache jq

ADD teamcityparams.xml /teamcityparams.xml
ADD run_vcr_tests.sh /run_vcr_tests.sh
ADD teamcityparamsrecording.xml /teamcityparamsrecording.xml
ENTRYPOINT ["/run_vcr_tests.sh"]
