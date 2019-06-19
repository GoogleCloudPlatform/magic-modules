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

if __name__ == '__main__':
    unittest.main()


