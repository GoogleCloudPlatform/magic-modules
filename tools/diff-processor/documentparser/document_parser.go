package documentparser

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const (
	nestedNamePattern = `\(#(nested_[a-z0-9_]+)\)`

	itemNamePattern   = "\\* `([a-z0-9_\\./]+)`"
	nestedLinkPattern = `<a\s+name="([a-z0-9_]+)">`

	sectionSeparator      = "## "
	nestedObjectSeparator = `<a name="nested_`
	listItemSeparator     = "* `"
)

// DocumentParser parse *.html.markdown resource doc files.
type DocumentParser struct {
	argumentRoot   *node
	attriibuteRoot *node
}

type node struct {
	name     string
	children []*node
	text     string
}

func NewParser() *DocumentParser {
	return &DocumentParser{}
}

func (d *DocumentParser) Arguments() []string {
	var paths []string
	traverse(
		&paths,
		"",
		d.argumentRoot,
	)
	sort.Strings(paths)
	return paths
}

func traverse(paths *[]string, path string, n *node) {
	if n == nil {
		return
	}
	var curPath string
	if path != "" {
		curPath = path + "." + n.name
	} else {
		curPath = n.name
	}
	if curPath != "" {
		*paths = append(*paths, curPath)
	}
	for _, c := range n.children {
		traverse(paths, curPath, c)
	}
}

func (d *DocumentParser) Attributes() []string {
	var paths []string
	traverse(
		&paths,
		"",
		d.attriibuteRoot,
	)
	sort.Strings(paths)
	return paths
}

// Parse parse a resource document markdown's arguments and attributes section.
// The parsed file format is defined in mmv1/templates/terraform/resource.html.markdown.tmpl.
func (d *DocumentParser) Parse(src []byte) error {
	var argument, attribute string
	for _, p := range strings.Split(string(src), "\n"+sectionSeparator) {
		if strings.HasPrefix(p, "Attributes Reference") {
			attribute = p
		}
		if strings.HasPrefix(p, "Argument Reference") {
			argument = p
		}
	}
	if len(argument) != 0 {
		argumentParts := strings.Split(argument, "- - -")
		for _, part := range argumentParts {
			n, err := d.parseSection(part)
			if err != nil {
				return err
			}
			if d.argumentRoot == nil {
				d.argumentRoot = n
			} else {
				d.argumentRoot.children = append(d.argumentRoot.children, n.children...)
			}
		}
	}
	if len(attribute) != 0 {
		n, err := d.parseSection(attribute)
		if err != nil {
			return err
		}
		d.attriibuteRoot = n
	}
	return nil
}

func (d *DocumentParser) parseSection(input string) (*node, error) {
	parts := strings.Split(input, "\n"+nestedObjectSeparator)
	nestedBlock := make(map[string]string)
	for _, p := range parts[1:] {
		nestedName, err := findPattern(nestedObjectSeparator+p, nestedLinkPattern)
		if err != nil {
			return nil, err
		}
		if nestedName == "" {
			return nil, fmt.Errorf("could not find nested object name in %s", nestedObjectSeparator+p)
		}
		nestedBlock[nestedName] = p
	}
	// bfs to traverse the first part without nested blocks.
	root := &node{
		text: parts[0],
	}
	if err := d.bfs(root, nestedBlock); err != nil {
		return nil, err
	}
	return root, nil
}

func (d *DocumentParser) bfs(root *node, nestedBlock map[string]string) error {
	if root == nil {
		return fmt.Errorf("no node to visit")
	}
	queue := []*node{root}

	for len(queue) > 0 {
		l := len(queue)
		for _, cur := range queue {
			// the separator should always at the beginning of the line
			items := strings.Split(cur.text, "\n"+listItemSeparator)
			for _, item := range items[1:] {
				text := listItemSeparator + item
				itemName, err := findItemName(text)
				if err != nil {
					return err
				}
				// There is a special case in some hand written resource eg. in compute_instance, where its attributes is in a.0.b.0.c format.
				itemName = strings.ReplaceAll(itemName, ".0.", ".")
				nestedName, err := findNestedName(text)
				if err != nil {
					return err
				}
				newNode := &node{
					name: itemName,
				}
				cur.children = append(cur.children, newNode)
				if text, ok := nestedBlock[nestedName]; ok {
					newNode.text = text
					queue = append(queue, newNode)
				}
			}

		}
		queue = queue[l:]
	}
	return nil
}

func findItemName(text string) (name string, err error) {
	name, err = findPattern(text, itemNamePattern)
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", fmt.Errorf("cannot find item name from %s", text)
	}
	return
}

func findPattern(text string, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	match := re.FindStringSubmatch(text)

	if match != nil {
		return match[1], nil
	}
	return "", nil
}

func findNestedName(text string) (string, error) {
	s := strings.ReplaceAll(text, "\n", "")
	return findPattern(s, nestedNamePattern)
}
