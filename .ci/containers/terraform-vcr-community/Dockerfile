FROM gcr.io/graphite-docker-images/terraform-vcr-tester:latest

# Same image, adds new script to do the opposite of check_membership.sh, skip if user is GCP org member
ADD reverse_check_membership.sh /reverse_check_membership.sh
ENTRYPOINT ["/reverse_check_membership.sh"]
