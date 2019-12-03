#!/usr/bin/env python
import os
import urllib
from github import Github
from pyutils import downstreams

def get_merged_patches(gh):
  """Download all merged patches for open upstream PRs.

  Args:
    gh: Github client to make calls to Github with.
  """
  open_pulls = gh.get_repo('GoogleCloudPlatform/magic-modules')\
                 .get_pulls(state='open')
  for open_pr in open_pulls:
    print 'Downloading patches for upstream PR %d...' % open_pr.number
    parsed_urls = downstreams.get_parsed_downstream_urls(gh, open_pr.number)
    for repo_name, pulls in parsed_urls:
      repo = gh.get_repo(repo_name)
      for r, pr_num in pulls:
          print 'Check to see if %s/%s is merged and should be downloaded\n' % (
            r, pr_num)
          downstream_pr = repo.get_pull(int(pr_num))
          if downstream_pr.is_merged():
            download_patch(r, downstream_pr)

def download_patch(repo, pr):
  """Download merged downstream PR patch.

  Args:
    pr: Github Pull request to download patch for
  """
  download_location = os.path.join('./patches', repo_name, '%d.patch' % pr.id)
  print download_location
  # Skip already downloaded patches
  if os.path.exists(download_location):
    return

  if not os.path.exists(os.path.dirname(download_location)):
      os.makedirs(os.path.dirname(download_location))
  urllib.urlretrieve(pr.patch_url, download_location)

if __name__ == '__main__':
  gh = Github(os.environ.get('GH_TOKEN'))
  get_merged_patches(gh)
