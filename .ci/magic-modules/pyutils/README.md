# Magic Modules CI Utils

This directory manages all Python utils that the Magician uses to take upstream Magic Module PRs and generate and manage PRs in various downstream repos.

What this shouldn't contain:

- Python scripts called directly by Concourse jobs.
- Non-Python code

## Tests

Currently we use the standard [unittest](https://docs.python.org/3/library/unittest.html) library. Because CI development is mostly done locally on your developer machine before being directly deployed, these tests are run manually.

This section reviews running/writing tests for someone fairly new to Python/unittest, so some of this information is just from unittest docs.

### Running tests

Set a test environment variable to make calls to Github:
```
export TEST_GITHUB_TOKEN=...
```

Otherwise, tests calling Github will be ignored (or likely be rate-limited).
```
cd pyutils

python -m unittest discover -p "*_test.py"
python ./changelog_utils_test.py
```

Read [unittest](https://docs.python.org/3/library/unittest.html#command-line-interface) docs to see how to run tests at finer granularity.

*NOTE*: Don't forget to delete .pyc files if you feel like tests aren't reflecting your changes!

### Writing Tests:

This is mostly a very shallow review of unittest, but your test should inherit from the `unittest.TestCase` class in some way (i.e. we haven't had the need to write our own TestCase-inheriting Test class but feel free to in the future if needed).

```
class MyModuleTest(unittest.TestCase):
```

Make sure to include the following at the bottom of your test file, so it defaults to running the tests in this file if run as a normal Python script.
```
if __name__ == '__main__':
  unittest.main()
```



