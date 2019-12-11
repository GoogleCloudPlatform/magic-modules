import re
from bs4 import BeautifulSoup
import mistune

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

def get_release_notes(body):
  """Parse release note blocks from a given text block.

  Each code-block with a "release-note:..." language class.
  Example:
    ```release-note:new-resource
    a_new_resource
    ```

    ```release-note:bug
    Fixed a bug
    ```
  Args:
    body (string): PR body to pull release note block from

  Returns:
    List of tuples of (`release-note` heading, release note)
  """
  release_notes = []

  # Parse markdown and find all code blocks
  md = mistune.markdown(body)
  soup = BeautifulSoup(md, 'html.parser')
  for codeblock in soup.find_all('code'):
    block_classes = codeblock.get('class')
    if not block_classes:
      continue

    note_type = get_release_note_type_from_class(block_classes[0])
    note_text = codeblock.get_text().strip()
    if note_type and note_text:
      release_notes.append((note_type, note_text))

  return release_notes

def get_release_note_type_from_class(class_str):
  # expected class is 'lang-release-note:...' for release notes
  prefix_len = len("lang-release-note:")
  if class_str[:prefix_len] == "lang-release-note:":
    return class_str[len("lang-"):]
  return None

def set_release_notes(release_notes, body):
  """Sanitize and adds the given release note block for PR body text.

  For a given text block, removes any existing "releasenote" markdown code
  blocks and adds the given release notes at the end.

  Args:
    release_note (list(Tuple(string)): List of
      (release-note heading, release note)
    body (string): Text body to find and edit release note blocks in

  Returns:
    Modified text
  """
  edited = ""
  md = mistune.markdown(body)
  soup = BeautifulSoup(md, 'html.parser')
  for blob in soup.find_all('p'):
    edited += blob.get_text().strip() + "\n\n"

  for heading, note in release_notes:
    edited += "\n```%s\n%s\n```\n" % (heading, note.strip())
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
