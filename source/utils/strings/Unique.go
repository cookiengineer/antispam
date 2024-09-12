package strings

import "sort"

func Unique(values []string) []string {

	var result []string

	hashmap := make(map[string]bool)

	for v := 0; v < len(values); v++ {

		value := values[v]

		_, ok := hashmap[value]

		if ok == false {
			result = append(result, value)
			hashmap[value] = true
		}

	}

	sort.Strings(result)

	return result

}
