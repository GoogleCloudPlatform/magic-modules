// Package docfixer makes mechanical fixes to handwritten resource docs so the
// missing doc detector can read them, using the resource's schema to resolve
// which block each heading documents.
//
// It only makes changes that don't alter wording:
//   - add a `<a name="nested_...">` anchor to an anchorless
//     "The `x` block supports:" heading, and
//   - link the parent field's bullet to that anchor, by fixing a broken
//     "[documented below](#...)" link, by turning "Structure is documented
//     below" into a link, or by appending
//     " Structure is [documented below](#nested_...)." to the bullet's first
//     paragraph.
//
// Anything it can't resolve unambiguously is reported, not changed.
package docfixer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	headingRe      = regexp.MustCompile("^(<a\\s+name=\"([^\"]+)\"></a>\\s*)?The `([a-z0-9_]+)` (?:block|blocks|subblock) (?:supports|contains|holds|includes|support)\\b")
	anyHeadingRe   = regexp.MustCompile("^(?:<a\\s+name=\"[^\"]+\"></a>\\s*)?The `[^`]+` ")
	bulletRe       = regexp.MustCompile("^[*-][ \\t]+`([a-z0-9_./]+)`")
	sectionRe      = regexp.MustCompile(`^##[ \t]+(.*)$`)
	fieldSectionRe = regexp.MustCompile(`^(?:Argument|Attributes|Ephemeral Attributes) Reference\b`)
	argSectionRe   = regexp.MustCompile(`^##[ \t]+Argument Reference\b`)
	hrRe           = regexp.MustCompile(`^[ \t]*(?:- - -|-{3,})[ \t]*$`)
	anchorLineRe   = regexp.MustCompile(`^<a\s+name="`)
	anchorNameRe   = regexp.MustCompile(`<a\s+name="([^"]+)"`)
	plainBelowRe   = regexp.MustCompile(`(?i)structure (?:is )?documented below`)
	linkedBelowRe  = regexp.MustCompile(`documented below\]\(#([a-z0-9_-]+)\)`)
	indentedListRe = regexp.MustCompile(`^\s+[*-]\s`)
)

// block is a nested block in a resource schema.
type block struct {
	path     string // dotted path, e.g. "serial_pipeline.stages"
	name     string // last path element
	parent   string // dotted path of the parent block, "" for top level
	children map[string]bool
}

// schemaBlocks returns all nested blocks in the schema, keyed by dotted path.
func schemaBlocks(r *schema.Resource) map[string]*block {
	out := make(map[string]*block)
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
			b := &block{path: path, name: name, parent: prefix, children: make(map[string]bool)}
			for c := range elem.Schema {
				b.children[c] = true
			}
			out[path] = b
			walk(path, elem.Schema)
		}
	}
	walk("", r.Schema)
	return out
}

