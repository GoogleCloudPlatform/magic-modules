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

  def test_get_release_notes(self):
    test_cases = [
      ("releasenote text not found", []),
      (
"""Empty release note:
```release-note:test

```
""", []),
      ("""
Random code block
```
This is not a release note
```
""", []),
      ("""
Empty release note with non-empty code block:
```release-note:test

```

```
This is not a release note
```
""", []),
      ("""
Empty code block with non-empty release note:

```invalid

```

```release-note:test
This is a release note
```
""", [("release-note:test", "This is a release note")]),
      ("""
Single release notes
```release-note:test
This is a release note
```
""", [("release-note:test", "This is a release note")])
      # ("""
      #   Multiple release notes
      #   ```release-note:foo
      #   note foo
      #   ```

      #   ```release-note:bar
      #   note bar
      #   ```

      #   ```release-note:baz
      #   note baz
      #   ```
      #   """, [
      #     ("release-note:foo", "note foo"),
      #     ("release-note:bar", "note bar"),
      #     ("release-note:baz", "note baz"),
      # ]),
    ]
    for k, expected in test_cases:
      actual = get_release_notes(k)
      self.assertEqual(len(actual), len(expected),
        "test %s\n: expected %d items, got %d: %s" % (k, len(expected), len(actual), actual))
      for idx, note_tuple in enumerate(expected):
        self.assertEqual(actual[idx][0], note_tuple[0],
          "test %s\n: expected note type %s, got %s" % (
            k, note_tuple[0], actual[idx][0]))

        self.assertEqual(actual[idx][1], note_tuple[1],
          "test %s\n: expected note type %s, got %s" % (
            k, note_tuple[1], actual[idx][1]))


  def test_set_release_notes(self):
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
    release_notes = [
      ("release-note:foo", "new message foo"),
      ("release-note:bar", "new message bar"),
    ]

    replaced = set_release_notes(release_notes, downstream_body)

    # Existing non-code-block text should still be in body
    self.assertIn("All of the blocks below should be replaced\n", replaced)
    self.assertIn("More text\n", replaced)

    # New release notes should have been added.
    self.assertIn("```release-note:foo\nnew message foo\n```\n", replaced)
    self.assertIn("```release-note:bar\nnew message bar\n```\n", replaced)

    # Old release notes and code blocks should be removed.
    self.assertEqual(len(re.findall("```.+\n", replaced)), 2,
      "expected only two release note blocks in text. Result:\n%s" % replaced)
    self.assertNotIn("This should be replaced", replaced)


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


