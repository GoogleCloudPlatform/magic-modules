#! /bin/bash

pushd magic-modules-branched

BRANCH_NAME=$(cat branchname)
{
    echo "## 3.0.0 diff report as of $(git rev-parse HEAD^2)";
    echo "[TPG Diff](https://github.com/modular-magician/terraform-provider-google/compare/$BRANCH_NAME-previous..$BRANCH_NAME)";
    echo "[TPGB Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/$BRANCH_NAME-previous..$BRANCH_NAME)";
    echo "[Mapper Diff](https://github.com/modular-magician/terraform-google-conversion/compare/$BRANCH_NAME-previous..$BRANCH_NAME)";
} > ../message/message.txt
