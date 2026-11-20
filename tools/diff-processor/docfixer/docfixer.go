// Package docfixer makes mechanical fixes to handwritten resource docs so the
// missing doc detector can read them, using the resource's schema to resolve
// which block each heading documents.
//
// It only makes changes that don't alter wording:
//   - add a `<a name="nested_..."></a>` anchor to an anchorless
//     "The `x` block supports:" heading, and
//   - link the parent field's bullet to that anchor, either by turning
//     "Structure is documented below" into a link or by appending
//     " Structure is [documented below](#nested_...)." to the bullet's first
//     paragraph.
//
// Anything it can't resolve unambiguously is reported, not changed.
package docfixer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	headingRe      = regexp.MustCompile("^(<a\\s+name=\"([^\"]+)\"></a>\\s*)?The `([a-z0-9_]+)` (?:block|blocks|subblock) (?:supports|contains|holds|includes|support)\\b")
	anyHeadingRe   = regexp.MustCompile("^(?:<a\\s+name=\"[^\"]+\"></a>\\s*)?The `[^`]+` ")
	bulletRe       = regexp.MustCompile("^[*-][ \\t]+`([a-z0-9_./]+)`")
	sectionRe      = regexp.MustCompile(`^##[ \t]+(.*)$`)
	fieldSectionRe = regexp.MustCompile(`^(?:Argument|Attributes|Ephemeral Attributes) Reference\b`)
	hrRe           = regexp.MustCompile(`^[ \t]*(?:- - -|-{3,})[ \t]*$`)
	anchorLineRe   = regexp.MustCompile(`^<a\s+name="`)
	plainBelowRe   = regexp.MustCompile(`(?i)structure is documented below`)
	linkedBelowRe  = regexp.MustCompile(`documented below\]\(#([a-z0-9_-]+)\)`)
	indentedListRe = regexp.MustCompile(`^\s+[*-]\s`)
)

// Block is a nested block in a resource schema.
type Block struct {
	Path     string // dotted path, e.g. "serial_pipeline.stages"
	Name     string // last path element
	Parent   string // dotted path of the parent block, "" for top level
	Children map[string]bool
}

// Blocks returns all nested blocks in the schema, keyed by dotted path.
func Blocks(r *schema.Resource) map[string]*Block {
	out := make(map[string]*Block)
	var walk func(prefix string, s map[string]*schema.Schema)
	walk = func(prefix string, s map[string]*schema.Schema) {
		for name, f := range s {
			elem, ok := f.Elem.(*schema.Resource)
			if !ok {
				continue
			}
			path := name
			if prefix != "" {
				path = prefix + "." + name
			}
			b := &Block{Path: path, Name: name, Parent: prefix, Children: make(map[string]bool)}
			for c := range elem.Schema {
				b.Children[c] = true
			}
			out[path] = b
			walk(path, elem.Schema)
		}
	}
	walk("", r.Schema)
	return out
}

type heading struct {
	line    int // index into lines
	end     int // exclusive end of the heading's region
	name    string
	anchor  string // existing anchor, if any
	bullets []string
	path    string // resolved block path
	shared  []string
}

// Result describes what Fix did.
type Result struct {
	Text    string
	Changes []string
	Notes   []string // things it couldn't fix
}

