#!/usr/bin/env bash
#
# run_pre_gen_checks.sh
#
# Runs Phase 1 pre-generation and static checks against magic-modules.
# Checks:
# 1. gofmt formatting check
# 2. tools/template-check (version-guard and unused-tmpl)
# 3. mmv1 unit tests
# 4. internal tools unit tests (go-changelog, issue-labeler, template-check, test-reader)
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

# 2. Go Formatting Check (gofmt)
echo -e "${BLUE}Running gofmt check...${NC}"
GOFMT_OUTPUT="$(gofmt -l $(git ls-files '*.go'))"
if [ -n "$GOFMT_OUTPUT" ]; then
  echo -e "${RED}The following files are not formatted properly:${NC}" >&2
  echo "$GOFMT_OUTPUT" >&2
  exit 1
fi

# 3. Template Validation Checks (tools/template-check)
echo -e "${BLUE}Running template validation checks...${NC}"
(
  cd tools/template-check
  tmpls=$(git diff --name-only --diff-filter=d origin/main ../../*.tmpl | sed 's=^=../../=g')
  if [ -n "$tmpls" ]; then
    go run main.go version-guard --file-list "${tmpls//$'\n'/,}"
  fi
)

(
  cd tools/template-check
  newtmplfiles=$(git diff --name-only --diff-filter=A origin/main HEAD -- ../../mmv1 | grep .tmpl | sed 's=^=../../=g')
  if [ -n "$newtmplfiles" ]; then
    go run main.go unused-tmpl --file-list "${newtmplfiles//$'\n'/,}"
  fi
)

# 4. MMv1 Core Unit Tests
echo -e "${BLUE}Running mmv1 unit tests...${NC}"
(cd mmv1 && go test ./...)

# 5. Tool Unit Tests
echo -e "${BLUE}Running tool unit tests...${NC}"
(cd tools/go-changelog && go test ./...)
(cd tools/issue-labeler && go test ./...)
(cd tools/template-check && go test ./...)
(cd tools/test-reader && go test ./...)

echo -e "${GREEN}=== All pre-generation & static checks passed successfully! ===${NC}"
