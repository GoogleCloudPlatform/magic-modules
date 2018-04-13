#!/usr/bin/env python
import functools
import os
import re
import sys
from github import Github

if __name__ == '__main__':
  g = Github(os.environ.get('GH_TOKEN'))
  all_prs = functools.reduce(lambda x,y: x+re.findall(r'\ndepends: (https://github.com/.*)', y.body),
      g.get_repo('GoogleCloudPlatform/magic-modules').get_pull(int(sys.argv[1])).get_issue_comments(), [])
  for downstream_pr in all_prs:
    print downstream_pr
