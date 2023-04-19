package support

import "sort"

func ksort(items map[string]string) map[string]string {
	keys := make([]string, 0, len(items))

	for key := range items {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	output := map[string]string{}

	for _, k := range keys {
		output[k] = items[k]
	}

	return output
}
