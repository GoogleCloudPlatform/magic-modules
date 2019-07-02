import re

RELEASE_NOTE_RE = r'```releasenote[\s]*(?P<release_note>[^`{3}]*)```'

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

def get_release_note(body):
  """Parse release note block from a given text block.

  Finds the first markdown code block with a "releasenote" language class.
  Example:
    ```releasenote
    This is the release note
    ```

  Args:
    body (string): PR body to pull release note block from

  Returns:
    Release note if found or empty string.
  """
  m = re.search(RELEASE_NOTE_RE, body, re.MULTILINE)
  return m.groupdict("")["release_note"].strip() if m else ""

def set_release_note(release_note, body):
  """Sanitize and adds the given release note block for PR body text.

  For a given text block, removes any existing "releasenote" markdown code
  blocks and adds the given release note at the end.

  Example:
    # Set a release note
    > print set_release_note(
        "This is the new release note",
        "``releasenote\nChanges to downstream\n```\n")
    "```releasenote\nThis is the new release note\n```\n"

    # Remove for empty release note
    > print set_release_note("",
        "PR description\n```releasenote\nChanges to downstream\n```\n")
    "PR description\n"

  Args:
    release_note (string): Release note to add
    body (string): Text body to find and edit release note blocks in

  Returns:
    Modified text
  """
  edited = re.sub(RELEASE_NOTE_RE, '', body)
  release_note = release_note.strip()
  if release_note:
    edited += "\n```releasenote\n%s\n```\n" % release_note.strip()
  return edited

def find_prefixed_labels(labels, prefix):
  """Util for filtering and cleaning labels that start with a given prefix.

  Given a list of labels, find only the specific labels with the given prefix.

  Args:
    prefix: String expected to be prefix of relevant labels
    labels: List of string labels

  Return:
    Filtered labels (i.e. all labels starting with prefix)
  """
  changelog_labels = []
  for l in labels:
    l = l.strip()
    if l.startswith(prefix) and len(l) > len(prefix):
      changelog_labels.append(l)
  return changelog_labels
