set -x

OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make terraform VERSION=beta OUTPUT_PATH="$OUTPUT_PATH" PRODUCT=firebasedatabase

cd $OUTPUT_PATH
git status

export GOOGLE_PROJECT=fredzqm-staging
export GOOGLE_USE_DEFAULT_CREDENTIALS=1
export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-c

TF_LOG=TRACE make testacc TEST=./google-beta TESTARGS='-run=TestAccFirebaseDatabaseDatabase' | tee tests.log

