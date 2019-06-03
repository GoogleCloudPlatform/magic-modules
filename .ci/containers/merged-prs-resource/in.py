#! /usr/local/bin/python
from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

from absl import app
import json
import sys
import os
from github import Github
import urllib

def main(argv):
  in_json = json.load(sys.stdin)
  g = Github(in_json['source']['token'])
  version = in_json.get('version', {})
  for repo_name, pr_numbers in version.iteritems():
    repo = g.get_repo(repo_name)
    if not pr_numbers: continue
    for pr_number in pr_numbers.split(','):
      download_location = os.path.join(argv[1], repo_name, pr_number + '.patch')
      if not os.path.exists(os.path.dirname(download_location)):
        os.makedirs(os.path.dirname(download_location))
      pr = repo.get_pull(int(pr_number))
      urllib.urlretrieve(pr.patch_url, download_location)
  print(json.dumps({"version": version}))

if __name__ == '__main__':
  app.run(main)
