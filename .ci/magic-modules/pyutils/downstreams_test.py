from downstreams import *
import unittest
import os
from github import Github

TOKEN_ENV_VAR = "TEST_GITHUB_TOKEN"

class TestUpstreamPullRequests(unittest.TestCase):
  """
    Terrible test data from scraping
    https://github.com/GoogleCloudPlatform/magic-modules/pull/1000
    TODO: If this test becomes load-bearing, mock out the Github client instead
    of using this.
  """
  TEST_PR_NUM = 1000
  EXPECTED_DOWNSTREAM_URLS = [
    "https://github.com/terraform-providers/terraform-provider-google-beta/pull/186",
    "https://github.com/terraform-providers/terraform-provider-google/pull/2591",
    "https://github.com/modular-magician/ansible/pull/142",
  ]
  EXPECTED_PARSED_DOWNSTREAMS = {
    "terraform-providers/terraform-provider-google-beta": [186],
    "terraform-providers/terraform-provider-google": [2591],
    "modular-magician/ansible": [142],
  }

  def setUp(self):
    gh_token = os.environ.get(TOKEN_ENV_VAR)
    if not gh_token:
      self.skipTest(
        "test env var %s not set, skip tests calling Github" % TOKEN_ENV_VAR)
    self.test_client = Github(gh_token)

  def test_find_unmerged_downstreams(self):
    self.assertFalse(find_unmerged_downstreams(self.test_client, self.TEST_PR_NUM))

  def test_parsed_downstream_urls(self):
    result = get_parsed_downstream_urls(self.test_client, self.TEST_PR_NUM)
    repo_cnt = 0
    for repo, pulls in result:
      # Verify each repo in result.
      self.assertIn(repo, self.EXPECTED_PARSED_DOWNSTREAMS,
        "unexpected repo %s in result" % repo)
      repo_cnt += 1

      # Verify each pull request in result.
      expected_pulls = self.EXPECTED_PARSED_DOWNSTREAMS[repo]
      pull_cnt = 0
      for repo, prid in pulls:
        self.assertIn(int(prid), expected_pulls)
        pull_cnt += 1
      # Verify exact count of pulls (here because iterator).
      self.assertEquals(pull_cnt, len(expected_pulls),
        "expected %d pull requests in result[%s]" % (len(expected_pulls), repo))

    # Verify exact count of repos (here because iterator).
    self.assertEquals(repo_cnt, len(self.EXPECTED_PARSED_DOWNSTREAMS),
        "expected %d pull requests in result[%s]" % (
          len(self.EXPECTED_PARSED_DOWNSTREAMS), repo))

  def test_downstream_urls(self):
    test_client = Github(os.environ.get(TOKEN_ENV_VAR))
    result = get_downstream_urls(self.test_client,self.TEST_PR_NUM)

    expected_len = len(self.EXPECTED_DOWNSTREAM_URLS)
    self.assertEquals(len(result), expected_len,
      "expected %d downstream urls, got %d" % (expected_cnt, len(result)))
    for url in result:
      self.assertIn(str(url), self.EXPECTED_DOWNSTREAM_URLS)


if __name__ == '__main__':
    unittest.main()


