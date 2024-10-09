package diff

import (
	"cmp"
)

// Return the union of two maps, overwriting any shared keys with the second map's values
func union[K comparable, V any](map1, map2 map[K]V) map[K]V {
	if len(map1) == 0 {
		return map2
	}
	if len(map2) == 0 {
		return map1
	}
	merged := make(map[K]V, len(map1)+len(map2))
	for k, v := range map1 {
		merged[k] = v
	}
	for k, v := range map2 {
		merged[k] = v
	}
	return merged
}

func sortedUnion[K cmp.Ordered](list1, list2 []K) []K {
	if len(list1) == 0 {
		return list2
	}
	if len(list2) == 0 {
		return list1
	}
	merged := make([]K, 0, len(list1)+len(list2))
	i, j := 0, 0
	for i < len(list1) && j < len(list2) {
		if list1[i] < list2[j] {
			merged = append(merged, list1[i])
			i++
		} else if list1[i] > list2[j] {
			merged = append(merged, list2[j])
			j++
		} else { // keys are equal
			merged = append(merged, list1[i])
			i++
			j++
		}
	}

	// Add any remaining elements from both slices
	merged = append(merged, list1[i:]...)
	merged = append(merged, list2[j:]...)
	return merged
}
