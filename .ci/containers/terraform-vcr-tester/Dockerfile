FROM gcr.io/graphite-docker-images/magic-modules-base:latest

ADD teamcityparams.xml /teamcityparams.xml
ADD run_vcr_tests.sh /run_vcr_tests.sh
ADD check_membership.sh /check_membership.sh
ADD teamcityparamsrecording.xml /teamcityparamsrecording.xml
ENTRYPOINT ["/check_membership.sh"]
