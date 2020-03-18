#!/bin/bash

set -x

ASSIGNEE=$(shuf -n 1 <(printf "danawillow\nrambleraptor\nemilymye\nrileykarson\nSirGitsalot\nslevenick\nchrisst\nc2thorn\nndmckinley"))

cat > comment/pr_comment << EOF
Hello!  I am a robot who works on Magic Modules PRs.

I have detected that you are a community contributor, so your PR will be assigned to someone with a commit-bit on this repo for initial review.

Thanks for your contribution!  A human will be with you soon.

@$ASSIGNEE, please review this PR or find an appropriate assignee.
EOF

# Something is preventing the magician from actually assigning the PRs.
# Leave this part in so we know what was supposed to happen, but the real
# logic is above.
echo $ASSIGNEE > comment/assignee
cat comment/assignee
