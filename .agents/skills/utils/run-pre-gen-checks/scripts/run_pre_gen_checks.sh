#!/usr/bin/env bash
#
# run_pre_gen_checks.sh
#
# Runs Phase 1 pre-generation and static checks against magic-modules in parallel.
# Checks:
# 1. Submodule check (disallows git submodules)
# 2. gofmt formatting check
# 3. tools/template-check (version-guard and unused-tmpl on modified .tmpl files)
# 4. YAML linting check on changed product YAML files (if yamllint is installed)
# 5. mmv1 unit tests
# 6. internal tools & magician unit tests (scoped to touched tools/ or .ci/ files, or all with --all)
#

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

RUN_ALL=0
for arg in "$@"; do
  if [ "$arg" == "--all" ]; then
    RUN_ALL=1
  fi
done

echo -e "${BLUE}=== Pre-Generation & Static Checks ===${NC}"

# 1. Validate environment & determine repo root reliably
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$REPO_ROOT" ] || [ ! -d "${REPO_ROOT}/mmv1" ]; then
  if [ -d "mmv1" ]; then
    REPO_ROOT="$(pwd)"
  else
    echo -e "${RED}Error: This script must be run within the magic-modules repository.${NC}"
    exit 1
  fi
fi

cd "$REPO_ROOT"

# Determine base ref for diff checks
BASE_REF="origin/main"
if ! git rev-parse --verify "$BASE_REF" >/dev/null 2>&1; then
  if git rev-parse --verify "upstream/main" >/dev/null 2>&1; then
    BASE_REF="upstream/main"
  else
    BASE_REF="HEAD"
  fi
fi

# Detect all modified files (staged, unstaged, or committed relative to base ref)
CHANGED_FILES="$(git diff --name-only "$BASE_REF" 2>/dev/null || git diff --name-only HEAD 2>/dev/null || true)"

SCRATCH_DIR="${REPO_ROOT}/scratch/pre-gen-checks"
LOGS_DIR="${SCRATCH_DIR}/logs"
mkdir -p "$LOGS_DIR"

echo -e "${BLUE}Running pre-generation checks in parallel (base: ${BASE_REF})...${NC}"

# Check 1: Submodule Check
(
  output=$(git submodule status --recursive 2>&1)
  if [ -n "$output" ]; then
    echo "$output" >&2
    echo "Submodules are not allowed" >&2
    exit 1
  else
    echo "No submodules found. Submodule check passed."
  fi
) > "${LOGS_DIR}/submodule.log" 2>&1 &
SUBMODULE_PID=$!

# Check 2: Go Formatting Check (gofmt)
(
  GOFMT_OUTPUT="$(gofmt -l $(git ls-files '*.go'))"
  if [ -n "$GOFMT_OUTPUT" ]; then
    echo "The following files are not formatted properly:" >&2
    echo "$GOFMT_OUTPUT" >&2
    exit 1
  fi
  echo "gofmt check passed."
) > "${LOGS_DIR}/gofmt.log" 2>&1 &
GOFMT_PID=$!

