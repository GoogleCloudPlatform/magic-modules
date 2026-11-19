#!/usr/bin/env bash
#
# validate_provider_changes.sh
#
# Automates local provider change validation checks for Magic Modules by
# comparing the current working tree against a base branch (e.g., main),
# running all 5 diff-processor checks from CI (.ci/magician/cmd/generate_comment.go)
# on both GA and Beta (also catching GA-only gaps such as Beta -> GA promotions):
#   - breaking-changes        (GA and Beta)
#   - detect-missing-tests    (GA and Beta; supports --exhaustive)
#   - detect-unreadable-tests (GA and Beta)
#   - detect-missing-docs     (GA and Beta)
#   - detect-missing-identity (GA and Beta)
#
# Usage:
#   ./.agents/skills/utils/validate-provider-changes/scripts/validate_provider_changes.sh [--exhaustive] [base_ref]
#

set -e

# Define color codes for pretty printing
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
  cat <<'EOF'
Usage: validate_provider_changes.sh [--exhaustive] [base_ref]

Validate generated provider changes against a base branch or commit, running
all 5 diff-processor checks from CI (.ci/magician/cmd/generate_comment.go)
on both GA and Beta (also catching GA-only gaps such as Beta -> GA promotions):
  - breaking-changes        (GA and Beta)
  - detect-missing-tests    (GA and Beta)
  - detect-unreadable-tests (GA and Beta)
  - detect-missing-docs     (GA and Beta)
  - detect-missing-identity (GA and Beta)

Options:
  --exhaustive    Pass --exhaustive to detect-missing-tests to verify SET,
                  MODIFY, and UNSET transitions at every nesting level.
  -h, --help      Show this help message and exit.

Arguments:
  base_ref        Git ref to compare against (default: merge base with
                  origin/main, upstream/main, or main).
EOF
}

# 1. Parse CLI arguments
BASE_REF=""
EXHAUSTIVE=false

while [ $# -gt 0 ]; do
  case "$1" in
    --exhaustive)
      EXHAUSTIVE=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      while [ $# -gt 0 ]; do
        if [ -n "$BASE_REF" ]; then
          echo -e "${RED}Error: Unexpected extra argument: $1${NC}" >&2
          usage >&2
          exit 1
        fi
        BASE_REF="$1"
        shift
      done
      ;;
    -*)
      echo -e "${RED}Error: Unknown option: $1${NC}" >&2
      usage >&2
      exit 1
      ;;
    *)
      if [ -n "$BASE_REF" ]; then
        echo -e "${RED}Error: Unexpected extra argument: $1${NC}" >&2
        usage >&2
        exit 1
      fi
      BASE_REF="$1"
      shift
      ;;
  esac
done

echo -e "${BLUE}=== Magic Modules Provider Change Validator ===${NC}"

# 2. Validate environment
if [ ! -d "mmv1" ] || [ ! -d "tools/diff-processor" ]; then
  echo -e "${RED}Error: This script must be run from the root of the magic-modules repository.${NC}"
  exit 1
fi

# 3. Determine base ref
if [ -z "$BASE_REF" ]; then
  # Try to find the merge base with origin/main, fallback to upstream/main, fallback to main
  BASE_REF=$(git merge-base HEAD origin/main 2>/dev/null || git merge-base HEAD upstream/main 2>/dev/null || git merge-base HEAD main 2>/dev/null || echo "main")
fi

echo -e "${BLUE}Comparing current changes against base ref: ${YELLOW}${BASE_REF}${NC}"
if [ "$EXHAUSTIVE" = true ]; then
  echo -e "${BLUE}Missing tests mode: ${YELLOW}--exhaustive${NC}"
fi

# 4. Set up paths
REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/schema-diff-check"
DIFF_PROCESSOR_DIR="${REPO_ROOT}/tools/diff-processor"

mkdir -p "$SCRATCH_DIR"

