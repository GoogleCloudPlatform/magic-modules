#!/bin/bash
# input: two environment variables
#  TPG_BREAKING - results of runing breaking change detector
#   against tpg
#  TPGB_BREAKING - results of runing breaking change detector
#   against tpgb
# output: echo to console
#  message section cotaining: a header,
#  tpg's unique messages, and all of tpgb's messages

tpgUnique=""
newline=$'\n'

# This while loop itterates over each individual
# line of TPG_BREAKING. The input to the while loop
# is through the <<< at the conclusion of the loop.
while read -r tpgi; do
  simpleTPG=$(sed 's/-.*//' <<< "$tpgi")
  found="false"
  while read -r tpgbi; do
    simpleTPGB=$(sed 's/-.*//' <<< "$tpgbi")
    if [ "$simpleTPG" == "$simpleTPGB" ]; then
      found="true"
   fi
  done <<< "$TPGB_BREAKING"
  if [ "$found" != "true" ]; then
    if [ "$tpgUnique" == "" ]; then
      tpgUnique="${tpgi}"
    else
      tpgUnique="${tpgUnique}${newline}${tpgi}"
    fi
  fi
done <<< "$TPG_BREAKING"


breakingchanges=""
if [ "$tpgUnique" != "" ]; then
  tpgUnique=$(sed 's/^/\* /' <<< "$tpgUnique")
  breakingchanges="${breakingchanges}${tpgUnique}${newline}"
fi

if [ "$TPGB_BREAKING" != "" ]; then
  tpgbBreaking=$(sed 's/^/\* /' <<< "$TPGB_BREAKING")
  breakingchanges="${breakingchanges}${tpgbBreaking}${newline}"
fi

if [ "$breakingchanges" != "" ]; then
message="## Breaking Change(s) Detected
The following breaking change(s) were detected within your pull request.

${breakingchanges}

If you believe this detection to be incorrect please raise the concern with your reviewer. If you intend to make this change you will need to wait for a [major release](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#example-major-number-increments) window. An \`override-breaking-change\` label can be added to allow merging.
"
fi

echo "$message"


