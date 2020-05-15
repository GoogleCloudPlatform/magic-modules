FROM alpine

RUN apk add --no-cache curl
RUN apk add --no-cache jq

ADD teamcityparams.xml /teamcityparams.xml
ADD vcr_test_terraform.sh /vcr_test_terraform.sh
ENTRYPOINT ["/vcr_test_terraform.sh"]