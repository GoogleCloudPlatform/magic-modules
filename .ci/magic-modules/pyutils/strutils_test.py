from strutils import *
import unittest
import os
from github import Github


class TestStringUtils(unittest.TestCase):
  def test_find_dependency_urls(self):
    test_urls = [
      "https://github.com/repo-owner/repo-A/pull/1",
      "https://github.com/repo-owner/repo-A/pull/2",
      "https://github.com/repo-owner/repo-B/pull/3",
    ]
    test_body = "".join(["\ndepends: %s\n" % u for u in test_urls])
    result = find_dependency_urls_in_comment(test_body)
    self.assertEquals(len(result), len(test_urls),
      "expected %d urls to be parsed from comment" % len(test_urls))
    for test_url in test_urls:
      self.assertIn(test_url, result)

  def test_parse_github_url(self):
    test_cases = {
      "https://github.com/repoowner/reponame/pull/1234": ("repoowner/reponame", 1234),
      "not a real url": None,
    }
    for k in test_cases:
      result = parse_github_url(k)
      expected = test_cases[k]
      if not expected:
        self.assertIsNone(result, "expected None, got %s" % result)
      else:
        self.assertEquals(result[0], expected[0])
        self.assertEquals(int(result[1]), expected[1])

  def test_get_release_note(self):
    upstream_body = """
      ```releasenote
      This is a release note
      ```
    """
    test_cases = {
      ("releasenote text not found", ""),
      ("""
        Empty release note:
        ```releasenote
        ```
        """, ""),
      ("""
        Random code block
        ```
        This is not a release note
        ```
        """, ""),
      ("""
        Empty release note with non-empty code block:
        ```releasenote
        ```

        ```
        This is not a release note
        ```
        """, ""),
      ("""
        Empty code block with non-empty release note:

        ```invalid
        ```

        ```releasenote
        This is a release note
        ```
        """, "This is a release note"),
      ("""```releasenote
        This is a release note
        ```
        """, "This is a release note"),
    }
    for k, v in test_cases:
      self.assertEqual(get_release_note(k), v)

  def test_set_release_note(self):
    downstream_body = """
      All of the blocks below should be replaced

      ```releasenote
      This should be replaced
      ```

      More text

      ```releasenote
      ```

      ```test
      ```
      """
    release_note = "The release note was replaced"

    replaced = set_release_note(release_note, downstream_body)
    self.assertIn(
      "```releasenote\nThe release note was replaced\n```\n",
      replaced)

    self.assertEqual(len(re.findall("```releasenote", replaced)), 1,
      "expected only one release note block in text. Result:\n%s" % replaced)

    self.assertNotIn("This should be replaced", replaced)
    self.assertIn("All of the blocks below should be replaced\n", replaced)
    self.assertIn("More text\n", replaced)

  def test_find_prefixed_labels(self):
    self.assertFalse(find_prefixed_labels([], "test: "))
    self.assertFalse(find_prefixed_labels(["", ""], "test: "))
    labels = find_prefixed_labels(["foo", "bar"], "")
    self.assertIn("foo", labels)
    self.assertIn("bar", labels)

    test_labels = [
      "test: foo",
      "test: bar",
      # Not valid changelog labels
      "not a changelog label",
      "test: "
    ]
    result = find_prefixed_labels(test_labels, prefix="test: ")

    self.assertEqual(len(result), 2, "expected only 2 labels returned")
    self.assertIn("test: foo", result)
    self.assertIn("test: bar", result)
if __name__ == '__main__':
    unittest.main()


