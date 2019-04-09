package anaml

import (
	"strconv"
)

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandIdentifierList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vv, _ := strconv.Atoi(v.(string))
			vs = append(vs, vv)
		}
	}
	return vs
}

func identifierList(ints []int) []string {
	vs := make([]string, 0, len(ints))
	for _, v := range ints {
		vs = append(vs, strconv.Itoa(v))
	}
	return vs
}
