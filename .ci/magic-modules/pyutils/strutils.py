import re

def find_dependency_urls_in_comment(body):
  """Util to parse downstream dependencies from a given comment body.

  Example:
    $ find_dependency_urls_in_comment(\"""
      This is a comment on an MM PR.

      depends: https://github.com/ownerFoo/repoFoo/pull/100
      depends: https://github.com/ownerBar/repoBar/pull/10
    \""")
    [https://github.com/ownerFoo/repoFoo/pull/100,
     https://github.com/ownerBar/repoBar/pull/10]

  Args:
    body (string): Text of comment in upstream PR

  Returns:
    List of PR URLs found.
  """
  return re.findall(
    r'^depends: (https://github.com/[^\s]*)', body, re.MULTILINE)

def parse_github_url(gh_url):
  """Util to parse Github repo/PR id from a Github PR URL.

  Args:
    gh_url (string): URL of Github pull request.

  Returns:
    Tuple of (repo name, pr number)
  """
  matches = re.match(r'https://github.com/([\w-]+/[\w-]+)/pull/(\d+)', gh_url)
  if matches:
    repo, prnum = matches.groups()
    return (repo, int(prnum))
  return None