# Always copy test-reader and issue-labeler into SCRATCH_DIR so go.mod replace directives resolve
rm -rf "${SCRATCH_DIR}/test-reader" "${SCRATCH_DIR}/issue-labeler"
cp -r "${REPO_ROOT}/tools/test-reader" "$SCRATCH_DIR/"
cp -r "${REPO_ROOT}/tools/issue-labeler" "$SCRATCH_DIR/"

# 5. Set up cleanup trap to ensure we never leave git worktrees behind
cleanup() {
  echo -e "${BLUE}Cleaning up temporary files and worktrees...${NC}"
  if [ -d "$SCRATCH_DIR/mm-base-worktree" ]; then
    git worktree remove --force "$SCRATCH_DIR/mm-base-worktree" 2>/dev/null || true
  fi
}
trap cleanup EXIT

# Helper function to pretty-print JSON output
pretty_print_json() {
  local json="$1"
  if command -v jq &> /dev/null; then
    echo "$json" | jq . 2>/dev/null || echo "$json"
  else
    echo "$json" | python3 -m json.tool 2>/dev/null || echo "$json"
  fi
}

# Helper function to check if the docs JSON output has missing docs
has_missing_docs() {
  local json="$1"
  if command -v jq &> /dev/null; then
    # Returns 0 (true) if there are missing docs, 1 (false) otherwise
    echo "$json" | jq -e '((.Resource // []) | length > 0) or ((.DataSource // []) | length > 0)' >/dev/null 2>&1
    return $?
  else
    # Simple Python fallback to check if either list is non-empty
    python3 -c "
import sys, json
data = json.loads(sys.argv[1]) or {}
resources = data.get('Resource') or []
datasources = data.get('DataSource') or []
if len(resources) > 0 or len(datasources) > 0:
    sys.exit(0)
sys.exit(1)
" "$json" 2>/dev/null
    return $?
  fi
}

# Helper function to check if the unreadable tests JSON output has new or stale entries
has_unreadable_tests() {
  local json="$1"
  if command -v jq &> /dev/null; then
    # Returns 0 (true) if there are new or stale unreadable tests, 1 (false) otherwise
    echo "$json" | jq -e '((.NewUnreadableTests // {}) | length > 0) or ((.StaleUnreadableTests // []) | length > 0)' >/dev/null 2>&1
    return $?
  else
    # Simple Python fallback to check if either collection is non-empty
    python3 -c "
import sys, json
data = json.loads(sys.argv[1]) or {}
new_unreadable = data.get('NewUnreadableTests') or {}
stale_unreadable = data.get('StaleUnreadableTests') or []
if len(new_unreadable) > 0 or len(stale_unreadable) > 0:
    sys.exit(0)
sys.exit(1)
" "$json" 2>/dev/null
    return $?
  fi
}

# 6. Build current mmv1 binary in the current context (once for all checks)
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

# 7. Create temporary worktree for base commit (once for all checks)
echo -e "${BLUE}Creating temporary worktree for base commit ${YELLOW}${BASE_REF}${BLUE}...${NC}"
WORKTREE_DIR="$SCRATCH_DIR/mm-base-worktree"
rm -rf "$WORKTREE_DIR"
git worktree add --detach "$WORKTREE_DIR" "$BASE_REF"



