#!/usr/bin/env python
import mistune
import sys
import bs4
import re
import argparse

class PrDescriptionTree(object):
  class PrDescriptionNode(object):
    def __init__(self, title, body, parent, depth):
      self.title = title
      self.body = body
      self.parent = parent
      self.depth = depth

  def __init__(self, pr_body):
    self._nodes = {}
    m = mistune.markdown(pr_body)
    bs = bs4.BeautifulSoup(m, "html.parser")
    current_element = bs.hr
    current_node = None
    stack = []

    while current_element is not None:
      # Create node if it's a good time to do that.
      if current_element.name and re.match('h[1-4]', current_element.name):
        # Store previous node if we have one.
        if current_node:
          stack.append(current_node)

        # Clear stack of nodes that don't need to be there any longer.
        new_depth = int(current_element.name[1:])
        while stack and new_depth <= stack[-1].depth:
          stack.pop()
        current_node = self.PrDescriptionNode(title=current_element.contents[0],
            body='', parent=stack[-1] if stack else None, depth=new_depth)
        self._nodes[current_node.title] = current_node
      elif current_node:
        current_node.body += current_element.string
      current_element = current_element.next_sibling

  def __getitem__(self, name):
    if name not in self._nodes:
      raise KeyError("%s not found in tree - existing nodes are '%s'." % (name, self._nodes.keys()))
    n = self._nodes[name]
    while n is not None:
      if n.body.strip():
        return n.body.strip()
      n = n.parent
    raise ValueError("No parents of '%s' contained content!" % name)

if __name__ == '__main__':
  parser = argparse.ArgumentParser(description="Extract a description from the full PR description.")
  parser.add_argument("--tag", type=str, required=True)
  args = parser.parse_args()
  print PrDescriptionTree(sys.stdin.read())[args.tag]
