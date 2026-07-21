#!/usr/bin/env bash
#
# run_acctests.sh
#
# Generates and builds downstream provider in scratch directory, then runs acceptance tests.
#
# Usage:
#   ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh [VERSION] [SERVICE_OR_TEST_PATH] [TEST_NAME]
#
# Examples:
#   ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh beta compute TestAccComputeInstance_basic
#   ./.agents/skills/utils/run-acctests/scripts/run_acctests.sh beta ./google-beta/services/compute TestAccComputeInstance_basic
#

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Magic Modules Acceptance Test Runner ===${NC}"

if [ ! -d "mmv1" ]; then
  echo -e "${RED}Error: This script must be run from the root of the magic-modules repository.${NC}"
  exit 1
fi

VERSION="${1:-beta}"
TEST_TARGET="$2"
TEST_NAME="$3"

if [ -z "$TEST_TARGET" ]; then
  echo -e "${RED}Error: Service or test path is required.${NC}"
  echo "Usage: $0 [beta|ga] <service_name_or_path> [test_name_pattern]"
  exit 1
fi

REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/acctest-${VERSION}"

# Resolve provider repo & cache
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

# Resolve service path
if [[ "$TEST_TARGET" == ./* ]] || [[ "$TEST_TARGET" == services/* ]]; then
  SERVICE_PATH="$TEST_TARGET"
else
  SERVICE_PATH="./${FOLDER_NAME}/services/${TEST_TARGET}"
fi

echo -e "${BLUE}Version: ${YELLOW}${VERSION}${NC}"
echo -e "${BLUE}Service Path: ${YELLOW}${SERVICE_PATH}${NC}"
if [ -n "$TEST_NAME" ]; then
  echo -e "${BLUE}Test Name Pattern: ${YELLOW}${TEST_NAME}${NC}"
fi

# 1. Ensure provider cache exists
if [ ! -d "$PROVIDER_CACHE" ]; then
  echo -e "${BLUE}Cloning downstream provider into cache...${NC}"
  git clone --depth 1 "$PROVIDER_REPO" "$PROVIDER_CACHE"
else
  echo -e "${BLUE}Updating cached downstream provider...${NC}"
  (
    cd "$PROVIDER_CACHE"
    git fetch --depth 1 origin main 2>/dev/null || git fetch --depth 1 origin master 2>/dev/null || true
    git reset --hard FETCH_HEAD 2>/dev/null || true
  )
fi

# 2. Build mmv1 binary
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

# 3. Copy cache to scratch & generate provider
echo -e "${BLUE}Generating provider code into scratch directory...${NC}"
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

# 4. Build downstream provider
echo -e "${BLUE}Building downstream provider binary...${NC}"
(
  cd "$SCRATCH_DIR"
  make build
)

# 5. Execute Acceptance Test
echo -e "${BLUE}Running acceptance test in ${SCRATCH_DIR}...${NC}"
LOG_FILE="${REPO_ROOT}/scratch/test_output.log"

set +e
(
  cd "$SCRATCH_DIR"
  if [ -n "$TEST_NAME" ]; then
    TF_LOG=DEBUG make testacc TEST="$SERVICE_PATH" TESTARGS="-run=${TEST_NAME}\$\$" > "$LOG_FILE" 2>&1
  else
    TF_LOG=DEBUG make testacc TEST="$SERVICE_PATH" > "$LOG_FILE" 2>&1
  fi
)
TEST_EXIT_CODE=$?
set -e

if [ $TEST_EXIT_CODE -eq 0 ]; then
  echo -e "${GREEN}=== Acceptance test passed successfully! ===${NC}"
  exit 0
else
  echo -e "${RED}=== Acceptance test failed! ===${NC}"
  echo -e "${YELLOW}Logs saved to: ${LOG_FILE}${NC}"
  echo -e "${YELLOW}Use the parse-debug-logs skill to analyze ${LOG_FILE}${NC}"
  exit $TEST_EXIT_CODE
fi
