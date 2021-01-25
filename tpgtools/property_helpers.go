package main

// Sorts properties to be in a standard order
func propComparator(props []Property) func(i, j int) bool {
	return func(i, j int) bool {
		l := props[i]
		r := props[j]

		// required < non-required
		if l.Required && !r.Required {
			return true
		}

		// conversely, non-required > required
		if r.Required && !l.Required {
			return false
		}

		// same deal- settable (optional / O+C) fields > Computed fields
		if l.Settable && !r.Settable {
			return true
		}
		if r.Settable && !l.Settable {
			return false
		}

		// finally, sort by name
		return l.Name() < r.Name()
	}
}