# Check 3: Template Validation Checks (tools/template-check)
(
  cd tools/template-check
  tmpls=$(git diff --name-only --diff-filter=d "$BASE_REF" ../../*.tmpl 2>/dev/null | sed 's=^=../../=g' || true)
  if [ -n "$tmpls" ]; then
    go run main.go version-guard --file-list "${tmpls//$'\n'/,}"
  fi

  newtmplfiles=$(git diff --name-only --diff-filter=A "$BASE_REF" HEAD -- ../../mmv1 2>/dev/null | grep '\.tmpl$' | sed 's=^=../../=g' || true)
  if [ -n "$newtmplfiles" ]; then
    go run main.go unused-tmpl --file-list "${newtmplfiles//$'\n'/,}"
  fi
  echo "Template validation checks passed."
) > "${LOGS_DIR}/template.log" 2>&1 &
TEMPLATE_PID=$!

# Check 4: YAML Linting Check (yamllint)
(
  if command -v yamllint >/dev/null 2>&1; then
    yamlfiles=$(echo "$CHANGED_FILES" | grep '^mmv1/products/' || true)
    if [ -n "$yamlfiles" ]; then
      yamllint -c .yamllint $yamlfiles
      echo "YAML lint check passed for modified product YAML files."
    else
      echo "No modified product YAML files to lint. YAML lint check passed."
    fi
  else
    echo "yamllint not found in PATH; skipping YAML lint check."
  fi
) > "${LOGS_DIR}/yamllint.log" 2>&1 &
YAMLLINT_PID=$!

# Check 5: MMv1 Core Unit Tests
(
  cd mmv1 && go test ./...
) > "${LOGS_DIR}/mmv1_unit.log" 2>&1 &
MMV1_PID=$!

# Check 6: Tool & CI Unit Tests (diff-scoped unless --all or RUN_ALL=1)
(
  ran_any=0
  if [ "$RUN_ALL" -eq 1 ] || echo "$CHANGED_FILES" | grep -q '^tools/go-changelog/'; then
    echo "Testing tools/go-changelog..."
    (cd tools/go-changelog && go test ./...)
    ran_any=1
  fi
  if [ "$RUN_ALL" -eq 1 ] || echo "$CHANGED_FILES" | grep -q '^tools/issue-labeler/'; then
    echo "Testing tools/issue-labeler..."
    (cd tools/issue-labeler && go test ./...)
    ran_any=1
  fi
  if [ "$RUN_ALL" -eq 1 ] || echo "$CHANGED_FILES" | grep -q '^tools/template-check/'; then
    echo "Testing tools/template-check..."
    (cd tools/template-check && go test ./...)
    ran_any=1
  fi
  if [ "$RUN_ALL" -eq 1 ] || echo "$CHANGED_FILES" | grep -q '^tools/test-reader/'; then
    echo "Testing tools/test-reader..."
    (cd tools/test-reader && go test ./...)
    ran_any=1
  fi
  if [ "$RUN_ALL" -eq 1 ] || echo "$CHANGED_FILES" | grep -q '^\.ci/'; then
    echo "Testing .ci/magician..."
    (cd .ci/magician && go test ./...)
    ran_any=1
  fi
  if [ "$RUN_ALL" -eq 1 ] || echo "$CHANGED_FILES" | grep -q '^\.agents/scripts/'; then
    echo "Testing .agents/scripts..."
    python3 -m unittest discover -s .agents/scripts -p "test_*.py"
    ran_any=1
  fi
  if [ "$ran_any" -eq 0 ]; then
    echo "No changes detected in tools/, .ci/, or .agents/scripts/; skipping tool unit tests (use --all to run)."
  fi
) > "${LOGS_DIR}/tools_unit.log" 2>&1 &
TOOLS_PID=$!

set +e
wait $SUBMODULE_PID
SUBMODULE_STATUS=$?

wait $GOFMT_PID
GOFMT_STATUS=$?

wait $TEMPLATE_PID
TEMPLATE_STATUS=$?

wait $YAMLLINT_PID
YAMLLINT_STATUS=$?

wait $MMV1_PID
MMV1_STATUS=$?

wait $TOOLS_PID
TOOLS_STATUS=$?
set -e

# Quiet on success: only output failure logs when something fails
if [ $SUBMODULE_STATUS -eq 0 ] && [ $GOFMT_STATUS -eq 0 ] && [ $TEMPLATE_STATUS -eq 0 ] && [ $YAMLLINT_STATUS -eq 0 ] && [ $MMV1_STATUS -eq 0 ] && [ $TOOLS_STATUS -eq 0 ]; then
  echo -e "\n${GREEN}=== All pre-generation & static checks passed successfully! ===${NC}"
  echo "  ✓ Submodule check passed"
  echo "  ✓ Go formatting check passed (gofmt)"
  echo "  ✓ Template validation checks passed"
  echo "  ✓ YAML lint check passed"
  echo "  ✓ MMv1 core unit tests passed"
  echo "  ✓ Tool & CI unit tests passed"
  exit 0
else
  echo -e "\n${RED}=== Pre-generation & static checks failed! ===${NC}"

  if [ $SUBMODULE_STATUS -ne 0 ]; then
    echo -e "\n${RED}=== Submodule Check Logs ===${NC}"
    cat "${LOGS_DIR}/submodule.log"
  fi

  if [ $GOFMT_STATUS -ne 0 ]; then
    echo -e "\n${RED}=== Go Formatting Check Logs (gofmt) ===${NC}"
    cat "${LOGS_DIR}/gofmt.log"
  fi

  if [ $TEMPLATE_STATUS -ne 0 ]; then
    echo -e "\n${RED}=== Template Validation Check Logs ===${NC}"
    cat "${LOGS_DIR}/template.log"
  fi

  if [ $YAMLLINT_STATUS -ne 0 ]; then
    echo -e "\n${RED}=== YAML Lint Check Logs ===${NC}"
    cat "${LOGS_DIR}/yamllint.log"
  fi

  if [ $MMV1_STATUS -ne 0 ]; then
    echo -e "\n${RED}=== MMv1 Unit Test Logs ===${NC}"
    cat "${LOGS_DIR}/mmv1_unit.log"
  fi

  if [ $TOOLS_STATUS -ne 0 ]; then
    echo -e "\n${RED}=== Tool & Magician Unit Test Logs ===${NC}"
    cat "${LOGS_DIR}/tools_unit.log"
  fi

  exit 1
fi
