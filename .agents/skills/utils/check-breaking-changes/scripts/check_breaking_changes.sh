#!/usr/bin/env bash
#
# check_breaking_changes.sh
#
# Automates local breaking change detection for Magic Modules by comparing
# the current working tree against a base branch (e.g., main).
#
# Usage:
#   ./.agents/skills/utils/check-breaking-changes/scripts/check_breaking_changes.sh [base_ref]
#

set -e

# Define color codes for pretty printing
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Magic Modules Breaking Changes Checker ===${NC}"

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
  echo -e "${GREEN}No changes detected in mmv1/products/. Skipping breaking changes check.${NC}"
  exit 0
fi

echo -e "${BLUE}Detected changed products:${NC}"
for PRODUCT in $PRODUCTS; do
  echo -e "  - ${YELLOW}${PRODUCT}${NC}"
done

# 4. Set up paths
REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/breaking-changes-check"
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

# 9. Generate "Old" (Base) Provider Code
echo -e "${BLUE}Creating temporary worktree for base commit ${YELLOW}${BASE_REF}${BLUE}...${NC}"
WORKTREE_DIR="$SCRATCH_DIR/mm-base-worktree"
rm -rf "$WORKTREE_DIR"
git worktree add --detach "$WORKTREE_DIR" "$BASE_REF"

echo -e "${BLUE}Generating base provider code for products: ${YELLOW}${PRODUCT_LIST}${BLUE}...${NC}"
(
  cd mmv1
  $MM_BINARY --output "$DIFF_PROCESSOR_DIR/old" --version beta --product "$PRODUCT_LIST" --base "${WORKTREE_DIR}/mmv1"
)

# 10. Generate "New" (Current) Provider Code
echo -e "${BLUE}Generating current provider code for products: ${YELLOW}${PRODUCT_LIST}${BLUE}...${NC}"
(
  cd mmv1
  $MM_BINARY --output "$DIFF_PROCESSOR_DIR/new" --version beta --product "$PRODUCT_LIST"
)

# 10. Perform package substitutions for side-by-side compilation
echo -e "${BLUE}Preparing provider code for compilation...${NC}"
(
  cd "$DIFF_PROCESSOR_DIR"
  
  real_package_name=github.com/hashicorp/terraform-provider-google-beta
  real_folder_name=google-beta

  # Old package substitution
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
  
  # New package substitution
  cd ../new/
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
  
  # Tidy and build diff-processor
  cd ..
  echo -e "${BLUE}Compiling diff-processor tool...${NC}"
  go mod tidy
  mkdir -p bin/
  go build -o ./bin/diff-processor .
)

# 11. Run the breaking changes check
echo -e "${BLUE}Running breaking changes check...${NC}"
set +e
OUTPUT=$("$DIFF_PROCESSOR_DIR/bin/diff-processor" breaking-changes 2>&1)
EXIT_CODE=$?
set -e

if [ $EXIT_CODE -ne 0 ]; then
  echo -e "${RED}Error: diff-processor failed to run or encountered a compilation issue:${NC}"
  echo "$OUTPUT"
  exit $EXIT_CODE
fi

if [ -z "$OUTPUT" ] || [ "$OUTPUT" == "[]" ] || [ "$OUTPUT" == "null" ]; then
  echo -e "${GREEN}No breaking changes detected! ${NC}"
else
  echo -e "${RED}Breaking changes detected!${NC}"
  if command -v jq &> /dev/null; then
    echo "$OUTPUT" | jq .
  else
    echo "$OUTPUT" | python3 -m json.tool
  fi
  exit 1
fi
