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
	root                 *node
	nestedBlockParagraph map[string]string
}

type node struct {
	name        string
	isAttribute bool
	children    []*node
	text        string
}

func NewParser() *DocumentParser {
	return &DocumentParser{
		nestedBlockParagraph: make(map[string]string),
	}
}

func (d *DocumentParser) Arguments() []string {
	var paths []string
	traverse(
		&paths,
		"",
		func(n *node) bool {
			return !n.isAttribute
		},
		d.root,
	)
	sort.Strings(paths)
	return paths
}

func traverse(paths *[]string, path string, shouldVisit func(*node) bool, n *node) {
	if n == nil {
		return
	}
	if !shouldVisit(n) {
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
		traverse(paths, curPath, shouldVisit, c)
	}
}

func (d *DocumentParser) Attributes() []string {
	if d == nil || d.root == nil {
		return nil
	}
	var attributes []string
	for _, c := range d.root.children {
		if !c.isAttribute {
			continue
		}
		attributes = append(attributes, c.name)
	}
	sort.Strings(attributes)
	return attributes
}

// Parse parse a resource document markdown's arguments and attributes section.
// It expects the markdown to contain specific format:
// - Section titles are identified like "## abcdefg".
// - Each item is identified like "* `abcdefg`"".
// - Nested objects are identified by a starting line contains <a name="nested_abcdefg">.
// - Attributes reference do not have nested layers.
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
		if err := d.parseArgument(argument); err != nil {
			return err
		}
	}
	if len(attribute) != 0 {
		if err := d.parseAttribute(attribute); err != nil {
			return err
		}
	}

	return nil
}

func (d *DocumentParser) parseArgument(input string) error {
	parts := strings.Split(input, "\n"+nestedObjectSeparator)
	for _, p := range parts[1:] {
		nestedName, err := findPattern(nestedObjectSeparator+p, nestedLinkPattern)
		if err != nil {
			return err
		}
		if nestedName == "" {
			return fmt.Errorf("could not find nested object name in %s", nestedObjectSeparator+p)
		}
		d.nestedBlockParagraph[nestedName] = p
	}
	return d.bfs(parts[0])
}

func (d *DocumentParser) bfs(input string) error {
	d.root = &node{
		name: "",
		text: input,
	}
	queue := []*node{d.root}

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
				nestedName, err := findNestedName(text)
				if err != nil {
					return err
				}
				newNode := &node{
					name: itemName,
				}
				cur.children = append(cur.children, newNode)
				if text, ok := d.nestedBlockParagraph[nestedName]; ok {
					newNode.text = text
					queue = append(queue, newNode)
				}
			}

		}
		queue = queue[l:]
	}
	return nil
}

func (d *DocumentParser) parseAttribute(input string) error {
	items := strings.Split(input, "\n"+listItemSeparator)
	for _, item := range items[1:] {
		itemName, err := findItemName(listItemSeparator + item)
		if err != nil {
			return err
		}
		itemName = strings.ReplaceAll(itemName, ".0.", ".")
		newNode := &node{
			name:        itemName,
			isAttribute: true,
		}
		d.root.children = append(d.root.children, newNode)
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
