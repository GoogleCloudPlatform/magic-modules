package detector

import (
	"fmt"
	"sort"
	"strings"
)

func formatTransitionCell(tr TransitionResult) string {
	switch tr.Status {
	case StatusCovered:
		if tr.CoveredBy != "" {
			if tr.UnsetMode != "" {
				return fmt.Sprintf("COVERED (%s, %s)", tr.CoveredBy, tr.UnsetMode)
			}
			return fmt.Sprintf("COVERED (%s)", tr.CoveredBy)
		}
		return "COVERED"
	case StatusExempt:
		if tr.ExemptCode != "" {
			return fmt.Sprintf("EXEMPT (%s)", tr.ExemptCode)
		}
		return "EXEMPT"
	default:
		return "MISSING"
	}
}

func buildSuggestedTestNote(
	resourceName string,
	fields map[string]*ExhaustiveFieldStatus,
	incompleteFields []string,
) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Resource %q has %d field(s) with incomplete exhaustive lifecycle test coverage:\n", resourceName, len(incompleteFields))
	for _, fp := range incompleteFields {
		st := fields[fp]
		var missing []string
		if st.Set.Status == StatusMissing {
			missing = append(missing, string(TransitionSet))
		}
		if st.Modify.Status == StatusMissing {
			missing = append(missing, string(TransitionModify))
		}
		if st.Unset.Status == StatusMissing {
			missing = append(missing, string(TransitionUnset))
		}
		fmt.Fprintf(&b, "- %s: missing %s\n", fp, strings.Join(missing, ", "))
	}
	return b.String()
}

// FormatExhaustiveMarkdown renders a human-readable Markdown table summarizing
// SET, MODIFY, and UNSET coverage and exemptions per resource and field path.
func FormatExhaustiveMarkdown(results map[string]*ExhaustiveMissingTestInfo) string {
	if len(results) == 0 {
		return "All target fields have exhaustive lifecycle test coverage.\n"
	}

	resourceNames := make([]string, 0, len(results))
	for r := range results {
		resourceNames = append(resourceNames, r)
	}
	sort.Strings(resourceNames)

	var b strings.Builder
	for i, resourceName := range resourceNames {
		if i > 0 {
			b.WriteString("\n")
		}
		info := results[resourceName]
		fmt.Fprintf(&b, "### `%s`\n\n", resourceName)
		b.WriteString("| Field Path | Kind | SET | MODIFY | UNSET |\n")
		b.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")

		for _, fp := range sortedFieldPaths(info.Fields) {
			st := info.Fields[fp]
			kind := "Leaf"
			if st.IsNestedBlock {
				kind = "NestedBlock"
			}
			fmt.Fprintf(
				&b,
				"| `%s` | %s | %s | %s | %s |\n",
				fp,
				kind,
				formatTransitionCell(st.Set),
				formatTransitionCell(st.Modify),
				formatTransitionCell(st.Unset),
			)
		}
	}
	return b.String()
}
