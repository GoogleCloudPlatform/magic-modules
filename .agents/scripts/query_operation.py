#!/usr/bin/env python3
import json
import re
import subprocess
import sys
import urllib.request

op_url = sys.argv[1]
assert re.match(r"^https://[a-z0-9.-]*googleapis\.com/[a-zA-Z0-9_./-]+$", op_url)
token = subprocess.run(
    ["gcloud", "auth", "print-access-token"],
    check=True,
    capture_output=True,
    text=True,
).stdout.strip()
req = urllib.request.Request(op_url, headers={"Authorization": f"Bearer {token}"})
with urllib.request.urlopen(req) as resp:
    print(json.dumps(json.load(resp), indent=2))