// heading is a "The `x` block supports:" line and the bullets under it.
type heading struct {
	line    int // index into lines
	end     int // exclusive end of the heading's region
	name    string
	anchor  string // existing anchor, if any
	bullets []string
	path    string   // resolved block path
	shared  []string // other block paths with the same shape that this heading documents
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
	lines := strings.Split(strings.ReplaceAll(doc, "\r\n", "\n"), "\n")
	res := Result{Text: doc}
	blocks := schemaBlocks(r)

	// Docs with one Argument Reference section per variant (e.g. per
	// environment generation) repeat the same blocks in each, so a heading
	// can't be matched to its parent bullet by position.
	argSections := 0
	for _, l := range lines {
		if argSectionRe.MatchString(l) {
			argSections++
		}
	}
	if argSections > 1 {
		res.Notes = append(res.Notes, fmt.Sprintf("%d Argument Reference sections; doc left alone", argSections))
		return res
	}

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

	// Root bullets are column-0 bullets in field sections outside any heading
	// region. Regions started by bare anchors (not "The `x` block" headings)
	// count as heading regions too.
	inRegion := make([]bool, len(lines))
	for _, h := range headings {
		for k := h.line; k < h.end; k++ {
			inRegion[k] = true
		}
	}
	for i, l := range lines {
		if inField[i] && anchorLineRe.MatchString(l) && !headingRe.MatchString(l) {
			for k := i; k < len(lines) && inField[k] && (k == i || !isBoundary(k)); k++ {
				inRegion[k] = true
			}
		}
	}

	// findBullets returns the lines in [from, to) with a bullet for name.
	findBullets := func(name string, from, to int, root bool) []int {
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
	parentResolved := func(b *block) bool { return b.parent == "" || resolved[b.parent] != nil }
	// parentBulletLine returns the line of b's bullet in its parent's
	// region, or -1 if there isn't exactly one.
	parentBulletLine := func(b *block) int {
		var cands []int
		if b.parent == "" {
			cands = findBullets(b.name, 0, len(lines), true)
		} else if ph := resolved[b.parent]; ph != nil {
			cands = findBullets(b.name, ph.line+1, ph.end, false)
		}
		if len(cands) == 1 {
			return cands[0]
		}
		return -1
	}
	// matches reports whether every bullet under h is a field of b.
	matches := func(h *heading, b *block) bool {
		if b.name != h.name || len(h.bullets) == 0 {
			return false
		}
		for _, c := range h.bullets {
			if !b.children[strings.SplitN(c, ".", 2)[0]] {
				return false
			}
		}
		return true
	}
	// Resolve outside in: a nested block's heading can only be resolved once
	// its parent's heading has been.
	for changed := true; changed; {
		changed = false
		for _, h := range headings {
			if h.path != "" {
				continue
			}
			var cands []*block
			for _, b := range blocks {
				if resolved[b.path] == nil && matches(h, b) && parentResolved(b) {
					cands = append(cands, b)
				}
			}
			// Prefer the candidate whose parent bullet is the closest one before the heading.
			var best *block
			bestLine, tie := -1, false
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
			h.path = best.path
			resolved[best.path] = h
			changed = true
		}
	}

	// Blocks with the same schema shape under other parents that have no
	// heading of their own share the heading (the doc documents the shape once).
	for _, b := range blocks {
		if resolved[b.path] != nil || !parentResolved(b) || parentBulletLine(b) < 0 {
			continue
		}
		var share *heading
		for _, h := range headings {
			if h.path != "" && matches(h, b) && sameChildren(blocks[h.path], b) {
				if share != nil && share != h {
					share = nil
					break
				}
				share = h
			}
		}
		if share != nil {
			share.shared = append(share.shared, b.path)
			resolved[b.path] = share
		}
	}

	existingAnchors := map[string]bool{}
	for _, l := range lines {
		for _, m := range anchorNameRe.FindAllStringSubmatch(l, -1) {
			existingAnchors[m[1]] = true
		}
	}

	// Add anchors and links. Edits keep each line's bullet or heading prefix,
	// so the line classification above stays valid.
	for _, h := range headings {
		switch {
		case h.path == "":
			if h.anchor == "" {
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: couldn't match `%s` heading (bullets %v) to exactly one schema block; left alone", h.line+1, h.name, h.bullets))
			}
			continue
		case h.anchor != "" && !strings.HasPrefix(h.anchor, "nested_"):
			res.Notes = append(res.Notes, fmt.Sprintf("line %d: heading has non-nested anchor %q; left alone", h.line+1, h.anchor))
			continue
		case h.anchor == "":
			want := "nested_" + strings.ReplaceAll(h.path, ".", "_")
			if existingAnchors[want] {
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: anchor %q already exists; left `%s` heading alone", h.line+1, want, h.name))
				continue
			}
			lines[h.line] = fmt.Sprintf(`<a name="%s"></a>`, want) + lines[h.line]
			res.Changes = append(res.Changes, fmt.Sprintf("line %d: added anchor %s to `%s` heading", h.line+1, want, h.name))
			existingAnchors[want] = true
			h.anchor = want
		}
		for _, p := range append([]string{h.path}, h.shared...) {
			b := blocks[p]
			pl := parentBulletLine(b)
			if pl < 0 {
				res.Notes = append(res.Notes, fmt.Sprintf("`%s`: couldn't find exactly one parent bullet for `%s`; not linked", p, b.name))
				continue
			}
			msg, ok := linkBullet(lines, pl, h.anchor, existingAnchors, isBoundary)
			switch {
			case !ok:
				res.Notes = append(res.Notes, fmt.Sprintf("line %d: %s", pl+1, msg))
			case msg != "":
				res.Changes = append(res.Changes, fmt.Sprintf("line %d: %s", pl+1, msg))
			}
		}
	}

	if len(res.Changes) > 0 {
		res.Text = strings.Join(lines, "\n")
		if crlf {
			res.Text = strings.ReplaceAll(res.Text, "\n", "\r\n")
		}
	}
	return res
}

func sameChildren(a, b *block) bool {
	if len(a.children) != len(b.children) {
		return false
	}
	for c := range a.children {
		if !b.children[c] {
			return false
		}
	}
	return true
}

// linkBullet makes the bullet starting at line pl link to anchor, editing
// lines in place. It returns a description of the change ("" if the bullet
// was already linked) and false if it couldn't link the bullet.
func linkBullet(lines []string, pl int, anchor string, anchors map[string]bool, isBoundary func(int) bool) (string, bool) {
	// Extent of the bullet: until the next column-0 bullet, boundary, or code fence.
	end := pl + 1
	for end < len(lines) && !bulletRe.MatchString(lines[end]) && !isBoundary(end) && !strings.HasPrefix(lines[end], "```") {
		end++
	}
	bullet := strings.Join(lines[pl:end], "\n")
	// Any existing link to the anchor counts, whatever its text (for
	// example "See [structure below](#nested_x)"), even across lines.
	if strings.Contains(bullet, "(#"+anchor+")") {
		return "", true
	}
	for k := pl; k < end; k++ {
		if m := linkedBelowRe.FindStringSubmatch(lines[k]); m != nil {
			if anchors[m[1]] {
				return fmt.Sprintf("bullet already links to a different anchor %q; not changed", m[1]), false
			}
			lines[k] = strings.Replace(lines[k], "(#"+m[1]+")", "(#"+anchor+")", 1)
			return fmt.Sprintf("fixed broken link #%s -> #%s", m[1], anchor), true
		}
	}
	for k := pl; k < end; k++ {
		if loc := plainBelowRe.FindStringIndex(lines[k]); loc != nil {
			// Keep the matched "Structure is" wording; link "documented below".
			l := lines[k]
			mid := loc[1] - len("documented below")
			lines[k] = l[:mid] + "[documented below](#" + anchor + ")" + l[loc[1]:]
			return `linked existing "Structure is documented below"`, true
		}
	}
	// Append to the end of the bullet's first paragraph.
	last := pl
	for k := pl + 1; k < end; k++ {
		if strings.TrimSpace(lines[k]) == "" || indentedListRe.MatchString(lines[k]) {
			break
		}
		last = k
	}
	l := strings.TrimRight(lines[last], " \t")
	sep := ". "
	if strings.ContainsAny(l[len(l)-1:], ".:)-") {
		sep = " "
	}
	lines[last] = l + sep + "Structure is [documented below](#" + anchor + ")."
	return "appended link to " + anchor, true
}
