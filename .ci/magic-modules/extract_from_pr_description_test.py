import unittest
import extract_from_pr_description

class TestExtraction(unittest.TestCase):
    def setUp(self):
      self.tree = extract_from_pr_description.PrDescriptionTree("""
A summary of changes goes here!

-----------------------------------------------------------------
# all
Foo
Bar
Baz
## terraform
Bar
## puppet
Baz
### puppet-dns
Qux
### puppet-compute
## chef
""")

    def testEmpty(self):
      self.assertEqual("Foo\nBar\nBaz", self.tree['chef'])

    def testThreeDeepEmpty(self):
      self.assertEqual("Baz", self.tree['puppet-compute'])

    def testThreeDeep(self):
      self.assertEqual("Qux", self.tree['puppet-dns'])

    def testTwoDeep(self):
      self.assertEqual("Bar", self.tree['terraform'])

if __name__ == '__main__':
    unittest.main()
