# TODO(fredzqm): remove before sending PR for review
set -x
set -e

OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
#make terraform VERSION=beta OUTPUT_PATH="$OUTPUT_PATH" PRODUCT=firebase
make terraform VERSION=beta OUTPUT_PATH="$OUTPUT_PATH" PRODUCT=firebasedatabase

cd $OUTPUT_PATH
git status

export GOOGLE_PROJECT=fredzqm-staging-b
export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-c
export GOOGLE_USE_DEFAULT_CREDENTIALS=True
export GOOGLE_IMPERSONATE_SERVICE_ACCOUNT=terraform-tester@fredzqm-staging.iam.gserviceaccount.com
# https://console.cloud.google.com/iam-admin/iam?authuser=2&organizationId=641623452985
export GOOGLE_ORG=641623452985
# export GOOGLE_BILLING_ACCOUNT=261046259366

#export TF_LOG=TRACE
make testacc TEST=./google-beta TESTARGS='-run=TestAccFirebaseDatabaseInstance' | tee tests.log
make testacc TEST=./google-beta TESTARGS='-sweep=us-central1 -sweep-run=FirebaseDatabaseInstance' > output.log