run_version_checks() {
  local VERSION="$1"
  local VERSION_UPPER
  VERSION_UPPER=$(echo "$VERSION" | tr '[:lower:]' '[:upper:]')
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
    "$MM_BINARY" --output "$TARGET_DIFF_PROC/old" --version "$VERSION" --base "${WORKTREE_DIR}/mmv1"
  ) &
  local OLD_GEN_PID=$!

  (
    cd mmv1
    "$MM_BINARY" --output "$TARGET_DIFF_PROC/new" --version "$VERSION"
  ) &
  local NEW_GEN_PID=$!

  wait $OLD_GEN_PID
  wait $NEW_GEN_PID

  # Compute changed services before package substitution rewrites imports.
  # Mirrors changedServices() in .ci/magician/cmd/generate_comment.go (run for both GA and Beta).
  local CHANGED_SERVICES=()
  local OLD_SERVICES_DIR="$TARGET_DIFF_PROC/old/$REAL_FOLDER_NAME/services"
  local NEW_SERVICES_DIR="$TARGET_DIFF_PROC/new/$REAL_FOLDER_NAME/services"
  local svc
  while IFS= read -r svc; do
    [ -z "$svc" ] && continue
    if [ ! -d "$OLD_SERVICES_DIR/$svc" ] || [ ! -d "$NEW_SERVICES_DIR/$svc" ] || ! diff -rq "$OLD_SERVICES_DIR/$svc" "$NEW_SERVICES_DIR/$svc" >/dev/null 2>&1; then
      CHANGED_SERVICES+=("$svc")
    fi
  done < <(
    for d in "$OLD_SERVICES_DIR"/*/ "$NEW_SERVICES_DIR"/*/; do
      [ -d "$d" ] || continue
      d="${d%/}"
      printf '%s\n' "${d##*/}"
    done | LC_ALL=C sort -u
  )

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
  local STDERR_FILE="$TARGET_DIFF_PROC/stderr.log"

  # 1. Run breaking changes check
  echo -e "\n${BLUE}=== Running Breaking Changes Check (${VERSION_UPPER}) ===${NC}"
  set +e
  local BREAKING_OUTPUT
  BREAKING_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" breaking-changes 2>"$STDERR_FILE")
  local BREAKING_EXIT_CODE=$?
  set -e

  if [ $BREAKING_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor breaking-changes failed:${NC}"
    [ -n "$BREAKING_OUTPUT" ] && echo "$BREAKING_OUTPUT"
    [ -s "$STDERR_FILE" ] && cat "$STDERR_FILE"
    HAS_ERR=1
  elif [ -n "$BREAKING_OUTPUT" ] && [ "$BREAKING_OUTPUT" != "[]" ] && [ "$BREAKING_OUTPUT" != "null" ]; then
    echo -e "${RED}[${VERSION_UPPER}] Breaking changes detected!${NC}"
    pretty_print_json "$BREAKING_OUTPUT"
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No breaking changes detected!${NC}"
  fi

  # 2. Run missing tests check
  local TEST_ARGS=(detect-missing-tests)
  local TESTS_LABEL="Missing Tests Check (${VERSION_UPPER})"
  if [ "$EXHAUSTIVE" = true ]; then
    TEST_ARGS+=(--exhaustive)
    TESTS_LABEL="Missing Tests Check (${VERSION_UPPER}, --exhaustive)"
  fi
  TEST_ARGS+=("$TARGET_DIFF_PROC/new/google/services")

  echo -e "\n${BLUE}=== Running ${TESTS_LABEL} ===${NC}"
  set +e
  local TESTS_OUTPUT
  TESTS_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" "${TEST_ARGS[@]}" 2>"$STDERR_FILE")
  local TESTS_EXIT_CODE=$?
  set -e

  if [ $TESTS_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor detect-missing-tests failed:${NC}"
    [ -n "$TESTS_OUTPUT" ] && echo "$TESTS_OUTPUT"
    [ -s "$STDERR_FILE" ] && cat "$STDERR_FILE"
    HAS_ERR=1
  elif [ -n "$TESTS_OUTPUT" ] && [ "$TESTS_OUTPUT" != "{}" ] && [ "$TESTS_OUTPUT" != "null" ]; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Missing tests detected!${NC}"
    pretty_print_json "$TESTS_OUTPUT"
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No missing tests detected!${NC}"
  fi

  # 3. Run unreadable tests check
  echo -e "\n${BLUE}=== Running Unreadable Tests Check (${VERSION_UPPER}) ===${NC}"
  if [ ${#CHANGED_SERVICES[@]} -gt 0 ]; then
    echo -e "${BLUE}[${VERSION_UPPER}] Checking changed services: ${YELLOW}${CHANGED_SERVICES[*]}${NC}"
  else
    echo -e "${BLUE}[${VERSION_UPPER}] No changed services detected.${NC}"
  fi
  set +e
  local UNREADABLE_OUTPUT
  UNREADABLE_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" detect-unreadable-tests "$TARGET_DIFF_PROC/new/google/services" ${CHANGED_SERVICES[@]+"${CHANGED_SERVICES[@]}"} 2>"$STDERR_FILE")
  local UNREADABLE_EXIT_CODE=$?
  set -e

  if [ $UNREADABLE_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor detect-unreadable-tests failed:${NC}"
    [ -n "$UNREADABLE_OUTPUT" ] && echo "$UNREADABLE_OUTPUT"
    [ -s "$STDERR_FILE" ] && cat "$STDERR_FILE"
    HAS_ERR=1
  elif has_unreadable_tests "$UNREADABLE_OUTPUT"; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Unreadable or stale unreadable tests detected!${NC}"
    pretty_print_json "$UNREADABLE_OUTPUT"
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No unreadable tests detected!${NC}"
  fi

  # 4. Run missing documentation check
  echo -e "\n${BLUE}=== Running Missing Documentation Check (${VERSION_UPPER}) ===${NC}"
  set +e
  local DOCS_OUTPUT
  DOCS_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" detect-missing-docs "$TARGET_DIFF_PROC/new" 2>"$STDERR_FILE")
  local DOCS_EXIT_CODE=$?
  set -e

  if [ $DOCS_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor detect-missing-docs failed:${NC}"
    [ -n "$DOCS_OUTPUT" ] && echo "$DOCS_OUTPUT"
    [ -s "$STDERR_FILE" ] && cat "$STDERR_FILE"
    HAS_ERR=1
  elif has_missing_docs "$DOCS_OUTPUT"; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Missing documentation detected!${NC}"
    pretty_print_json "$DOCS_OUTPUT"
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No missing documentation detected!${NC}"
  fi

  # 5. Run missing ResourceIdentity coverage check
  echo -e "\n${BLUE}=== Running Missing Identity Check (${VERSION_UPPER}) ===${NC}"
  set +e
  local IDENTITY_OUTPUT
  IDENTITY_OUTPUT=$("$TARGET_DIFF_PROC/bin/diff-processor" detect-missing-identity "$TARGET_DIFF_PROC/new/google/services" 2>"$STDERR_FILE")
  local IDENTITY_EXIT_CODE=$?
  set -e

  if [ $IDENTITY_EXIT_CODE -ne 0 ]; then
    echo -e "${RED}[${VERSION_UPPER}] Error: diff-processor detect-missing-identity failed:${NC}"
    [ -n "$IDENTITY_OUTPUT" ] && echo "$IDENTITY_OUTPUT"
    [ -s "$STDERR_FILE" ] && cat "$STDERR_FILE"
    HAS_ERR=1
  elif [ -n "$IDENTITY_OUTPUT" ] && [ "$IDENTITY_OUTPUT" != "{}" ] && [ "$IDENTITY_OUTPUT" != "null" ]; then
    echo -e "${YELLOW}[${VERSION_UPPER}] Missing ResourceIdentity coverage detected!${NC}"
    pretty_print_json "$IDENTITY_OUTPUT"
    HAS_ERR=1
  else
    echo -e "${GREEN}[${VERSION_UPPER}] No missing ResourceIdentity coverage detected!${NC}"
  fi

  return $HAS_ERR
}

echo -e "${BLUE}Running GA and Beta provider validation checks in parallel...${NC}"

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
echo -e "${BLUE}=== GA Provider Validation Results ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "$SCRATCH_DIR/check-ga.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Beta Provider Validation Results ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "$SCRATCH_DIR/check-beta.log"

# 8. Final status report
echo -e "\n${BLUE}===========================================${NC}"
if [ $GA_STATUS -ne 0 ] || [ $BETA_STATUS -ne 0 ]; then
  echo -e "${RED}=== Provider validation checks failed! ===${NC}"
  exit 1
else
  echo -e "${GREEN}=== All provider validation checks passed successfully! ===${NC}"
  exit 0
fi
