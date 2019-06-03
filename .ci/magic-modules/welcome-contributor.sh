#!/bin/bash

set -x

cat > comment/pr_comment << EOF
Hello!  I am a robot who works on Magic Modules PRs.

I have detected that you are a community contributor, so your PR will be assigned to someone with a commit-bit on this repo for initial review.  They will authorize it to run through our CI pipeline, which will generate downstream PRs.

Thanks for your contribution!  A human will be with you soon.
EOF

shuf -n 1 <(printf "danawillow\nrambleraptor\nchrisst\nrileykarson\nSirGitsalot\nslevenick\ntysen") > comment/assignee

cat comment/assignee
