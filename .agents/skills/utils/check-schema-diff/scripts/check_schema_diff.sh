#!/usr/bin/env bash
#
# check_schema_diff.sh
#
# Automates local schema difference checks (breaking changes, missing tests,
# and missing documentation) for Magic Modules by comparing the current
# working tree against a base branch (e.g., main).
#
# Usage:
#   ./.agents/skills/utils/check-schema-diff/scripts/check_schema_diff.sh [base_ref]
#

set -e

# Define color codes for pretty printing
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Magic Modules Schema Diff Checker ===${NC}"

# 1. Validate environment
if [ ! -d "mmv1" ] || [ ! -d "tools/diff-processor" ]; then
  echo -e "${RED}Error: This script must be run from the root of the magic-modules repository.${NC}"
  exit 1
fi

# 2. Determine base ref
BASE_REF="$1"
if [ -z "$BASE_REF" ]; then
  # Try to find the merge base with origin/main, fallback to main, fallback to origin/master
  BASE_REF=$(git merge-base HEAD origin/main 2>/dev/null || git merge-base HEAD main 2>/dev/null || echo "main")
fi

echo -e "${BLUE}Comparing current changes against base ref: ${YELLOW}${BASE_REF}${NC}"

# 3. Set up paths
REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/schema-diff-check"
PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache"
DIFF_PROCESSOR_DIR="${REPO_ROOT}/tools/diff-processor"

mkdir -p "$SCRATCH_DIR"

# 4. Set up cleanup trap to ensure we never leave git worktrees behind
cleanup() {
  echo -e "${BLUE}Cleaning up temporary files and worktrees...${NC}"
  if [ -d "$SCRATCH_DIR/mm-base-worktree" ]; then
    git worktree remove --force "$SCRATCH_DIR/mm-base-worktree" 2>/dev/null || true
  fi
}
trap cleanup EXIT

# Helper function to check if the docs JSON output has missing docs
has_missing_docs() {
  local json="$1"
  if command -v jq &> /dev/null; then
    # Returns 0 (true) if there are missing docs, 1 (false) otherwise
    echo "$json" | jq -e '(.Resource | length > 0) or (.DataSource | length > 0)' >/dev/null 2>&1
    return $?
  else
    # Simple Python fallback to check if either list is non-empty
    python3 -c "
import sys, json
data = json.loads(sys.argv[1])
resources = data.get('Resource') or []
datasources = data.get('DataSource') or []
if len(resources) > 0 or len(datasources) > 0:
    sys.exit(0)
sys.exit(1)
" "$json" 2>/dev/null
    return $?
  fi
}

# 5. Build current mmv1 binary in the current context (once for all checks)
echo -e "${BLUE}Building current mmv1 binary...${NC}"
if [ -f "MODULE.bazel" ] && command -v bazel &> /dev/null; then
  bazel build //mmv1
  MM_BINARY="$(pwd)/bazel-bin/mmv1/mmv1_/mmv1"
else
  (
    cd mmv1
    go build -o ../bin/mmv1 .
  )
  MM_BINARY="$(pwd)/bin/mmv1"
fi

# 6. Create temporary worktree for base commit (once for all checks)
echo -e "${BLUE}Creating temporary worktree for base commit ${YELLOW}${BASE_REF}${BLUE}...${NC}"
WORKTREE_DIR="$SCRATCH_DIR/mm-base-worktree"
rm -rf "$WORKTREE_DIR"
git worktree add --detach "$WORKTREE_DIR" "$BASE_REF"



