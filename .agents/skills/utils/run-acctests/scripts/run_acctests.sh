#!/usr/bin/env bash
#
# run_acctests.sh
#
# Generates, builds, and runs acceptance tests for Beta, then GA sequentially.
# Short-circuits immediately if Beta tests fail.
#
# Usage:
#   ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh <SERVICE_OR_TEST_PATH> [TEST_NAME]
#
# Examples:
#   ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh compute TestAccComputeInstance_basic
#   ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh ./google-beta/services/compute TestAccComputeInstance_basic
#

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Acceptance Test Runner (Beta then GA) ===${NC}"

if [ ! -d "mmv1" ]; then
  echo -e "${RED}Error: This script must be run from the root of the magic-modules repository.${NC}"
  exit 1
fi

FIRST_ARG="$1"
SECOND_ARG="$2"
THIRD_ARG="$3"

if [ "$FIRST_ARG" = "beta" ] || [ "$FIRST_ARG" = "ga" ]; then
  TEST_TARGET="$SECOND_ARG"
  TEST_NAME="$THIRD_ARG"
else
  TEST_TARGET="$FIRST_ARG"
  TEST_NAME="$SECOND_ARG"
fi

if [ -z "$TEST_TARGET" ]; then
  echo -e "${RED}Error: Service or test path is required.${NC}"
  echo "Usage: $0 <service_name_or_path> [test_name_pattern]"
  exit 1
fi

REPO_ROOT=$(pwd)

# Build mmv1 binary once
echo -e "${BLUE}Building mmv1 binary...${NC}"
if [ -f "MODULE.bazel" ] && command -v bazel &> /dev/null; then
  bazel build //mmv1
  MM_BINARY="${REPO_ROOT}/bazel-bin/mmv1/mmv1_/mmv1"
else
  (
    cd mmv1
    go build -o ../bin/mmv1 .
  )
  MM_BINARY="${REPO_ROOT}/bin/mmv1"
fi

run_acctest_for_version() {
  local VERSION="$1"
  local VERSION_UPPER=$(echo "$VERSION" | tr '[:lower:]' '[:upper:]')
  local PROVIDER_REPO
  local PROVIDER_CACHE
  local FOLDER_NAME
  local SCRATCH_DIR="${REPO_ROOT}/scratch/acctest-${VERSION}"

  if [ "$VERSION" = "ga" ]; then
    PROVIDER_REPO="https://github.com/hashicorp/terraform-provider-google.git"
    PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-ga"
    FOLDER_NAME="google"
  else
    VERSION="beta"
    PROVIDER_REPO="https://github.com/hashicorp/terraform-provider-google-beta.git"
    PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-beta"
    FOLDER_NAME="google-beta"
  fi

  # Resolve service path for this version
  local SERVICE_PATH
  if [[ "$TEST_TARGET" == ./* ]] || [[ "$TEST_TARGET" == services/* ]]; then
    SERVICE_PATH=$(echo "$TEST_TARGET" | sed -E "s/google(-beta)?/google/g")
    if [ "$VERSION" = "beta" ]; then
      SERVICE_PATH=$(echo "$SERVICE_PATH" | sed -E "s/\.\/google\//\.\/google-beta\//g")
    fi
  else
    SERVICE_PATH="./${FOLDER_NAME}/services/${TEST_TARGET}"
  fi

  echo -e "\n${BLUE}=====================================================${NC}"
  echo -e "${BLUE}=== Starting ${VERSION_UPPER} Acceptance Test ===${NC}"
  echo -e "${BLUE}=====================================================${NC}"
  echo -e "${BLUE}Version: ${YELLOW}${VERSION}${NC}"
  echo -e "${BLUE}Service Path: ${YELLOW}${SERVICE_PATH}${NC}"
  if [ -n "$TEST_NAME" ]; then
    echo -e "${BLUE}Test Name Pattern: ${YELLOW}${TEST_NAME}${NC}"
  fi

  # 1. Ensure provider cache exists
  if [ ! -d "$PROVIDER_CACHE" ]; then
    echo -e "${BLUE}[${VERSION_UPPER}] Cloning downstream provider into cache...${NC}"
    git clone --depth 1 "$PROVIDER_REPO" "$PROVIDER_CACHE"
  else
    echo -e "${BLUE}[${VERSION_UPPER}] Updating cached downstream provider...${NC}"
    (
      cd "$PROVIDER_CACHE"
      git fetch --depth 1 origin main 2>/dev/null || git fetch --depth 1 origin master 2>/dev/null || true
      git reset --hard FETCH_HEAD 2>/dev/null || true
    )
  fi

  # 2. Copy cache to scratch & generate provider
  echo -e "${BLUE}[${VERSION_UPPER}] Generating provider code into scratch directory...${NC}"
  rm -rf "$SCRATCH_DIR"
  cp -R "$PROVIDER_CACHE" "$SCRATCH_DIR"

  (
    cd mmv1
    if [ "$VERSION" = "ga" ]; then
      "$MM_BINARY" --output "$SCRATCH_DIR" --version ga --no-docs && "$MM_BINARY" --output "$SCRATCH_DIR" --version beta --no-code
    else
      "$MM_BINARY" --output "$SCRATCH_DIR" --version "$VERSION"
    fi
  )

  # 3. Build downstream provider
  echo -e "${BLUE}[${VERSION_UPPER}] Building downstream provider binary...${NC}"
  (
    cd "$SCRATCH_DIR"
    make build
  )

  # 4. Execute Acceptance Test
  echo -e "${BLUE}[${VERSION_UPPER}] Running acceptance test in ${SCRATCH_DIR}...${NC}"
  local LOGS_DIR="${SCRATCH_DIR}/logs"
  mkdir -p "$LOGS_DIR"
  local LOG_FILE="${LOGS_DIR}/test_output_${VERSION}.log"

  set +e
  (
    cd "$SCRATCH_DIR"
    if [ -n "$TEST_NAME" ]; then
      TF_LOG=DEBUG make testacc TEST="$SERVICE_PATH" TESTARGS="-run=${TEST_NAME}\$\$" > "$LOG_FILE" 2>&1
    else
      TF_LOG=DEBUG make testacc TEST="$SERVICE_PATH" > "$LOG_FILE" 2>&1
    fi
  )
  local TEST_EXIT_CODE=$?
  set -e

  if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}=== [${VERSION_UPPER}] Acceptance test passed successfully! ===${NC}"
    return 0
  else
    echo -e "${RED}=== [${VERSION_UPPER}] Acceptance test failed! ===${NC}"
    echo -e "${YELLOW}Logs saved to: ${LOG_FILE}${NC}"
    echo -e "${YELLOW}Use the parse-debug-logs skill to analyze ${LOG_FILE}${NC}"
    return $TEST_EXIT_CODE
  fi
}

# Step 1: Run Beta Acceptance Test
run_acctest_for_version beta
BETA_STATUS=$?

if [ $BETA_STATUS -ne 0 ]; then
  echo -e "\n${RED}=== Beta acceptance test failed! Short-circuiting execution (skipping GA test). ===${NC}"
  exit $BETA_STATUS
fi

# Step 2: Run GA Acceptance Test
run_acctest_for_version ga
GA_STATUS=$?

if [ $GA_STATUS -ne 0 ]; then
  echo -e "\n${RED}=== GA acceptance test failed! ===${NC}"
  exit $GA_STATUS
fi

echo -e "\n${GREEN}=== Sequential Acceptance Tests (Beta & GA) passed successfully! ===${NC}"
exit 0
