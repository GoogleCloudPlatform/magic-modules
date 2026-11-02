package documentparser

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	// * `xxx` or - `xxx` at the start of a line, possibly indented
	fieldNameRegex      = regexp.MustCompile("(?m)^([ \\t]*)[*-][ \\t]+`([a-z0-9_\\./]+)`")
	nestedObjectRegex   = regexp.MustCompile(`<a\s+name="([a-z0-9_]+)">`) // <a name="xxx">
	nestedHashTagRegex  = regexp.MustCompile(`\(#(nested_[a-z0-9_]+)\)`)  // #(nested_xxx)
	horizontalLineRegex = regexp.MustCompile("- - -|-{3,}")               // - - - or ---

	sectionSeparator = "## "
)

// DocumentParser parse *.html.markdown resource doc files.
type DocumentParser struct {
	root        *node
	nestedBlock map[string]string
}

type node struct {
	name     string
	children []*node
	text     string
}

func NewParser() *DocumentParser {
	return &DocumentParser{
		nestedBlock: make(map[string]string),
	}
}

func (d *DocumentParser) FlattenFields() []string {
	var paths []string
	traverse(
		&paths,
		"",
		d.root,
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

// Parse parse a resource document markdown's arguments and attributes section.
// The parsed file format is defined in mmv1/templates/terraform/resource.html.markdown.tmpl.
func (d *DocumentParser) Parse(src []byte) error {
	var argument, attribute, ephemeralAttribute string
	for _, p := range strings.Split(string(src), "\n"+sectionSeparator) {
		if strings.HasPrefix(p, "Attributes Reference") {
			attribute = p
		}
		if strings.HasPrefix(p, "Argument Reference") {
			argument = p
		}
		if strings.HasPrefix(p, "Ephemeral Attributes Reference") {
			ephemeralAttribute = p
		}
	}
	for _, text := range []string{argument, attribute, ephemeralAttribute} {
		if len(text) != 0 {
			sections := horizontalLineRegex.Split(text, -1)
			var allTopLevelFieldSections string
			for _, part := range sections {
				topLevelPropertySection, err := d.extractNestedObject(part)
				if err != nil {
					return err
				}
				allTopLevelFieldSections += topLevelPropertySection
			}
			root := &node{
				text: allTopLevelFieldSections,
			}
			if err := d.bfs(root, d.nestedBlock); err != nil {
				return err
			}
			if d.root == nil {
				d.root = root
			} else {
				d.root.children = append(d.root.children, root.children...)
			}
		}
	}
	return nil
}

func (d *DocumentParser) extractNestedObject(input string) (string, error) {
	parts := splitWithRegexp(input, nestedObjectRegex)
	for _, p := range parts[1:] {
		nestedName := findPattern(p, nestedObjectRegex)
		if nestedName == "" {
			return "", fmt.Errorf("could not find nested object name in %s", p)
		}
		d.nestedBlock[nestedName] = p
	}
	return parts[0], nil
}

func (d *DocumentParser) bfs(root *node, nestedBlock map[string]string) error {
	if root == nil {
		return fmt.Errorf("no node to visit")
	}
	queue := []*node{root}

	for len(queue) > 0 {
		l := len(queue)
		for _, cur := range queue {
			parts := splitFieldBullets(cur.text)
			for _, p := range parts[1:] {
				p = strings.ReplaceAll(p, "\n", "")
				match := fieldNameRegex.FindStringSubmatch(p)
				if match == nil {
					return fmt.Errorf("could not find field name in %s", p)
				}
				fieldName := match[2]
				// There is a special case in some hand written resource eg. in compute_instance, where its attributes is in a.0.b.0.c format.
				fieldName = strings.ReplaceAll(fieldName, ".0.", ".")
				newNode := &node{
					name: fieldName,
				}
				cur.children = append(cur.children, newNode)

				nestedHashTag := findPattern(p, nestedHashTagRegex)
				if text, ok := nestedBlock[nestedHashTag]; ok {
					newNode.text = text
					queue = append(queue, newNode)
				}
			}

		}
		queue = queue[l:]
	}
	return nil
}

// Split text at the bullets of the fields it documents. The fields are the
// least indented bullets at the start of a line, usually unindented. More
// indented bullets, such as lists of example values, and backticked names in
// the middle of a line are part of the description of the field before them.
func splitFieldBullets(text string) []string {
	matches := fieldNameRegex.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return []string{text}
	}
	minIndent := -1
	for _, m := range matches {
		if indent := m[3] - m[2]; minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	var parts []string
	start := 0
	for _, m := range matches {
		if m[3]-m[2] != minIndent {
			continue
		}
		parts = append(parts, text[start:m[0]])
		start = m[0]
	}
	return append(parts, text[start:])
}

func findPattern(text string, re *regexp.Regexp) string {
	match := re.FindStringSubmatch(text)
	if match != nil {
		return match[1]
	}
	return ""
}

func splitWithRegexp(text string, re *regexp.Regexp) []string {
	matches := re.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return []string{text}
	}
	var parts []string
	start := 0
	for _, match := range matches {
		end := match[0]

		parts = append(parts, text[start:end])
		start = end
	}
	parts = append(parts, text[start:])
	return parts
}
