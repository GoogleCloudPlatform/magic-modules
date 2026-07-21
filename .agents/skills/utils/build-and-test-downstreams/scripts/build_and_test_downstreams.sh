#!/usr/bin/env bash
#
# build_and_test_downstreams.sh
#
# Generates downstream provider code from Magic Modules into scratch directories
# and runs build, unit tests, lint, and docscheck in parallel.
#

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Downstream Provider Build & Unit Test Runner ===${NC}"

# 1. Validate environment
if [ ! -d "mmv1" ]; then
  echo -e "${RED}Error: This script must be run from the root of the magic-modules repository.${NC}"
  exit 1
fi

REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/downstream-build"
LOGS_DIR="${SCRATCH_DIR}/logs"
mkdir -p "$SCRATCH_DIR" "$LOGS_DIR"

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

# Function to setup, generate, build, and test a downstream provider version
run_provider_build_and_test() {
  local VERSION="$1"
  local VERSION_UPPER=$(echo "$VERSION" | tr '[:lower:]' '[:upper:]')
  local PROVIDER_REPO
  local PROVIDER_CACHE
  local DEST_DIR="${SCRATCH_DIR}/downstream-${VERSION}"
  local LOG_FILE="${LOGS_DIR}/build-test-${VERSION}.log"

  exec > "$LOG_FILE" 2>&1

  if [ "$VERSION" = "ga" ]; then
    PROVIDER_REPO="https://github.com/hashicorp/terraform-provider-google.git"
    PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-ga"
  elif [ "$VERSION" = "beta" ]; then
    PROVIDER_REPO="https://github.com/hashicorp/terraform-provider-google-beta.git"
    PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-beta"
  elif [ "$VERSION" = "docs-examples" ]; then
    PROVIDER_REPO="https://github.com/terraform-google-modules/docs-examples.git"
    PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-docs-examples"
  fi

  echo "=== [${VERSION_UPPER}] Starting Downstream Build & Test ==="

  # 1. Ensure provider cache is up to date
  if [ ! -d "$PROVIDER_CACHE" ]; then
    echo "[${VERSION_UPPER}] Cloning downstream provider repository..."
    git clone --depth 1 "$PROVIDER_REPO" "$PROVIDER_CACHE"
  else
    echo "[${VERSION_UPPER}] Updating cached downstream provider..."
    (
      cd "$PROVIDER_CACHE"
      git fetch --depth 1 origin main 2>/dev/null || git fetch --depth 1 origin master 2>/dev/null || true
      git reset --hard FETCH_HEAD 2>/dev/null || true
    )
  fi

  # 2. Copy cache to dest dir
  rm -rf "$DEST_DIR"
  cp -R "$PROVIDER_CACHE" "$DEST_DIR"

  # 3. Generate provider code
  echo "[${VERSION_UPPER}] Generating provider code from Magic Modules..."
  if [ "$VERSION" = "docs-examples" ]; then
    (
      cd mmv1
      "$MM_BINARY" --version ga --provider oics --output "$DEST_DIR"
    )
  else
    (
      cd mmv1
      if [ "$VERSION" = "ga" ]; then
        "$MM_BINARY" --output "$DEST_DIR" --version ga --no-docs && "$MM_BINARY" --output "$DEST_DIR" --version beta --no-code
      else
        "$MM_BINARY" --output "$DEST_DIR" --version "$VERSION"
      fi
    )
  fi

  # 4. Build and run tests
  if [ "$VERSION" != "docs-examples" ]; then
    echo "[${VERSION_UPPER}] Compiling provider binary..."
    (
      cd "$DEST_DIR"
      go build
      echo "[${VERSION_UPPER}] Running provider unit tests..."
      make testnolint
      echo "[${VERSION_UPPER}] Running provider lint..."
      make lint
      echo "[${VERSION_UPPER}] Running docs check..."
      make docscheck
    )
  fi

  echo "=== [${VERSION_UPPER}] Downstream Build & Test Completed Successfully ==="
}

echo -e "${BLUE}Running GA, Beta, and docs-examples downstream generation and unit tests in parallel...${NC}"

run_provider_build_and_test ga &
GA_PID=$!

run_provider_build_and_test beta &
BETA_PID=$!

run_provider_build_and_test docs-examples &
DOCS_PID=$!

set +e
wait $GA_PID
GA_STATUS=$?

wait $BETA_PID
BETA_STATUS=$?

wait $DOCS_PID
DOCS_STATUS=$?
set -e

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== GA Provider Build & Test Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/build-test-ga.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Beta Provider Build & Test Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/build-test-beta.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Docs Examples Generation Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/build-test-docs-examples.log"

if [ $GA_STATUS -ne 0 ] || [ $BETA_STATUS -ne 0 ] || [ $DOCS_STATUS -ne 0 ]; then
  echo -e "\n${RED}=== Downstream provider build and unit tests failed! ===${NC}"
  exit 1
else
  echo -e "\n${GREEN}=== All downstream provider builds and unit tests passed successfully! ===${NC}"
  exit 0
fi