run_version_checks() {
  local VERSION="$1"
  local VERSION_UPPER=$(echo "$VERSION" | tr '[:lower:]' '[:upper:]')
  local TARGET_DIFF_PROC="$SCRATCH_DIR/diff-proc-${VERSION}"
  local PROVIDER_REPO
  local PROVIDER_CACHE
  local REAL_PACKAGE_NAME
  local REAL_FOLDER_NAME

  if [ "$VERSION" = "ga" ]; then
    PROVIDER_REPO="https://github.com/hashicorp/terraform-provider-google.git"
    PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-ga"
    REAL_PACKAGE_NAME="github.com/hashicorp/terraform-provider-google"
    REAL_FOLDER_NAME="google"
  else
    PROVIDER_REPO="https://github.com/hashicorp/terraform-provider-google-beta.git"
    if [ -d "${REPO_ROOT}/scratch/provider-cache" ] && [ ! -d "${REPO_ROOT}/scratch/provider-cache-beta" ]; then
      PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache"
    else
      PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache-beta"
    fi
    REAL_PACKAGE_NAME="github.com/hashicorp/terraform-provider-google-beta"
    REAL_FOLDER_NAME="google-beta"
  fi

  # Prepare isolated diff-processor directory
  rm -rf "$TARGET_DIFF_PROC"
  mkdir -p "$TARGET_DIFF_PROC"
  (cd "$DIFF_PROCESSOR_DIR" && tar cf - --exclude=old --exclude=new --exclude=bin .) | (cd "$TARGET_DIFF_PROC" && tar xf -)

  # Ensure we have a cached downstream provider clone
  if [ ! -d "$PROVIDER_CACHE" ]; then
    echo -e "${BLUE}[${VERSION_UPPER}] Cloning downstream provider into cache...${NC}"
    git clone --depth 500 "$PROVIDER_REPO" "$PROVIDER_CACHE" >/dev/null 2>&1
  else
    echo -e "${BLUE}[${VERSION_UPPER}] Updating cached downstream provider...${NC}"
    (
      cd "$PROVIDER_CACHE"
      git fetch --depth 500 origin main >/dev/null 2>&1
      git reset --hard origin/main >/dev/null 2>&1
    )
  fi

  # Find downstream commit matching base ref
  echo -e "${BLUE}[${VERSION_UPPER}] Finding downstream commit matching base ref ${YELLOW}${BASE_REF}${BLUE}...${NC}"
  local MATCHING_COMMIT=""
  for MM_SHA in $(git log --format="%H" -n 100 "$BASE_REF"); do
    MATCHING_COMMIT=$(cd "$PROVIDER_CACHE" && git log --grep="\[upstream:${MM_SHA}\]" --format="%H" -n 1)
    if [ -n "$MATCHING_COMMIT" ]; then
      echo -e "${GREEN}[${VERSION_UPPER}] Found matching provider commit: ${YELLOW}${MATCHING_COMMIT}${GREEN} (for MM commit ${MM_SHA:0:8})${NC}"
      break
    fi
  done

  if [ -z "$MATCHING_COMMIT" ]; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Warning: No matching provider commit found in recent ancestors of ${BASE_REF}; using provider HEAD.${NC}"
    MATCHING_COMMIT=$(cd "$PROVIDER_CACHE" && git rev-parse HEAD)
  fi

  # Prepare old and new directories in diff-processor
  echo -e "${BLUE}[${VERSION_UPPER}] Preparing diff-processor environment...${NC}"
  cp -R "$PROVIDER_CACHE" "$TARGET_DIFF_PROC/old"
  cp -R "$PROVIDER_CACHE" "$TARGET_DIFF_PROC/new"
  (cd "$TARGET_DIFF_PROC/old" && git checkout "$MATCHING_COMMIT" >/dev/null 2>&1)
  (cd "$TARGET_DIFF_PROC/new" && git checkout "$MATCHING_COMMIT" >/dev/null 2>&1)

  # Generate base and current provider code in parallel
  echo -e "${BLUE}[${VERSION_UPPER}] Generating base and current provider code...${NC}"
  (
    cd mmv1
    $MM_BINARY --output "$TARGET_DIFF_PROC/old" --version "$VERSION" --base "${WORKTREE_DIR}/mmv1"
  ) &
  local OLD_GEN_PID=$!

  (
    cd mmv1
    $MM_BINARY --output "$TARGET_DIFF_PROC/new" --version "$VERSION"
  ) &
  local NEW_GEN_PID=$!

  wait $OLD_GEN_PID
  wait $NEW_GEN_PID

  # Perform package substitutions for side-by-side compilation
  echo -e "${BLUE}[${VERSION_UPPER}] Preparing provider code for compilation...${NC}"
  (
    cd "$TARGET_DIFF_PROC"
    
    # Old package substitution
    (
      cd old/
      if [ -d "$REAL_FOLDER_NAME" ] && [ "$REAL_FOLDER_NAME" != "google" ]; then
        mv "$REAL_FOLDER_NAME" google
      fi
      local fake_package_name=google/provider/old

      if [ "$(uname)" = "Darwin" ]; then
        find . -type f -name "*.go" -exec sed -i "" "s~${REAL_PACKAGE_NAME}/${REAL_FOLDER_NAME}~${fake_package_name}/google~g" {} +
        sed -i "" "s|${REAL_PACKAGE_NAME}|${fake_package_name}|g" go.mod
      else
        find . -type f -name "*.go" -exec sed -i "s~${REAL_PACKAGE_NAME}/${REAL_FOLDER_NAME}~${fake_package_name}/google~g" {} +
        sed -i "s|${REAL_PACKAGE_NAME}|${fake_package_name}|g" go.mod
      fi
    ) &
    local OLD_SUB_PID=$!
    
    # New package substitution
    (
      cd new/
      if [ -d "$REAL_FOLDER_NAME" ] && [ "$REAL_FOLDER_NAME" != "google" ]; then
        mv "$REAL_FOLDER_NAME" google
      fi
      local fake_package_name=google/provider/new

      if [ "$(uname)" = "Darwin" ]; then
        find . -type f -name "*.go" -exec sed -i "" "s~${REAL_PACKAGE_NAME}/${REAL_FOLDER_NAME}~${fake_package_name}/google~g" {} +
        sed -i "" "s|${REAL_PACKAGE_NAME}|${fake_package_name}|g" go.mod
      else
        find . -type f -name "*.go" -exec sed -i "s~${REAL_PACKAGE_NAME}/${REAL_FOLDER_NAME}~${fake_package_name}/google~g" {} +
        sed -i "s|${REAL_PACKAGE_NAME}|${fake_package_name}|g" go.mod
      fi
    ) &
    local NEW_SUB_PID=$!

    wait $OLD_SUB_PID
    wait $NEW_SUB_PID
    
    # Tidy and build diff-processor
    echo -e "${BLUE}[${VERSION_UPPER}] Compiling diff-processor tool...${NC}"
    go mod tidy >/dev/null 2>&1
    mkdir -p bin/
    go build -o ./bin/diff-processor .
  )

  local HAS_ERR=0

  # Run breaking changes check
  echo -e "\n${BLUE}=== Running Breaking Changes Check (${VERSION_UPPER}) ===${NC}"
  set +e
  local BREAKING_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" breaking-changes 2>&1)
  local BREAKING_EXIT_CODE=$?
  set -e

  if [ $BREAKING_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor breaking-changes failed:${NC}"
    echo "$BREAKING_OUTPUT"
    return $BREAKING_EXIT_CODE
  fi

  if [ -n "$BREAKING_OUTPUT" ] && [ "$BREAKING_OUTPUT" != "[]" ] && [ "$BREAKING_OUTPUT" != "null" ]; then
    echo -e "${RED}[${VERSION_UPPER}] Breaking changes detected!${NC}"
    if command -v jq &> /dev/null; then
      echo "$BREAKING_OUTPUT" | jq .
    else
      echo "$BREAKING_OUTPUT" | python3 -m json.tool
    fi
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No breaking changes detected!${NC}"
  fi

  # Run missing tests check
  echo -e "\n${BLUE}=== Running Missing Tests Check (${VERSION_UPPER}) ===${NC}"
  set +e
  local TESTS_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" detect-missing-tests "$TARGET_DIFF_PROC/new/google/services" 2>&1)
  local TESTS_EXIT_CODE=$?
  set -e

  if [ $TESTS_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor detect-missing-tests failed:${NC}"
    echo "$TESTS_OUTPUT"
    return $TESTS_EXIT_CODE
  fi

  if [ -n "$TESTS_OUTPUT" ] && [ "$TESTS_OUTPUT" != "{}" ] && [ "$TESTS_OUTPUT" != "null" ]; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Missing tests detected!${NC}"
    if command -v jq &> /dev/null; then
      echo "$TESTS_OUTPUT" | jq .
    else
      echo "$TESTS_OUTPUT" | python3 -m json.tool
    fi
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No missing tests detected!${NC}"
  fi

  # Run missing documentation check
  echo -e "\n${BLUE}=== Running Missing Documentation Check (${VERSION_UPPER}) ===${NC}"
  set +e
  local DOCS_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" detect-missing-docs "$TARGET_DIFF_PROC/new" 2>&1)
  local DOCS_EXIT_CODE=$?
  set -e

  if [ $DOCS_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor detect-missing-docs failed:${NC}"
    echo "$DOCS_OUTPUT"
    return $DOCS_EXIT_CODE
  fi

  if has_missing_docs "$DOCS_OUTPUT"; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Missing documentation detected!${NC}"
    if command -v jq &> /dev/null; then
      echo "$DOCS_OUTPUT" | jq .
    else
      echo "$DOCS_OUTPUT" | python3 -m json.tool
    fi
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No missing documentation detected!${NC}"
  fi

  return $HAS_ERR
}

echo -e "${BLUE}Running GA and Beta schema diff checks in parallel...${NC}"

run_version_checks ga > "$SCRATCH_DIR/check-ga.log" 2>&1 &
GA_PID=$!

run_version_checks beta > "$SCRATCH_DIR/check-beta.log" 2>&1 &
BETA_PID=$!

set +e
wait $GA_PID
GA_STATUS=$?

wait $BETA_PID
BETA_STATUS=$?
set -e

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== GA Provider Schema Diff Results ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "$SCRATCH_DIR/check-ga.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Beta Provider Schema Diff Results ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "$SCRATCH_DIR/check-beta.log"

# 7. Final status report
echo -e "\n${BLUE}===========================================${NC}"
if [ $GA_STATUS -ne 0 ] || [ $BETA_STATUS -ne 0 ]; then
  echo -e "${RED}=== Schema change checks failed! ===${NC}"
  exit 1
else
  echo -e "${GREEN}=== All schema change checks passed successfully! ===${NC}"
  exit 0
fi
