set -x

OUTPUT_PATH=$GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make terraform VERSION=beta OUTPUT_PATH="$OUTPUT_PATH" PRODUCT=firebasedatabase

cd $OUTPUT_PATH && git status