// Fix applies the mechanical fixes to a doc for the resource described by r.
func Fix(doc string, r *schema.Resource) Result {
	crlf := strings.Contains(doc, "\r\n")
	text := strings.ReplaceAll(doc, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	res := Result{}
	blocks := Blocks(r)

	// Mark which lines are in field reference sections.
	inField := make([]bool, len(lines))
	cur := false
	for i, l := range lines {
		if m := sectionRe.FindStringSubmatch(l); m != nil {
			cur = fieldSectionRe.MatchString(m[1])
		}
		inField[i] = cur
	}

	isBoundary := func(i int) bool {
		l := lines[i]
		return sectionRe.MatchString(l) || hrRe.MatchString(l) || anchorLineRe.MatchString(l) || anyHeadingRe.MatchString(l)
	}

	// Find headings and their regions.
	var headings []*heading
	for i, l := range lines {
		if !inField[i] {
			continue
		}
		m := headingRe.FindStringSubmatch(l)
		if m == nil {
			if anyHeadingRe.MatchString(l) && !strings.HasPrefix(l, "<a") {
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: heading not in a recognized form, left alone: %q", i+1, l))
			}
			continue
		}
		h := &heading{line: i, name: m[3], anchor: m[2]}
		j := i + 1
		for j < len(lines) && inField[j] && !isBoundary(j) {
			if bm := bulletRe.FindStringSubmatch(lines[j]); bm != nil {
				h.bullets = append(h.bullets, bm[1])
			}
			j++
		}
		h.end = j
		headings = append(headings, h)
	}

	// Root bullets: column-0 bullets in field sections that aren't inside any heading region.
	inRegion := make([]bool, len(lines))
	for _, h := range headings {
		for k := h.line; k < h.end; k++ {
			inRegion[k] = true
		}
	}
	// Also treat regions started by bare anchors (not "The `x` block" headings) as non-root.
	for i, l := range lines {
		if inField[i] && anchorLineRe.MatchString(l) && headingRe.FindStringSubmatch(l) == nil {
			for k := i; k < len(lines) && inField[k] && (k == i || !isBoundary(k)); k++ {
				inRegion[k] = true
			}
		}
	}

	// bulletLine finds the line of the bullet for `name` among the given candidate lines.
	findBullet := func(name string, from, to int, root bool) []int {
		var out []int
		for k := from; k < to; k++ {
			if !inField[k] || (root && inRegion[k]) {
				continue
			}
			if bm := bulletRe.FindStringSubmatch(lines[k]); bm != nil && bm[1] == name {
				out = append(out, k)
			}
		}
		return out
	}

	// Resolve each heading to a block path.
	resolved := map[string]*heading{} // path -> heading
	parentResolved := func(b *Block) bool { return b.Parent == "" || resolved[b.Parent] != nil }
	parentBulletLine := func(b *Block) int {
		var cands []int
		if b.Parent == "" {
			cands = findBullet(b.Name, 0, len(lines), true)
		} else if ph := resolved[b.Parent]; ph != nil {
			cands = findBullet(b.Name, ph.line+1, ph.end, false)
		}
		if len(cands) == 1 {
			return cands[0]
		}
		return -1
	}
	matches := func(h *heading, b *Block) bool {
		if b.Name != h.name || len(h.bullets) == 0 {
			return false
		}
		for _, c := range h.bullets {
			if !b.Children[strings.SplitN(c, ".", 2)[0]] {
				return false
			}
		}
		return true
	}
	for changed := true; changed; {
		changed = false
		for _, h := range headings {
			if h.path != "" {
				continue
			}
			var cands []*Block
			for _, b := range blocks {
				if resolved[b.Path] == nil && matches(h, b) && parentResolved(b) {
					cands = append(cands, b)
				}
			}
			if len(cands) == 0 {
				continue
			}
			// Prefer the candidate whose parent bullet is the closest one before the heading.
			best, bestLine, tie := (*Block)(nil), -1, false
			for _, b := range cands {
				pl := parentBulletLine(b)
				if pl < 0 || pl > h.line {
					continue
				}
				if pl > bestLine {
					best, bestLine, tie = b, pl, false
				} else if pl == bestLine {
					tie = true
				}
			}
			if best == nil && len(cands) == 1 {
				best = cands[0]
			}
			if best == nil || tie {
				continue
			}
			h.path = best.Path
			resolved[best.Path] = h
			changed = true
		}
	}

	// Blocks with the same schema shape under other parents that have no heading of their own
	// share the heading (the doc documents the shape once).
	for _, b := range blocks {
		if resolved[b.Path] != nil || !parentResolved(b) || parentBulletLine(b) < 0 {
			continue
		}
		var share *heading
		for _, h := range headings {
			if h.path != "" && h.name == b.Name && matches(h, b) && sameChildren(blocks[h.path], b) {
				if share != nil && share != h {
					share = nil
					break
				}
				share = h
			}
		}
		if share != nil {
			share.shared = append(share.shared, b.Path)
			resolved[b.Path] = share
		}
	}

	existingAnchors := map[string]bool{}
	for _, l := range lines {
		for _, m := range regexp.MustCompile(`<a\s+name="([^"]+)"`).FindAllStringSubmatch(l, -1) {
			existingAnchors[m[1]] = true
		}
	}

	// Assign anchors and edit.
	edits := map[int]string{} // line -> new content
	for _, h := range headings {
		if h.path == "" {
			if h.anchor == "" {
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: couldn't match `%s` heading (bullets %v) to exactly one schema block; left alone", h.line+1, h.name, h.bullets))
			}
			continue
		}
		anchor := h.anchor
		if anchor == "" || !strings.HasPrefix(anchor, "nested_") {
			want := "nested_" + strings.ReplaceAll(h.path, ".", "_")
			if existingAnchors[want] {
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: anchor %q already exists; left `%s` heading alone", h.line+1, want, h.name))
				continue
			}
			if h.anchor == "" {
				edits[h.line] = fmt.Sprintf(`<a name="%s"></a>`, want) + lines[h.line]
				res.Changes = append(res.Changes, fmt.Sprintf("line %d: added anchor %s to `%s` heading", h.line+1, want, h.name))
			} else {
				// Non-nested_ anchor: rename it and its links.
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: heading has non-nested anchor %q; left alone", h.line+1, h.anchor))
				continue
			}
			existingAnchors[want] = true
			anchor = want
		}
		for _, p := range append([]string{h.path}, h.shared...) {
			b := blocks[p]
			pl := parentBulletLine(b)
			if pl < 0 {
				res.Notes = append(res.Notes, fmt.Sprintf("`%s`: couldn't find exactly one parent bullet for `%s`; not linked", p, b.Name))
				continue
			}
			msg, ok := linkBullet(lines, edits, pl, anchor, existingAnchors, isBoundary)
			if ok {
				if msg != "" {
					res.Changes = append(res.Changes, fmt.Sprintf("line %d: %s", pl+1, msg))
				}
			} else {
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: %s", pl+1, msg))
			}
		}
	}

	keys := make([]int, 0, len(edits))
	for k := range edits {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		lines[k] = edits[k]
	}
	out := strings.Join(lines, "\n")
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	res.Text = out
	return res
}

func sameChildren(a, b *Block) bool {
	if len(a.Children) != len(b.Children) {
		return false
	}
	for c := range a.Children {
		if !b.Children[c] {
			return false
		}
	}
	return true
}

// linkBullet makes the bullet starting at line pl link to anchor.
func linkBullet(lines []string, edits map[int]string, pl int, anchor string, anchors map[string]bool, isBoundary func(int) bool) (string, bool) {
	get := func(k int) string {
		if e, ok := edits[k]; ok {
			return e
		}
		return lines[k]
	}
	// Extent of the bullet: until the next column-0 bullet or boundary.
	end := pl + 1
	for end < len(lines) && !bulletRe.MatchString(lines[end]) && !isBoundary(end) && !strings.HasPrefix(lines[end], "```") {
		end++
	}
	for k := pl; k < end; k++ {
		if m := linkedBelowRe.FindStringSubmatch(get(k)); m != nil {
			if m[1] == anchor {
				return "", true
			}
			if anchors[m[1]] {
				return fmt.Sprintf("bullet already links to a different anchor %q; not changed", m[1]), false
			}
			edits[k] = strings.Replace(get(k), "(#"+m[1]+")", "(#"+anchor+")", 1)
			return fmt.Sprintf("fixed broken link #%s -> #%s", m[1], anchor), true
		}
	}
	for k := pl; k < end; k++ {
		if loc := plainBelowRe.FindStringIndex(get(k)); loc != nil {
			l := get(k)
			phrase := l[loc[0]:loc[1]]
			edits[k] = l[:loc[0]] + phrase[:len(phrase)-len("documented below")] + "[documented below](#" + anchor + ")" + l[loc[1]:]
			return "linked existing \"Structure is documented below\"", true
		}
	}
	// Append to the end of the bullet's first paragraph.
	last := pl
	for k := pl + 1; k < end; k++ {
		l := get(k)
		if strings.TrimSpace(l) == "" || indentedListRe.MatchString(l) {
			break
		}
		last = k
	}
	l := strings.TrimRight(get(last), " \t")
	sep := " "
	if !strings.HasSuffix(l, ".") && !strings.HasSuffix(l, ":") && !strings.HasSuffix(l, ")") && !strings.HasSuffix(l, "-") {
		sep = ". "
	}
	if strings.HasSuffix(l, "-") || strings.HasSuffix(l, ")") && strings.HasPrefix(strings.TrimSpace(l), "(") {
		sep = " "
	}
	edits[last] = l + sep + "Structure is [documented below](#" + anchor + ")."
	return "appended link to " + anchor, true
}
