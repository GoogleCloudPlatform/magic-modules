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

# 3. Detect changed products
PRODUCTS=$(git diff --name-only "$BASE_REF" | grep '^mmv1/products/' | cut -d'/' -f3 | sort -u || true)

if [ -z "$PRODUCTS" ]; then
  echo -e "${GREEN}No changes detected in mmv1/products/. Skipping schema diff checks.${NC}"
  exit 0
fi

echo -e "${BLUE}Detected changed products:${NC}"
for PRODUCT in $PRODUCTS; do
  echo -e "  - ${YELLOW}${PRODUCT}${NC}"
done

# 4. Set up paths
REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/schema-diff-check"
PROVIDER_CACHE="${REPO_ROOT}/scratch/provider-cache"
DIFF_PROCESSOR_DIR="${REPO_ROOT}/tools/diff-processor"

mkdir -p "$SCRATCH_DIR"

# 5. Set up cleanup trap to ensure we never leave git worktrees behind
cleanup() {
  echo -e "${BLUE}Cleaning up temporary files and worktrees...${NC}"
  if [ -d "$SCRATCH_DIR/mm-base-worktree" ]; then
    git worktree remove --force "$SCRATCH_DIR/mm-base-worktree" 2>/dev/null || true
  fi
  # Clean up backup files created by sed
  find "$DIFF_PROCESSOR_DIR/old" -name "*.bak" -delete 2>/dev/null || true
  find "$DIFF_PROCESSOR_DIR/new" -name "*.bak" -delete 2>/dev/null || true
  # Clean up changes to go.mod and go.sum in diff-processor
  git checkout -- "$DIFF_PROCESSOR_DIR/go.mod" "$DIFF_PROCESSOR_DIR/go.sum" 2>/dev/null || true
}
trap cleanup EXIT

# 6. Ensure we have a cached downstream provider clone
if [ ! -d "$PROVIDER_CACHE" ]; then
  echo -e "${BLUE}Cloning downstream provider (depth 1) into cache...${NC}"
  git clone --depth 1 https://github.com/hashicorp/terraform-provider-google-beta.git "$PROVIDER_CACHE"
else
  echo -e "${BLUE}Updating cached downstream provider...${NC}"
  (
    cd "$PROVIDER_CACHE"
    git fetch --depth 1 origin main
    git reset --hard origin/main
  )
fi

# 7. Prepare old and new directories in diff-processor
echo -e "${BLUE}Preparing diff-processor environment...${NC}"
rm -rf "$DIFF_PROCESSOR_DIR/old" "$DIFF_PROCESSOR_DIR/new"
cp -R "$PROVIDER_CACHE" "$DIFF_PROCESSOR_DIR/old"
cp -R "$PROVIDER_CACHE" "$DIFF_PROCESSOR_DIR/new"

# Construct a comma-separated list of changed products
PRODUCT_LIST=$(echo "$PRODUCTS" | tr '\n' ',' | sed 's/,$//')

# Clean only the directories of the changed products in old and new directories
echo -e "${BLUE}Cleaning target provider directories for changed products: ${YELLOW}${PRODUCT_LIST}${BLUE}...${NC}"
for PRODUCT in $PRODUCTS; do
  rm -rf "$DIFF_PROCESSOR_DIR/old/google-beta/services/${PRODUCT}"
  rm -rf "$DIFF_PROCESSOR_DIR/new/google-beta/services/${PRODUCT}"
done

# 8. Build current mmv1 binary in the current context
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

# 9. Generate "Old" (Base) and "New" (Current) Provider Code in Parallel
echo -e "${BLUE}Creating temporary worktree for base commit ${YELLOW}${BASE_REF}${BLUE}...${NC}"
WORKTREE_DIR="$SCRATCH_DIR/mm-base-worktree"
rm -rf "$WORKTREE_DIR"
git worktree add --detach "$WORKTREE_DIR" "$BASE_REF"

echo -e "${BLUE}Generating base and current provider code in parallel for products: ${YELLOW}${PRODUCT_LIST}${BLUE}...${NC}"
(
  cd mmv1
  $MM_BINARY --output "$DIFF_PROCESSOR_DIR/old" --version beta --product "$PRODUCT_LIST" --base "${WORKTREE_DIR}/mmv1"
) &
OLD_GEN_PID=$!

(
  cd mmv1
  $MM_BINARY --output "$DIFF_PROCESSOR_DIR/new" --version beta --product "$PRODUCT_LIST"
) &
NEW_GEN_PID=$!

wait $OLD_GEN_PID
wait $NEW_GEN_PID

