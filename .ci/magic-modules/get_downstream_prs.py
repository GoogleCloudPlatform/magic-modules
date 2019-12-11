#!/usr/bin/env python
import os
import sys
from github import Github
from pyutils import downstreams

if __name__ == '__main__':
  assert len(sys.argv) == 2, "expected a Github PR ID as argument"
  upstream_pr = int(sys.argv[1])

  downstream_urls = downstreams.get_downstream_urls(
    Github(os.environ.get('GH_TOKEN')), upstream_pr)
  for url in downstream_urls:
    print url
