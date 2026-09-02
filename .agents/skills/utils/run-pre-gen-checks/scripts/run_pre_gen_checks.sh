#!/usr/bin/env bash
#
# run_pre_gen_checks.sh
#
# Runs Phase 1 pre-generation and static checks against magic-modules in parallel.
# Checks:
# 1. Submodule check (disallows git submodules)
# 2. gofmt formatting check
# 3. tools/template-check (version-guard and unused-tmpl)
# 4. YAML linting check on changed product YAML files (if yamllint is installed)
# 5. mmv1 unit tests
# 6. internal tools & magician unit tests (go-changelog, issue-labeler, template-check, test-reader, .ci/magician)
#

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Pre-Generation & Static Checks ===${NC}"

# 1. Validate environment
if [ ! -d "mmv1" ]; then
  echo -e "${RED}Error: This script must be run from the root of the magic-modules repository.${NC}"
  exit 1
fi

REPO_ROOT=$(pwd)
SCRATCH_DIR="${REPO_ROOT}/scratch/pre-gen-checks"
LOGS_DIR="${SCRATCH_DIR}/logs"
mkdir -p "$LOGS_DIR"

echo -e "${BLUE}Running pre-generation checks in parallel...${NC}"

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
  tmpls=$(git diff --name-only --diff-filter=d origin/main ../../*.tmpl | sed 's=^=../../=g')
  if [ -n "$tmpls" ]; then
    go run main.go version-guard --file-list "${tmpls//$'\n'/,}"
  fi

  newtmplfiles=$(git diff --name-only --diff-filter=A origin/main HEAD -- ../../mmv1 | grep .tmpl | sed 's=^=../../=g')
  if [ -n "$newtmplfiles" ]; then
    go run main.go unused-tmpl --file-list "${newtmplfiles//$'\n'/,}"
  fi
  echo "Template validation checks passed."
) > "${LOGS_DIR}/template.log" 2>&1 &
TEMPLATE_PID=$!

# Check 4: YAML Linting Check (yamllint)
(
  if command -v yamllint >/dev/null 2>&1; then
    yamlfiles=$(git diff --name-only --diff-filter=d origin/main -- mmv1/products 2>/dev/null || true)
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

# Check 6: Tool & CI Unit Tests
(
  (cd tools/go-changelog && go test ./...)
  (cd tools/issue-labeler && go test ./...)
  (cd tools/template-check && go test ./...)
  (cd tools/test-reader && go test ./...)
  (cd .ci/magician && go test ./...)
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

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Submodule Check Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/submodule.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Go Formatting Check Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/gofmt.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Template Validation Check Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/template.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== YAML Lint Check Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/yamllint.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== MMv1 Unit Test Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/mmv1_unit.log"

echo -e "\n${BLUE}=====================================================${NC}"
echo -e "${BLUE}=== Tool & Magician Unit Test Logs ===${NC}"
echo -e "${BLUE}=====================================================${NC}"
cat "${LOGS_DIR}/tools_unit.log"

if [ $SUBMODULE_STATUS -ne 0 ] || [ $GOFMT_STATUS -ne 0 ] || [ $TEMPLATE_STATUS -ne 0 ] || [ $YAMLLINT_STATUS -ne 0 ] || [ $MMV1_STATUS -ne 0 ] || [ $TOOLS_STATUS -ne 0 ]; then
  echo -e "\n${RED}=== Pre-generation & static checks failed! ===${NC}"
  exit 1
else
  echo -e "\n${GREEN}=== All pre-generation & static checks passed successfully! ===${NC}"
  exit 0
fi
