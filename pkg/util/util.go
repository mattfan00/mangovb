package util

import "strconv"

func SliceAtoi(sarr []string) ([]int, error) {
	iarr := make([]int, len(sarr))
	for i, s := range sarr {
		j, err := strconv.Atoi(s)
		if err != nil {
			return iarr, err
		}
		iarr[i] = j
	}

	return iarr, nil
}
