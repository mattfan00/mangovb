package store

import (
	"fmt"
	"strings"
)

func setMap(clauses map[string]interface{}) string {
	resultArr := make([]string, len(clauses))
	i := 0
	for key, value := range clauses {
		resultArr[i] = fmt.Sprintf("%s = %v", key, value)
		i++
	}

	return strings.Join(resultArr, ", ")
}