# 10. Perform package substitutions for side-by-side compilation in parallel
echo -e "${BLUE}Preparing provider code for compilation...${NC}"
(
  cd "$DIFF_PROCESSOR_DIR"
  
  real_package_name=github.com/hashicorp/terraform-provider-google-beta
  real_folder_name=google-beta
  
  # Old package substitution
  (
    cd old/
    if [ -d "google-beta" ]; then
      mv google-beta google
    fi
    fake_package_name=google/provider/old

    if [ "$(uname)" = "Darwin" ]; then
      find . -type f -name "*.go" -exec sed -i "" "s~${real_package_name}/${real_folder_name}~${fake_package_name}/google~g" {} +
      sed -i "" "s|${real_package_name}|${fake_package_name}|g" go.mod
    else
      find . -type f -name "*.go" -exec sed -i "s~${real_package_name}/${real_folder_name}~${fake_package_name}/google~g" {} +
      sed -i "s|${real_package_name}|${fake_package_name}|g" go.mod
    fi
  ) &
  OLD_SUB_PID=$!
  
  # New package substitution
  (
    cd new/
    if [ -d "google-beta" ]; then
      mv google-beta google
    fi
    fake_package_name=google/provider/new

    if [ "$(uname)" = "Darwin" ]; then
      find . -type f -name "*.go" -exec sed -i "" "s~${real_package_name}/${real_folder_name}~${fake_package_name}/google~g" {} +
      sed -i "" "s|${real_package_name}|${fake_package_name}|g" go.mod
    else
      find . -type f -name "*.go" -exec sed -i "s~${real_package_name}/${real_folder_name}~${fake_package_name}/google~g" {} +
      sed -i "s|${real_package_name}|${fake_package_name}|g" go.mod
    fi
  ) &
  NEW_SUB_PID=$!

  wait $OLD_SUB_PID
  wait $NEW_SUB_PID
  
  # Tidy and build diff-processor
  echo -e "${BLUE}Compiling diff-processor tool...${NC}"
  go mod tidy
  mkdir -p bin/
  go build -o ./bin/diff-processor .
)

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

HAS_ISSUES=0

# 11. Run the checks

# 11.1 Run breaking changes check
echo -e "\n${BLUE}=== Running Breaking Changes Check ===${NC}"
set +e
BREAKING_OUTPUT=$("$DIFF_PROCESSOR_DIR/bin/diff-processor" breaking-changes 2>&1)
BREAKING_EXIT_CODE=$?
set -e

if [ $BREAKING_EXIT_CODE -ne 0 ]; then
  echo -e "${RED}Error: diff-processor breaking-changes failed to run or encountered a compilation issue:${NC}"
  echo "$BREAKING_OUTPUT"
  exit $BREAKING_EXIT_CODE
fi

if [ -n "$BREAKING_OUTPUT" ] && [ "$BREAKING_OUTPUT" != "[]" ] && [ "$BREAKING_OUTPUT" != "null" ]; then
  echo -e "${RED}Breaking changes detected!${NC}"
  if command -v jq &> /dev/null; then
    echo "$BREAKING_OUTPUT" | jq .
  else
    echo "$BREAKING_OUTPUT" | python3 -m json.tool
  fi
  HAS_ISSUES=1
else
  echo -e "${GREEN}No breaking changes detected!${NC}"
fi

# 11.2 Run missing tests check
echo -e "\n${BLUE}=== Running Missing Tests Check ===${NC}"
set +e
TESTS_OUTPUT=$("$DIFF_PROCESSOR_DIR/bin/diff-processor" detect-missing-tests "$DIFF_PROCESSOR_DIR/new/google/services" 2>&1)
TESTS_EXIT_CODE=$?
set -e

if [ $TESTS_EXIT_CODE -ne 0 ]; then
  echo -e "${RED}Error: diff-processor detect-missing-tests failed to run:${NC}"
  echo "$TESTS_OUTPUT"
  exit $TESTS_EXIT_CODE
fi

if [ -n "$TESTS_OUTPUT" ] && [ "$TESTS_OUTPUT" != "{}" ] && [ "$TESTS_OUTPUT" != "null" ]; then
  echo -e "${YELLOW}Missing tests detected!${NC}"
  if command -v jq &> /dev/null; then
    echo "$TESTS_OUTPUT" | jq .
  else
    echo "$TESTS_OUTPUT" | python3 -m json.tool
  fi
  HAS_ISSUES=1
else
  echo -e "${GREEN}No missing tests detected!${NC}"
fi

# 11.3 Run missing documentation check
echo -e "\n${BLUE}=== Running Missing Documentation Check ===${NC}"
set +e
DOCS_OUTPUT=$("$DIFF_PROCESSOR_DIR/bin/diff-processor" detect-missing-docs "$DIFF_PROCESSOR_DIR/new" 2>&1)
DOCS_EXIT_CODE=$?
set -e

if [ $DOCS_EXIT_CODE -ne 0 ]; then
  echo -e "${RED}Error: diff-processor detect-missing-docs failed to run:${NC}"
  echo "$DOCS_OUTPUT"
  exit $DOCS_EXIT_CODE
fi

if has_missing_docs "$DOCS_OUTPUT"; then
  echo -e "${YELLOW}Missing documentation detected!${NC}"
  if command -v jq &> /dev/null; then
    echo "$DOCS_OUTPUT" | jq .
  else
    echo "$DOCS_OUTPUT" | python3 -m json.tool
  fi
  HAS_ISSUES=1
else
  echo -e "${GREEN}No missing documentation detected!${NC}"
fi

# 12. Final status report
echo -e "\n${BLUE}===========================================${NC}"
if [ $HAS_ISSUES -eq 1 ]; then
  echo -e "${RED}=== Schema change checks failed! ===${NC}"
  exit 1
else
  echo -e "${GREEN}=== All schema change checks passed successfully! ===${NC}"
  exit 0
fi
