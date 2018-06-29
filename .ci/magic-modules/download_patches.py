#!/usr/bin/env python
import get_downstream_prs
import itertools
import re
import operator
import os
import urllib

from github import Github

if __name__ == '__main__':
  g = Github(os.environ.get('GH_TOKEN'))
  open_pulls = g.get_repo('GoogleCloudPlatform/magic-modules').get_pulls(state='open')
  depends = [item for sublist in [get_downstream_prs.get_github_dependencies(g, open_pull.number) for open_pull in open_pulls] for item in sublist]
  parsed_dependencies = [re.match(r'https://github.com/([\w-]+/[\w-]+)/pull/(\d+)', d).groups() for d in depends]
  for r, pulls in itertools.groupby(parsed_dependencies, key=operator.itemgetter(0)):
    repo = g.get_repo(r)
    for pull in pulls:
      pr = repo.get_pull(int(pull[1]))
      print 'Checking %s to see if it should be downloaded.' % (pr,)
      if pr.is_merged():
        download_location = os.path.join('./patches', pull[0], pull[1] + '.patch')
        if not os.path.exists(os.path.dirname(download_location)):
          os.makedirs(os.path.dirname(download_location))
        urllib.urlretrieve(pr.patch_url, download_location)
