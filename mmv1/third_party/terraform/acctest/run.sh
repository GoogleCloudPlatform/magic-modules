

## outline of what this script does

## 1. set the release_diff env var to have a non nil value
# SERVICE_ACCOUNT_KEY_PATH=../../../../ryanoaksnightly2-4466c5daa7a1.json
export RELEASE_DIFF=true
# export GOOGLE_REGION=us-central1
# export GOOGLE_ZONE=us-central1-a
# export ORG_ID=280476229921
# export GOOGLE_PROJECT=ryanoaksnightly2
# export GOOGLE_ORG=280476229921
# export GOOGLE_CUST_ID=C02jqjwhe
# export GOOGLE_ORG_DOMAIN=ryanoakstestco.joonix.net
# export GOOGLE_PROJECT_NUMBER=1011002275372
# export GOOGLE_IDENTITY_USER=ryantest
# export GOOGLE_BILLING_ACCOUNT=01129A-55361F-811C70
# # export GCLOUD_KEYFILE_JSON=$SERVICE_ACCOUNT_KEY_PATH
# export GCLOUD_SERVICE_ACCOUNT_KEY_PATH=$SERVICE_ACCOUNT_KEY_PATH
# #export GOOGLE_USE_DEFAULT_CREDENTIALS=true
# export TF_ACC=true
export GOOGLE_REGION=us-central1
export GOOGLE_ZONE=us-central1-a
export ORG_ID=280476229921
export GOOGLE_PROJECT=ryanoaksnightly2
export GOOGLE_ORG=280476229921
export GOOGLE_CUST_ID=C02jqjwhe
export GOOGLE_ORG_DOMAIN=ryanoakstestco.joonix.net
export GOOGLE_PROJECT_NUMBER=1011002275372
export GOOGLE_IDENTITY_USER=ryantest
export GOOGLE_BILLING_ACCOUNT=01129A-55361F-811C70
export GOOGLE_USE_DEFAULT_CREDENTIALS=true
export TF_ACC=true
## 2. github actions to get to correct provider(irrelevant for now)

##3. run make test command for something(hardcoded for now)

declare -a TEST_COMMANDS

TEST_COMMANDS+=( "make testacc TEST=./google/services/alloydb TESTARGS='-run=TestAccAlloydbCluster_alloydbClusterBasicExample'" )
TEST_COMMANDS+=( "make testacc TEST=./google/services/bigtable TESTARGS='-run=TestAccBigtableInstance_basic'" )
TEST_COMMANDS+=( "make testacc TEST=./google/services/alloydb TESTARGS='-run=TestAccAlloydbCluster_withMaintenanceWindowsMissingFields'" )
TEST_COMMANDS+=( "make testacc TEST=./google/services/bigquery TESTARGS='-run=TestAccDataSourceGoogleBigqueryDefaultServiceAccount_basic'" )
OVERALL_STATUS=0

echo "--- Running Terraform Acceptance Tests (Sequentially) ---"

# 3. Run make test command for something
for i in "${!TEST_COMMANDS[@]}"; do
    TEST_CMD="${TEST_COMMANDS[$i]}"
    LOG_FILE="output_${i}.log"
    TEST_NAME=$(echo "$TEST_CMD" | sed -n "s/.*TESTARGS='-run=\([^']*\)'.*/\1/p")

    echo "Running test: $TEST_NAME (logging to $LOG_FILE)"
    $TEST_CMD > "$LOG_FILE" 2>&1
    if [ $? -ne 0 ]; then
        echo "FAIL: Test '$TEST_NAME' failed. Check $LOG_FILE for details."
        OVERALL_STATUS=1
    else
        echo "PASS: Test '$TEST_NAME' completed successfully."
    fi
done

echo -e "\n--- Parsing Test Outputs for Diffs ---"
for i in "${!TEST_COMMANDS[@]}"; do
    LOG_FILE="output_${i}.log"
    TEST_NAME=$(echo "${TEST_COMMANDS[$i]}" | sed -n "s/.*TESTARGS='-run=\([^']*\)'.*/\1/p")
    echo "Diffs for $TEST_NAME:"
    grep -E "Error running|Error: Invalid provider configuration|error unmarshaling credentials|\[Diff\]" "$LOG_FILE" || echo "  No specific diff flags found in $TEST_NAME log."
done

## 4. now that we have an output file(given changes to vcr), parse it and run grep FLAG to get all diff tests

## 5. print those diffs

unset release_diff